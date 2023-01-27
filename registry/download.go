package registry

import (
	"archive/zip"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/schollz/progressbar/v3"
)

type PluginType string

const (
	PluginTypeSource      PluginType = "source"
	PluginTypeDestination PluginType = "destination"
	DefaultDownloadDir               = ".cq"
	RetryAttempts                    = 5
	RetryWaitTime                    = 1 * time.Second
)

func DownloadPluginFromGithub(ctx context.Context, localPath string, org string, name string, version string, typ PluginType) error {
	downloadDir := filepath.Dir(localPath)
	pluginZipPath := localPath + ".zip"
	// https://github.com/cloudquery/cloudquery/releases/download/plugins-source-test-v1.1.5/test_darwin_amd64.zip
	var urls []string
	if org == "cloudquery" {
		urls = []string{
			// monorepo plugin
			fmt.Sprintf("https://github.com/cloudquery/cloudquery/releases/download/plugins-%s-%s-%s/%s_%s_%s.zip", typ, name, version, name, runtime.GOOS, runtime.GOARCH),
			// community plugin under CloudQuery org
			fmt.Sprintf("https://github.com/cloudquery/cq-%s-%s/releases/download/%s/cq-%s-%s_%s_%s.zip", typ, name, version, typ, name, runtime.GOOS, runtime.GOARCH),
		}
	} else {
		// https://github.com/yevgenypats/cq-source-test/releases/download/v1.0.1/cq-source-test_darwin_amd64.zip
		urls = []string{
			// community plugin under user org
			fmt.Sprintf("https://github.com/%s/cq-%s-%s/releases/download/%s/cq-%s-%s_%s_%s.zip", org, typ, name, version, typ, name, runtime.GOOS, runtime.GOARCH),
		}
	}

	if _, err := os.Stat(localPath); err == nil {
		return nil
	}

	if err := os.MkdirAll(downloadDir, 0755); err != nil {
		return fmt.Errorf("failed to create plugin directory %s: %w", downloadDir, err)
	}

	urlIndex, err := downloadFile(ctx, pluginZipPath, urls...)
	if err != nil {
		return fmt.Errorf("failed to download plugin: %w", err)
	}

	archive, err := zip.OpenReader(pluginZipPath)
	if err != nil {
		return fmt.Errorf("failed to open plugin archive: %w", err)
	}
	defer archive.Close()

	var pathInArchive string
	if org == "cloudquery" && urlIndex == 0 {
		pathInArchive = fmt.Sprintf("plugins/%s/%s", typ, name)
	} else {
		pathInArchive = fmt.Sprintf("cq-%s-%s", typ, name)
	}
	pathInArchive = WithBinarySuffix(pathInArchive)

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

func downloadFile(ctx context.Context, localPath string, urls ...string) (urlIndex int, err error) {
	// Create the file
	out, err := os.Create(localPath)
	if err != nil {
		return 0, fmt.Errorf("failed to create file %s: %w", localPath, err)
	}
	defer out.Close()

	for r := 0; r < RetryAttempts; r++ {
		for i, url := range urls {
			// Get the data
			req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
			if err != nil {
				return 0, fmt.Errorf("failed create request %s: %w", url, err)
			}

			// Do http request
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				return 0, fmt.Errorf("failed to get url %s: %w", url, err)
			}
			// Check server response
			if resp.StatusCode == http.StatusNotFound && i < len(urls)-1 {
				// check alternative url
				resp.Body.Close()
				continue
			} else if resp.StatusCode != http.StatusOK {
				fmt.Printf("Failed downloading %s with status code %d. Retrying\n", url, resp.StatusCode)
				resp.Body.Close()
				break
			}

			fmt.Printf("Downloading %s\n", url)
			bar := downloadProgressBar(resp.ContentLength, "Downloading")

			// Writer the body to file
			_, err = io.Copy(io.MultiWriter(out, bar), resp.Body)
			if err != nil {
				return 0, fmt.Errorf("failed to copy body to file %s: %w", localPath, err)
			}
			resp.Body.Close()
			return i, nil
		}
		time.Sleep(RetryWaitTime)
	}

	return 0, errors.New("failed to download plugin")
}

func downloadProgressBar(maxBytes int64, description ...string) *progressbar.ProgressBar {
	desc := ""
	if len(description) > 0 {
		desc = description[0]
	}
	return progressbar.NewOptions64(
		maxBytes,
		progressbar.OptionSetDescription(desc),
		progressbar.OptionSetWriter(os.Stdout),
		progressbar.OptionShowBytes(true),
		progressbar.OptionSetWidth(10),
		progressbar.OptionThrottle(65*time.Millisecond),
		progressbar.OptionShowCount(),
		progressbar.OptionOnCompletion(func() {
			fmt.Fprint(os.Stdout, "\n")
		}),
		progressbar.OptionSpinnerType(14),
		progressbar.OptionFullWidth(),
		progressbar.OptionSetRenderBlankState(true),
	)
}

func WithBinarySuffix(filePath string) string {
	if runtime.GOOS == "windows" {
		return filePath + ".exe"
	}
	return filePath
}
