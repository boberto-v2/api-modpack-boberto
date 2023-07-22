package modpack_models

type MinecraftModPack struct {
	Id      string `json:"id"`
	Name    string `json:"name"`
	Premium bool   `json:"premium"`
	Address string `json:"address"`
	FileUrl string `json:"file_url"`
	Author  string `json:"author"`
	Version string `json:"version"`
}
