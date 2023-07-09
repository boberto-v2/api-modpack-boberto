package routes

import (
	"fmt"
	"net/http"

	"github.com/brutalzinn/boberto-modpack-api/common"
	upload_service "github.com/brutalzinn/boberto-modpack-api/services/upload"
	rest "github.com/brutalzinn/go-easy-rest"
	"github.com/gin-gonic/gin"
)

// TODO: Show daniel how we will handle with files for all necessaries uploads

func CreateUploadRoute(router gin.IRouter) {
	router.POST("/upload/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")
		_, err := upload_service.GetUploadPath(id)
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
		upload_service.SaveFiles(id, files)
		url := common.GetUrl(ctx)
		resourceData := rest.NewResData()
		resourceData.Add(rest.Resource{
			Object:    "file_object",
			Attribute: "",
			Link: []rest.Link{
				{
					Rel:    "_self",
					Href:   fmt.Sprintf("%s/application/status/%s", url, id),
					Method: "GET",
				},
			},
		})
		ctx.JSON(http.StatusAccepted, resourceData)
	})
}
