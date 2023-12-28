package premium

import (
	"crypto/ed25519"
	"encoding/hex"
	"encoding/json"
	"errors"
	"os"
	"time"

	"github.com/rs/zerolog"
)

type License struct {
	LicensedTo string    `json:"licensed_to"` // Customers name, e.g. "Acme Inc"
	Plugin     string    `json:"plugin"`      // This is <team>/<kind>/<name>
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
	ErrInvalidLicenseDetails   = errors.New("invalid license details")
	ErrLicenseNotValidYet      = errors.New("license not valid yet")
	ErrLicenseExpired          = errors.New("license expired")
)

func ValidateLicense(logger zerolog.Logger, licenseFile, pluginTeam, pluginKind, pluginName string) error {
	licenseContents, err := os.ReadFile(licenseFile)
	if err != nil {
		return err
	}

	l, err := UnpackLicense(licenseContents)
	if err != nil {
		return err
	}

	return l.IsValid(logger, pluginTeam, pluginKind, pluginName)
}

func UnpackLicense(lic []byte) (*License, error) {
	const publicKey = "96e4749b7550d33bd776cb7fb056d74cb16ed69dfe8e59c16e8c2500a94162b1"

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
func (l *License) IsValid(logger zerolog.Logger, pluginTeam, pluginKind, pluginName string) error {
	expectPlugin := pluginTeam + "/" + pluginKind + "/" + pluginName
	if l.Plugin != expectPlugin {
		return ErrInvalidLicenseDetails
	}
	now := time.Now().UTC()
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
