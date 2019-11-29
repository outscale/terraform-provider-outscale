package outscale

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-outscale/osc/fcu"
)

func dataSourceOutscaleOAPIVpcEndpointServices() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOutscaleOAPIVpcEndpointServicesRead,

		Schema: map[string]*schema.Schema{
			"request_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"prefix_list_name": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataSourceOutscaleOAPIVpcEndpointServicesRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	params := &fcu.DescribeVpcEndpointServicesInput{}
	var res *fcu.DescribeVpcEndpointServicesOutput
	var err error
	err = resource.Retry(40*time.Minute, func() *resource.RetryError {
		res, err = conn.VM.DescribeVpcEndpointServices(params)

		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded") {
				fmt.Printf("[INFO] Request limit exceeded")
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}

		return resource.RetryableError(err)
	})

	if err != nil {
		return err
	}

	if len(res.ServiceNames) < 1 {
		return fmt.Errorf("your query returned no results, please change your search criteria and try again")
	}

	d.SetId(resource.UniqueId())
	d.Set("request_id", res.RequestID)

	return d.Set("prefix_list_name", flattenStringList(res.ServiceNames))
}
