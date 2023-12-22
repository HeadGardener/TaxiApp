package models

import (
	"errors"
	"time"
)

var ErrInvalidTaxiType = errors.New("invalid taxi type")

type TaxiType int

const (
	economy TaxiType = iota
	comfort
	business
)

const (
	economyStr  = "economy"
	comfortStr  = "comfort"
	businessStr = "business"
)

var TaxiTypes = map[TaxiType]string{
	economy:  economyStr,
	comfort:  comfortStr,
	business: businessStr,
}

var TaxiTypesStr = map[string]TaxiType{
	economyStr:  economy,
	comfortStr:  comfort,
	businessStr: business,
}

func (tt TaxiType) String() string {
	switch tt {
	case economy:
		return economyStr
	case comfort:
		return comfortStr
	case business:
		return businessStr
	default:
		return "undefined"
	}
}

func (tt TaxiType) FromString(taxiType string) TaxiType {
	return TaxiTypesStr[taxiType]
}

type Trip struct {
	ID     string    `db:"id"`
	UserID string    `db:"user_id"`
	From   string    `db:"from"`
	To     string    `db:"to"`
	Rating float32   `db:"rating"`
	Date   time.Time `db:"date"`
}
