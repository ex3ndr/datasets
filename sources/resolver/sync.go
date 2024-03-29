package resolver

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ex3ndr/datasets/project"
	"github.com/ex3ndr/datasets/utils"
)

func Sync(src project.ProjectFile) {
	start := time.Now()
	error := doSync(src)
	if error != nil {
		fmt.Println(("❌  " + utils.Failure("error "+error.Error())))
		os.Exit(1)
	} else {
		fmt.Println("✨  Done in " + time.Since(start).String() + ".")
	}
}

func doSync(src project.ProjectFile) error {

	// Create datasets directory if not exists
	err := os.MkdirAll(filepath.Join("external_datasets", ".downloads"), 0755)
	if err != nil {
		return err
	}

	// Resolving datasets
	fmt.Println(utils.Faint("[1/3]") + " 🔎  Resolving datasets...")
	resolved := make([]*Resolved, 0)
	for _, dataset := range src.Datasets {
		resolve, err := ResolveDataset(dataset)
		if err != nil {
			return err
		}
		resolved = append(resolved, resolve)
	}

	// Fetching packages
	fmt.Println(utils.Faint("[2/3]") + " 🚚  Fetching datasets...")

	// Create temporary directory
	tempDir, err := os.MkdirTemp(filepath.Join("external_datasets", ".downloads"), "datasets-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tempDir)

	// Downloading files
	for _, resolved := range resolved {

		// Target directory
		target := filepath.Join("external_datasets", resolved.ID)

		// Check if directory exists
		if _, err := os.Stat(target); !os.IsNotExist(err) {
			continue
		}

		// Check if file
		if strings.HasPrefix(resolved.Endpoint, "file:") {
			path := resolved.Endpoint[5:]
			resolved.Location = &path
			continue
		}

		// Download file
		tempFilePath := filepath.Join(tempDir, resolved.ID+".tar.gz")
		err = utils.DownloadFile(tempFilePath, resolved.Endpoint, "          "+resolved.ID)
		if err != nil {
			fmt.Println(utils.ClearLine() + "          " + utils.Failure("failure") + " " + resolved.ID)
			return err
		}
		fmt.Println(utils.ClearLine() + "          " + utils.Success("success") + " " + resolved.ID)
		resolved.Location = &tempFilePath
	}

	// Unpacking packages
	fmt.Println(utils.Faint("[3/3]") + " 📦  Unpacking datasets...")
	for _, resolved := range resolved {

		// Check if unpack required
		if resolved.Location == nil {
			continue
		}

		// Create directory
		target := filepath.Join("external_datasets", resolved.ID)

		// File path
		tempFilePath := *resolved.Location

		// Unpack file
		if resolved.Format == "tar-gz" {
			err = utils.UnpackTarGz(tempFilePath, target, 1, "          "+resolved.ID)
			if err != nil {
				fmt.Println(utils.ClearLine() + "          " + utils.Failure("failure") + " " + resolved.ID)
				return err
			}
			fmt.Println(utils.ClearLine() + "          " + utils.Success("success") + " " + resolved.ID)
		} else if resolved.Format == "tar" {
			err = utils.UnpackTar(tempFilePath, target, 1, "          "+resolved.ID)
			if err != nil {
				fmt.Println(utils.ClearLine() + "          " + utils.Failure("failure") + " " + resolved.ID)
				return err
			}
			fmt.Println(utils.ClearLine() + "          " + utils.Success("success") + " " + resolved.ID)
		} else {
			return fmt.Errorf("unknown format: " + resolved.Format)
		}
	}

	return nil
}
