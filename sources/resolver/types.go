package resolver

type Resolved struct {
	ID       string
	Endpoint string
	Location *string
	Format   string
}

type DatasetDescriptorHashes struct {
	SHA1   string `yaml:"sha1"`
	SHA256 string `yaml:"sha256"`
	MD5    string `yaml:"md5"`
}

type DatasetDescriptorData struct {
	URL    string                  `yaml:"url"`
	Hashes DatasetDescriptorHashes `yaml:"hashes"`
	Format *string                 `yaml:"format"`
}

type DatasetDescriptorExtra struct {
	ID          string                `yaml:"id"`
	Name        string                `yaml:"name"`
	Description string                `yaml:"description"`
	Dataset     DatasetDescriptorData `yaml:"dataset"`
}

type DatasetDescriptor struct {
	ID          string                            `yaml:"id"`
	Name        string                            `yaml:"name"`
	URL         string                            `yaml:"url"`
	Description string                            `yaml:"description"`
	Lisence     string                            `yaml:"lisence"`
	Dataset     DatasetDescriptorData             `yaml:"dataset"`
	Extras      map[string]DatasetDescriptorExtra `yaml:"extras"`
}
