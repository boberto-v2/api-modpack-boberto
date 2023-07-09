package authentication_apikey

import (
	"time"

	_ "github.com/brutalzinn/boberto-modpack-api/database"
	apikey_database "github.com/brutalzinn/boberto-modpack-api/database/apikey"
	entities_apikey "github.com/brutalzinn/boberto-modpack-api/database/apikey/entities"
	entities_user "github.com/brutalzinn/boberto-modpack-api/database/user/entities"
)

type UserApiKey struct {
	AppName  string
	Duration time.Duration
	User     entities_user.User
}

type ApiKeyResult struct {
	Key    string
	Result string
}

func (userApiKey *UserApiKey) New() *ApiKey {
	apiKey := &ApiKey{
		AppName:  userApiKey.AppName,
		Duration: userApiKey.Duration,
		User:     userApiKey.User,
	}
	return apiKey
}

func (userApiKey *UserApiKey) Generate() (*ApiKeyResult, error) {
	apiKey, err := generate(userApiKey.AppName)
	if err != nil {
		return nil, err
	}
	apiKeyEntity := entities_apikey.ApiKey{
		Key:      apiKey.Result,
		UserId:   userApiKey.User.ID,
		AppName:  userApiKey.AppName,
		Duration: 365,
		Scopes:   "defaut",
		Enabled:  true,
	}
	_, err = apikey_database.Insert(apiKeyEntity)
	if err != nil {
		return nil, err
	}
	return apiKey, nil
}
