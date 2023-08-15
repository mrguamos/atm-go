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
	message.Mti = getMti(*message, reversal)
	if reversal {
		originalDataElements := serializeOriginalDataElements(message.Mti, message.TraceNumber, message.LocalTransactionDateTime, padLeftWithZeros(message.AcquiringInstitutionCode, 10))
		message.OriginalDataElements = originalDataElements
		message.Transaction = "REVERSAL " + message.Transaction
		return
	}
	message.TransmissionDateTime = generateTransmissionDateTime()
	message.TraceNumber = generateStan()
	message.Rrn = generateRrn()
	message.LocalTransactionDateTime = generateLocalTransactionDateTime(message.TransmissionDateTime)
	message.ProcessCode = getProcessCode(*message) + "0000"
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
			Field102: message.SourceAccount,
			Field103: message.DestinationAccount,
		},
	}
	xmlData, err := xml.MarshalIndent(iso, "", "    ")
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
