package client_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/outscale/osc-sdk-go/v3/pkg/profile"
	"github.com/outscale/terraform-provider-outscale/internal/client"
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
	if err != nil {
		t.Fatalf("NewProfile: %v", err)
	}

	assertEqual(t, "access key", p.AccessKey, "config-ak")
	assertEqual(t, "secret key", p.SecretKey, "config-sk")
	assertEqual(t, "region", p.Region, "config-region")
	assertEqual(t, "api endpoint", p.Endpoints.API, "https://config.api")
	assertEqual(t, "x509 cert", p.X509ClientCert, "/config/cert.pem")
	assertEqual(t, "x509 key", p.X509ClientKey, "/config/key.pem")
	assertEqual(t, "protocol", p.Protocol, "https")
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
	if err != nil {
		t.Fatalf("NewProfile: %v", err)
	}

	assertEqual(t, "access key", p.AccessKey, "config-ak")
	assertEqual(t, "secret key", p.SecretKey, "config-sk")
	assertEqual(t, "region", p.Region, "config-region")
	assertEqual(t, "oks endpoint", p.Endpoints.OKS, "https://config.oks")
	assertEqual(t, "protocol", p.Protocol, "https")
	assertEqual(t, "x509 cert", p.X509ClientCert, "/env/cert.pem")
	assertEqual(t, "x509 key", p.X509ClientKey, "/env/key.pem")
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
	if err := configFile.Save(); err != nil {
		t.Fatalf("save config: %v", err)
	}

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
	if err != nil {
		t.Fatalf("NewProfile: %v", err)
	}

	assertEqual(t, "access key", p.AccessKey, "config-ak")
	assertEqual(t, "secret key", p.SecretKey, "env-sk")
	assertEqual(t, "region", p.Region, "main-region")
	assertEqual(t, "api endpoint", p.Endpoints.API, "https://main.api")
	assertEqual(t, "x509 cert", p.X509ClientCert, "/main/cert.pem")
	assertEqual(t, "x509 key", p.X509ClientKey, "/main/key.pem")
	assertEqual(t, "protocol", p.Protocol, "https")
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
		if err != nil {
			t.Fatalf("Unsetenv: %v", err)
		}
	}
}

func assertEqual(t *testing.T, field, got, want string) {
	t.Helper()

	if got != want {
		t.Fatalf("expected %s %q, got %q", field, want, got)
	}
}
