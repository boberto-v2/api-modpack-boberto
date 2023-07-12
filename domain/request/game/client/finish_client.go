package game_client_request

type ModPackFtp struct {
	ServerFtp      *Ftp            `json:"server_ftp"`
	ClientFtp      *Ftp            `json:"client_ftp"`
	PteroDactylApi *PterodactylApi `json:"pterodactyl_api"`
}

type Ftp struct {
	Address   string `json:"address"`
	User      string `json:"user"`
	Password  string `json:"password"`
	Directory string `json:"directory"`
}

type PterodactylApi struct {
	BaseUrl  string `json:"base_url"`
	ServerId string `json:"server_id"`
	ApiKey   string `json:"api_key"`
}
