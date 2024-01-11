package sync

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/ex3ndr/datasets/resolver"
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
	fmt.Println("Downloading " + resolved.ID)

	// Create temporary directory
	tempDir, err := os.MkdirTemp("", "datasets")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tempDir)

	// Downloading file
	resp, err := http.Get(resolved.Endpoint)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	tempFilePath := filepath.Join(tempDir, "download.tar.gz")
	out, err := os.Create(tempFilePath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the data to file
	_, err = io.Copy(out, resp.Body)
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
	err = unpackTarGz(file, target)
	if err != nil {
		return err
	}

	return nil
}

// unpackTarGz unpacks a tar.gz file to a specified destination, skipping the top-level directory
func unpackTarGz(gzipStream io.Reader, dst string) error {
	unzippedStream, err := gzip.NewReader(gzipStream)
	if err != nil {
		return err
	}
	defer unzippedStream.Close()

	tarReader := tar.NewReader(unzippedStream)

	for {
		header, err := tarReader.Next()

		switch {
		case err == io.EOF:
			return nil
		case err != nil:
			return err
		case header == nil:
			continue
		}

		target := filepath.Join(dst, strings.Join(strings.Split(header.Name, "/")[1:], "/"))

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(target, os.FileMode(header.Mode)); err != nil {
				return err
			}
		case tar.TypeReg:
			outFile, err := os.Create(target)
			if err != nil {
				return err
			}
			if _, err := io.Copy(outFile, tarReader); err != nil {
				outFile.Close()
				return err
			}
			outFile.Close()
		}
	}
}
