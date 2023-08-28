package docs

import (
	"github.com/cloudquery/plugin-sdk/v4/caser"
	"github.com/cloudquery/plugin-sdk/v4/schema"
)

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
