package rest_object

import (
	"github.com/brutalzinn/boberto-modpack-api/common"
	rest "github.com/brutalzinn/go-easy-rest"
	"github.com/gin-gonic/gin"
)

const (
	MODPACK_OBJECT = "modpack_object"
	FILE_OBJECT    = "file_object"
)

type RestObject struct {
	globalUrl string
	Attribute any
	Link      []rest.Link
}

func New(ctx *gin.Context) RestObject {
	restObject := RestObject{
		globalUrl: common.GetUrl(ctx),
	}
	return restObject
}

func (restObject RestObject) CreateHyperMedia(ctx *gin.Context) RestObject {
	urlHost := common.GetUrl(ctx)
	for i, href := range restObject.Link {
		href = restObject.resolveHref(urlHost, href)
		restObject.Link[i] = href
	}
	return restObject
}

func (restObject RestObject) resolveHref(url string, restLink rest.Link) rest.Link {
	restLink.Href = url + restLink.Href
	return restLink
}
