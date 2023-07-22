package application_routes

import (
	event_service "github.com/brutalzinn/boberto-modpack-api/src/src/src/infra/services/event"
	"github.com/gin-gonic/gin"
)

func CreateEventRoute(router gin.IRouter) {
	router.GET("/event", func(ctx *gin.Context) {
		event_service.WebSocketHandler(ctx)
	})
}
