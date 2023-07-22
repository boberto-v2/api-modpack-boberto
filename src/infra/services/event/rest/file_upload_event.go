package event_rest

import (
	"encoding/json"
)

type FileUploadEventObject struct {
	Progress float64 `json:"progress"`
}

func CreateFileUploadEventObject(progress float64) ([]byte, error) {
	fileUpload := FileUploadEventObject{
		Progress: progress,
	}
	result, err := json.Marshal(fileUpload)
	if err != nil {
		return nil, err
	}
	return result, nil
}
