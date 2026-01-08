package oapi

import (
	"context"
	"fmt"
	"time"

	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceOutscaleProductTypes() *schema.Resource {
	return &schema.Resource{
		Read: DataSourceOutscaleProductTypesRead,

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

func DataSourceOutscaleProductTypesRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*client.OutscaleClient).OSCAPI

	req := oscgo.ReadProductTypesRequest{}

	var err error
	filters, filtersOk := d.GetOk("filter")

	if filtersOk {
		req.Filters, err = buildOutscaleProductTypeDataSourceFilters(filters.(*schema.Set))
		if err != nil {
			return err
		}
	}

	var resp oscgo.ReadProductTypesResponse
	err = retry.Retry(120*time.Second, func() *retry.RetryError {
		rp, httpResp, err := conn.ProductTypeApi.ReadProductTypes(context.Background()).ReadProductTypesRequest(req).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		resp = rp
		return nil
	})

	var errString string
	if err != nil {
		errString = err.Error()
		return fmt.Errorf("error reading product types (%s)", errString)
	}

	if len(resp.GetProductTypes()) == 0 {
		return ErrNoResults
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
	d.SetId(id.UniqueId())

	return nil
}
