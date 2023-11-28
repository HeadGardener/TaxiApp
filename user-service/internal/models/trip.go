package models

import (
	"errors"
	"time"
)

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

var taxiTypesStr = map[string]TaxiType{
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
	return taxiTypesStr[taxiType]
}

type Trip struct {
	ID       string    `json:"id" db:"taxi_type"`
	TaxiType TaxiType  `json:"taxi_type" db:"taxi_type"`
	Driver   string    `json:"driver" db:"driver"` // placeholder
	From     string    `json:"from" db:"from"`
	To       string    `json:"to" db:"to"`
	Rating   float32   `json:"rating" db:"rating"` // placeholder
	Date     time.Time `json:"date" db:"date"`
}

func (t *Trip) Validate() error {
	if _, ok := TaxiTypes[t.TaxiType]; !ok {
		return errors.New("invalid taxi type: only economy, comfort and business are available")
	}

	return nil
}
