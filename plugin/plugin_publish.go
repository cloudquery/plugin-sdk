package plugin

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
)

const (
	goos_linux = "linux"
	goos_windows = "windows"
	goos_darwin = "darwin"

	goarch_amd64 = "amd64"
	goarch_arm64 = "arm64"
)

type PackageType string

const (
	PackageTypeNative PackageType = "native"
	PackageTypeDocker PackageType = "docker"
)

type BuildTarget struct {
	OS string `json:"os"`
	Arch string `json:"arch"`
}

// manifest is the plugin.json file inside the dist directory. It is used by CloudQuery registry
// to be able to publish the plugin with all the needed metadata.
type Manifest struct {
	Name string `json:"name"`
	Version string `json:"version"`
	Title string `json:"title"`
	ShortDescription string `json:"short_description"`
	Description string `json:"description"`
	Categories []string `json:"categories"`
	Protocols []int `json:"protocols"`
	SupportedTargets []BuildTarget `json:"supported_targets"`
	PackageType PackageType `json:"package_type"`
}

var buildTargets = []BuildTarget{
	{goos_linux, goarch_amd64},
	{goos_windows, goarch_amd64},
	{goos_darwin, goarch_amd64},
	{goos_darwin, goarch_arm64},
}

func isDirectoryExist(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}

func (p *Plugin) writeTablesJson(ctx context.Context, dir string) error {
	tables, err := p.Tables(ctx, TableOptions{
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

func (p *Plugin) writeManifest(ctx context.Context, dir string) error {
	manifest := Manifest{
		Name: p.Name(),
		Version: p.Version(),
		Title: p.title,
		ShortDescription: p.shortDescription,
		Description: p.description,
		Categories: p.categories,
		Protocols: []int{3},
		SupportedTargets: p.targets,
		PackageType: PackageTypeNative,
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

func (p *Plugin) Publish(ctx context.Context, pluginDirectory string) error {
	distPath := pluginDirectory + "/dist"
	if isDirectoryExist(distPath) {
		return fmt.Errorf("dist directory already exist: %s", distPath)
	}
	if err := os.MkdirAll(distPath, 0755); err != nil {
		return fmt.Errorf("failed to create dist directory: %w", err)
	}

	if err := p.writeTablesJson(ctx, distPath); err != nil {
		return fmt.Errorf("failed to write tables.json: %w", err)
	}

	for _, target := range p.targets {
		fmt.Println("Building for OS: " + target.OS + ", ARCH: " + target.Arch)
		if err := p.build(ctx, pluginDirectory, target.OS, target.Arch); err != nil {
			return fmt.Errorf("failed to build plugin for %s/%s: %w", target.OS, target.Arch, err)
		}
	}
	if err := p.writeManifest(ctx, distPath); err != nil {
		return fmt.Errorf("failed to write manifest: %w", err)
	}
	return nil
}



func (p *Plugin) build(ctx context.Context, pluginDirectory string, goos string, goarch string) error {
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