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
	"strings"

	"github.com/cloudquery/plugin-sdk/v4/plugin"
	"github.com/spf13/cobra"
)

const (
	pluginPackageShort = "Package plugin for publishing to CloudQuery registry."
	pluginPackageLong  = `Package plugin for publishing to CloudQuery registry.

This creates a directory with the plugin binaries, package.json and documentation.
`
)

// PackageJSON is the package.json file inside the dist directory. It is used by the CloudQuery package command
// to be able to package the plugin with all the needed metadata.
type PackageJSON struct {
	SchemaVersion    int                `json:"schema_version"`
	Name             string             `json:"name"`
	Message          string             `json:"message"`
	Version          string             `json:"version"`
	Protocols        []int              `json:"protocols"`
	SupportedTargets []TargetBuild      `json:"supported_targets"`
	PackageType      plugin.PackageType `json:"package_type"`
}

type TargetBuild struct {
	OS   string `json:"os"`
	Arch string `json:"arch"`
	Path string `json:"path"`
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
	return strings.TrimSpace(importPath), nil
}

func (s *PluginServe) writePackageJSON(dir, pluginVersion, message string) error {
	targets := []TargetBuild{}
	for _, target := range s.plugin.Targets() {
		pluginName := fmt.Sprintf("plugin-%s-%s-%s-%s", s.plugin.Name(), pluginVersion, target.OS, target.Arch)
		targets = append(targets, TargetBuild{
			OS:   target.OS,
			Arch: target.Arch,
			Path: pluginName + ".zip",
		})
	}
	packageJSON := PackageJSON{
		SchemaVersion:    1,
		Name:             s.plugin.Name(),
		Message:          message,
		Version:          pluginVersion,
		Protocols:        s.versions,
		SupportedTargets: targets,
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

func (*PluginServe) copyDocs(distPath, docsPath string) error {
	err := os.MkdirAll(filepath.Join(distPath, "docs"), 0755)
	if err != nil {
		return err
	}
	dirEntry, err := os.ReadDir(docsPath)
	if err != nil {
		return err
	}
	for _, entry := range dirEntry {
		if entry.IsDir() {
			continue
		}
		if strings.HasSuffix(entry.Name(), ".md") {
			src := filepath.Join(docsPath, entry.Name())
			dst := filepath.Join(distPath, "docs", entry.Name())
			err := copyFile(src, dst)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func copyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()
	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()
	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return err
	}
	return nil
}

func (s *PluginServe) newCmdPluginPackage() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "package -m <message> <plugin_directory> <version>",
		Short: pluginPackageShort,
		Long:  pluginPackageLong,
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			pluginDirectory := args[0]
			pluginVersion := args[1]
			distPath := path.Join(pluginDirectory, "dist")
			if cmd.Flag("dist-dir").Changed {
				distPath = cmd.Flag("dist-dir").Value.String()
			}
			docsPath := path.Join(pluginDirectory, "docs")
			if cmd.Flag("docs-dir").Changed {
				docsPath = cmd.Flag("docs-dir").Value.String()
			}
			message := ""
			if !cmd.Flag("message").Changed {
				return fmt.Errorf("message is required")
			}
			message = cmd.Flag("message").Value.String()
			if strings.HasPrefix(message, "@") {
				messageFile := strings.TrimPrefix(message, "@")
				messageBytes, err := os.ReadFile(messageFile)
				if err != nil {
					return err
				}
				message = string(messageBytes)
			}
			message = normalizeMessage(message)

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
			if err := s.writePackageJSON(distPath, pluginVersion, message); err != nil {
				return fmt.Errorf("failed to write manifest: %w", err)
			}
			if err := s.copyDocs(distPath, docsPath); err != nil {
				return fmt.Errorf("failed to copy docs: %w", err)
			}
			return nil
		},
	}
	cmd.Flags().StringP("dist-dir", "D", "", "dist directory to output the built plugin. (default: <plugin_directory>/dist)")
	cmd.Flags().StringP("docs-dir", "", "", "docs directory containing markdown files to copy to the dist directory. (default: <plugin_directory>/docs)")
	cmd.Flags().StringP("message", "m", "", "message that summarizes what is new or changed in this version. Use @<file> to read from file. Supports markdown.")
	return cmd
}

func normalizeMessage(s string) string {
	s = strings.TrimSpace(s)
	s = strings.ReplaceAll(s, "\r\n", "\n")
	s = strings.ReplaceAll(s, "\r", "\n")
	return s
}
