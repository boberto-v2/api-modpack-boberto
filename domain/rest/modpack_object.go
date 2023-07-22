package rest_object

import (
	"github.com/brutalzinn/boberto-modpack-api/common/hypermedia"
	modpack_cache_models "github.com/brutalzinn/boberto-modpack-api/infra/services/modpack/cache/models"
	rest "github.com/brutalzinn/go-easy-rest"
)

type ModPackObject struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Environment string `json:"environment"`
	Status      string `json:"status"`
}

func (restObject *RestObject) CreateModPackObject(modPackCache modpack_cache_models.ModPackCache) *RestObject {
	hyperlink := hypermedia.New(restObject.ctx)
	restObject.Resource = rest.Resource{
		Object: MODPACK_OBJECT,
		Attribute: ModPackObject{
			Id:          modPackCache.Id,
			Name:        modPackCache.Name,
			Environment: modPackCache.Environment,
			Status:      modPackCache.Status.GetModPackStatus(),
		},
		Link: hyperlink.GetCurrentHyperLink(),
	}
	return restObject
}
