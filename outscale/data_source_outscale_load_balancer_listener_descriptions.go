package outscale

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceOutscaleOAPILoadBalancerLDs() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOutscaleOAPILoadBalancerLDsRead,

		Schema: getDataSourceSchemas(
			map[string]*schema.Schema{
				"load_balancer_names": {
					Type:     schema.TypeList,
					Required: true,
					Elem:     &schema.Schema{Type: schema.TypeString},
				},
				"listeners": {
					Type:     schema.TypeSet,
					Computed: true,
					Elem: &schema.Resource{
						Schema: lb_listener_schema(true),
					},
				},
				"request_id": {
					Type:     schema.TypeString,
					Computed: true,
				},
			}),
	}
}

func dataSourceOutscaleOAPILoadBalancerLDsRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	lb, resp, err := readLbs0(conn, d)
	if err != nil {
		return err
	}

	if lb.Listeners != nil {
		if err := d.Set("listeners", flattenOAPIListeners(lb.Listeners)); err != nil {
			return err
		}
	} else {
		if err := d.Set("listeners", make([]map[string]interface{}, 0)); err != nil {
			return err
		}
	}
	d.Set("request_id", resp.ResponseContext.RequestId)
	d.SetId(resource.UniqueId())

	return nil
}
