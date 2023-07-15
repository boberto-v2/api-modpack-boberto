package middlewares

import (
	"net/http"

	config "github.com/brutalzinn/boberto-modpack-api/configs"
	authentication_apikey "github.com/brutalzinn/boberto-modpack-api/services/authentication/apikey"
	"github.com/gin-gonic/gin"
)

// TODO: Show to daniel how context works with Goroutines and how goroutines works and how to share data across routes
func ApiKeyMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		cfg := config.GetConfig()
		apiKeyHeader := ctx.GetHeader(cfg.API.ApiKeyHeader)
		if ctx.Value("user_id") != nil {
			ctx.Next()
			return
		}
		apiKey, err := authentication_apikey.GetApiKeyByHeaderValue(apiKeyHeader)
		if err != nil {
			ctx.AbortWithError(http.StatusUnauthorized, err)
			return
		}
		ctx.Set("user_id", apiKey.User.ID)
		ctx.Next()
	}
}
