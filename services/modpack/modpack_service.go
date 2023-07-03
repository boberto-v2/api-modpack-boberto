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
	"github.com/jlaffaye/ftp"
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
func GetFTPManifest(ftpClient *ftp.ServerConn) (*models_manifest.ModPackFileManifest, error) {
	ftpManifestJson, err := file_service.ReadFileFTP("manifest.json", ftpClient)
	if err != nil {
		return nil, err
	}
	var ftpManifest models_manifest.ModPackFileManifest
	err = json.Unmarshal(ftpManifestJson, &ftpManifest)
	if err != nil {
		log.Printf(err.Error())
		return nil, err
	}
	return &ftpManifest, nil
}

func UploadOrDeleteSyncModPack(
	oldManifest models_manifest.ModPackFileManifest,
	newManifest models_manifest.ModPackFileManifest,
	modPackPath string,
	ftpClient *ftp.ServerConn) error {
	log.Println("Start sync files between manifest at ftp server and local manifest")
	newFiles, deleteFiles := manifest_service.CompareManifest(oldManifest, newManifest)
	for _, newFile := range newFiles {
		filePath := filepath.Join(modPackPath, newFile.Path)
		err := file_service.UploadFileFTP(filePath, modPackPath, ftpClient)
		log.Printf("Upload file ftp %v", filePath)
		if err != nil {
			log.Printf(err.Error())
			return err
		}
	}
	for _, deleteFile := range deleteFiles {
		log.Printf("Delete file ftp %v", deleteFile)
		err := file_service.DeleteFileFromFTP(deleteFile.Path, ftpClient)
		if err != nil {
			log.Printf(err.Error())
			return err
		}
	}
	err := file_service.UploadFileFTP("manifest.json", modPackPath, ftpClient)
	if err != nil {
		log.Printf("new manifest uploaded")
		return err
	}
	//now we can start the token system to transform this in a true api.
	//this kind of upload files to server or client files need be done by a token to client.
	//this is the moment that minecraft launcher needs to call to read the list of download files
	log.Printf("upload cycle is finish here")
	return nil
}

func UploadServer(modpack models.Modpack,
	ftpCredentials models.ModPackFtp) error {
	config := config.GetConfig()
	modPackPath := filepath.Join(config.PublicPath, modpack.NormalizedName, models.Server.GetFolderName())

	ftpClient, err := file_service.OpenFtpConnection(
		ftpCredentials.ServerFtp.Directory,
		ftpCredentials.ServerFtp.Address,
		ftpCredentials.ServerFtp.User,
		ftpCredentials.ServerFtp.Password,
	)
	if err != nil {
		log.Println("Somethings goes wrong at ftp connecion. Put this on queue", err)
		return err
	}
	manifest := manifest_service.ReadManifest(modpack, models.Server)
	oldManifest, err := GetFTPManifest(ftpClient)
	if err != nil {
		for _, item := range manifest.Files {
			file_service.UploadFileFTPWithDirectories(item.Path, modPackPath, ftpClient)
		}
		log.Printf("Uploading all server modpacks files to ftp server")
		log.Printf("Uploading server manifest file")
		err = file_service.UploadFileFTP("manifest.json", modPackPath, ftpClient)
		if err != nil {
			log.Printf("Manifest server uploaded")
			return err
		}
		return nil
	}
	if &oldManifest != nil {
		UploadOrDeleteSyncModPack(*oldManifest, manifest, modPackPath, ftpClient)
	}
	modpack_cache.SetStatus(modpack.Id, models.Finish)
	return nil
}

func UploadClient(modpack models.Modpack,
	ftpCredentials models.ModPackFtp) error {
	config := config.GetConfig()
	modPackPath := filepath.Join(config.PublicPath, modpack.NormalizedName, models.Client.GetFolderName())

	ftpClient, err := file_service.OpenFtpConnection(
		ftpCredentials.ClientFtp.Directory,
		ftpCredentials.ClientFtp.Address,
		ftpCredentials.ClientFtp.User,
		ftpCredentials.ClientFtp.Password,
	)
	if err != nil {
		log.Println("Somethings goes wrong at ftp connecion. Put this on queue", err)
		return err
	}
	manifest := manifest_service.ReadManifest(modpack, models.Client)
	oldManifest, err := GetFTPManifest(ftpClient)
	if err != nil {
		for _, item := range manifest.Files {
			file_service.UploadFileFTPWithDirectories(item.Path, modPackPath, ftpClient)
		}
		log.Printf("Uploading all client modpacks files to ftp server")
		log.Printf("Uploading client manifest file")
		err = file_service.UploadFileFTP("manifest.json", modPackPath, ftpClient)
		if err != nil {
			log.Printf("Manifest server uploaded")
			return err
		}
		return nil
	}
	if &oldManifest != nil {
		UploadOrDeleteSyncModPack(*oldManifest, manifest, modPackPath, ftpClient)
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
