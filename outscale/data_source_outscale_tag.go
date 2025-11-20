package outscale

import (
	"context"
	"fmt"
	"time"

	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/outscale/terraform-provider-outscale/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
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
	conn := meta.(*OutscaleClient).OSCAPI

	// Build up search parameters
	params := oscgo.ReadTagsRequest{}

	filters, filtersOk := d.GetOk("filter")

	var err error
	if filtersOk {
		params.Filters, err = oapiBuildOutscaleDataSourceFilters(filters.(*schema.Set))
		if err != nil {
			return err
		}
	}

	var resp oscgo.ReadTagsResponse
	err = resource.Retry(60*time.Second, func() *resource.RetryError {
		rp, httpResp, err := conn.TagApi.ReadTags(context.Background()).ReadTagsRequest(params).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		resp = rp
		return nil
	})

	if err != nil {
		return err
	}

	if len(resp.GetTags()) < 1 {
		return fmt.Errorf("your query returned no results, please change your search criteria and try again")
	}

	if len(resp.GetTags()) > 1 {
		return fmt.Errorf("your query returned more than one result, Please try a more " +
			"specific search criteria")
	}

	tag := resp.GetTags()[0]

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

func oapiBuildOutscaleDataSourceFilters(set *schema.Set) (*oscgo.FiltersTag, error) {
	filters := oscgo.FiltersTag{}
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
			return nil, utils.UnknownDataSourceFilterError(context.Background(), name)
		}
	}

	return &filters, nil
}
