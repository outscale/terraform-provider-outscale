package provider

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/framework/fwhelpers"
	"github.com/outscale/terraform-provider-outscale/internal/framework/fwhelpers/to"
	"github.com/outscale/terraform-provider-outscale/internal/utils"
	"github.com/tidwall/gjson"
)

func (c *FrameworkProvider) ClientFW(ctx context.Context, data *ProviderModel, diags *diag.Diagnostics) (*client.OutscaleClient, error) {
	loadConfigFromEnv(data)
	err := mergeProfileConfig(data)
	if err != nil {
		return nil, err
	}

	var apiEndpoint, apiRegion, apiCertPath, apiKeyPath string
	var apiInsecure bool
	if fwhelpers.IsSet(data.API) {
		apiModel, diag := to.Slice[APIModel](ctx, data.API)
		diags.Append(diag...)
		if diags.HasError() {
			return nil, errors.New("failed to extract api configuration")
		}
		if len(apiModel) > 0 {
			apiEndpoint = apiModel[0].Endpoint.ValueString()
			apiRegion = apiModel[0].Region.ValueString()
			apiCertPath = apiModel[0].X509CertPath.ValueString()
			apiKeyPath = apiModel[0].X509KeyPath.ValueString()
			apiInsecure = apiModel[0].Insecure.ValueBool()
		}
	}
	// fallback to deprecated configuration
	if apiEndpoint == "" && len(data.Endpoints) > 0 {
		apiEndpoint = data.Endpoints[0].API.ValueString()
	}
	if apiCertPath == "" {
		apiCertPath = data.X509CertPath.ValueString()
	}
	if apiKeyPath == "" {
		apiKeyPath = data.X509KeyPath.ValueString()
	}
	if apiInsecure == false {
		apiInsecure = data.Insecure.ValueBool()
	}
	if apiRegion == "" {
		apiRegion = data.Region.ValueString()
	}

	oksEndpoint := ""
	oksRegion := ""
	if fwhelpers.IsSet(data.OKS) {
		oksModels, diag := to.Slice[OKSModel](ctx, data.OKS)
		diags.Append(diag...)
		if diags.HasError() {
			return nil, errors.New("failed to extract oks configuration")
		}
		if len(oksModels) > 0 {
			oksEndpoint = oksModels[0].Endpoint.ValueString()
			oksRegion = oksModels[0].Region.ValueString()
		}
	}
	// fallback to deprecated configuration
	if oksEndpoint == "" && len(data.Endpoints) > 0 {
		oksEndpoint = data.Endpoints[0].OKS.ValueString()
	}
	if oksRegion == "" {
		oksRegion = data.Region.ValueString()
	}

	oscClient, err := client.NewOAPIClient(client.Config{
		AccessKeyID:  data.AccessKeyId.ValueString(),
		SecretKeyID:  data.SecretKeyId.ValueString(),
		Region:       apiRegion,
		APIEndpoint:  apiEndpoint,
		X509CertPath: apiCertPath,
		X509KeyPath:  apiKeyPath,
		Insecure:     apiInsecure,
		UserAgent:    UserAgent,
	})
	if err != nil {
		return nil, err
	}

	oksClient, err := client.NewOKSClient(client.Config{
		AccessKeyID: data.AccessKeyId.ValueString(),
		SecretKeyID: data.SecretKeyId.ValueString(),
		Region:      oksRegion,
		OKSEndpoint: oksEndpoint,
		UserAgent:   UserAgent,
	})
	if err != nil {
		return nil, err
	}

	outscaleClient := &client.OutscaleClient{
		OSCAPI: oscClient,
		OKS:    oksClient,
	}
	return outscaleClient, nil
}

func mergeProfileConfig(data *ProviderModel) error {
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
				return err
			}
			configFilePath = homePath + utils.SuffixConfigFilePath
		}
		jsonFile, err := os.ReadFile(configFilePath)
		if err != nil {
			return fmt.Errorf("unable to read config file '%v', Error: %w", configFilePath, err)
		}
		profile := gjson.GetBytes(jsonFile, profileName)
		if !gjson.Valid(profile.String()) {
			return errors.New("invalid json profile file")
		}
		if !profile.Get("access_key").Exists() ||
			!profile.Get("secret_key").Exists() {
			return errors.New("profile 'access_key' or 'secret_key' are not defined! ")
		}
		setProfile(data, profile)
	}
	return nil
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
			if endpoint, ok := endpoints["api"].(string); ok && endpoint != "" {
				endp[0].API = types.StringValue(endpoint)
			}
			if endpoint, ok := endpoints["oks"].(string); ok && endpoint != "" {
				endp[0].OKS = types.StringValue(endpoint)
			}
			data.Endpoints = endp
		}
	}
}

func loadConfigFromEnv(data *ProviderModel) {
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
		hasEndpoint := false
		if endpoint := utils.GetEnvVariableValue([]string{"OSC_ENDPOINT_API", "OUTSCALE_OAPI_URL"}); endpoint != "" {
			endp[0].API = types.StringValue(endpoint)
			hasEndpoint = true
		}
		if endpoint := utils.GetEnvVariableValue([]string{"OSC_ENDPOINT_OKS"}); endpoint != "" {
			endp[0].OKS = types.StringValue(endpoint)
			hasEndpoint = true
		}
		if hasEndpoint {
			data.Endpoints = endp
		}
	}
}
