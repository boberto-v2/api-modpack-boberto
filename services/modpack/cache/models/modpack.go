package modpack_cache_models

import (
	"github.com/brutalzinn/boberto-modpack-api/common"
	modpack_models "github.com/brutalzinn/boberto-modpack-api/services/modpack/models"
)

type ModPackCache struct {
	Id             string                       `json:"id"`
	Name           string                       `json:"name"`
	Environment    string                       `json:"environment"`
	Status         modpack_models.ModPackStatus `json:"status"`
	ManifestUrl    string                       `json:"manifest_url"`
	NormalizedName string                       `json:"normalized_name"`
}

func New() ModPackCache {
	modpack := ModPackCache{}
	modpack.Id = common.GenerateUUID()
	modpack.Status = modpack_models.Created
	return modpack
}
