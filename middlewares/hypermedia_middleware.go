/// lets do the version 2 of the hypermedia middleware https://github.com/brutalzinn/api-task-list/blob/main/src/middlewares/hypermedia/hypermedia_middleware.go

package middlewares

import (
	rest "github.com/brutalzinn/go-easy-rest"

	"github.com/gin-gonic/gin"
)

func (hypermedia Hypermedia) HypermediaMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Set("links", hypermedia.Links)
		ctx.Next()
	}
}
func New(hypermedia Hypermedia) *Hypermedia {
	opt := &Hypermedia{
		Links: hypermedia.Links,
	}
	return opt
}

type Hypermedia struct {
	Links []rest.Link `json:"links"`
}
