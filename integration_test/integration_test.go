//go:build integration

package integration_test

import (
	"log"
	"os"
	"path/filepath"
	"runtime"
)

var pluginsDir = func() (pluginsDir string) {
	_, filename, _, _ := runtime.Caller(0)
	prevPath, err := filepath.Abs(filename)
	if err != nil {
		log.Fatalf("Failed to get working dir: %s", err)
	}
	path := filepath.Clean(filepath.Join(prevPath, ".."))
	for prevPath != path {
		_, err := os.Stat(filepath.Join(path, "go.mod"))
		if err == nil {
			pluginsDir = filepath.Join(path, "dist", "plugins")
			return
		}
		prevPath = path
		path = filepath.Clean(filepath.Join(path, ".."))
	}
	log.Fatal("Failed to find go.mod in parent folders")
	return
}()
