package request

type ApiKeyRegisterRequest struct {
	AppName string `json:"app_name"`
	Days    uint32 `json:"days"`
}
