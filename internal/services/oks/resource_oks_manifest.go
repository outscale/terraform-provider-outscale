package oks

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/mapvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/outscale/osc-sdk-go/v3/pkg/oks"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/framework/fwhelpers"
	"github.com/outscale/terraform-provider-outscale/internal/framework/fwhelpers/to"
	"github.com/outscale/terraform-provider-outscale/internal/framework/validators/validatorstring"
	"github.com/outscale/terraform-provider-outscale/internal/services/oks/okshelpers"

	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	k8stypes "k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/util/jsonpath"
)

var (
	_ resource.Resource               = &manifestResource{}
	_ resource.ResourceWithConfigure  = &manifestResource{}
	_ resource.ResourceWithModifyPlan = &manifestResource{}
)

const (
	manifestErrCreate    = "Unable to create manifest resource"
	manifestErrRead      = "Unable to read manifest resource"
	manifestErrUpdate    = "Unable to update manifest resource"
	manifestErrDelete    = "Unable to delete manifest resource"
	manifestErrWait      = "Unable to wait for resource deletion"
	manifestErrWaitFor   = "Unable to wait for resource wait_for condition to be met"
	manifestFieldManager = "terraform-provider-outscale"

	manifestDeleteBlockedWarning = `Terraform is planning to delete this Kubernetes object. The Kubernetes API may reject the deletion if the object is still required by the cluster, for example when deleting the last remaining NodePool.

If deletion fails and the Kubernetes object is expected to disappear with its cluster, such as when destroying the whole configuration or deleting the OKS Cluster and its associated manifests together, set 'skip_delete = true'. After setting 'skip_delete = true', run 'terraform apply' so the value is saved in state before retrying the destroy.

With 'skip_delete' enabled, a delete removes only the Terraform state entry, it does not delete the Kubernetes object. Removing this resource from state with 'terraform state rm' has the same state-only effect.`

	manifestReplaceWarning = `Terraform planned this manifest change as a replacement. By default, Terraform deletes the current Kubernetes object before creating the replacement.

If the current object is still required by the cluster, configure Terraform to create the replacement before deleting the current object:

lifecycle {
  create_before_destroy = true
}

The current and replacement objects must be able to coexist during the replacement. For resources that require unique names, use a different metadata.name in the replacement manifest.`
)

type manifestResource struct {
	Client *oks.Client
}

type manifestModel struct {
	ClusterId  types.String   `tfsdk:"cluster_id"`
	Manifest   types.String   `tfsdk:"manifest"`
	Object     types.String   `tfsdk:"object"`
	SkipDelete types.Bool     `tfsdk:"skip_delete"`
	Wait       types.Bool     `tfsdk:"wait"`
	WaitFor    types.Object   `tfsdk:"wait_for"`
	Id         types.String   `tfsdk:"id"`
	Timeouts   timeouts.Value `tfsdk:"timeouts"`
}

type manifestWaitForModel struct {
	Timeout timetypes.GoDuration `tfsdk:"timeout"`
	Fields  types.Map            `tfsdk:"fields"`
}

func NewResourceManifest() resource.Resource {
	return &manifestResource{}
}

func (r *manifestResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(client.OutscaleClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *oks.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.Client = client.OKS
}

func (r *manifestResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_oks_manifest"
}

func (r *manifestResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	if req.Plan.Raw.IsNull() {
		// destroy: show warning about lifecyle
		if !req.State.Raw.IsNull() {
			var state manifestModel
			resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
			if resp.Diagnostics.HasError() {
				return
			}

			// do not show the warning is skip_delete is set to true
			if !fwhelpers.IsSet(state.SkipDelete) || !state.SkipDelete.ValueBool() {
				resp.Diagnostics.AddWarning("Manifest resource deletion may be blocked", manifestDeleteBlockedWarning)
			}
		}
		return
	}

	var plan manifestModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// dri is retrievable only if we have a cluster_id and manifest
	if !fwhelpers.IsSet(plan.ClusterId) || !fwhelpers.IsSet(plan.Manifest) {
		return
	}

	timeout, diags := plan.Timeouts.Create(ctx, CreateDefaultTimeout)
	if fwhelpers.CheckDiags(resp, diags) {
		return
	}

	planObj, dri, err := okshelpers.GetResourceInterfaceFromManifest(ctx, r.Client, plan.ClusterId.ValueString(), plan.Manifest.ValueString(), timeout)
	if err != nil {
		resp.Diagnostics.AddError("Invalid manifest", err.Error())
		return
	}

	// update: check if the resource needs to be replaced
	if !req.State.Raw.IsNull() {
		var state manifestModel
		resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
		if resp.Diagnostics.HasError() {
			return
		}
		if !fwhelpers.IsSet(state.Manifest) {
			return
		}

		stateObj, err := okshelpers.FromYAML(state.Manifest.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("Invalid manifest", err.Error())
			return
		}

		if planObj.GetKind() != stateObj.GetKind() || planObj.GetName() != stateObj.GetName() || planObj.GetNamespace() != stateObj.GetNamespace() {
			resp.RequiresReplace = append(resp.RequiresReplace, path.Root("manifest"))
			resp.Diagnostics.AddWarning("Manifest resource replacement may require 'create_before_destroy'", manifestReplaceWarning)
		}
	}

	jsonBytes, err := planObj.MarshalJSON()
	if err != nil {
		resp.Diagnostics.AddError("Invalid manifest", err.Error())
		return
	}

	// validate the manifest using a dry-run patch apply
	patchOptions := metav1.PatchOptions{
		FieldManager: manifestFieldManager,
		DryRun:       []string{metav1.DryRunAll},
	}
	_, err = dri.Patch(ctx, planObj.GetName(), k8stypes.ApplyPatchType, jsonBytes, patchOptions)
	if err != nil {
		resp.Diagnostics.AddError("Dry-run validation failed", err.Error())
		return
	}
}

func (r *manifestResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
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
			"cluster_id": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"manifest": schema.StringAttribute{
				Required: true,
			},
			"object": schema.StringAttribute{
				Computed: true,
			},
			"skip_delete": schema.BoolAttribute{
				Computed: true,
				Optional: true,
				Default:  booldefault.StaticBool(false),
			},
			"wait": schema.BoolAttribute{
				Computed: true,
				Optional: true,
				Default:  booldefault.StaticBool(false),
			},
			"wait_for": schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{
					"timeout": schema.StringAttribute{
						Optional:   true,
						CustomType: timetypes.GoDurationType{},
					},
					"fields": schema.MapAttribute{
						Required:    true,
						ElementType: types.StringType,
						Validators: []validator.Map{
							mapvalidator.KeysAre(
								stringvalidator.LengthAtLeast(1),
								validatorstring.IsJSONPath(toJSONPathExpression),
							),
							mapvalidator.ValueStringsAre(
								validatorstring.IsRegex(),
							),
						},
					},
				},
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

func (r *manifestResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data manifestModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	timeout, diag := data.Timeouts.Create(ctx, CreateDefaultTimeout)
	if fwhelpers.CheckDiags(resp, diag) {
		return
	}

	obj, dri, err := okshelpers.GetResourceInterfaceFromManifest(ctx, r.Client, data.ClusterId.ValueString(), data.Manifest.ValueString(), timeout)
	if err != nil {
		resp.Diagnostics.AddError(manifestErrCreate, err.Error())
		return
	}

	objResp, err := dri.Get(ctx, obj.GetName(), metav1.GetOptions{})
	if err != nil && !k8serrors.IsNotFound(err) {
		resp.Diagnostics.AddError(manifestErrCreate, err.Error())
		return
	}
	if objResp != nil {
		resp.Diagnostics.AddError(manifestErrCreate, "Manifest resource with the same name already exists.")
		return
	}

	err = r.applyManifest(ctx, &data, obj, dri)
	if err != nil {
		resp.Diagnostics.AddError(manifestErrCreate, err.Error())
		return
	}
	diag = r.waitForFields(ctx, data, obj, dri, timeout)
	if fwhelpers.CheckDiags(resp, diag) {
		return
	}

	stateData, err := r.read(ctx, data, obj, dri)
	if err != nil {
		resp.Diagnostics.AddError(manifestErrRead, err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &stateData)...)
}

func (r *manifestResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data manifestModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	timeout, diag := data.Timeouts.Read(ctx, ReadDefaultTimeout)
	if fwhelpers.CheckDiags(resp, diag) {
		return
	}

	obj, dri, err := okshelpers.GetResourceInterfaceFromManifest(ctx, r.Client, data.ClusterId.ValueString(), data.Manifest.ValueString(), timeout)
	if err != nil {
		resp.Diagnostics.AddError(manifestErrRead, err.Error())
		return
	}

	stateData, err := r.read(ctx, data, obj, dri)
	if err != nil {
		if errors.Is(err, ErrResourceEmpty) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(manifestErrRead, err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &stateData)...)
}

func (r *manifestResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state manifestModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	timeout, diag := plan.Timeouts.Update(ctx, UpdateDefaultTimeout)
	if fwhelpers.CheckDiags(resp, diag) {
		return
	}

	obj, dri, err := okshelpers.GetResourceInterfaceFromManifest(ctx, r.Client, plan.ClusterId.ValueString(), plan.Manifest.ValueString(), timeout)
	if err != nil {
		resp.Diagnostics.AddError(manifestErrUpdate, err.Error())
		return
	}

	if fwhelpers.HasChange(plan.Manifest, state.Manifest) {
		err := r.applyManifest(ctx, &plan, obj, dri)
		if err != nil {
			resp.Diagnostics.AddError(manifestErrUpdate, err.Error())
			return
		}

		diag := r.waitForFields(ctx, plan, obj, dri, timeout)
		if fwhelpers.CheckDiags(resp, diag) {
			return
		}
	}

	stateData, err := r.read(ctx, plan, obj, dri)
	if err != nil {
		resp.Diagnostics.AddError(manifestErrRead, err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &stateData)...)
}

func (r *manifestResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data manifestModel
	if fwhelpers.CheckDiags(resp, req.State.Get(ctx, &data)) {
		return
	}

	timeout, diag := data.Timeouts.Delete(ctx, DeleteDefaultTimeout)
	if fwhelpers.CheckDiags(resp, diag) {
		return
	}

	obj, dri, err := okshelpers.GetResourceInterfaceFromManifest(ctx, r.Client, data.ClusterId.ValueString(), data.Manifest.ValueString(), timeout)
	if err != nil {
		resp.Diagnostics.AddError(manifestErrDelete, err.Error())
		return
	}

	err = dri.Delete(ctx, obj.GetName(), metav1.DeleteOptions{})

	resourceGone := k8serrors.IsGone(err) || k8serrors.IsNotFound(err)
	if err != nil && !resourceGone {
		if fwhelpers.IsSet(data.SkipDelete) && data.SkipDelete.ValueBool() {
			return
		}

		resp.Diagnostics.AddError(manifestErrDelete, err.Error())
		return
	}

	if fwhelpers.IsSet(data.Wait) && data.Wait.ValueBool() {
		err = wait.PollUntilContextTimeout(ctx, 5*time.Second, timeout, true, func(ctx context.Context) (bool, error) {
			_, err := dri.Get(ctx, obj.GetName(), metav1.GetOptions{})
			if k8serrors.IsNotFound(err) || k8serrors.IsGone(err) {
				return true, nil
			}

			return false, err
		})
		if err != nil {
			resp.Diagnostics.AddError(manifestErrWait, err.Error())
		}
	}
}

func (r *manifestResource) read(ctx context.Context, data manifestModel, obj *unstructured.Unstructured, dri dynamic.ResourceInterface) (manifestModel, error) {
	respManifest, err := dri.Get(ctx, obj.GetName(), metav1.GetOptions{})
	if k8serrors.IsNotFound(err) || k8serrors.IsGone(err) {
		return data, ErrResourceEmpty
	}
	if err != nil {
		return data, fmt.Errorf("get manifest: %w", err)
	}

	yamlObj, err := okshelpers.ToYAML(respManifest.Object)
	if err != nil {
		return data, err
	}

	data.Id = to.String(respManifest.GetUID())
	data.Object = to.String(yamlObj)

	return data, nil
}

func (r *manifestResource) applyManifest(ctx context.Context, data *manifestModel, obj *unstructured.Unstructured, dri dynamic.ResourceInterface) error {
	jsonBytes, err := obj.MarshalJSON()
	if err != nil {
		return err
	}

	patchOptions := metav1.PatchOptions{FieldManager: manifestFieldManager}
	resp, err := dri.Patch(ctx, obj.GetName(), k8stypes.ApplyPatchType, jsonBytes, patchOptions)
	if err != nil {
		return fmt.Errorf("apply manifest: %w", err)
	}

	data.Id = to.String(resp.GetUID())

	return nil
}

type fieldMatcher struct {
	Field string
	Path  *jsonpath.JSONPath
	Regex *regexp.Regexp
}

func (r *manifestResource) waitForFields(ctx context.Context, data manifestModel, obj *unstructured.Unstructured, dri dynamic.ResourceInterface, timeout time.Duration) (diags diag.Diagnostics) {
	if !fwhelpers.IsSet(data.WaitFor) {
		return nil
	}

	fields, waitForTimeout, diag := expandWaitFor(ctx, data.WaitFor)
	if diag.HasError() {
		return diag
	}
	if waitForTimeout > 0 {
		timeout = waitForTimeout
	}

	matchers := make([]fieldMatcher, 0, len(fields))
	for key, value := range fields {
		re, err := regexp.Compile(value)
		if err != nil {
			diags.AddAttributeError(
				path.Root("wait_for").AtName("fields").AtMapKey(key),
				"Invalid wait_for field pattern",
				err.Error(),
			)
			return
		}

		jp := jsonpath.New(key).AllowMissingKeys(true)
		err = jp.Parse(toJSONPathExpression(key))
		if err != nil {
			diags.AddAttributeError(
				path.Root("wait_for").AtName("fields").AtMapKey(key),
				"Invalid JSONPath expression",
				err.Error(),
			)
			return
		}

		matchers = append(matchers, fieldMatcher{
			Field: key,
			Path:  jp,
			Regex: re,
		})
	}

	err := wait.PollUntilContextTimeout(ctx, 5*time.Second, timeout, true, func(ctx context.Context) (bool, error) {
		resp, err := dri.Get(ctx, obj.GetName(), metav1.GetOptions{})
		if err != nil {
			return false, err
		}

		for _, m := range matchers {
			var output bytes.Buffer
			if err := m.Path.Execute(&output, resp.Object); err != nil {
				return false, fmt.Errorf("evaluate jsonpath %q: %w", m.Field, err)
			}

			value := output.String()
			if value == "" {
				return false, nil
			}

			if !m.Regex.MatchString(value) {
				return false, nil
			}
		}

		return true, nil
	})
	if err != nil {
		diags.AddError(manifestErrWaitFor, err.Error())
		return
	}

	return
}

func expandWaitFor(ctx context.Context, waitFor types.Object) (map[string]string, time.Duration, diag.Diagnostics) {
	model, diags := to.Model[manifestWaitForModel](ctx, waitFor)
	if diags.HasError() {
		return nil, 0, diags
	}

	var timeout time.Duration
	if fwhelpers.IsSet(model.Timeout) {
		parsedTimeout, diags := model.Timeout.ValueGoDuration()
		if diags.HasError() {
			return nil, 0, diags
		}
		timeout = parsedTimeout
	}

	fields := make(map[string]string)
	diags = model.Fields.ElementsAs(ctx, &fields, false)
	if diags.HasError() {
		return nil, 0, diags
	}

	return fields, timeout, nil
}

// Converts a field string to a JSONPath expression
// We accept field name to omit the leading dot and the wrapping braces
func toJSONPathExpression(field string) string {
	field = strings.TrimSpace(field)
	if strings.HasPrefix(field, "{") {
		return field
	}
	if !strings.HasPrefix(field, ".") {
		field = "." + field
	}

	return "{" + field + "}"
}
