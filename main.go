package main

import (
	"fmt"

	config "github.com/brutalzinn/go-multiple-file/configs"
	"github.com/brutalzinn/go-multiple-file/routes"
	modpack_cache "github.com/brutalzinn/go-multiple-file/services/cache"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	err := config.Load()
	if err != nil {
		panic(err)
	}
	config := config.GetConfig()
	modpack_cache.New()
	router.GET("/healthchecker", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"status": "success", "message": "How to Upload Single and Multiple Files in Golang"})
	})
	routes.CreateModPackRoute(router)
	router.Run(fmt.Sprintf(":%s", config.Port))
}
