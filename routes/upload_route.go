package routes

import (
	"fmt"
	"net/http"
	"time"

	"github.com/brutalzinn/boberto-modpack-api/common"
	event_service "github.com/brutalzinn/boberto-modpack-api/services/event"
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
		go upload_service.SaveFiles(id, files)

		go func() {

			for {
				fmt.Print("testtttt")
				event_service.Emit(id, []byte("test message"))
				time.Sleep(3 * time.Second)
			}

		}()
		url := common.GetSocketUrl(ctx)
		resourceData := rest.NewResData()
		resourceData.Add(rest.Resource{
			Object:    "file_object",
			Attribute: "",
			Link: []rest.Link{
				{
					Rel:    "event_upload",
					Href:   fmt.Sprintf("%s/application/event?name=%s", url, id),
					Method: "GET",
				},
			},
		})
		ctx.JSON(http.StatusAccepted, resourceData)
	})
}
