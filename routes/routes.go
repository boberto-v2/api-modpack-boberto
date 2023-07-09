package routes

import (
	"github.com/gin-gonic/gin"
)

// TODO: Show daniel how we will use versions to route in agree with RFC 7231
// https://datatracker.ietf.org/doc/html/rfc7231
func CreateRoutes(router gin.IRouter) {
	//TODO: Explain why we will divide the routes by client, server and users routes.
	user := router.Group("/user")
	{
		CreateUserRoutes(user)
	}
	client := router.Group("/client")
	{
		CreateClientRoute(client)
	}

	application := router.Group("/application")
	{
		CreateUploadRoute(application)
	}
}
