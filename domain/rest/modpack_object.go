package rest_object

import (
	modpack_cache_models "github.com/brutalzinn/boberto-modpack-api/services/modpack/cache/models"
	rest "github.com/brutalzinn/go-easy-rest"
)

type ModPackObject struct {
	Name        string `json:"name"`
	Environment string `json:"environment"`
	Status      string `json:"status"`
}

func (restObject RestObject) CreateModPackObject(data modpack_cache_models.ModPackCache) RestObject {
	restObject.Resource = rest.Resource{
		Object: MODPACK_OBJECT,
		Attribute: ModPackObject{
			Name:        data.Name,
			Environment: data.Environment,
			Status:      data.Status.GetModPackStatus(),
		},
		Link: restObject.Link,
	}
	restObject.create()
	return restObject
}
