package routes

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/brutalzinn/boberto-modpack-api/common"
	config "github.com/brutalzinn/boberto-modpack-api/configs"
	file_service "github.com/brutalzinn/boberto-modpack-api/services/file"
	modpack_cache "github.com/brutalzinn/boberto-modpack-api/services/modpack/cache"
	modpack_cache_models "github.com/brutalzinn/boberto-modpack-api/services/modpack/cache/models"
	modpack_models "github.com/brutalzinn/boberto-modpack-api/services/modpack/models"
	"github.com/gin-gonic/gin"
)

var cfg = config.GetConfig()

func CreateModPackRoute(router gin.IRouter) {

	router.GET("/modpacks/:id", func(ctx *gin.Context) {
		id := ctx.Params.ByName("id")
		modpack, _ := modpack_cache.GetModpackCacheById(id)
		ctx.JSON(http.StatusOK, &modpack)
	})

	router.POST("/modpacks/create", func(ctx *gin.Context) {
		var modpack modpack_models.MinecraftModPack
		if err := ctx.ShouldBindJSON(&modpack); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		nameNormalized := common.NormalizeString(modpack.Name)
		modpackPath := filepath.Join(cfg.API.PublicPath, nameNormalized)
		dirExists := file_service.IsDirectoryExists(modpackPath)
		if dirExists {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Modpack is already in process of creation"})
			return
		}
		file_service.CreateDirectoryIfNotExists(modpackPath)
		modpackCache := modpack_cache_models.
			ModPackCache{
			Name: modpack.Name,
		}.New()
		modpack_cache.CreateModpackCache(modpackCache)
		modpack_cache.SetStatus(modpack.Id, modpack_models.PendingClientFiles)
		ctx.JSON(http.StatusOK, gin.H{"data": modpack})
	})

	router.POST("/modpacks/files/upload", func(ctx *gin.Context) {
		idQuery := ctx.Query("id")
		envQuery := ctx.Query("environment")
		modpack, found := modpack_cache.GetModpackCacheById(idQuery)
		if !found {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Modpack is not ready to upload files"})
			return
		}

		environment := modpack_models.ParseMinecraftEnvironment(envQuery)
		if modpack.Status == modpack_models.PendingClientFiles &&
			environment != modpack_models.Client {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Modpack is " + modpack_models.PendingClientFiles.GetModPackStatus()})
			return
		}
		form, _ := ctx.MultipartForm()
		if form == nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "You need upload a file or array of files"})
			return
		}
		config := config.GetConfig()
		finalModPackPath := fmt.Sprintf("%s/%s/%s", config.API.PublicPath, modpack.NormalizedName, environment.GetFolderName())
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

		if modpack.Status == modpack_models.PendingClientFiles &&
			environment == modpack_models.Client {
			modpack_cache.SetStatus(modpack.Id, modpack_models.PendingServerFiles)
		}

		if modpack.Status == modpack_models.PendingServerFiles &&
			environment == modpack_models.Server {
			modpack_cache.SetStatus(modpack.Id, modpack_models.Waiting)
		}
		//modPackFiles := modpack_service.CreateModPackFilesManifest(modpack, environment)
		//manifest_service.WriteModPackManifestFiles(modpack, modPackFiles, environment)
		//ctx.JSON(http.StatusOK, gin.H{"files": modPackFiles})
	})

	router.POST("/modpacks/:id/finish", func(ctx *gin.Context) {
		id := ctx.Params.ByName("id")
		var modpackFtp modpack_models.ModPackFtp
		if err := ctx.ShouldBindJSON(&modpackFtp); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		modpack, found := modpack_cache.GetModpackCacheById(id)
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

		/*
			[BLOG.ROBERTOCPAES.DEV - HYPERFOCUS - PERSONAL DATA COMMENTARY - IGNORE]
				CRITICAL ISSUE. I NEED STOP THIS. MY HANDS IS SHAKING TOO.
			[BLOG.ROBERTOCPAES.DEV - HYPERFOCUS - PERSONAL DATA COMMENTARY - IGNORE]
		*/

		// modpack_cache.SetStatus(modpack.Id, modpack_models.PendingFileUpload)
		// go modpack_service.UploadServer(modpack, modpackFtp)
		// go modpack_service.UploadClient(modpack, modpackFtp)
		ctx.JSON(http.StatusOK, gin.H{"data": modpack})
	})

}
