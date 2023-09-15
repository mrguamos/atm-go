package main

import (
	"encoding/xml"
	"fmt"
	"io"

	"github.com/moov-io/iso8583/network"
	"github.com/rs/zerolog/log"
)

type postbridgeSwitch struct {
}

var postbridge = &postbridgeSwitch{}

func (s *postbridgeSwitch) build(message *Message, reversal bool) {
	originalMti := message.Mti
	message.Mti = s.getMti(*message, reversal)
	if reversal {
		originalDataElements := s.serializeOriginalDataElements(originalMti, message.TraceNumber, message.TransmissionDateTime, padLeftWithZeros(message.AcquiringInstitutionCode, 10))
		message.OriginalDataElements = originalDataElements
		message.Transaction = "REVERSAL " + message.Transaction
		return
	}
	message.TransmissionDateTime = generateTransmissionDateTime()
	message.TraceNumber = generateStan()
	message.Rrn = generateRrn()
	message.LocalTransactionDateTime = generateLocalTransactionDateTime(message.TransmissionDateTime)
	message.ProcessCode = s.getProcessCode(*message) + "0000"
}

func (s *postbridgeSwitch) pack(message Message) ([]byte, error) {
	t := message.TransmissionDateTime
	l := message.LocalTransactionDateTime
	iso := Iso8583PostXml{
		MsgType: message.Mti,
		Fields: &Fields{
			Field002: message.PrimaryAccountNumber,
			Field003: message.ProcessCode,
			Field004: padLeftWithZeros(moveDecimalRight(message.TransactionAmount), 12),
			Field007: t,
			Field011: message.TraceNumber,
			Field012: t[4:],
			Field013: t[:4],
			Field014: fmt.Sprintf("%02s%02s", l[:2], l[2:4]),
			Field015: t[0:4],
			Field018: string(message.Device),
			Field028: "D" + padLeftWithZeros(moveDecimalRight(message.TransactionFee), 8),
			Field032: padLeftWithZeros(message.AcquiringInstitutionCode, 10),
			Field037: message.Rrn,
			Field041: message.TerminalID,
			Field043: message.TerminalNameAndLocation,
			Field049: string(message.CurrencyCode),
			Field090: message.OriginalDataElements,
			Field100: message.ReceivingInstitutionCode,
			Field102: message.SourceAccount,
			Field103: message.DestinationAccount,
			Field127025: &IccData{
				IccRequest: &IccRequestType{
					AmountAuthorized: padLeftWithZeros(moveDecimalRight(message.TransactionAmount), 12),
				},
			},
		},
	}
	xmlData, err := xml.MarshalIndent(iso, "", "    ")
	if err != nil {
		return nil, err
	}
	withHeader := xml.Header + string(xmlData)
	log.Print(withHeader)
	originalLength := len(withHeader)
	lengthPrefix := make([]byte, 2)
	lengthValue := int16(originalLength)
	lengthPrefix[0] = byte(lengthValue >> 8)
	lengthPrefix[1] = byte(lengthValue)
	withPrefix := append(lengthPrefix, withHeader...)
	return withPrefix, nil
}

func (s *postbridgeSwitch) unpack(r io.Reader) (AtmResponse, error) {
	lengthBuffer := network.NewBinary2BytesHeader()
	lengthBuffer.ReadFrom(r)
	length := lengthBuffer.Len
	response := make([]byte, length)
	n, err := io.ReadFull(r, response)
	log.Printf("%v", string(response))
	var iso8583PostXml Iso8583PostXml
	if err != nil {
		return AtmResponse{}, err
	}
	if n < 2 {
		return AtmResponse{}, err
	}
	err = xml.Unmarshal(response, &iso8583PostXml)
	if err != nil {
		return AtmResponse{}, err
	}
	balance := balanceDeserializer(iso8583PostXml.Fields.Field054)
	return AtmResponse{
		Balance:      fmt.Sprintf("%.2f", balance),
		TraceNumber:  iso8583PostXml.Fields.Field011,
		ResponseCode: iso8583PostXml.Fields.Field039,
		RRN:          iso8583PostXml.Fields.Field037,
	}, nil
}

func (s *postbridgeSwitch) packEchoTest() ([]byte, error) {
	t := generateTransmissionDateTime()
	message := Iso8583PostXml{
		MsgType: "0800",
		Fields: &Fields{
			Field007: t,
			Field011: generateStan(),
			Field012: t[4:],
			Field013: t[0:4],
			Field070: "301",
		},
	}
	xmlData, err := xml.MarshalIndent(message, "", "    ")
	if err != nil {
		return nil, err
	}
	withHeader := xml.Header + string(xmlData)
	originalLength := len(withHeader)
	lengthPrefix := make([]byte, 2)
	lengthValue := int16(originalLength)
	lengthPrefix[0] = byte(lengthValue >> 8)
	lengthPrefix[1] = byte(lengthValue)
	withPrefix := append(lengthPrefix, withHeader...)
	return withPrefix, nil
}

func (s *postbridgeSwitch) getProcessCode(message Message) string {
	var processCode string
	transaction := message.Transaction
	device := message.Device
	switch transaction {
	case PURCHASE:
		processCode = "00"
	case WITHDRAW, IBFTD, ELOAD:
		processCode = "01"
	case IBFTC:
		if message.TargetBank == INTER_SYSTEM {
			processCode = "26"
		} else {
			processCode = "21"
		}
	case BAL_INQ:
		if device == ATM {
			processCode = "30"
		} else {
			processCode = "31"
		}
	case FT:
		processCode = "40"
	case BILLS:
		processCode = "50"
	default:
		panic("Unable to get Process Code")
	}

	return processCode
}

func (s *postbridgeSwitch) getMti(message Message, reversal bool) string {
	if reversal {
		return fmt.Sprintf("0%d", FinancialReversal)
	}
	switch message.Channel {
	case MASTERCARD:
		return fmt.Sprintf("0%d", FinancialRequestMasterVisa)
	default:
		return fmt.Sprintf("0%d", FinancialRequest)
	}
}

func (s *postbridgeSwitch) serializeOriginalDataElements(mti string, traceNumber string, transmissionDateTime string, acquiringCode string) string {
	return fmt.Sprint(mti, traceNumber, transmissionDateTime, "01", acquiringCode)
}
