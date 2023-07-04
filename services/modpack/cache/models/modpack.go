package modpack_cache_models

import (
	"github.com/brutalzinn/boberto-modpack-api/common"
	modpack_models "github.com/brutalzinn/boberto-modpack-api/services/modpack/models"
)

type ModPackCache struct {
	Id             string                       `json:"id"`
	Name           string                       `json:"name"`
	Status         modpack_models.ModPackStatus `json:"status"`
	ManifestUrl    string                       `json:"manifest_url"`
	NormalizedName string                       `json:"normalized_name"`
}

func (modpack ModPackCache) New() ModPackCache {
	modpack.Id = common.GenerateUUID()
	modpack.NormalizedName = common.NormalizeString(modpack.NormalizedName)
	modpack.Status = modpack_models.Created
	return modpack
}
