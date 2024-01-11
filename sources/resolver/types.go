package resolver

type Resolved struct {
	ID       string
	Name     string
	Endpoint string
}

type DatasetDescriptorHashes struct {
	SHA1   string `yaml:"sha1"`
	SHA256 string `yaml:"sha256"`
	MD5    string `yaml:"md5"`
}

type DatasetDescriptorVariants struct {
	Name   string                  `yaml:"name"`
	URL    string                  `yaml:"url"`
	Hashes DatasetDescriptorHashes `yaml:"hashes"`
}

type DatasetDescriptorFiles struct {
	Default  string                               `yaml:"default"`
	Variants map[string]DatasetDescriptorVariants `yaml:"variants"`
}

type DatasetDescriptor struct {
	Name        string                 `yaml:"name"`
	URL         string                 `yaml:"url"`
	Description string                 `yaml:"description"`
	Lisence     string                 `yaml:"lisence"`
	Files       DatasetDescriptorFiles `yaml:"files"`
}
