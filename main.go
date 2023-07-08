package main

import (
	"fmt"
	"net/http"

	config "github.com/brutalzinn/boberto-modpack-api/configs"
	"github.com/brutalzinn/boberto-modpack-api/routes"
	modpack_cache "github.com/brutalzinn/boberto-modpack-api/services/modpack/cache"
	"github.com/gin-gonic/gin"
)

//TODO: Explain the goals of this project to Daniel and explain the break changes about Boberto with C# and Boberto with GO
//Because now we can undestand what goes wrong with Boberto project at past

// 1. Explain about Gin and draw the goals at trello.
// 2. Create a draw.io fluxogram with the schemas of resource.
// 3. Explain the arch limits to determine the restful at this API.
// 4. Choose a primary security level to be used as authentication
// 5. Choose a secondary authentication ( API KEY OR OAUTH ) to be used by the modpack manager
// 6. Choose web front Engine ( I think its to be more appreciate if uses nextjs)
// 7. Configure the Daniel ambient to our first CODE PAIR TRULLY PROGRAMMING!!!!!! :) ( NO MORE THAN THREE YEARS TO WE DONE THIS. )
// 8. Explain how we apply rest at the routes and how he will consume the hypermedia routes at nextjs backend processor.

func main() {
	err := config.Load()
	if err != nil {
		panic(err)
	}
	config := config.GetConfig()
	router := gin.Default()
	router.GET("/healthchecker", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	modpack_cache.New()
	routes.CreateRoutes(router)
	router.Run(fmt.Sprintf(":%s", config.API.Port))
}
