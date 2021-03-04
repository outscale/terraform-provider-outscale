package outscale

import (
	"context"
	"fmt"
	"strings"
	"time"

	oscgo "github.com/outscale/osc-sdk-go/v2"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceOutscaleOAPISubregion() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOutscaleOAPISubregionRead,
		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			"region_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"subregion_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"request_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceOutscaleOAPISubregionRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	filters, filtersOk := d.GetOk("filter")
	if !filtersOk {
		return fmt.Errorf("filters must be provided")
	}

	filtersReq := &oscgo.FiltersSubregion{}
	if filtersOk {
		filtersReq = buildOutscaleOAPIDataSourceSubregionsFilters(filters.(*schema.Set))
	}
	req := oscgo.ReadSubregionsRequest{Filters: filtersReq}

	var resp oscgo.ReadSubregionsResponse
	var err error
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, _, err = conn.SubregionApi.ReadSubregions(context.Background()).ReadSubregionsRequest(req).Execute()
		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded:") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return err
	}

	subregions := resp.GetSubregions()

	if len(subregions) < 1 {
		return fmt.Errorf("your query returned no results, please change your search criteria and try again")
	}
	if len(subregions) > 1 {
		return fmt.Errorf("your query returned more than one result, please try a more specific search criteria")
	}

	return resourceDataAttrSetter(d, func(set AttributeSetter) error {
		subregion := subregions[0]
		d.SetId(resource.UniqueId())

		if err = set("subregion_name", subregion.SubregionName); err != nil {
			return err
		}
		if err = set("region_name", subregion.RegionName); err != nil {
			return err
		}
		if err = set("state", subregion.State); err != nil {
			return err
		}

		return d.Set("request_id", resp.ResponseContext.RequestId)
	})
}
