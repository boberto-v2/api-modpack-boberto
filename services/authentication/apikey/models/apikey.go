package authentication_apikey

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/brutalzinn/boberto-modpack-api/common"
	config "github.com/brutalzinn/boberto-modpack-api/configs"
	"github.com/brutalzinn/boberto-modpack-api/database/user/entities"
)

var cfg = config.GetConfig()

type ApiKey struct {
	Key      string
	AppName  string
	User     entities.User
	Duration time.Duration
	ExpireAt time.Time
	CreateAt time.Time
	UpdateAt time.Time
}

func New(duration time.Duration, user entities.User) *ApiKey {
	apiKey := &ApiKey{
		Duration: duration,
		User:     user,
	}
	return apiKey
}

func (apiKey *ApiKey) Mount() *ApiKey {
	key, _ := createKey()
	apiKey.Key = key

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
func createApiPrefix(apiKeyCrypt string, appName string) string {
	return fmt.Sprintf("%s-%s", appName, apiKeyCrypt)
}

func removeApiPrefix(apiKeyCrypt string) (string, error) {
	apikeyformat := strings.Split(apiKeyCrypt, "-")
	if len(apikeyformat) != 2 {
		return "", errors.New("Api key invalid")
	}
	apiKey := apikeyformat[1]
	return apiKey, nil
}

func createKey() (string, error) {
	uuid := common.GenerateUUID()
	keyhash, err := common.Encrypt(uuid, cfg.Authentication.AesKey)
	if err != nil {
		return "", errors.New("Invalid aes key")
	}
	return keyhash, nil
}

func (apiKey *ApiKey) AddExpire(expireAt time.Duration) *ApiKey {
	apiKey.ExpireAt = time.Now().Add(expireAt)
	return apiKey
}

func (apiKey *ApiKey) Regenerate() *ApiKey {
	apiKey.ExpireAt = time.Now().Add(apiKey.Duration)
	return apiKey
}

func (apiKey *ApiKey) Revoke() {

}
