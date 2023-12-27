package models

const (
	undefined = "undefined"
)

type AcceptOrderStatus int

const (
	Accepted = iota
	Rejected
)

const (
	ConsumedStatus = "CONSUMED"
	AcceptedStatus = "ACCEPTED"
	RejectedStatus = "REJECTED"
)

var AcceptOrderStatuses = map[AcceptOrderStatus]string{
	Accepted: AcceptedStatus,
	Rejected: RejectedStatus,
}

var AcceptOrderStatusesStr = map[string]AcceptOrderStatus{
	AcceptedStatus: Accepted,
	RejectedStatus: Rejected,
}

func (aos AcceptOrderStatus) String() string {
	switch aos {
	case Accepted:
		return AcceptedStatus
	case Rejected:
		return RejectedStatus
	default:
		return undefined
	}
}

func (aos AcceptOrderStatus) FromString(orderStatus string) AcceptOrderStatus {
	return AcceptOrderStatusesStr[orderStatus]
}

type CompleteOrderStatus int

const (
	Completed = iota
	Canceled
)

const (
	CompletedStatus = "COMPLETED"
	CanceledStatus  = "CANCELED"
)

var CompleteOrderStatuses = map[CompleteOrderStatus]string{
	Completed: CompletedStatus,
	Canceled:  CanceledStatus,
}

var CompleteOrderStatusesStr = map[string]CompleteOrderStatus{
	CompletedStatus: Completed,
	CanceledStatus:  Canceled,
}

func (cos CompleteOrderStatus) String() string {
	switch cos {
	case Completed:
		return CompletedStatus
	case Canceled:
		return CanceledStatus
	default:
		return undefined
	}
}

func (cos CompleteOrderStatus) FromString(orderStatus string) CompleteOrderStatus {
	return CompleteOrderStatusesStr[orderStatus]
}
