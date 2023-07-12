package game_routes

import (
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/brutalzinn/boberto-modpack-api/common"
	config "github.com/brutalzinn/boberto-modpack-api/configs"
	game_client_request "github.com/brutalzinn/boberto-modpack-api/domain/request/game/client"
	rest_object "github.com/brutalzinn/boberto-modpack-api/domain/rest"
	file_service "github.com/brutalzinn/boberto-modpack-api/services/file"
	modpack_cache "github.com/brutalzinn/boberto-modpack-api/services/modpack/cache"
	modpack_cache_models "github.com/brutalzinn/boberto-modpack-api/services/modpack/cache/models"
	modpack_models "github.com/brutalzinn/boberto-modpack-api/services/modpack/models"
	upload_service "github.com/brutalzinn/boberto-modpack-api/services/upload"
	rest "github.com/brutalzinn/go-easy-rest"
	"github.com/gin-gonic/gin"
)

func CreateClientRoute(router gin.IRouter) {

	router.GET("/modpack/:id", func(ctx *gin.Context) {
		id := ctx.Params.ByName("id")
		modpack, _ := modpack_cache.GetById(id)
		ctx.JSON(http.StatusOK, &modpack)
	})

	router.POST("/modpack/create", func(ctx *gin.Context) {
		var cfg = config.GetConfig()
		var createClientModPackRequest game_client_request.CreateClientModPackRequest
		if err := ctx.ShouldBindJSON(&createClientModPackRequest); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		nameNormalized := common.NormalizeString(createClientModPackRequest.Name)
		modpackPath := filepath.Join(cfg.ModPacks.PublicPath, nameNormalized)
		file_service.CreateAndDestroyDirectory(modpackPath)

		modpackCache := modpack_cache_models.
			ModPackCache{
			Environment: modpack_models.Client.GetFolderName(),
			Name:        createClientModPackRequest.Name,
		}.New()
		modpackCache.Status = modpack_models.PendingClientFiles
		modpack_cache.Create(modpackCache)
		outputDir := filepath.Join(modpackPath, modpackCache.Environment)
		uploadCache := upload_service.Create(outputDir)
		///form one to create rest
		restmodPackObject := rest_object.New(ctx)
		restmodPackObject.Link = []rest.Link{
			{
				Rel:    "_self",
				Href:   fmt.Sprintf("/game/client/modpack/%s", modpackCache.Id),
				Method: "GET",
			},
			{
				Rel:    "delete",
				Href:   fmt.Sprintf("/game/client/modpack/%s", modpackCache.Id),
				Method: "DELETE",
			},
			{
				Rel:    "update",
				Href:   fmt.Sprintf("/game/client/modpack/%s", modpackCache.Id),
				Method: "PUT",
			},
		}
		restmodPackObject.CreateModPackObject(modpackCache)
		// create a rest object to represent a upload object
		// form with more idiomatic sintax
		restUploadFileObject := rest_object.RestObject{
			Link: []rest.Link{
				{
					Rel:    "upload_file_zip",
					Href:   fmt.Sprintf("/application/upload/%s", uploadCache.Id),
					Method: "POST",
				},
			},
		}.CreateFileObject(&uploadCache)

		restWaitingObject := rest_object.New(ctx).CreateWaitingObject(rest_object.WaitingObject{
			Message: rest_object.WAITING_CLIENT_MESSAGE,
		})

		restResourceData := rest.NewResData()
		restResourceData.Add(restmodPackObject.Resource)
		restResourceData.Add(restUploadFileObject.Resource)
		restResourceData.Add(restWaitingObject.Resource)
		ctx.JSON(http.StatusOK, restResourceData)
	})

	router.POST("/modpack/finish/:id", func(ctx *gin.Context) {
		// id := ctx.Params.ByName("id")
		// var modpackFtp modpack_models.ModPackFtp
		// if err := ctx.ShouldBindJSON(&modpackFtp); err != nil {
		// 	ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		// 	return
		// }
		// modpackCache, found := modpack_cache.GetById(id)
		// if !found {
		// 	ctx.JSON(http.StatusBadRequest, gin.H{"error": "The token provided is invalid or expired"})
		// 	return
		// }

		//create modpack rest object
		// restModpackObject := rest_object.RestObject{
		// 	Attribute: modpackCache,
		// }.CreateModPackObject()

		//	ctx.JSON(http.StatusOK, gin.H{"data": resourceData})
	})
}
