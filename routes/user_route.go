package routes

import (
	"net/http"

	"github.com/brutalzinn/boberto-modpack-api/common"
	user_database "github.com/brutalzinn/boberto-modpack-api/database/user"
	entities_user "github.com/brutalzinn/boberto-modpack-api/database/user/entities"
	"github.com/brutalzinn/boberto-modpack-api/domain/request"
	authentication_user "github.com/brutalzinn/boberto-modpack-api/services/authentication/user"
	"github.com/gin-gonic/gin"
)

// TODO: Define how we will handle with user token.
func CreateAuthRoutes(router gin.IRouter) {
	router.POST("/login", func(ctx *gin.Context) {
		var userRequest request.LoginRequest
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

	router.POST("/register", func(ctx *gin.Context) {
		var registerRequest request.RegisterRequest
		if err := ctx.ShouldBindJSON(&registerRequest); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		_, err := user_database.FindByEmail(registerRequest.Email)
		if err == nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Something goes wrong"})
			return
		}
		password, _ := common.BcryptHash(registerRequest.Password, 8)
		user := entities_user.User{
			Password: password,
			Email:    registerRequest.Email,
			Username: registerRequest.Username,
		}
		_, err = user_database.Insert(user)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, nil)
	})
}
