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
	"path"
	"path/filepath"
	"regexp"

	"github.com/cloudquery/plugin-sdk/v4/plugin"
	"github.com/spf13/cobra"
)

const (
	pluginPublishShort = "Package plugin for publishing to CloudQuery registry."
	pluginPublishLong  = `Package plugin for publishing to CloudQuery registry.

This creates a directory with the plugin binaries, package.json and documentation.
`
)

// PackageJSON is the package.json file inside the dist directory. It is used by the CloudQuery package command
// to be able to package the plugin with all the needed metadata.
type PackageJSON struct {
	Name             string               `json:"name"`
	Version          string               `json:"version"`
	Protocols        []int                `json:"protocols"`
	SupportedTargets []plugin.BuildTarget `json:"supported_targets"`
	PackageType      plugin.PackageType   `json:"package_type"`
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

func (s *PluginServe) build(pluginDirectory, goos, goarch, distPath, pluginVersion string) error {
	pluginName := fmt.Sprintf("plugin-%s-%s-%s-%s", s.plugin.Name(), pluginVersion, goos, goarch)
	pluginPath := path.Join(distPath, pluginName)
	args := []string{"build", "-o", pluginPath}
	importPath, err := s.getModuleName(pluginDirectory)
	if err != nil {
		return err
	}
	args = append(args, "-buildmode=exe")
	args = append(args, "-ldflags", fmt.Sprintf("-s -w -X %s/plugin.Version=%s", importPath, pluginVersion))
	cmd := exec.Command("go", args...)
	cmd.Dir = pluginDirectory
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
	_, err = io.Copy(pluginZip, pluginFile)
	if err != nil {
		return fmt.Errorf("failed to copy plugin file to zip archive: %w", err)
	}
	if err := pluginFile.Close(); err != nil {
		return err
	}
	if err := os.Remove(pluginPath); err != nil {
		return fmt.Errorf("failed to remove plugin file: %w", err)
	}
	return nil
}

func (*PluginServe) getModuleName(pluginDirectory string) (string, error) {
	goMod, err := os.ReadFile(path.Join(pluginDirectory, "go.mod"))
	if err != nil {
		return "", fmt.Errorf("failed to open go.mod: %w", err)
	}
	reMod := regexp.MustCompile(`module\s+(.+)\n`)
	importPathMatches := reMod.FindStringSubmatch(string(goMod))
	if len(importPathMatches) != 2 {
		return "", fmt.Errorf("failed to parse import path from go.mod")
	}
	importPath := importPathMatches[1]
	if err != nil {
		return "", fmt.Errorf("failed to get import path: %w", err)
	}
	return importPath, nil
}

func (s *PluginServe) writePackageJSON(dir, pluginVersion string) error {
	packageJSON := PackageJSON{
		Name:             s.plugin.Name(),
		Version:          pluginVersion,
		Protocols:        s.versions,
		SupportedTargets: s.plugin.Targets(),
		PackageType:      plugin.PackageTypeNative,
	}
	buffer := &bytes.Buffer{}
	m := json.NewEncoder(buffer)
	m.SetIndent("", "  ")
	m.SetEscapeHTML(false)
	err := m.Encode(packageJSON)
	if err != nil {
		return err
	}
	outputPath := filepath.Join(dir, "package.json")
	return os.WriteFile(outputPath, buffer.Bytes(), 0644)
}

func (s *PluginServe) newCmdPluginPackage() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "package <plugin_directory> <version>",
		Short: pluginPublishShort,
		Long:  pluginPublishLong,
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			pluginDirectory := args[0]
			pluginVersion := args[1]
			distPath := path.Join(pluginDirectory, "dist")
			if cmd.Flag("dist-dir").Changed {
				distPath = cmd.Flag("dist-dir").Value.String()
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
				if err := s.build(pluginDirectory, target.OS, target.Arch, distPath, pluginVersion); err != nil {
					return fmt.Errorf("failed to build plugin for %s/%s: %w", target.OS, target.Arch, err)
				}
			}
			if err := s.writePackageJSON(distPath, pluginVersion); err != nil {
				return fmt.Errorf("failed to write manifest: %w", err)
			}
			return nil
		},
	}
	cmd.Flags().String("dist-dir", "", "dist directory to output the built plugin. (default: <plugin_directory>/dist)")
	return cmd
}
