package routes

import (
	"net/http"

	login_request "github.com/brutalzinn/boberto-modpack-api/domain/models"
	authentication_user "github.com/brutalzinn/boberto-modpack-api/services/authentication/user"
	"github.com/gin-gonic/gin"
)

//TODO: Define how we will handle with user token.
//
func CreateUserRoutes(router gin.IRouter) {
	router.POST("/login", func(ctx *gin.Context) {
		var userRequest login_request.LoginRequest
		if err := ctx.ShouldBindJSON(&userRequest); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		userId, err := authentication_user.Authentication(userRequest.Email, userRequest.Password)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Email or password is invalid"})
			return
		}
		token, err := authentication_user.GenerateJWT(userId)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Email or password is invalid"})
			return
		}
		ctx.JSON(http.StatusUnauthorized, gin.H{"access_token": token})
	})
}
