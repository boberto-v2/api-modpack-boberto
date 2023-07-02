package modpack_cache

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	config "github.com/brutalzinn/boberto-modpack-api/configs"
	"github.com/brutalzinn/boberto-modpack-api/models"
	"github.com/patrickmn/go-cache"
)

var modpack_cache = cache.New(5*time.Minute, 10*time.Minute)

func New() {
	config := config.GetConfig()
	modpack_cache.OnEvicted(func(s string, modpackExpire any) {
		oldModPack := modpackExpire.(*models.Modpack)
		if oldModPack.Status == models.Aborted ||
			oldModPack.Status == models.Canceled ||
			oldModPack.Status == models.Error {
			finalPath := filepath.Join(config.PublicPath, oldModPack.NormalizedName)
			fmt.Printf("trying to remove %s", finalPath)
			os.RemoveAll(finalPath)
		}
		fmt.Printf("Delete cache for %s", s)
	})
}

func SetStatus(id string, status models.ModPackStatus) models.Modpack {
	modpack, _ := GetModpackCacheById(id)
	modpack.Status = status
	modpack_cache.SetDefault(id, modpack)
	return modpack
}

func GetStatus(id string) models.ModPackStatus {
	modpack, _ := GetModpackCacheById(id)
	status := modpack.Status
	return status
}

func GetModpackCacheById(id string) (modpack models.Modpack, found bool) {
	if modpack, found := modpack_cache.Get(id); found {
		return modpack.(models.Modpack), true
	}
	return models.Modpack{}, false
}

func CreateModpackCache(modpack models.Modpack) string {
	modpack_cache.Set(modpack.Id, modpack, cache.DefaultExpiration)
	return modpack.Id
}
