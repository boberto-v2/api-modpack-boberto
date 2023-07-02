package manifest_service

import (
	"fmt"
	"time"

	config "github.com/brutalzinn/boberto-modpack-api/configs"
	"github.com/brutalzinn/boberto-modpack-api/models"
	models_manifest "github.com/brutalzinn/boberto-modpack-api/models/manifest"
	json_service "github.com/brutalzinn/boberto-modpack-api/services/json"
)

func NewModPackManifest(modpack models.Modpack) {
	config := config.GetConfig()
	modpackPath := fmt.Sprintf("%s/%s/%s", config.PublicPath, modpack.NormalizedName, "manifest.json")
	clientManifestUrl := fmt.Sprintf("http://localhost:%s/%s/%s/%s", config.Port, modpack.NormalizedName, models.Client.GetFolderName(), "manifest.json")
	serverManifestUrl := fmt.Sprintf("http://localhost:%s/%s/%s/%s", config.Port, modpack.NormalizedName, models.Server.GetFolderName(), "manifest.json")
	manifest := models_manifest.Manifest{
		Name:              modpack.Name,
		Author:            modpack.Author,
		Visible:           true,
		Version:           modpack.Version,
		ClientManifestUrl: clientManifestUrl,
		ServerManifestUrl: serverManifestUrl,
	}
	jsonManifest := json_service.JsonWritter{
		Path: modpackPath,
		Data: manifest,
	}
	jsonManifest.CreateFile()
}

func NewManifest(modpack *models.Modpack, files []models.ModPackFile, environment models.MinecraftEnvironment) {
	config := config.GetConfig()
	modpackPath := fmt.Sprintf("%s/%s/%s/%s", config.PublicPath, modpack.NormalizedName, environment.GetFolderName(), "manifest.json")
	manifest := models_manifest.ModPackFileManifest{
		CreateAt: time.Now(),
		Files:    files,
	}
	jsonManifest := json_service.JsonWritter{
		Path: modpackPath,
		Data: manifest,
	}
	jsonManifest.CreateFile()
}
