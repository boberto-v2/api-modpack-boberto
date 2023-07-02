package file_services

import (
	"fmt"
	"hash/crc32"
	"io"
	"os"
)

func GetCRC32(filePath string) (uint32, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return 0, err
	}
	defer file.Close()
	hash := crc32.NewIEEE()
	if _, err := io.Copy(hash, file); err != nil {
		return 0, err
	}
	checksum := hash.Sum32()
	return checksum, nil
}

func CreateDirectoryIfNotExists(dirPath string) error {
	_, err := os.Stat(dirPath)
	if os.IsNotExist(err) {
		err := os.MkdirAll(dirPath, 0755)
		if err != nil {
			return err
		}
		fmt.Println("Directory created:", dirPath)
	}
	return err
}

func IsDirectoryExists(dirPath string) bool {
	_, err := os.Stat(dirPath)
	dirExists := os.IsExist(err)
	return dirExists

}
