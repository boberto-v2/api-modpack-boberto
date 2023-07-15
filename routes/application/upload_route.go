package application_routes

import (
	"net/http"

	rest_object "github.com/brutalzinn/boberto-modpack-api/domain/rest"
	event_service "github.com/brutalzinn/boberto-modpack-api/services/event"
	upload_service "github.com/brutalzinn/boberto-modpack-api/services/upload"
	upload_cache "github.com/brutalzinn/boberto-modpack-api/services/upload/cache"
	"github.com/gin-gonic/gin"
)

// TODO: Show daniel how we will handle with files for all necessaries uploads
func CreateUploadRoute(router gin.IRouter) {
	router.POST("/upload/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")
		uploadCache, err := upload_service.GetById(id)
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
		eventId := form.Value["event"][0]
		go upload_service.SaveFiles(uploadCache.OutputDir, files, func(file string) {
			EmitIfNecessary(eventId, "saving.. "+file)
			UnzipIfNecessary(eventId, file, uploadCache.OutputDir)
		})
		uploadCache.Status = upload_cache.UPLOAD_COMPLETED
		uploadCache.Save()
		restUploadFileObject := rest_object.New(ctx)
		restUploadFileObject.CreateUploadFileObject(uploadCache)
		ctx.JSON(http.StatusAccepted, restUploadFileObject.Resource)
	})

}

func EmitIfNecessary(eventId, message string) {
	event, found := event_service.GetById(eventId)
	if !found {
		return
	}
	event.Emit(message)
}
func UnzipIfNecessary(eventId, file, outputDir string) {
	isZip := upload_service.IsZip(file)
	if isZip {
		upload_service.UnZip(file, outputDir, func(s string) {
			EmitIfNecessary(eventId, "unzip.. "+s)
		})
	}
}
