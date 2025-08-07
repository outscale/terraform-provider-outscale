package outscale

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	sdkv3_oks "github.com/outscale/osc-sdk-go/v3/pkg/oks"
)

var (
	_ datasource.DataSource              = &oksQuotasDataSource{}
	_ datasource.DataSourceWithConfigure = &oksQuotasDataSource{}
)

func NewDataSourceOKSQuotas() datasource.DataSource {
	return &oksQuotasDataSource{}
}

func (d *oksQuotasDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	client, ok := req.ProviderData.(OutscaleClientFW)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected DataSource Configure Type",
			fmt.Sprintf("Expected *oks.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}
	d.Client = client.OKS
}

type oksQuotasDataSource struct {
	Client *sdkv3_oks.Client
}

type oksQuotasModel struct {
	CPSubregions       types.Set    `tfsdk:"cp_subregions"`
	ClustersPerProject types.Int32  `tfsdk:"clusters_per_project"`
	KubeVersions       types.Set    `tfsdk:"kube_versions"`
	Projects           types.Int32  `tfsdk:"projects"`
	RequestId          types.String `tfsdk:"request_id"`
}

func (d *oksQuotasDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_oks_quotas"
}

func (d *oksQuotasDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"cp_subregions": schema.SetAttribute{
				Computed:    true,
				ElementType: types.StringType,
			},
			"clusters_per_project": schema.Int64Attribute{
				Computed: true,
			},
			"kube_versions": schema.SetAttribute{
				Computed:    true,
				ElementType: types.StringType,
			},
			"projects": schema.Int64Attribute{
				Computed: true,
			},
			"request_id": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

func (d *oksQuotasDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data oksQuotasModel

	quotas, err := d.Client.GetQuotas(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to get OKS Quotas",
			"Error: "+err.Error(),
		)
		return
	}

	cpSubregions, diags := types.SetValueFrom(ctx, types.StringType, quotas.Quotas.CPSubregions)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	data.CPSubregions = cpSubregions
	data.ClustersPerProject = types.Int32Value(int32(quotas.Quotas.ClustersPerProject))
	kubeVer, diags := types.SetValueFrom(ctx, types.StringType, quotas.Quotas.KubeVersions)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	data.KubeVersions = kubeVer
	data.Projects = types.Int32Value(int32(quotas.Quotas.Projects))
	data.RequestId = types.StringValue(quotas.ResponseContext.RequestId)
	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
