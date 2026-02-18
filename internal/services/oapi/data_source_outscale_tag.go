package oapi

import (
	"time"

	"github.com/outscale/osc-sdk-go/v3/pkg/osc"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceOutscaleTag() *schema.Resource {
	return &schema.Resource{
		Read: DataSourceOutscaleTagRead,
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

func DataSourceOutscaleTagRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*client.OutscaleClient).OSC

	// Build up search parameters
	params := osc.ReadTagsRequest{}

	filters, filtersOk := d.GetOk("filter")

	var err error
	if filtersOk {
		params.Filters, err = oapiBuildOutscaleDataSourceFilters(filters.(*schema.Set))
		if err != nil {
			return err
		}
	}

	var resp osc.ReadTagsResponse
	err = retry.Retry(60*time.Second, func() *retry.RetryError {
		rp, httpResp, err := client.TagApi.ReadTags(ctx).ReadTagsRequest(params).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		resp = rp
		return nil
	})
	if err != nil {
		return err
	}

	if len(resp.Tags) < 1 {
		return ErrNoResults
	}

	if len(resp.Tags) > 1 {
		return ErrMultipleResults
	}

	tag := resp.Tags[0]

	if err := d.Set("key", tag.GetKey()); err != nil {
		return err
	}
	if err := d.Set("value", tag.GetValue()); err != nil {
		return err
	}
	if err := d.Set("resource_id", tag.GetResourceId()); err != nil {
		return err
	}

	if err := d.Set("resource_type", tag.GetResourceType()); err != nil {
		return err
	}

	d.SetId(id.UniqueId())

	return err
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
			filters.SetKeys(filterValues)
		case "resource_ids":
			filters.SetResourceIds(filterValues)
		case "resource_types":
			filters.SetResourceTypes(filterValues)
		case "values":
			filters.SetValues(filterValues)
		default:
			return nil, utils.UnknownDataSourceFilterError(ctx, name)
		}
	}

	return &filters, nil
}
