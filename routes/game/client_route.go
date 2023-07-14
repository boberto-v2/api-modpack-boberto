package game_routes

import (
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

		modpackCache := modpack_cache_models.New()
		modpackCache.Environment = ""
		modpackCache.Name = createClientModPackRequest.Name
		modpackCache.NormalizedName = common.NormalizeString(createClientModPackRequest.Name)
		modpackCache.Status = modpack_models.PendingClientFiles

		modpack_cache.Create(modpackCache)
		// outputDir := filepath.Join(modpackPath, modpackCache.Environment)
		// uploadCache := upload_service.Create(outputDir)

		///form one to create rest
		restModPackFileObject := rest_object.New(ctx)
		restModPackFileObject.CreateModPackObject(modpackCache)
		// create a rest object to represent a upload object
		// form with more idiomatic sintax
		// restUploadFileObject := rest_object.New(ctx)
		// restUploadFileObject.CreateFileObject(&uploadCache)
		// restResourceData := rest.NewResData()
		// restResourceData.Add(restModPackFileObject.Resource)
		// restResourceData.Add(restUploadFileObject.Resource)
		ctx.JSON(http.StatusOK, restModPackFileObject)
	})

	router.POST("/modpack/finish/:id", func(ctx *gin.Context) {
		// id := ctx.Params.ByName("id")
		// var finishClientModpackRequest game_client_request.FinishClientModPackRequest
		// if err := ctx.ShouldBindJSON(&finishClientModpackRequest); err != nil {
		// 	ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		// 	return
		// }
		// modpackCache, found := modpack_cache.GetById(id)
		// if !found {
		// 	ctx.JSON(http.StatusBadRequest, gin.H{"error": "The token provided is invalid or expired"})
		// 	return
		// }
		// modpack := modpack_models.MinecraftModPack{
		// 	Name: modpackCache.Name,
		// }
		// ftpClientConnection := ftp_models.Ftp{
		// 	Address:   finishClientModpackRequest.ClientFtp.Address,
		// 	User:      finishClientModpackRequest.ClientFtp.User,
		// 	Password:  finishClientModpackRequest.ClientFtp.Password,
		// 	Directory: finishClientModpackRequest.ClientFtp.Directory,
		// }
		// modpack_service.UploadClient(modpack, ftpClientConnection)
		// // create modpack rest object
		// restObject := rest_object.New(ctx).CreateModPackObject(modpackCache)
		// ctx.JSON(http.StatusOK, restObject)
	})
}
