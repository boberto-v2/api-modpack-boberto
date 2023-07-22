package pterodactyl_service

import (
	"os"
	"testing"

	"github.com/brutalzinn/boberto-modpack-api/src/src/src/test_utils"
)

// Integration test for SendCommandToServer function
func TestSendCommandToServer(t *testing.T) {
	test_utils.SkipCI(t)
	pteroConn := getPteroConnection()
	command := "say hello server"
	err := SendCommandToServer(pteroConn.ApiUrl, pteroConn.ApiKey, pteroConn.ServerId, command)
	if err != nil {
		t.Errorf("SendCommandToServer returned an error: %s", err.Error())
	}
}

// Integration test for SendSignalPower function
func TestSendSignalRestart(t *testing.T) {
	test_utils.SkipCI(t)
	pteroConn := getPteroConnection()
	err := SendSignalPower(pteroConn.ApiUrl, pteroConn.ApiKey, pteroConn.ServerId, RESTART)
	if err != nil {
		t.Errorf("SendCommandToServer returned an error: %s", err.Error())
	}
}

func TestGetServerResources(t *testing.T) {
	test_utils.SkipCI(t)
	pteroConn := getPteroConnection()
	resources, err := GetResources(pteroConn.ApiUrl, pteroConn.ApiKey, pteroConn.ServerId)
	if err != nil {
		t.Errorf("SendCommandToServer returned an error: %s", err.Error())
	}
	t.Logf(resources.Object)
}

type PteroApiInfo struct {
	ApiUrl   string
	ApiKey   string
	ServerId string
}

func getPteroConnection() PteroApiInfo {
	pteroApiConnection := PteroApiInfo{
		ApiUrl:   os.Getenv("ptero_api_url"),
		ApiKey:   os.Getenv("ptero_api_key"),
		ServerId: os.Getenv("ptero_server_id"),
	}
	return pteroApiConnection
}
