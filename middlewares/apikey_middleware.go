package middlewares

import (
	"net/http"

	authentication_apikey "github.com/brutalzinn/boberto-modpack-api/services/authentication/apikey"
	"github.com/gin-gonic/gin"
)

// TODO: Show to daniel how context works with Goroutines and how goroutines works and how to share data across routes
func ApiKeyMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		apiKeyHeader, err := authentication_apikey.GetApiKeyByHeaderValue(ctx.GetHeader("x-api-key"))
		if err != nil {
			ctx.AbortWithError(http.StatusUnauthorized, err)
			return
		}
		ctx.Set("user_id", apiKeyHeader.User.ID)
		ctx.Next()
	}
}
