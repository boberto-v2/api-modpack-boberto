package application_routes

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	rest_object "github.com/brutalzinn/boberto-modpack-api/domain/rest"
	event_service "github.com/brutalzinn/boberto-modpack-api/services/event"
	file_service "github.com/brutalzinn/boberto-modpack-api/services/file"
	upload_service "github.com/brutalzinn/boberto-modpack-api/services/upload"
	upload_cache "github.com/brutalzinn/boberto-modpack-api/services/upload/cache"
	"github.com/gin-gonic/gin"
)

// TODO: Show daniel how we will handle with files for all necessaries uploads
func CreateUploadRoute(router gin.IRouter) {
	router.POST("/upload/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")
		eventId := ctx.Query("event")
		uploadCache, err := upload_service.GetById(id)
		event, eventFound := event_service.GetById(eventId)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		form, _ := ctx.MultipartForm()
		if form == nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "You need upload a file or array of files"})
			return
		}
		files := form.File["files"]
		for _, file := range files {
			filename := filepath.Base(file.Filename)
			filePath := filepath.Join(uploadCache.OutputDir, filename)
			if err := ctx.SaveUploadedFile(file, filePath); err != nil {
				if eventFound {
					event.Emit(fmt.Sprint("upload file err: %s", err.Error()))
				}
				ctx.String(http.StatusBadRequest, "upload file err: %s", err.Error())
				return
			}
			if eventFound {
				event.Emit(fmt.Sprint("uploaded file %s", filename))
			}
			isZip := file_service.IsZip(filename)
			if isZip {
				file_service.UnZip(filePath, uploadCache.OutputDir)
				os.Remove(filePath)
			}

		}
		uploadCache.Status = upload_cache.UPLOAD_COMPLETED
		uploadCache.Save()
		restUploadFileObject := rest_object.New(ctx)
		restUploadFileObject.CreateUploadFileObject(uploadCache)
		ctx.JSON(http.StatusAccepted, restUploadFileObject.Resource)
	})
}
