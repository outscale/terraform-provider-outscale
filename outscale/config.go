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

// Client ...
func (c *Config) Client() (*OutscaleClient, error) {
	tlsconfig := &tls.Config{InsecureSkipVerify: c.Insecure}
	cert, err := tls.LoadX509KeyPair(c.X509CertPath, c.X509KeyPath)
	if err == nil {
		tlsconfig = &tls.Config{
			InsecureSkipVerify: false,
			Certificates:       []tls.Certificate{cert},
		}
	}

	skipClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: tlsconfig,
			Proxy:           http.ProxyFromEnvironment,
		},
	}

	skipClient.Transport = logging.NewSubsystemLoggingHTTPTransport("Outscale", skipClient.Transport)
	endpoint := c.Endpoints["api"]
	if c.Region == "" && endpoint == "" {
		return nil, errors.New("'region' or 'endpoints' must be set for provider configuration")
	}

	basePath := fmt.Sprintf("api.%s.outscale.com", c.Region)
	oscConfig := oscgo.NewConfiguration()

	if endpoint != "" {
		basePath = endpoint
		if strings.Contains(basePath, "://") {
			if scheme, host, found := strings.Cut(basePath, "://"); found {
				oscConfig.Scheme = scheme
				basePath = host
			}
		}
		endpointSplit := strings.Split(basePath, ".")
		c.Region = endpointSplit[1]
	}

	skipClient.Transport = NewTransport(c.AccessKeyID, c.SecretKeyID, c.Region, skipClient.Transport)
	oscConfig.Debug = true
	oscConfig.HTTPClient = skipClient
	oscConfig.Host = basePath
	oscConfig.UserAgent = fmt.Sprintf("terraform-provider-outscale/%s", version.GetVersion())

	oscClient := oscgo.NewAPIClient(oscConfig)
	client := &OutscaleClient{
		OSCAPI: oscClient,
	}
	return client, nil
}
