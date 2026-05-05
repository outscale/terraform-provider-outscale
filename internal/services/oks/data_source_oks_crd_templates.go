package oks

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/outscale/osc-sdk-go/v3/pkg/oks"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/framework/fwhelpers"
	"github.com/outscale/terraform-provider-outscale/internal/framework/fwhelpers/to"
	"github.com/samber/lo"
	"sigs.k8s.io/yaml"
)

var (
	_ datasource.DataSource              = &crdTemplatesDataSource{}
	_ datasource.DataSourceWithConfigure = &crdTemplatesDataSource{}
)

type crdTemplatesDataSource struct {
	Client *oks.Client
}

type crdTemplatesModel struct {
	Manifests types.Set    `tfsdk:"manifests"`
	Id        types.String `tfsdk:"id"`
}

func NewDataSourceCRDTemplates() datasource.DataSource {
	return &crdTemplatesDataSource{}
}

func (d *crdTemplatesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *crdTemplatesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_oks_crd_templates"
}

func (d *crdTemplatesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"manifests": schema.SetAttribute{
				Computed:    true,
				ElementType: types.StringType,
			},
			"id": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

func (d *crdTemplatesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data crdTemplatesModel
	if fwhelpers.CheckDiags(resp, req.Config.Get(ctx, &data)) {
		return
	}

	nodepool, err := d.Client.GetNodepoolTemplate(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Unable to get NodePool template", err.Error())
		return
	}

	netPeering, err := d.Client.GetNetPeeringRequestTemplate(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Unable to get NetPeeringRequest template", err.Error())
		return
	}

	netPeeringAcceptance, err := d.Client.GetNetPeeringAcceptanceTemplate(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Unable to get NetPeeringAcceptance template", err.Error())
		return
	}

	manifests, err := lo.MapErr([]any{nodepool.Template, netPeering.Template, netPeeringAcceptance.Template}, func(item any, _ int) (string, error) {
		yamlManifest, err := yaml.Marshal(item)
		if err != nil {
			return "", err
		}

		return strings.TrimSpace(string(yamlManifest)) + "\n", nil
	})
	if err != nil {
		resp.Diagnostics.AddError("Unable to convert manifest to YAML", err.Error())
		return
	}

	manifestsSet, diag := to.Set(ctx, manifests)
	if fwhelpers.CheckDiags(resp, diag) {
		return
	}

	data.Manifests = manifestsSet
	data.Id = to.String("oks-crd-templates")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
