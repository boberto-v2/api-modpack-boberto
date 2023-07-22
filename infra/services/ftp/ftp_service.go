package ftp_service

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	file_service "github.com/brutalzinn/boberto-modpack-api/infra/services/file"
	"github.com/jlaffaye/ftp"
)

const (
	FTP_SEPARATOR = "/"
)

// / TODO: Explain to Daniel why i choose FTP upload files
func OpenFtpConnection(ftpDir string, ftpHost string, ftpUser string, ftpPass string) (*ftp.ServerConn, error) {
	client, err := ftp.Dial(ftpHost)
	if err != nil {
		return client, err
	}
	// defer client.Quit()
	err = client.Login(ftpUser, ftpPass)
	if err != nil {
		return client, err
	}

	err = client.ChangeDir(ftpDir)
	if err != nil {
		return client, err
	}
	return client, err
}

func ReadFileFTP(ftpFile string, client *ftp.ServerConn) ([]byte, error) {
	r, err := client.Retr(ftpFile)
	if err != nil {
		log.Printf("Error on read file using client.Retr %s", err.Error())
		return nil, err
	}
	defer r.Close()
	buf, err := ioutil.ReadAll(r)
	fmt.Printf("Read: %s\n", ftpFile)
	return buf, nil
}

func UploadFileFTP(fileFtp string, relativeToLocalPath string, client *ftp.ServerConn) error {
	client.ChangeDir(FTP_SEPARATOR)
	directory, filename := filepath.Split(fileFtp)
	dirs := strings.Split(directory, string(os.PathSeparator))
	for _, dir := range dirs {
		if dir == "" {
			continue
		}
		err := client.ChangeDir(dir)
		if err != nil {
			err = client.MakeDir(dir)
			if err != nil {
				return fmt.Errorf("failed to create directory: %v", err)
			}
			err = client.ChangeDir(dir)
			if err != nil {
				return fmt.Errorf("failed to change directory: %v", err)
			}
		}
	}
	file, err := os.Open(filepath.Join(relativeToLocalPath, fileFtp))
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()
	err = client.Stor(filename, file)
	if err != nil {
		return fmt.Errorf("failed to upload file: %v", err)
	}
	return nil
}

func UploadMultipleFilesFTP(files []string, relativeToLocalPath string, client *ftp.ServerConn) error {
	for _, filePath := range files {
		err := UploadFileFTP(filePath, relativeToLocalPath, client)
		if err != nil {
			return err
		}
	}
	return nil
}

func DeleteMultipleFilesFTP(filesToDelete []string, client *ftp.ServerConn) error {
	for _, file := range filesToDelete {
		err := DeleteFileFTP(file, client)
		if err != nil {
			return err
		}
	}
	for _, file := range filesToDelete {
		parentDir := file_service.GetParentDirectory(file, FTP_SEPARATOR)
		err := deleteEmptyParentDirectories(parentDir, client)
		if err != nil {
			return err
		}
	}
	return nil
}

func DeleteFileFTP(fileFtp string, client *ftp.ServerConn) error {
	err := client.Delete(fileFtp)
	if err != nil {
		return err
	}
	fmt.Printf("Deleted: %s\n", fileFtp)
	parentDir := file_service.GetParentDirectory(fileFtp, FTP_SEPARATOR)
	err = deleteEmptyParentDirectories(parentDir, client)
	if err != nil {
		return err
	}
	return nil
}

func deleteEmptyParentDirectories(directory string, client *ftp.ServerConn) error {
	files, err := client.List(directory)
	if err != nil {
		return err
	}
	filesCount := 1
	for _, entry := range files {
		if entry.Type == 0 {
			filesCount++
		}
	}
	if filesCount == 0 {
		err := client.RemoveDir(directory)
		if err != nil {
			return err
		}
		fmt.Printf("Deleted directory: %s\n", directory)
		parentDir := file_service.GetParentDirectory(directory, FTP_SEPARATOR)
		if parentDir != "" {
			err := deleteEmptyParentDirectories(parentDir, client)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
