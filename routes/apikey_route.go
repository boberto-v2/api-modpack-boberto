package routes

import (
	"fmt"
	"net/http"
	"time"

	"github.com/brutalzinn/boberto-modpack-api/common"
	apikey_database "github.com/brutalzinn/boberto-modpack-api/database/apikey"
	"github.com/brutalzinn/boberto-modpack-api/domain/request"
	authentication_apikey "github.com/brutalzinn/boberto-modpack-api/services/authentication/apikey"
	authentication_user "github.com/brutalzinn/boberto-modpack-api/services/authentication/user"
	rest "github.com/brutalzinn/go-easy-rest"
	"github.com/gin-gonic/gin"
)

func CreateApiKeyRoute(router gin.IRouter) {

	router.POST("/apikey/generate", func(ctx *gin.Context) {
		currentUser, err := authentication_user.GetCurrentUser(ctx)
		var apiKeyGenerateRequest request.ApiKeyRegisterRequest
		if err := ctx.ShouldBindJSON(&apiKeyGenerateRequest); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		userApiKey := authentication_apikey.UserApiKey{
			AppName:  apiKeyGenerateRequest.AppName,
			ExpireAt: time.Now().Add(time.Duration(apiKeyGenerateRequest.Days) * time.Hour * 24),
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
					Href:   fmt.Sprintf("%s/user/apikey", url),
					Method: "GET",
				},
				{
					Rel:    "delete",
					Href:   fmt.Sprintf("%s/user/apikey/delete/%s", url, result.ID),
					Method: "DELETE",
				},
				{
					Rel:    "regenerate",
					Href:   fmt.Sprintf("%s/user/apikey/regenerate/%s", url, result.ID),
					Method: "PUT",
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
	router.PUT("/apikey/regenerate/:id", func(ctx *gin.Context) {
		id := ctx.Params.ByName("id")
		currentUser, err := authentication_user.GetCurrentUser(ctx)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		apiKeyEntity, err := apikey_database.Get(id, currentUser.ID)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		userApiKey := authentication_apikey.UserApiKey{
			AppName:  apiKeyEntity.AppName,
			ExpireAt: time.Now().Add(time.Duration(apiKeyEntity.Duration) * time.Hour * 24),
			User:     *currentUser,
		}
		newApiKey, err := userApiKey.Regenerate(id)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		// TODO: Show daniel how we will separate rest objects.
		url := common.GetUrl(ctx)
		resourceData := rest.NewResData()
		resourceData.Add(rest.Resource{
			Object: "apikey_object",
			Attribute: map[string]any{
				"api_key": newApiKey.Key,
			},
			Link: []rest.Link{
				{
					Rel:    "_self",
					Href:   fmt.Sprintf("%s/user/apikey", url),
					Method: "GET",
				},
				{
					Rel:    "delete",
					Href:   fmt.Sprintf("%s/user/apikey/delete/%s", url, id),
					Method: "DELETE",
				},
				{
					Rel:    "regenerate",
					Href:   fmt.Sprintf("%s/user/apikey/regenerate/%s", url, id),
					Method: "PUT",
				},
			},
		})
		ctx.JSON(http.StatusOK, resourceData)
	})

	router.POST("/apikey/delete/:id", func(ctx *gin.Context) {
		id := ctx.Params.ByName("id")
		currentUser, err := authentication_user.GetCurrentUser(ctx)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		apiKeyEntity, err := apikey_database.Get(id, currentUser.ID)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		_, err = apikey_database.Delete(apiKeyEntity)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, nil)
	})
}
