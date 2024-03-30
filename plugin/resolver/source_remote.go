package resolver

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
	"time"
)

const (
	maxDownloadSize = 50 * 1024 * 1024 // 50MB
	downloadTimeout = 5 * time.Minute
	regAPITimeout   = 10 * time.Second
)

// RemoteSource is a plugin source that looks up plugins from a remote registry.
// The registry should implement the Fabric Registry API.
type RemoteSource struct {
	// BaseURL is the base URL of the registry.
	BaseURL string
	// DownloadDir is the directory where the plugins are downloaded.
	DownloadDir string
	// UserAgent is the http user agent to use for the requests.
	// Useful for debugging and statistics on the registry side.
	UserAgent string
}

// regVersion represents a version of a plugin in the registry
type regVersion struct {
	Version   Version       `json:"version"`
	Platforms []regPlatform `json:"platforms"`
}

// regPlatform represents an available platform for a plugin version in the registry
type regPlatform struct {
	OS   string `json:"os"`
	Arch string `json:"arch"`
}

// regDownloadInfo represents the download info for a specific platform of a plugin version in the registry
type regDownloadInfo struct {
	OS          string `json:"os"`
	Arch        string `json:"arch"`
	DownloadURL string `json:"download_url"`
}

// regError represents an error response.
type regError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// registryError implements the error interface.
func (err regError) Error() string {
	return fmt.Sprintf("[code=%s]: %s", err.Code, err.Message)
}

// Lookup returns the versions found of the plugin with the given name.
func (source RemoteSource) Lookup(ctx context.Context, name Name) ([]Version, error) {
	versions, err := source.fetchVersions(ctx, name)
	if err != nil {
		if rerr, ok := err.(regError); ok && rerr.Code == "not_found" {
			return nil, ErrPluginNotFound
		} else {
			return nil, fmt.Errorf("failed to lookup plugin versions in the registry: %w", err)
		}
	}
	var matches []Version
	for _, version := range versions {
		hasPlatform := slices.ContainsFunc(version.Platforms, func(platform regPlatform) bool {
			return platform.OS == runtime.GOOS && platform.Arch == runtime.GOARCH
		})
		if hasPlatform {
			matches = append(matches, version.Version)
		}
	}
	return matches, nil
}

// call makes a http request to the registry with the given timeout.
func (source RemoteSource) call(req *http.Request, timeout time.Duration) (*http.Response, error) {
	if source.UserAgent != "" {
		req.Header.Set("User-Agent", source.UserAgent)
	}
	client := &http.Client{
		Timeout: timeout,
	}
	return client.Do(req)
}

// decodeBody decodes the http response from the registry into the provided value.
func (source RemoteSource) decodeBody(resp *http.Response, v interface{}) error {
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var errResp struct {
			Error regError `json:"error"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&errResp); err != nil {
			return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
		}
		return errResp.Error
	}
	return json.NewDecoder(resp.Body).Decode(v)
}

// fetchVersions looks up the plugin versions in the registry.
func (source RemoteSource) fetchVersions(ctx context.Context, name Name) ([]regVersion, error) {
	url := fmt.Sprintf("%s/v1/plugins/%s/%s/versions", source.BaseURL, name.Namespace(), name.Short())
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	resp, err := source.call(req, regAPITimeout)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var respData struct {
		Versions []regVersion `json:"versions"`
	}
	if err := source.decodeBody(resp, &respData); err != nil {
		return nil, err
	}
	return respData.Versions, nil
}

// Resolve returns the binary path and checksum for the given plugin version.
func (source RemoteSource) Resolve(ctx context.Context, name Name, version Version, checksums []Checksum) (*ResolvedPlugin, error) {
	downloadInfo, err := source.fetchDownloadInfo(ctx, name, version)
	if err != nil {
		return nil, fmt.Errorf("failed to get plugin from the registry: %w", err)
	}
	return source.download(ctx, name, version, downloadInfo, checksums)
}

// fetchDownloadInfo resolves the download info for sthe given plugin version from the registry.
func (source RemoteSource) fetchDownloadInfo(ctx context.Context, name Name, version Version) (*regDownloadInfo, error) {
	url := fmt.Sprintf("%s/v1/plugins/%s/%s/%s/download/%s/%s",
		source.BaseURL,
		name.Namespace(),
		name.Short(),
		version.String(),
		runtime.GOOS,
		runtime.GOARCH,
	)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	resp, err := source.call(req, regAPITimeout)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var info regDownloadInfo
	if err := source.decodeBody(resp, &info); err != nil {
		return nil, err
	}
	return &info, nil
}

// fetchChecksums fetches the plugin checksums from the registry.
func (source RemoteSource) fetchChecksums(ctx context.Context, name Name, version Version) ([]Checksum, error) {
	url := fmt.Sprintf("%s/v1/plugins/%s/%s/%s/checksums", source.BaseURL, name.Namespace(), name.Short(), version.String())
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	resp, err := source.call(req, regAPITimeout)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var respData struct {
		Checksums []Checksum `json:"checksums"`
	}
	if err := source.decodeBody(resp, &respData); err != nil {
		return nil, err
	}
	return respData.Checksums, nil
}

// download downloads the plugin from the registry and returns the binary path and checksum.
func (source RemoteSource) download(ctx context.Context, name Name, version Version, info *regDownloadInfo, checksums []Checksum) (_ *ResolvedPlugin, err error) {
	// If the checksums are not provided it means plugin version is not locked and we need to fetch the checksums from the registry.
	if len(checksums) == 0 {
		var err error
		checksums, err = source.fetchChecksums(ctx, name, version)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch plugin checksums: %w", err)
		}
	}
	// make a http request to download the plugin
	req, err := http.NewRequestWithContext(ctx, "GET", info.DownloadURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create download request: %w", err)
	}
	req.Header.Set("Accept", "application/octet-stream")
	resp, err := source.call(req, downloadTimeout)
	if err != nil {
		return nil, fmt.Errorf("failed to download plugin: %w", err)
	}
	defer resp.Body.Close()
	// verify download response headers
	if err = source.verifyDownloadHeaders(resp); err != nil {
		return nil, err
	}
	// calculate checksum of the downloaded archive while writing to the file
	h := sha256.New()
	buf := io.TeeReader(resp.Body, h)
	// extract plugin while downloading without saving to disk
	binaryPath, checksumPath, err := source.extract(name, version, buf, checksums)
	if err != nil {
		return nil, fmt.Errorf("failed to extract plugin: %w", err)
	}
	// cleanup extracted files if there is an error during checksum verification
	defer func() {
		if err == nil {
			return
		}
		// if there is an error, remove extracted binary file
		os.Remove(binaryPath)
		os.Remove(checksumPath)
		// remove directory if it is empty
		entries, err := os.ReadDir(filepath.Dir(binaryPath))
		if err == nil && len(entries) == 0 {
			os.Remove(filepath.Dir(binaryPath))
		}
	}()
	// read remaining data from the response body to verify the checksum of the downloaded archive
	if _, err := io.Copy(io.Discard, buf); err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	// verify checksum of the downloaded archive
	sum := Checksum{
		Object: "archive",
		OS:     runtime.GOOS,
		Arch:   runtime.GOARCH,
		Sum:    h.Sum(nil),
	}
	if !sum.Match(checksums) {
		return nil, fmt.Errorf("invalid plugin archive checksum: '%s'", sum)
	}
	return &ResolvedPlugin{
		BinaryPath: binaryPath,
		Checksums:  checksums,
	}, nil
}

// verifyDownloadHeaders verifies the download response headers.
func (source RemoteSource) verifyDownloadHeaders(res *http.Response) error {
	// verify the download size
	if res.ContentLength > maxDownloadSize {
		return fmt.Errorf("plugin download size exceeds the limit, got = %d, expect < %d", res.ContentLength, maxDownloadSize)
	}
	disposition, params, err := mime.ParseMediaType(res.Header.Get("Content-Disposition"))
	if err != nil {
		return fmt.Errorf("failed to parse content disposition: %w", err)
	}
	if disposition != "attachment" {
		return fmt.Errorf("unsupported content disposition: %s", disposition)
	}
	fn := params["filename"]
	if fn == "" {
		return fmt.Errorf("missing filename in content disposition")
	}
	if !strings.HasSuffix(fn, ".tar.gz") {
		return fmt.Errorf("unsupported archive type: %s", fn)
	}
	return nil
}

// extract the plugin from the tar.gz file and returns the binary and checksum file path.
func (source RemoteSource) extract(name Name, version Version, archive io.Reader, checksums []Checksum) (binPath, sumPath string, err error) {
	read, err := gzip.NewReader(archive)
	if err != nil {
		return "", "", fmt.Errorf("failed to create gzip reader: %w", err)
	}
	defer read.Close()
	reader := tar.NewReader(read)
	var found *tar.Header
	for {
		header, err := reader.Next()
		if err == io.EOF {
			break
		}
		if header.Typeflag != tar.TypeReg {
			continue
		}
		if header.Name != fmt.Sprintf("%s@%s", name.Short(), version.String()) &&
			header.Name != fmt.Sprintf("%s@%s.exe", name.Short(), version.String()) {
			continue
		}
		found = header
		break
	}
	if found == nil {
		return "", "", fmt.Errorf("plugin binary not found in tar.gz file")
	}
	binaryPath := filepath.Join(source.DownloadDir, name.Namespace(), filepath.Base(found.Name))
	checksumPath := strings.TrimSuffix(binaryPath, ".exe") + "_checksums.txt"
	if err := os.MkdirAll(filepath.Dir(binaryPath), 0o755); err != nil {
		return "", "", fmt.Errorf("failed to create plugin directory: %w", err)
	}
	binaryFile, err := os.OpenFile(binaryPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o755)
	if err != nil {
		return "", "", fmt.Errorf("failed to create plugin file: %w", err)
	}
	// cleanup the downloaded binary on error
	defer func() {
		binaryFile.Close()
		if err != nil {
			// if there is an error, remove extracted binary file and checksum file
			os.Remove(binaryPath)
			// remove directory if it is empty
			entries, err := os.ReadDir(filepath.Dir(binaryPath))
			if err == nil && len(entries) == 0 {
				os.Remove(filepath.Dir(binaryPath))
			}
		}
	}()
	// calculate checksum of the plugin binary while writing to the file
	h := sha256.New()
	buf := io.MultiWriter(h, binaryFile)
	// write the plugin binary
	if _, err := io.Copy(buf, reader); err != nil { //nolint:gosec // lint issue not possible here
		return "", "", fmt.Errorf("failed to write plugin file: %w", err)
	}
	sum := Checksum{
		Object: "binary",
		OS:     runtime.GOOS,
		Arch:   runtime.GOARCH,
		Sum:    h.Sum(nil),
	}
	if !sum.Match(checksums) {
		return "", "", fmt.Errorf("invalid plugin binary checksum: '%s'", sum)
	}
	// Create checksums file to be used for the following installs when plugin is installed from the local source.
	checksumFile, err := os.Create(checksumPath)
	if err != nil {
		return "", "", fmt.Errorf("failed to create plugin meta file: %w", err)
	}
	// cleanup checksum file operation
	defer func() {
		checksumFile.Close()
		if err != nil { // if there is an error, remove checksum file
			os.Remove(checksumPath)
		}
	}()
	if err := encodeChecksums(checksumFile, checksums); err != nil {
		return "", "", fmt.Errorf("failed to write plugin meta file: %w", err)
	}
	return binaryPath, checksumPath, nil
}
