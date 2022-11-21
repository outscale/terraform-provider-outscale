package outscale

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceLoadBalancerVms() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceLoadBalancerVmsRead,
		Schema: getDataSourceSchemas(map[string]*schema.Schema{
			"load_balancer_name": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},

			"backend_vm_ids": {
				Type:     schema.TypeList,
				ForceNew: true,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"request_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		}),
	}
}

func dataSourceLoadBalancerVmsRead(d *schema.ResourceData,
	meta interface{}) error {
	conn := meta.(*Client).OSCAPI

	lb, _, err := readLbs0(conn, d)
	if err != nil {
		return err
	}

	d.Set("backend_vm_ids", flattenStringList(lb.BackendVmIds))
	return nil
}
