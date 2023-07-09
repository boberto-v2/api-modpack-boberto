package routes

import (
	"fmt"
	"net/http"
	"time"

	apikey_database "github.com/brutalzinn/boberto-modpack-api/database/apikey"
	"github.com/brutalzinn/boberto-modpack-api/domain/request"
	rest_object "github.com/brutalzinn/boberto-modpack-api/domain/rest"
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
		restResourceData := rest.NewResData()
		apiKeyObject := rest_object.ApiKeyCredentialObject{
			Id:     result.ID,
			Key:    result.Key,
			Header: "x-api-key",
			Link: []rest.Link{
				{
					Rel:    "_self",
					Href:   fmt.Sprintf("/user/apikey/%s", result.ID),
					Method: "GET",
				},
				{
					Rel:    "delete",
					Href:   fmt.Sprintf("/user/apikey/delete/%s", result.ID),
					Method: "DELETE",
				},
				{
					Rel:    "regenerate",
					Href:   fmt.Sprintf("/user/apikey/regenerate/%s", result.ID),
					Method: "PUT",
				},
			},
		}
		restObject := rest_object.New(ctx).CreateApiKeycredentialObject(apiKeyObject)
		restResourceData.Add(restObject.Resource)
		ctx.JSON(http.StatusOK, restResourceData)
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

		restResourceData := rest.NewResData()
		for _, apiKey := range apiKeys {
			apiKeyObject := rest_object.ApiKeyCredentialObject{
				Id:     apiKey.ID,
				Header: "x-api-key",
				Scopes: apiKey.Scopes,
				Link: []rest.Link{
					{
						Rel:    "_self",
						Href:   fmt.Sprintf("/user/apikey/%s", apiKey.ID),
						Method: "GET",
					},
					{
						Rel:    "delete",
						Href:   fmt.Sprintf("/user/apikey/delete/%s", apiKey.ID),
						Method: "DELETE",
					},
					{
						Rel:    "regenerate",
						Href:   fmt.Sprintf("/user/apikey/regenerate/%s", apiKey.ID),
						Method: "PUT",
					},
				},
			}
			restObject := rest_object.New(ctx).CreateApiKeycredentialObject(apiKeyObject)
			restResourceData.Add(restObject.Resource)
		}

		ctx.JSON(http.StatusOK, restResourceData)
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
			Id:       apiKeyEntity.ID,
			AppName:  apiKeyEntity.AppName,
			ExpireAt: time.Now().Add(time.Duration(apiKeyEntity.Duration) * time.Hour * 24),
			User:     *currentUser,
		}
		newApiKey, err := userApiKey.Regenerate()
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		// TODO: Show daniel how we will separate rest objects.
		restResourceData := rest.NewResData()
		apiKeyObject := rest_object.ApiKeyCredentialObject{
			Id:     newApiKey.ID,
			Header: "x-api-key",
			Link: []rest.Link{
				{
					Rel:    "_self",
					Href:   fmt.Sprintf("/user/apikey/%s", apiKeyEntity.ID),
					Method: "GET",
				},
				{
					Rel:    "delete",
					Href:   fmt.Sprintf("/user/apikey/delete/%s", apiKeyEntity.ID),
					Method: "DELETE",
				},
				{
					Rel:    "regenerate",
					Href:   fmt.Sprintf("/user/apikey/regenerate/%s", apiKeyEntity.ID),
					Method: "PUT",
				},
			},
		}
		restObject := rest_object.New(ctx).CreateApiKeycredentialObject(apiKeyObject)
		restResourceData.Add(restObject.Resource)
		ctx.JSON(http.StatusOK, restResourceData)
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
