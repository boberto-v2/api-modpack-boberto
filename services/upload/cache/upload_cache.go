package upload_cache

import (
	"time"

	upload_cache_models "github.com/brutalzinn/boberto-modpack-api/services/upload/cache/models"
	"github.com/patrickmn/go-cache"
)

var upload_cache = cache.New(5*time.Minute, 10*time.Minute)

func GetById(id string) (uploadCache upload_cache_models.UploadCache, found bool) {
	if uploadCache, found := upload_cache.Get(id); found {
		return uploadCache.(upload_cache_models.UploadCache), true
	}
	return upload_cache_models.UploadCache{}, false
}

func Create(uploadCache upload_cache_models.UploadCache) {
	upload_cache.Set(uploadCache.Id, uploadCache, cache.DefaultExpiration)
}
