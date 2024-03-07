package resolver

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/ex3ndr/datasets/project"
	"gopkg.in/yaml.v3"
)

type InternalResolved struct {
	ID       *string
	Endpoint string
}

func resolveHTTPDataset(dataset string) (*InternalResolved, error) {
	// Return resolved
	return &InternalResolved{
		Endpoint: dataset,
	}, nil
}

func resolveFileDataset(path string) (*InternalResolved, error) {

	// Check if file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, errors.New("file not found: " + path)
	}

	// Check if file is a tar.gz archive
	if !strings.HasSuffix(strings.ToLower(path), ".tar.gz") {
		return nil, errors.New("file is not a tar.gz archive: " + path)
	}

	// Get file name from path and without extension
	resolvedName := filepath.Base(path)
	resolvedName = resolvedName[:len(resolvedName)-7]

	// Return resolved
	return &InternalResolved{
		ID:       &resolvedName,
		Endpoint: "file:" + path,
	}, nil
}

func resolveGithubDataset(dataset string) (*InternalResolved, error) {
	return nil, errors.New("not implemented")
}

func resolveHugginFaceDataset(dataset string) (*InternalResolved, error) {
	return nil, errors.New("not implemented")
}

func resolveStandardDataset(dataset string) (*InternalResolved, error) {

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

	// Applying mirror
	url := ResolveMirror(data.URL)

	// Resolving hashes
	return &InternalResolved{
		ID:       &id,
		Endpoint: url,
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

func doResolveDataset(dataset project.DatasetReference) (*InternalResolved, error) {

	// Local file
	if strings.HasPrefix(dataset.Source, "file:") {
		return resolveFileDataset(dataset.Source[5:])
	}

	// Remote file
	if strings.HasPrefix(dataset.Source, "http://") || strings.HasPrefix(dataset.Source, "https://") {
		return resolveHTTPDataset(dataset.Source)
	}

	// Github
	if strings.HasPrefix(dataset.Source, "github.com/") {
		return resolveGithubDataset(dataset.Source)
	}

	// HuggingFace
	if strings.HasPrefix(dataset.Source, "huggingface.co") {
		return resolveHugginFaceDataset(dataset.Source)
	}

	// Standard
	if isValidStandartName(dataset.Source) {
		return resolveStandardDataset(dataset.Source)
	}

	return nil, errors.New("unknown dataset type: " + dataset.Source)
}

func ResolveDataset(dataset project.DatasetReference) (*Resolved, error) {

	// Internal resolve
	ir, err := doResolveDataset(dataset)
	if err != nil {
		return nil, err
	}

	// Finalize
	var id string
	if ir.ID == nil && dataset.Name == nil {
		return nil, errors.New("unable to resolve dataset name automatically, please provide it manually")
	} else if dataset.Name != nil {
		id = *dataset.Name
	} else {
		id = *ir.ID
	}

	// Format
	var format string = "tar-gz"
	if dataset.Format != nil {
		format = *dataset.Format
	}

	// Return resolved
	return &Resolved{
		ID:       id,
		Endpoint: ir.Endpoint,
		Format:   format,
	}, nil
}
