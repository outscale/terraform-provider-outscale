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

func dataSourceOutscaleOAPIProductType() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOutscaleOAPIProductTypeRead,

		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"product_type_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vendor": {
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

func dataSourceOutscaleOAPIProductTypeRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	req := oscgo.ReadProductTypesRequest{}

	filters, filtersOk := d.GetOk("filter")

	if filtersOk {
		req.Filters = buildOutscaleOAPIProductTypeDataSourceFilters(filters.(*schema.Set))
	}

	var resp oscgo.ReadProductTypesResponse
	err := resource.Retry(120*time.Second, func() *resource.RetryError {
		var err error
		rp, httpResp, err := conn.ProductTypeApi.ReadProductTypes(context.Background()).ReadProductTypesRequest(req).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		resp = rp
		return nil
	})

	if err != nil {
		errString := err.Error()
		return fmt.Errorf("[DEBUG] Error reading ProductType (%s)", errString)
	}

	if len(resp.GetProductTypes()) == 0 {
		return fmt.Errorf("no matching product type found")
	}

	if len(resp.GetProductTypes()) > 1 {
		return fmt.Errorf("multiple product type matched; use additional constraints to reduce matches to a single product type")
	}

	productType := resp.GetProductTypes()[0]

	d.SetId(productType.GetProductTypeId())
	if err := d.Set("product_type_id", productType.GetProductTypeId()); err != nil {
		return err
	}
	if err := d.Set("description", productType.GetDescription()); err != nil {
		return err
	}
	if err := d.Set("vendor", productType.GetVendor()); err != nil {
		return err
	}

	return nil
}

func buildOutscaleOAPIProductTypeDataSourceFilters(set *schema.Set) *oscgo.FiltersProductType {
	var filters oscgo.FiltersProductType
	for _, v := range set.List() {
		m := v.(map[string]interface{})
		var filterValues []string
		for _, e := range m["values"].([]interface{}) {
			filterValues = append(filterValues, e.(string))
		}

		switch name := m["name"].(string); name {
		case "product_type_ids":
			filters.ProductTypeIds = &filterValues
		default:
			log.Printf("[Debug] Unknown Filter Name: %s.", name)
		}
	}
	return &filters
}
