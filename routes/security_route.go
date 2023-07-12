// I will use https://pkg.go.dev/crypto/ecdh to communicate sensitive data between with third parties integration.
// i dont take responsability to store any kind of sensitive information. The client needs send to me and save this data at local with secure storage.

package routes

import (
	event_service "github.com/brutalzinn/boberto-modpack-api/services/event"
	"github.com/gin-gonic/gin"
)

func CreateSecurityRoute(router gin.IRouter) {
	router.GET("/event", func(ctx *gin.Context) {
		event_service.WebSocketHandler(ctx)
	})
}
