package models

type ModPackStatus int

const (
	Created  ModPackStatus = 1
	Pending  ModPackStatus = 2
	Canceled ModPackStatus = 3
	Error    ModPackStatus = 4
	Aborted  ModPackStatus = 5
)

func (modPackStatus ModPackStatus) GetModPackStatus() string {
	switch modPackStatus {
	case Created:
		return "created"
	case Pending:
		return "pending"
	case Canceled:
		return "canceled"
	case Error:
		return "error"
	default:
		return "aborted"
	}
}

func ParseModPackStatus(modPackStatus string) ModPackStatus {
	switch modPackStatus {
	case "created":
		return Created
	case "pending":
		return Pending
	case "canceled":
		return Canceled
	case "error":
		return Error
	default:
		return Aborted
	}
}
