package serve

import (
	"archive/zip"
	"bufio"
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	cloudquery_api "github.com/cloudquery/cloudquery-api-go"
	"github.com/cloudquery/plugin-sdk/v4/plugin"
	"github.com/cloudquery/plugin-sdk/v4/schema"
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
	Team             string             `json:"team"`
	Kind             plugin.Kind        `json:"kind"`
	Name             string             `json:"name"`
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

func (s *PluginServe) writeTablesJSON(ctx context.Context, dir string) error {
	tables, err := s.plugin.Tables(ctx, plugin.TableOptions{
		Tables: []string{"*"},
	})
	if err != nil {
		return err
	}
	flattenedTables := tables.FlattenTables()
	tablesToEncode := make([]cloudquery_api.PluginTableCreate, 0, len(flattenedTables))
	for _, t := range flattenedTables {
		table := tables.Get(t.Name)
		var parent *string
		if table.Parent != nil {
			parent = &table.Parent.Name
		}
		relations := make([]string, 0, len(table.Relations))
		if table.Relations != nil {
			for _, relation := range table.Relations {
				relations = append(relations, relation.Name)
			}
		}
		columns := make([]cloudquery_api.PluginTableColumn, 0, len(table.Columns))
		for _, column := range table.Columns {
			c := cloudquery_api.PluginTableColumn{
				Name:           column.Name,
				Description:    column.Description,
				Type:           column.Type.String(),
				IncrementalKey: column.IncrementalKey,
				NotNull:        column.NotNull,
				// PrimaryKey Will be set to true Under the following conditions:
				// 1. If the column is a `PrimaryKeyComponent`
				// 2. If the column is a `PrimaryKey` and both of the following are true column name is NOT `_cq_id`  and there are other columns that are a PrimaryKeyComponent
				PrimaryKey: (column.PrimaryKey && !(column.Name == schema.CqIDColumn.Name && len(table.PrimaryKeyComponents()) > 0)) || column.PrimaryKeyComponent, //nolint:staticcheck
				Unique:     column.Unique,
			}
			if column.TypeSchema != "" {
				typeSchema := column.TypeSchema
				c.TypeSchema = &typeSchema
			}
			columns = append(columns, c)
		}
		tablesToEncode = append(tablesToEncode, cloudquery_api.PluginTableCreate{
			Description:       &table.Description,
			IsIncremental:     &table.IsIncremental,
			IsPaid:            &table.IsPaid,
			Name:              table.Name,
			Parent:            parent,
			Relations:         &relations,
			Title:             &table.Title,
			Columns:           &columns,
			PermissionsNeeded: &table.PermissionsNeeded,
			SensitiveColumns:  &table.SensitiveColumns,
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

func (s *PluginServe) build(pluginDirectory string, target plugin.BuildTarget, distPath, pluginVersion string) (*TargetBuild, error) {
	pluginFileName := fmt.Sprintf("plugin-%s-%s-%s-%s", s.plugin.Name(), pluginVersion, target.OS, target.Arch)
	pluginPath := path.Join(distPath, pluginFileName)
	importPath, err := s.getModuleName(pluginDirectory)
	if err != nil {
		return nil, err
	}
	stripSymbols := "-s "
	if target.IncludeSymbols {
		stripSymbols = ""
	}
	ldFlags := fmt.Sprintf("%[1]s -w -X %[2]s/plugin.Version=%[3]s -X %[2]s/resources/plugin.Version=%[3]s", stripSymbols, importPath, pluginVersion)
	args := []string{"build", "-trimpath", "-buildvcs=false", "-mod=readonly", "-o", pluginPath, "-buildmode=exe", "-ldflags", ldFlags}
	cmd := exec.Command("go", args...)
	cmd.Dir = pluginDirectory
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = append(os.Environ(), target.EnvVariables()...)
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

	pluginZip, err := zipWriter.Create(pluginFileName)
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

	targetZip := pluginFileName + ".zip"
	checksum, err := calcChecksum(path.Join(distPath, targetZip))
	if err != nil {
		return nil, fmt.Errorf("failed to calculate checksum: %w", err)
	}

	return &TargetBuild{
		OS:       target.OS,
		Arch:     target.Arch,
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
		return "", errors.New("failed to parse import path from go.mod")
	}
	importPath := importPathMatches[1]
	return strings.TrimSpace(importPath), nil
}

func (s *PluginServe) writePackageJSON(dir, version, message string, targets []TargetBuild) error {
	packageJSON := PackageJSON{
		SchemaVersion:    1,
		Name:             s.plugin.Name(),
		Message:          message,
		Team:             s.plugin.Team(),
		Kind:             s.plugin.Kind(),
		Version:          version,
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

func (s *PluginServe) writeSpecJSONSchema(dir string) error {
	if s.plugin.JSONSchema() == "" {
		return nil
	}

	return os.WriteFile(filepath.Join(dir, "spec_json_schema.json"), []byte(s.plugin.JSONSchema()), 0644)
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

func (*PluginServe) versionRegex() *regexp.Regexp {
	return regexp.MustCompile(`^(var)?\s?Version\s*=`)
}

func (s *PluginServe) validatePluginExports(pluginPath string) error {
	st, err := os.Stat(pluginPath)
	if err != nil {
		return err
	}
	if !st.IsDir() {
		return errors.New("plugin path must be a directory")
	}

	checkRelativeDirs := []string{"resources" + string(filepath.Separator) + "plugin", "plugin"}
	foundDirs := []string{}
	for _, dir := range checkRelativeDirs {
		p := filepath.Join(pluginPath, dir)
		s, err := os.Stat(p)
		if err == nil && s.IsDir() {
			foundDirs = append(foundDirs, dir)
		}
	}
	if len(foundDirs) == 0 {
		return fmt.Errorf("plugin directory must contain at least one of the following directories: %s", strings.Join(checkRelativeDirs, ", "))
	}

	findVersion := s.versionRegex()

	foundVersion := false
	for _, dir := range foundDirs {
		p := filepath.Join(pluginPath, dir)
		if err := filepath.WalkDir(p, func(path string, d os.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() || foundVersion {
				return nil
			}
			if !strings.HasSuffix(strings.ToLower(d.Name()), ".go") {
				return nil
			}

			f, err := os.Open(path)
			if err != nil {
				return err
			}
			defer f.Close()
			if ok, err := containsRegex(f, findVersion); err != nil {
				return err
			} else if ok {
				foundVersion = true
			}
			return nil
		}); err != nil {
			return err
		}
	}
	if !foundVersion {
		return errors.New("could not find `Version` global variable in package")
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

func containsRegex(r io.Reader, needle *regexp.Regexp) (bool, error) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		if needle.MatchString(scanner.Text()) {
			return true, nil
		}
	}
	return false, scanner.Err()
}

func (s *PluginServe) newCmdPluginPackage() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "package -m <message> <version> <plugin_directory>",
		Short: pluginPackageShort,
		Long:  pluginPackageLong,
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			pluginVersion := args[0]
			pluginDirectory := args[1]
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
				return errors.New("message is required")
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

			if s.plugin.Name() == "" {
				return errors.New("plugin name is required for packaging")
			}
			if s.plugin.Team() == "" {
				return errors.New("plugin team is required (hint: use the plugin.WithTeam() option)")
			}
			if s.plugin.Kind() == "" {
				return errors.New("plugin kind is required (hint: use the plugin.WithKind() option)")
			}

			if err := s.validatePluginExports(pluginDirectory); err != nil {
				return err
			}

			if s.plugin.Kind() == plugin.KindSource {
				if err := s.plugin.Init(cmd.Context(), nil, plugin.NewClientOptions{
					NoConnection: true,
				}); err != nil {
					return err
				}
				if err := s.writeTablesJSON(cmd.Context(), distPath); err != nil {
					return err
				}
			}

			targets := []TargetBuild{}
			for _, target := range s.plugin.Targets() {
				fmt.Println("Building for OS: " + target.OS + ", ARCH: " + target.Arch)
				targetBuild, err := s.build(pluginDirectory, target, distPath, pluginVersion)
				if err != nil {
					return fmt.Errorf("failed to build plugin for %s/%s: %w", target.OS, target.Arch, err)
				}
				targets = append(targets, *targetBuild)
			}
			if err := s.writePackageJSON(distPath, pluginVersion, message, targets); err != nil {
				return fmt.Errorf("failed to write manifest: %w", err)
			}
			if err := s.copyDocs(distPath, docsPath); err != nil {
				return fmt.Errorf("failed to copy docs: %w", err)
			}
			if err := s.writeSpecJSONSchema(distPath); err != nil {
				return fmt.Errorf("failed to write spec json schema: %w", err)
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
