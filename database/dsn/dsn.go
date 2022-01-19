package dsn

import (
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/xo/dburl"
)

func init() {
	dburl.Register(dburl.Scheme{
		Driver:    "timescale",
		Generator: dburl.GenPostgres,
		Transport: dburl.TransportTCP | dburl.TransportUnix,
		Opaque:    false,
		Aliases:   []string{"timescaledb", "tsdb", "ts"},
		Override:  "",
	})
}

// ParseConnectionString will try and parse any type of connection string and return a dburl
func ParseConnectionString(connString string) (*dburl.URL, error) {
	u, err := dburl.Parse(connString)
	if err == dburl.ErrInvalidDatabaseScheme {
		// connString may be a database URL or a DSN
		connString, err = convertDSNToURL(connString)
		if err != nil {
			return nil, fmt.Errorf("failed to parse dsn string: %w", RedactParseError(err))
		}
		u, err = dburl.Parse(connString)
	}

	return u, RedactParseError(err)
}

// SetDSNElement parses the given DSN and sets/adds the given map values as query parameters, returning a URI DSN
func SetDSNElement(dsn string, elems map[string]string) (string, error) {
	u, err := ParseConnectionString(dsn)
	if err != nil {
		return "", err
	}

	vals := u.Query()
	for k, v := range elems {
		vals.Set(k, v)
	}
	u.RawQuery = vals.Encode()
	return u.String(), nil
}

// RedactParseError looks for url.Error and redacts the URL (DSN)
// The original message wrapping the url.Error is lost.
func RedactParseError(err error) error {
	if err == nil {
		return nil
	}

	// pgconn does a good job of redacting the password, but it wraps the net/url's parse error, which includes the original URI
	// dburl returns url.Parse errors as-is

	var e *url.Error
	if errors.As(err, &e) {
		e.URL = "DSN redacted"
		return e
	}

	return err
}

var asciiSpace = [256]uint8{'\t': 1, '\n': 1, '\v': 1, '\f': 1, '\r': 1, ' ': 1}

// ParseDSNSettings taken from https://github.com/jackc/pgconn
//nolint
func parseDSNSettings(s string) (map[string]string, error) {
	settings := make(map[string]string)

	nameMap := map[string]string{
		"dbname": "database",
	}

	for len(s) > 0 {
		var key, val string
		eqIdx := strings.IndexRune(s, '=')
		if eqIdx < 0 {
			return nil, errors.New("invalid dsn")
		}

		key = strings.Trim(s[:eqIdx], " \t\n\r\v\f")
		s = strings.TrimLeft(s[eqIdx+1:], " \t\n\r\v\f")
		if len(s) == 0 {
		} else if s[0] != '\'' {
			end := 0
			for ; end < len(s); end++ {
				if asciiSpace[s[end]] == 1 {
					break
				}
				if s[end] == '\\' {
					end++
					if end == len(s) {
						return nil, errors.New("invalid backslash")
					}
				}
			}
			val = strings.Replace(strings.Replace(s[:end], "\\\\", "\\", -1), "\\'", "'", -1)
			if end == len(s) {
				s = ""
			} else {
				s = s[end+1:]
			}
		} else { // quoted string
			s = s[1:]
			end := 0
			for ; end < len(s); end++ {
				if s[end] == '\'' {
					break
				}
				if s[end] == '\\' {
					end++
				}
			}
			if end == len(s) {
				return nil, errors.New("unterminated quoted string in connection info string")
			}
			val = strings.Replace(strings.Replace(s[:end], "\\\\", "\\", -1), "\\'", "'", -1)
			if end == len(s) {
				s = ""
			} else {
				s = s[end+1:]
			}
		}

		if k, ok := nameMap[key]; ok {
			key = k
		}

		if key == "" {
			return nil, errors.New("invalid dsn")
		}

		settings[key] = val
	}
	return settings, nil
}

var nonQueryKeys = []string{"host", "port", "database", "password", "user"}

func convertDSNToURL(connString string) (string, error) {
	settings, err := parseDSNSettings(connString)
	if err != nil {
		return "", fmt.Errorf("failed to parse dsn string, %w", err)
	}
	host, ok := settings["host"]
	if !ok {
		host = "localhost"
	}
	delete(settings, "host")
	port, ok := settings["port"]
	if !ok {
		port = "5432"
	}
	delete(settings, "port")
	baseURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", settings["user"], settings["password"], host, port, settings["database"])

	for _, k := range nonQueryKeys {
		delete(settings, k)
	}
	if len(settings) == 0 {
		return baseURL, nil
	}
	queryParams := make([]string, 0)
	for k, v := range settings {
		queryParams = append(queryParams, fmt.Sprintf("%s=%s", k, v))

	}
	return fmt.Sprintf("%s?%s", baseURL, strings.Join(queryParams, "&")), nil
}
