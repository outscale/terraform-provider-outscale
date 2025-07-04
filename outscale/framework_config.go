package outscale

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/logging"
	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/outscale/terraform-provider-outscale/utils"
	"github.com/outscale/terraform-provider-outscale/version"
	"github.com/tidwall/gjson"
)

// OutscaleClient client
type OutscaleClient_fw struct {
	OSCAPI *oscgo.APIClient
}

// Client ...
func (c *frameworkProvider) Client_fw(ctx context.Context, data *ProviderModel, diags *diag.Diagnostics) (*OutscaleClient_fw, error) {
	ok, err := isProfileSet(data)
	if err != nil {
		return nil, err
	}
	if !ok {
		setDefaultEnv(data)
	}

	tlsconfig := &tls.Config{InsecureSkipVerify: data.Insecure.ValueBool()}
	cert, err := tls.LoadX509KeyPair(data.X509CertPath.ValueString(), data.X509KeyPath.ValueString())
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

	if data.Region.IsNull() && len(data.Endpoints) == 0 {
		return nil, errors.New("'region' or 'endpoints' must be set for provider configuration")
	}

	oscConfig := oscgo.NewConfiguration()
	basePath := fmt.Sprintf("api.%s.outscale.com", data.Region.ValueString())

	if len(data.Endpoints) > 0 {
		basePath = data.Endpoints[0].API.ValueString()
		if strings.Contains(basePath, "://") {
			if scheme, host, found := strings.Cut(basePath, "://"); found {
				oscConfig.Scheme = scheme
				basePath = host
			}
		}
		endpointSplit := strings.Split(basePath, ".")
		data.Region = types.StringValue(endpointSplit[1])
	}

	skipClient.Transport = NewTransport(data.AccessKeyId.ValueString(), data.SecretKeyId.ValueString(), data.Region.ValueString(), skipClient.Transport)
	oscConfig.Debug = true
	oscConfig.HTTPClient = skipClient
	oscConfig.Host = basePath
	oscConfig.UserAgent = fmt.Sprintf("terraform-provider-outscale/%s", version.GetVersion())
	oscClient := oscgo.NewAPIClient(oscConfig)
	client := &OutscaleClient_fw{
		OSCAPI: oscClient,
	}
	return client, nil
}

func isProfileSet(data *ProviderModel) (bool, error) {
	isProfSet := false
	if profileName, ok := os.LookupEnv("OSC_PROFILE"); ok || !data.Profile.IsNull() {
		if data.Profile.ValueString() != "" {
			profileName = data.Profile.ValueString()
		}

		var configFilePath string
		if envPath, ok := os.LookupEnv("OSC_CONFIG_FILE"); ok || !data.ConfigFile.IsNull() {
			if data.ConfigFile.ValueString() != "" {
				configFilePath = data.ConfigFile.ValueString()
			} else {
				configFilePath = envPath
			}
		} else {
			homePath, err := os.UserHomeDir()
			if err != nil {
				return isProfSet, err
			}
			configFilePath = homePath + utils.SuffixConfigFilePath
		}
		jsonFile, err := os.ReadFile(configFilePath)
		if err != nil {
			return isProfSet, fmt.Errorf("unable to read config file '%v', Error: %w", configFilePath, err)
		}
		profile := gjson.GetBytes(jsonFile, profileName)
		if !gjson.Valid(profile.String()) {
			return isProfSet, errors.New("invalid json profile file")
		}
		if !profile.Get("access_key").Exists() ||
			!profile.Get("secret_key").Exists() {
			return isProfSet, errors.New("profile 'access_key' or 'secret_key' are not defined! ")
		}
		setProfile(data, profile)
		isProfSet = true
	}
	return isProfSet, nil
}

func setProfile(data *ProviderModel, profile gjson.Result) {
	if data.AccessKeyId.IsNull() {
		if accessKeyId := profile.Get("access_key").String(); accessKeyId != "" {
			data.AccessKeyId = types.StringValue(accessKeyId)
		}
	}
	if data.SecretKeyId.IsNull() {
		if secretKeyId := profile.Get("secret_key").String(); secretKeyId != "" {
			data.SecretKeyId = types.StringValue(secretKeyId)
		}
	}
	if data.Region.IsNull() {
		if profile.Get("region").Exists() {
			if region := profile.Get("region").String(); region != "" {
				data.Region = types.StringValue(region)
			}
		}
	}
	if data.X509CertPath.IsNull() {
		if profile.Get("x509_cert_path").Exists() {
			if x509Cert := profile.Get("x509_cert_path").String(); x509Cert != "" {
				data.X509CertPath = types.StringValue(x509Cert)
			}
		}
	}
	if data.X509KeyPath.IsNull() {
		if profile.Get("x509_key_path").Exists() {
			if x509Key := profile.Get("x509_key_path").String(); x509Key != "" {
				data.X509KeyPath = types.StringValue(x509Key)
			}
		}
	}
	if len(data.Endpoints) == 0 {
		if profile.Get("endpoints").Exists() {
			endpoints := profile.Get("endpoints").Value().(map[string]interface{})
			if endpoint := endpoints["api"].(string); endpoint != "" {
				endp := make([]Endpoints, 1)
				endp[0].API = types.StringValue(endpoint)
				data.Endpoints = endp
			}
		}
	}
}

func setDefaultEnv(data *ProviderModel) {
	if data.AccessKeyId.IsNull() {
		if accessKeyId := utils.GetEnvVariableValue([]string{"OSC_ACCESS_KEY", "OUTSCALE_ACCESSKEYID"}); accessKeyId != "" {
			data.AccessKeyId = types.StringValue(accessKeyId)
		}
	}
	if data.SecretKeyId.IsNull() {
		if secretKeyId := utils.GetEnvVariableValue([]string{"OSC_SECRET_KEY", "OUTSCALE_SECRETKEYID"}); secretKeyId != "" {
			data.SecretKeyId = types.StringValue(secretKeyId)
		}
	}

	if data.Region.IsNull() {
		if region := utils.GetEnvVariableValue([]string{"OSC_REGION", "OUTSCALE_REGION"}); region != "" {
			data.Region = types.StringValue(region)
		}
	}

	if data.X509CertPath.IsNull() {
		if x509Cert := utils.GetEnvVariableValue([]string{"OSC_X509_CLIENT_CERT", "OUTSCALE_X509CERT"}); x509Cert != "" {
			data.X509CertPath = types.StringValue(x509Cert)
		}
	}

	if data.X509KeyPath.IsNull() {
		if x509Key := utils.GetEnvVariableValue([]string{"OSC_X509_CLIENT_KEY", "OUTSCALE_X509KEY"}); x509Key != "" {
			data.X509KeyPath = types.StringValue(x509Key)
		}
	}
	if len(data.Endpoints) == 0 {
		if endpoint := utils.GetEnvVariableValue([]string{"OSC_ENDPOINT_API", "OUTSCALE_OAPI_URL"}); endpoint != "" {
			endp := make([]Endpoints, 1)
			endp[0].API = types.StringValue(endpoint)
			data.Endpoints = endp
		}
	}
}
