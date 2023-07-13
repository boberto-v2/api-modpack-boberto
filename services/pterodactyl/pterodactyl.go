package pterodactyl_service

import (
	"encoding/json"
	"fmt"

	"github.com/go-resty/resty/v2"
)

func GetUploadToken(apiURL, apiKey string) (*UploadTokenResponse, error) {
	client := resty.New()
	resp, err := client.R().
		SetHeader("Authorization", "Bearer "+apiKey).
		Post(apiURL + "/api/client/servers/{serverId}/files/upload")
	if err != nil {
		return nil, fmt.Errorf("failed to get upload token: %v", err)
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("failed to get upload token with status code: %d", resp.StatusCode())
	}

	var uploadTokenResp UploadTokenResponse
	err = json.Unmarshal(resp.Body(), &uploadTokenResp)
	if err != nil {
		return nil, fmt.Errorf("failed to parse upload token response: %v", err)
	}

	return &uploadTokenResp, nil
}

func UploadFile(apiURL, uploadToken, filePath string) error {
	client := resty.New()
	resp, err := client.R().
		SetFile("file", filePath).
		SetQueryParam("token", uploadToken).
		Post(apiURL + "/api/client/servers/{serverId}/files/upload")
	if err != nil {
		return fmt.Errorf("file upload failed: %v", err)
	}

	if resp.StatusCode() != 200 {
		return fmt.Errorf("file upload failed with status code: %d", resp.StatusCode())
	}

	return nil
}

func SendCommandToServer(apiURL, apiKey, serverID, command string) error {
	client := resty.New()
	resp, err := client.R().
		SetHeader("Authorization", "Bearer "+apiKey).
		SetQueryParam("command", command).
		Post(apiURL + "/api/client/servers/" + serverID + "/command")
	if err != nil {
		return fmt.Errorf("failed to send command: %v", err)
	}

	if resp.StatusCode() != 204 {
		return fmt.Errorf("command execution failed with status code: %d", resp.StatusCode())
	}

	var commandResp CommandResponse
	err = json.Unmarshal(resp.Body(), &commandResp)
	if err != nil {
		return fmt.Errorf("failed to parse command response: %v", err)
	}

	if !commandResp.Success {
		return fmt.Errorf("command execution failed: %s", commandResp.Message)
	}

	return nil
}

func SendSignalPower(apiURL string, apiKey string, serverID string, signal Signal) error {
	client := resty.New()
	resp, err := client.R().
		SetHeader("Authorization", "Bearer "+apiKey).
		SetQueryParam("signal", signal.GetName()).
		Post(apiURL + "/api/client/servers/" + serverID + "/power")
	if err != nil {
		return fmt.Errorf("failed to send command: %v", err)
	}

	if resp.StatusCode() != 204 {
		return fmt.Errorf("command execution failed with status code: %d", resp.StatusCode())
	}

	var commandResp CommandResponse
	err = json.Unmarshal(resp.Body(), &commandResp)
	if err != nil {
		return fmt.Errorf("failed to parse command response: %v", err)
	}

	if !commandResp.Success {
		return fmt.Errorf("command execution failed: %s", commandResp.Message)
	}

	return nil
}

func GetResources(apiURL, apiKey string, serverID string) (*StatsResponse, error) {
	client := resty.New()
	resp, err := client.R().
		SetHeader("Authorization", "Bearer "+apiKey).
		Get(apiURL + "/api/client/servers/" + serverID + "/resources")
	if err != nil {
		return nil, fmt.Errorf("failed to get stats response: %v", err)
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("failed to get stats response with status code: %d", resp.StatusCode())
	}

	var statsResponse StatsResponse
	err = json.Unmarshal(resp.Body(), &statsResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to parse stats response: %v", err)
	}

	return &statsResponse, nil
}
