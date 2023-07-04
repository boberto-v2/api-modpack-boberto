package modpack_service

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	file_service "github.com/brutalzinn/boberto-modpack-api/services/file"
	modpack_cache_models "github.com/brutalzinn/boberto-modpack-api/services/modpack/cache/models"
	manifest_models "github.com/brutalzinn/boberto-modpack-api/services/modpack/manifest/models"
	modpack_models "github.com/brutalzinn/boberto-modpack-api/services/modpack/models"
)

func CreateModPackFilesManifest(modpackCache modpack_cache_models.ModPackCache,
	environment modpack_models.MinecraftEnvironment) []manifest_models.ManifestFile {

	modpackPath := filepath.Join(cfg.PublicPath,
		modpackCache.NormalizedName,
		environment.GetFolderName())

	modpackFiles := []manifest_models.ManifestFile{}
	err := filepath.Walk(modpackPath,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				return nil
			}
			relativePath := strings.ReplaceAll(path, modpackPath+string(os.PathSeparator), "")
			checksum, _ := file_service.GetCRC32(path)
			fileType := GetType(relativePath)
			modpackFile := manifest_models.ManifestFile{
				Name:        info.Name(),
				Size:        info.Size(),
				Path:        relativePath,
				Url:         "",
				Checksum:    checksum,
				Environment: environment,
				Type:        fileType,
			}
			modpackFiles = append(modpackFiles, modpackFile)
			fmt.Println(path, info.Size())
			return nil
		})
	if err != nil {
		log.Println(err)
	}
	return modpackFiles
}

func GetType(file string) modpack_models.MinecraftFileType {
	parts := strings.Split(file, string(os.PathSeparator))
	switch parts[0] {
	case "mods":
		return modpack_models.Mod
	case "data":
		return modpack_models.Data
	case "saves":
		return modpack_models.World
	case "config":
		return modpack_models.Config
	case "natives":
		return modpack_models.Library
	case "shaderpacks":
		return modpack_models.Texture
	default:
		return modpack_models.Other
	}
}
