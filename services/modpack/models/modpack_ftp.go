package modpack_models

type ModPackFtp struct {
	ServerFtp Ftp
	ClientFtp Ftp
}

type Ftp struct {
	Address   string
	User      string
	Password  string
	Directory string
}
