package routes

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	config "github.com/brutalzinn/go-multiple-file/configs"
	"github.com/brutalzinn/go-multiple-file/models"
	services_cache "github.com/brutalzinn/go-multiple-file/services/cache"
	file_services "github.com/brutalzinn/go-multiple-file/services/file"
	manifest_service "github.com/brutalzinn/go-multiple-file/services/manifest"
	"github.com/gin-gonic/gin"
)

func CreateModPackRoute(router gin.IRouter) {

	router.GET("/modpacks/:id", func(ctx *gin.Context) {
		id := ctx.Params.ByName("id")
		modpack := services_cache.GetModpackCacheById(id)
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
		modpackPath := fmt.Sprintf("%s/%s", config.PublicPath, modpack.NormalizedName)
		dirExists := file_services.IsDirectoryExists(modpackPath)
		if dirExists {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Modpack is already in process of creation"})
			return
		}
		services_cache.CreateModpackCache(&modpack)
		file_services.CreateDirectoryIfNotExists(modpackPath)
		manifest_service.NewModPackManifest(modpack)
		ctx.JSON(http.StatusOK, gin.H{"data": modpack})
	})

	router.POST("/modpacks/files/upload", func(ctx *gin.Context) {
		idQuery := ctx.Query("id")
		envQuery := ctx.Query("environment")

		modpack := services_cache.GetModpackCacheById(idQuery)
		if modpack == nil {
			ctx.Status(http.StatusBadRequest)
			return
		}
		environment := models.ParseMinecraftEnvironment(envQuery)
		form, _ := ctx.MultipartForm()
		config := config.GetConfig()
		files := form.File["files"]
		modpackFiles := []models.ModPackFile{}
		for _, file := range files {
			originalFileName := strings.TrimSuffix(filepath.Base(file.Filename),
				filepath.Ext(file.Filename))
			finalModPackPath := fmt.Sprintf("%s/%s/%s", config.PublicPath, modpack.NormalizedName, environment.GetFolderName())
			file_services.CreateDirectoryIfNotExists(finalModPackPath)
			finalPath := fmt.Sprintf("%s/%s", finalModPackPath, originalFileName)
			out, err := os.Create(finalPath)
			if err != nil {
				log.Fatal(err)
			}
			defer out.Close()

			readerFile, _ := file.Open()
			_, err = io.Copy(out, readerFile)
			if err != nil {
				log.Fatal(err)
			}
			checksum, err := file_services.GetCRC32(finalPath)
			if err != nil {
				log.Fatal(err)
			}
			modpackFile := models.ModPackFile{
				Name:        originalFileName,
				Path:        finalPath,
				Checksum:    checksum,
				Environment: environment,
			}
			modpackFiles = append(modpackFiles, modpackFile)
		}

		manifest_service.NewManifest(modpack, modpackFiles, environment)
		ctx.JSON(http.StatusOK, gin.H{"files": modpackFiles})
	})
}
