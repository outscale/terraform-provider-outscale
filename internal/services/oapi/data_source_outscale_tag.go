package oapi

import (
	"context"
	"time"

	"github.com/outscale/osc-sdk-go/v3/pkg/options"
	"github.com/outscale/osc-sdk-go/v3/pkg/osc"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceOutscaleTag() *schema.Resource {
	return &schema.Resource{
		ReadContext: DataSourceOutscaleTagRead,
		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
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
	}
}

func DataSourceOutscaleTagRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.OutscaleClient).OSC

	// Build up search parameters
	params := osc.ReadTagsRequest{}

	filters, filtersOk := d.GetOk("filter")

	var err error
	if filtersOk {
		params.Filters, err = oapiBuildOutscaleDataSourceFilters(filters.(*schema.Set))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	resp, err := client.ReadTags(ctx, params, options.WithRetryTimeout(60*time.Second))
	if err != nil {
		return diag.FromErr(err)
	}

	if resp.Tags == nil || len(*resp.Tags) < 1 {
		return diag.FromErr(ErrNoResults)
	}

	if len(*resp.Tags) > 1 {
		return diag.FromErr(ErrMultipleResults)
	}

	tag := (*resp.Tags)[0]

	if err := d.Set("key", tag.Key); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("value", tag.Value); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("resource_id", tag.ResourceId); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("resource_type", tag.ResourceType); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(id.UniqueId())

	return diag.FromErr(err)
}

func oapiBuildOutscaleDataSourceFilters(set *schema.Set) (*osc.FiltersTag, error) {
	filters := osc.FiltersTag{}
	for _, v := range set.List() {
		m := v.(map[string]interface{})
		var filterValues []string

		for _, e := range m["values"].([]interface{}) {
			filterValues = append(filterValues, e.(string))
		}

		switch name := m["name"].(string); name {
		case "keys":
			filters.Keys = &filterValues
		case "resource_ids":
			filters.ResourceIds = &filterValues
		case "resource_types":
			filters.ResourceTypes = &filterValues
		case "values":
			filters.Values = &filterValues
		default:
			return nil, utils.UnknownDataSourceFilterError(name)
		}
	}

	return &filters, nil
}
