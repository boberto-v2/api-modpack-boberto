package application_routes

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/brutalzinn/boberto-modpack-api/common"
	rest_object "github.com/brutalzinn/boberto-modpack-api/domain/rest"
	event_service "github.com/brutalzinn/boberto-modpack-api/infra/services/event"
	event_rest "github.com/brutalzinn/boberto-modpack-api/infra/services/event/rest"
	file_service "github.com/brutalzinn/boberto-modpack-api/infra/services/file"
	upload_service "github.com/brutalzinn/boberto-modpack-api/infra/services/upload"
	upload_cache "github.com/brutalzinn/boberto-modpack-api/infra/services/upload/cache"
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
					messageEvent, err := event_rest.CreateMessageEventObject(err.Error())
					if err != nil {
						return
					}
					event.Emit(messageEvent)
				}
				ctx.String(http.StatusBadRequest, "upload file err: %s", err.Error())
				return
			}

			isZip := file_service.IsZip(filename)
			if isZip {
				file_service.UnZip(filePath, uploadCache.OutputDir, func(percentage common.ProgressCalculator) {
					fileUploadEvent, err := event_rest.CreateFileUploadEventObject(percentage.Progress)
					if err != nil {
						return
					}
					if eventFound {
						event.Emit(fileUploadEvent)
					}
				})
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
