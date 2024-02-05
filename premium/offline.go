package premium

import (
	"crypto/ed25519"
	_ "embed"
	"encoding/hex"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/cloudquery/plugin-sdk/v4/plugin"
	"github.com/rs/zerolog"
)

type License struct {
	LicensedTo string    `json:"licensed_to"`       // Customers name, e.g. "Acme Inc"
	Plugins    []string  `json:"plugins,omitempty"` // List of plugins, each in the format <org>/<kind>/<name>, e.g. "cloudquery/source/aws". Optional, if empty all plugins are allowed.
	IssuedAt   time.Time `json:"issued_at"`
	ValidFrom  time.Time `json:"valid_from"`
	ExpiresAt  time.Time `json:"expires_at"`
}

type LicenseWrapper struct {
	LicenseBytes []byte `json:"license"`
	Signature    string `json:"signature"` // crypto
}

var (
	ErrInvalidLicenseSignature = errors.New("invalid license signature")
	ErrLicenseNotValidYet      = errors.New("license not valid yet")
	ErrLicenseExpired          = errors.New("license expired")
	ErrLicenseNotApplicable    = errors.New("license not applicable to this plugin")
)

//go:embed offline.key
var publicKey string

var timeFunc = time.Now

func ValidateLicense(logger zerolog.Logger, meta plugin.Meta, licenseFileOrDirectory string) error {
	fi, err := os.Stat(licenseFileOrDirectory)
	if err != nil {
		return err
	}
	if !fi.IsDir() {
		return validateLicenseFile(logger, meta, licenseFileOrDirectory)
	}

	found := false
	var lastError error
	err = filepath.WalkDir(licenseFileOrDirectory, func(path string, d os.DirEntry, err error) error {
		if d.IsDir() {
			if path == licenseFileOrDirectory {
				return nil
			}
			return filepath.SkipDir
		}
		if err != nil {
			return err
		}

		if filepath.Ext(path) != ".cqlicense" {
			return nil
		}

		logger.Debug().Str("path", path).Msg("considering license file")
		lastError = validateLicenseFile(logger, meta, path)
		switch lastError {
		case nil:
			found = true
			return filepath.SkipAll
		case ErrLicenseNotApplicable:
			return nil
		default:
			return lastError
		}
	})
	if err != nil {
		return err
	}
	if found {
		return nil
	}
	if lastError != nil {
		return lastError
	}
	return errors.New("failed to validate license directory")
}

func validateLicenseFile(logger zerolog.Logger, meta plugin.Meta, licenseFile string) error {
	licenseContents, err := os.ReadFile(licenseFile)
	if err != nil {
		return err
	}

	l, err := UnpackLicense(licenseContents)
	if err != nil {
		return err
	}

	if len(l.Plugins) > 0 {
		ref := strings.Join([]string{meta.Team, string(meta.Kind), meta.Name}, "/")
		teamRef := meta.Team + "/*"
		if !slices.Contains(l.Plugins, ref) && !slices.Contains(l.Plugins, teamRef) {
			return ErrLicenseNotApplicable
		}
	}

	return l.IsValid(logger)
}

func UnpackLicense(lic []byte) (*License, error) {
	publicKeyBytes, err := hex.DecodeString(publicKey)
	if err != nil {
		return nil, err
	}

	var lw LicenseWrapper
	if err := json.Unmarshal(lic, &lw); err != nil {
		return nil, err
	}

	signatureBytes, err := hex.DecodeString(lw.Signature)
	if err != nil {
		return nil, err
	}

	if !ed25519.Verify(publicKeyBytes, lw.LicenseBytes, signatureBytes) {
		return nil, ErrInvalidLicenseSignature
	}

	var l License
	if err := json.Unmarshal(lw.LicenseBytes, &l); err != nil {
		return nil, err
	}

	return &l, nil
}

func (l *License) IsValid(logger zerolog.Logger) error {
	now := timeFunc().UTC()
	if now.Before(l.ValidFrom) {
		return ErrLicenseNotValidYet
	}
	if now.After(l.ExpiresAt) {
		return ErrLicenseExpired
	}

	msg := logger.Info()
	if now.Add(15 * 24 * time.Hour).After(l.ExpiresAt) {
		msg = logger.Warn()
	}

	msg.Time("expires_at", l.ExpiresAt).Msgf("Offline license for %s loaded.", l.LicensedTo)
	return nil
}
