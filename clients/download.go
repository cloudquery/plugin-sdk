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
	"time"

	"github.com/avast/retry-go/v4"

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

	err := downloadFile(ctx, pluginZipPath, downloadURL)
	if err != nil {
		return fmt.Errorf("failed to download plugin: %w", err)
	}

	archive, err := zip.OpenReader(pluginZipPath)
	if err != nil {
		return fmt.Errorf("failed to open plugin archive: %w", err)
	}
	defer archive.Close()

	pathInArchive := fmt.Sprintf("plugins/%s/%s", typ, name)
	if org != "cloudquery" {
		pathInArchive = fmt.Sprintf("cq-%s-%s", typ, name)
	}
	pathInArchive = withBinarySuffix(pathInArchive)

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

func downloadFile(ctx context.Context, localPath string, url string) (err error) {
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

	err = retry.Do(
		func() error {
			// Do http request
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				return fmt.Errorf("failed to get url %s: %w", url, err)
			}

			// Check server response
			if resp.StatusCode != http.StatusOK {
				fmt.Printf("Failed downloading %s with status code %d. Retrying\n", url, resp.StatusCode)
				return fmt.Errorf("statusCode != 200")
			}
			defer resp.Body.Close()

			fmt.Printf("Downloading %s\n", url)
			bar := downloadProgressBar(resp.ContentLength, "Downloading")

			// Writer the body to file
			_, err = io.Copy(io.MultiWriter(out, bar), resp.Body)
			if err != nil {
				return fmt.Errorf("failed to copy body to file %s: %w", localPath, err)
			}

			return nil
		},
		retry.RetryIf(func(err error) bool {
			return err.Error() == "statusCode != 200"
		}),
		retry.Attempts(RetryAttempts),
		retry.Delay(RetryWaitTime),
	)

	if err != nil {
		return fmt.Errorf("failed downloading: %s", url)
	}

	return nil
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

func withBinarySuffix(filePath string) string {
	if runtime.GOOS == "windows" {
		return filePath + ".exe"
	}
	return filePath
}
