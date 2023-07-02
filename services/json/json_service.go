package json_service

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type JsonWritter struct {
	Path string
	Data any
}

func (jsonWritter JsonWritter) CreateFile() {
	data, _ := json.MarshalIndent(jsonWritter.Data, "", " ")
	_ = ioutil.WriteFile(jsonWritter.Path, data, 0644)
}

func ReadFile(jsonWritter JsonWritter) any {
	content, err := ioutil.ReadFile(jsonWritter.Path)
	if err != nil {
		log.Fatal("Error when opening file: ", err)
	}
	var payload any
	err = json.Unmarshal(content, &payload)
	if err != nil {
		log.Fatal("Error during Unmarshal(): ", err)
	}
	return payload
}
