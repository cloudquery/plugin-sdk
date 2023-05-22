package source

import (
	"bytes"
	"embed"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"text/template"

	"github.com/apache/arrow/go/v13/arrow"
	"github.com/cloudquery/plugin-sdk/v3/caser"
	"github.com/cloudquery/plugin-sdk/v3/schema"
	"github.com/cloudquery/plugin-sdk/v3/types"
)

//go:embed templates/*.go.tpl
var templatesFS embed.FS

var reMatchNewlines = regexp.MustCompile(`\n{3,}`)
var reMatchHeaders = regexp.MustCompile(`(#{1,6}.+)\n+`)

var DefaultTitleExceptions = map[string]string{
	// common abbreviations
	"acl":   "ACL",
	"acls":  "ACLs",
	"api":   "API",
	"apis":  "APIs",
	"ca":    "CA",
	"cidr":  "CIDR",
	"cidrs": "CIDRs",
	"db":    "DB",
	"dbs":   "DBs",
	"dhcp":  "DHCP",
	"iam":   "IAM",
	"iot":   "IOT",
	"ip":    "IP",
	"ips":   "IPs",
	"ipv4":  "IPv4",
	"ipv6":  "IPv6",
	"mfa":   "MFA",
	"ml":    "ML",
	"oauth": "OAuth",
	"vpc":   "VPC",
	"vpcs":  "VPCs",
	"vpn":   "VPN",
	"vpns":  "VPNs",
	"waf":   "WAF",
	"wafs":  "WAFs",

	// cloud providers
	"aws": "AWS",
	"gcp": "GCP",
}

func DefaultTitleTransformer(table *schema.Table) string {
	if table.Title != "" {
		return table.Title
	}
	csr := caser.New(caser.WithCustomExceptions(DefaultTitleExceptions))
	return csr.ToTitle(table.Name)
}

func sortTables(tables schema.Tables) {
	sort.SliceStable(tables, func(i, j int) bool {
		return tables[i].Name < tables[j].Name
	})

	for _, table := range tables {
		sortTables(table.Relations)
	}
}

type templateData struct {
	PluginName string
	Tables     schema.Tables
}

// GeneratePluginDocs creates table documentation for the source plugin based on its list of tables
func (p *Plugin) GeneratePluginDocs(dir, format string) error {
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}

	setDestinationManagedCqColumns(p.Tables())

	sortedTables := make(schema.Tables, 0, len(p.Tables()))
	for _, t := range p.Tables() {
		sortedTables = append(sortedTables, t.Copy(nil))
	}
	sortTables(sortedTables)

	switch format {
	case "markdown":
		return p.renderTablesAsMarkdown(dir, p.name, sortedTables)
	case "json":
		return p.renderTablesAsJSON(dir, sortedTables)
	default:
		return fmt.Errorf("unsupported format: %v", format)
	}
}

// setDestinationManagedCqColumns overwrites or adds the CQ columns that are managed by the destination plugins (_cq_sync_time, _cq_source_name).
func setDestinationManagedCqColumns(tables []*schema.Table) {
	for _, table := range tables {
		table.OverwriteOrAddColumn(&schema.CqSyncTimeColumn)
		table.OverwriteOrAddColumn(&schema.CqSourceNameColumn)
		setDestinationManagedCqColumns(table.Relations)
	}
}

type jsonTable struct {
	Name        string       `json:"name"`
	Title       string       `json:"title"`
	Description string       `json:"description"`
	Columns     []jsonColumn `json:"columns"`
	Relations   []jsonTable  `json:"relations"`
}

type jsonColumn struct {
	Name             string `json:"name"`
	Type             string `json:"type"`
	IsPrimaryKey     bool   `json:"is_primary_key,omitempty"`
	IsIncrementalKey bool   `json:"is_incremental_key,omitempty"`
}

func (p *Plugin) renderTablesAsJSON(dir string, tables schema.Tables) error {
	jsonTables := p.jsonifyTables(tables)
	buffer := &bytes.Buffer{}
	m := json.NewEncoder(buffer)
	m.SetIndent("", "  ")
	m.SetEscapeHTML(false)
	err := m.Encode(jsonTables)
	if err != nil {
		return err
	}
	outputPath := filepath.Join(dir, "__tables.json")
	return os.WriteFile(outputPath, buffer.Bytes(), 0644)
}

func (p *Plugin) jsonifyTables(tables schema.Tables) []jsonTable {
	jsonTables := make([]jsonTable, len(tables))
	for i, table := range tables {
		jsonColumns := make([]jsonColumn, len(table.Columns))
		for c, col := range table.Columns {
			jsonColumns[c] = jsonColumn{
				Name:             col.Name,
				Type:             formatType(col.Type),
				IsPrimaryKey:     col.PrimaryKey,
				IsIncrementalKey: col.IncrementalKey,
			}
		}
		jsonTables[i] = jsonTable{
			Name:        table.Name,
			Title:       p.titleTransformer(table),
			Description: table.Description,
			Columns:     jsonColumns,
			Relations:   p.jsonifyTables(table.Relations),
		}
	}
	return jsonTables
}

func (p *Plugin) renderTablesAsMarkdown(dir string, pluginName string, tables schema.Tables) error {
	for _, table := range tables {
		if err := p.renderAllTables(table, dir); err != nil {
			return err
		}
	}
	t, err := template.New("all_tables.md.go.tpl").Funcs(template.FuncMap{
		"indentToDepth": indentToDepth,
	}).ParseFS(templatesFS, "templates/all_tables*.md.go.tpl")
	if err != nil {
		return fmt.Errorf("failed to parse template for README.md: %v", err)
	}

	var b bytes.Buffer
	if err := t.Execute(&b, templateData{PluginName: pluginName, Tables: tables}); err != nil {
		return fmt.Errorf("failed to execute template: %v", err)
	}
	content := formatMarkdown(b.String())
	outputPath := filepath.Join(dir, "README.md")
	f, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create file %v: %v", outputPath, err)
	}
	f.WriteString(content)
	return nil
}

func (p *Plugin) renderAllTables(t *schema.Table, dir string) error {
	if err := p.renderTable(t, dir); err != nil {
		return err
	}
	for _, r := range t.Relations {
		if err := p.renderAllTables(r, dir); err != nil {
			return err
		}
	}
	return nil
}

func (p *Plugin) renderTable(table *schema.Table, dir string) error {
	t := template.New("").Funcs(map[string]any{
		"formatType": formatType,
		"title":      p.titleTransformer,
	})
	t, err := t.New("table.md.go.tpl").ParseFS(templatesFS, "templates/table.md.go.tpl")
	if err != nil {
		return fmt.Errorf("failed to parse template: %v", err)
	}

	outputPath := filepath.Join(dir, fmt.Sprintf("%s.md", table.Name))

	var b bytes.Buffer
	if err := t.Execute(&b, table); err != nil {
		return fmt.Errorf("failed to execute template: %v", err)
	}
	content := formatMarkdown(b.String())
	f, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create file %v: %v", outputPath, err)
	}
	f.WriteString(content)
	return f.Close()
}

func formatMarkdown(s string) string {
	s = reMatchNewlines.ReplaceAllString(s, "\n\n")
	return reMatchHeaders.ReplaceAllString(s, `$1`+"\n\n")
}

func formatType(v arrow.DataType) string {
	if arrow.IsNested(v.ID()) {
		switch v.ID() {
		case arrow.STRUCT:
			s := "struct<"
			for i, f := range v.(*arrow.StructType).Fields() {
				if i > 0 {
					s += ","
				}
				s += f.Name + ":" + formatType(f.Type)
			}
			return s + ">"
		case arrow.LIST:
			return "list<" + formatType(v.(*arrow.ListType).Elem()) + ">"
		case arrow.LARGE_LIST:
			return "large_list<" + formatType(v.(*arrow.LargeListType).Elem()) + ">"
		case arrow.FIXED_SIZE_LIST:
			return fmt.Sprintf("fixed_size_list<%s,%d>", formatType(v.(*arrow.FixedSizeListType).Elem()), v.(*arrow.FixedSizeListType).Len())
		case arrow.MAP:
			return fmt.Sprintf("map<%s,%s>", formatType(v.(*arrow.MapType).KeyType()), formatType(v.(*arrow.MapType).ItemType()))
		default:
			// TODO: support other nested types like unions
			panic("unsupported nested type: " + v.String())
		}
	}
	switch {
	case arrow.TypeEqual(v, arrow.FixedWidthTypes.Boolean):
		return "boolean"
	case arrow.TypeEqual(v, arrow.FixedWidthTypes.Duration_s):
		return "duration[seconds]"
	case arrow.TypeEqual(v, arrow.FixedWidthTypes.Duration_ms):
		return "duration[milliseconds]"
	case arrow.TypeEqual(v, arrow.FixedWidthTypes.Duration_us):
		return "duration[microseconds]"
	case arrow.TypeEqual(v, arrow.FixedWidthTypes.Duration_ns):
		return "duration[nanoseconds]"
	case arrow.TypeEqual(v, arrow.FixedWidthTypes.MonthDayNanoInterval):
		return "month_day_nano_interval"
	case arrow.TypeEqual(v, arrow.FixedWidthTypes.DayTimeInterval):
		return "day_time_interval"
	case arrow.TypeEqual(v, arrow.FixedWidthTypes.Time32s):
		return "time32[seconds]"
	case arrow.TypeEqual(v, arrow.FixedWidthTypes.Time32ms):
		return "time32[milliseconds]"
	case arrow.TypeEqual(v, arrow.FixedWidthTypes.Time64us):
		return "time64[microseconds]"
	case arrow.TypeEqual(v, arrow.FixedWidthTypes.Time64ns):
		return "time64[nanoseconds]"
	case arrow.TypeEqual(v, arrow.FixedWidthTypes.Date32):
		return "date32"
	case arrow.TypeEqual(v, arrow.FixedWidthTypes.Date64):
		return "date64"
	case arrow.TypeEqual(v, arrow.FixedWidthTypes.Timestamp_s):
		return "timestamp[seconds]"
	case arrow.TypeEqual(v, arrow.FixedWidthTypes.Timestamp_ms):
		return "timestamp[milliseconds]"
	case arrow.TypeEqual(v, arrow.FixedWidthTypes.Timestamp_us):
		return "timestamp[microseconds]"
	case arrow.TypeEqual(v, arrow.FixedWidthTypes.Timestamp_ns):
		return "timestamp[nanoseconds]"
	case arrow.TypeEqual(v, arrow.PrimitiveTypes.Uint8):
		return "uint8"
	case arrow.TypeEqual(v, arrow.PrimitiveTypes.Uint16):
		return "uint16"
	case arrow.TypeEqual(v, arrow.PrimitiveTypes.Uint32):
		return "uint32"
	case arrow.TypeEqual(v, arrow.PrimitiveTypes.Uint64):
		return "uint64"
	case arrow.TypeEqual(v, arrow.PrimitiveTypes.Int8):
		return "int8"
	case arrow.TypeEqual(v, arrow.PrimitiveTypes.Int16):
		return "int16"
	case arrow.TypeEqual(v, arrow.PrimitiveTypes.Int32):
		return "int32"
	case arrow.TypeEqual(v, arrow.PrimitiveTypes.Int64):
		return "int64"
	case arrow.TypeEqual(v, arrow.PrimitiveTypes.Float32):
		return "float32"
	case arrow.TypeEqual(v, arrow.PrimitiveTypes.Float64):
		return "float64"
	case arrow.TypeEqual(v, arrow.BinaryTypes.String):
		return "string"
	case arrow.TypeEqual(v, arrow.BinaryTypes.LargeString):
		return "large_string"
	case arrow.TypeEqual(v, arrow.BinaryTypes.Binary):
		return "binary"
	case arrow.TypeEqual(v, arrow.BinaryTypes.LargeBinary):
		return "large_binary"
	case arrow.TypeEqual(v, types.ExtensionTypes.UUID):
		return "uuid"
	case arrow.TypeEqual(v, types.ExtensionTypes.JSON):
		return "json"
	case arrow.TypeEqual(v, types.ExtensionTypes.Inet):
		return "inet_address"
	case arrow.TypeEqual(v, types.ExtensionTypes.MAC):
		return "mac_address"
	}
	return v.String()
}

func indentToDepth(table *schema.Table) string {
	s := ""
	t := table
	for t.Parent != nil {
		s += "  "
		t = t.Parent
	}
	return s
}
