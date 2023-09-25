package serve

import (
	"archive/zip"
	"bytes"
	"context"
	"crypto/sha256"
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
	Kind             plugin.Kind        `json:"kind"`
	Message          string             `json:"message"`
	Version          string             `json:"version"`
	Protocols        []int              `json:"protocols"`
	SupportedTargets []TargetBuild      `json:"supported_targets"`
	PackageType      plugin.PackageType `json:"package_type"`
}

type TargetBuild struct {
	OS       string `json:"os"`
	Arch     string `json:"arch"`
	Path     string `json:"path"`
	Checksum string `json:"checksum"`
}

// This is the structure the CLI publish command expects
type pluginTable struct {
	Description   string    `json:"description,omitempty"`
	IsIncremental bool      `json:"is_incremental,omitempty"`
	Name          string    `json:"name,omitempty"`
	Parent        *string   `json:"parent,omitempty"`
	Relations     *[]string `json:"relations,omitempty"`
	Title         string    `json:"title,omitempty"`
}

func (s *PluginServe) writeTablesJSON(ctx context.Context, dir string) error {
	tables, err := s.plugin.Tables(ctx, plugin.TableOptions{
		Tables: []string{"*"},
	})
	if err != nil {
		return err
	}
	flattenedTables := tables.FlattenTables()
	tablesToEncode := make([]pluginTable, 0, len(flattenedTables))
	for _, t := range flattenedTables {
		table := tables.Get(t.Name)
		var parent *string
		if table.Parent != nil {
			parent = &table.Parent.Name
		}
		var relations *[]string
		if table.Relations != nil {
			names := table.Relations.TableNames()
			relations = &names
		}
		tablesToEncode = append(tablesToEncode, pluginTable{
			Description:   table.Description,
			IsIncremental: table.IsIncremental,
			Name:          table.Name,
			Parent:        parent,
			Relations:     relations,
			Title:         table.Title,
		})
	}
	buffer := &bytes.Buffer{}
	m := json.NewEncoder(buffer)
	m.SetIndent("", "")
	m.SetEscapeHTML(false)
	err = m.Encode(tablesToEncode)
	if err != nil {
		return err
	}
	outputPath := filepath.Join(dir, "tables.json")
	return os.WriteFile(outputPath, buffer.Bytes(), 0644)
}

func (s *PluginServe) build(pluginDirectory, goos, goarch, distPath, pluginVersion string) (*TargetBuild, error) {
	pluginName := fmt.Sprintf("plugin-%s-%s-%s-%s", s.plugin.Name(), pluginVersion, goos, goarch)
	pluginPath := path.Join(distPath, pluginName)
	args := []string{"build", "-o", pluginPath}
	importPath, err := s.getModuleName(pluginDirectory)
	if err != nil {
		return nil, err
	}
	args = append(args, "-buildmode=exe")
	args = append(args, "-ldflags", fmt.Sprintf("-s -w -X %s/plugin.Version=%s", importPath, pluginVersion))
	cmd := exec.Command("go", args...)
	cmd.Dir = pluginDirectory
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, fmt.Sprintf("GOOS=%s", goos))
	cmd.Env = append(cmd.Env, fmt.Sprintf("GOARCH=%s", goarch))
	cmd.Env = append(cmd.Env, fmt.Sprintf("CGO_ENABLED=%v", getEnvOrDefault("CGO_ENABLED", "0"))) // default to CGO_ENABLED=0
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("failed to build plugin with `go %v`: %w", args, err)
	}

	pluginFile, err := os.Open(pluginPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open plugin file: %w", err)
	}
	defer pluginFile.Close()

	zipPluginPath := pluginPath + ".zip"
	zipPluginFile, err := os.Create(zipPluginPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create zip file: %w", err)
	}
	defer zipPluginFile.Close()

	zipWriter := zip.NewWriter(zipPluginFile)
	defer zipWriter.Close()

	pluginZip, err := zipWriter.Create(pluginName)
	if err != nil {
		zipWriter.Close()
		return nil, fmt.Errorf("failed to create file in zip archive: %w", err)
	}
	_, err = io.Copy(pluginZip, pluginFile)
	if err != nil {
		zipWriter.Close()
		return nil, fmt.Errorf("failed to copy plugin file to zip archive: %w", err)
	}
	err = zipWriter.Close()
	if err != nil {
		return nil, fmt.Errorf("failed to close zip archive: %w", err)
	}

	if err := pluginFile.Close(); err != nil {
		return nil, err
	}
	if err := os.Remove(pluginPath); err != nil {
		return nil, fmt.Errorf("failed to remove plugin file: %w", err)
	}

	targetZip := fmt.Sprintf(pluginName + ".zip")
	checksum, err := calcChecksum(path.Join(distPath, targetZip))
	if err != nil {
		return nil, fmt.Errorf("failed to calculate checksum: %w", err)
	}

	return &TargetBuild{
		OS:       goos,
		Arch:     goarch,
		Path:     targetZip,
		Checksum: "sha256:" + checksum,
	}, nil
}

func calcChecksum(p string) (string, error) {
	// calculate SHA-256 checksum
	f, err := os.Open(p)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer f.Close()
	hash := sha256.New()
	if _, err := io.Copy(hash, f); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", hash.Sum(nil)), nil
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

func (s *PluginServe) writePackageJSON(dir string, pluginKind plugin.Kind, pluginVersion, message string, targets []TargetBuild) error {
	packageJSON := PackageJSON{
		SchemaVersion:    1,
		Name:             s.plugin.Name(),
		Message:          message,
		Kind:             pluginKind,
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
		Use:   "package -m <message> <source|destination> <version> <plugin_directory>",
		Short: pluginPackageShort,
		Long:  pluginPackageLong,
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			pluginKind := plugin.Kind(args[0])
			if err := pluginKind.Validate(); err != nil {
				return err
			}
			pluginVersion := args[1]
			pluginDirectory := args[2]
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
			if pluginKind == plugin.KindSource {
				if err := s.writeTablesJSON(cmd.Context(), distPath); err != nil {
					return err
				}
			}
			targets := []TargetBuild{}
			for _, target := range s.plugin.Targets() {
				fmt.Println("Building for OS: " + target.OS + ", ARCH: " + target.Arch)
				targetBuild, err := s.build(pluginDirectory, target.OS, target.Arch, distPath, pluginVersion)
				if err != nil {
					return fmt.Errorf("failed to build plugin for %s/%s: %w", target.OS, target.Arch, err)
				}
				targets = append(targets, *targetBuild)
			}
			if err := s.writePackageJSON(distPath, pluginKind, pluginVersion, message, targets); err != nil {
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
