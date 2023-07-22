package authorization_routes

import (
	"net/http"

	"github.com/brutalzinn/boberto-modpack-api/common"
	"github.com/brutalzinn/boberto-modpack-api/domain/request"
	rest_object "github.com/brutalzinn/boberto-modpack-api/domain/rest"
	authentication_user "github.com/brutalzinn/boberto-modpack-api/infra/services/authentication/user"
	user_database "github.com/brutalzinn/boberto-modpack-api/repository/database/user"
	entities_user "github.com/brutalzinn/boberto-modpack-api/repository/database/user/entities"
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
		//The resource can change. REMEMBER IT.
		//At this case, the user can be blocked because many attemps to login, in this case.. waiting object will received too.
		//Soo.. the client needs to wait for a array of object and handle for itself.
		userCredentialsObject := rest_object.New(ctx)
		userCredentialsObject.CreateUserCredentialsObject(token)
		ctx.JSON(http.StatusOK, gin.H{"data": userCredentialsObject.Resource})
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
