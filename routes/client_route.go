package routes

import (
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/brutalzinn/boberto-modpack-api/common"
	config "github.com/brutalzinn/boberto-modpack-api/configs"
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
			Environment: modpack_models.Client.GetFolderName(),
			Name:        modpack.Name,
		}.New()
		modpack_cache.Create(modpackCache)
		modpackCache.Status = modpack_models.PendingClientFiles
		modpack_cache.Replace(modpack.Id, modpackCache)
		outputDir := filepath.Join(modpackPath, modpackCache.Environment)
		uploadCache := upload_service.Create(outputDir)
		url := common.GetUrl(ctx)
		resourceData := rest.NewResData()
		resourceData.Add(rest.Resource{
			Object:    "modpack_object",
			Attribute: modpackCache,
			Link: []rest.Link{
				{
					Rel:    "_self",
					Href:   fmt.Sprintf("%s/game/modpack/%s", url, modpackCache.Id),
					Method: "GET",
				},
				{
					Rel:    "upload_file",
					Href:   fmt.Sprintf("%s/application/upload/%s", url, uploadCache.Id),
					Method: "POST",
				},
			},
		})
		ctx.JSON(http.StatusOK, resourceData)
	})

	router.POST("/modpack/finish/:id", func(ctx *gin.Context) {
		id := ctx.Params.ByName("id")
		var modpackFtp modpack_models.ModPackFtp
		if err := ctx.ShouldBindJSON(&modpackFtp); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		modpackCache, found := modpack_cache.GetById(id)
		if !found {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "The token provided is invalid or expired"})
			return
		}
		url := common.GetUrl(ctx)
		resourceData := rest.NewResData()
		resourceData.Add(rest.Resource{
			Object:    "modpack_object",
			Attribute: modpackCache,
			Link: []rest.Link{
				{
					Rel:    "_self",
					Href:   fmt.Sprintf("%s/game/client/modpack/%s", url, modpackCache.Id),
					Method: "GET",
				},
				{
					Rel:    "delete",
					Href:   fmt.Sprintf("%s/game/client/modpack/%s", url, modpackCache.Id),
					Method: "DELETE",
				},
				{
					Rel:    "update",
					Href:   fmt.Sprintf("%s/game/client/modpack/%s", url, modpackCache.Id),
					Method: "PUT",
				},
			},
		})
		ctx.JSON(http.StatusOK, gin.H{"data": resourceData})
	})
}
