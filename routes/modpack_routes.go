package routes

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	config "github.com/brutalzinn/boberto-modpack-api/configs"
	"github.com/brutalzinn/boberto-modpack-api/models"
	modpack_cache "github.com/brutalzinn/boberto-modpack-api/services/cache"
	services_cache "github.com/brutalzinn/boberto-modpack-api/services/cache"
	file_service "github.com/brutalzinn/boberto-modpack-api/services/file"
	manifest_service "github.com/brutalzinn/boberto-modpack-api/services/manifest"
	modpack_service "github.com/brutalzinn/boberto-modpack-api/services/modpack"
	"github.com/gin-gonic/gin"
)

func CreateModPackRoute(router gin.IRouter) {

	router.GET("/modpacks/:id", func(ctx *gin.Context) {
		id := ctx.Params.ByName("id")
		modpack, _ := services_cache.GetModpackCacheById(id)
		ctx.JSON(http.StatusOK, &modpack)
	})

	router.POST("/modpacks/create", func(ctx *gin.Context) {
		var modpack models.Modpack
		if err := ctx.ShouldBindJSON(&modpack); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		modpack.New()
		config := config.GetConfig()
		modpackPath := filepath.Join(config.PublicPath, modpack.NormalizedName)
		dirExists := file_service.IsDirectoryExists(modpackPath)
		if dirExists {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Modpack is already in process of creation"})
			return
		}
		file_service.CreateDirectoryIfNotExists(modpackPath)
		services_cache.CreateModpackCache(modpack)
		manifest_service.NewModPackManifest(modpack)
		services_cache.SetStatus(modpack.Id, models.PendingClientFiles)
		ctx.JSON(http.StatusOK, gin.H{"data": modpack})
	})

	router.POST("/modpacks/files/upload", func(ctx *gin.Context) {
		idQuery := ctx.Query("id")
		envQuery := ctx.Query("environment")
		modpack, found := services_cache.GetModpackCacheById(idQuery)
		if !found {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Modpack is not ready to upload files"})
			return
		}

		environment := models.ParseMinecraftEnvironment(envQuery)
		if modpack.Status == models.PendingClientFiles &&
			environment != models.Client {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Modpack is " + models.PendingClientFiles.GetModPackStatus()})
			return
		}
		form, _ := ctx.MultipartForm()
		if form == nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "You need upload a file or array of files"})
			return
		}
		config := config.GetConfig()
		finalModPackPath := fmt.Sprintf("%s/%s/%s", config.PublicPath, modpack.NormalizedName, environment.GetFolderName())
		file_service.CreateDirectoryIfNotExists(finalModPackPath)
		files := form.File["files"]
		for _, file := range files {
			finalZipFile := filepath.Join(finalModPackPath, file.Filename)
			out, err := os.Create(finalZipFile)
			if err != nil {
				log.Fatal(err)
			}
			defer out.Close()
			readerFile, _ := file.Open()
			_, err = io.Copy(out, readerFile)
			if err != nil {
				log.Fatal(err)
			}
			file_service.Unzip(finalZipFile, finalModPackPath)
			os.Remove(finalZipFile)
		}

		if modpack.Status == models.PendingClientFiles &&
			environment == models.Client {
			modpack_cache.SetStatus(modpack.Id, models.PendingServerFiles)
		}

		if modpack.Status == models.PendingServerFiles &&
			environment == models.Server {
			modpack_cache.SetStatus(modpack.Id, models.Waiting)
		}
		modPackFiles := modpack_service.CreateModPackFile(modpack, environment)
		manifest_service.NewManifest(modpack, modPackFiles, environment)
		ctx.JSON(http.StatusOK, gin.H{"files": modPackFiles})
	})

	router.POST("/modpacks/:id/finish", func(ctx *gin.Context) {
		id := ctx.Params.ByName("id")
		var modpackFtp models.ModPackFtp
		if err := ctx.ShouldBindJSON(&modpackFtp); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		modpack, found := services_cache.GetModpackCacheById(id)
		if !found {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Modpack is not ready to finish"})
			return
		}
		// if modpack.Status == models.PendingClientFiles ||
		// 	modpack.Status == models.PendingServerFiles ||
		// 	modpack.Status == models.PendingFileUpload {
		// 	ctx.JSON(http.StatusBadRequest, gin.H{"error": "Modpack is " + modpack.Status.GetModPackStatus()})
		// 	return
		// }
		services_cache.SetStatus(modpack.Id, models.PendingFileUpload)
		go modpack_service.UploadServer(modpack, modpackFtp)
		//go modpack_service.UploadClient(modpack, modpackFtp)
		ctx.JSON(http.StatusOK, gin.H{"data": modpack})
	})
}
