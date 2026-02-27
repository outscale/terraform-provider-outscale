package oapi

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/outscale/goutils/sdk/ptr"
	"github.com/outscale/osc-sdk-go/v3/pkg/options"
	"github.com/outscale/osc-sdk-go/v3/pkg/osc"

	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/framework/fwhelpers"
	"github.com/outscale/terraform-provider-outscale/internal/framework/fwhelpers/from"
	"github.com/outscale/terraform-provider-outscale/internal/framework/fwhelpers/to"
)

var (
	_ resource.Resource              = &resourcePolicyVersion{}
	_ resource.ResourceWithConfigure = &resourcePolicyVersion{}
)

type PolicyVersionModel struct {
	Body           types.String   `tfsdk:"body"`
	CreationDate   types.String   `tfsdk:"creation_date"`
	DefaultVersion types.Bool     `tfsdk:"default_version"`
	VersionId      types.String   `tfsdk:"version_id"`
	PolicyOrn      types.String   `tfsdk:"policy_orn"`
	Document       types.String   `tfsdk:"document"`
	SetAsDefault   types.Bool     `tfsdk:"set_as_default"`
	Timeouts       timeouts.Value `tfsdk:"timeouts"`
	Id             types.String   `tfsdk:"id"`
}

type resourcePolicyVersion struct {
	Client *osc.Client
}

func NewResourcePolicyVersion() resource.Resource {
	return &resourcePolicyVersion{}
}

func (r *resourcePolicyVersion) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	client, ok := req.ProviderData.(client.OutscaleClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *osc.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}
	r.Client = client.OSC
}

func (r *resourcePolicyVersion) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_policy_version"
}

func (r *resourcePolicyVersion) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Blocks: map[string]schema.Block{
			"timeouts": timeouts.Block(ctx, timeouts.Opts{
				Create: true,
				Read:   true,
				Delete: true,
			}),
		},
		Attributes: map[string]schema.Attribute{
			"body": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"creation_date": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"default_version": schema.BoolAttribute{
				Computed: true,
			},
			"version_id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"policy_orn": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"document": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"set_as_default": schema.BoolAttribute{
				Optional: true,
			},
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *resourcePolicyVersion) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data PolicyVersionModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	createTimeout, diag := data.Timeouts.Create(ctx, CreateDefaultTimeout)
	if fwhelpers.CheckDiags(resp, diag) {
		return
	}

	createReq := osc.CreatePolicyVersionRequest{
		Document:  data.Document.ValueString(),
		PolicyOrn: data.PolicyOrn.ValueString(),
	}

	if fwhelpers.IsSet(data.SetAsDefault) {
		createReq.SetAsDefault = data.SetAsDefault.ValueBoolPointer()
	}

	createResp, err := r.Client.CreatePolicyVersion(ctx, createReq, options.WithRetryTimeout(createTimeout))
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to create Policy Version",
			err.Error(),
		)
		return
	}
	policyVersion := ptr.From(createResp.PolicyVersion)
	data.Id = to.String(id.UniqueId())
	data.VersionId = to.String(policyVersion.VersionId)

	stateData, err := r.read(ctx, createTimeout, data)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to set Policy Version state",
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &stateData)...)
}

func (r *resourcePolicyVersion) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data PolicyVersionModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	to, diag := data.Timeouts.Read(ctx, ReadDefaultTimeout)
	if fwhelpers.CheckDiags(resp, diag) {
		return
	}

	stateData, err := r.read(ctx, to, data)
	if err != nil {
		if errors.Is(err, ErrResourceEmpty) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Unable to set Policy Version state",
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &stateData)...)
}

func (r *resourcePolicyVersion) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

func (r *resourcePolicyVersion) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data PolicyVersionModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	deleteTimeout, diag := data.Timeouts.Delete(ctx, DeleteDefaultTimeout)
	if fwhelpers.CheckDiags(resp, diag) {
		return
	}

	// Refresh state to check current default version status
	// This is necessary because user/user_group resources can change the default version
	stateData, err := r.read(ctx, deleteTimeout, data)
	if err != nil {
		if errors.Is(err, ErrResourceEmpty) {
			return
		}
		resp.Diagnostics.AddError(
			"Unable to refresh Policy Version state before deletion",
			err.Error(),
		)
		return
	}
	data = stateData

	// If this version is currently the default, we need to set v1 as default first
	if data.DefaultVersion.ValueBool() {
		req := osc.SetDefaultPolicyVersionRequest{
			PolicyOrn: data.PolicyOrn.ValueString(),
			VersionId: "v1",
		}
		_, err := r.Client.SetDefaultPolicyVersion(ctx, req, options.WithRetryTimeout(deleteTimeout))
		if err != nil {
			resp.Diagnostics.AddError(
				"Unable to set Policy Version v1 as default before deletion",
				err.Error(),
			)
			return
		}
	}

	deleteReq := osc.DeletePolicyVersionRequest{
		PolicyOrn: data.PolicyOrn.ValueString(),
		VersionId: data.VersionId.ValueString(),
	}
	_, err = r.Client.DeletePolicyVersion(ctx, deleteReq, options.WithRetryTimeout(deleteTimeout))
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to delete Policy Version",
			err.Error(),
		)
	}
}

func (r *resourcePolicyVersion) read(ctx context.Context, timeout time.Duration, data PolicyVersionModel) (PolicyVersionModel, error) {
	req := osc.ReadPolicyVersionRequest{
		PolicyOrn: data.PolicyOrn.ValueString(),
		VersionId: data.VersionId.ValueString(),
	}
	resp, err := r.Client.ReadPolicyVersion(ctx, req, options.WithRetryTimeout(timeout))
	if err != nil {
		return data, err
	}

	if resp.PolicyVersion == nil {
		return data, ErrResourceEmpty
	}

	policyVersion := ptr.From(resp.PolicyVersion)

	data.DefaultVersion = to.Bool(ptr.From(policyVersion.DefaultVersion))
	data.CreationDate = to.String(from.ISO8601(policyVersion.CreationDate))
	data.Body = to.String(ptr.From(policyVersion.Body))

	return data, nil
}
