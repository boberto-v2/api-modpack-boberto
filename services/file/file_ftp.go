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

const (
	initialFtpDir = "/"
)

func DeleteFilesFromFTP(filesToDelete []string, c *ftp.ServerConn) error {
	for _, file := range filesToDelete {
		err := c.Delete(file)
		if err != nil {
			return err
		}
		fmt.Printf("Deleted: %s\n", file)
	}
	for _, file := range filesToDelete {
		parentDir := GetParentDirectory(file)
		err := DeleteEmptyParentDirectories(parentDir, c)
		if err != nil {
			return err
		}
	}

	return nil
}

func DeleteFileFromFTP(fileFtp string, c *ftp.ServerConn) error {
	err := c.Delete(fileFtp)
	if err != nil {
		return err
	}
	fmt.Printf("Deleted: %s\n", fileFtp)
	parentDir := GetParentDirectory(fileFtp)
	err = DeleteEmptyParentDirectories(parentDir, c)
	if err != nil {
		return err
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
func DeleteEmptyParentDirectories(directory string, client *ftp.ServerConn) error {
	files, err := client.List(directory)
	if err != nil {
		return err
	}

	if len(files) == 0 {
		err := client.RemoveDir(directory)
		if err != nil {
			return err
		}

		fmt.Printf("Deleted directory: %s\n", directory)
		parentDir := GetParentDirectory(directory)
		if parentDir != "" {
			err := DeleteEmptyParentDirectories(parentDir, client)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
func UploadFileFTPWithDirectories(fileFtp string, relativeToPath string, client *ftp.ServerConn) error {
	client.ChangeDir(initialFtpDir)
	test := initialFtpDir + fileFtp
	directory, filename := filepath.Split(test)
	dirs := strings.Split(directory, string(os.PathSeparator))
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
	file, err := os.Open(filepath.Join(relativeToPath, fileFtp))
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

func UploadFileFTP(fileFtp string, relativeToPath string, client *ftp.ServerConn) error {
	err := client.ChangeDir(initialFtpDir)
	file, err := os.Open(filepath.Join(relativeToPath, fileFtp))
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()
	err = client.Stor(fileFtp, file)
	if err != nil {
		return fmt.Errorf("failed to upload file: %v", err)
	}
	return nil
}

func UploadFilesFTPWithDirectories(files []string, relativeToPath string, client *ftp.ServerConn) error {
	for _, filePath := range files {
		err := UploadFileFTPWithDirectories(filePath, relativeToPath, client)
		if err != nil {
			return err
		}
	}
	return nil
}
