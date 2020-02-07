package outscale

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-outscale/osc/fcu"
)

func dataSourceOutscaleOAPIProductTypes() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOutscaleOAPIProductTypesRead,

		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			"product_type": &schema.Schema{
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
						"product_type_vendor": {
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
	conn := meta.(*OutscaleClient).FCU

	filters, filtersOk := d.GetOk("filter")

	params := &fcu.DescribeProductTypesInput{}

	if filtersOk {
		params.Filters = buildOutscaleDataSourceFilters(filters.(*schema.Set))
	}

	var resp *fcu.DescribeProductTypesOutput
	var err error

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, err = conn.VM.DescribeProductTypes(params)
		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded:") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return resource.NonRetryableError(err)
	})

	if err != nil {
		return err
	}
	if resp == nil || len(resp.ProductTypeSet) == 0 {
		return fmt.Errorf("no matching Product Types found: %#v", params)
	}

	vcs := make([]map[string]interface{}, len(resp.ProductTypeSet))

	for k, v := range resp.ProductTypeSet {
		vc := make(map[string]interface{})
		vc["description"] = *v.Description
		vc["product_type_id"] = *v.ProductTypeId
		if v.Vendor != nil {
			vc["product_type_vendor"] = *v.Vendor
		} else {
			vc["product_type_vendor"] = ""
		}
		vcs[k] = vc
	}

	if err := d.Set("product_type", vcs); err != nil {
		return err
	}
	d.Set("request_id", resp.RequestId)
	d.SetId(resource.UniqueId())

	return nil
}
