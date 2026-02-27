package oapi

import (
	"context"
	"time"

	"github.com/outscale/goutils/sdk/ptr"

	"github.com/outscale/osc-sdk-go/v3/pkg/options"
	"github.com/outscale/osc-sdk-go/v3/pkg/osc"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceOutscaleProductType() *schema.Resource {
	return &schema.Resource{
		ReadContext: DataSourceOutscaleProductTypeRead,

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

func DataSourceOutscaleProductTypeRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.OutscaleClient).OSC

	req := osc.ReadProductTypesRequest{}

	filters, filtersOk := d.GetOk("filter")

	var err error
	if filtersOk {
		req.Filters, err = buildOutscaleProductTypeDataSourceFilters(filters.(*schema.Set))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	resp, err := client.ReadProductTypes(ctx, req, options.WithRetryTimeout(120*time.Second))
	if err != nil {
		errString := err.Error()
		return diag.Errorf("error reading producttype (%s)", errString)
	}

	if resp.ProductTypes == nil || len(*resp.ProductTypes) == 0 {
		return diag.FromErr(ErrNoResults)
	}

	if len(*resp.ProductTypes) > 1 {
		return diag.FromErr(ErrMultipleResults)
	}

	productType := (*resp.ProductTypes)[0]

	d.SetId(ptr.From(productType.ProductTypeId))
	if err := d.Set("product_type_id", ptr.From(productType.ProductTypeId)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("description", ptr.From(productType.Description)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("vendor", ptr.From(productType.Vendor)); err != nil {
		return diag.FromErr(err)
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
			return nil, utils.UnknownDataSourceFilterError(name)
		}
	}
	return &filters, nil
}
