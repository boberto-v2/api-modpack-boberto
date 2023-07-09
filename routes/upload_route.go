package routes

import (
	"fmt"
	"net/http"

	rest_object "github.com/brutalzinn/boberto-modpack-api/domain/rest"
	upload_service "github.com/brutalzinn/boberto-modpack-api/services/upload"
	rest "github.com/brutalzinn/go-easy-rest"
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
		go upload_service.SaveFiles(id, files)
		restUploadFileObject := rest_object.RestObject{
			Link: []rest.Link{
				{
					Rel:    "upload_file",
					Href:   fmt.Sprintf("/application/upload/%s", uploadCache.Id),
					Method: "POST",
				},
			},
		}.CreateFileObject(uploadCache)

		ctx.JSON(http.StatusAccepted, restUploadFileObject)
	})
}
