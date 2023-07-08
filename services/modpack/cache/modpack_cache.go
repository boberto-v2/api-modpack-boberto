package modpack_cache

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	config "github.com/brutalzinn/boberto-modpack-api/configs"
	modpack_cache_models "github.com/brutalzinn/boberto-modpack-api/services/modpack/cache/models"
	modpack_models "github.com/brutalzinn/boberto-modpack-api/services/modpack/models"
	"github.com/patrickmn/go-cache"
)

var cfg = config.GetConfig()
var modpack_cache = cache.New(5*time.Minute, 10*time.Minute)

func New() {
	modpack_cache.OnEvicted(func(s string, modpackExpire any) {
		oldModPack := modpackExpire.(*modpack_cache_models.ModPackCache)
		if oldModPack.Status == modpack_models.Aborted ||
			oldModPack.Status == modpack_models.Canceled ||
			oldModPack.Status == modpack_models.Error {
			finalPath := filepath.Join(cfg.API.PublicPath, oldModPack.NormalizedName)
			fmt.Printf("trying to remove %s", finalPath)
			os.RemoveAll(finalPath)
		}
		fmt.Printf("Delete cache for %s", s)
	})
}

func Replace(id string, modpack modpack_cache_models.ModPackCache) modpack_cache_models.ModPackCache {
	modpack_cache.SetDefault(id, modpack)
	return modpack
}

func GetById(id string) (modpack modpack_cache_models.ModPackCache, found bool) {
	if modpack, found := modpack_cache.Get(id); found {
		return modpack.(modpack_cache_models.ModPackCache), true
	}
	return modpack_cache_models.ModPackCache{}, false
}

func Create(modpack modpack_cache_models.ModPackCache) string {
	modpack_cache.Set(modpack.Id, modpack, cache.DefaultExpiration)
	return modpack.Id
}
