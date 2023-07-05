package modpack_service

import (
	"encoding/json"
	"log"
	"path/filepath"

	"github.com/brutalzinn/boberto-modpack-api/common"
	config "github.com/brutalzinn/boberto-modpack-api/configs"
	ftp_service "github.com/brutalzinn/boberto-modpack-api/services/ftp"
	manifest_service "github.com/brutalzinn/boberto-modpack-api/services/modpack/manifest"
	manifest_compare "github.com/brutalzinn/boberto-modpack-api/services/modpack/manifest/comparer"
	manifest_models "github.com/brutalzinn/boberto-modpack-api/services/modpack/manifest/models"
	modpack_models "github.com/brutalzinn/boberto-modpack-api/services/modpack/models"
	"github.com/jlaffaye/ftp"
)

var cfg = config.GetConfig()

func UploadManifest(modPack modpack_models.MinecraftModPack, ftpClient *ftp.ServerConn) error {

	nameNormalized := common.NormalizeString(modPack.Name)

	modPackPath := filepath.Join(cfg.API.PublicPath,
		nameNormalized,
		modpack_models.Client.GetFolderName())

	err := ftp_service.UploadFileFTP(
		cfg.API.ManifestName,
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
	ftpManifestJson, err := ftp_service.ReadFileFTP(cfg.API.ManifestName, ftpClient)
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
			log.Printf(err.Error())
			return err
		}
	}
	for _, deleteFile := range result.ToDelete {
		err := ftp_service.DeleteFileFTP(deleteFile.Path, ftpClient)
		if err != nil {
			log.Printf(err.Error())
			return err
		}
	}

	//now we can start the token system to transform this in a true api.
	//this kind of upload files to server or client files need be done by a token to client.
	//this is the moment that minecraft launcher needs to call to read the list of download files
	log.Printf("upload cycle is finish here")
	return nil
}

func UploadServer(modpack modpack_models.MinecraftModPack,
	ftpCredentials modpack_models.ModPackFtp) error {

	normalizeName := common.NormalizeString(modpack.Name)
	modPackPath := filepath.Join(
		cfg.API.PublicPath,
		normalizeName,
		modpack_models.Server.GetFolderName())

	ftpClient, err := ftp_service.OpenFtpConnection(
		ftpCredentials.ServerFtp.Directory,
		ftpCredentials.ServerFtp.Address,
		ftpCredentials.ServerFtp.User,
		ftpCredentials.ServerFtp.Password,
	)
	if err != nil {
		log.Println("Something goes wrong at ftp connecion. Put this at queue", err)
		return err
	}

	manifest := manifest_service.ReadModPackManifestFiles(modpack, modpack_models.Client)
	oldManifest := manifest_service.ReadModPackManifestFiles(modpack, modpack_models.Server)

	// oldManifest, err := GetFTPManifest(ftpClient)
	if err != nil {
		for _, item := range manifest.Files {
			ftp_service.UploadFileFTP(item.Path, modPackPath, ftpClient)
		}
		log.Printf("Uploading all server modpacks files to ftp server")
		log.Printf("Uploading server manifest file")
		err = ftp_service.UploadFileFTP(cfg.API.ManifestName, modPackPath, ftpClient)
		if err != nil {
			log.Printf("Manifest server uploaded")
			return err
		}
		return nil
	}
	if &oldManifest != nil {
		UploadOrDeleteSyncModPack(oldManifest, manifest, modPackPath, ftpClient)
	}
	return nil
}

func UploadClient(modpack modpack_models.MinecraftModPack,
	ftpCredentials modpack_models.ModPackFtp) error {

	nameNormalized := common.NormalizeString(modpack.Name)

	modPackPath := filepath.Join(
		cfg.API.PublicPath,
		nameNormalized,
		modpack_models.Client.GetFolderName())

	ftpClient, err := ftp_service.OpenFtpConnection(
		ftpCredentials.ClientFtp.Directory,
		ftpCredentials.ClientFtp.Address,
		ftpCredentials.ClientFtp.User,
		ftpCredentials.ClientFtp.Password,
	)
	if err != nil {
		log.Println("Somethings goes wrong at ftp connecion. Put this on queue", err)
		return err
	}
	manifest := manifest_service.ReadModPackManifestFiles(modpack, modpack_models.Client)
	oldManifest, err := GetModPackManifest(ftpClient)
	if err != nil {
		for _, item := range manifest.Files {
			ftp_service.UploadFileFTP(item.Path, modPackPath, ftpClient)
		}
		log.Printf("Uploading all client modpacks files to ftp server")
		return nil
	}
	if &oldManifest != nil {
		UploadOrDeleteSyncModPack(*oldManifest, manifest, modPackPath, ftpClient)
	}
	return nil
}
