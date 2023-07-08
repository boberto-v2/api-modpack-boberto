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

// TODO: Show daniel how we will handle with files for all necessaries uploads
func CreateClientRoute(router gin.IRouter) {

	router.GET("/modpack/:id", func(ctx *gin.Context) {
		id := ctx.Params.ByName("id")
		modpack, _ := modpack_cache.GetById(id)
		ctx.JSON(http.StatusOK, &modpack)
	})

	router.POST("/modpacks/create", func(ctx *gin.Context) {
		var cfg = config.GetConfig()
		var modpack modpack_models.MinecraftModPack
		if err := ctx.ShouldBindJSON(&modpack); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		nameNormalized := common.NormalizeString(modpack.Name)
		modpackPath := filepath.Join(cfg.API.PublicPath, nameNormalized)
		file_service.CreateDirectoryIfNotExists(modpackPath)
		modpackCache := modpack_cache_models.
			ModPackCache{
			Name: modpack.Name,
		}.New()
		modpack_cache.Create(modpackCache)
		modpackCache.Status = modpack_models.PendingClientFiles
		modpack_cache.Replace(modpack.Id, modpackCache)
		uploadCache := upload_service.Create(modpack_models.Client.GetFolderName())
		//mount robject esource representation of modpack
		url := common.GetUrl(ctx)
		resourceData := rest.NewResData()
		resourceData.Add(rest.Resource{
			Object:    "modpack",
			Attribute: modpackCache,
			Link: []rest.Link{
				{
					Rel:  "_self",
					Href: fmt.Sprintf("%s/modpack/%s", url, modpackCache.Id),
				},
				{
					Rel:  "upload_client_file",
					Href: fmt.Sprintf("%s/upload/%s", url, uploadCache.Id),
				},
			},
		})
		ctx.JSON(http.StatusOK, resourceData)
	})
}
