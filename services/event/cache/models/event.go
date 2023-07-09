package event_cache_models

import "time"

type EventCache struct {
	Id       string    `json:"id"`
	Name     string    `json:"name"`
	CreateAt time.Time `json:"create_at"`
	ExpireAt time.Time `json:"expire_at"`
}
