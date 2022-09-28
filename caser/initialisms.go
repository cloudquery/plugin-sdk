package caser

// commonInitialisms taken from https://github.com/golang/lint/blob/master/lint.go
var commonInitialisms = map[string]bool{
	"ACL":   true,
	"API":   true,
	"ASCII": true,
	"CIDR":  true,
	"CPU":   true,
	"CSS":   true,
	"DNS":   true,
	"EOF":   true,
	"FQDN":  true,
	"GUID":  true,
	"HTML":  true,
	"HTTP":  true,
	"HTTPS": true,
	"ID":    true,
	"IP":    true,
	"IPC":   true,
	"IPv4":  true,
	"IPv6":  true,
	"JSON":  true,
	"LHS":   true,
	"PID":   true,
	"QOS":   true,
	"QPS":   true,
	"RAM":   true,
	"RHS":   true,
	"RPC":   true,
	"SLA":   true,
	"SMTP":  true,
	"SQL":   true,
	"SSH":   true,
	"TCP":   true,
	"TLS":   true,
	"TTL":   true,
	"UDP":   true,
	"UI":    true,
	"UID":   true,
	"UUID":  true,
	"URI":   true,
	"URL":   true,
	"UTF8":  true,
	"VM":    true,
	"XML":   true,
	"XMPP":  true,
	"XSRF":  true,
	"XSS":   true,
}

// add exceptions here for things that are not automatically convertable
var snakeToCamelExceptions = map[string]string{
	"oauth": "OAuth",
	"ipv4":  "IPv4",
	"ipv6":  "IPv6",
}

// add exceptions here for things that are not automatically convertable
var camelToSnakeExceptions = map[string]string{
	"IPv4": "ipv4",
	"IPv6": "ipv6",
}

// startsWithInitialism returns the initialism if the given string begins with it
func startsWithInitialism(s string) string {
	var initialism string
	// the longest initialism is 5 char, the shortest 2
	// we choose the longest match
	for i := 1; i <= len(s) && i <= 5; i++ {
		if len(s) > i-1 && commonInitialisms[s[:i]] && len(s[:i]) > len(initialism) {
			initialism = s[:i]
		}
	}
	return initialism
}

func ConfigureInitialisms(c map[string]bool) {
	for k, v := range c {
		commonInitialisms[k] = v
	}
}
