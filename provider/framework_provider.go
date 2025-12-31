package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/services/oapi"
	"github.com/outscale/terraform-provider-outscale/internal/services/oks"
)

var (
	_ provider.Provider                       = &FrameworkProvider{}
	_ provider.ProviderWithEphemeralResources = &FrameworkProvider{}
)

type FrameworkProvider struct {
	version     string
	onConfigure func(*client.OutscaleClient)
}

func New(version string) provider.Provider {
	return &FrameworkProvider{
		version: version,
	}
}

func NewWithConfigure(version string, on func(*client.OutscaleClient)) provider.Provider {
	return &FrameworkProvider{
		version:     version,
		onConfigure: on,
	}
}

type ProviderModel struct {
	AccessKeyId  types.String `tfsdk:"access_key_id"`
	SecretKeyId  types.String `tfsdk:"secret_key_id"`
	Region       types.String `tfsdk:"region"`
	Endpoints    []Endpoints  `tfsdk:"endpoints"`
	X509CertPath types.String `tfsdk:"x509_cert_path"`
	X509KeyPath  types.String `tfsdk:"x509_key_path"`
	ConfigFile   types.String `tfsdk:"config_file"`
	Profile      types.String `tfsdk:"profile"`
	Insecure     types.Bool   `tfsdk:"insecure"`
}

type Endpoints struct {
	API types.String `tfsdk:"api"`
	OKS types.String `tfsdk:"oks"`
}

func (p *FrameworkProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "outscale"
	resp.Version = p.version
}

func (p *FrameworkProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Blocks: map[string]schema.Block{
			"endpoints": schema.SetNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"api": schema.StringAttribute{
							Optional:    true,
							Description: "The Endpoint for Outscale API operations.",
						},
						"oks": schema.StringAttribute{
							Optional:    true,
							Description: "The Endpoint for OKS API operations.",
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
			"config_file": schema.StringAttribute{
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

func (p *FrameworkProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
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
	if config.ConfigFile.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("config_file"),
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
	client, err := p.ClientFW(ctx, &config, &diags)
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
	resp.EphemeralResourceData = *client

	if p.onConfigure != nil {
		p.onConfigure(client)
	}
}

func (p *FrameworkProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		oapi.NewDataSourceQuota,
		oks.NewDataSourceOKSQuotas,
		oks.NewDataSourceOKSKubeconfig,
	}
}

func (p *FrameworkProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		oapi.NewResourceNet,
		oapi.NewResourceAccessKey,
		oapi.NewResourcefGPU,
		oapi.NewResourceKeypair,
		oapi.NewResourceSubnet,
		oapi.NewResourceNetPeering,
		oapi.NewResourceNetPeeringAcceptation,
		oapi.NewResourceNetAttributes,
		oapi.NewResourceInternetService,
		oapi.NewResourceInternetServiceLink,
		oapi.NewResourceNetAccessPoint,
		oapi.NewResourceRoute,
		oapi.NewResourceRouteTable,
		oapi.NewResourceRouteTableLink,
		oapi.NewResourceMainRouteTableLink,
		oapi.NewResourceVolume,
		oapi.NewResourceVolumeLink,
		oapi.NewResourceLBUVms,
		oks.NewResourceProject,
		oks.NewResourceCluster,
		oapi.NewResourceSecurityGroup,
		oapi.NewResourceSecurityGroupRule,
	}
}

func (p *FrameworkProvider) EphemeralResources(_ context.Context) []func() ephemeral.EphemeralResource {
	return []func() ephemeral.EphemeralResource{
		oapi.NewKeypairEphemeralResource,
	}
}
