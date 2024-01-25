package resolver

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

func resolveHTTPDataset(dataset string) (*Resolved, error) {
	return nil, errors.New("not implemented")
}

func resolveFileDataset(path string) (*Resolved, error) {

	// Check if file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, errors.New("file not found: " + path)
	}

	// Check if file is a tar.gz archive
	if !strings.HasSuffix(strings.ToLower(path), ".tar.gz") {
		return nil, errors.New("file is not a tar.gz archive: " + path)
	}

	// Get file name from path and without extension
	name := filepath.Base(path)
	name = name[:len(name)-7]

	// Return resolved
	return &Resolved{
		ID:       name,
		Endpoint: "file:" + path,
	}, nil
}

func resolveGithubDataset(dataset string) (*Resolved, error) {
	return nil, errors.New("not implemented")
}

func resolveHugginFaceDataset(dataset string) (*Resolved, error) {
	return nil, errors.New("not implemented")
}

func resolveStandardDataset(dataset string) (*Resolved, error) {

	// Resolve name and version
	name := dataset
	version := ""
	if strings.Contains(dataset, "@") {
		parts := strings.Split(dataset, "@")
		name = parts[0]
		version = parts[1]
	}

	// Downloading descriptor
	resp, err := http.Get("https://raw.githubusercontent.com/ex3ndr/datasets/main/collection/" + name + ".yaml")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Checking status
	if resp.StatusCode != 200 {
		return nil, errors.New("dataset " + name + " not found")
	}

	// Parsing descriptor
	var descriptor DatasetDescriptor
	err = yaml.NewDecoder(resp.Body).Decode(&descriptor)
	if err != nil {
		return nil, err
	}

	// Loading id
	id := descriptor.ID

	// Resolving
	data := descriptor.Dataset
	if version != "" {
		d := descriptor.Extras[version]
		data = d.Dataset
		id = d.ID
	}

	// Resolving hashes
	return &Resolved{
		ID:       id,
		Endpoint: data.URL,
	}, nil
}

func isValidStandartName(name string) bool {
	// Must start with an alphanumeric character
	// Can contain periods, dashes, and underscores
	// Must not end with a period or a dash
	// Optional @version part follows similar rules
	validNamePattern := `^[a-zA-Z0-9]+[a-zA-Z0-9_.-]*[a-zA-Z0-9]+(@[a-zA-Z0-9]+[a-zA-Z0-9_.-]*[a-zA-Z0-9]+)?$`

	matched, err := regexp.MatchString(validNamePattern, name)
	if err != nil {
		fmt.Println("Error in regex pattern: ", err)
		return false
	}
	return matched
}

func ResolveDataset(dataset string) (*Resolved, error) {

	// Local file
	if strings.HasPrefix(dataset, "file:") {
		return resolveFileDataset(dataset[5:])
	}

	// Remote file
	if strings.HasPrefix(dataset, "http://") || strings.HasPrefix(dataset, "https://") {
		return resolveHTTPDataset(dataset)
	}

	// Github
	if strings.HasPrefix(dataset, "github.com/") {
		return resolveGithubDataset(dataset)
	}

	// HuggingFace
	if strings.HasPrefix(dataset, "huggingface.co") {
		return resolveHugginFaceDataset(dataset)
	}

	// Standard
	if isValidStandartName(dataset) {
		return resolveStandardDataset(dataset)
	}

	return nil, errors.New("unknown dataset type: " + dataset)
}
