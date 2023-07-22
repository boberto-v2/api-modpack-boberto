package ftp_service_test

import (
	"testing"

	ftp_service "github.com/brutalzinn/boberto-modpack-api/services/ftp"
	"github.com/brutalzinn/boberto-modpack-api/test_utils"
	"github.com/jlaffaye/ftp"
)

// i think for case these tests... i suppose i will attach a volume using docker.. but.. i need to pause this project now
// i need a little pause to put my brain at head again! :)
func TestOpenFtpConnection(t *testing.T) {
	test_utils.SkipCI(t)
	ftpDir := "/"
	ftpHost := "localhost:21"
	ftpUser := "test"
	ftpPass := "test"
	conn, err := ftp_service.OpenFtpConnection(ftpDir, ftpHost, ftpUser, ftpPass)
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

	content, err := ftp_service.ReadFileFTP(ftpFile, conn)

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
	err := ftp_service.UploadFileFTP("subpath/file_1.txt", "files", conn)
	if err != nil {
		t.Errorf("UploadFileFTP returned an error: %s", err.Error())
	}
}

// Integration test for UploadMultipleFilesFTP function
// files/subpath
func TestUploadMultipleFilesFTP(t *testing.T) {
	conn := getTestFtpConnection()
	files := []string{
		//"file_1.txt",
		//"file_2.txt",
		//"file_3.txt",
		"subpath/file_1.txt",
	}
	err := ftp_service.UploadMultipleFilesFTP(files, "files", conn)

	if err != nil {
		t.Errorf("UploadMultipleFilesFTP returned an error: %s", err.Error())
	}
}

// Integration test for DeleteMultipleFilesFTP function
func TestDeleteMultipleFilesFTP(t *testing.T) {
	test_utils.SkipCI(t)
	conn := getTestFtpConnection()
	filesToDelete := []string{
		"file_1.txt",
		"file_2.txt",
		"file_3.txt",
	}
	err := ftp_service.DeleteMultipleFilesFTP(filesToDelete, conn)
	if err != nil {
		t.Errorf("DeleteMultipleFilesFTP returned an error: %s", err.Error())
	}
}

// Helper function to initialize a test FTP connection
func getTestFtpConnection() *ftp.ServerConn {
	conn, err := ftp_service.OpenFtpConnection("/", "localhost:21", "test", "test")
	if err != nil {
		return nil
	}
	return conn
}
