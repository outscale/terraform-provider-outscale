package oapi

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	fwschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	oscgo "github.com/outscale/osc-sdk-go/v2"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/outscale/terraform-provider-outscale/internal/framework/fwhelpers"
	"github.com/outscale/terraform-provider-outscale/internal/framework/fwhelpers/to"
	"github.com/outscale/terraform-provider-outscale/internal/utils"
	"github.com/samber/lo"
)

type TagsModel struct {
	Tags types.Set `tfsdk:"tags"`
}

type TagsComputedModel struct {
	Tags types.List `tfsdk:"tags"`
}

type ResourceTag struct {
	Key   types.String `tfsdk:"key"`
	Value types.String `tfsdk:"value"`
}

func TagsSchemaFW() *fwschema.SetNestedBlock {
	return &fwschema.SetNestedBlock{
		NestedObject: fwschema.NestedBlockObject{
			Attributes: map[string]fwschema.Attribute{
				"key": fwschema.StringAttribute{
					Required: true,
				},
				"value": fwschema.StringAttribute{
					Optional: true,
					Computed: true,
				},
			},
		},
	}
}

func TagsSchemaComputedFW() *fwschema.ListAttribute {
	return &fwschema.ListAttribute{
		Computed: true,
		ElementType: types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"key":   types.StringType,
				"value": types.StringType,
			},
		},
	}
}

func TagsSchemaSDK() *schema.Schema {
	return &schema.Schema{
		Type: schema.TypeSet,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"key": {
					Type:     schema.TypeString,
					Optional: true,
					Computed: true,
				},
				"value": {
					Type:     schema.TypeString,
					Computed: true,
					Optional: true,
				},
			},
		},
		Optional: true,
	}
}

func TagsSchemaComputedSDK() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeSet,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"key": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"value": {
					Type:     schema.TypeString,
					Computed: true,
				},
			},
		},
	}
}

var OAPITagAttrTypes = fwhelpers.GetAttrTypes(ResourceTag{})

func ComputedTagsNull() types.List {
	return types.ListNull(types.ObjectType{AttrTypes: OAPITagAttrTypes})
}

func TagsNull() types.Set {
	return types.SetNull(types.ObjectType{AttrTypes: OAPITagAttrTypes})
}

func createOAPITags(ctx context.Context, client *oscgo.APIClient, tags []oscgo.ResourceTag, resourceId string) error {
	req := oscgo.NewCreateTagsRequest([]string{resourceId}, tags)

	err := retry.RetryContext(ctx, 60*time.Second, func() *retry.RetryError {
		_, httpResp, err := client.TagApi.CreateTags(context.Background()).CreateTagsRequest(*req).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		return nil
	})

	return err
}

func createOAPITagsFW(ctx context.Context, client *oscgo.APIClient, tagsSet types.Set, resourceId string) diag.Diagnostics {
	if fwhelpers.IsSet(tagsSet) {
		var diags diag.Diagnostics
		tagsModel, diag := to.Slice[ResourceTag](ctx, tagsSet)
		tags := expandOAPITagsFW(tagsModel)
		if diag.HasError() {
			return diag
		}

		err := createOAPITags(ctx, client, tags, resourceId)
		if err != nil {
			diags.AddError(
				"Unable to create tags",
				err.Error(),
			)
			return diags
		}
	}

	return nil
}

func createOAPITagsSDK(client *oscgo.APIClient, d *schema.ResourceData) error {
	if tagsSchema, ok := d.GetOk("tags"); ok {
		set := tagsSchema.(*schema.Set)
		tags := expandOAPITagsSDK(set)
		resourceId := d.Id()

		err := createOAPITags(context.Background(), client, tags, resourceId)
		if err != nil {
			return fmt.Errorf("unable to create tags: %s", err)
		}
	}

	return nil
}

func updateOAPITags(ctx context.Context, client *oscgo.APIClient, toCreate, toRemove []oscgo.ResourceTag, resourceId string) error {
	if len(toRemove) > 0 {
		err := retry.RetryContext(ctx, 60*time.Second, func() *retry.RetryError {
			_, httpResp, err := client.TagApi.DeleteTags(ctx).DeleteTagsRequest(oscgo.DeleteTagsRequest{
				ResourceIds: []string{resourceId},
				Tags:        toRemove,
			}).Execute()
			if err != nil {
				return utils.CheckThrottling(httpResp, err)
			}
			return nil
		})
		if err != nil {
			return fmt.Errorf("unable to delete tags: %s", err.Error())
		}
	}
	if len(toCreate) > 0 {
		err := retry.RetryContext(ctx, 60*time.Second, func() *retry.RetryError {
			_, httpResp, err := client.TagApi.CreateTags(ctx).CreateTagsRequest(oscgo.CreateTagsRequest{
				ResourceIds: []string{resourceId},
				Tags:        toCreate,
			}).Execute()
			if err != nil {
				return utils.CheckThrottling(httpResp, err)
			}
			return nil
		})
		if err != nil {
			return fmt.Errorf("unable to create tags: %s", err.Error())
		}
	}

	return nil
}

func diffOAPITagsFW(ctx context.Context, oldSet, newSet types.Set) ([]oscgo.ResourceTag, []oscgo.ResourceTag, diag.Diagnostics) {
	var diags diag.Diagnostics

	oldTagsModel, diag := to.Slice[ResourceTag](ctx, oldSet)
	diags.Append(diag...)
	newTagsModel, diag := to.Slice[ResourceTag](ctx, newSet)
	diags.Append(diag...)
	if diag.HasError() {
		return nil, nil, diag
	}

	oldTags := expandOAPITagsFW(oldTagsModel)
	newTags := expandOAPITagsFW(newTagsModel)

	toCreate, toRemove := diffOAPITags(oldTags, newTags)
	return toCreate, toRemove, nil
}

func diffOAPITags(oldTags, newTags []oscgo.ResourceTag) ([]oscgo.ResourceTag, []oscgo.ResourceTag) {
	return lo.Difference(newTags, oldTags)
}

func updateOAPITagsSDK(client *oscgo.APIClient, d *schema.ResourceData) error {
	if d.HasChange("tags") {
		oldRaw, newRaw := d.GetChange("tags")
		old := oldRaw.(*schema.Set)
		new := newRaw.(*schema.Set)
		create, remove := diffOAPITags(expandOAPITagsSDK(old), expandOAPITagsSDK(new))
		resourceId := d.Id()

		return updateOAPITags(context.Background(), client, create, remove, resourceId)
	}

	return nil
}

func updateOAPITagsFW(ctx context.Context, client *oscgo.APIClient, oldSet, newSet types.Set, resourceId string) diag.Diagnostics {
	if oldSet.Equal(newSet) {
		return nil
	}
	var diags diag.Diagnostics

	create, remove, diag := diffOAPITagsFW(ctx, oldSet, newSet)
	diags.Append(diag...)
	if diags.HasError() {
		return diags
	}

	err := updateOAPITags(ctx, client, create, remove, resourceId)
	if err != nil {
		diags.AddError("Error updating tags", err.Error())
	}

	return diags
}

func flattenOAPITagsFW(ctx context.Context, tags []oscgo.ResourceTag) (types.Set, diag.Diagnostics) {
	tagsModel := lo.Map(tags, func(tag oscgo.ResourceTag, _ int) ResourceTag {
		return ResourceTag{
			Key:   to.String(tag.GetKey()),
			Value: to.String(tag.GetValue()),
		}
	})
	return types.SetValueFrom(ctx, types.ObjectType{AttrTypes: OAPITagAttrTypes}, tagsModel)
}

func flattenOAPIComputedTagsFW(ctx context.Context, tags []oscgo.ResourceTag) (types.List, diag.Diagnostics) {
	tagsModel := lo.Map(tags, func(tag oscgo.ResourceTag, _ int) ResourceTag {
		return ResourceTag{
			Key:   to.String(tag.GetKey()),
			Value: to.String(tag.GetValue()),
		}
	})
	return types.ListValueFrom(ctx, types.ObjectType{AttrTypes: OAPITagAttrTypes}, tagsModel)
}

func FlattenOAPITagsSDK(tags []oscgo.ResourceTag) []map[string]string {
	return lo.Map(tags, func(tag oscgo.ResourceTag, _ int) map[string]string {
		return map[string]string{
			"key":   tag.Key,
			"value": tag.Value,
		}
	})
}

func expandOAPITagsSDK(tags *schema.Set) []oscgo.ResourceTag {
	return lo.Map(tags.List(), func(v any, _ int) oscgo.ResourceTag {
		tag := v.(map[string]any)
		return oscgo.ResourceTag{
			Key:   tag["key"].(string),
			Value: tag["value"].(string),
		}
	})
}

func expandOAPITagsFW(tags []ResourceTag) []oscgo.ResourceTag {
	return lo.Map(tags, func(tag ResourceTag, _ int) oscgo.ResourceTag {
		return oscgo.ResourceTag{
			Key:   tag.Key.ValueString(),
			Value: tag.Value.ValueString(),
		}
	})
}

func oapiTagDescIgnored(t *oscgo.Tag) bool {
	filter := []string{"^outscale:"}
	for _, v := range filter {
		if r, _ := regexp.MatchString(v, t.GetKey()); r {
			return true
		}
	}
	return false
}

func flattenOAPITagsDescSDK(tags []oscgo.Tag) []map[string]any {
	res := make([]map[string]any, len(tags))

	for i, t := range tags {
		if !oapiTagDescIgnored(&t) {
			res[i] = map[string]any{
				"key":           t.Key,
				"value":         t.Value,
				"resource_id":   t.ResourceId,
				"resource_type": t.ResourceType,
			}
		}
	}
	return res
}
