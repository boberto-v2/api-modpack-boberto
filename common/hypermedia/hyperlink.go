package hypermedia

import (
	goeasyrest "github.com/brutalzinn/go-easy-rest"
	"github.com/gin-gonic/gin"
)

const (
	CTX_LINK_KEY = "links"
)

type UrlType int

const (
	WEBSOCKET UrlType = 1
	HTTP      UrlType = 2
	FTP       UrlType = 3
)

type HyperLink struct {
	context *gin.Context
	Options *HyperOptions
}

type HyperOptions struct {
	UrlType UrlType
	Id      string
}

func New(ctx *gin.Context) *HyperLink {
	hyperLink := HyperLink{
		context: ctx,
		Options: &HyperOptions{
			UrlType: HTTP,
		},
	}
	return &hyperLink
}

func (hyperLink *HyperLink) SetOptions(options HyperOptions) {
	hyperLink.Options = &options
}

func (hyperLink *HyperLink) GetCurrentHyperLink() []goeasyrest.Link {
	ctxLinks := hyperLink.context.Value(CTX_LINK_KEY).([]goeasyrest.Link)
	links := make([]goeasyrest.Link, 0)
	for _, item := range ctxLinks {
		item.Href = item.Href + hyperLink.Options.Id
		links = append(links, item)
	}
	return links
}

func (hyperLink *HyperLink) AddHyperLink(link goeasyrest.Link) goeasyrest.Link {
	hostUrl := hyperLink.getCurrentUrl()
	newLink := goeasyrest.Link{
		Rel:    link.Rel,
		Href:   hostUrl + link.Href + hyperLink.Options.Id,
		Method: link.Method,
	}
	addUrlContext(hyperLink.context, newLink)
	return newLink
}

func (hyperLink *HyperLink) getCurrentUrl() string {
	hostUrl := ""
	switch hyperLink.Options.UrlType {
	case WEBSOCKET:
		return GetSocketUrl(hyperLink.context)
	case HTTP:
		return GetUrl(hyperLink.context)
	}
	return hostUrl
}
