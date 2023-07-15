package rest_object

import (
	"fmt"

	"github.com/brutalzinn/boberto-modpack-api/common"
	upload_cache "github.com/brutalzinn/boberto-modpack-api/services/upload/cache"
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
		Link: []rest.Link{
			{
				Rel:    UPLOAD_FILE,
				Href:   fmt.Sprintf("%s/application/upload/%s", common.GetUrl(restObject.ctx), uploadCache.Id),
				Method: "POST",
			},
		},
	}
	return restObject
}
