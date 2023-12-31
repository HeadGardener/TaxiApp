package models

var (
	AcceptStatus   = "ACCEPTED"
	ProcessStatus  = "PROCESSING"
	ConsumedStatus = "CONSUMED"
)

type Order struct {
	TaxiType TaxiType `json:"taxi_type"`
	From     string   `json:"from"`
	To       string   `json:"to"`
}

type UserOrder struct {
	UserID   string
	TaxiType TaxiType
	From     string
	To       string
	Status   string
}
