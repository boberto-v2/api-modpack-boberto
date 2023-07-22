package modpack_service

import (
	"encoding/json"
	"log"
	"path/filepath"

	"github.com/brutalzinn/boberto-modpack-api/common"
	config "github.com/brutalzinn/boberto-modpack-api/configs"
	ftp_service "github.com/brutalzinn/boberto-modpack-api/infra/services/ftp"
	ftp_models "github.com/brutalzinn/boberto-modpack-api/infra/services/ftp/models"
	manifest_service "github.com/brutalzinn/boberto-modpack-api/infra/services/manifest"
	manifest_compare "github.com/brutalzinn/boberto-modpack-api/infra/services/manifest/comparer"
	manifest_models "github.com/brutalzinn/boberto-modpack-api/infra/services/manifest/models"
	modpack_models "github.com/brutalzinn/boberto-modpack-api/infra/services/modpack/models"
	"github.com/jlaffaye/ftp"
)

func UploadManifest(modPack modpack_models.MinecraftModPack, ftpClient *ftp.ServerConn) error {
	cfg := config.GetConfig()
	nameNormalized := common.NormalizeString(modPack.Name)
	modPackPath := filepath.Join(cfg.ModPacks.PublicPath,
		nameNormalized,
		modpack_models.Client.GetFolderName())

	err := ftp_service.UploadFileFTP(
		cfg.ModPacks.ManifestName,
		modPackPath,
		ftpClient)

	log.Printf("Manifest uploaded with sucess")
	if err != nil {
		log.Printf("Manifest uploaded wrong")
		return err
	}
	return nil
}

func GetModPackManifest(ftpClient *ftp.ServerConn) (*manifest_models.ManifestFiles, error) {
	cfg := config.GetConfig()
	ftpManifestJson, err := ftp_service.ReadFileFTP(cfg.ModPacks.ManifestName, ftpClient)
	if err != nil {
		return nil, err
	}
	var ftpManifest manifest_models.ManifestFiles
	err = json.Unmarshal(ftpManifestJson, &ftpManifest)
	if err != nil {
		log.Printf(err.Error())
		return nil, err
	}
	return &ftpManifest, nil
}

func UploadOrDeleteSyncModPack(
	oldManifest manifest_models.ManifestFiles,
	newManifest manifest_models.ManifestFiles,
	modPackPath string,
	ftpClient *ftp.ServerConn) error {
	log.Println("Start sync files between manifest at ftp server and local manifest")
	manifestComparer := manifest_compare.New(oldManifest, newManifest)
	result := manifestComparer.Compare()
	for _, newFile := range result.ToUpload {
		filePath := filepath.Join(modPackPath, newFile.Path)
		err := ftp_service.UploadFileFTP(filePath, modPackPath, ftpClient)
		if err != nil {
			return err
		}
	}
	for _, deleteFile := range result.ToDelete {
		err := ftp_service.DeleteFileFTP(deleteFile.Path, ftpClient)
		if err != nil {
			return err
		}
	}
	log.Printf("upload cycle is finish here")
	return nil
}

func UploadClient(modpack modpack_models.MinecraftModPack,
	ftpCredentials ftp_models.Ftp) error {
	cfg := config.GetConfig()
	nameNormalized := common.NormalizeString(modpack.Name)
	modPackPath := filepath.Join(
		cfg.ModPacks.PublicPath,
		nameNormalized,
		modpack_models.Client.GetFolderName())

	ftpClient, err := ftp_service.OpenFtpConnection(
		ftpCredentials.Directory,
		ftpCredentials.Address,
		ftpCredentials.User,
		ftpCredentials.Password,
	)
	if err != nil {
		log.Println("Somethings goes wrong at ftp connecion. Close everyting", err)
		return err
	}
	manifest := manifest_service.ReadModPackManifestFiles(modpack, modpack_models.Client)
	oldManifest, err := GetModPackManifest(ftpClient)
	if err != nil {
		for _, item := range manifest.Files {
			ftp_service.UploadFileFTP(item.Path, modPackPath, ftpClient)
		}
		log.Printf("Uploading all client modpacks files to ftp server")
		err := UploadManifest(modpack, ftpClient)
		if err != nil {
			log.Println("Somethings goes wrong at ftp connecion. Close everyting", err)
			return err
		}
		return err
	}
	UploadOrDeleteSyncModPack(*oldManifest, manifest, modPackPath, ftpClient)
	err = UploadManifest(modpack, ftpClient)
	if err != nil {
		log.Println("Somethings goes wrong at ftp connecion. Close everyting", err)
		return err
	}

	return nil
}
