package oapi

import (
	"context"
	"time"

	"github.com/outscale/goutils/sdk/ptr"
	"github.com/outscale/osc-sdk-go/v3/pkg/options"
	"github.com/outscale/osc-sdk-go/v3/pkg/osc"
	"github.com/outscale/terraform-provider-outscale/internal/client"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceOutscaleProductTypes() *schema.Resource {
	return &schema.Resource{
		ReadContext: DataSourceOutscaleProductTypesRead,

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

func DataSourceOutscaleProductTypesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.OutscaleClient).OSC

	req := osc.ReadProductTypesRequest{}

	var err error
	filters, filtersOk := d.GetOk("filter")

	if filtersOk {
		req.Filters, err = buildOutscaleProductTypeDataSourceFilters(filters.(*schema.Set))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	resp, err := client.ReadProductTypes(ctx, req, options.WithRetryTimeout(120*time.Second))

	var errString string
	if err != nil {
		errString = err.Error()
		return diag.Errorf("error reading product types (%s)", errString)
	}

	if resp.ProductTypes == nil || len(*resp.ProductTypes) == 0 {
		return diag.FromErr(ErrNoResults)
	}

	productTypes := make([]map[string]interface{}, len(*resp.ProductTypes))

	for k, productType := range *resp.ProductTypes {
		productTypeMap := make(map[string]interface{})

		if ptr.From(productType.Description) != "" {
			productTypeMap["description"] = productType.Description
		}
		if ptr.From(productType.ProductTypeId) != "" {
			productTypeMap["product_type_id"] = productType.ProductTypeId
		}
		if ptr.From(productType.Vendor) != "" {
			productTypeMap["vendor"] = productType.Vendor
		}

		productTypes[k] = productTypeMap
	}

	if err := d.Set("product_types", productTypes); err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id.UniqueId())

	return nil
}
