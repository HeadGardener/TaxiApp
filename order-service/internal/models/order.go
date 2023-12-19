package models

import "time"

type TaxiType int

type OrderStatus int

const (
	Economy TaxiType = iota
	Comfort
	Business
)

const (
	Progress OrderStatus = iota
	Finished
	Creating
)

const (
	economyStr  = "economy"
	comfortStr  = "comfort"
	businessStr = "business"
)

const (
	progressStr = "progress"
	finishedStr = "finished"
	creatingStr = "creating"
)

var TaxiTypes = map[TaxiType]string{
	Economy:  economyStr,
	Comfort:  comfortStr,
	Business: businessStr,
}

var OrderStatuses = map[OrderStatus]string{
	Progress: progressStr,
	Finished: finishedStr,
	Creating: creatingStr,
}

var TaxiTypesStr = map[string]TaxiType{
	economyStr:  Economy,
	comfortStr:  Comfort,
	businessStr: Business,
}

var OrderStatusesStr = map[string]OrderStatus{
	progressStr: Progress,
	finishedStr: Finished,
	creatingStr: Creating,
}

func (tt TaxiType) String() string {
	switch tt {
	case Economy:
		return economyStr
	case Comfort:
		return comfortStr
	case Business:
		return businessStr
	default:
		return "undefined"
	}
}

func (tt TaxiType) FromString(taxiType string) TaxiType {
	return TaxiTypesStr[taxiType]
}

func (os OrderStatus) String() string {
	switch os {
	case Progress:
		return progressStr
	case Finished:
		return finishedStr
	case Creating:
		return creatingStr
	default:
		return "undefined"
	}
}

func (os OrderStatus) FromString(status string) OrderStatus {
	return OrderStatusesStr[status]
}

type Order struct {
	ID       string      `bson:"id"`
	UserID   string      `bson:"user_id"`
	DriverID string      `bson:"driver_id"`
	From     string      `bson:"from"`
	To       string      `bson:"to"`
	TaxiType TaxiType    `bson:"taxi_type"`
	FromDate time.Time   `bson:"from_date"`
	ToDate   time.Time   `bson:"to_date"`
	Status   OrderStatus `bson:"status"`
	Comment  string      `bson:"comment"`
}
