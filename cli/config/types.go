package config

type ZarfFile struct {
	Source     string
	Shasum     string
	Target     string
	Executable bool
}

type ZarfChart struct {
	Name    string
	Url     string
	Version string
}

type ZarfFeature struct {
	Name        string
	Description string
	Default     bool
	Manifests   string
	Images      []string
	Files       []ZarfFile
	Charts      []ZarfChart
}

type ZarfMetatdata struct {
	Name         string
	Description  string
	Version      string
	Uncompressed bool
}

type ZarfContainerTarget struct {
	Namespace string
	Selector  string
	Container string
	Path      string
}

type ZarfData struct {
	Source string
	Target ZarfContainerTarget
}

type ZarfConfig struct {
	Kind     string
	Metadata ZarfMetatdata
	Features []ZarfFeature
	Data     []ZarfData
	Local    struct {
		Manifests string
		Images    []string
		Files     []ZarfFile
		Charts    []ZarfChart
	}
	Remote struct {
		Images []string
		Repos  []string
	}
	Package struct {
		Terminal   string
		User       string
		Timestamp string
	}
}
