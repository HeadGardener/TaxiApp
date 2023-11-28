package models

import (
	"fmt"
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

var TaxiTypesStr = map[string]TaxiType{
	economyStr:  economy,
	comfortStr:  comfort,
	businessStr: business,
}

func (tt TaxiType) String() string {
	switch tt {
	case 0:
		return economyStr
	case 1:
		return comfortStr
	case 2:
		return businessStr
	default:
		return fmt.Sprintf("undefined")
	}
}

func (tt TaxiType) FromString(taxiType string) TaxiType {
	return TaxiTypesStr[taxiType]
}

type Trip struct {
	DriverID string    `db:"driver_id"`
	UserID   string    `db:"user_id"`
	From     string    `db:"from"`
	To       string    `db:"to"`
	Rating   float32   `db:"rating"`
	Date     time.Time `db:"date"`
}
