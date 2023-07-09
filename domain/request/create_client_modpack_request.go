package request

type CreateClientModPackRequest struct {
	Name             string `json:"name"`
	Author           string `json:"author"`
	Premium          bool   `json:"premium"`
	MinecraftVersion string `json:"minecraft_version"`
	ModPackVersion   string `json:"modpack_version"`
}
