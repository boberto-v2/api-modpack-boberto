package rest_object

import (
	"github.com/brutalzinn/boberto-modpack-api/src/src/src/common/hypermedia"
	upload_cache "github.com/brutalzinn/boberto-modpack-api/src/src/src/infra/services/upload/cache"
	rest "github.com/brutalzinn/go-easy-rest"
)

const (
	UPLOAD_FILE = "upload_file"
)

type FileObject struct {
	Id string `json:"id"`
}

func (restObject *RestObject) CreateUploadFileObject(uploadCache *upload_cache.UploadCache) *RestObject {
	restObject.Resource = rest.Resource{
		Object: FILE_OBJECT,
		Attribute: FileObject{
			Id: uploadCache.Id,
		},
	}
	hyperlink := hypermedia.New(restObject.ctx)
	hyperlink.SetOptions(hypermedia.HyperOptions{UrlType: hypermedia.HTTP, Id: uploadCache.Id})
	hyperlink.AddHyperLink(rest.Link{Rel: UPLOAD_FILE, Href: "/application/upload/", Method: "POST"})
	return restObject
}
