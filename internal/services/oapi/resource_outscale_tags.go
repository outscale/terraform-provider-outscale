package oapi

import (
	"context"
	"time"

	"github.com/outscale/goutils/sdk/ptr"
	"github.com/outscale/osc-sdk-go/v3/pkg/options"
	"github.com/outscale/osc-sdk-go/v3/pkg/osc"
	"github.com/outscale/terraform-provider-outscale/internal/client"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceOutscaleTags() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceOutscaleTagsCreate,
		ReadContext:   ResourceOutscaleTagsRead,
		DeleteContext: ResourceOutscaleTagsDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		DeprecationMessage: "This resource is deprecated and will be removed in the next major version. Use the tags block of the specific resource instead.",

		Schema: getOAPITagsSchema(),
	}
}

func ResourceOutscaleTagsCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.OutscaleClient).OSC

	request := osc.CreateTagsRequest{}

	tag, tagsOk := d.GetOk("tag")
	resourceIds, resourceIdsOk := d.GetOk("resource_ids")
	if !tagsOk && !resourceIdsOk {
		return diag.Errorf("one tag and resource id, must be assigned")
	}

	if tagsOk {
		request.Tags = expandOAPITagsSDK(tag.(*schema.Set))
	}
	if resourceIdsOk {
		var rids []string
		sgs := resourceIds.(*schema.Set).List()
		for _, v := range sgs {
			str := v.(string)
			rids = append(rids, str)
		}

		request.ResourceIds = rids
	}

	_, err := client.CreateTags(ctx, request, options.WithRetryTimeout(60*time.Second))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(id.UniqueId())

	return ResourceOutscaleTagsRead(ctx, d, meta)
}

func ResourceOutscaleTagsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.OutscaleClient).OSC

	// Build up search parameters
	params := osc.ReadTagsRequest{
		Filters: &osc.FiltersTag{},
	}

	tag, tagsOk := d.GetOk("tag")
	filter := osc.FiltersTag{}
	if tagsOk {
		tgs := expandOAPITagsSDK(tag.(*schema.Set))
		keys := make([]string, 0, len(tgs))
		values := make([]string, 0, len(tgs))
		for _, t := range tgs {
			keys = append(keys, t.Key)
			values = append(values, t.Value)
		}
		filter.Keys = &keys
		filter.Values = &values
		params.Filters = &filter

	}

	resourceIds, resourceIdsOk := d.GetOk("resource_ids")
	if resourceIdsOk {
		var rids []string
		sgs := resourceIds.(*schema.Set).List()
		for _, v := range sgs {
			str := v.(string)
			rids = append(rids, str)
		}

		filter.ResourceIds = &rids
		params.Filters = &filter
	}

	resp, err := client.ReadTags(ctx, params, options.WithRetryTimeout(60*time.Second))
	if err != nil {
		return diag.FromErr(err)
	}

	tg := flattenOAPITagsDescSDK(ptr.From(resp.Tags))
	if err := d.Set("tags", tg); err != nil {
		return diag.FromErr(err)
	}

	return diag.FromErr(err)
}

func ResourceOutscaleTagsDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.OutscaleClient).OSC

	request := osc.DeleteTagsRequest{}

	tag, tagsOk := d.GetOk("tag")

	resourceIds, resourceIdsOk := d.GetOk("resource_ids")

	if !tagsOk && !resourceIdsOk {
		return diag.Errorf("one tag and resource id, must be assigned")
	}

	if tagsOk {
		request.Tags = expandOAPITagsSDK(tag.(*schema.Set))
	}
	if resourceIdsOk {
		var rids []string
		sgs := resourceIds.(*schema.Set).List()
		for _, v := range sgs {
			str := v.(string)
			rids = append(rids, str)
		}

		request.ResourceIds = rids
	}

	_, err := client.DeleteTags(ctx, request, options.WithRetryTimeout(60*time.Second))
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func getOAPITagsSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"resource_ids": {
			Type:     schema.TypeSet,
			Optional: true,
			ForceNew: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
		"tag": {
			Type:     schema.TypeSet,
			Optional: true,
			ForceNew: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"key": {
						Type:     schema.TypeString,
						Computed: true,
						Optional: true,
						ForceNew: true,
					},
					"value": {
						Type:     schema.TypeString,
						Optional: true,
						Computed: true,
						ForceNew: true,
					},
				},
			},
		},
		"tags": {
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
					"resource_id": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"resource_type": {
						Type:     schema.TypeString,
						Computed: true,
					},
				},
			},
		},
		"request_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}
