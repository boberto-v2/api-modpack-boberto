package routes

import (
	"github.com/brutalzinn/boberto-modpack-api/middlewares"
	"github.com/gin-gonic/gin"
)

// TODO: Show daniel how we will use versions to route in agree with RFC 7231
// https://datatracker.ietf.org/doc/html/rfc7231

// for me this is not sooooo readable but is agree with gin documentation. we can do better after
func CreateRoutes(router gin.IRouter) {
	CreateAuthRoutes(router)
	user := router.Group("/user", middlewares.JWTMiddleware())
	{
		CreateApiKeyRoute(user)
	}
	game := router.Group("/game", middlewares.ApiKeyMiddleware())
	{
		server := game.Group("/server")
		{
			CreateClientRoute(server)
		}
		client := game.Group("/client")
		{
			CreateClientRoute(client)
		}
	}
	application := router.Group("/application", middlewares.ApiKeyMiddleware())
	{
		CreateUploadRoute(application)
		CreateEventRoute(application)
	}
}
