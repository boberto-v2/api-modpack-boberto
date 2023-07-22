package authentication_apikey

import (
	"time"

	_ "github.com/brutalzinn/boberto-modpack-api/src/src/src/repository/database"
	apikey_database "github.com/brutalzinn/boberto-modpack-api/src/src/src/repository/database/apikey"
	entities_apikey "github.com/brutalzinn/boberto-modpack-api/src/src/src/repository/database/apikey/entities"
	entities_user "github.com/brutalzinn/boberto-modpack-api/src/src/src/repository/database/user/entities"
)

const (
	Default = "default"
	Member  = "modpack_read, modpack_create, modpack_update, modpack_delete"
)

type UserApiKey struct {
	Id       string
	AppName  string
	ExpireAt time.Time
	User     entities_user.User
}

type ApiKeyResult struct {
	ID     string
	Key    string
	Result string
}

func (userApiKey *UserApiKey) New() *ApiKey {
	apiKey := &ApiKey{
		AppName:  userApiKey.AppName,
		ExpireAt: userApiKey.ExpireAt,
		User:     userApiKey.User,
	}
	return apiKey
}

func (userApiKey *UserApiKey) Generate() (*ApiKeyResult, error) {
	appNameNormalized := normalizeAppName(userApiKey.AppName)
	apiKey, err := generate(appNameNormalized)
	if err != nil {
		return nil, err
	}
	apiKeyEntity := entities_apikey.ApiKey{
		Key:      apiKey.Result,
		UserId:   userApiKey.User.ID,
		AppName:  appNameNormalized,
		ExpireAt: userApiKey.ExpireAt,
		Duration: int64(userApiKey.ExpireAt.Sub(time.Now()).Hours() / 24),
		Scopes:   Default,
		Enabled:  true,
	}
	id, err := apikey_database.Insert(apiKeyEntity)
	if err != nil {
		return nil, err
	}
	apiKey.ID = id
	return apiKey, nil
}

func (userApiKey *UserApiKey) Regenerate() (*ApiKeyResult, error) {
	appNameNormalized := normalizeAppName(userApiKey.AppName)
	apiKey, err := generate(userApiKey.AppName)
	if err != nil {
		return nil, err
	}
	currentTime := time.Now()
	duration := int64(userApiKey.ExpireAt.Sub(currentTime).Hours() / 24)
	apiKeyEntity := entities_apikey.ApiKey{
		ID:       userApiKey.Id,
		Key:      apiKey.Result,
		AppName:  appNameNormalized,
		Duration: duration,
		ExpireAt: userApiKey.ExpireAt,
		Scopes:   Default,
		Enabled:  true,
	}
	_, err = apikey_database.Update(apiKeyEntity)
	if err != nil {
		return nil, err
	}

	return apiKey, nil
}
