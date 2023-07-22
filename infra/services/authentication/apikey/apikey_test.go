package authentication_apikey

import (
	"testing"
	"time"

	"github.com/brutalzinn/boberto-modpack-api/common"
	user_database "github.com/brutalzinn/boberto-modpack-api/repository/database/user"
	entities_user "github.com/brutalzinn/boberto-modpack-api/repository/database/user/entities"
	"github.com/stretchr/testify/assert"
)

// TODO: Explain to Daniel how we start create unit test and integration tests
func TestGenerate(t *testing.T) {
	t.Setenv("PG_URI", "postgres://root:test@127.0.0.1:5555/test?sslmode=disable")
	user, err := CreateFakeUser()
	if err != nil {
		t.Error(err)
	}
	// one year api key
	userApiKey := UserApiKey{
		AppName:  "bricks",
		ExpireAt: time.Now().Add(time.Duration(365) * time.Hour * 24),
		User:     *user,
	}
	result, err := userApiKey.Generate()
	if err != nil {
		t.Log(err)
		return
	}
	isValid := common.BcryptCheckHash(result.Key, result.Result)
	assert.True(t, isValid)
	t.Log(result.Key, result.Result)
}

func TestValid(t *testing.T) {
	t.Setenv("PG_URI", "postgres://root:test@127.0.0.1:5555/test?sslmode=disable")
	user, err := CreateFakeUser()
	if err != nil {
		t.Error(err)
	}
	// one year api key
	userApiKey := UserApiKey{
		AppName:  "dirt",
		ExpireAt: time.Now().Add(time.Duration(365) * time.Hour * 24),
		User:     *user,
	}
	result, err := userApiKey.Generate()
	if err != nil {
		return
	}

	apiKeyUser := result.Key
	t.Log(apiKeyUser)
	apiKey, err := GetApiKeyByHeaderValue(apiKeyUser)
	t.Logf("err %s", err)
	t.Log(apiKey)
	assert.NoError(t, err)
}

func CreateFakeUser() (*entities_user.User, error) {
	password, _ := common.BcryptHash("123", 4)
	user := entities_user.User{
		Password: password,
		Email:    "test@test.com",
		Username: "test",
	}
	id, err := user_database.Insert(user)
	user.ID = id
	if err != nil {
		return nil, err
	}
	return &user, nil
}
