package rest_object

import (
	"github.com/brutalzinn/boberto-modpack-api/common"
	modpack_cache_models "github.com/brutalzinn/boberto-modpack-api/services/modpack/cache/models"
	rest "github.com/brutalzinn/go-easy-rest"
)

type ModPackObject struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Environment string `json:"environment"`
	Status      string `json:"status"`
}

func (restObject *RestObject) CreateModPackObject(modPackCache modpack_cache_models.ModPackCache) *RestObject {
	restObject.Resource = rest.Resource{
		Object: MODPACK_OBJECT,
		Attribute: ModPackObject{
			Id:          modPackCache.Id,
			Name:        modPackCache.Name,
			Environment: modPackCache.Environment,
			Status:      modPackCache.Status.GetModPackStatus(),
		},
		Link: common.GetCurrentHypermedia(restObject.ctx, modPackCache.Id),
	}
	return restObject
}
