package outscale

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ provider.Provider = &frameworkProvider{}
)

func New(version string) provider.Provider {
	return &frameworkProvider{
		version: version,
	}
}

type frameworkProvider struct {
	accessKeyId    types.String
	secretKeyId    types.String
	region         types.String
	endpoints      []Endpoints
	x509CertPath   string
	x509KeyPath    string
	configFilePath string
	insecure       bool
	profile        string
	version        string
}

type ProviderModel struct {
	AccessKeyId    types.String `tfsdk:"access_key_id"`
	SecretKeyId    types.String `tfsdk:"secret_key_id"`
	Region         types.String `tfsdk:"region"`
	Endpoints      []Endpoints  `tfsdk:"endpoints"`
	X509CertPath   types.String `tfsdk:"x509_cert_path"`
	X509KeyPath    types.String `tfsdk:"x509_key_path"`
	ConfigFilePath types.String `tfsdk:"config_file_path"`
	Profile        types.String `tfsdk:"profile"`
	Insecure       types.Bool   `tfsdk:"insecure"`
}

type Endpoints struct {
	API types.String `tfsdk:"api"`
}

func (p *frameworkProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "outscale"
	resp.Version = p.version
}

func (p *frameworkProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Blocks: map[string]schema.Block{
			"endpoints": schema.SetNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"api": schema.StringAttribute{
							Optional:    true,
							Description: "The Endpoints for API operations.",
						},
					},
				},
			},
		},

		Attributes: map[string]schema.Attribute{
			"access_key_id": schema.StringAttribute{
				Optional:    true,
				Description: "The Access Key ID for API operations.",
			},
			"secret_key_id": schema.StringAttribute{
				Optional:    true,
				Description: "The Secret Key ID for API operations.",
			},
			"region": schema.StringAttribute{
				Optional:    true,
				Description: "The Region for API operations.",
			},
			"x509_cert_path": schema.StringAttribute{
				Optional:    true,
				Description: "The path to your x509 cert",
			},
			"x509_key_path": schema.StringAttribute{
				Optional:    true,
				Description: "The path to your x509 key",
			},
			"config_file_path": schema.StringAttribute{
				Optional:    true,
				Description: "The path to your configuration file in which you have defined your credentials.",
			},
			"profile": schema.StringAttribute{
				Optional:    true,
				Description: "The name of your profile in which you define your credencial",
			},
			"insecure": schema.BoolAttribute{
				Optional:    true,
				Description: "tls insecure connection",
			},
		},
	}
}

func (p *frameworkProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {

	var config ProviderModel

	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	if config.AccessKeyId.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("access_key_id"),
			"Unknown Outscale API AccessKeyId",
			"The provider cannot create the Outscale API client as there is an unknown configuration value for the Outscale API access_key_id. "+
				"Either target apply the source Outscale of the value first, set the value statically in the configuration, or use the 'OSC_ACCESS_KEY or OUTSCALE_ACCESSKEYID' environment variable.",
		)
	}

	if config.SecretKeyId.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("secret_key_id"),
			"Unknown HashiCups API SecretKeyId",
			"The provider cannot create the Outscale API client as there is an unknown configuration value for the Outscale API secret_key_id. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the 'OSC_SECRET_KEY or OUTSCALE_SECRETKEYID' environment variable.",
		)
	}

	if config.Region.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("region"),
			"Unknown Outscale API Region",
			"The provider cannot create the Outscale API client as there is an unknown configuration value for the Outscale API region. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the 'OSC_REGION or OUTSCALE_REGION' environment variable.",
		)
	}
	if config.X509CertPath.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("x509_cert_path"),
			"Unknown Outscale API X509CertPath",
			"The provider cannot create the Outscale API client as there is an unknown configuration value for the Outscale API x509_cert_path. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the 'OSC_X509_CLIENT_CERT or OUTSCALE_X509CERT' environment variable.",
		)
	}

	if config.X509KeyPath.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("x509_key_path"),
			"Unknown Outscale API X509KeyPath",
			"The provider cannot create the Outscale API client as there is an unknown configuration value for the Outscale API x509_key_path. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the 'OSC_X509_CLIENT_KEY or OUTSCALE_X509KEY' environment variable.",
		)
	}
	if config.ConfigFilePath.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("config_file_path"),
			"Unknown Outscale API ConfigFilePath",
			"The provider cannot create the Outscale API client as there is an unknown configuration value for the Outscale API profile. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the 'OSC_CONFIG_FILE' environment variable.",
		)
	}
	if config.Profile.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("profile"),
			"Unknown Outscale API profile",
			"The provider cannot create the Outscale API client as there is an unknown configuration value for the Outscale API profile. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the 'OSC_PROFILE' environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}
	client, err := p.Client_fw(ctx, &config, &diags)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create Outscale API Client",
			"An unexpected error occurred when creating the Outscale API client. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"Outscale Client Error: "+err.Error(),
		)
		return
	}
	resp.DataSourceData = *client
	resp.ResourceData = *client
}

func (p *frameworkProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewDataSourceQuota,
	}
}

func (p *frameworkProvider) Resources(ctx context.Context) []func() resource.Resource {
	/*return []func() resource.Resource{
		NewResource,
	}*/
	return nil
}
