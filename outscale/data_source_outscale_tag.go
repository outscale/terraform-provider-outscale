package outscale

import (
	"context"
	"fmt"
	"log"
	"time"

	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/terraform-providers/terraform-provider-outscale/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceOutscaleOAPITag() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOutscaleOAPITagRead,
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

func dataSourceOutscaleOAPITagRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	// Build up search parameters
	params := oscgo.ReadTagsRequest{}

	filters, filtersOk := d.GetOk("filter")

	if filtersOk {
		params.SetFilters(oapiBuildOutscaleDataSourceFilters(filters.(*schema.Set)))
	}

	var resp oscgo.ReadTagsResponse
	var err error

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

	d.SetId(resource.UniqueId())

	return err
}

func oapiBuildOutscaleDataSourceFilters(set *schema.Set) oscgo.FiltersTag {
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
			log.Printf("[Debug] Unknown Filter Name: %s.", name)
		}
	}

	return filters
}
