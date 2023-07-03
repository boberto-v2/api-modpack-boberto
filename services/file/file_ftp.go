package file_service

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/jlaffaye/ftp"
)

func UploadFilesToFTP(files []string, modPackPath string, client *ftp.ServerConn) error {
	for _, filePath := range files {
		directory, filename := filepath.Split(filePath)
		dirs := strings.Split(directory, "/")
		for _, dir := range dirs {
			if dir == "" {
				continue
			}
			_, err := client.List(dir)
			if err == nil {
				err = client.ChangeDir(dir)
				if err != nil {
					return fmt.Errorf("failed to change directory: %v", err)
				}
			} else {
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
		file, err := os.Open(filepath.Join(modPackPath, filePath))
		if err != nil {
			return fmt.Errorf("failed to open file: %v", err)
		}
		defer file.Close()
		err = client.Stor(filename, file)
		if err != nil {
			return fmt.Errorf("failed to upload file: %v", err)
		}
	}
	return nil
}

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

func UploadFileFtp(localFile string, client *ftp.ServerConn) error {
	fileReader, err := os.Open(localFile)
	if err != nil {
		log.Printf("Cant open local file %s", err.Error())
		return err
	}
	defer fileReader.Close()
	err = client.Stor(localFile, fileReader)
	log.Printf("Error on upload file using client.Stor %s", err.Error())
	if err != nil {
		return err
	}
	log.Printf("Uploaded file %s", localFile)
	return err
}

func DeleteFileFTP(ftpFile string, client *ftp.ServerConn) error {
	err := client.Delete(ftpFile)
	if err != nil {
		log.Printf("Error on delete file %s", err.Error())
		return err
	}
	fmt.Printf("Deleted: %s\n", ftpFile)
	return nil
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
