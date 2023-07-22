package pterodactyl_service

type Signal int

const (
	STOP    Signal = 1
	START   Signal = 2
	RESTART Signal = 3
	KILL    Signal = 4
)

func (signal Signal) GetName() string {
	switch signal {
	case STOP:
		return "stop"
	case START:
		return "start"
	case RESTART:
		return "restart"
	case KILL:
		return "kill"
	default:
		return "other"
	}
}
