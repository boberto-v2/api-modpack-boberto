package event_service

//thanks to https://gist.github.com/crosstyan/47e7d3fa1b9e4716c0d6c76760a4a70c
import (
	"time"

	"github.com/brutalzinn/boberto-modpack-api/common"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/patrickmn/go-cache"
)

var event_cache = cache.New(5*time.Minute, 10*time.Minute)
var sessionGroupMap = make(map[string]map[uuid.UUID]*websocket.Conn)

type Event struct {
	Id       string      `json:"id"`
	Name     string      `json:"name"`
	Status   EventStatus `json:"status"`
	CreateAt time.Time   `json:"create_at"`
	ExpireAt time.Time   `json:"expire_at"`
}

func Create(eventName string) Event {
	id := common.GenerateUUID()
	new_event := Event{
		Id:       id,
		Name:     eventName,
		CreateAt: time.Now(),
		Status:   EVENT_STARTED,
	}
	event_cache.Add(id, new_event, cache.DefaultExpiration)
	return new_event
}

func GetById(id string) (eventCache Event, found bool) {
	if eventCache, found := event_cache.Get(id); found {
		return eventCache.(Event), true
	}
	return Event{}, false
}

func Remove(key string) {
	event_cache.Delete(key)
}
