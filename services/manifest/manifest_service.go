package manifest_service

import (
	"fmt"
	"path/filepath"
	"time"

	config "github.com/brutalzinn/boberto-modpack-api/configs"
	"github.com/brutalzinn/boberto-modpack-api/models"
	models_manifest "github.com/brutalzinn/boberto-modpack-api/models/manifest"
	json_service "github.com/brutalzinn/boberto-modpack-api/services/json"
)

func NewModPackManifest(modpack models.Modpack) {
	config := config.GetConfig()
	modPackManifestFile := filepath.Join(config.PublicPath, modpack.NormalizedName, "manifest.json")
	clientManifestUrl := fmt.Sprintf("%s/%s/%s/%s", modpack.ManifestUrl, modpack.NormalizedName, models.Client.GetFolderName(), "manifest.json")
	serverManifestUrl := fmt.Sprintf("%s/%s/%s/%s", modpack.ManifestUrl, modpack.NormalizedName, models.Server.GetFolderName(), "manifest.json")
	manifest := models_manifest.Manifest{
		Name:              modpack.Name,
		Author:            modpack.Author,
		Visible:           true,
		Version:           modpack.Version,
		ClientManifestUrl: clientManifestUrl,
		ServerManifestUrl: serverManifestUrl,
	}
	jsonManifest := json_service.JsonWritter{
		Path: modPackManifestFile,
		Data: manifest,
	}
	jsonManifest.CreateFile()
}

func NewManifest(modpack models.Modpack, files []models.ModPackFile, environment models.MinecraftEnvironment) {
	config := config.GetConfig()
	maifestFile := filepath.Join(config.PublicPath, modpack.NormalizedName, environment.GetFolderName(), "manifest.json")
	manifest := models_manifest.ModPackFileManifest{
		CreateAt: time.Now(),
		Files:    files,
	}
	jsonManifest := json_service.JsonWritter{
		Path: maifestFile,
		Data: manifest,
	}
	jsonManifest.CreateFile()
}
