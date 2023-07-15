package upload_cache

import (
	"time"

	"github.com/patrickmn/go-cache"
)

var upload_cache = cache.New(5*time.Minute, 10*time.Minute)

type UploadStatus int

type UploadCache struct {
	Id        string       `json:"id"`
	OutputDir string       `json:"output_dir"`
	Status    UploadStatus `json:"status"`
	CreateAt  time.Time    `json:"create_at"`
	ExpireAt  time.Time    `json:"expire_at"`
}

const (
	UPLOAD_CREATED   UploadStatus = 1
	UPLOAD_PENDING   UploadStatus = 2
	UPLOAD_COMPLETED UploadStatus = 3
	UPLOAD_CANCELED  UploadStatus = 4
)

func GetById(id string) (uploadCache UploadCache, found bool) {
	if uploadCache, found := upload_cache.Get(id); found {
		return uploadCache.(UploadCache), true
	}
	return UploadCache{}, false
}

func Create(uploadCache UploadCache) {
	upload_cache.SetDefault(uploadCache.Id, uploadCache)
}

func (uploadCache UploadCache) Save() {
	upload_cache.SetDefault(uploadCache.Id, uploadCache)
}
