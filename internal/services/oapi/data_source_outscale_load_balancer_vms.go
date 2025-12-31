package oapi

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/utils"
)

func DataSourceOutscaleLoadBalancerVms() *schema.Resource {
	return &schema.Resource{
		Read: DataSourceOutscaleLoadBalancerVmsRead,
		Schema: getDataSourceSchemas(map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			"load_balancer_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"backend_vm_ids": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"request_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		}),
	}
}

func DataSourceOutscaleLoadBalancerVmsRead(d *schema.ResourceData,
	meta interface{}) error {
	conn := meta.(*client.OutscaleClient).OSCAPI

	lb, _, err := readLbs0(conn, d)
	if err != nil {
		return err
	}

	d.Set("backend_vm_ids", utils.StringSlicePtrToInterfaceSlice(lb.BackendVmIds))
	return nil
}
