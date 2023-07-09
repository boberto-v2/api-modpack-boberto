package event_cache

import (
	"time"

	event_cache_models "github.com/brutalzinn/boberto-modpack-api/services/event/cache/models"
	"github.com/patrickmn/go-cache"
)

var event_cache = cache.New(5*time.Minute, 10*time.Minute)

func GetById(id string) (eventCache event_cache_models.EventCache, found bool) {
	if eventCache, found := event_cache.Get(id); found {
		return eventCache.(event_cache_models.EventCache), true
	}
	return event_cache_models.EventCache{}, false
}

func Add(eventCache event_cache_models.EventCache) {
	event_cache.Add(eventCache.Id, eventCache, cache.DefaultExpiration)
}

func Remove(key string) {
	event_cache.Delete(key)
}
