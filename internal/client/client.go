package client

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/logging"
	oscgo "github.com/outscale/osc-sdk-go/v2"
	sdkv3_oks "github.com/outscale/osc-sdk-go/v3/pkg/oks"
	sdkv3_profile "github.com/outscale/osc-sdk-go/v3/pkg/profile"
	sdkv3_utils "github.com/outscale/osc-sdk-go/v3/pkg/utils"
	tflogging "github.com/outscale/terraform-provider-outscale/internal/logging"
	"github.com/outscale/terraform-provider-outscale/internal/transport"
)

type OutscaleClient struct {
	OSCAPI *oscgo.APIClient
	OKS    *sdkv3_oks.Client
}

type Config struct {
	AccessKeyID  string
	SecretKeyID  string
	Region       string
	APIEndpoint  string
	OKSEndpoint  string
	X509CertPath string
	X509KeyPath  string
	Insecure     bool
	UserAgent    string
}

func tlsConfig(insecure bool, certFile string, keyFile string) (tlsconfig *tls.Config) {
	tlsconfig = &tls.Config{InsecureSkipVerify: insecure}
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err == nil {
		tlsconfig = &tls.Config{
			InsecureSkipVerify: false,
			Certificates:       []tls.Certificate{cert},
		}
	}
	return
}

func httpConfig(tlsConfig *tls.Config, accessKeyID, secretKeyID, region string) (httpClient *http.Client) {
	httpClient = &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
			Proxy:           http.ProxyFromEnvironment,
		},
	}
	httpClient.Transport = logging.NewSubsystemLoggingHTTPTransport("Outscale", httpClient.Transport)
	securityProvider := transport.NewSecurityProviderAWSv4(accessKeyID, secretKeyID, "", "oapi", region)
	httpClient.Transport = transport.NewTransport(httpClient.Transport, securityProvider)

	return
}

func endpointConfig(oscConfig *oscgo.Configuration, endpoint string, region string) string {
	if endpoint == "" {
		return fmt.Sprintf("api.%s.outscale.com", region)
	}
	if strings.Contains(endpoint, "://") {
		if scheme, host, found := strings.Cut(endpoint, "://"); found {
			oscConfig.Scheme = scheme
			endpoint = host
		}
	}
	return endpoint
}

func NewOAPIClient(cfg Config) (*oscgo.APIClient, error) {
	oscConfig := oscgo.NewConfiguration()

	endpoint := endpointConfig(oscConfig, cfg.APIEndpoint, cfg.Region)
	tlsConfig := tlsConfig(cfg.Insecure, cfg.X509CertPath, cfg.X509KeyPath)
	httpClient := httpConfig(tlsConfig, cfg.AccessKeyID, cfg.SecretKeyID, cfg.Region)

	oscConfig.Host = endpoint
	oscConfig.HTTPClient = httpClient
	oscConfig.Debug = true
	oscConfig.UserAgent = cfg.UserAgent

	return oscgo.NewAPIClient(oscConfig), nil
}

func NewOKSClient(cfg Config) (*sdkv3_oks.Client, error) {
	profile := sdkv3_profile.Profile{
		AccessKey:     cfg.AccessKeyID,
		SecretKey:     cfg.SecretKeyID,
		Region:        cfg.Region,
		TlsSkipVerify: cfg.Insecure,
		Protocol:      "https",
		Endpoints: sdkv3_profile.Endpoint{
			OKS: cfg.OKSEndpoint,
			API: cfg.APIEndpoint,
		},
	}

	logger := sdkv3_utils.WithLogging(tflogging.NewTflogWrapper())
	return sdkv3_oks.NewClient(&profile, sdkv3_utils.WithUseragent(cfg.UserAgent), logger)
}
