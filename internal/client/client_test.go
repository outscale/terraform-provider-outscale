package client_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/outscale/osc-sdk-go/v3/pkg/osc"
	"github.com/outscale/osc-sdk-go/v3/pkg/profile"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOSCProfileConfigPriority(t *testing.T) {
	clearProfileEnv(t)

	t.Setenv("OSC_ACCESS_KEY", "env-ak")
	t.Setenv("OSC_SECRET_KEY", "env-sk")
	t.Setenv("OSC_REGION", "env-region")
	t.Setenv("OSC_ENDPOINT_API", "https://env.api")
	t.Setenv("OSC_X509_CLIENT_CERT", "/env/cert.pem")
	t.Setenv("OSC_X509_CLIENT_KEY", "/env/key.pem")

	cfg := client.Config{
		AccessKey:    "config-ak",
		SecretKey:    "config-sk",
		Region:       "config-region",
		APIEndpoint:  "https://config.api",
		X509CertPath: "/config/cert.pem",
		X509KeyPath:  "/config/key.pem",
	}

	p, err := cfg.NewProfile(cfg.ToOSCOption())
	require.NoError(t, err)

	assert.Equal(t, "config-ak", p.AccessKey, "access key")
	assert.Equal(t, "config-sk", p.SecretKey, "secret key")
	assert.Equal(t, "config-region", p.Region, "region")
	assert.Equal(t, "https://config.api", p.Endpoints.API, "api endpoint")
	assert.Equal(t, "/config/cert.pem", p.X509ClientCert, "x509 cert")
	assert.Equal(t, "/config/key.pem", p.X509ClientKey, "x509 key")
	assert.Equal(t, "https", p.Protocol, "protocol")
}

func TestOKSProfileConfigPriority(t *testing.T) {
	clearProfileEnv(t)

	t.Setenv("OSC_ACCESS_KEY", "env-ak")
	t.Setenv("OSC_SECRET_KEY", "env-sk")
	t.Setenv("OSC_REGION", "env-region")
	t.Setenv("OSC_ENDPOINT_OKS", "https://env.oks")
	t.Setenv("OSC_X509_CLIENT_CERT", "/env/cert.pem")
	t.Setenv("OSC_X509_CLIENT_KEY", "/env/key.pem")
	t.Setenv("OSC_TLS_SKIP_VERIFY", "true")

	cfg := client.Config{
		AccessKey:    "config-ak",
		SecretKey:    "config-sk",
		Region:       "config-region",
		OKSEndpoint:  "https://config.oks",
		X509CertPath: "/config/cert.pem",
		X509KeyPath:  "/config/key.pem",
	}

	p, err := cfg.NewProfile(cfg.ToOKSOption())
	require.NoError(t, err)

	assert.Equal(t, "config-ak", p.AccessKey, "access key")
	assert.Equal(t, "config-sk", p.SecretKey, "secret key")
	assert.Equal(t, "config-region", p.Region, "region")
	assert.Equal(t, "https://config.oks", p.Endpoints.OKS, "oks endpoint")
	assert.Equal(t, "https", p.Protocol, "protocol")
	assert.Equal(t, "/env/cert.pem", p.X509ClientCert, "x509 cert")
	assert.Equal(t, "/env/key.pem", p.X509ClientKey, "x509 key")
}

func TestProfileMerges(t *testing.T) {
	clearProfileEnv(t)

	configPath := filepath.Join(t.TempDir(), "config.json")
	configFile := profile.ConfigFile{
		Path: configPath,
		Profiles: map[string]profile.Profile{
			"default": {
				AccessKey:      "default-ak",
				SecretKey:      "default-sk",
				Region:         "default-region",
				X509ClientCert: "/default/cert.pem",
				X509ClientKey:  "/default/key.pem",
				Endpoints: profile.Endpoint{
					API: "https://default.api",
				},
			},
			"main": {
				AccessKey:      "main-ak",
				Region:         "main-region",
				X509ClientCert: "/main/cert.pem",
				X509ClientKey:  "/main/key.pem",
				Endpoints: profile.Endpoint{
					API: "https://main.api",
				},
			},
		},
	}
	require.NoError(t, configFile.Save())

	t.Setenv("OSC_ACCESS_KEY", "env-ak")
	t.Setenv("OSC_SECRET_KEY", "env-sk")
	t.Setenv("OSC_REGION", "env-region")
	t.Setenv("OSC_ENDPOINT_API", "https://env.api")
	t.Setenv("OSC_X509_CLIENT_CERT", "/env/cert.pem")
	t.Setenv("OSC_X509_CLIENT_KEY", "/env/key.pem")

	cfg := client.Config{
		AccessKey:  "config-ak",
		ConfigFile: configPath,
		Profile:    "main",
	}

	p, err := cfg.NewProfile(cfg.ToOSCOption())
	require.NoError(t, err)

	assert.Equal(t, "config-ak", p.AccessKey, "access key")
	assert.Equal(t, "env-sk", p.SecretKey, "secret key")
	assert.Equal(t, "main-region", p.Region, "region")
	assert.Equal(t, "https://main.api", p.Endpoints.API, "api endpoint")
	assert.Equal(t, "/main/cert.pem", p.X509ClientCert, "x509 cert")
	assert.Equal(t, "/main/key.pem", p.X509ClientKey, "x509 key")
	assert.Equal(t, "https", p.Protocol, "protocol")
}

func TestOSCApiEndpoint(t *testing.T) {
	clearProfileEnv(t)

	newClient := func(endpoint string) *osc.Client {
		cfg := client.Config{
			Region:      "eu-west-2",
			APIEndpoint: endpoint,
		}

		oscClient, err := client.NewOSCClient(cfg)
		require.NoError(t, err)

		return oscClient
	}

	t.Run("host as endpoint", func(t *testing.T) {
		host := "api.eu-west-2.outscale.com"
		client := newClient(host)

		_, err := client.ReadRegions(t.Context(), osc.ReadRegionsRequest{})
		require.NoError(t, err)
	})
	t.Run("https host as endpoint", func(t *testing.T) {
		host := "https://api.eu-west-2.outscale.com"
		client := newClient(host)

		_, err := client.ReadRegions(t.Context(), osc.ReadRegionsRequest{})
		require.NoError(t, err)
	})
	t.Run("full api endpoint", func(t *testing.T) {
		endpoint := "https://api.eu-west-2.outscale.com/api/v1"
		client := newClient(endpoint)

		_, err := client.ReadRegions(t.Context(), osc.ReadRegionsRequest{})
		require.NoError(t, err)
	})
	t.Run("https host as endpoint", func(t *testing.T) {
		host := "https://api.eu-west-2.outscale.com"
		client := newClient(host)

		_, err := client.ReadRegions(t.Context(), osc.ReadRegionsRequest{})
		require.NoError(t, err)
	})
}

func clearProfileEnv(t *testing.T) {
	t.Helper()

	for _, key := range []string{
		"OSC_ACCESS_KEY",
		"OSC_SECRET_KEY",
		"OSC_ACCESS_KEY_V2",
		"OSC_SECRET_KEY_V2",
		"OSC_X509_CLIENT_CERT",
		"OSC_X509_CLIENT_CERT_B64",
		"OSC_X509_CLIENT_KEY",
		"OSC_X509_CLIENT_KEY_B64",
		"OSC_TLS_SKIP_VERIFY",
		"OSC_LOGIN",
		"OSC_PASSWORD",
		"OSC_PROTOCOL",
		"OSC_REGION",
		"OSC_ENDPOINT_API",
		"OSC_ENDPOINT_OKS",
		"OSC_ENDPOINT_LBU",
		"OSC_ENDPOINT_OOS",
		"OSC_ENDPOINT_FCU",
		"OSC_ENDPOINT_EIM",
		"OSC_ENDPOINT_DIRECT_LINK",
		"OSC_PROFILE",
		"OSC_CONFIG_FILE",
	} {
		t.Setenv(key, "")
		err := os.Unsetenv(key)
		require.NoError(t, err)
	}
}
