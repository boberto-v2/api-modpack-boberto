package rest_object

import (
	rest "github.com/brutalzinn/go-easy-rest"
)

type ApiKeyCredentialObject struct {
	Id     string      `json:"id"`
	Key    string      `json:"key"`
	Header string      `json:"header"`
	Scopes string      `json:"scopes"`
	Link   []rest.Link `json:"link"`
}

func (restObject *RestObject) CreateApiKeycredentialObject(data ApiKeyCredentialObject) *RestObject {
	restObject.Resource = rest.Resource{
		Object:    APIKEY_CREDENTIAL_OBJECT,
		Attribute: data,
		// Link:      hypermedia.GetCurrentHyperLink(restObject.ctx, data.Id),
	}
	return restObject
}
