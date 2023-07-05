package manifest_service

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/brutalzinn/boberto-modpack-api/common"
	config "github.com/brutalzinn/boberto-modpack-api/configs"
	manifest_models "github.com/brutalzinn/boberto-modpack-api/services/modpack/manifest/models"
	modpack_models "github.com/brutalzinn/boberto-modpack-api/services/modpack/models"
)

var cfg = config.GetConfig()

func CreateModPackManifest(modpack manifest_models.ManifestModPack) {

	nameNormalize := common.NormalizeString(modpack.Name)
	modPackManifestFile := filepath.Join(cfg.API.PublicPath, nameNormalize, cfg.API.ManifestName)
	clientManifestUrl := fmt.Sprintf("%s/%s/%s/%s", modpack.ClientManifestUrl, nameNormalize, modpack_models.Client.GetFolderName(), cfg.API.ManifestName)
	serverManifestUrl := fmt.Sprintf("%s/%s/%s/%s", modpack.ServerManifestUrl, nameNormalize, modpack_models.Server.GetFolderName(), cfg.API.ManifestName)
	manifest := manifest_models.ManifestModPack{
		Id:                modpack.Id,
		Name:              modpack.Name,
		Author:            modpack.Author,
		Visible:           true,
		Version:           modpack.Version,
		ClientManifestUrl: clientManifestUrl,
		ServerManifestUrl: serverManifestUrl,
	}
	WriteModPackJsonManifest(modPackManifestFile, manifest)
}

func WriteModPackManifestFiles(
	modpack modpack_models.MinecraftModPack,
	files []manifest_models.ManifestFile,
	environment modpack_models.MinecraftEnvironment) {

	nameNormalize := common.NormalizeString(modpack.Name)
	manifestFile := filepath.Join(cfg.API.PublicPath, nameNormalize, environment.GetFolderName(), cfg.API.ManifestName)
	manifest := manifest_models.ManifestFiles{
		CreateAt: time.Now(),
		Files:    files,
	}

	WriteManifestJsonFiles(manifestFile, manifest)
}

func ReadModPackManifestFiles(modpack modpack_models.MinecraftModPack, environment modpack_models.MinecraftEnvironment) manifest_models.ManifestFiles {
	nameNormalize := common.NormalizeString(modpack.Name)

	manifestFile := filepath.Join(cfg.API.PublicPath,
		nameNormalize,
		environment.GetFolderName(),
		cfg.API.ManifestName)

	manifest := ReadManifestJsonFiles(manifestFile)
	return manifest
}
