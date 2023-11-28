package models

import (
	"errors"
	"time"
)

type TransactionType int

const (
	Refill TransactionType = iota
	Spent
)

const (
	refillStr = "refill"
	spentStr  = "spent"
)

var transactionTypes = map[TransactionType]string{
	Refill: refillStr,
	Spent:  spentStr,
}

var transactionTypesStr = map[string]TransactionType{
	refillStr: Refill,
	spentStr:  Spent,
}

func (tt TransactionType) String() string {
	switch tt {
	case 0:
		return refillStr
	case 1:
		return spentStr
	default:
		return "unknown"
	}
}

func (tt TransactionType) FromString(transactionType string) TransactionType {
	return transactionTypesStr[transactionType]
}

const (
	Create   = "create"
	Blocked  = "blocked"
	Canceled = "canceled"
	Success  = "success"
)

type Transaction struct {
	ID                int             `db:"id"`
	WalletID          string          `db:"wallet_id"`
	Money             float64         `db:"spent"`
	TransactionType   TransactionType `db:"transaction_type"`
	TransactionStatus string          `db:"transaction_status"`
	Date              time.Time       `db:"date"`
}

type Credentials struct {
	WalletID string
	Money    float64
}

func BuildTransaction(transactionType TransactionType, status string, money float64) (Transaction, error) {
	if _, ok := transactionTypes[transactionType]; !ok {
		return Transaction{}, errors.New("invalid transaction type: only refill and spent are available")
	}

	return Transaction{
		WalletID:          "undefined",
		Money:             money,
		TransactionType:   transactionType,
		TransactionStatus: status,
		Date:              time.Now(),
	}, nil
}
