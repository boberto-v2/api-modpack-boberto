package upload_service

import (
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"github.com/brutalzinn/boberto-modpack-api/common"
	file_service "github.com/brutalzinn/boberto-modpack-api/services/file"
	upload_cache "github.com/brutalzinn/boberto-modpack-api/services/upload/cache"
	upload_cache_models "github.com/brutalzinn/boberto-modpack-api/services/upload/cache/models"
)

func Create(outputDir string) upload_cache_models.UploadCache {
	id := common.GenerateUUID()
	uploadCache := upload_cache_models.UploadCache{
		Id:        id,
		CreateAt:  time.Now(),
		OutputDir: outputDir,
		ExpireAt:  time.Now().Add(time.Duration(time.Hour * 1)),
	}
	upload_cache.Create(uploadCache)
	return uploadCache
}
func GetById(id string) (*upload_cache_models.UploadCache, error) {
	uploadCache, found := upload_cache.GetById(id)
	if !found {
		return nil, errors.New("The token provided is invalid or expired.")
	}
	return &uploadCache, nil
}

func SaveFiles(id string, files []*multipart.FileHeader, callback func(string)) error {
	uploadCache, err := GetById(id)
	if err != nil {
		return err
	}
	for _, file := range files {
		err := saveFile(uploadCache.OutputDir, file)
		if err != nil {
			return err
		}
		callback(fmt.Sprintf("saving...", file.Filename))
	}
	return nil
}

func saveFile(fileDir string, file *multipart.FileHeader) error {
	filePath := filepath.Join(fileDir, file.Filename)
	out, err := os.Create(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()
	readerFile, _ := file.Open()
	_, err = io.Copy(out, readerFile)
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

func UnZip(zipFilePath string, outputPath string, callback func(string)) {
	callback(fmt.Sprintf("unziping.. %s", zipFilePath))
	file_service.Unzip(zipFilePath, outputPath)
	callback(fmt.Sprintf("unzip completed %s", zipFilePath))
}

func IsZip(filePath string) bool {
	fileExtension := filepath.Ext(filePath)
	isZipExtenion := fileExtension == ".zip"
	return isZipExtenion
}
