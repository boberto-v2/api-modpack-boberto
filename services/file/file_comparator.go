package file_service

// I LOVE GOROUTINES.
import (
	"hash/crc32"
	"io"
	"os"
	"sync"
)

type checksumResult struct {
	Checksum uint32
	Err      error
}

func calculateChecksum(filePath string, wg *sync.WaitGroup, result chan<- checksumResult) {
	defer wg.Done()

	file, err := os.Open(filePath)
	if err != nil {
		result <- checksumResult{0, err}
		return
	}
	defer file.Close()

	hasher := crc32.NewIEEE()
	if _, err := io.Copy(hasher, file); err != nil {
		result <- checksumResult{0, err}
		return
	}

	result <- checksumResult{hasher.Sum32(), nil}
}
func CompareChecksum(file1Path, file2Path string) (bool, error) {
	var wg sync.WaitGroup
	resultCh := make(chan checksumResult, 2)

	wg.Add(2)
	go calculateChecksum(file1Path, &wg, resultCh)
	go calculateChecksum(file2Path, &wg, resultCh)

	wg.Wait()
	close(resultCh)

	var checksum1, checksum2 uint32
	var err1, err2 error

	for res := range resultCh {
		if res.Err != nil {
			return false, res.Err
		}

		if checksum1 == 0 {
			checksum1 = res.Checksum
			err1 = res.Err
		} else {
			checksum2 = res.Checksum
			err2 = res.Err
		}
	}

	if err1 != nil {
		return false, err1
	}
	if err2 != nil {
		return false, err2
	}

	return checksum1 == checksum2, nil
}
