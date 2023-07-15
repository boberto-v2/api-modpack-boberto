package game_client_request

type FinishClientModPackRequest struct {
	Name      string `json:"name"`
	FileUrl   string `json:"file_url"`
	ClientFtp *Ftp   `json:"client_ftp"`
}

type Ftp struct {
	Address   string `json:"address"`
	User      string `json:"user"`
	Password  string `json:"password"`
	Directory string `json:"directory"`
}
