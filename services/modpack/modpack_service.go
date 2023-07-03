package modpack_service

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	config "github.com/brutalzinn/boberto-modpack-api/configs"
	"github.com/brutalzinn/boberto-modpack-api/models"
	models_manifest "github.com/brutalzinn/boberto-modpack-api/models/manifest"
	modpack_cache "github.com/brutalzinn/boberto-modpack-api/services/cache"
	file_service "github.com/brutalzinn/boberto-modpack-api/services/file"
	manifest_service "github.com/brutalzinn/boberto-modpack-api/services/manifest"
)

func CreateModPackFile(modpack models.Modpack, environment models.MinecraftEnvironment) []models.ModPackFile {
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
			relativePath := strings.ReplaceAll(path, modpackPath+string(os.PathSeparator), "")
			checksum, _ := file_service.GetCRC32(path)
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

func UploadServer(modpack models.Modpack,
	ftpCredentials models.ModPackFtp) error {
	config := config.GetConfig()
	modpackPath := filepath.Join(config.PublicPath, modpack.NormalizedName, models.Server.GetFolderName())

	ftpClient, err := file_service.OpenFtpConnection(
		ftpCredentials.ServerFtp.Directory,
		ftpCredentials.ServerFtp.Address,
		ftpCredentials.ServerFtp.User,
		ftpCredentials.ServerFtp.Password,
	)
	if err != nil {
		log.Printf(err.Error())
		return err
	}
	manifest := manifest_service.ReadManifest(modpack, models.Server)
	ftpManifest, err := file_service.ReadFileFTP("manifest.json", ftpClient)
	if err != nil {
		var files []string
		for _, item := range manifest.Files {
			files = append(files, item.Path)
		}
		file_service.UploadFilesToFTP(files, modpackPath, ftpClient)
		log.Printf("Uploading all init files to ftp server")
		return nil
	}
	var oldManifest models_manifest.ModPackFileManifest
	err = json.Unmarshal(ftpManifest, &oldManifest)
	if err != nil {
		log.Printf(err.Error())
		return err
	}
	err = file_service.UploadFileFtp("manifest.json", modpackPath, ftpClient)
	if err != nil {
		log.Printf("Manifest server uploaded")
		return err
	}
	newFiles, deleteFiles := manifest_service.CompareManifest(oldManifest, manifest)
	log.Printf("new server files %v delete %v files", len(newFiles), len(deleteFiles))
	for _, newFile := range newFiles {
		filePath := filepath.Join(modpackPath, newFile.Path)
		err := file_service.UploadFileFtp(filePath, modpackPath, ftpClient)
		if err != nil {
			log.Printf(err.Error())
			return err
		}
	}
	for _, deleteFile := range deleteFiles {
		err := file_service.DeleteFileFTP(deleteFile.Path, ftpClient)
		if err != nil {
			log.Printf(err.Error())
			return err
		}
	}
	modpack_cache.SetStatus(modpack.Id, models.Finish)
	return nil
}

func UploadClient(modpack models.Modpack,
	ftpCredentials models.ModPackFtp) error {
	config := config.GetConfig()
	modpackPath := filepath.Join(config.PublicPath, modpack.NormalizedName, models.Client.GetFolderName())

	ftpClient, err := file_service.OpenFtpConnection(
		ftpCredentials.ClientFtp.Directory,
		ftpCredentials.ClientFtp.Address,
		ftpCredentials.ClientFtp.User,
		ftpCredentials.ClientFtp.Password,
	)
	if err != nil {
		log.Printf(err.Error())
		return err
	}
	manifest := manifest_service.ReadManifest(modpack, models.Client)
	ftpManifest, err := file_service.ReadFileFTP("manifest.json", ftpClient)
	if err != nil {
		var files []string
		for _, item := range manifest.Files {
			files = append(files, item.Path)
		}
		file_service.UploadFilesToFTP(files, modpackPath, ftpClient)
		log.Printf("Uploading all init files to ftp server")
		return nil
	}
	var oldManifest models_manifest.ModPackFileManifest
	err = json.Unmarshal(ftpManifest, &oldManifest)
	if err != nil {
		log.Printf(err.Error())
		return err
	}
	err = file_service.UploadFileFtp("manifest.json", modpackPath, ftpClient)
	if err != nil {
		log.Printf("Manifest client uploaded")
		return err
	}
	newFiles, deleteFiles := manifest_service.CompareManifest(oldManifest, manifest)
	log.Printf("new client files %v delete %v files", len(newFiles), len(deleteFiles))
	for _, newFile := range newFiles {
		filePath := filepath.Join(modpackPath, newFile.Path)
		err := file_service.UploadFileFtp(filePath, modpackPath, ftpClient)
		if err != nil {
			log.Printf(err.Error())
			return err
		}
	}
	for _, deleteFile := range deleteFiles {
		err := file_service.DeleteFileFTP(deleteFile.Path, ftpClient)
		if err != nil {
			log.Printf(err.Error())
			return err
		}
	}
	modpack_cache.SetStatus(modpack.Id, models.Finish)
	return nil
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
