package middlewares

import (
	"errors"
	"net/http"
	"strings"

	authentication_user "github.com/brutalzinn/boberto-modpack-api/services/authentication/user"
	"github.com/gin-gonic/gin"
)

// TODO: Show to daniel how context works with Goroutines and how goroutines works and how to share data across routes
func JWTMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//go to next middleware if x-api-key is provided
		authHeaderApiKey := ctx.GetHeader("x-api-key")
		if authHeaderApiKey != "" {
			ctx.Next()
		}
		authHeaderBaerer, err := extractBearerToken(ctx.GetHeader("Authorization"))
		if err != nil {
			ctx.AbortWithError(http.StatusUnauthorized, err)
			return
		}
		claims, err := authentication_user.VerifyJWT(authHeaderBaerer)
		if err != nil {
			ctx.AbortWithError(http.StatusUnauthorized, err)
			return
		}
		ctx.Set("user_id", claims.ID)
		ctx.Next()
	}
}

func extractBearerToken(header string) (string, error) {
	if header == "" {
		return "", errors.New("No header provided")
	}
	jwtToken := strings.Split(header, " ")
	if len(jwtToken) != 2 {
		return "", errors.New("Invalid JWT")
	}

	return jwtToken[1], nil
}
