package main

type ReleaserConfig struct {
	Builds   []ReleaserBuild   `yaml:"builds"`
	Archives []ReleaserArchive `yaml:"archives"`
}

type ReleaserBuild struct {
	ID   string   `yaml:"id"`
	GOOS []string `yaml:"goos"`
}

type ReleaserArchive struct {
	ID           string   `yaml:"id"`
	Format       string   `yaml:"format"`
	Builds       []string `yaml:"builds"`
	NameTemplate string   `yaml:"name_template"`
}

type ReleaserFormatOverride struct {
	GOOS   string `yaml:"goos"`
	Format string `yaml:"format"`
}
