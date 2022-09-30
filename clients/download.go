package clients

import (
	"archive/zip"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

type PluginType string

const (
	PluginTypeSource      PluginType = "source"
	PluginTypeDestination PluginType = "destination"
)

func DownloadPluginFromGithub(ctx context.Context, localPath string, githubPath string, version string, typ PluginType, writers ...io.Writer) error {
	pathSplit := strings.Split(githubPath, "/")
	if len(pathSplit) != 2 {
		return fmt.Errorf("invalid github path. should be in format: owner/repo")
	}
	org, name := pathSplit[0], pathSplit[1]
	downloadDir := filepath.Dir(localPath)
	pluginZipPath := localPath + ".zip"
	// https://github.com/cloudquery/cloudquery/releases/download/plugins-source-test-v1.1.5/test_darwin_amd64.zip
	downloadURL := fmt.Sprintf("https://github.com/cloudquery/cloudquery/releases/download/plugins-%s-%s-%s/%s_%s_%s.zip", typ, name, version, name, runtime.GOOS, runtime.GOARCH)
	if org != "cloudquery" {
		// https://github.com/yevgenypats/cq-source-test/releases/download/v1.0.1/cq-source-test_darwin_amd64.zip
		downloadURL = fmt.Sprintf("https://github.com/%s/cq-%s-%s/releases/download/%s/cq-%s-%s_%s_%s.zip", org, typ, name, version, typ, name, runtime.GOOS, runtime.GOARCH)
	}

	if _, err := os.Stat(localPath); err == nil {
		return nil
	}

	if err := os.MkdirAll(downloadDir, 0755); err != nil {
		return fmt.Errorf("failed to create plugin directory %s: %w", downloadDir, err)
	}

	err := downloadFile(ctx, pluginZipPath, downloadURL, writers...)
	if err != nil {
		return fmt.Errorf("failed to download plugin: %w", err)
	}

	archive, err := zip.OpenReader(pluginZipPath)
	if err != nil {
		return fmt.Errorf("failed to open plugin archive: %w", err)
	}

	pathInArchive := fmt.Sprintf("plugins/%s/%s", typ, name)
	if org != "cloudquery" {
		pathInArchive = fmt.Sprintf("cq-%s-%s", typ, name)
	}

	fileInArchive, err := archive.Open(pathInArchive)
	if err != nil {
		return fmt.Errorf("failed to open plugin archive plugins/source/%s: %w", name, err)
	}
	out, err := os.OpenFile(localPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0744)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", localPath, err)
	}
	_, err = io.Copy(out, fileInArchive)
	if err != nil {
		return fmt.Errorf("failed to copy body to file: %w", err)
	}
	err = out.Close()
	if err != nil {
		return fmt.Errorf("failed to close file: %w", err)
	}
	return nil
}

func downloadFile(ctx context.Context, localPath string, url string, writers ...io.Writer) (err error) {
	// Create the file
	out, err := os.Create(localPath)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", localPath, err)
	}
	defer out.Close()

	// Get the data
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("failed create request %s: %w", url, err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to get url %s: %w", url, err)
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s. downloading %s", resp.Status, url)
	}
	var w []io.Writer
	w = append(w, out)
	w = append(w, writers...)
	
	// Writer the body to file
	_, err = io.Copy(io.MultiWriter(w...), resp.Body)
	if err != nil {
		return fmt.Errorf("failed to copy body to file %s: %w", localPath, err)
	}

	return nil
}
