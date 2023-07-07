package middlewares

import (
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
)

// TODO: Show to daniel how context works with Goroutines and how goroutines works and how to share data across routes
func ApiKeyMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//if x-api-key is not provided break access here
		// apiKeyHeader, err := extractApiKeyHeader(ctx.GetHeader("x-api-key"))
		// if err != nil {
		// 	ctx.AbortWithError(http.StatusUnauthorized, err)
		// 	return
		// }
		// claims, err := authentication_user.VerifyJWT(authHeaderBaerer)
		// if err != nil {
		// 	ctx.AbortWithError(http.StatusUnauthorized, err)
		// 	return
		// }
		//ctx.Set("user_id", claims.ID)
	}
}

func extractApiKeyHeader(header string) (string, error) {
	if header == "" {
		return "", errors.New("No api key provided")
	}
	jwtToken := strings.Split(header, "_")
	if len(jwtToken) != 2 {
		return "", errors.New("Invalid api key")
	}
	return jwtToken[1], nil
}
