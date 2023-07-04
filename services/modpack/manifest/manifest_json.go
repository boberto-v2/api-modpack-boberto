package manifest_service

import (
	"encoding/json"
	"io/ioutil"
	"log"

	manifest_models "github.com/brutalzinn/boberto-modpack-api/services/modpack/manifest/models"
)

func WriteModPackJsonManifest(fileName string, data manifest_models.ManifestModPack) {
	jsonData, _ := json.MarshalIndent(data, "", " ")
	_ = ioutil.WriteFile(fileName, jsonData, 0644)
}

func ReadModPackJsonManifest(fileName string) manifest_models.ManifestModPack {
	content, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatal("Error when opening file: ", err)
	}
	var payload manifest_models.ManifestModPack
	err = json.Unmarshal(content, &payload)
	if err != nil {
		log.Fatal("Error during Unmarshal(): ", err)
	}
	return payload
}

func WriteManifestJsonFiles(fileName string, data manifest_models.ManifestFiles) {
	jsonData, _ := json.MarshalIndent(data, "", " ")
	_ = ioutil.WriteFile(fileName, jsonData, 0644)
}

func ReadManifestJsonFiles(fileName string) manifest_models.ManifestFiles {
	content, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatal("Error when opening file: ", err)
	}
	var payload manifest_models.ManifestFiles
	err = json.Unmarshal(content, &payload)
	if err != nil {
		log.Fatal("Error during Unmarshal(): ", err)
	}
	return payload
}
