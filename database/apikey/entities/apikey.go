package entities

import "time"

type ApiKey struct {
	ID          string
	Key         string
	UserId      string
	Scopes      string
	Description string
	ExpireAt    time.Time
	CreateAt    time.Time
	UpdateAt    time.Time
}
