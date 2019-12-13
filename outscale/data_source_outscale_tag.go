package outscale

import (
	"context"
	"fmt"
	"github.com/antihax/optional"
	oscgo "github.com/marinsalinas/osc-sdk-go"
	"strings"
	"time"

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
		resp, _, err = conn.TagApi.ReadTags(context.Background(), &oscgo.ReadTagsOpts{ReadTagsRequest: optional.NewInterface(params)})
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

	if len(resp.GetTags()) < 1 {
		return fmt.Errorf("your query returned no results, please change your search criteria and try again")
	}

	if len(resp.GetTags()) > 1 {
		return fmt.Errorf("your query returned more than one result, Please try a more " +
			"specific search criteria")
	}

	tag := resp.GetTags()[0]

	d.Set("key", tag.GetKey())
	d.Set("value", tag.GetValue())
	d.Set("resource_id", tag.GetResourceId())
	d.Set("resource_type", tag.GetResourceType())

	d.SetId(resource.UniqueId())

	return err
}

func oapiBuildOutscaleDataSourceFilters(set *schema.Set) oscgo.FiltersTag {
	var filterKeys []string
	var filterValues []string
	for _, v := range set.List() {
		m := v.(map[string]interface{})

		for _, e := range m["values"].([]interface{}) {
			filterValues = append(filterValues, e.(string))
		}

		filterKeys = append(filterKeys, m["name"].(string))
	}

	filters := oscgo.FiltersTag{
		Keys:   &filterKeys,
		Values: &filterValues,
	}
	return filters
}
