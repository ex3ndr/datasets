package project

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

type ProjectFile struct {
	Datasets []DatasetReference
}

type DatasetReference struct {
	Source string
	Name   *string
}

func UnmarshalProject(d []byte) (*ProjectFile, error) {

	// Unmarshal YAML
	var source map[string]yaml.Node
	err := yaml.Unmarshal(d, &source)
	if err != nil {
		return nil, err
	}

	// Create project file
	res := &ProjectFile{Datasets: make([]DatasetReference, 0)}

	// Check if source is a document (always true?)
	if datasets, ok := source["datasets"]; ok {

		// Check if datasets is a mapping node
		if datasets.Kind != yaml.SequenceNode {
			return nil, fmt.Errorf("datasets is not a list")
		}

		// Iterate over datasets
		for i := 0; i < len(datasets.Content); i++ {
			if datasets.Content[i].Kind == yaml.ScalarNode { // String-ish node
				res.Datasets = append(res.Datasets, DatasetReference{Source: datasets.Content[i].Value})
			} else if datasets.Content[i].Kind == yaml.MappingNode {
				node := datasets.Content[i]
				var name *string
				var source *string = nil
				for j := 0; j < len(node.Content); j += 2 {
					if node.Content[j].Value == "name" {
						name = &node.Content[j+1].Value
					} else if node.Content[j].Value == "source" {
						source = &node.Content[j+1].Value
					} else {
						return nil, fmt.Errorf("unknown key: %s", node.Content[j].Value)
					}
				}
				if source == nil {
					return nil, fmt.Errorf("source is required")
				}
				res.Datasets = append(res.Datasets, DatasetReference{Source: *source, Name: name})
			} else {
				return nil, fmt.Errorf("dataset must be a string or a mapping")
			}
		}
	}

	return res, nil
}
