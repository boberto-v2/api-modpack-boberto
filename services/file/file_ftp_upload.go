package file_service

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/jlaffaye/ftp"
)

func UploadFilesAndDirsToFTP(server, username, password, destinationDir string, paths []string) error {
	client, err := ftp.Dial(server, ftp.DialWithTimeout(5*time.Second))
	if err != nil {
		return err
	}
	defer client.Quit()
	err = client.Login(username, password)
	if err != nil {
		return err
	}
	err = client.ChangeDir(destinationDir)
	if err != nil {
		return err
	}
	var wg sync.WaitGroup
	errCh := make(chan error, len(paths))
	for _, path := range paths {
		wg.Add(1)
		go func(p string) {
			defer wg.Done()
			err := uploadFileOrDir(client, p)
			if err != nil {
				errCh <- fmt.Errorf("error uploading %s: %s", p, err)
			}
		}(path)
	}
	go func() {
		wg.Wait()
		close(errCh)
	}()

	for err := range errCh {
		fmt.Println(err)
	}
	return nil
}

func uploadFileOrDir(client *ftp.ServerConn, path string) error {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return err
	}

	if fileInfo.IsDir() {
		// Upload directory
		return uploadDir(client, path)
	} else {
		// Upload file
		return uploadFile(client, path)
	}
}

func uploadFile(client *ftp.ServerConn, filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Get the base filename from the file path
	baseFilename := filepath.Base(filePath)

	// Upload the file using STOR command
	err = client.Stor(baseFilename, file)
	if err != nil {
		return err
	}

	return nil
}

func uploadDir(client *ftp.ServerConn, dirPath string) error {
	// Get the base directory name from the path
	baseDir := filepath.Base(dirPath)

	// Create the base directory on the FTP server
	err := client.MakeDir(baseDir)
	if err != nil {
		return err
	}

	// Get the list of files and subdirectories in the local directory
	fileInfos, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return err
	}

	// Recursively upload files and subdirectories
	for _, fileInfo := range fileInfos {
		childPath := filepath.Join(dirPath, fileInfo.Name())
		err := uploadFileOrDir(client, childPath)
		if err != nil {
			return err
		}
	}

	return nil
}
