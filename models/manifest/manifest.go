package models_manifest

type Manifest struct {
	Id                string `json:"id"`
	Name              string `json:"name"`
	Version           string `json:"version"`
	Author            string `json:"author"`
	Visible           bool   `json:"visible"`
	ClientManifestUrl string `json:"client_manifest_url"`
	ServerManifestUrl string `json:"server_manifest_url"`
}
