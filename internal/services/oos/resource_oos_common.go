package oos

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	s3types "github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/hashicorp/terraform-plugin-framework-validators/objectvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/outscale/goutils/sdk/ptr"
	"github.com/outscale/osc-sdk-go/v3/pkg/oos"
	"github.com/outscale/terraform-provider-outscale/internal/framework/fwhelpers"
	"github.com/outscale/terraform-provider-outscale/internal/framework/fwhelpers/from"
	"github.com/outscale/terraform-provider-outscale/internal/framework/fwhelpers/to"
	"github.com/outscale/terraform-provider-outscale/internal/framework/fwtypes"
	"github.com/samber/lo"
	"golang.org/x/sync/errgroup"
	"k8s.io/utils/keymutex"
)

type resourceCommon struct {
	Client *oos.Client
	mu     keymutex.KeyMutex
}

type grantModel struct {
	FullControl types.Object `tfsdk:"full_control"`
	Read        types.Object `tfsdk:"read"`
	Write       types.Object `tfsdk:"write"`
	ReadACP     types.Object `tfsdk:"read_acp"`
	WriteACP    types.Object `tfsdk:"write_acp"`
}

type grantPermissionModel struct {
	Ids            types.Set `tfsdk:"ids"`
	EmailAddresses types.Set `tfsdk:"email_addresses"`
}

type permissionsModel struct {
	Grants types.Set    `tfsdk:"grants"`
	Owner  types.Object `tfsdk:"owner"`
}

func (permissionsModel) AttributeTypes() map[string]attr.Type {
	return fwtypes.AttributeTypesAtPath(permissionsSchemaType)
}

type grantsModel struct {
	Permission types.String `tfsdk:"permission"`
	Grantee    types.Object `tfsdk:"grantee"`
}

func (grantsModel) AttributeTypes() map[string]attr.Type {
	return fwtypes.AttributeTypesAtPath(permissionsSchemaType, "grants")
}

type granteeModel struct {
	Type         types.String `tfsdk:"type"`
	DisplayName  types.String `tfsdk:"display_name"`
	EmailAddress types.String `tfsdk:"email_address"`
	Id           types.String `tfsdk:"id"`
	Uri          types.String `tfsdk:"uri"`
}

func (granteeModel) AttributeTypes() map[string]attr.Type {
	return fwtypes.AttributeTypesAtPath(permissionsSchemaType, "grants", "grantee")
}

type ownerModel struct {
	DisplayName types.String `tfsdk:"display_name"`
	Id          types.String `tfsdk:"id"`
}

func (ownerModel) AttributeTypes() map[string]attr.Type {
	return fwtypes.AttributeTypesAtPath(permissionsSchemaType, "owner")
}

type grantHeaders struct {
	FullControl *string
	Read        *string
	Write       *string
	ReadACP     *string
	WriteACP    *string
}

var permissionsSchemaType = permissionsAttributes().GetType()

func permissionsAttributes() schema.Attribute {
	return schema.SingleNestedAttribute{
		Computed: true,
		Attributes: map[string]schema.Attribute{
			"grants": schema.SetNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"permission": schema.StringAttribute{
							Computed: true,
						},
						"grantee": schema.SingleNestedAttribute{
							Computed: true,
							Attributes: map[string]schema.Attribute{
								"type": schema.StringAttribute{
									Computed: true,
								},
								"display_name": schema.StringAttribute{
									Computed: true,
								},
								"email_address": schema.StringAttribute{
									Computed: true,
								},
								"id": schema.StringAttribute{
									Computed: true,
								},
								"uri": schema.StringAttribute{
									Computed: true,
								},
							},
						},
					},
				},
			},
			"owner": schema.SingleNestedAttribute{
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"display_name": schema.StringAttribute{
						Computed: true,
					},
					"id": schema.StringAttribute{
						Computed: true,
					},
				},
			},
		},
	}
}

func grantAttributes() schema.Attribute {
	return schema.SingleNestedAttribute{
		Optional: true,
		Attributes: map[string]schema.Attribute{
			"full_control": schema.SingleNestedAttribute{
				Optional: true,
				Validators: []validator.Object{
					objectvalidator.AtLeastOneOf(
						path.MatchRoot("grant").AtName("read"),
						path.MatchRoot("grant").AtName("read_acp"),
						path.MatchRoot("grant").AtName("write"),
						path.MatchRoot("grant").AtName("write_acp"),
					),
				},
				Attributes: grantPermissionAttributes("full_control"),
			},
			"read": schema.SingleNestedAttribute{
				Optional:   true,
				Attributes: grantPermissionAttributes("read"),
			},
			"read_acp": schema.SingleNestedAttribute{
				Optional:   true,
				Attributes: grantPermissionAttributes("read_acp"),
			},
			"write": schema.SingleNestedAttribute{
				Optional:   true,
				Attributes: grantPermissionAttributes("write"),
			},
			"write_acp": schema.SingleNestedAttribute{
				Optional:   true,
				Attributes: grantPermissionAttributes("write_acp"),
			},
		},
	}
}

func grantPermissionAttributes(permission string) map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"ids": schema.SetAttribute{
			Optional:    true,
			ElementType: types.StringType,
			Validators: []validator.Set{
				setvalidator.SizeAtLeast(1),
				setvalidator.AtLeastOneOf(path.MatchRoot("grant").AtName(permission).AtName("email_addresses")),
			},
		},
		"email_addresses": schema.SetAttribute{
			Optional:    true,
			ElementType: types.StringType,
			Validators: []validator.Set{
				setvalidator.SizeAtLeast(1),
			},
		},
	}
}

func (r *resourceCommon) flattenPermissions(ctx context.Context, owner *s3types.Owner, aclGrants []s3types.Grant) (types.Object, error) {
	var obj types.Object
	model := permissionsModel{
		Owner: to.ObjectNull[ownerModel](),
	}

	grantObjects, err := lo.MapErr(aclGrants, func(grant s3types.Grant, _ int) (types.Object, error) {
		var grantObject types.Object
		model := grantsModel{
			Permission: to.String(grant.Permission),
			Grantee:    to.ObjectNull[granteeModel](),
		}

		if grant.Grantee != nil {
			granteeData := granteeModel{
				Type:         to.String(grant.Grantee.Type),
				DisplayName:  to.String(grant.Grantee.DisplayName),
				EmailAddress: to.String(grant.Grantee.EmailAddress),
				Id:           to.String(grant.Grantee.ID),
				Uri:          to.String(grant.Grantee.URI),
			}

			granteeObj, diags := to.Object(ctx, granteeData)
			if diags.HasError() {
				return grantObject, from.Diag(diags)
			}
			model.Grantee = granteeObj
		}

		grantObject, diags := to.Object(ctx, model)
		if diags.HasError() {
			return grantObject, from.Diag(diags)
		}

		return grantObject, nil
	})
	if err != nil {
		return obj, err
	}

	uniqueGrants := make([]attr.Value, 0, len(grantObjects))
	for _, grant := range grantObjects {
		if !lo.ContainsBy(uniqueGrants, grant.Equal) {
			uniqueGrants = append(uniqueGrants, grant)
		}
	}

	grants, diags := to.SetFromAttrType(ctx, uniqueGrants, to.ObjType(grantsModel{}.AttributeTypes()))
	if diags.HasError() {
		return obj, from.Diag(diags)
	}
	model.Grants = grants

	if owner != nil {
		ownerData := ownerModel{
			DisplayName: to.String(owner.DisplayName),
			Id:          to.String(owner.ID),
		}

		ownerObj, diags := to.Object(ctx, ownerData)
		if diags.HasError() {
			return obj, from.Diag(diags)
		}
		model.Owner = ownerObj
	}

	obj, diags = to.Object(ctx, model)
	if diags.HasError() {
		return obj, from.Diag(diags)
	}

	return obj, nil
}

func (r *resourceCommon) expandGrantHeaders(ctx context.Context, grant types.Object) (grantHeaders, diag.Diagnostics) {
	var diags diag.Diagnostics
	var headers grantHeaders

	model, diag := to.Model[grantModel](ctx, grant)
	diags.Append(diag...)
	if diags.HasError() {
		return headers, diags
	}

	format := func(value types.Object) *string {
		if !fwhelpers.IsSet(value) {
			return nil
		}

		permission, diag := to.Model[grantPermissionModel](ctx, value)
		diags.Append(diag...)
		if diag.HasError() {
			return nil
		}

		var grantees []string
		if fwhelpers.IsSet(permission.Ids) {
			ids, diag := to.Slice[string](ctx, permission.Ids)
			diags.Append(diag...)

			grantees = append(grantees, lo.Map(ids, func(id string, _ int) string {
				return "id=" + id
			})...)
		}
		if fwhelpers.IsSet(permission.EmailAddresses) {
			emails, diag := to.Slice[string](ctx, permission.EmailAddresses)
			diags.Append(diag...)

			grantees = append(grantees, lo.Map(emails, func(email string, _ int) string {
				return "emailaddress=" + email
			})...)
		}
		if diags.HasError() || len(grantees) == 0 {
			return nil
		}

		return new(strings.Join(grantees, ","))
	}

	headers.FullControl = format(model.FullControl)
	headers.Read = format(model.Read)
	headers.ReadACP = format(model.ReadACP)
	headers.Write = format(model.Write)
	headers.WriteACP = format(model.WriteACP)

	return headers, diags
}

func (r *resourceCommon) cleanObjects(ctx context.Context, timeout time.Duration, bucket string, key ...string) error {
	input := &s3.ListObjectVersionsInput{
		Bucket: &bucket,
	}
	if len(key) > 0 {
		input.Prefix = &key[0]
	}

	paginator := s3.NewListObjectVersionsPaginator(r.Client, input)
	var versions []s3types.ObjectVersion
	var objectIDs []s3types.ObjectIdentifier

	for paginator.HasMorePages() {
		objects, err := paginator.NextPage(ctx, oos.WithRetryTimeout(timeout))
		if err != nil {
			return err
		}

		for _, object := range objects.Versions {
			if len(key) > 0 && ptr.From(object.Key) != key[0] {
				continue
			}

			versions = append(versions, object)
			objectIDs = append(objectIDs, s3types.ObjectIdentifier{
				Key:       object.Key,
				VersionId: object.VersionId,
			})
		}
		for _, marker := range objects.DeleteMarkers {
			if len(key) > 0 && ptr.From(marker.Key) != key[0] {
				continue
			}

			objectIDs = append(objectIDs, s3types.ObjectIdentifier{
				Key:       marker.Key,
				VersionId: marker.VersionId,
			})
		}
	}
	if len(objectIDs) == 0 {
		return nil
	}

	now := time.Now()
	var retentionErr error
	var retentionErrMu sync.Mutex
	var group errgroup.Group

	// Verify that no object has a blocking retention policy before starting deletion
	for _, version := range versions {
		group.Go(func() error {
			output, err := r.Client.GetObjectRetention(ctx, &s3.GetObjectRetentionInput{
				Bucket:    &bucket,
				Key:       version.Key,
				VersionId: version.VersionId,
			}, oos.WithRetryTimeout(timeout))
			if err != nil {
				// If the Object Lock Configuration is not found, ignore the error
				if IsErrorMessage(err, "InvalidRequest", "Bucket is missing Object Lock Configuration") {
					return nil
				}

				retentionErrMu.Lock()
				retentionErr = errors.Join(retentionErr, fmt.Errorf(
					"check object %q version %q retention: %w",
					ptr.From(version.Key),
					ptr.From(version.VersionId),
					err,
				))
				retentionErrMu.Unlock()
				return nil
			}
			if output == nil || output.Retention == nil {
				return nil
			}

			var objectErr error
			switch {
			case output.Retention.RetainUntilDate == nil:
				objectErr = fmt.Errorf(
					"object %q version %q has %s retention without an expiration date",
					ptr.From(version.Key),
					ptr.From(version.VersionId),
					output.Retention.Mode,
				)
			case now.Before(*output.Retention.RetainUntilDate):
				objectErr = fmt.Errorf(
					"object %q version %q is protected by %s retention until %s",
					ptr.From(version.Key),
					ptr.From(version.VersionId),
					output.Retention.Mode,
					output.Retention.RetainUntilDate.Format(time.RFC3339),
				)
			}
			if objectErr == nil {
				return nil
			}

			retentionErrMu.Lock()
			retentionErr = errors.Join(retentionErr, objectErr)
			retentionErrMu.Unlock()
			return nil
		})
	}
	_ = group.Wait()
	if retentionErr != nil {
		return retentionErr
	}

	const batchSize = 1000
	for start := 0; start < len(objectIDs); start += batchSize {
		end := min(start+batchSize, len(objectIDs))
		output, err := r.Client.DeleteObjects(ctx, &s3.DeleteObjectsInput{
			Bucket: &bucket,
			Delete: &s3types.Delete{
				Objects: objectIDs[start:end],
			},
		}, oos.WithRetryTimeout(timeout))
		if err != nil {
			return fmt.Errorf("delete bucket objects: %w", err)
		}
		if output == nil {
			return ErrEmptyResponse
		}

		deleteErrors := lo.Map(output.Errors, func(objectErr s3types.Error, _ int) error {
			return fmt.Errorf(
				"delete object %q version %q: %s: %s",
				ptr.From(objectErr.Key),
				ptr.From(objectErr.VersionId),
				ptr.From(objectErr.Code),
				ptr.From(objectErr.Message),
			)
		})
		if err := errors.Join(deleteErrors...); err != nil {
			return err
		}
	}

	return nil
}
