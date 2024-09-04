package premium

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/licensemanager"
	"github.com/aws/aws-sdk-go-v2/service/licensemanager/types"
	"github.com/cloudquery/plugin-sdk/v4/faker"
	"github.com/cloudquery/plugin-sdk/v4/plugin"
	"github.com/cloudquery/plugin-sdk/v4/premium/mocks"
	"github.com/golang/mock/gomock"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnpackLicense(t *testing.T) {
	publicKey = "eacdff4866c8bc0d97de3c2d7668d0970c61aa16c3f12d6ba8083147ff92c9a6"

	t.Run("Success", func(t *testing.T) {
		licData := `{"license":"eyJsaWNlbnNlZF90byI6IlVOTElDRU5TRUQgVEVTVCIsImlzc3VlZF9hdCI6IjIwMjMtMTItMjhUMTk6MDI6MjguODM4MzY3WiIsInZhbGlkX2Zyb20iOiIyMDIzLTEyLTI4VDE5OjAyOjI4LjgzODM2N1oiLCJleHBpcmVzX2F0IjoiMjAyMy0xMi0yOVQxOTowMjoyOC44MzgzNjdaIn0=","signature":"8687a858463764b052455b3c783d979d364b5fb653b86d88a7463e495480db62fdec7ae1a84d1e30dddee77eb769a0e498ecfc836538c53e410aeb1a0c04d102"}`

		l, err := UnpackLicense([]byte(licData))
		require.NoError(t, err)
		require.Equal(t, "UNLICENSED TEST", l.LicensedTo)
		require.Equal(t, l.ExpiresAt.Add(-24*time.Hour).Truncate(time.Hour), l.ValidFrom.Truncate(time.Hour))
	})
	t.Run("Fail", func(t *testing.T) {
		licData := `{"license":"eyJsaWNlbnNlZF90byI6IlVOTElDRU5TRUQgVEVTVCIsImlzc3VlZF9hdCI6IjIwMjMtMTItMjhUMTk6MDI6MjguODM4MzY3WiIsInZhbGlkX2Zyb20iOiIyMDIzLTEyLTI4VDE5OjAyOjI4LjgzODM2N1oiLCJleHBpcmVzX2F0IjoiMjAyMy0xMi0yOVQxOTowMjoyOC44MzgzNjdaIn0=","signature":"9687a858463764b052455b3c783d979d364b5fb653b86d88a7463e495480db62fdec7ae1a84d1e30dddee77eb769a0e498ecfc836538c53e410aeb1a0c04d102"}`
		l, err := UnpackLicense([]byte(licData))
		require.ErrorIs(t, err, ErrInvalidLicenseSignature)
		require.Nil(t, l)
	})
}

func TestValidateLicense(t *testing.T) {
	publicKey = "eacdff4866c8bc0d97de3c2d7668d0970c61aa16c3f12d6ba8083147ff92c9a6"
	licData := `{"license":"eyJsaWNlbnNlZF90byI6IlVOTElDRU5TRUQgVEVTVCIsImlzc3VlZF9hdCI6IjIwMjMtMTItMjhUMTk6MDI6MjguODM4MzY3WiIsInZhbGlkX2Zyb20iOiIyMDIzLTEyLTI4VDE5OjAyOjI4LjgzODM2N1oiLCJleHBpcmVzX2F0IjoiMjAyMy0xMi0yOVQxOTowMjoyOC44MzgzNjdaIn0=","signature":"8687a858463764b052455b3c783d979d364b5fb653b86d88a7463e495480db62fdec7ae1a84d1e30dddee77eb769a0e498ecfc836538c53e410aeb1a0c04d102"}`
	validTime := time.Date(2023, 12, 29, 12, 0, 0, 0, time.UTC)
	expiredTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	nopMeta := plugin.Meta{Team: "cloudquery", Kind: "source", Name: "test"}

	t.Run("SingleFile", func(t *testing.T) {
		dir := t.TempDir()
		f := filepath.Join(dir, "testlicense.cqlicense")
		if err := os.WriteFile(f, []byte(licData), 0644); err != nil {
			require.NoError(t, err)
		}

		t.Run("Expired", licenseTest(f, nopMeta, expiredTime, ErrLicenseExpired))
		t.Run("Success", licenseTest(f, nopMeta, validTime, nil))
	})
	t.Run("Dir", func(t *testing.T) {
		dir := t.TempDir()
		f := filepath.Join(dir, "testlicense.cqlicense")
		if err := os.WriteFile(f, []byte(licData), 0644); err != nil {
			require.NoError(t, err)
		}
		t.Run("Expired", licenseTest(dir, nopMeta, expiredTime, ErrLicenseExpired))
		t.Run("Success", licenseTest(dir, nopMeta, validTime, nil))
	})
}

func TestValidateSpecificLicense(t *testing.T) {
	publicKey = `de452e6028fe488f56ee0dfcf5b387ee773f03d24de66f00c40ec5b17085c549`
	licData := `{"license":"eyJsaWNlbnNlZF90byI6IlVOTElDRU5TRUQgVEVTVCIsInBsdWdpbnMiOlsiY2xvdWRxdWVyeS9zb3VyY2UvdGVzdDEiLCJjbG91ZHF1ZXJ5L3NvdXJjZS90ZXN0MiJdLCJpc3N1ZWRfYXQiOiIyMDI0LTAxLTAyVDExOjEwOjA5LjE0OTYwNVoiLCJ2YWxpZF9mcm9tIjoiMjAyNC0wMS0wMlQxMToxMDowOS4xNDk2MDVaIiwiZXhwaXJlc19hdCI6IjIwMjQtMDEtMDNUMTE6MTA6MDkuMTQ5NjA1WiJ9","signature":"e5752577c2b2c5a8920b3277fd11504d9c6820e8acb22bc17ccda524857c1d9fc7534f39b9a122376069ad682a2b616a10d1cfae40a984fb57fee31f13a15302"}`
	validTime := time.Date(2024, 1, 2, 12, 0, 0, 0, time.UTC)
	expiredTime := time.Date(2024, 1, 3, 12, 0, 0, 0, time.UTC)
	invalidMeta := plugin.Meta{Team: "cloudquery", Kind: "source", Name: "test"}
	validMeta := plugin.Meta{Team: "cloudquery", Kind: "source", Name: "test1"}

	t.Run("SingleFile", func(t *testing.T) {
		dir := t.TempDir()
		f := filepath.Join(dir, "testlicense.cqlicense")
		if err := os.WriteFile(f, []byte(licData), 0644); err != nil {
			require.NoError(t, err)
		}

		t.Run("Expired", licenseTest(f, validMeta, expiredTime, ErrLicenseExpired))
		t.Run("Success", licenseTest(f, validMeta, validTime, nil))
		t.Run("NotApplicable", licenseTest(f, invalidMeta, validTime, ErrLicenseNotApplicable))
	})
	t.Run("SingleDir", func(t *testing.T) {
		dir := t.TempDir()
		if err := os.WriteFile(filepath.Join(dir, "testlicense.cqlicense"), []byte(licData), 0644); err != nil {
			require.NoError(t, err)
		}
		t.Run("Expired", licenseTest(dir, validMeta, expiredTime, ErrLicenseExpired))
		t.Run("Success", licenseTest(dir, validMeta, validTime, nil))
		t.Run("NotApplicable", licenseTest(dir, invalidMeta, validTime, ErrLicenseNotApplicable))
	})
}

func TestValidateSpecificLicenseMultiFile(t *testing.T) {
	publicKey = `de452e6028fe488f56ee0dfcf5b387ee773f03d24de66f00c40ec5b17085c549`
	licData1 := `{"license":"eyJsaWNlbnNlZF90byI6IlVOTElDRU5TRUQgVEVTVCIsInBsdWdpbnMiOlsiY2xvdWRxdWVyeS9zb3VyY2UvdGVzdDEiLCJjbG91ZHF1ZXJ5L3NvdXJjZS90ZXN0MiJdLCJpc3N1ZWRfYXQiOiIyMDI0LTAxLTAyVDExOjEwOjA5LjE0OTYwNVoiLCJ2YWxpZF9mcm9tIjoiMjAyNC0wMS0wMlQxMToxMDowOS4xNDk2MDVaIiwiZXhwaXJlc19hdCI6IjIwMjQtMDEtMDNUMTE6MTA6MDkuMTQ5NjA1WiJ9","signature":"e5752577c2b2c5a8920b3277fd11504d9c6820e8acb22bc17ccda524857c1d9fc7534f39b9a122376069ad682a2b616a10d1cfae40a984fb57fee31f13a15302"}`
	licData3 := `{"license":"eyJsaWNlbnNlZF90byI6IlVOTElDRU5TRUQgVEVTVDMiLCJwbHVnaW5zIjpbImNsb3VkcXVlcnkvc291cmNlL3Rlc3QzIl0sImlzc3VlZF9hdCI6IjIwMjQtMDEtMDJUMTE6MjA6NTcuMzE2NDE0WiIsInZhbGlkX2Zyb20iOiIyMDI0LTAxLTAyVDExOjIwOjU3LjMxNjQxNFoiLCJleHBpcmVzX2F0IjoiMjAyNC0wMS0wM1QxMToyMDo1Ny4zMTY0MTRaIn0=","signature":"9be752d46010af84ec7295ede29915950dab13d4eca3b82b5645f793b39a03a6eef6bc653bee26e2a4f148b4d0fd54df6401059fda6104bc207f6dec2127850f"}`

	validTime := time.Date(2024, 1, 2, 12, 0, 0, 0, time.UTC)
	expiredTime := time.Date(2024, 1, 3, 12, 0, 0, 0, time.UTC)
	invalidMeta := plugin.Meta{Team: "cloudquery", Kind: "source", Name: "test"}
	validMeta1 := plugin.Meta{Team: "cloudquery", Kind: "source", Name: "test1"}
	validMeta3 := plugin.Meta{Team: "cloudquery", Kind: "source", Name: "test3"}

	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "testlicense1.cqlicense"), []byte(licData1), 0644); err != nil {
		require.NoError(t, err)
	}
	if err := os.WriteFile(filepath.Join(dir, "testlicense3.cqlicense"), []byte(licData3), 0644); err != nil {
		require.NoError(t, err)
	}

	t.Run("Expired", licenseTest(dir, validMeta1, expiredTime, ErrLicenseExpired))
	t.Run("Success", licenseTest(dir, validMeta1, validTime, nil))
	t.Run("SuccessOther", licenseTest(dir, validMeta3, validTime, nil))
	t.Run("NotApplicable", licenseTest(dir, invalidMeta, validTime, ErrLicenseNotApplicable))
}

func TestValidateTeamLicense(t *testing.T) {
	publicKey = `de452e6028fe488f56ee0dfcf5b387ee773f03d24de66f00c40ec5b17085c549`
	licData := `{"license":"eyJsaWNlbnNlZF90byI6IlVOTElDRU5TRUQgVEVTVCIsInBsdWdpbnMiOlsidGVzdC10ZWFtLyoiXSwiaXNzdWVkX2F0IjoiMjAyNC0wMi0wNVQxNjozOTozMy4zMzkxMjZaIiwidmFsaWRfZnJvbSI6IjIwMjQtMDItMDVUMTY6Mzk6MzMuMzM5MTI2WiIsImV4cGlyZXNfYXQiOiIyMDI0LTAyLTA2VDE2OjM5OjMzLjMzOTEyNloifQ==","signature":"cba85dcbd48d909f92d6e84d1d56b47075484efb2a7db1c478fc09659bb498e2a761add3c743c2d9a50b82b29b1730600cd8f68d6571896ca7d08f3107751e07"}`
	validTime := time.Date(2024, 2, 5, 18, 0, 0, 0, time.UTC)
	expiredTime := time.Date(2024, 2, 6, 18, 0, 0, 0, time.UTC)
	invalidMeta := plugin.Meta{Team: "cloudquery", Kind: "source", Name: "test"}
	validMeta1 := plugin.Meta{Team: "test-team", Kind: "source", Name: "some-plugin"}
	validMeta2 := plugin.Meta{Team: "test-team", Kind: "destination", Name: "some-plugin2"}

	t.Run("SingleFile", func(t *testing.T) {
		dir := t.TempDir()
		f := filepath.Join(dir, "testlicense.cqlicense")
		if err := os.WriteFile(f, []byte(licData), 0644); err != nil {
			require.NoError(t, err)
		}

		t.Run("Expired1", licenseTest(f, validMeta1, expiredTime, ErrLicenseExpired))
		t.Run("Expired2", licenseTest(f, validMeta2, expiredTime, ErrLicenseExpired))
		t.Run("Success1", licenseTest(f, validMeta1, validTime, nil))
		t.Run("Success2", licenseTest(f, validMeta2, validTime, nil))
		t.Run("NotApplicable", licenseTest(f, invalidMeta, validTime, ErrLicenseNotApplicable))
	})
	t.Run("SingleDir", func(t *testing.T) {
		dir := t.TempDir()
		if err := os.WriteFile(filepath.Join(dir, "testlicense.cqlicense"), []byte(licData), 0644); err != nil {
			require.NoError(t, err)
		}
		t.Run("Expired1", licenseTest(dir, validMeta1, expiredTime, ErrLicenseExpired))
		t.Run("Expired2", licenseTest(dir, validMeta2, expiredTime, ErrLicenseExpired))
		t.Run("Success1", licenseTest(dir, validMeta1, validTime, nil))
		t.Run("Success2", licenseTest(dir, validMeta2, validTime, nil))
		t.Run("NotApplicable", licenseTest(dir, invalidMeta, validTime, ErrLicenseNotApplicable))
	})
}

func licenseTest(inputPath string, meta plugin.Meta, timeIs time.Time, expectError error) func(t *testing.T) {
	return func(t *testing.T) {
		timeFunc = func() time.Time {
			return timeIs
		}
		licenseClient, err := NewLicenseClient(context.TODO(), zerolog.Nop(), WithMeta(meta), WithLicenseFileOrDirectory(inputPath))
		require.NoError(t, err)
		err = licenseClient.ValidateLicense(context.TODO())
		if expectError == nil {
			require.NoError(t, err)
		} else {
			require.ErrorIs(t, err, expectError)
		}
	}
}

func TestValidateMarketplaceLicense(t *testing.T) {
	ctrl := gomock.NewController(t)
	m := mocks.NewMockAWSLicenseManagerInterface(ctrl)
	out := licensemanager.CheckoutLicenseOutput{}
	in := licenseInput{
		CheckoutLicenseInput: licensemanager.CheckoutLicenseInput{
			CheckoutType: types.CheckoutTypeProvisional,
			ProductSKU:   aws.String(awsProductSKU),
			Entitlements: []types.EntitlementData{
				{
					Name: aws.String("Unlimited"),
					Unit: types.EntitlementDataUnitNone,
				},
			},
			KeyFingerprint: aws.String("aws:294406891311:AWS/Marketplace:issuer-fingerprint"),
		},
	}

	assert.NoError(t, faker.FakeObject(&out))
	m.EXPECT().CheckoutLicense(gomock.Any(), in).Return(&out, nil)
	t.Setenv("CQ_AWS_MARKETPLACE_LICENSE", "true")

	licenseClient, err := NewLicenseClient(context.TODO(), zerolog.Nop(), WithAWSLicenseManagerClient(m))
	require.NoError(t, err)
	require.NoError(t, licenseClient.ValidateLicense(context.TODO()))
}

type licenseInput struct {
	licensemanager.CheckoutLicenseInput
}

func (li licenseInput) Matches(x any) bool {
	testInput, ok := x.(*licensemanager.CheckoutLicenseInput)
	if !ok {
		return false
	}

	if testInput.CheckoutType != li.CheckoutType {
		return false
	}

	for i, ent := range testInput.Entitlements {
		if aws.ToString(ent.Name) != aws.ToString(li.Entitlements[i].Name) {
			return false
		}
		if aws.ToString(ent.Value) != aws.ToString(li.Entitlements[i].Value) {
			return false
		}
	}

	if aws.ToString(testInput.KeyFingerprint) != aws.ToString(li.KeyFingerprint) {
		return false
	}
	if aws.ToString(testInput.ProductSKU) != aws.ToString(li.ProductSKU) {
		return false
	}
	return true
}

func (li licenseInput) String() string {
	return fmt.Sprintf("{CheckoutType:%s Entitlements:%v KeyFingerprint:%s ProductSKU:%s}", li.CheckoutType, li.Entitlements, *li.KeyFingerprint, *li.ProductSKU)
}
