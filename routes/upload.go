package routes

import (
	"github.com/gin-gonic/gin"
)

// TODO: Show daniel how we will handle with files for all necessaries uploads

func CreateUploadRoute(router gin.IRouter) {
	router.POST("/upload/:token", func(ctx *gin.Context) {
		// token := ctx.Query("token")
	})
	router.POST("/upload/create", func(ctx *gin.Context) {

	})
}
