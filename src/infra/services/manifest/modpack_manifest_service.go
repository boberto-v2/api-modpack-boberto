package manifest_service

import (
	"path/filepath"
	"time"

	"github.com/brutalzinn/boberto-modpack-api/src/src/src/common"
	config "github.com/brutalzinn/boberto-modpack-api/src/src/src/configs"
	manifest_models "github.com/brutalzinn/boberto-modpack-api/src/src/src/infra/services/manifest/models"
	modpack_models "github.com/brutalzinn/boberto-modpack-api/src/src/src/infra/services/modpack/models"
)

func WriteModPackManifestFiles(
	modpack modpack_models.MinecraftModPack,
	files []manifest_models.ManifestFile,
	environment modpack_models.MinecraftEnvironment) {
	cfg := config.GetConfig()
	nameNormalize := common.NormalizeString(modpack.Name)
	manifestFile := filepath.Join(cfg.ModPacks.PublicPath, nameNormalize, environment.GetFolderName(), cfg.ModPacks.ManifestName)
	manifest := manifest_models.ManifestFiles{
		CreateAt: time.Now(),
		Files:    files,
	}

	WriteManifestJsonFiles(manifestFile, manifest)
}

func ReadModPackManifestFiles(modpack modpack_models.MinecraftModPack, environment modpack_models.MinecraftEnvironment) manifest_models.ManifestFiles {
	cfg := config.GetConfig()
	nameNormalize := common.NormalizeString(modpack.Name)
	manifestFile := filepath.Join(cfg.ModPacks.PublicPath,
		nameNormalize,
		environment.GetFolderName(),
		cfg.ModPacks.ManifestName)

	manifest := ReadManifestJsonFiles(manifestFile)
	return manifest
}
