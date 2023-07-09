package routes

import (
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/brutalzinn/boberto-modpack-api/common"
	config "github.com/brutalzinn/boberto-modpack-api/configs"
	rest_object "github.com/brutalzinn/boberto-modpack-api/domain/rest"
	file_service "github.com/brutalzinn/boberto-modpack-api/services/file"
	modpack_cache "github.com/brutalzinn/boberto-modpack-api/services/modpack/cache"
	modpack_cache_models "github.com/brutalzinn/boberto-modpack-api/services/modpack/cache/models"
	modpack_models "github.com/brutalzinn/boberto-modpack-api/services/modpack/models"
	upload_service "github.com/brutalzinn/boberto-modpack-api/services/upload"
	rest "github.com/brutalzinn/go-easy-rest"
	"github.com/gin-gonic/gin"
)

func CreateServerRoute(router gin.IRouter) {

	router.GET("/modpack/:id", func(ctx *gin.Context) {
		id := ctx.Params.ByName("id")
		modpack, _ := modpack_cache.GetById(id)
		ctx.JSON(http.StatusOK, &modpack)
	})

	router.POST("/modpack/create", func(ctx *gin.Context) {
		var cfg = config.GetConfig()
		var modpack modpack_models.MinecraftModPack
		if err := ctx.ShouldBindJSON(&modpack); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		nameNormalized := common.NormalizeString(modpack.Name)
		modpackPath := filepath.Join(cfg.API.PublicPath, nameNormalized)
		file_service.CreateAndDestroyDirectory(modpackPath)
		modpackCache := modpack_cache_models.
			ModPackCache{
			Environment: modpack_models.Server.GetFolderName(),
			Name:        modpack.Name,
		}.New()
		modpack_cache.Create(modpackCache)
		modpackCache.Status = modpack_models.PendingServerFiles
		modpack_cache.Replace(modpack.Id, modpackCache)
		//create upload ticket
		outputDir := filepath.Join(modpackPath, modpackCache.Environment)
		uploadCache := upload_service.Create(outputDir)
		restModPackFileObject := rest_object.RestObject{
			Link: []rest.Link{
				{
					Rel:    "_self",
					Href:   fmt.Sprintf("/game/server/modpack/%s", modpackCache.Id),
					Method: "GET",
				},
				{
					Rel:    "delete",
					Href:   fmt.Sprintf("/game/server/modpack/%s", modpackCache.Id),
					Method: "DELETE",
				},
				{
					Rel:    "update",
					Href:   fmt.Sprintf("/game/server/modpack/%s", modpackCache.Id),
					Method: "PUT",
				},
			},
		}.CreateModPackObject(modpackCache)

		// create a rest object to represent a upload object
		restUploadFileObject := rest_object.RestObject{
			Link: []rest.Link{
				{
					Rel:    "upload_file",
					Href:   fmt.Sprintf("/application/upload/%s", uploadCache.Id),
					Method: "POST",
				},
			},
		}.CreateFileObject(&uploadCache)

		restWaitingObject := rest_object.New(ctx).CreateWaitingObject(rest_object.WaitingObject{
			Message: rest_object.WAITING_SERVER_MESSAGE,
		})

		restResourceData := rest.NewResData()
		restResourceData.Add(restModPackFileObject.Resource)
		restResourceData.Add(restUploadFileObject.Resource)
		restResourceData.Add(restWaitingObject.Resource)
		ctx.JSON(http.StatusOK, restResourceData)
	})

	router.POST("/modpack/finish/:id", func(ctx *gin.Context) {
		id := ctx.Params.ByName("id")
		var modpackFtp modpack_models.ModPackFtp
		if err := ctx.ShouldBindJSON(&modpackFtp); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		modpack, found := modpack_cache.GetById(id)
		if !found {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "The token provided is invalid or expired"})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"data": modpack})
	})
}
