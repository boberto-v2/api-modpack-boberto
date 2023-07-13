package pterodactyl

import (
	"testing"

	"github.com/brutalzinn/boberto-modpack-api/test_utils"
)

// Integration test for SendCommandToServer function
func TestSendCommandToServer(t *testing.T) {
	test_utils.SkipCI(t)
	apiURL := "your_api_url"
	apiKey := "your_api_key"
	serverID := "your_server_id"
	command := "your_command"
	err := SendCommandToServer(apiURL, apiKey, serverID, command)
	if err != nil {
		t.Errorf("SendCommandToServer returned an error: %s", err.Error())
	}
}
