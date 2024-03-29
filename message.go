package main

import (
	"fmt"
	"math/big"
	"math/rand"
	"net"
	"time"

	"github.com/jmoiron/sqlx"
)

type Message struct {
	Transaction              Transaction `db:"transaction" json:"transaction"`
	Switch                   AtmSwitch   `db:"switch" json:"switch"`
	PrimaryAccountNumber     string      `db:"primary_account_number" json:"primaryAccountNumber,omitempty"`
	TransactionAmount        float64     `db:"transaction_amount" json:"transactionAmount,omitempty"`
	AcquiringInstitutionCode string      `db:"acquiring_institution_code" json:"acquiringInstitutionCode"`
	ReceivingInstitutionCode string      `db:"receiving_institution_code" json:"receivingInstitutionCode,omitempty"`
	TransactionFee           float64     `db:"transaction_fee" json:"transactionFee,omitempty"`
	TerminalNameAndLocation  string      `db:"terminal_name_location" json:"terminalNameAndLocation"`
	CurrencyCode             Currency    `db:"currency_code" json:"currencyCode"`
	TerminalID               string      `db:"terminal_id" json:"terminalId"`
	SourceAccount            string      `db:"source_account" json:"sourceAccount,omitempty"`
	DestinationAccount       string      `db:"destination_account" json:"destinationAccount,omitempty"`
	Channel                  Channel     `db:"channel" json:"channel"`
	Device                   Device      `db:"device" json:"device"`
	TargetBank               Bank        `db:"target_bank" json:"targetBank,omitempty"`
	Id                       int         `db:"id" json:"id,omitempty"`
	Rrn                      string      `db:"rrn" json:"rrn,omitempty"`
	TraceNumber              string      `db:"trace_number" json:"traceNumber,omitempty"`
	TransmissionDateTime     string      `db:"transmission_date_time" json:"transmissionDateTime,omitempty"`
	LocalTransactionDateTime string      `db:"local_transaction_date_time" json:"localTransactionDateTime,omitempty"`
	OriginalDataElements     string      `db:"original_data_elements,omitempty" json:"originalDataElements,omitempty"`
	Mti                      string      `db:"mti" json:"mti,omitempty"`
	ProcessCode              string      `db:"process_code" json:"processCod,omitempty"`
}

type AtmResponse struct {
	TraceNumber  string `json:"traceNumber"`
	ResponseCode string `json:"responseCode"`
	Balance      string `json:"balance"`
	RRN          string `json:"rrn"`
}

type Currency string

const (
	PHP Currency = "608"
	USD Currency = "840"
)

type Device string

const (
	ATM Device = "6011"
	POS Device = "6012"
	NAD Device = "6016"
)

type AtmSwitch string

const (
	CORTEX     AtmSwitch = "CORTEX"
	POSTBRIDGE AtmSwitch = "POSTBRIDGE"
	COREWARE   AtmSwitch = "COREWARE"
	NARADA     AtmSwitch = "NARADA"
)

type Transaction string

const (
	WITHDRAW Transaction = "WITHDRAW"
	BAL_INQ  Transaction = "BAL_INQ"
	FT       Transaction = "FT"
	IBFTC    Transaction = "IBFTC"
	IBFTD    Transaction = "IBFTD"
	ELOAD    Transaction = "ELOAD"
	BILLS    Transaction = "BILLS"
	PURCHASE Transaction = "PURCHASE"
)

type Channel string

const (
	ON_US      Channel = "ON_US"
	OFF_US     Channel = "OFF_US"
	MASTERCARD Channel = "MASTERCARD"
)

type Bank string

const (
	OTHER_BANK   Bank = "OTHER_BANK"
	INTER_SYSTEM Bank = "INTER_SYSTEM"
)

type MTI string

const (
	FinancialRequestMasterVisa    MTI = "100"
	FinancialRequest              MTI = "200"
	FinancialAdvice               MTI = "220"
	FinancialReversal             MTI = "400"
	FinancialReversalAdvice       MTI = "420"
	FinancialReversalRepeatAdvice MTI = "421"
	NetworkManagementRequest      MTI = "800"
)

type messageService struct {
	db *sqlx.DB
}

func (s *messageService) getMessage(id int) (Message, error) {
	message := Message{}
	err := s.db.Get(&message, "SELECT * FROM atm_message WHERE id=$1", id)
	return message, err
}

func (s *messageService) getMessages(page int) ([]Message, error) {
	message := []Message{}
	err := s.db.Select(&message, "SELECT * FROM atm_message ORDER BY id DESC LIMIT 50 OFFSET $1", (page-1)*10)
	return message, err
}

func (s *messageService) saveMessage(message Message) error {
	_, err := s.db.NamedExec(`INSERT INTO atm_message (
		mti,
		"transaction", 
		primary_account_number, 
		transaction_amount, 
		acquiring_institution_code, 
		receiving_institution_code,
		terminal_name_location, 
		currency_code,
		terminal_id,
		source_account,
		destination_account,
		channel,
		device,
		target_bank,
		rrn,
		trace_number,
		transmission_date_time,
		local_transaction_date_time,
		original_data_elements,
		process_code,
		switch
	  ) VALUES (
		:mti, :transaction, :primary_account_number, :transaction_amount, :acquiring_institution_code, :receiving_institution_code, 
		:terminal_name_location, :currency_code, :terminal_id, :source_account, :destination_account, :channel, :device, 
		:target_bank, :rrn, :trace_number, :transmission_date_time, :local_transaction_date_time, :original_data_elements, :process_code, :switch
	  )`, message)
	return err
}

func (s *messageService) sendTcpMessage(conn net.Conn, packed []byte) error {
	conn.SetWriteDeadline(time.Now().Add(30 * time.Second))
	_, err := conn.Write(packed)
	if err != nil {
		return err
	}
	err = conn.SetReadDeadline(time.Now().Add(30 * time.Second))
	if err != nil {
		return err
	}
	return nil
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

func generateStan() string {
	min := 0
	max := 999999
	randomSixDigitNumber := rand.Intn(max-min+1) + min
	return fmt.Sprintf("%06d", randomSixDigitNumber)
}

func generateLocalTransactionDateTime(transmissionDateTime string) string {
	currentDate := time.Now()
	year := fmt.Sprintf("%02d", currentDate.Year()%100)
	return year + transmissionDateTime
}

func balanceDeserializer(input string) *big.Float {
	b := new(big.Float)
	if len(input) >= 20 {
		balance := input[8:20]
		b.SetString(balance)
		return b.Quo(b, big.NewFloat(100))
	}
	big.NewFloat(0)
	return b
}
