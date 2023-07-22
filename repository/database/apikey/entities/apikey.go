package entities_apikey

import "time"

type ApiKey struct {
	ID       string
	Key      string
	UserId   string
	AppName  string
	Duration int64
	Scopes   string
	Enabled  bool
	ExpireAt time.Time
	CreateAt time.Time
	UpdateAt *time.Time
}
