package main

import (
	"fmt"
	"io"
	"math/big"
	"math/rand"
	"sort"
	"time"

	"github.com/moov-io/iso8583"
	"github.com/moov-io/iso8583/encoding"
	"github.com/moov-io/iso8583/field"
	"github.com/moov-io/iso8583/network"
	"github.com/rs/zerolog/log"
)

const header = "ISO8583-1993001000000"

type cortexSwitch struct {
	spec iso8583.MessageSpec
}

var fisGlobalSpec = &iso8583.MessageSpec{
	Fields: map[int]field.Field{
		0: field.NewNumeric(&field.Spec{
			Length:      4,
			Description: "Message Type Indicator",
			Enc:         encoding.BCD,
			Pref:        BCDPrefixer.Fixed,
		}),
		1: field.NewBitmap(&field.Spec{
			Description: "Bitmap",
			Enc:         encoding.Binary,
			Pref:        BCDPrefixer.Fixed,
		}),
		2: field.NewString(&field.Spec{
			Length:      19,
			Description: "Primary Account Number",
			Enc:         encoding.BCD,
			Pref:        BCDPrefixer.LL,
		}),
		3: field.NewString(&field.Spec{
			Length:      6,
			Description: "Processing Code",
			Enc:         encoding.BCD,
			Pref:        BCDPrefixer.Fixed,
		}),
		4: field.NewString(&field.Spec{
			Length:      12,
			Description: "Transaction Amount",
			Enc:         encoding.BCD,
			Pref:        BCDPrefixer.Fixed,
		}),
		6: field.NewString(&field.Spec{
			Length:      12,
			Description: "Card Holder Billing Amount",
			Enc:         encoding.BCD,
			Pref:        BCDPrefixer.Fixed,
		}),
		7: field.NewString(&field.Spec{
			Length:      10,
			Description: "Transmission Date and Time",
			Enc:         encoding.BCD,
			Pref:        BCDPrefixer.Fixed,
		}),
		11: field.NewString(&field.Spec{
			Length:      6,
			Description: "Trace Number",
			Enc:         encoding.BCD,
			Pref:        BCDPrefixer.Fixed,
		}),
		12: field.NewString(&field.Spec{
			Length:      12,
			Description: "Local Transaction Date and Time",
			Enc:         encoding.BCD,
			Pref:        BCDPrefixer.Fixed,
		}),
		26: field.NewString(&field.Spec{
			Length:      4,
			Description: "Merchant Code",
			Enc:         encoding.BCD,
			Pref:        BCDPrefixer.Fixed,
		}),
		30: field.NewString(&field.Spec{
			Length:      12,
			Description: "Original Amount",
			Enc:         encoding.BCD,
			Pref:        BCDPrefixer.Fixed,
		}),
		32: field.NewString(&field.Spec{
			Length:      11,
			Description: "Acquirer Code",
			Enc:         encoding.BCD,
			Pref:        BCDPrefixer.LL,
		}),
		33: field.NewString(&field.Spec{
			Length:      11,
			Description: "Forwarding Code",
			Enc:         encoding.BCD,
			Pref:        BCDPrefixer.LL,
		}),
		37: field.NewString(&field.Spec{
			Length:      12,
			Description: "Retrieval Reference Number",
			Enc:         encoding.ASCII,
			Pref:        BCDPrefixer.Fixed,
		}),
		39: field.NewString(&field.Spec{
			Length:      3,
			Description: "Response Code",
			Enc:         encoding.BCD,
			Pref:        BCDPrefixer.Fixed,
		}),
		41: field.NewString(&field.Spec{
			Length:      8,
			Description: "Terminal ID",
			Enc:         encoding.ASCII,
			Pref:        BCDPrefixer.Fixed,
		}),
		43: field.NewString(&field.Spec{
			Length:      99,
			Description: "Terminal Name and Location",
			Enc:         encoding.ASCII,
			Pref:        BCDPrefixer.LL,
		}),
		46: field.NewString(&field.Spec{
			Length:      204,
			Description: "Transaction Fee",
			Enc:         encoding.ASCII,
			Pref:        BCDPrefixer.LLL,
		}),
		49: field.NewString(&field.Spec{
			Length:      3,
			Description: "Transaction Currency Code",
			Enc:         encoding.BCD,
			Pref:        BCDPrefixer.Fixed,
		}),
		51: field.NewString(&field.Spec{
			Length:      3,
			Description: "Card Holder Currency Code",
			Enc:         encoding.BCD,
			Pref:        BCDPrefixer.Fixed,
		}),
		54: field.NewString(&field.Spec{
			Length:      120,
			Description: "Account Balance",
			Enc:         encoding.ASCII,
			Pref:        BCDPrefixer.LLL,
		}),
		56: field.NewString(&field.Spec{
			Length:      35,
			Description: "Original Data Elements",
			Enc:         encoding.BCD,
			Pref:        BCDPrefixer.LL,
		}),
		100: field.NewString(&field.Spec{
			Length:      11,
			Description: "Receiving Code",
			Enc:         encoding.BCD,
			Pref:        BCDPrefixer.LL,
		}),
		102: field.NewString(&field.Spec{
			Length:      28,
			Description: "Source Account",
			Enc:         encoding.ASCII,
			Pref:        BCDPrefixer.LL,
		}),
		103: field.NewString(&field.Spec{
			Length:      28,
			Description: "Destination Account",
			Enc:         encoding.ASCII,
			Pref:        BCDPrefixer.LL,
		}),
	},
}

func getMti(message Message, reversal bool) string {
	if reversal {
		return fmt.Sprintf("1%d", FinancialReversal)
	}
	switch message.Channel {
	case MASTERCARD:
		return fmt.Sprintf("1%d", FinancialRequestMasterVisa)
	default:
		return fmt.Sprintf("1%d", FinancialRequest)
	}
}

func getProcessCode(message Message) string {
	var processCode string
	transaction := message.Transaction
	device := message.Device
	switch transaction {
	case PURCHASE, ELOAD:
		processCode = "00"
	case WITHDRAW:
		processCode = "01"
	case IBFTC:
		processCode = "10"
	case IBFTD:
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

var cortex = &cortexSwitch{
	spec: *fisGlobalSpec,
}

func (s *cortexSwitch) build(message *Message, reversal bool) {
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

func (s *cortexSwitch) pack(message Message) ([]byte, error) {
	isoMesage := iso8583.NewMessage(fisGlobalSpec)

	isoMesage.MTI(message.Mti)

	if len(message.PrimaryAccountNumber) > 0 {
		isoMesage.Field(2, message.PrimaryAccountNumber)
	}

	isoMesage.Field(3, message.ProcessCode)
	isoMesage.Field(4, padLeftWithZeros(moveDecimalRight(message.TransactionAmount), fisGlobalSpec.Fields[4].Spec().Length))
	isoMesage.Field(6, padLeftWithZeros(moveDecimalRight(message.TransactionAmount), fisGlobalSpec.Fields[6].Spec().Length))
	isoMesage.Field(7, message.TransmissionDateTime)
	isoMesage.Field(11, message.TraceNumber)
	isoMesage.Field(12, message.LocalTransactionDateTime)
	isoMesage.Field(26, string(message.Device))
	isoMesage.Field(30, padLeftWithZeros(moveDecimalRight(message.TransactionAmount), fisGlobalSpec.Fields[4].Spec().Length))
	isoMesage.Field(32, padLeftWithZeros(message.AcquiringInstitutionCode, 10))
	isoMesage.Field(37, message.Rrn)

	if len(message.TerminalID) > 0 {
		isoMesage.Field(41, message.TerminalID)
	}

	if len(message.TerminalNameAndLocation) > 0 {
		isoMesage.Field(43, string(message.TerminalNameAndLocation))
	}

	isoMesage.Field(46, serializeFee(message.TransactionFee))

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

	headerBytes := []byte(header)
	result := append(headerBytes, rawMessage...)
	originalLength := len(result)
	lengthPrefix := make([]byte, 2)
	lengthValue := int16(originalLength)
	lengthPrefix[0] = byte(lengthValue >> 8)
	lengthPrefix[1] = byte(lengthValue)
	withPrefix := append(lengthPrefix, result...)
	return withPrefix, err
}

func (s *cortexSwitch) unpack(r io.Reader) (AtmResponse, error) {
	lengthBuffer := network.NewBinary2BytesHeader()
	lengthBuffer.ReadFrom(r)
	length := lengthBuffer.Len
	response := make([]byte, length)
	n, err := io.ReadFull(r, response)
	if err != nil {
		return AtmResponse{}, err
	}
	if n < 21 {
		return AtmResponse{}, err
	}
	header := make([]byte, 21)
	copy(header, response[:21])
	responseMessage := iso8583.NewMessage(fisGlobalSpec)
	unpacked := make([]byte, length)
	copy(unpacked, response[21:])
	responseMessage.Unpack(unpacked)
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

func moveDecimalRight(amount float64) string {
	return fmt.Sprint(int(amount * 100))
}

func padLeftWithZeros(input string, length int) string {
	return fmt.Sprintf("%0*s", length, input)
}

func generateTransmissionDateTime() string {
	currentDate := time.Now()
	month := fmt.Sprintf("%02d", int(currentDate.Month()))
	day := fmt.Sprintf("%02d", currentDate.Day())
	hours := fmt.Sprintf("%02d", currentDate.Hour())
	minutes := fmt.Sprintf("%02d", currentDate.Minute())
	seconds := fmt.Sprintf("%02d", currentDate.Second())
	return month + day + hours + minutes + seconds
}

func generateLocalTransactionDateTime(transmissionDateTime string) string {
	currentDate := time.Now()
	year := fmt.Sprintf("%02d", currentDate.Year()%100)
	return year + transmissionDateTime
}

func generateStan() string {
	min := 0
	max := 999999
	randomSixDigitNumber := rand.Intn(max-min+1) + min
	return fmt.Sprintf("%06d", randomSixDigitNumber)
}

func generateRrn() string {
	min := 0
	max := 999999999999
	randomTwelveDigitNumber := rand.Intn(max-min+1) + min
	return fmt.Sprintf("%012d", randomTwelveDigitNumber)
}

func serializeFee(fee float64) string {
	serializedFee := fmt.Sprint("00608D-", padLeftWithZeros(moveDecimalRight(fee), 7))
	return addTrailingSpaces(serializedFee, 34)
}

func serializeOriginalDataElements(mti string, traceNumber string, localTransactionDateTime string, acquiringCode string) string {
	return fmt.Sprint(mti, traceNumber, localTransactionDateTime, "01", acquiringCode)
}

func addTrailingSpaces(input string, length int) string {
	return fmt.Sprintf("%-*s", length, input)
}

func balanceDeserializer(input string) *big.Float {
	b := new(big.Float)
	if len(input) >= 20 {
		balance := input[7:18]
		b.SetString(balance)
		return b.Quo(b, big.NewFloat(100))
	}
	big.NewFloat(0)
	return b
}
