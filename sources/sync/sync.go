package sync

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ex3ndr/datasets/resolver"
	"github.com/ex3ndr/datasets/utils"
)

func Sync(src ProjectFile) error {

	// Resolving datasets
	resolved := make([]resolver.Resolved, 0)
	for _, dataset := range src.Datasets {
		fmt.Println("Resolving " + dataset)
		resolve, err := resolver.ResolveDataset(dataset)
		if err != nil {
			return err
		}
		resolved = append(resolved, *resolve)
	}

	// Create datasets directory
	err := os.MkdirAll("external_datasets", 0755)
	if err != nil {
		return err
	}

	// Sync datasets
	for _, dataset := range resolved {
		fmt.Println("Syncing " + dataset.Name)
		err = syncDataset(dataset, "external_datasets")
		if err != nil {
			return err
		}
	}

	return nil
}

func syncDataset(resolved resolver.Resolved, dir string) error {

	// Create directory
	target := filepath.Join(dir, resolved.ID)

	// Check if directory exists
	if _, err := os.Stat(target); !os.IsNotExist(err) {
		fmt.Println("Directory " + target + " already exists. Skipping sync.")
		return nil
	}
	fmt.Println("Downloading " + resolved.ID + " from " + resolved.Endpoint)

	// Create temporary directory
	tempDir, err := os.MkdirTemp("", "datasets")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tempDir)

	// Downloading file
	tempFilePath := filepath.Join(tempDir, "download.tar.gz")
	err = utils.DownloadFile(tempFilePath, resolved.Endpoint, resolved.ID)
	if err != nil {
		return err
	}

	// Unpack file
	fmt.Println("Unpacking " + resolved.ID)
	file, err := os.Open(tempFilePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	err = utils.UnpackTarGz(file, target, 1)
	if err != nil {
		return err
	}

	return nil
}
