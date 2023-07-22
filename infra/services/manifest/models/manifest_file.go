package manifest_models

import (
	"time"

	modpack_models "github.com/brutalzinn/boberto-modpack-api/infra/services/modpack/models"
)

type ManifestFiles struct {
	CreateAt time.Time      `json:"create_at"`
	UpdateAt time.Time      `json:"update_at"`
	Files    []ManifestFile `json:"files"`
}

type ManifestFile struct {
	Name        string                              `json:"name"`
	Path        string                              `json:"path"`
	Size        int64                               `json:"size"`
	Checksum    uint32                              `json:"checksum"`
	Url         string                              `json:"url"`
	Environment modpack_models.MinecraftEnvironment `json:"environment"`
	Type        modpack_models.MinecraftFileType    `json:"type"`
}
