package modpack_cache

import (
	"fmt"
	"os"
	"time"

	config "github.com/brutalzinn/go-multiple-file/configs"
	"github.com/brutalzinn/go-multiple-file/models"
	"github.com/patrickmn/go-cache"
)

var modpack_cache = cache.New(5*time.Minute, 10*time.Minute)

func New() {
	config := config.GetConfig()
	modpack_cache.OnEvicted(func(s string, modpackExpire any) {
		oldModPack := modpackExpire.(*models.Modpack)
		finalPath := fmt.Sprintf("%s/%s", config.PublicPath, oldModPack.NormalizedName)
		fmt.Printf("trying to remove %s", finalPath)
		os.Remove(finalPath)
	})
}

func GetModpackCacheById(id string) *models.Modpack {
	if modpack, found := modpack_cache.Get(id); found {
		return modpack.(*models.Modpack)
	}
	return nil
}

func CreateModpackCache(modpack *models.Modpack) string {
	modpack_cache.Set(modpack.Id, modpack, cache.DefaultExpiration)
	return modpack.Id
}
