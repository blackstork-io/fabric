package main

type Metadata struct {
	Plugins []*PluginMetadata `json:"plugins"`
}

type PluginMetadata struct {
	Name     string                   `json:"name"`
	Version  string                   `json:"version"`
	Archives []*PluginArchiveMetadata `json:"archives"`
}

type PluginArchiveMetadata struct {
	Filename       string `json:"filename"`
	OS             string `json:"os"`
	Arch           string `json:"arch"`
	BinaryChecksum string `json:"binary_checksum"`
}
