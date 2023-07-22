package modpack_models

type ModPackStatus uint

const (
	Created            ModPackStatus = 1 //modpack created
	PendingClientFiles ModPackStatus = 2 // modpack is pending for client file upload
	PendingServerFiles ModPackStatus = 3 // modpack is pending for server file upload
	PendingFileUpload  ModPackStatus = 4 // pending file upload ( UPLOAD BY FTP TO OTHER LOCATION AFTER SERVER PENDING)
	Waiting            ModPackStatus = 5 // waiting something.. ( pterodactyl integration enters here)
	Finish             ModPackStatus = 6 // mod pack is created
	Canceled           ModPackStatus = 7 // canceled and cache need be clear
	Error              ModPackStatus = 8 // something goes wrong, cache needs be clear and process needs be restarted
	Aborted            ModPackStatus = 9 // process aborted and cache need be clear
)

func (modPackStatus ModPackStatus) GetModPackStatus() string {
	switch modPackStatus {
	case Created:
		return "created"
	case PendingClientFiles:
		return "pending_client_files"
	case PendingServerFiles:
		return "pending_server_files"
	case PendingFileUpload:
		return "pending_file_upload"
	case Waiting:
		return "waiting"
	case Finish:
		return "finish"
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
	case "pending_client_files":
		return PendingClientFiles
	case "pending_server_files":
		return PendingServerFiles
	case "pending_file_upload":
		return PendingFileUpload
	case "waiting":
		return Waiting
	case "finish":
		return Finish
	case "canceled":
		return Canceled
	case "error":
		return Error
	default:
		return Aborted
	}
}
