package json

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

type testFS struct {
	FS   fs.FS
	path string
}

func makeTestFS(tb testing.TB) testFS {
	tb.Helper()

	path, err := filepath.EvalSymlinks(tb.TempDir())
	if err != nil {
		tb.Fatalf("failed to create testFS: %s", err)
	}

	path = filepath.ToSlash(path)

	tb.Logf("creating testFS at %s", path)
	return testFS{
		FS:   os.DirFS(path),
		path: path,
	}
}

func (t testFS) Open(name string) (fs.File, error) {
	return t.FS.Open(filepath.ToSlash(name))
}

func (t testFS) Path() string {
	return t.path
}

func (t testFS) WriteFile(name string, data []byte, perm os.FileMode) error {
	name = filepath.ToSlash(name)
	if filepath.IsAbs(name) {
		if strings.HasPrefix(name, t.path) {
			return os.WriteFile(name, data, perm)
		}
		return fmt.Errorf("path is outside test fs root folder")
	}
	return os.WriteFile(filepath.ToSlash(filepath.Join(t.path, name)), data, perm)
}

func (t testFS) MkdirAll(path string, perm os.FileMode) error {
	path = filepath.ToSlash(path)
	if filepath.IsAbs(path) {
		if strings.HasPrefix(path, t.path) {
			return os.MkdirAll(path, perm)
		}
		return fmt.Errorf("path is outside test fs root folder")
	}
	return os.MkdirAll(filepath.ToSlash(filepath.Join(t.path, path)), perm)
}

func testJSON(m any) json.RawMessage {
	contents, err := json.Marshal(m)
	if err != nil {
		panic(err)
	}
	return contents
}
