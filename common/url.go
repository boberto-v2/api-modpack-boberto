package common

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

func GetCurrentHypermedia(ctx *gin.Context, resourceId string) []goeasyrest.Link {
	hostUrl := GetUrl(ctx)
	ctxLinks := ctx.Value("links").([]goeasyrest.Link)
	links := make([]goeasyrest.Link, 0)
	for _, item := range ctxLinks {
		item.Href = hostUrl + item.Href + resourceId
		links = append(links, item)
	}

	return links

}
