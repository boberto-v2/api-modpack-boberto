package routes

import (
	"net/http"

	rest_object "github.com/brutalzinn/boberto-modpack-api/domain/rest"
	event_service "github.com/brutalzinn/boberto-modpack-api/services/event"
	upload_service "github.com/brutalzinn/boberto-modpack-api/services/upload"
	"github.com/gin-gonic/gin"
)

// TODO: Show daniel how we will handle with files for all necessaries uploads
func CreateUploadRoute(router gin.IRouter) {
	router.POST("/upload/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")
		_, err := upload_service.GetById(id)
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
		event := event_service.Create(event_service.FILE_UPLOAD)
		event.Emit("Saving local files..")
		go upload_service.SaveFiles(id, files)
		restUploadFileObject := rest_object.New(ctx)
		restUploadFileObject.CreateEventObject(event)
		ctx.JSON(http.StatusAccepted, restUploadFileObject)
	})
}
