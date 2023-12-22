package models

const (
	queuesStartSize = 100
)

type OrderInfo struct {
	UserID  string
	OrderID string
	From    string
	To      string
}

type UsersQueues struct {
	Economy  chan OrderInfo
	Comfort  chan OrderInfo
	Business chan OrderInfo
}

func NewUsersQueues() *UsersQueues {
	return &UsersQueues{
		Economy:  make(chan OrderInfo, queuesStartSize),
		Comfort:  make(chan OrderInfo, queuesStartSize),
		Business: make(chan OrderInfo, queuesStartSize),
	}
}

type DriversQueues struct {
	Economy  chan string
	Comfort  chan string
	Business chan string
}

func NewDriversQueues() *DriversQueues {
	return &DriversQueues{
		Economy:  make(chan string, queuesStartSize),
		Comfort:  make(chan string, queuesStartSize),
		Business: make(chan string, queuesStartSize),
	}
}
