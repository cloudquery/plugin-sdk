package docs

import (
	"embed"
	"errors"
	"fmt"
	"os"
	"regexp"
	"sort"

	"github.com/cloudquery/plugin-sdk/v4/caser"
	"github.com/cloudquery/plugin-sdk/v4/schema"
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

type Format int

const (
	FormatMarkdown Format = iota
	FormatJSON
)

func (r Format) String() string {
	return [...]string{"markdown", "json"}[r]
}

func FormatFromString(s string) (Format, error) {
	switch s {
	case "markdown":
		return FormatMarkdown, nil
	case "json":
		return FormatJSON, nil
	default:
		return FormatMarkdown, fmt.Errorf("unknown format %s", s)
	}
}

type Generator struct {
	tables           schema.Tables
	titleTransformer func(*schema.Table) string
	pluginName       string
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

// NewGenerator creates a new generator for the given tables.
// The tables are sorted by name. pluginName is optional and is used in markdown only
func NewGenerator(pluginName string, tables schema.Tables) *Generator {
	sortedTables := make(schema.Tables, 0, len(tables))
	for _, t := range tables {
		sortedTables = append(sortedTables, t.Copy(nil))
	}
	sortTables(sortedTables)

	return &Generator{
		tables:           sortedTables,
		titleTransformer: DefaultTitleTransformer,
		pluginName:       pluginName,
	}
}

func (g *Generator) Generate(dir string, format Format) error {
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}

	switch format {
	case FormatMarkdown:
		return errors.New("markdown format is not supported directly via the plugin, use the `cloudquery tables` command instead")
	case FormatJSON:
		return g.renderTablesAsJSON(dir)
	default:
		return fmt.Errorf("unsupported format: %v", format)
	}
}

// setDestinationManagedCqColumns overwrites or adds the CQ columns that are managed by the destination plugins (_cq_sync_time, _cq_source_name).
// func setDestinationManagedCqColumns(tables []*schema.Table) {
// 	for _, table := range tables {
// 		table.OverwriteOrAddColumn(&schema.CqSyncTimeColumn)
// 		table.OverwriteOrAddColumn(&schema.CqSourceNameColumn)
// 		setDestinationManagedCqColumns(table.Relations)
// 	}
// }
