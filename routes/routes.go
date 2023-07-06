package routes

import (
	config "github.com/brutalzinn/boberto-modpack-api/configs"
	"github.com/gin-gonic/gin"
)

var cfg = config.GetConfig()

// TODO: Show daniel how we will use versions to route in agree with RFC 7231
// https://datatracker.ietf.org/doc/html/rfc7231
func CreateRoutes(router gin.IRouter) {

	user := router.Group("/user")
	{
		CreateUserRoutes(user)
	}
	//TODO: Explain why we will divide the routes by client, server and users routes.

	// server := router.Group("/server")
	// {
	// 	server.POST("/login", loginEndpoint)
	// }

	// client := router.Group("/client")
	// {
	// 	client.POST("/login", loginEndpoint)
	// }
}
