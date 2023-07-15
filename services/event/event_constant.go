package event_service

type EventStatus int

const (
	MODPACK_FEEDBACK_EVENT = "modpack_feedback_event"
	UPLOAD_FILE_EVENT      = "upload_file_event"
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
