package main

import (
	"errors"
	"fmt"
	"io"
	"sort"

	"github.com/moov-io/iso8583"
	"github.com/moov-io/iso8583/encoding"
	"github.com/moov-io/iso8583/field"
	"github.com/moov-io/iso8583/network"
	"github.com/moov-io/iso8583/prefix"
	"github.com/rs/zerolog/log"
)

type naradaSwitch struct {
	spec iso8583.MessageSpec
}

var naradaSpec = &iso8583.MessageSpec{
	Fields: map[int]field.Field{
		0: field.NewString(&field.Spec{
			Length:      4,
			Description: "Message Type Indicator",
			Enc:         encoding.EBCDIC,
			Pref:        prefix.EBCDIC.Fixed,
		}),
		1: field.NewBitmap(&field.Spec{
			Description: "Bitmap",
			Enc:         encoding.Binary,
			Pref:        prefix.EBCDIC.Fixed,
		}),
		2: field.NewString(&field.Spec{
			Length:      99,
			Description: "Primary Account Number",
			Enc:         encoding.EBCDIC,
			Pref:        prefix.EBCDIC.LL,
		}),
		3: field.NewString(&field.Spec{
			Length:      6,
			Description: "Processing Code",
			Enc:         encoding.EBCDIC,
			Pref:        prefix.EBCDIC.Fixed,
		}),
		4: field.NewString(&field.Spec{
			Length:      12,
			Description: "Transaction Amount",
			Enc:         encoding.EBCDIC,
			Pref:        prefix.EBCDIC.Fixed,
		}),
		6: field.NewString(&field.Spec{
			Length:      12,
			Description: "Card Holder Billing Amount",
			Enc:         encoding.EBCDIC,
			Pref:        prefix.EBCDIC.Fixed,
		}),
		7: field.NewString(&field.Spec{
			Length:      10,
			Description: "Transmission Date and Time",
			Enc:         encoding.EBCDIC,
			Pref:        prefix.EBCDIC.Fixed,
		}),
		11: field.NewString(&field.Spec{
			Length:      6,
			Description: "Trace Number",
			Enc:         encoding.EBCDIC,
			Pref:        prefix.EBCDIC.Fixed,
		}),
		13: field.NewString(&field.Spec{
			Length:      4,
			Description: "Local Transaction Date and Time",
			Enc:         encoding.EBCDIC,
			Pref:        prefix.EBCDIC.Fixed,
		}),
		14: field.NewString(&field.Spec{
			Length:      4,
			Description: "Expiration Date",
			Enc:         encoding.EBCDIC,
			Pref:        prefix.EBCDIC.Fixed,
		}),
		15: field.NewString(&field.Spec{
			Length:      4,
			Description: "Settlement Date",
			Enc:         encoding.EBCDIC,
			Pref:        prefix.EBCDIC.Fixed,
		}),
		18: field.NewString(&field.Spec{
			Length:      4,
			Description: "Merchant Code",
			Enc:         encoding.EBCDIC,
			Pref:        prefix.EBCDIC.Fixed,
		}),
		19: field.NewString(&field.Spec{
			Length:      3,
			Description: "Acquiring Institution Country Code",
			Enc:         encoding.EBCDIC,
			Pref:        prefix.EBCDIC.Fixed,
		}),
		22: field.NewString(&field.Spec{
			Length:      3,
			Description: "Pos Entry Mode",
			Enc:         encoding.EBCDIC,
			Pref:        prefix.EBCDIC.Fixed,
		}),
		25: field.NewString(&field.Spec{
			Length:      3,
			Description: "Pos Condition Code",
			Enc:         encoding.EBCDIC,
			Pref:        prefix.EBCDIC.Fixed,
		}),
		28: field.NewString(&field.Spec{
			Length:      9,
			Description: "Transaction Fee Amount",
			Enc:         encoding.EBCDIC,
			Pref:        prefix.EBCDIC.Fixed,
		}),
		30: field.NewString(&field.Spec{
			Length:      12,
			Description: "Original Amount",
			Enc:         encoding.EBCDIC,
			Pref:        prefix.EBCDIC.Fixed,
		}),
		32: field.NewString(&field.Spec{
			Length:      99,
			Description: "Acquirer Code",
			Enc:         encoding.EBCDIC,
			Pref:        prefix.EBCDIC.LL,
		}),
		35: field.NewString(&field.Spec{
			Length:      99,
			Description: "Track Date 2",
			Enc:         encoding.EBCDIC,
			Pref:        prefix.EBCDIC.LL,
		}),
		37: field.NewString(&field.Spec{
			Length:      12,
			Description: "Retrieval Reference Number",
			Enc:         encoding.EBCDIC,
			Pref:        prefix.EBCDIC.Fixed,
		}),
		39: field.NewString(&field.Spec{
			Length:      2,
			Description: "Response Code",
			Enc:         encoding.EBCDIC,
			Pref:        prefix.EBCDIC.Fixed,
		}),
		40: field.NewString(&field.Spec{
			Length:      3,
			Description: "Service Restriction Code",
			Enc:         encoding.EBCDIC,
			Pref:        prefix.EBCDIC.Fixed,
		}),
		41: field.NewString(&field.Spec{
			Length:      8,
			Description: "Terminal ID",
			Enc:         encoding.EBCDIC,
			Pref:        prefix.EBCDIC.Fixed,
		}),
		42: field.NewString(&field.Spec{
			Length:      15,
			Description: "Terminal Code",
			Enc:         encoding.EBCDIC,
			Pref:        prefix.EBCDIC.Fixed,
		}),
		43: field.NewString(&field.Spec{
			Length:      40,
			Description: "Terminal Name and Location",
			Enc:         encoding.EBCDIC,
			Pref:        prefix.EBCDIC.Fixed,
		}),
		47: field.NewString(&field.Spec{
			Length:      999,
			Description: "Additional Data National",
			Enc:         encoding.EBCDIC,
			Pref:        prefix.EBCDIC.LLL,
		}),
		49: field.NewString(&field.Spec{
			Length:      3,
			Description: "Transaction Currency Code",
			Enc:         encoding.EBCDIC,
			Pref:        prefix.EBCDIC.Fixed,
		}),
		51: field.NewString(&field.Spec{
			Length:      3,
			Description: "Card Holder Currency Code",
			Enc:         encoding.EBCDIC,
			Pref:        prefix.EBCDIC.Fixed,
		}),
		52: field.NewString(&field.Spec{
			Length:      16,
			Description: "Card Holder Currency Code",
			Enc:         encoding.EBCDIC,
			Pref:        prefix.EBCDIC.Fixed,
		}),
		54: field.NewString(&field.Spec{
			Length:      999,
			Description: "Account Balance",
			Enc:         encoding.EBCDIC,
			Pref:        prefix.EBCDIC.LLL,
		}),
		56: field.NewString(&field.Spec{
			Length:      99,
			Description: "Original Data Elements",
			Enc:         encoding.EBCDIC,
			Pref:        prefix.EBCDIC.LL,
		}),
		63: field.NewString(&field.Spec{
			Length:      99,
			Description: "Citi Share Data",
			Enc:         encoding.EBCDIC,
			Pref:        prefix.EBCDIC.LL,
		}),
		94: field.NewString(&field.Spec{
			Length:      2,
			Description: "Service Indicator",
			Enc:         encoding.EBCDIC,
			Pref:        prefix.EBCDIC.Fixed,
		}),
		102: field.NewString(&field.Spec{
			Length:      99,
			Description: "Source Account",
			Enc:         encoding.EBCDIC,
			Pref:        prefix.EBCDIC.LL,
		}),
		103: field.NewString(&field.Spec{
			Length:      99,
			Description: "Destination Account",
			Enc:         encoding.EBCDIC,
			Pref:        prefix.EBCDIC.LL,
		}),
	},
}

func (s *naradaSwitch) getMti(message Message, reversal bool) string {
	if reversal {
		return fmt.Sprintf("0%s", FinancialReversal)
	}
	switch message.Channel {
	case MASTERCARD:
		return fmt.Sprintf("0%s", FinancialRequestMasterVisa)
	default:
		return fmt.Sprintf("0%s", FinancialRequest)
	}
}

func (s *naradaSwitch) getProcessCode(message Message) string {
	var processCode string
	transaction := message.Transaction
	device := message.Device
	switch transaction {
	case PURCHASE, ELOAD:
		processCode = "00"
	case WITHDRAW:
		processCode = "01"
	case IBFTD:
		processCode = "10"
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
		processCode = "48"
	case BILLS:
		processCode = "51"
	default:
		panic("Unable to get Process Code")
	}

	return processCode
}

var narada = &naradaSwitch{
	spec: *naradaSpec,
}

func (s *naradaSwitch) build(message *Message, reversal bool) {
	originalMti := message.Mti
	message.Mti = s.getMti(*message, reversal)
	if reversal {
		originalDataElements := s.serializeOriginalDataElements(originalMti, message.TraceNumber, message.LocalTransactionDateTime, padLeftWithZeros(message.AcquiringInstitutionCode, 10))
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

func (s *naradaSwitch) pack(message Message) ([]byte, error) {
	isoMesage := iso8583.NewMessage(naradaSpec)

	isoMesage.MTI(message.Mti)

	if len(message.PrimaryAccountNumber) > 0 {
		isoMesage.Field(2, message.PrimaryAccountNumber)
	}

	isoMesage.Field(3, message.ProcessCode)
	isoMesage.Field(4, padLeftWithZeros(moveDecimalRight(message.TransactionAmount), naradaSpec.Fields[4].Spec().Length))
	isoMesage.Field(6, padLeftWithZeros(moveDecimalRight(message.TransactionAmount), naradaSpec.Fields[6].Spec().Length))
	isoMesage.Field(7, message.TransmissionDateTime)
	isoMesage.Field(11, message.TraceNumber)
	isoMesage.Field(13, message.LocalTransactionDateTime[2:4]+message.LocalTransactionDateTime[:2])
	isoMesage.Field(18, string(message.Device))
	isoMesage.Field(32, padLeftWithZeros(message.AcquiringInstitutionCode, 10))
	isoMesage.Field(37, message.Rrn)

	if len(message.TerminalID) > 0 {
		isoMesage.Field(41, message.TerminalID)
	}

	if len(message.TerminalNameAndLocation) > 0 {
		isoMesage.Field(43, string(message.TerminalNameAndLocation))
	}

	isoMesage.Field(46, padLeftWithZeros(moveDecimalRight(message.TransactionFee), 10))

	if len(message.CurrencyCode) > 0 {
		isoMesage.Field(49, string(message.CurrencyCode))
	}

	isoMesage.Field(51, string(PHP))

	if len(message.OriginalDataElements) > 0 {
		isoMesage.Field(56, message.OriginalDataElements)
	}

	if len(message.ReceivingInstitutionCode) > 0 {
		isoMesage.Field(100, message.ReceivingInstitutionCode)
	}

	if len(message.SourceAccount) > 0 {
		isoMesage.Field(102, message.SourceAccount)
	}

	if len(message.DestinationAccount) > 0 {
		isoMesage.Field(103, message.DestinationAccount)
	}

	keys := make([]int, 0, len(isoMesage.GetFields()))

	for k := range isoMesage.GetFields() {
		keys = append(keys, k)
	}

	sort.Ints(keys)

	for _, k := range keys {
		val, _ := isoMesage.GetString(k)
		log.Printf("Field %d: %s", k, val)
	}

	rawMessage, err := isoMesage.Pack()
	if err != nil {
		return nil, err
	}

	originalLength := len(rawMessage)
	lengthPrefix := make([]byte, 2)
	lengthValue := int16(originalLength)
	lengthPrefix[0] = byte(lengthValue >> 8)
	lengthPrefix[1] = byte(lengthValue)
	withPrefix := append(lengthPrefix, rawMessage...)
	return withPrefix, err
}

func (s *naradaSwitch) unpack(r io.Reader) (AtmResponse, error) {
	lengthBuffer := network.NewBinary2BytesHeader()
	lengthBuffer.ReadFrom(r)
	length := lengthBuffer.Len
	response := make([]byte, length)
	n, err := io.ReadFull(r, response)
	if err != nil {
		return AtmResponse{}, err
	}
	if n < 2 {
		return AtmResponse{}, err
	}
	responseMessage := iso8583.NewMessage(naradaSpec)
	responseMessage.Unpack(response)
	traceNumber, _ := responseMessage.GetField(11).String()
	responseCode, _ := responseMessage.GetField(39).String()
	rrn, _ := responseMessage.GetField(37).String()
	balanceField, _ := responseMessage.GetField(54).String()
	balance := balanceDeserializer(balanceField)
	atmResponse := AtmResponse{
		TraceNumber:  traceNumber,
		ResponseCode: responseCode,
		RRN:          rrn,
		Balance:      fmt.Sprintf("%.2f", balance),
	}
	keys := make([]int, 0, len(responseMessage.GetFields()))

	for k := range responseMessage.GetFields() {
		keys = append(keys, k)
	}

	sort.Ints(keys)

	for _, k := range keys {
		val, _ := responseMessage.GetString(k)
		log.Printf("Field %d: %s", k, val)
	}
	return atmResponse, nil
}

func (s *naradaSwitch) serializeOriginalDataElements(mti string, traceNumber string, localTransactionDateTime string, acquiringCode string) string {
	return fmt.Sprint(mti, traceNumber, localTransactionDateTime, "01", acquiringCode)
}

func (s *naradaSwitch) packEchoTest() ([]byte, error) {
	return nil, errors.New("echo test not supported")
}
