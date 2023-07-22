package application_routes

import (
	"net/http"
	"time"

	"github.com/brutalzinn/boberto-modpack-api/src/src/src/domain/request"
	rest_object "github.com/brutalzinn/boberto-modpack-api/src/src/src/domain/rest"
	"github.com/brutalzinn/boberto-modpack-api/src/src/src/infra/middlewares"
	authentication_apikey "github.com/brutalzinn/boberto-modpack-api/src/src/src/infra/services/authentication/apikey"
	authentication_user "github.com/brutalzinn/boberto-modpack-api/src/src/src/infra/services/authentication/user"
	apikey_database "github.com/brutalzinn/boberto-modpack-api/src/src/src/repository/database/apikey"
	rest "github.com/brutalzinn/go-easy-rest"
	"github.com/gin-gonic/gin"
)

func CreateApiKeyRoute(router gin.IRouter) {
	router.Use(createHypermediaUrl().HypermediaMiddleware())
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
		apiKeyObject := rest_object.ApiKeyCredentialObject{
			Id:     result.ID,
			Key:    result.Key,
			Header: "x-api-key",
		}

		apiKeyCredentialsObject := rest_object.New(ctx)
		apiKeyCredentialsObject.CreateApiKeycredentialObject(apiKeyObject)
		ctx.JSON(http.StatusOK, apiKeyCredentialsObject.Resource)
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
			}
			apiKeyCredentialsObject := rest_object.New(ctx)
			apiKeyCredentialsObject.CreateApiKeycredentialObject(apiKeyObject)
			restResourceData.Add(apiKeyCredentialsObject.Resource)
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

		apiKeyObject := rest_object.ApiKeyCredentialObject{
			Id:     newApiKey.ID,
			Key:    newApiKey.Key,
			Header: "x-api-key",
		}
		apiKeyCredentialsObject := rest_object.New(ctx)
		apiKeyCredentialsObject.CreateApiKeycredentialObject(apiKeyObject)
		ctx.JSON(http.StatusOK, apiKeyCredentialsObject.Resource)
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

func createHypermediaUrl() *middlewares.Hypermedia {
	var links []rest.Link
	// links = append(links, middlewares.CreateHypermediaUrl("_self", "GET", "/user/apikey/"))
	// links = append(links, middlewares.CreateHypermediaUrl("delete", "DELETE", "/user/apikey/delete/"))
	// links = append(links, middlewares.CreateHypermediaUrl("regenerate", "PUT", "/user/apikey/regenerate/"))
	options := middlewares.Hypermedia{
		Links: links,
	}
	hyperMedia := middlewares.New(options)
	return hyperMedia
}
