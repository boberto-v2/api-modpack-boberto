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

func (manifestComparer ManifestComparer) Compare() ManifestComparerResult {
	result := ManifestComparerResult{
		ToDelete: []manifest_models.ManifestFile{},
		ToUpload: []manifest_models.ManifestFile{},
	}
	checksumsOld := make(map[uint32]bool)
	for _, fileOld := range manifestComparer.OldManifest.Files {
		checksumsOld[fileOld.Checksum] = true
	}
	for _, fileNew := range manifestComparer.NewManifest.Files {
		if _, exists := checksumsOld[fileNew.Checksum]; !exists {
			result.ToUpload = append(result.ToUpload, fileNew)
		}
	}
	for _, fileOld := range manifestComparer.OldManifest.Files {
		found := false
		for _, fileNew := range manifestComparer.NewManifest.Files {
			if fileNew.Checksum == fileOld.Checksum {
				found = true
				break
			}
		}
		if !found {
			result.ToDelete = append(result.ToDelete, fileOld)
		}
	}
	return result
}

//TODO: Old comparer ( legacy and hard to read )

// func (manifestComparer ManifestComparer) Compare() (result ManifestComparerResult) {
// 	result.ToUpload = []manifest_models.ManifestFile{}
// 	for _, oldFile := range manifestComparer.OldManifest.Files {
// 		found := false
// 		for _, newFile := range manifestComparer.NewManifest.Files {
// 			if oldFile.Name != newFile.Name {
// 				continue
// 			}
// 			if oldFile.Checksum != newFile.Checksum {
// 				result.ToUpload = append(result.ToUpload, newFile)
// 			}
// 			found = true
// 		}
// 		if !found {
// 			result.ToDelete = append(result.ToDelete, oldFile)
// 		}
// 	}
// 	return
// }
