package rest_object

import (
	"fmt"

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
		Link: []rest.Link{
			{
				Rel:    "_self",
				Href:   fmt.Sprintf("/game/client/modpack/%s", modPackCache.Id),
				Method: "GET",
			},
			{
				Rel:    "delete",
				Href:   fmt.Sprintf("/game/client/modpack/%s", modPackCache.Id),
				Method: "DELETE",
			},
			{
				Rel:    "update",
				Href:   fmt.Sprintf("/game/client/modpack/%s", modPackCache.Id),
				Method: "PUT",
			},
		},
	}
	return restObject
}
