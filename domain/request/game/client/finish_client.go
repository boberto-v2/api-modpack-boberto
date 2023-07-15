package game_client_request

type FinishClientModPackRequest struct {
	Name           string                  `json:"name"`
	ServerFtp      *Ftp                    `json:"server_ftp"`
	ClientFtp      *Ftp                    `json:"client_ftp"`
	PterodactylApi *PterodactylIntegration `json:"pterodactyl_api"`
}

type Ftp struct {
	Address   string `json:"address"`
	User      string `json:"user"`
	Password  string `json:"password"`
	Directory string `json:"directory"`
}

type PterodactylIntegration struct {
	BaseUrl  string `json:"base_url"`
	ServerId string `json:"server_id"`
	ApiKey   string `json:"api_key"`
}
