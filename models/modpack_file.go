package models

type ModPackFile struct {
	Name        string               `json:"name"`
	Path        string               `json:"path"`
	Size        int64                `json:"size"`
	Checksum    uint32               `json:"checksum"`
	Environment MinecraftEnvironment `json:"environment"`
	Type        MinecraftFileType    `json:"type"`
}
