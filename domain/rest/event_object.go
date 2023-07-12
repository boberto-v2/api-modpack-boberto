package rest_object

import (
	"github.com/brutalzinn/boberto-modpack-api/common"
	event_service "github.com/brutalzinn/boberto-modpack-api/services/event"
	rest "github.com/brutalzinn/go-easy-rest"
)

type EventObject struct {
	Id   string    `json:"id"`
	Link rest.Link `json:"link"`
}

func (restObject RestObject) CreateEventObject(event event_service.Event) RestObject {
	restObject.Resource = rest.Resource{
		Object: EVENT_OBJECT,
		Attribute: EventObject{
			Id: event.Id,
			Link: rest.Link{
				Rel:  "listen_event",
				Href: common.GetSocketUrl(restObject.ctx),
			},
		},
		Link: restObject.Link,
	}
	return restObject
}