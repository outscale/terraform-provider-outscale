package client

import (
	"fmt"

	"github.com/outscale/goutils/sdk/batch"
	"github.com/outscale/osc-sdk-go/v3/pkg/oks"
	"github.com/outscale/osc-sdk-go/v3/pkg/options"
	"github.com/outscale/osc-sdk-go/v3/pkg/osc"
	"github.com/outscale/osc-sdk-go/v3/pkg/profile"
	"github.com/outscale/terraform-provider-outscale/internal/logging"
)

type OutscaleClient struct {
	OKS *oks.Client
	OSC *osc.Client

	// OSC API Batchers
	VmBatcher            *batch.BatcherByID[osc.Vm]
	VolumeBatcher        *batch.BatcherByID[osc.Volume]
	SecurityGroupBatcher *batch.BatcherByID[osc.SecurityGroup]
	NetBatcher           *batch.BatcherByID[osc.Net]
	SubnetBatcher        *batch.BatcherByID[osc.Subnet]
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

func (cfg Config) ToOSCOption() profile.Option {
	return func(profile *profile.Profile) error {
		profile.AccessKey = cfg.AccessKey
		profile.SecretKey = cfg.SecretKey
		profile.Region = cfg.Region
		profile.X509ClientCert = cfg.X509CertPath
		profile.X509ClientKey = cfg.X509KeyPath
		profile.TlsSkipVerify = cfg.Insecure
		profile.Endpoints.API = cfg.APIEndpoint
		profile.Protocol = "https"

		return nil
	}
}

// NewProfile merges profile values in the following order:
// provider config -> selected profile file -> environment values for any fields that are still unset
func (cfg Config) NewProfile(configOption profile.Option) (*profile.Profile, error) {
	opts := []profile.Option{configOption, profile.MergeWith(profile.FromEnv())}

	if cfg.Profile != "" || cfg.ConfigFile != "" {
		opts = []profile.Option{configOption, profile.MergeWith(profile.FromFile(cfg.Profile, cfg.ConfigFile)), profile.MergeWith(profile.FromEnv())}
	}

	return profile.New(opts...)
}

func NewOSCClient(cfg Config) (*osc.Client, error) {
	profile, err := cfg.NewProfile(cfg.ToOSCOption())
	if err != nil {
		return nil, fmt.Errorf("new profile: %w", err)
	}

	logger := options.WithLogging(logging.NewTflogWrapper())
	userAgent := options.WithUseragent(cfg.UserAgent)

	return osc.NewClient(profile, userAgent, logger)
}

func (cfg Config) ToOKSOption() profile.Option {
	return func(profile *profile.Profile) error {
		profile.AccessKey = cfg.AccessKey
		profile.SecretKey = cfg.SecretKey
		profile.Region = cfg.Region
		profile.Endpoints.OKS = cfg.OKSEndpoint
		profile.TlsSkipVerify = false
		profile.Protocol = "https"
		profile.X509ClientCert = ""
		profile.X509ClientKey = ""

		return nil
	}
}

func NewOKSClient(cfg Config) (*oks.Client, error) {
	profile, err := cfg.NewProfile(cfg.ToOKSOption())
	if err != nil {
		return nil, fmt.Errorf("new profile: %w", err)
	}

	logger := options.WithLogging(logging.NewTflogWrapper())
	userAgent := options.WithUseragent(cfg.UserAgent)

	return oks.NewClient(profile, userAgent, logger)
}
