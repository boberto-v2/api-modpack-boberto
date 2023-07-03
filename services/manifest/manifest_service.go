package manifest_service

import (
	"fmt"
	"path/filepath"
	"time"

	config "github.com/brutalzinn/boberto-modpack-api/configs"
	"github.com/brutalzinn/boberto-modpack-api/models"
	models_manifest "github.com/brutalzinn/boberto-modpack-api/models/manifest"
)

func NewModPackManifest(modpack models.Modpack) {
	config := config.GetConfig()
	modPackManifestFile := filepath.Join(config.PublicPath, modpack.NormalizedName, "manifest.json")
	clientManifestUrl := fmt.Sprintf("%s/%s/%s/%s", modpack.ManifestUrl, modpack.NormalizedName, models.Client.GetFolderName(), "manifest.json")
	serverManifestUrl := fmt.Sprintf("%s/%s/%s/%s", modpack.ManifestUrl, modpack.NormalizedName, models.Server.GetFolderName(), "manifest.json")
	manifest := models_manifest.Manifest{
		Id:                modpack.Id,
		Name:              modpack.Name,
		Author:            modpack.Author,
		Visible:           true,
		Version:           modpack.Version,
		ClientManifestUrl: clientManifestUrl,
		ServerManifestUrl: serverManifestUrl,
	}
	CreateModPackManifest(modPackManifestFile, manifest)
}

func NewManifest(modpack models.Modpack,
	files []models.ModPackFile,
	environment models.MinecraftEnvironment) {
	config := config.GetConfig()
	manifestFile := filepath.Join(config.PublicPath, modpack.NormalizedName, environment.GetFolderName(), "manifest.json")
	manifest := models_manifest.ModPackFileManifest{
		CreateAt: time.Now(),
		Files:    files,
	}

	CreateManifestFile(manifestFile, manifest)
}

func ReadManifest(modpack models.Modpack, environment models.MinecraftEnvironment) models_manifest.ModPackFileManifest {
	config := config.GetConfig()
	manifestFile := filepath.Join(config.PublicPath, modpack.NormalizedName, environment.GetFolderName(), "manifest.json")
	manifest := ReadManifestFile(manifestFile)
	return manifest
}

func CompareManifest(oldManifest models_manifest.ModPackFileManifest,
	newManifest models_manifest.ModPackFileManifest) (newFiles []models.ModPackFile, deleteFiles []models.ModPackFile) {
	newFiles = []models.ModPackFile{}
	for _, oldFile := range oldManifest.Files {
		found := false
		for _, newFile := range newManifest.Files {
			if oldFile.Name != newFile.Name {
				continue
			}
			if oldFile.Checksum != newFile.Checksum {
				newFiles = append(newFiles, newFile)
			}
			found = true
		}
		if !found {
			deleteFiles = append(deleteFiles, oldFile)
		}
	}
	return newFiles, deleteFiles
}
