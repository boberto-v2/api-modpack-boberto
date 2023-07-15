package game_server_request

type FinishServerModPackRequest struct {
	ServerFtp      *Ftp                    `json:"server_ftp"`
	PterodactylApi *PterodactylIntegration `json:"pterodactyl_api"`
}

type Ftp struct {
	Address   string `json:"address"`
	User      string `json:"user"`
	Password  string `json:"password"`
	Directory string `json:"directory"`
}

type PterodactylIntegration struct {
	BaseUrl    string `json:"base_url"`
	ServerId   string `json:"server_id"`
	ApiKey     string `json:"api_key"`
	StartupCMD string `json:"startup_cmd"`
}
