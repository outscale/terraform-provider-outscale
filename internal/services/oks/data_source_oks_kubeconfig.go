package oks

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/outscale/osc-sdk-go/v3/pkg/oks"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/framework/fwhelpers"
	"github.com/outscale/terraform-provider-outscale/internal/framework/fwhelpers/to"
)

var (
	_ datasource.DataSource              = &oksKubeconfigDataSource{}
	_ datasource.DataSourceWithConfigure = &oksKubeconfigDataSource{}
)

func NewDataSourceOKSKubeconfig() datasource.DataSource {
	return &oksKubeconfigDataSource{}
}

func (d *oksKubeconfigDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	client, ok := req.ProviderData.(client.OutscaleClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected DataSource Configure Type",
			fmt.Sprintf("Expected *oks.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}
	d.Client = client.OKS
}

type oksKubeconfigDataSource struct {
	Client *oks.Client
}

type oksKubeconfigModel struct {
	ClusterId    types.String `tfsdk:"cluster_id"`
	User         types.String `tfsdk:"user"`
	Group        types.String `tfsdk:"group"`
	Ttl          types.String `tfsdk:"ttl"`
	XEncryptNacl types.String `tfsdk:"x_encrypt_nacl"`
	Kubeconfig   types.String `tfsdk:"kubeconfig"`
	RequestId    types.String `tfsdk:"request_id"`
	Id           types.String `tfsdk:"id"`
}

func (d *oksKubeconfigDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_oks_kubeconfig"
}

func (d *oksKubeconfigDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"cluster_id": schema.StringAttribute{
				Required: true,
			},
			"user": schema.StringAttribute{
				Optional: true,
			},
			"group": schema.StringAttribute{
				Optional: true,
			},
			"ttl": schema.StringAttribute{
				Optional: true,
			},
			"x_encrypt_nacl": schema.StringAttribute{
				Optional:  true,
				Sensitive: true,
			},
			"kubeconfig": schema.StringAttribute{
				Computed:  true,
				Sensitive: true,
			},
			"request_id": schema.StringAttribute{
				Computed: true,
			},
			"id": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

func (d *oksKubeconfigDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var (
		data, config     oksKubeconfigModel
		kubeconfigResp   *oks.KubeconfigResponse
		err              error
		user, group, ttl *string
	)
	diags := req.Config.Get(ctx, &config)
	if fwhelpers.CheckDiags(resp, diags) {
		return
	}

	if fwhelpers.IsSet(config.User) {
		user = config.User.ValueStringPointer()
	}
	if fwhelpers.IsSet(config.Group) {
		group = config.Group.ValueStringPointer()
	}
	if fwhelpers.IsSet(config.Ttl) {
		ttl = config.Ttl.ValueStringPointer()
	}
	clusterId := config.ClusterId.ValueString()

	if fwhelpers.IsSet(config.XEncryptNacl) {
		params := &oks.GetKubeconfigWithPubkeyNACLParams{
			User:         user,
			Group:        group,
			Ttl:          ttl,
			XEncryptNacl: config.XEncryptNacl.ValueStringPointer(),
		}
		kubeconfigResp, err = d.Client.GetKubeconfigWithPubkeyNACL(ctx, clusterId, params)
	} else {
		params := &oks.GetKubeconfigParams{
			User:  user,
			Group: group,
			Ttl:   ttl,
		}
		kubeconfigResp, err = d.Client.GetKubeconfig(ctx, clusterId, params)
	}

	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to get OKS Kubeconfig",
			"Error: "+err.Error(),
		)
		return
	}

	data.Kubeconfig = to.String(kubeconfigResp.Cluster.Data.Kubeconfig)
	data.RequestId = to.String(kubeconfigResp.Cluster.RequestId)
	data.Id = to.String(id.UniqueId())

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}
