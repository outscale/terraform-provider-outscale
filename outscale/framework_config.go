package outscale

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	oscgo "github.com/outscale/osc-sdk-go/v2"
	sdkv3_oks "github.com/outscale/osc-sdk-go/v3/pkg/oks"
	sdkv3_profile "github.com/outscale/osc-sdk-go/v3/pkg/profile"
	sdkv3_utils "github.com/outscale/osc-sdk-go/v3/pkg/utils"
	"github.com/outscale/terraform-provider-outscale/utils"
	"github.com/tidwall/gjson"
)

// OutscaleClient client
type OutscaleClientFW struct {
	OSCAPI *oscgo.APIClient
	OKS    *sdkv3_oks.Client
}

func newOKSClientFW(data *ProviderModel) (oksClient *sdkv3_oks.Client, err error) {
	profile := sdkv3_profile.Profile{
		AccessKey:      data.AccessKeyId.ValueString(),
		SecretKey:      data.SecretKeyId.ValueString(),
		Region:         data.Region.ValueString(),
		X509ClientCert: data.X509CertPath.ValueString(),
		X509ClientKey:  data.X509CertPath.ValueString(),
		TlsSkipVerify:  data.Insecure.ValueBool(),
		Protocol:       "https",
		Endpoints: sdkv3_profile.Endpoint{
			OKS: data.Endpoints[0].OKS.ValueString(),
			API: data.Endpoints[0].API.ValueString(),
		},
	}
	logger := sdkv3_utils.WithLogging(utils.NewTflogWrapper())

	return sdkv3_oks.NewClient(&profile, sdkv3_utils.WithUseragent(UserAgent), logger)
}

func newAPIClientFW(data *ProviderModel) (apiClient *oscgo.APIClient, err error) {
	if data.Region.IsNull() && len(data.Endpoints) == 0 {
		return nil, errors.New("'region' or 'endpoints' must be set for provider configuration")
	}
	oscConfig := oscgo.NewConfiguration()

	tlsConfig := ClientTLSConfig(data.Insecure.ValueBool(), data.X509CertPath.ValueString(), data.X509KeyPath.ValueString())
	httpClient := ClientHTTPConfig(tlsConfig)
	httpClient.Transport = NewTransport(data.AccessKeyId.ValueString(), data.SecretKeyId.ValueString(), data.Region.ValueString(), httpClient.Transport)
	endpoint := ClientEndpointConfig(oscConfig, data.Endpoints[0].API.ValueString(), data.Region.ValueString())
	data.Endpoints[0].API = types.StringValue(endpoint)

	oscConfig.Debug = true
	oscConfig.HTTPClient = httpClient
	oscConfig.Host = endpoint
	oscConfig.UserAgent = UserAgent

	return oscgo.NewAPIClient(oscConfig), nil
}

// Client ...
func (c *frameworkProvider) ClientFW(ctx context.Context, data *ProviderModel, diags *diag.Diagnostics) (*OutscaleClientFW, error) {
	ok, err := isProfileSet(data)
	if err != nil {
		return nil, err
	}
	if !ok {
		setDefaultEnv(data)
	}

	oscClient, err := newAPIClientFW(data)
	if err != nil {
		return nil, err
	}
	oksClient, err := newOKSClientFW(data)
	if err != nil {
		return nil, err
	}
	client := &OutscaleClientFW{
		OSCAPI: oscClient,
		OKS:    oksClient,
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
			endp := make([]Endpoints, 1)
			if endpoint := endpoints["api"].(string); endpoint != "" {
				endp[0].API = types.StringValue(endpoint)
			}
			if endpoint := endpoints["oks"].(string); endpoint != "" {
				endp[0].OKS = types.StringValue(endpoint)
			}
			data.Endpoints = endp
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
		endp := make([]Endpoints, 1)
		if endpoint := utils.GetEnvVariableValue([]string{"OSC_ENDPOINT_API", "OUTSCALE_OAPI_URL"}); endpoint != "" {
			endp[0].API = types.StringValue(endpoint)
		}
		if endpoint := utils.GetEnvVariableValue([]string{"OSC_ENDPOINT_OKS"}); endpoint != "" {
			endp[0].OKS = types.StringValue(endpoint)
		}
		data.Endpoints = endp
	}
}
