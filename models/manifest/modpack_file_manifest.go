package models_manifest

import (
	"time"

	"github.com/brutalzinn/boberto-modpack-api/models"
)

type ModPackFileManifest struct {
	CreateAt time.Time            `json:"create_at"`
	Files    []models.ModPackFile `json:"files"`
}
