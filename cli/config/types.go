package config

type ZarfFile struct {
	Source     string   `yaml:"source"`
	Shasum     string   `yaml:"shasum,omitempty"`
	Target     string   `yaml:"target"`
	Executable bool     `yaml:"executable,omitempty"`
	Symlinks   []string `yaml:"symlinks,omitempty"`
	Template   bool     `yaml:"template,omitempty"`
}

type ZarfChart struct {
	Name    string `yaml:"name"`
	Url     string `yaml:"url"`
	Version string `yaml:"version"`
}

type ZarfComponent struct {
	Name          string               `yaml:"name"`
	Description   string               `yaml:"description,omitempty"`
	Default       bool                 `yaml:"default,omitempty"`
	Required      bool                 `yaml:"required,omitempty"`
	Files         []ZarfFile           `yaml:"files,omitempty"`
	ManifestsPath string               `yaml:"manifests,omitempty"`
	Images        []string             `yaml:"images,omitempty"`
	Charts        []ZarfChart          `yaml:"charts,omitempty"`
	Repos         []string             `yaml:"repos,omitempty"`
	Scripts       ZarfComponentScripts `yaml:"scripts,omitempty"`
}

type ZarfComponentScripts struct {
	Retry  bool     `yaml:"retry,omitempty"`
	Before []string `yaml:"before,omitempty"`
	After  []string `yaml:"after,omitempty"`
}

type ZarfMetadata struct {
	Name         string `yaml:"name,omitempty"`
	Description  string `yaml:"description,omitempty"`
	Version      string `yaml:"version,omitempty"`
	Uncompressed bool   `yaml:"uncompressed,omitempty"`
}

type ZarfContainerTarget struct {
	Namespace string `yaml:"namespace"`
	Selector  string `yaml:"selector"`
	Container string `yaml:"container,omitempty"`
	Path      string `yaml:"path"`
}

type ZarfData struct {
	Source string              `yaml:"source"`
	Target ZarfContainerTarget `yaml:"target"`
}

type ZarfBuildData struct {
	Terminal  string `yaml:"terminal"`
	User      string `yaml:"user"`
	Timestamp string `yaml:"timestamp"`
	Version   string `yaml:"string"`
}

type ZarfPackage struct {
	Kind       string          `yaml:"kind,omitempty"`
	Metadata   ZarfMetadata    `yaml:"metadata,omitempty"`
	Build      ZarfBuildData   `yaml:"build,omitempty"`
	Data       []ZarfData      `yaml:"data,omitempty"`
	Components []ZarfComponent `yaml:"components,omitempty"`
}

type ZarfState struct {
	Kind string `yaml:"kind"`
	TLS  struct {
		CertPublicPath  string `yaml:"certPublicPath"`
		CertPrivatePath string `yaml:"certPrivatePath"`
		Host            string `yaml:"host"`
	} `yaml:"tls"`
}
