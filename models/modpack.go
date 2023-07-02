package models

import (
	"strings"

	"github.com/gofrs/uuid"
)

type Modpack struct {
	Id             string        `json:"id"`
	Name           string        `json:"name"`
	Premium        bool          `json:"premium"`
	Status         ModPackStatus `json:"status"`
	Address        string        `json:"address"`
	NormalizedName string        `json:"normalized_name"`
	Author         string        `json:"author"`
	Version        string        `json:"version"`
}

func (modpack *Modpack) New() {
	id, _ := uuid.NewV4()
	modpack.Id = id.String()
	modpackNameLower := strings.ToLower(modpack.Name)
	modpack.NormalizedName = strings.ReplaceAll(modpackNameLower, " ", "_")
}
