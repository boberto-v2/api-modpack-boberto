package game_routes

import (
	"net/http"
	"path/filepath"

	"github.com/brutalzinn/boberto-modpack-api/common"
	config "github.com/brutalzinn/boberto-modpack-api/configs"
	game_client_request "github.com/brutalzinn/boberto-modpack-api/domain/request/game/client"
	rest_object "github.com/brutalzinn/boberto-modpack-api/domain/rest"
	"github.com/brutalzinn/boberto-modpack-api/middlewares"
	event_service "github.com/brutalzinn/boberto-modpack-api/services/event"
	file_service "github.com/brutalzinn/boberto-modpack-api/services/file"
	ftp_models "github.com/brutalzinn/boberto-modpack-api/services/ftp/models"
	modpack_service "github.com/brutalzinn/boberto-modpack-api/services/modpack"
	modpack_cache "github.com/brutalzinn/boberto-modpack-api/services/modpack/cache"
	modpack_cache_models "github.com/brutalzinn/boberto-modpack-api/services/modpack/cache/models"
	manifest_service "github.com/brutalzinn/boberto-modpack-api/services/modpack/manifest"
	modpack_models "github.com/brutalzinn/boberto-modpack-api/services/modpack/models"
	upload_service "github.com/brutalzinn/boberto-modpack-api/services/upload"
	rest "github.com/brutalzinn/go-easy-rest"
	"github.com/gin-gonic/gin"
)

func CreateClientRoute(router gin.IRouter) {

	router.Use(createHypermediaUrl().HypermediaMiddleware())

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
		modpackCache.Name = createClientModPackRequest.Name
		modpackCache.NormalizedName = common.NormalizeString(createClientModPackRequest.Name)
		modpackCache.Status = modpack_models.PendingClientFiles

		modpack_cache.Create(modpackCache)
		outputDir := filepath.Join(modpackPath, modpack_models.Client.GetFolderName())
		uploadCache := upload_service.Create(outputDir)

		///form one to create rest
		restModPackFileObject := rest_object.New(ctx)
		restModPackFileObject.CreateModPackObject(modpackCache)
		// create a rest object to represent a upload object
		// form with more idiomatic sintax
		restUploadFileObject := rest_object.New(ctx)
		restUploadFileObject.CreateUploadFileObject(uploadCache)
		event := event_service.Create(event_service.MODPACK_PROGRESS_EVENT)
		restEventObject := rest_object.New(ctx)
		restEventObject.CreateEventObject(event)

		//i think this is more readable now.
		ctx.JSON(http.StatusOK, gin.H{
			"resource": restModPackFileObject.Resource,
		})
	})

	router.POST("/modpack/finish/:id", func(ctx *gin.Context) {
		id := ctx.Params.ByName("id")
		var finishClientModpackRequest game_client_request.FinishClientModPackRequest
		if err := ctx.ShouldBindJSON(&finishClientModpackRequest); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		modpackCache, found := modpack_cache.GetById(id)
		if !found {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "The token provided is invalid or expired"})
			return
		}
		modpack := modpack_models.MinecraftModPack{
			Name:    finishClientModpackRequest.Name,
			FileUrl: finishClientModpackRequest.FileUrl,
		}
		ftpClientConnection := ftp_models.Ftp{
			Address:   finishClientModpackRequest.ClientFtp.Address,
			User:      finishClientModpackRequest.ClientFtp.User,
			Password:  finishClientModpackRequest.ClientFtp.Password,
			Directory: finishClientModpackRequest.ClientFtp.Directory,
		}
		files := modpack_service.GetModPackFiles(modpack, modpack_models.Client)
		manifest_service.WriteModPackManifestFiles(modpack, files, modpack_models.Client)
		go modpack_service.UploadClient(modpack, ftpClientConnection)
		restObject := rest_object.New(ctx)
		restObject.CreateModPackObject(modpackCache)
		ctx.JSON(http.StatusOK, restObject)
	})
}
func createHypermediaUrl() *middlewares.Hypermedia {
	var links []rest.Link
	links = append(links, middlewares.CreateHypermediaUrl("_self", "GET", "/game/client/modpack/"))
	links = append(links, middlewares.CreateHypermediaUrl("finish", "POST", "/game/client/modpack/finish/"))
	options := middlewares.Hypermedia{
		Links: links,
	}
	hyperMedia := middlewares.New(options)
	return hyperMedia
}
