package rest_object

import (
	"fmt"

	upload_cache_models "github.com/brutalzinn/boberto-modpack-api/services/upload/cache/models"
	rest "github.com/brutalzinn/go-easy-rest"
)

const (
	UPLOAD_FILE = "upload_file"
)

type FileObject struct {
	Id string `json:"id"`
}

func (restObject RestObject) CreateFileObject(uploadCache *upload_cache_models.UploadCache) RestObject {
	restObject.Resource = rest.Resource{
		Object: FILE_OBJECT,
		Attribute: FileObject{
			Id: uploadCache.Id,
		},
		Link: []rest.Link{
			{
				Rel:    "upload_file",
				Href:   fmt.Sprintf("/application/upload/%s", uploadCache.Id),
				Method: "POST",
			},
		},
	}
	return restObject
}
