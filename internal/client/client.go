package client

import (
	"github.com/outscale/osc-sdk-go/v3/pkg/oks"
	"github.com/outscale/osc-sdk-go/v3/pkg/options"
	"github.com/outscale/osc-sdk-go/v3/pkg/osc"
	"github.com/outscale/osc-sdk-go/v3/pkg/profile"
	"github.com/outscale/terraform-provider-outscale/internal/logging"
)

type OutscaleClient struct {
	OKS *oks.Client
	OSC *osc.Client
}

type Config struct {
	AccessKey    string
	SecretKey    string
	Region       string
	APIEndpoint  string
	OKSEndpoint  string
	X509CertPath string
	X509KeyPath  string
	Insecure     bool
	UserAgent    string
	ConfigFile   string
	Profile      string
}

func NewOSCClient(cfg Config) (*osc.Client, error) {
	profile, err := profile.NewFrom(cfg.Profile, cfg.ConfigFile)
	if err != nil {
		return nil, err
	}

	if cfg.AccessKey != "" {
		profile.AccessKey = cfg.AccessKey
	}
	if cfg.SecretKey != "" {
		profile.SecretKey = cfg.SecretKey
	}
	if cfg.Region != "" {
		profile.Region = cfg.Region
	}
	if cfg.X509CertPath != "" {
		profile.X509ClientCert = cfg.X509CertPath
	}
	if cfg.X509KeyPath != "" {
		profile.X509ClientKey = cfg.X509KeyPath
	}
	if cfg.Insecure {
		profile.TlsSkipVerify = cfg.Insecure
	}
	if cfg.APIEndpoint != "" {
		profile.Endpoints.API = cfg.APIEndpoint
	}

	profile.Protocol = "https"

	logger := options.WithLogging(logging.NewTflogWrapper())
	userAgent := options.WithUseragent(cfg.UserAgent)

	return osc.NewClient(profile, userAgent, logger)
}

func NewOKSClient(cfg Config) (*oks.Client, error) {
	profile, err := profile.NewFrom(cfg.Profile, cfg.ConfigFile)
	if err != nil {
		return nil, err
	}

	if cfg.AccessKey != "" {
		profile.AccessKey = cfg.AccessKey
	}
	if cfg.SecretKey != "" {
		profile.SecretKey = cfg.SecretKey
	}
	if cfg.Region != "" {
		profile.Region = cfg.Region
	}
	if cfg.OKSEndpoint != "" {
		profile.Endpoints.OKS = cfg.OKSEndpoint
	}

	profile.Protocol = "https"
	profile.TlsSkipVerify = false
	profile.X509ClientCert = ""
	profile.X509ClientKey = ""

	logger := options.WithLogging(logging.NewTflogWrapper())
	userAgent := options.WithUseragent(cfg.UserAgent)

	return oks.NewClient(profile, userAgent, logger)
}
