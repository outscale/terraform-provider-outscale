package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/boolvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/framework/fwhelpers"
	"github.com/outscale/terraform-provider-outscale/internal/framework/fwhelpers/to"
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
	AccessKey  types.String `tfsdk:"access_key_id"`
	SecretKey  types.String `tfsdk:"secret_key_id"`
	Region     types.String `tfsdk:"region"`
	API        types.List   `tfsdk:"api"`
	OKS        types.List   `tfsdk:"oks"`
	ConfigFile types.String `tfsdk:"config_file"`
	Profile    types.String `tfsdk:"profile"`

	// Deprecated
	X509KeyPath  types.String `tfsdk:"x509_key_path"`
	X509CertPath types.String `tfsdk:"x509_cert_path"`
	Endpoints    []Endpoints  `tfsdk:"endpoints"`
	Insecure     types.Bool   `tfsdk:"insecure"`
}

type APIModel struct {
	X509CertPath types.String `tfsdk:"x509_cert_path"`
	X509KeyPath  types.String `tfsdk:"x509_key_path"`
	Insecure     types.Bool   `tfsdk:"insecure"`
	Endpoint     types.String `tfsdk:"endpoint"`
	Region       types.String `tfsdk:"region"`
}

type OKSModel struct {
	Endpoint types.String `tfsdk:"endpoint"`
	Region   types.String `tfsdk:"region"`
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
				Validators: []validator.Set{
					setvalidator.ConflictsWith(path.MatchRoot("api"), path.MatchRoot("oks")),
				},
				DeprecationMessage: deprecatedMsg("endpoints"),
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
			"api": schema.ListNestedBlock{
				Validators: []validator.List{
					listvalidator.SizeAtMost(1),
				},
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"endpoint": schema.StringAttribute{
							Optional: true,
						},
						"region": schema.StringAttribute{
							Optional: true,
						},
						"x509_cert_path": schema.StringAttribute{
							Optional:    true,
							Description: "Path to the x509 certificate",
						},
						"x509_key_path": schema.StringAttribute{
							Optional:    true,
							Description: "Path to the x509 key",
						},
						"insecure": schema.BoolAttribute{
							Optional:    true,
							Description: "TLS insecure connection",
						},
					},
				},
			},
			"oks": schema.ListNestedBlock{
				Validators: []validator.List{
					listvalidator.SizeAtMost(1),
				},
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"endpoint": schema.StringAttribute{
							Optional: true,
						},
						"region": schema.StringAttribute{
							Optional: true,
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
				Optional:           true,
				DeprecationMessage: deprecatedMsg("region"),
				Description:        "The Region for API operations.",
			},
			"config_file": schema.StringAttribute{
				Optional:    true,
				Description: "Path to the configuration file in which you have defined your credentials.",
			},
			"profile": schema.StringAttribute{
				Optional:    true,
				Description: "Name of your profile in which you define your credencial",
			},
			// Deprecated attributes
			"x509_cert_path": schema.StringAttribute{
				Optional:           true,
				DeprecationMessage: deprecatedMsg("x509_cert_path"),
				Description:        "Path to the x509 certificate for IaaS API operations.",
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.MatchRoot("api"), path.MatchRoot("oks")),
				},
			},
			"x509_key_path": schema.StringAttribute{
				Optional:           true,
				DeprecationMessage: deprecatedMsg("x509_key_path"),
				Description:        "Path to the x509 key for IaaS API operations.",
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.MatchRoot("api"), path.MatchRoot("oks")),
				},
			},
			"insecure": schema.BoolAttribute{
				Optional:           true,
				DeprecationMessage: deprecatedMsg("insecure"),
				Description:        "TLS insecure connection for IaaS API operations.",
				Validators: []validator.Bool{
					boolvalidator.ConflictsWith(path.MatchRoot("api"), path.MatchRoot("oks")),
				},
			},
		},
	}
}

func (p *FrameworkProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config ProviderModel
	diag := req.Config.Get(ctx, &config)
	if fwhelpers.CheckDiags(resp, diag) {
		return
	}

	client, err := config.newClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to create Client",
			"If the error is not clear, please contact the provider developers.\n\n"+
				err.Error(),
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
		oapi.NewResourceSecurityGroup,
		oapi.NewResourceSecurityGroupRule,
		oapi.NewResourcePolicy,
		oapi.NewResourcePolicyVersion,
		oapi.NewResourceUser,
		oapi.NewResourceUserGroup,
		oapi.NewResourceCa,
		oapi.NewResourceApiAccessRule,
		oapi.NewResourceApiAccessPolicy,

		oks.NewResourceProject,
		oks.NewResourceCluster,
	}
}

func (p *FrameworkProvider) EphemeralResources(_ context.Context) []func() ephemeral.EphemeralResource {
	return []func() ephemeral.EphemeralResource{
		oapi.NewKeypairEphemeralResource,
	}
}

func (data *ProviderModel) newClient(ctx context.Context) (*client.OutscaleClient, error) {
	oscConfig, diag := data.buildOSCConfig(ctx)
	if diag.HasError() {
		return nil, fmt.Errorf("failed to build osc config: %v", diag.Errors())
	}
	oksConfig, diag := data.buildOKSConfig(ctx)
	if diag.HasError() {
		return nil, fmt.Errorf("failed to build oks config: %v", diag.Errors())
	}

	if fwhelpers.IsSet(data.AccessKey) {
		oscConfig.AccessKey = data.AccessKey.ValueString()
		oksConfig.AccessKey = data.AccessKey.ValueString()
	}
	if fwhelpers.IsSet(data.SecretKey) {
		oscConfig.SecretKey = data.SecretKey.ValueString()
		oksConfig.SecretKey = data.SecretKey.ValueString()
	}
	oscConfig.UserAgent = UserAgent
	oksConfig.UserAgent = UserAgent

	osc, err := client.NewOSCClient(oscConfig)
	if err != nil {
		return nil, err
	}
	oks, err := client.NewOKSClient(oksConfig)
	if err != nil {
		return nil, err
	}

	client := &client.OutscaleClient{
		OKS: oks,
		OSC: osc,
	}

	return client, nil
}

func (data *ProviderModel) buildOSCConfig(ctx context.Context) (config client.Config, diags diag.Diagnostics) {
	if fwhelpers.IsSet(data.API) {
		apiModel, diag := to.Slice[APIModel](ctx, data.API)
		diags.Append(diag...)
		if diags.HasError() {
			return
		}

		if len(apiModel) > 0 {
			config.APIEndpoint = apiModel[0].Endpoint.ValueString()
			config.Region = apiModel[0].Region.ValueString()
			config.X509CertPath = apiModel[0].X509CertPath.ValueString()
			config.X509KeyPath = apiModel[0].X509KeyPath.ValueString()
			config.Insecure = apiModel[0].Insecure.ValueBool()
		}
	}
	// fallback to deprecated configuration
	if config.APIEndpoint == "" && len(data.Endpoints) > 0 && fwhelpers.IsSet(data.Endpoints[0].API) {
		config.APIEndpoint = data.Endpoints[0].API.ValueString()
	}
	if config.X509CertPath == "" && fwhelpers.IsSet(data.X509CertPath) {
		config.X509CertPath = data.X509CertPath.ValueString()
	}
	if config.X509KeyPath == "" && fwhelpers.IsSet(data.X509KeyPath) {
		config.X509KeyPath = data.X509KeyPath.ValueString()
	}
	if fwhelpers.IsSet(data.Insecure) {
		config.Insecure = data.Insecure.ValueBool()
	}
	if config.Region == "" && fwhelpers.IsSet(data.Region) {
		config.Region = data.Region.ValueString()
	}

	return
}

func (data *ProviderModel) buildOKSConfig(ctx context.Context) (config client.Config, diags diag.Diagnostics) {
	if fwhelpers.IsSet(data.OKS) {
		oksModel, diag := to.Slice[OKSModel](ctx, data.OKS)
		diags.Append(diag...)
		if diags.HasError() {
			return
		}

		if len(oksModel) > 0 {
			config.OKSEndpoint = oksModel[0].Endpoint.ValueString()
			config.Region = oksModel[0].Region.ValueString()
		}
	}
	// fallback to deprecated configuration
	if config.OKSEndpoint == "" && len(data.Endpoints) > 0 && fwhelpers.IsSet(data.Endpoints[0].OKS) {
		config.OKSEndpoint = data.Endpoints[0].OKS.ValueString()
	}
	if config.Region == "" && fwhelpers.IsSet(data.Region) {
		config.Region = data.Region.ValueString()
	}

	return
}
