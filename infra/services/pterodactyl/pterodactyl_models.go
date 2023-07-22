package pterodactyl_service

type TokenResponse struct {
	Token string `json:"token"`
}

type UploadTokenResponse struct {
	Token string `json:"token"`
}

type CommandResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type StatsResponse struct {
	Object     string       `json:"object"`
	Attributes StatsAttribs `json:"attributes"`
}

type StatsAttribs struct {
	CurrentState string        `json:"current_state"`
	IsSuspended  bool          `json:"is_suspended"`
	Resources    ResourceUsage `json:"resources"`
}

type ResourceUsage struct {
	Memory    int64 `json:"memory_bytes"`
	CPU       int64 `json:"cpu_absolute"`
	Disk      int64 `json:"disk_bytes"`
	NetworkRx int64 `json:"network_rx_bytes"`
	NetworkTx int64 `json:"network_tx_bytes"`
}
