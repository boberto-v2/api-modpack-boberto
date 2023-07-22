package game_routes

import (
	"github.com/gin-gonic/gin"
)

func CreateServerRoute(router gin.IRouter) {

	router.GET("/modpack/:id", func(ctx *gin.Context) {

	})

	router.POST("/modpack/create", func(ctx *gin.Context) {

	})

	router.POST("/modpack/finish/:id", func(ctx *gin.Context) {

	})
}
