package rest_object

import (
	rest "github.com/brutalzinn/go-easy-rest"
)

const (
	WAITING_CLIENT_MESSAGE       = "WAITING_CLIENT_FILE"
	WAITING_SERVER_MESSAGE       = "WAITING_SERVER_FILE"
	WAITING_INTEGRATION_PARTENER = "WAITING_INTEGRATION"
)

type WaitingObject struct {
	Message string `json:"message"`
}

func (restObject RestObject) CreateWaitingObject(data WaitingObject) RestObject {
	restObject.Resource = rest.Resource{
		Object:    WAITING_OBJECT,
		Attribute: data,
		Link:      restObject.Link,
	}
	restObject.create()
	return restObject
}
