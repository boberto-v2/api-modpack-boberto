package manifest_compare

import (
	manifest_models "github.com/brutalzinn/boberto-modpack-api/services/modpack/manifest/models"
)

type ManifestComparerResult struct {
	ToDelete []manifest_models.ManifestFile
	ToUpload []manifest_models.ManifestFile
}

type ManifestComparer struct {
	OldManifest manifest_models.ManifestFiles
	NewManifest manifest_models.ManifestFiles
}

func New(old manifest_models.ManifestFiles, new manifest_models.ManifestFiles) ManifestComparer {
	result := ManifestComparer{
		NewManifest: new,
		OldManifest: old}

	return result
}

func (manifestComparer ManifestComparer) Compare() (result ManifestComparerResult) {
	result.ToUpload = []manifest_models.ManifestFile{}
	for _, oldFile := range manifestComparer.OldManifest.Files {
		found := false
		for _, newFile := range manifestComparer.NewManifest.Files {
			if oldFile.Name != newFile.Name {
				continue
			}
			if oldFile.Checksum != newFile.Checksum {
				result.ToUpload = append(result.ToUpload, newFile)
			}
			found = true
		}
		if !found {
			result.ToDelete = append(result.ToDelete, oldFile)
		}
	}
	return
}
