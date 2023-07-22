package authentication_apikey

import (
	"errors"
	"fmt"
	"strings"

	"github.com/brutalzinn/boberto-modpack-api/common"
)

func createApiPrefix(appName string, apiKeyCrypt string) string {
	return fmt.Sprintf("%s&%s", appName, apiKeyCrypt)
}

func extractAppName(apiKeyCrypt string) (string, error) {
	appName, _, found := strings.Cut(apiKeyCrypt, "&")
	if !found {
		return "", errors.New("invalid api key")
	}
	return appName, nil
}

func generate(appName string) (*ApiKeyResult, error) {
	uuid := common.GenerateUUID()
	base64 := common.EncodeBase64([]byte(uuid))
	key := createApiPrefix(appName, base64)
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

func normalizeAppName(appName string) string {
	randomFactor := common.CreateRandomFactor()
	result := fmt.Sprintf("%s_%s", appName, randomFactor)
	return result
}
