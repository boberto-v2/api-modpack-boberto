package ftp_service

import (
	"testing"

	"github.com/brutalzinn/boberto-modpack-api/test_utils"
	"github.com/jlaffaye/ftp"
)

func TestOpenFtpConnection(t *testing.T) {
	test_utils.SkipCI(t)
	ftpDir := "your_ftp_directory"
	ftpHost := "your_ftp_host"
	ftpUser := "your_ftp_username"
	ftpPass := "your_ftp_password"
	conn, err := OpenFtpConnection(ftpDir, ftpHost, ftpUser, ftpPass)
	if err != nil {
		t.Errorf("OpenFtpConnection returned an error: %s", err.Error())
	}
	if conn == nil {
		t.Error("OpenFtpConnection returned nil connection")
	}
}

// Integration test for ReadFileFTP function
func TestReadFileFTP(t *testing.T) {
	test_utils.SkipCI(t)
	conn := getTestFtpConnection()

	ftpFile := "path/to/ftp_file"

	content, err := ReadFileFTP(ftpFile, conn)

	if err != nil {
		t.Errorf("ReadFileFTP returned an error: %s", err.Error())
	}
	if content == nil {
		t.Error("ReadFileFTP returned nil content")
	}
}

// Integration test for UploadFileFTP function
func TestUploadFileFTP(t *testing.T) {
	test_utils.SkipCI(t)
	conn := getTestFtpConnection()

	fileFtp := "path/to/ftp_file"
	relativeToPath := "path/to/local_file"
	err := UploadFileFTP(fileFtp, relativeToPath, conn)
	if err != nil {
		t.Errorf("UploadFileFTP returned an error: %s", err.Error())
	}
}

// Integration test for UploadMultipleFilesFTP function
func TestUploadMultipleFilesFTP(t *testing.T) {
	conn := getTestFtpConnection()
	files := []string{
		"path/to/ftp_file1",
		"path/to/ftp_file2",
	}
	err := UploadMultipleFilesFTP(files, "", conn)

	if err != nil {
		t.Errorf("UploadMultipleFilesFTP returned an error: %s", err.Error())
	}
}

// Integration test for DeleteMultipleFilesFTP function
func TestDeleteMultipleFilesFTP(t *testing.T) {
	test_utils.SkipCI(t)
	conn := getTestFtpConnection()

	filesToDelete := []string{
		"path/to/ftp_file1",
		"path/to/ftp_file2",
	}
	err := DeleteMultipleFilesFTP(filesToDelete, conn)
	if err != nil {
		t.Errorf("DeleteMultipleFilesFTP returned an error: %s", err.Error())
	}
}

// Helper function to initialize a test FTP connection
func getTestFtpConnection() *ftp.ServerConn {
	return nil
}
