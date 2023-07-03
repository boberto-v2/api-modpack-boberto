package manifest_service

import (
	"encoding/json"
	"io/ioutil"
	"log"

	models_manifest "github.com/brutalzinn/boberto-modpack-api/models/manifest"
)

func CreateModPackManifest(fileName string, data models_manifest.Manifest) {
	jsonData, _ := json.MarshalIndent(data, "", " ")
	_ = ioutil.WriteFile(fileName, jsonData, 0644)
}

func ReadModPackManifestFile(fileName string) models_manifest.Manifest {
	content, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatal("Error when opening file: ", err)
	}
	var payload models_manifest.Manifest
	err = json.Unmarshal(content, &payload)
	if err != nil {
		log.Fatal("Error during Unmarshal(): ", err)
	}
	return payload
}

func CreateManifestFile(fileName string, data models_manifest.ModPackFileManifest) {
	jsonData, _ := json.MarshalIndent(data, "", " ")
	_ = ioutil.WriteFile(fileName, jsonData, 0644)
}

func ReadManifestFile(fileName string) models_manifest.ModPackFileManifest {
	content, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatal("Error when opening file: ", err)
	}
	var payload models_manifest.ModPackFileManifest
	err = json.Unmarshal(content, &payload)
	if err != nil {
		log.Fatal("Error during Unmarshal(): ", err)
	}
	return payload
}
