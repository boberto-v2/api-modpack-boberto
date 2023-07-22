package event_service

type EventStatus int

const (
	MODPACK_PROGRESS_EVENT = "modpack_progress_event"
	FILE_PROGRESS_EVENT    = "file_progress_event"
	ERROR_EVENT            = "error_event"
)

const (
	EVENT_STARTED   EventStatus = 1
	EVENT_PENDING   EventStatus = 2
	EVENT_COMPLETED EventStatus = 3
	EVENT_CANCELED  EventStatus = 4
	EVENT_ABORTED   EventStatus = 5
)

func (event EventStatus) Parse() string {
	switch event {
	case EVENT_STARTED:
		return "started"
	case EVENT_PENDING:
		return "pending"
	case EVENT_CANCELED:
		return "canceled"
	case EVENT_COMPLETED:
		return "completed"
	case EVENT_ABORTED:
		return "aborted"
	default:
		return "other"
	}
}
