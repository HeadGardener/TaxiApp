package models

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
		Economy:  make(chan OrderInfo, 100),
		Comfort:  make(chan OrderInfo, 100),
		Business: make(chan OrderInfo, 100),
	}
}

type DriversQueues struct {
	Economy  chan string
	Comfort  chan string
	Business chan string
}

func NewDriversQueues() *DriversQueues {
	return &DriversQueues{
		Economy:  make(chan string, 100),
		Comfort:  make(chan string, 100),
		Business: make(chan string, 100),
	}
}
