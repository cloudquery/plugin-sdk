package serve

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/cloudquery/plugin-sdk/v4/plugin"
	"github.com/spf13/cobra"
)

const (
	pluginPublishShort = "Publish plugin to CloudQuery registry"
	pluginPublishLong  = `Publish plugin to CloudQuery registry

To just build the plugin without publishing, use the --dry-run flag.
Example:
go run main.go publish --dry-run
`
)

type PackageType string

const (
	PackageTypeNative PackageType = "native"
	PackageTypeDocker PackageType = "docker"
)

// manifest is the plugin.json file inside the dist directory. It is used by CloudQuery registry
// to be able to publish the plugin with all the needed metadata.
type Manifest struct {
	Name             string               `json:"name"`
	Version          string               `json:"version"`
	Title            string               `json:"title"`
	ShortDescription string               `json:"short_description"`
	Description      string               `json:"description"`
	Categories       []string             `json:"categories"`
	Protocols        []int                `json:"protocols"`
	SupportedTargets []plugin.BuildTarget `json:"supported_targets"`
	PackageType      PackageType          `json:"package_type"`
}

func isDirectoryExist(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}

func (s *PluginServe) writeTablesJSON(ctx context.Context, dir string) error {
	tables, err := s.plugin.Tables(ctx, plugin.TableOptions{
		Tables: []string{"*"},
	})
	if err != nil {
		return err
	}
	buffer := &bytes.Buffer{}
	m := json.NewEncoder(buffer)
	m.SetIndent("", "  ")
	m.SetEscapeHTML(false)
	err = m.Encode(tables)
	if err != nil {
		return err
	}
	outputPath := filepath.Join(dir, "tables.json")
	return os.WriteFile(outputPath, buffer.Bytes(), 0644)
}

func (*PluginServe) build(pluginDirectory string, goos string, goarch string) error {
	pluginName := "plugin" + "_" + goos + "_" + goarch
	distPath := pluginDirectory + "/dist"

	pluginPath := distPath + "/" + pluginName
	args := []string{"build", "-C", pluginDirectory, "-o", pluginPath}
	cmd := exec.Command("go", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ()
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to build plugin with `go %v`: %w", args, err)
	}

	pluginFile, err := os.Open(pluginPath)
	if err != nil {
		return fmt.Errorf("failed to open plugin file: %w", err)
	}
	defer pluginFile.Close()

	zipPluginPath := pluginPath + ".zip"
	zipPluginFile, err := os.Create(zipPluginPath)
	if err != nil {
		return fmt.Errorf("failed to create zip file: %w", err)
	}
	defer zipPluginFile.Close()

	zipWriter := zip.NewWriter(zipPluginFile)
	defer zipWriter.Close()

	pluginZip, err := zipWriter.Create(pluginName)
	if err != nil {
		return fmt.Errorf("failed to create file in zip archive: %w", err)
	}
	if err := os.Remove(pluginPath); err != nil {
		return fmt.Errorf("failed to remove plugin file: %w", err)
	}
	_, err = io.Copy(pluginZip, pluginFile)
	if err != nil {
		return fmt.Errorf("failed to copy plugin file to zip archive: %w", err)
	}
	return nil
}

func (s *PluginServe) writeManifest(dir string) error {
	manifest := Manifest{
		Name:             s.plugin.Name(),
		Version:          s.plugin.Version(),
		Title:            s.plugin.Title(),
		ShortDescription: s.plugin.ShortDescription(),
		Description:      s.plugin.Description(),
		Categories:       s.plugin.Categories(),
		Protocols:        s.versions,
		SupportedTargets: s.plugin.Targets(),
		PackageType:      PackageTypeNative,
	}
	buffer := &bytes.Buffer{}
	m := json.NewEncoder(buffer)
	m.SetIndent("", "  ")
	m.SetEscapeHTML(false)
	err := m.Encode(manifest)
	if err != nil {
		return err
	}
	outputPath := filepath.Join(dir, "plugin.json")
	return os.WriteFile(outputPath, buffer.Bytes(), 0644)
}

func (s *PluginServe) newCmdPluginPublish() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "publish <plugin_directory>",
		Short: pluginPublishShort,
		Long:  pluginPublishLong,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			pluginDirectory := args[0]
			distPath := pluginDirectory + "/dist"
			if isDirectoryExist(distPath) {
				return fmt.Errorf("dist directory already exist: %s", distPath)
			}
			if err := os.MkdirAll(distPath, 0755); err != nil {
				return err
			}
			if err := s.plugin.Init(cmd.Context(), nil, plugin.NewClientOptions{
				NoConnection: true,
			}); err != nil {
				return err
			}
			if err := s.writeTablesJSON(cmd.Context(), distPath); err != nil {
				return err
			}
			for _, target := range s.plugin.Targets() {
				fmt.Println("Building for OS: " + target.OS + ", ARCH: " + target.Arch)
				if err := s.build(pluginDirectory, target.OS, target.Arch); err != nil {
					return fmt.Errorf("failed to build plugin for %s/%s: %w", target.OS, target.Arch, err)
				}
			}
			if err := s.writeManifest(distPath); err != nil {
				return fmt.Errorf("failed to write manifest: %w", err)
			}
			return nil
		},
	}
	return cmd
}
