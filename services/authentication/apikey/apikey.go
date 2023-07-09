package authentication_apikey

import (
	"errors"
	"time"

	"github.com/brutalzinn/boberto-modpack-api/common"
	apikey_database "github.com/brutalzinn/boberto-modpack-api/database/apikey"
	entities_user "github.com/brutalzinn/boberto-modpack-api/database/user/entities"
)

type ApiKey struct {
	Key      string
	AppName  string
	Enabled  bool
	User     entities_user.User
	Duration int64
	ExpireAt time.Time
	CreateAt time.Time
	UpdateAt time.Time
}

func New(appName string, key string) ApiKey {
	apiKey := ApiKey{
		Key:     key,
		AppName: appName,
	}
	return apiKey
}

func GetApiKeyByHeaderValue(headerValue string) (*ApiKey, error) {
	appName, err := extractAppName(headerValue)
	if err != nil {
		return nil, err
	}
	result, err := apikey_database.GetByAppName(appName)
	if err != nil {
		return nil, err
	}
	apiKey := ApiKey{
		AppName:  result.AppName,
		Enabled:  result.Enabled,
		ExpireAt: result.ExpireAt,
		Key:      result.Key,
		User: entities_user.User{
			ID: result.UserId,
		},
	}
	isValid := apiKey.Validate(headerValue)
	if !isValid {
		return nil, errors.New("api key expired or invalid")
	}
	return &apiKey, nil
}

func (apiKey ApiKey) Validate(key string) bool {
	enabled := apiKey.Enabled
	isExpired := apiKey.IsKeyExpired()
	isValidHash := common.BcryptCheckHash(key, apiKey.Key)
	//why this breaks my heart?
	isValid := enabled && !isExpired && isValidHash
	return isValid
}

func (apiKey *ApiKey) AddExpire(expireAt time.Duration) *ApiKey {
	apiKey.ExpireAt = time.Now().Add(expireAt)
	return apiKey
}

func (apiKey *ApiKey) Revoke() {
	apiKey.Enabled = false
}

func (apiKey *ApiKey) IsKeyExpired() bool {
	currentDateTime := time.Now()
	result := currentDateTime.After(apiKey.ExpireAt)
	return result
}
