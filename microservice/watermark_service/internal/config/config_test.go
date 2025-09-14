package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoadConfig(t *testing.T) {
	t.Setenv("APP_PORT", "test")
	t.Setenv("AWS_REGION", "test")
	t.Setenv("LEAPCELL_BASE_ENDPOINT", "test")
	t.Setenv("LEAPCELL_CDN", "test")
	t.Setenv("AWS_ACCESS_KEY_ID", "test")
	t.Setenv("AWS_SECRET_ACCESS_KEY", "test")
	t.Setenv("AWS_BUCKET", "test")

	_, err := LoadConfig()

	require.NoError(t, err)
}

func TestLoadConfig_Fail(t *testing.T) {
	_, err := LoadConfig()
	require.Error(t, err)
}
