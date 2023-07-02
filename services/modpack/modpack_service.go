package modpack_service

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	config "github.com/brutalzinn/boberto-modpack-api/configs"
	"github.com/brutalzinn/boberto-modpack-api/models"
	file_service "github.com/brutalzinn/boberto-modpack-api/services/file"
)

func CreateModPackFile(modpack *models.Modpack, environment models.MinecraftEnvironment) []models.ModPackFile {
	config := config.GetConfig()
	modpackPath := filepath.Join(config.PublicPath, modpack.NormalizedName, environment.GetFolderName())
	modpackFiles := []models.ModPackFile{}
	err := filepath.Walk(modpackPath,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				return nil
			}
			checksum, _ := file_service.GetCRC32(path)
			relativePath := strings.ReplaceAll(path, modpackPath+string(os.PathSeparator), "")
			fileType := GetType(relativePath)
			modpackFile := models.ModPackFile{
				Name:        info.Name(),
				Size:        info.Size(),
				Path:        relativePath,
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

func GetType(file string) models.MinecraftFileType {
	parts := strings.Split(file, string(os.PathSeparator))
	switch parts[0] {
	case "mods":
		return models.Mod
	case "data":
		return models.Data
	case "saves":
		return models.World
	case "config":
		return models.Config
	case "natives":
		return models.Library
	case "shaderpacks":
		return models.Texture
	default:
		return models.Other
	}
}
