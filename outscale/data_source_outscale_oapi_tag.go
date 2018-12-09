package outscale

import (
	"fmt"
	"strings"
	"time"

	"github.com/terraform-providers/terraform-provider-outscale/osc/oapi"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
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
	conn := meta.(*OutscaleClient).OAPI

	// Build up search parameters
	params := oapi.ReadTagsRequest{}

	filters, filtersOk := d.GetOk("filter")

	if filtersOk {
		params.Filters = oapiBuildOutscaleDataSourceFilters(filters.(*schema.Set))
	}

	var resp *oapi.POST_ReadTagsResponses
	var err error

	err = resource.Retry(60*time.Second, func() *resource.RetryError {
		resp, err = conn.POST_ReadTags(params)
		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded") {
				return resource.RetryableError(err)
			}
		}

		return resource.NonRetryableError(err)
	})

	if err != nil {
		return err
	}

	if len(resp.OK.Tags) < 1 {
		return fmt.Errorf("your query returned no results, please change your search criteria and try again")
	}

	if len(resp.OK.Tags) > 1 {
		return fmt.Errorf("your query returned more than one result, Please try a more " +
			"specific search criteria")
	}

	tag := resp.OK.Tags[0]

	d.Set("key", tag.Key)
	d.Set("value", tag.Value)
	d.Set("resource_id", tag.ResourceId)
	d.Set("resource_type", tag.ResourceType)

	d.SetId(resource.UniqueId())

	return err
}

func oapiBuildOutscaleDataSourceFilters(set *schema.Set) oapi.FiltersTag {
	var filterKeys []string
	var filterValues []string
	for _, v := range set.List() {
		m := v.(map[string]interface{})

		for _, e := range m["values"].([]interface{}) {
			filterValues = append(filterValues, e.(string))
		}

		filterKeys = append(filterKeys, m["name"].(string))
	}

	filters := oapi.FiltersTag{
		Keys:   filterKeys,
		Values: filterValues,
	}
	return filters
}
