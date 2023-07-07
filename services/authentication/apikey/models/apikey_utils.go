package authentication_apikey

import (
	"errors"
	"fmt"
	"strings"

	"github.com/brutalzinn/boberto-modpack-api/common"
)

func createApiPrefix(appName string, apiKeyCrypt string) string {
	return fmt.Sprintf("%s_%s", appName, apiKeyCrypt)
}

func extractApiKey(apiKeyCrypt string) (*ApiKey, error) {
	apikeyformat := strings.Split(apiKeyCrypt, "_")
	if len(apikeyformat) != 2 {
		return nil, errors.New("Api key invalid")
	}
	apiKeyDecoded, err := common.DecodeBase64(apikeyformat[1])
	if err != nil {
		return nil, errors.New("Api key invalid")
	}
	apiKey := ApiKey{
		AppName: apikeyformat[0],
		Key:     string(apiKeyDecoded),
	}
	return &apiKey, nil
}

func generate(appName string) (*ApiKeyResult, error) {
	uuid := common.GenerateUUID()
	uuidBase64 := common.EncodeBase64([]byte(uuid))
	appNameNormalized := common.NormalizeString(appName)
	key := createApiPrefix(appNameNormalized, uuidBase64)
	keyHash, err := common.BcryptHash(key, 4)
	if err != nil {
		return nil, errors.New("cant generate api key")
	}
	result := &ApiKeyResult{
		Key:    key,
		Result: keyHash,
	}
	return result, nil
}
