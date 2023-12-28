package premium

import (
	"testing"
	"time"

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
