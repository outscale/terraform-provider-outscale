package oos

import (
	"context"
	"fmt"
	"time"

	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/outscale/osc-sdk-go/v3/pkg/oos"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/framework/fwhelpers"
	"github.com/outscale/terraform-provider-outscale/internal/framework/fwhelpers/to"
)

var (
	_ resource.Resource                   = &presignedURLResource{}
	_ resource.ResourceWithConfigure      = &presignedURLResource{}
	_ resource.ResourceWithValidateConfig = &presignedURLResource{}
	_ resource.ResourceWithModifyPlan     = &presignedURLResource{}
)

const (
	presignedURLErrCreate = "Unable to create Pre-Signed URL"
)

type presignedURLResource struct {
	Client *oos.Client
}

type presignedURLModel struct {
	Bucket     types.String         `tfsdk:"bucket"`
	Key        types.String         `tfsdk:"key"`
	Method     types.String         `tfsdk:"method"`
	Expiration timetypes.GoDuration `tfsdk:"expiration"`
	ExpiresAt  timetypes.RFC3339    `tfsdk:"expires_at"`
	URL        types.String         `tfsdk:"url"`
	Id         types.String         `tfsdk:"id"`
}

func NewResourcePresignedURL() resource.Resource {
	return &presignedURLResource{}
}

func (r *presignedURLResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(client.OutscaleClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *oos.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.Client = client.OOS
}

func (r *presignedURLResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_oos_presigned_url"
}

func (r *presignedURLResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var data presignedURLModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	method := data.Method.ValueString()
	if (method == "GET" || method == "PUT") && data.Key.IsNull() {
		resp.Diagnostics.AddAttributeError(
			path.Root("key"),
			"Missing Object Key",
			fmt.Sprintf("A non-empty key must be configured when method is %s.", method),
		)
	}
}

func (r *presignedURLResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	if req.State.Raw.IsNull() || req.Plan.Raw.IsNull() {
		return
	}

	var state presignedURLModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() || !fwhelpers.IsSet(state.ExpiresAt) {
		return
	}

	expiresAt, diags := state.ExpiresAt.ValueRFC3339Time()
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !time.Now().Before(expiresAt) {
		// We need to set the value to unknown for the `RequiresReplace` to be triggered
		resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, path.Root("expires_at"), timetypes.NewRFC3339Unknown())...)
		if resp.Diagnostics.HasError() {
			return
		}

		resp.RequiresReplace = append(resp.RequiresReplace, path.Root("expires_at"))
		resp.Diagnostics.AddAttributeWarning(
			path.Root("expires_at"),
			"Pre-Signed URL Expired",
			fmt.Sprintf(
				"The pre-signed URL expired at %s and will be replaced.\n\n(This warning is expected to appear twice: during `terraform plan` and again after a successful `terraform apply`. The newly generated pre-signed URL is valid once apply completes.)",
				expiresAt.Format(time.RFC3339),
			),
		)
	}
}

func (r *presignedURLResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"bucket": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"key": schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"method": schema.StringAttribute{
				Computed: true,
				Optional: true,
				Default:  stringdefault.StaticString("GET"),
				Validators: []validator.String{
					stringvalidator.OneOf("GET", "PUT", "HEAD", "DELETE"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"expiration": schema.StringAttribute{
				Computed:   true,
				Optional:   true,
				CustomType: timetypes.GoDurationType{},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"url": schema.StringAttribute{
				Computed: true,
			},
			"expires_at": schema.StringAttribute{
				Computed:   true,
				CustomType: timetypes.RFC3339Type{},
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

func (r *presignedURLResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data presignedURLModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Default Pre-Signed URL expiration of the client is 900 seconds (15 minutes)
	expiration := 15 * time.Minute
	switch {
	case fwhelpers.IsSet(data.Expiration):
		expirationDuration, d := data.Expiration.ValueGoDuration()
		if fwhelpers.CheckDiags(resp, d) {
			return
		}

		expiration = expirationDuration
	default:
		data.Expiration = timetypes.NewGoDurationValue(expiration)
	}

	presign := oos.NewPresignClient(r.Client, func(opts *s3.PresignOptions) {
		opts.Expires = expiration
	})

	var signedReq *v4.PresignedHTTPRequest
	var err error
	switch data.Method.ValueString() {
	case "GET":
		signedReq, err = presign.PresignGetObject(ctx, &s3.GetObjectInput{
			Bucket: data.Bucket.ValueStringPointer(),
			Key:    data.Key.ValueStringPointer(),
		})
	case "PUT":
		signedReq, err = presign.PresignPutObject(ctx, &s3.PutObjectInput{
			Bucket: data.Bucket.ValueStringPointer(),
			Key:    data.Key.ValueStringPointer(),
		})
	case "HEAD":
		if fwhelpers.IsSet(data.Key) {
			signedReq, err = presign.PresignHeadObject(ctx, &s3.HeadObjectInput{
				Bucket: data.Bucket.ValueStringPointer(),
				Key:    data.Key.ValueStringPointer(),
			})
		} else {
			signedReq, err = presign.PresignHeadBucket(ctx, &s3.HeadBucketInput{
				Bucket: data.Bucket.ValueStringPointer(),
			})
		}
	case "DELETE":
		if fwhelpers.IsSet(data.Key) {
			signedReq, err = presign.PresignDeleteObject(ctx, &s3.DeleteObjectInput{
				Bucket: data.Bucket.ValueStringPointer(),
				Key:    data.Key.ValueStringPointer(),
			})
		} else {
			signedReq, err = presign.PresignDeleteBucket(ctx, &s3.DeleteBucketInput{
				Bucket: data.Bucket.ValueStringPointer(),
			})
		}
	}
	if err != nil {
		resp.Diagnostics.AddError(presignedURLErrCreate, err.Error())
		return
	}
	if signedReq == nil {
		resp.Diagnostics.AddError(presignedURLErrCreate, "The pre-signer did not return a request.")
		return
	}

	data.ExpiresAt = timetypes.NewRFC3339TimeValue(time.Now().UTC().Add(expiration))
	data.URL = to.String(signedReq.URL)
	data.Id = types.StringValue(id.UniqueId())

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *presignedURLResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data presignedURLModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *presignedURLResource) Update(_ context.Context, _ resource.UpdateRequest, resp *resource.UpdateResponse) {
}

func (r *presignedURLResource) Delete(_ context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
}
