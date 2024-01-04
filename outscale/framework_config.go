package outscale

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/logging"
	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/terraform-providers/terraform-provider-outscale/utils"
	"github.com/terraform-providers/terraform-provider-outscale/version"
)

// OutscaleClient client
type OutscaleClient_fw struct {
	OSCAPI *oscgo.APIClient
}

// Client ...
func (c *frameworkProvider) Client_fw(ctx context.Context, data *ProviderModel, diags *diag.Diagnostics) (*OutscaleClient_fw, error) {

	//setdefaut environnemant varibles
	setDefaultEnv(data)
	tlsconfig := &tls.Config{InsecureSkipVerify: c.insecure}
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

	skipClient.Transport = logging.NewTransport("Outscale", skipClient.Transport)

	skipClient.Transport = NewTransport(data.AccessKeyId.ValueString(), data.SecretKeyId.ValueString(), data.Region.ValueString(), skipClient.Transport)

	basePath := fmt.Sprintf("api.%s.outscale.com", data.Region.ValueString())
	fmt.Printf("\n CONF_CLI: basePath=  %##++v\n", basePath)
	if endpoint, ok := data.Endpoints["api"]; ok {
		basePath = endpoint.(string)
	}
	fmt.Printf("\n APRES CONF_CLI: basePath=  %##++v\n", basePath)

	oscConfig := oscgo.NewConfiguration()
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
	/*
		if data.Endpoints.IsNull() {
			if endpoints := getEnvVariableValue([]string{"OSC_ENDPOINT_API", "OUTSCALE_OAPI_URL"}); endpoints != "" {
				data.Endpoints = types.StringValue(endpoints)
			}
		}
	*/
}
