package oapi

import (
	"context"
	"fmt"
	"time"

	"github.com/outscale/osc-sdk-go/v3/pkg/osc"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceOutscaleProductType() *schema.Resource {
	return &schema.Resource{
		Read: DataSourceOutscaleProductTypeRead,

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

func DataSourceOutscaleProductTypeRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*client.OutscaleClient).OSC

	req := osc.ReadProductTypesRequest{}

	filters, filtersOk := d.GetOk("filter")

	var err error
	if filtersOk {
		req.Filters, err = buildOutscaleProductTypeDataSourceFilters(filters.(*schema.Set))
		if err != nil {
			return err
		}
	}

	var resp osc.ReadProductTypesResponse
	err = retry.Retry(120*time.Second, func() *retry.RetryError {
		var err error
		rp, httpResp, err := client.ProductTypeApi.ReadProductTypes(ctx).ReadProductTypesRequest(req).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		resp = rp
		return nil
	})
	if err != nil {
		errString := err.Error()
		return fmt.Errorf("error reading producttype (%s)", errString)
	}

	if len(resp.GetProductTypes()) == 0 {
		return ErrNoResults
	}

	if len(resp.GetProductTypes()) > 1 {
		return ErrMultipleResults
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

func buildOutscaleProductTypeDataSourceFilters(set *schema.Set) (*osc.FiltersProductType, error) {
	var filters osc.FiltersProductType
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
			return nil, utils.UnknownDataSourceFilterError(ctx, name)
		}
	}
	return &filters, nil
}
