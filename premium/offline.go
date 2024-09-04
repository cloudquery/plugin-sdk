package premium

import (
	"context"
	"crypto/ed25519"
	_ "embed"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/licensemanager"
	"github.com/aws/aws-sdk-go-v2/service/licensemanager/types"
	"github.com/cloudquery/plugin-sdk/v4/plugin"
	"github.com/google/uuid"
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

//go:generate mockgen -package=mocks -destination=../premium/mocks/licensemanager.go -source=offline.go AWSLicenseManagerInterface
type AWSLicenseManagerInterface interface {
	CheckoutLicense(ctx context.Context, params *licensemanager.CheckoutLicenseInput, optFns ...func(*licensemanager.Options)) (*licensemanager.CheckoutLicenseOutput, error)
}

type CQLicenseClient struct {
	logger                  zerolog.Logger
	meta                    plugin.Meta
	licenseFileOrDirectory  string
	awsLicenseManagerClient AWSLicenseManagerInterface
	isMarketplaceLicense    bool
}

type LicenseClientOptions func(updater *CQLicenseClient)

func WithMeta(meta plugin.Meta) LicenseClientOptions {
	return func(cl *CQLicenseClient) {
		cl.meta = meta
	}
}

func WithLicenseFileOrDirectory(licenseFileOrDirectory string) LicenseClientOptions {
	return func(cl *CQLicenseClient) {
		cl.licenseFileOrDirectory = licenseFileOrDirectory
	}
}

func WithAWSLicenseManagerClient(awsLicenseManagerClient AWSLicenseManagerInterface) LicenseClientOptions {
	return func(cl *CQLicenseClient) {
		cl.awsLicenseManagerClient = awsLicenseManagerClient
	}
}

func NewLicenseClient(ctx context.Context, logger zerolog.Logger, ops ...LicenseClientOptions) (CQLicenseClient, error) {
	cl := CQLicenseClient{
		logger:               logger,
		isMarketplaceLicense: os.Getenv("CQ_AWS_MARKETPLACE_LICENSE") == "true",
	}

	for _, op := range ops {
		op(&cl)
	}

	if cl.isMarketplaceLicense && cl.awsLicenseManagerClient == nil {
		cfg, err := awsConfig.LoadDefaultConfig(ctx)
		if err != nil {
			return cl, fmt.Errorf("failed to load AWS config: %w", err)
		}
		cl.awsLicenseManagerClient = licensemanager.NewFromConfig(cfg)
	}

	return cl, nil
}

func (lc CQLicenseClient) ValidateLicense(ctx context.Context) error {
	// License can be provided via environment variable for AWS Marketplace or CLI flag
	switch {
	case lc.isMarketplaceLicense:
		return lc.validateMarketplaceLicense(ctx)
	case lc.licenseFileOrDirectory != "":
		lc.validateCQLicense()
	default:
		return ErrLicenseNotApplicable
	}

}

func (lc CQLicenseClient) validateCQLicense() error {
	fi, err := os.Stat(lc.licenseFileOrDirectory)
	if err != nil {
		return err
	}
	if !fi.IsDir() {
		return lc.validateLicenseFile(lc.licenseFileOrDirectory)
	}

	found := false
	var lastError error
	err = filepath.WalkDir(lc.licenseFileOrDirectory, func(path string, d os.DirEntry, err error) error {
		if d.IsDir() {
			if path == lc.licenseFileOrDirectory {
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

		lc.logger.Debug().Str("path", path).Msg("considering license file")
		lastError = lc.validateLicenseFile(path)
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

func (lc CQLicenseClient) validateLicenseFile(licenseFile string) error {
	licenseContents, err := os.ReadFile(licenseFile)
	if err != nil {
		return err
	}

	l, err := UnpackLicense(licenseContents)
	if err != nil {
		return err
	}

	if len(l.Plugins) > 0 {
		ref := strings.Join([]string{lc.meta.Team, string(lc.meta.Kind), lc.meta.Name}, "/")
		teamRef := lc.meta.Team + "/*"
		if !slices.Contains(l.Plugins, ref) && !slices.Contains(l.Plugins, teamRef) {
			return ErrLicenseNotApplicable
		}
	}

	return l.IsValid(lc.logger)
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

func (lc CQLicenseClient) validateMarketplaceLicense(ctx context.Context) error {
	clientToken := uuid.New()

	resp, err := lc.awsLicenseManagerClient.CheckoutLicense(ctx, &licensemanager.CheckoutLicenseInput{
		CheckoutType: types.CheckoutTypeProvisional,
		ClientToken:  aws.String(clientToken.String()),
		ProductSKU:   aws.String("55ukc0d5qv3gebks148tjr62j"),
		Entitlements: []types.EntitlementData{
			{
				Name: aws.String("Unlimited"),
				Unit: types.EntitlementDataUnitNone,
			},
		},
		// This is hardcoded for AWS Marketplace, because this is the only supported value for marketplace licenses
		KeyFingerprint: aws.String("aws:294406891311:AWS/Marketplace:issuer-fingerprint"),
	})
	if err != nil {
		return fmt.Errorf("failed to checkout license: %w", err)
	}
	if len(resp.EntitlementsAllowed) == 0 {
		return errors.New("no entitlements provisioned")
	}
	return nil
}
