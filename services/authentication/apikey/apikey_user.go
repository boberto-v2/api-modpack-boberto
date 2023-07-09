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
		Scopes:   "defaut",
		Enabled:  true,
	}
	id, err := apikey_database.Insert(apiKeyEntity)
	if err != nil {
		return nil, err
	}
	apiKey.ID = id
	return apiKey, nil
}

func (userApiKey *UserApiKey) Regenerate(apiKeyId string) (*ApiKeyResult, error) {
	appNameNormalized := normalizeAppName(userApiKey.AppName)
	apiKey, err := generate(userApiKey.AppName)
	if err != nil {
		return nil, err
	}

	apiKeyEntity := entities_apikey.ApiKey{
		ID:       apiKeyId,
		Key:      apiKey.Result,
		AppName:  appNameNormalized,
		ExpireAt: time.Now().Add(time.Duration(365) * time.Hour * 24),
		Scopes:   "defaut",
		Enabled:  true,
	}
	_, err = apikey_database.Update(apiKeyEntity)
	if err != nil {
		return nil, err
	}

	return apiKey, nil
}
