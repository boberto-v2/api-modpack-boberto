package rest_object

import (
	upload_cache_models "github.com/brutalzinn/boberto-modpack-api/services/upload/cache/models"
	rest "github.com/brutalzinn/go-easy-rest"
)

const (
	UPLOAD_FILE = "upload_file"
)

type FileObject struct {
	Id string `json:"id"`
}

func (restObject RestObject) CreateFileObject(data *upload_cache_models.UploadCache) RestObject {
	restObject.Resource = rest.Resource{
		Object: FILE_OBJECT,
		Attribute: FileObject{
			Id: data.Id,
		},
		Link: restObject.Link,
	}
	restObject.Create()
	return restObject
}
