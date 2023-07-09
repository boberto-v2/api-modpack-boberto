package common

import "github.com/gin-gonic/gin"

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
