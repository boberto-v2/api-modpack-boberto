package rest_object

import (
	rest "github.com/brutalzinn/go-easy-rest"
)

type EventObject struct {
	Id   string    `json:"id"`
	Link rest.Link `json:"link"`
}

func (restObject RestObject) CreateEventObject(event EventObject) RestObject {
	restObject.Resource = rest.Resource{
		Object:    EVENT_OBJECT,
		Attribute: event,
		Link:      restObject.Link,
	}
	restObject.Create()
	return restObject
}
