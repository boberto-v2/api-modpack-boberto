package upload_service

import (
	"errors"
	"path/filepath"
	"time"

	"github.com/brutalzinn/boberto-modpack-api/common"
	file_service "github.com/brutalzinn/boberto-modpack-api/services/file"
	upload_cache "github.com/brutalzinn/boberto-modpack-api/services/upload/cache"
)

func Create(outputDir string) *upload_cache.UploadCache {
	id := common.GenerateUUID()
	uploadCache := upload_cache.UploadCache{
		Id:        id,
		CreateAt:  time.Now(),
		OutputDir: outputDir,
		Status:    upload_cache.UPLOAD_CREATED,
		ExpireAt:  time.Now().Add(time.Duration(time.Hour * 1)),
	}
	upload_cache.Create(uploadCache)
	return &uploadCache
}

func GetById(id string) (*upload_cache.UploadCache, error) {
	uploadCache, found := upload_cache.GetById(id)
	if !found {
		return nil, errors.New("The token provided is invalid or expired.")
	}
	return &uploadCache, nil
}

func UnZip(zipFilePath string, outputPath string) {
	file_service.Unzip(zipFilePath, outputPath)
}

func IsZip(filePath string) bool {
	fileExtension := filepath.Ext(filePath)
	isZipExtenion := fileExtension == ".zip"
	return isZipExtenion
}
