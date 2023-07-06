package routes

import (
	"net/http"

	modpack_models "github.com/brutalzinn/boberto-modpack-api/services/modpack/models"
	"github.com/gin-gonic/gin"
)

func CreateUserRoutes(router gin.IRouter) {
	router.POST("/login", func(ctx *gin.Context) {
		var modpack modpack_models.MinecraftModPack
		if err := ctx.ShouldBindJSON(&modpack); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	})
}
