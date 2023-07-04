package modpack_models

type ModPackFtp struct {
	ServerFtp Ftp `json:"server_ftp"`
	ClientFtp Ftp `json:"client_ftp"`
}

type Ftp struct {
	Address   string `json:"address"`
	User      string `json:"user"`
	Password  string `json:"password"`
	Directory string `json:"directory"`
}
