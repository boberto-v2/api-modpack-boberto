package upload_service

import (
	"errors"
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
func GetUploadPath(id string) (string, error) {
	uploadCache, found := upload_cache.GetById(id)
	if !found {
		return "", errors.New("upload token expired or invalid")
	}
	return uploadCache.OutputDir, nil
}

func SaveFiles(id string, files []*multipart.FileHeader) error {
	for _, file := range files {
		err := SaveFile(id, file)
		if err != nil {
			return err
		}
	}
	return nil
}

func SaveFile(id string, file *multipart.FileHeader) error {
	outputPath, err := GetUploadPath(id)
	if err != nil {
		return err
	}
	finalOutputFile := filepath.Join(outputPath, file.Filename)
	out, err := os.Create(finalOutputFile)
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()
	readerFile, _ := file.Open()
	_, err = io.Copy(out, readerFile)
	if err != nil {
		log.Fatal(err)
	}
	unZipIfNecessary(finalOutputFile, outputPath)
	return nil
}

func unZipIfNecessary(finalOutputFile string, outputPath string) {
	fileExtension := filepath.Ext(finalOutputFile)
	if fileExtension != ".zip" {
		return
	}
	file_service.Unzip(finalOutputFile, outputPath)
	os.Remove(finalOutputFile)
}
