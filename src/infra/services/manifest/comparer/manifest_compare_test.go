package manifest_compare

import (
	"hash/crc32"
	"testing"

	manifest_models "github.com/brutalzinn/boberto-modpack-api/src/src/src/infra/services/manifest/models"
	"github.com/stretchr/testify/assert"
)

// To sync mod file to multiple client, we need to check every jar file with crc32 because is the best cheap and cost beneficy of a simple short algorithm that go already implements
// If the new manifest generated not have a especify file, we need to delete this file at ftp server too.
func TestManifestCompareToDeleteFiles(t *testing.T) {
	oldManifest := manifest_models.ManifestFiles{
		Files: []manifest_models.ManifestFile{
			{
				Name:     "test file 1",
				Checksum: crc32.ChecksumIEEE([]byte("somerandomcrc1")),
				Path:     "/path/to/file/test_file_1",
			},
			{
				Name:     "test file 2",
				Checksum: crc32.ChecksumIEEE([]byte("somerandomcrc2")),
				Path:     "/path/to/file/test_file_2",
			},
			{
				Name:     "test file 3",
				Checksum: crc32.ChecksumIEEE([]byte("somerandomcrc3")),
				Path:     "/path/to/file/test_file_3",
			},
		},
	}
	newManifest := manifest_models.ManifestFiles{
		Files: []manifest_models.ManifestFile{
			{
				Name:     "test file 2",
				Checksum: crc32.ChecksumIEEE([]byte("somerandomcrc2")),
				Path:     "/path/to/file/test_file_2",
			},
			{
				Name:     "test file 3",
				Checksum: crc32.ChecksumIEEE([]byte("somerandomcrc3")),
				Path:     "/path/to/file/test_file_3",
			},
		},
	}
	comparer := ManifestComparer{
		NewManifest: newManifest,
		OldManifest: oldManifest,
	}
	result := comparer.Compare()
	assert.Len(t, result.ToDelete, 1)
	assert.Len(t, result.ToUpload, 0)
}

// if the new manifest file have a file that doesnt found at ftp manifest file, we need upload it
func TestManifestCompareToUploadFiles(t *testing.T) {
	oldManifest := manifest_models.ManifestFiles{
		Files: []manifest_models.ManifestFile{
			{
				Name:     "test file 1",
				Checksum: crc32.ChecksumIEEE([]byte("somerandomcrc1")),
				Path:     "/path/to/file/test_file_1",
			},
			{
				Name:     "test file 2",
				Checksum: crc32.ChecksumIEEE([]byte("somerandomcrc2")),
				Path:     "/path/to/file/test_file_2",
			},
		},
	}
	newManifest := manifest_models.ManifestFiles{
		Files: []manifest_models.ManifestFile{
			{
				Name:     "test file 1",
				Checksum: crc32.ChecksumIEEE([]byte("somerandomcrc1")),
				Path:     "/path/to/file/test_file_1",
			},
			{
				Name:     "test file 2",
				Checksum: crc32.ChecksumIEEE([]byte("somerandomcrc2")),
				Path:     "/path/to/file/test_file_2",
			},
			{
				Name:     "test file 3",
				Checksum: crc32.ChecksumIEEE([]byte("somerandomcrc3")),
				Path:     "/path/to/file/test_file_3",
			},
		},
	}
	comparer := ManifestComparer{
		NewManifest: newManifest,
		OldManifest: oldManifest,
	}
	result := comparer.Compare()
	assert.Len(t, result.ToDelete, 0)
	assert.Len(t, result.ToUpload, 1)
}
