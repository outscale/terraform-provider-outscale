package outscale

import (
	"context"
	"fmt"
	"time"

	oscgo "github.com/outscale/osc-sdk-go/v2"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-outscale/utils"
)

func dataSourceOutscaleOAPIProductTypes() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOutscaleOAPIProductTypesRead,

		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			"product_types": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
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
					},
				},
			},
			"request_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceOutscaleOAPIProductTypesRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	req := oscgo.ReadProductTypesRequest{}

	filters, filtersOk := d.GetOk("filter")

	if filtersOk {
		req.Filters = buildOutscaleOAPIProductTypeDataSourceFilters(filters.(*schema.Set))
	}

	var resp oscgo.ReadProductTypesResponse
	var err error
	err = resource.Retry(120*time.Second, func() *resource.RetryError {
		resp, _, err = conn.ProductTypeApi.ReadProductTypes(context.Background()).ReadProductTypesRequest(req).Execute()
		if err != nil {
			return utils.CheckThrottling(err)
		}
		return nil
	})

	var errString string
	if err != nil {
		errString = err.Error()
		return fmt.Errorf("[DEBUG] Error reading product types (%s)", errString)
	}

	if len(resp.GetProductTypes()) == 0 {
		return fmt.Errorf("no matching product types found")
	}

	productTypes := make([]map[string]interface{}, len(resp.GetProductTypes()))

	for k, productType := range resp.GetProductTypes() {
		productTypeMap := make(map[string]interface{})

		if productType.GetDescription() != "" {
			productTypeMap["description"] = productType.GetDescription()
		}
		if productType.GetProductTypeId() != "" {
			productTypeMap["product_type_id"] = productType.GetProductTypeId()
		}
		if productType.GetVendor() != "" {
			productTypeMap["vendor"] = productType.GetVendor()
		}

		productTypes[k] = productTypeMap
	}

	if err := d.Set("product_types", productTypes); err != nil {
		return err
	}
	d.SetId(resource.UniqueId())

	return nil
}
