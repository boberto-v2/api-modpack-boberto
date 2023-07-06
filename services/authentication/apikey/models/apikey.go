package authentication_apikey

import (
	"time"
)

type ApiKey struct {
	Key      string
	UserId   string
	ExpireAt time.Time
	CreateAt time.Time
	UpdateAt time.Time
}

func (apiKey ApiKey) Verify() (bool, error) {
	return false, nil
}

func isKeyExpired(expireAt time.Time) bool {
	currentDateTime := time.Now()
	if expireAt.After(currentDateTime) {
		return false
	}
	return true
}
