package oapi

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/outscale/goutils/sdk/ptr"
	"github.com/outscale/osc-sdk-go/v3/pkg/options"
	"github.com/outscale/osc-sdk-go/v3/pkg/osc"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/framework/fwhelpers"
	"github.com/outscale/terraform-provider-outscale/internal/framework/fwhelpers/from"
	"github.com/outscale/terraform-provider-outscale/internal/framework/fwhelpers/to"
	"github.com/samber/lo"
)

var (
	_ resource.Resource                = &serverCertificateResource{}
	_ resource.ResourceWithConfigure   = &serverCertificateResource{}
	_ resource.ResourceWithImportState = &serverCertificateResource{}
)

const (
	serverCertErrCreate = "Unable to create Server Certificate"
	serverCertErrRead   = "Unable to read Server Certificate"
	serverCertErrUpdate = "Unable to update Server Certificate"
	serverCertErrDelete = "Unable to delete Server Certificate"
	serverCertErrState  = "Unable to set Cerver Certificate state"
)

type serverCertificateModel struct {
	Id             types.String   `tfsdk:"id"`
	Body           types.String   `tfsdk:"body"`
	Chain          types.String   `tfsdk:"chain"`
	ExpirationDate types.String   `tfsdk:"expiration_date"`
	Name           types.String   `tfsdk:"name"`
	Orn            types.String   `tfsdk:"orn"`
	Path           types.String   `tfsdk:"path"`
	PrivateKey     types.String   `tfsdk:"private_key"`
	RequestId      types.String   `tfsdk:"request_id"`
	UploadDate     types.String   `tfsdk:"upload_date"`
	Timeouts       timeouts.Value `tfsdk:"timeouts"`
}

type serverCertificateResource struct {
	Client *osc.Client
}

func NewResourceServerCertificate() resource.Resource {
	return &serverCertificateResource{}
}

func (r *serverCertificateResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *serverCertificateResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_server_certificate"
}

func (r *serverCertificateResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)

	certId := req.ID
	if certId == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import server certificate identifier. Got: %v", req.ID),
		)
		return
	}

	var data serverCertificateModel
	var timeoutsVal timeouts.Value
	data.Id = to.String(certId)

	resp.Diagnostics.Append(resp.State.GetAttribute(ctx, path.Root("timeouts"), &timeoutsVal)...)
	if resp.Diagnostics.HasError() {
		return
	}
	data.Timeouts = timeoutsVal

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *serverCertificateResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Blocks: map[string]schema.Block{
			"timeouts": timeouts.Block(ctx, timeouts.Opts{
				Create: true,
				Read:   true,
				Update: true,
				Delete: true,
			}),
		},
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"body": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"chain": schema.StringAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"expiration_date": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"orn": schema.StringAttribute{
				Computed: true,
			},
			"path": schema.StringAttribute{
				Computed: true,
				Optional: true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"private_key": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"request_id": schema.StringAttribute{
				Computed: true,
			},
			"upload_date": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *serverCertificateResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data serverCertificateModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := data.Timeouts.Create(ctx, CreateDefaultTimeout)
	if fwhelpers.CheckDiags(resp, diags) {
		return
	}

	createReq := osc.CreateServerCertificateRequest{
		Body:       data.Body.ValueString(),
		Name:       data.Name.ValueString(),
		PrivateKey: data.PrivateKey.ValueString(),
	}

	if fwhelpers.IsSet(data.Chain) {
		createReq.Chain = data.Chain.ValueStringPointer()
	}
	if fwhelpers.IsSet(data.Path) {
		createReq.Path = data.Path.ValueStringPointer()
	}

	createResp, err := r.Client.CreateServerCertificate(ctx, createReq, options.WithRetryTimeout(timeout))
	if err != nil {
		resp.Diagnostics.AddError(serverCertErrCreate, err.Error())
		return
	}

	data.Id = to.String(createResp.ServerCertificate.Id)

	stateData, err := r.read(ctx, timeout, data)
	if err != nil {
		resp.Diagnostics.AddError(serverCertErrState, err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &stateData)...)
}

func (r *serverCertificateResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data serverCertificateModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	timeout, diags := data.Timeouts.Read(ctx, ReadDefaultTimeout)
	if fwhelpers.CheckDiags(resp, diags) {
		return
	}

	stateData, err := r.read(ctx, timeout, data)
	if err != nil {
		if errors.Is(err, ErrResourceEmpty) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(serverCertErrRead, err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &stateData)...)
}

func (r *serverCertificateResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var stateData, planData serverCertificateModel
	resp.Diagnostics.Append(req.State.Get(ctx, &stateData)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &planData)...)
	if resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := planData.Timeouts.Update(ctx, UpdateDefaultTimeout)
	if fwhelpers.CheckDiags(resp, diags) {
		return
	}

	updateReq := osc.UpdateServerCertificateRequest{
		Name: stateData.Name.ValueString(),
	}

	if fwhelpers.HasChange(planData.Name, stateData.Name) {
		updateReq.NewName = planData.Name.ValueStringPointer()
	}
	if fwhelpers.HasChange(planData.Path, stateData.Path) {
		updateReq.NewPath = planData.Path.ValueStringPointer()
	}

	_, err := r.Client.UpdateServerCertificate(ctx, updateReq, options.WithRetryTimeout(timeout))
	if err != nil {
		resp.Diagnostics.AddError(serverCertErrUpdate, err.Error())
		return
	}

	newData, err := r.read(ctx, timeout, planData)
	if err != nil {
		resp.Diagnostics.AddError(serverCertErrState, err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &newData)...)
}

func (r *serverCertificateResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data serverCertificateModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := data.Timeouts.Delete(ctx, DeleteDefaultTimeout)
	if fwhelpers.CheckDiags(resp, diags) {
		return
	}

	deleteReq := osc.DeleteServerCertificateRequest{
		Name: data.Name.ValueString(),
	}

	_, err := r.Client.DeleteServerCertificate(ctx, deleteReq, options.WithRetryTimeout(timeout))
	if err != nil {
		resp.Diagnostics.AddError(serverCertErrDelete, err.Error())
	}
}

func (r *serverCertificateResource) read(ctx context.Context, timeout time.Duration, data serverCertificateModel) (serverCertificateModel, error) {
	resp, err := r.Client.ReadServerCertificates(ctx, osc.ReadServerCertificatesRequest{}, options.WithRetryTimeout(timeout))
	if err != nil {
		return data, err
	}
	if resp.ServerCertificates == nil || len(*resp.ServerCertificates) == 0 {
		return data, ErrResourceEmpty
	}

	server, ok := lo.Find(*resp.ServerCertificates, func(s osc.ServerCertificate) bool {
		return ptr.From(s.Id) == data.Id.ValueString()
	})
	if !ok {
		return data, ErrResourceEmpty
	}

	data.ExpirationDate = to.String(from.ISO8601(server.ExpirationDate))
	data.Name = to.String(server.Name)
	data.Orn = to.String(server.Orn)
	data.Path = to.String(server.Path)
	data.UploadDate = to.String(from.ISO8601(server.UploadDate))
	data.RequestId = to.String(resp.ResponseContext.RequestId)

	return data, nil
}
