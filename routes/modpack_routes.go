package routes

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	config "github.com/brutalzinn/boberto-modpack-api/configs"
	"github.com/brutalzinn/boberto-modpack-api/models"
	services_cache "github.com/brutalzinn/boberto-modpack-api/services/cache"
	file_service "github.com/brutalzinn/boberto-modpack-api/services/file"
	manifest_service "github.com/brutalzinn/boberto-modpack-api/services/manifest"
	modpack_service "github.com/brutalzinn/boberto-modpack-api/services/modpack"
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
		dirExists := file_service.IsDirectoryExists(modpackPath)
		if dirExists {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Modpack is already in process of creation"})
			return
		}
		services_cache.CreateModpackCache(&modpack)
		file_service.CreateDirectoryIfNotExists(modpackPath)
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
		finalModPackPath := fmt.Sprintf("%s/%s/%s", config.PublicPath, modpack.NormalizedName, environment.GetFolderName())
		file_service.CreateDirectoryIfNotExists(finalModPackPath)

		files := form.File["files"]
		for _, file := range files {
			// originalFileName := strings.TrimSuffix(filepath.Base(file.Filename),
			// 	filepath.Ext(file.Filename))
			finalZipFile := fmt.Sprintf("%s/%s", finalModPackPath, file.Filename)
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
			// checksum, err := file_service.GetCRC32(finalZipFile)
			// if err != nil {
			// 	log.Fatal(err)
			// }
			file_service.Unzip(finalZipFile, finalModPackPath)
			// modpackFile := models.ModPackFile{
			// 	Name:        originalFileName,
			// 	Path:        finalPath,
			// 	Checksum:    checksum,
			// 	Environment: environment,
			// }
			// modpackFiles = append(modpackFiles, modpackFile)
		}
		modPackFiles := modpack_service.CreateModPackFile(modpack, environment)
		manifest_service.NewManifest(modpack, modPackFiles, environment)
		ctx.JSON(http.StatusOK, gin.H{"files": modPackFiles})
	})
}
