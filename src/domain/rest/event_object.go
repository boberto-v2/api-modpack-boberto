package rest_object

import (
	"github.com/brutalzinn/boberto-modpack-api/src/src/src/common/hypermedia"
	event_service "github.com/brutalzinn/boberto-modpack-api/src/src/src/infra/services/event"
	rest "github.com/brutalzinn/go-easy-rest"
)

type EventObject struct {
	Id   string    `json:"id"`
	Link rest.Link `json:"link"`
}

func (restObject *RestObject) CreateEventObject(event event_service.Event) *RestObject {
	restObject.Resource = rest.Resource{
		Object: EVENT_OBJECT,
		Attribute: EventObject{
			Id: event.Id,
			Link: rest.Link{
				Rel:  "listen_event",
				Href: hypermedia.GetSocketUrl(restObject.ctx) + "/application/event?id=" + event.Id,
			},
		},
		Link: restObject.Link,
	}
	return restObject
}
