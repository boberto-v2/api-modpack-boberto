package routes

import (
	"net/http"

	upload_service "github.com/brutalzinn/boberto-modpack-api/services/upload"
	"github.com/gin-gonic/gin"
)

// TODO: Show daniel how we will handle with files for all necessaries uploads

func CreateUploadRoute(router gin.IRouter) {
	router.POST("/upload/:id", func(ctx *gin.Context) {
		id := ctx.Query("id")
		form, _ := ctx.MultipartForm()
		if form == nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "You need upload a file or array of files"})
			return
		}
		files := form.File["files"]
		upload_service.SaveFiles(id, files)
	})

	// router.POST("/upload/create/:outputdir", func(ctx *gin.Context) {
	// 	outputdir := ctx.Query("outputdir")
	// 	uploadCache := upload_service.Create(outputdir)
	// 	url := common.GetUrl(ctx)
	// 	resourceData := rest.NewResData()
	// 	resourceData.Add(rest.Resource{
	// 		Object:    "upload",
	// 		Attribute: uploadCache,
	// 		Link: []rest.Link{
	// 			{
	// 				Rel:  "_self",
	// 				Href: fmt.Sprintf("%s/upload/%s", url, uploadCache.Id),
	// 			},
	// 		},
	// 	})
	// 	ctx.JSON(http.StatusOK, resourceData)
	// })
}
