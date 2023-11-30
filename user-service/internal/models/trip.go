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
	ID       string    `db:"id"`
	TaxiType TaxiType  `db:"taxi_type"`
	DriverID string    `db:"driver_id"` // placeholder
	From     string    `db:"from"`
	To       string    `db:"to"`
	Rating   float32   `db:"rating"` // placeholder
	Date     time.Time `db:"date"`
}

func (t *Trip) Validate() error {
	if _, ok := TaxiTypes[t.TaxiType]; !ok {
		return errors.New("invalid taxi type: only economy, comfort and business are available")
	}

	return nil
}
