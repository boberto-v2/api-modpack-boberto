package game_client_request

import rest_object "github.com/brutalzinn/boberto-modpack-api/src/src/src/domain/rest"

type FinishClientModPackRequest struct {
	Name      string `json:"name"`
	FileUrl   string `json:"file_url"`
	ClientFtp *Ftp   `json:"client_ftp"`
	Event     rest_object.EventObject
}

type Ftp struct {
	Address   string `json:"address"`
	User      string `json:"user"`
	Password  string `json:"password"`
	Directory string `json:"directory"`
}
