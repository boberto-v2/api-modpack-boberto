package routes

import (
	"github.com/brutalzinn/boberto-modpack-api/middlewares"
	application_routes "github.com/brutalzinn/boberto-modpack-api/routes/application"
	authorization_routes "github.com/brutalzinn/boberto-modpack-api/routes/authorization"
	game_routes "github.com/brutalzinn/boberto-modpack-api/routes/game"
	"github.com/gin-gonic/gin"
)

// TODO: Show daniel how we will use versions to route in agree with RFC 7231
// https://datatracker.ietf.org/doc/html/rfc7231

// for me this is not sooooo readable but is agree with gin documentation. we can do better after
func CreateRoutes(router gin.IRouter) {
	authorization := router.Group("/auth")
	{
		authorization_routes.CreateAuthRoutes(authorization)
	}
	game := router.Group("/game", middlewares.JWTMiddleware(), middlewares.ApiKeyMiddleware())
	{
		server := game.Group("/server")
		{
			game_routes.CreateClientRoute(server)
		}
		client := game.Group("/client")
		{
			game_routes.CreateServerRoute(client)
		}
	}
	application := router.Group("/application", middlewares.JWTMiddleware(), middlewares.ApiKeyMiddleware())
	{
		application_routes.CreateUploadRoute(application)
		application_routes.CreateEventRoute(application)
		application_routes.CreateApiKeyRoute(application)
	}
}
