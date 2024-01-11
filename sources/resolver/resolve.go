package resolver

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

func resolveHTTPDataset(dataset string) (*Resolved, error) {
	return nil, errors.New("not implemented")
}

func resolveFileDataset(dataset string) (*Resolved, error) {
	return nil, errors.New("not implemented")
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
		name = d.Name
		id = d.ID
	}

	// Resolving hashes
	return &Resolved{
		ID:       id,
		Name:     name,
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
