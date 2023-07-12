package rest_object

import (
	"github.com/brutalzinn/boberto-modpack-api/common"
	rest "github.com/brutalzinn/go-easy-rest"
	"github.com/gin-gonic/gin"
)

//TODO: Get Daniel help to do some improviments here. Now we are parting with view of a oriented object language. This is realllllllllllly rude to do with Go.
//First step is mitigate this procedural steps to create object resources to a extern module called rest.
// at this time the rest module already created. But.. needs some improviments to do this more flexible organization.
//create rest object to represent a modpack

const (
	MODPACK_OBJECT           = "modpack_object"
	FILE_OBJECT              = "file_object"
	WAITING_OBJECT           = "waiting_object"
	EVENT_OBJECT             = "event_object"
	USER_CREDENTIAL_OBJECT   = "user_credential_object"
	APIKEY_CREDENTIAL_OBJECT = "apikey_credential_object"
)

type Rest interface {
	CreateObject()
}

type RestObject struct {
	globalUrl string
	Resource  rest.Resource
	Link      []rest.Link
	ctx       *gin.Context
}

func New(ctx *gin.Context) RestObject {
	restObject := RestObject{
		globalUrl: common.GetUrl(ctx),
		ctx:       ctx,
	}
	return restObject
}

func (restObject RestObject) create() RestObject {
	urlHost := common.GetUrl(restObject.ctx)
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
