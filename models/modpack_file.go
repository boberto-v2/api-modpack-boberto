package models

type ModPackFile struct {
	Name        string               `json:"name"`
	Path        string               `json:"path"`
	Checksum    uint32               `json:"checksum"`
	Environment MinecraftEnvironment `json:"environment"`
	Type        MinecraftFileType    `json:"type"`
}
