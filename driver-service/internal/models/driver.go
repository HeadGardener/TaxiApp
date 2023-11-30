package models

import (
	"fmt"
	"time"
)

const (
	ConsumedStatus = "CONSUMED"
)

type DriverStatus int

const (
	Busy DriverStatus = iota
	Free
	Disable
)

const (
	busyStr    = "busy"
	freeStr    = "free"
	disableStr = "disable"
)

var DriverStatuses = map[DriverStatus]string{
	Busy:    busyStr,
	Free:    freeStr,
	Disable: disableStr,
}

var DriverStatusesStr = map[string]DriverStatus{
	busyStr:    Busy,
	freeStr:    Free,
	disableStr: Disable,
}

func (ds DriverStatus) String() string {
	switch ds {
	case 0:
		return busyStr
	case 1:
		return freeStr
	case 2:
		return disableStr
	default:
		return fmt.Sprintf("undefined")
	}
}

func (ds DriverStatus) FromString(driverStatus string) DriverStatus {
	return DriverStatusesStr[driverStatus]
}

type Driver struct {
	ID           string       `db:"id"`
	Name         string       `db:"name"`
	Surname      string       `db:"surname"`
	Phone        string       `db:"phone"`
	Email        string       `db:"email"`
	TaxiType     TaxiType     `db:"taxi_type"`
	Balance      float64      `db:"balance"`
	Password     string       `db:"password_hash"`
	Rating       float32      `db:"rating"`
	DriverStatus DriverStatus `db:"driver_status"`
	Registration time.Time    `db:"registration"`
	IsActive     bool         `db:"is_active"`
}

type Credentials struct {
	Card  string
	Money float64
}

type Order struct {
	ID       string
	DriverID string
	UserID   string
	From     string
	To       string
	Status   string
}
