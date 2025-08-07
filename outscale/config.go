package outscale

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/logging"
	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/outscale/terraform-provider-outscale/version"
)

// Config ...
type Config struct {
	AccessKeyID  string
	SecretKeyID  string
	Region       string
	TokenID      string
	Endpoints    map[string]string
	X509CertPath string
	X509KeyPath  string
	Insecure     bool
	ConfigFile   string
	Profile      string
}

// OutscaleClient client
type OutscaleClient struct {
	OSCAPI *oscgo.APIClient
}

func ClientTLSConfig(insecure bool, certFile string, keyFile string) (tlsconfig *tls.Config) {
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

func ClientHTTPConfig(tlsConfig *tls.Config) (httpClient *http.Client) {
	httpClient = &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
			Proxy:           http.ProxyFromEnvironment,
		},
	}
	httpClient.Transport = logging.NewSubsystemLoggingHTTPTransport("Outscale", httpClient.Transport)

	return
}

func ClientEndpointConfig(oscConfig *oscgo.Configuration, endpoint string, region string) string {
	if endpoint != "" {
		if strings.Contains(endpoint, "://") {
			if scheme, host, found := strings.Cut(endpoint, "://"); found {
				oscConfig.Scheme = scheme
				endpoint = host
			}
		}
		return endpoint
	} else {
		return fmt.Sprintf("api.%s.outscale.com", region)
	}
}

// Client ...
func (c *Config) Client() (*OutscaleClient, error) {
	endpoint := c.Endpoints["api"]
	if c.Region == "" && endpoint == "" {
		return nil, errors.New("'region' or 'endpoints' must be set for provider configuration")
	}
	oscConfig := oscgo.NewConfiguration()

	endpoint = ClientEndpointConfig(oscConfig, endpoint, c.Region)
	c.Endpoints["api"] = endpoint
	tlsConfig := ClientTLSConfig(c.Insecure, c.X509CertPath, c.X509KeyPath)
	httpClient := ClientHTTPConfig(tlsConfig)
	httpClient.Transport = NewTransport(c.AccessKeyID, c.SecretKeyID, c.Region, httpClient.Transport)

	oscConfig.Host = endpoint
	oscConfig.HTTPClient = httpClient
	oscConfig.Debug = true
	oscConfig.UserAgent = fmt.Sprintf("terraform-provider-outscale/%s", version.GetVersion())

	oscClient := oscgo.NewAPIClient(oscConfig)
	client := &OutscaleClient{
		OSCAPI: oscClient,
	}
	return client, nil
}
