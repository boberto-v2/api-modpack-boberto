package hypermedia

import (
	goeasyrest "github.com/brutalzinn/go-easy-rest"
	"github.com/gin-gonic/gin"
)

func GetUrl(ctx *gin.Context) string {
	scheme := "http://"
	if ctx.Request.TLS != nil {
		scheme = "https://"
	}
	url := scheme + ctx.Request.Host
	return url
}

func GetSocketUrl(ctx *gin.Context) string {
	scheme := "ws://"
	if ctx.Request.TLS != nil {
		scheme = "wss://"
	}
	url := scheme + ctx.Request.Host
	return url
}

func getUrlContext(context *gin.Context) []goeasyrest.Link {
	ctxLinks := context.Value(CTX_LINK_KEY).([]goeasyrest.Link)
	return ctxLinks
}
func addUrlContext(context *gin.Context, link goeasyrest.Link) {
	ctxLinks := context.Value(CTX_LINK_KEY).([]goeasyrest.Link)
	ctxLinks = append(ctxLinks, link)
	context.Set(CTX_LINK_KEY, ctxLinks)
}
