package models

import "time"

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
	case Busy:
		return busyStr
	case Free:
		return freeStr
	case Disable:
		return disableStr
	default:
		return undefined
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
	DriverStatus DriverStatus `db:"status"`
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
