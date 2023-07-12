package application_routes

import (
	event_service "github.com/brutalzinn/boberto-modpack-api/services/event"
	"github.com/gin-gonic/gin"
)

func CreateEventRoute(router gin.IRouter) {
	router.GET("/event", func(ctx *gin.Context) {
		event_service.WebSocketHandler(ctx)
	})
}
