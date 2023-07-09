package routes

import (
	"fmt"
	"net/http"
	"time"

	"github.com/brutalzinn/boberto-modpack-api/common"
	apikey_database "github.com/brutalzinn/boberto-modpack-api/database/apikey"
	authentication_apikey "github.com/brutalzinn/boberto-modpack-api/services/authentication/apikey"
	authentication_user "github.com/brutalzinn/boberto-modpack-api/services/authentication/user"
	rest "github.com/brutalzinn/go-easy-rest"
	"github.com/gin-gonic/gin"
)

func CreateApiKeyRoute(router gin.IRouter) {

	router.POST("/apikey/generate", func(ctx *gin.Context) {
		currentUser, err := authentication_user.GetCurrentUser(ctx)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		userApiKey := authentication_apikey.UserApiKey{
			AppName:  "dirt",
			Duration: time.Duration(time.Hour * 24 * 365),
			User:     *currentUser,
		}
		result, err := userApiKey.Generate()
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		url := common.GetUrl(ctx)
		resourceData := rest.NewResData()
		resourceData.Add(rest.Resource{
			Object:    "apikey_object",
			Attribute: result.Key,
			Link: []rest.Link{
				{
					Rel:    "_self",
					Href:   fmt.Sprintf("%s/user/apikey/%s", url),
					Method: "GET",
				},
			},
		})
		ctx.JSON(http.StatusOK, resourceData)
	})

	router.GET("/apikey", func(ctx *gin.Context) {
		currentUser, err := authentication_user.GetCurrentUser(ctx)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		apiKeys, err := apikey_database.GetAll(currentUser.ID)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		url := common.GetUrl(ctx)
		resourceData := rest.NewResData()
		for _, apiKey := range apiKeys {
			resourceData.Add(rest.Resource{
				Object:    "apikey_object",
				Attribute: apiKey,
				Link: []rest.Link{
					{
						Rel:    "delete",
						Href:   fmt.Sprintf("%s/user/apikey/delete/%s", url, apiKey.ID),
						Method: "DELETE",
					},
					{
						Rel:    "regenerate",
						Href:   fmt.Sprintf("%s/user/apikey/regenerate/%s", url, apiKey.ID),
						Method: "PUT",
					},
				},
			})

		}
		ctx.JSON(http.StatusOK, resourceData)
	})
	router.POST("/apikey/regenerate/:id", func(ctx *gin.Context) {

	})

	router.POST("/apikey/delete/:id", func(ctx *gin.Context) {

	})
}
