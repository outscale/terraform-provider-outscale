package outscale

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func attrLBListenerDesc() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"load_balancer_name": {
			Type:     schema.TypeString,
			Required: true,
		},
		"listener": {
			Type:     schema.TypeMap,
			Computed: true,
			Elem: &schema.Resource{
				Schema: lb_listener_schema(true),
			},
		},
		"request_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}

func dataSourceOutscaleOAPILoadBalancerLD() *schema.Resource {
	return &schema.Resource{
		Read:   dataSourceOutscaleOAPILoadBalancerRead,
		Schema: getDataSourceSchemas(attrLBListenerDesc()),
	}
}

func dataSourceOutscaleOAPILoadBalancerLDRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	resp, elbName, err := readLbs(conn, d)
	if err != nil {
		return err
	}
	lbs := *resp.LoadBalancers
	if len(lbs) != 1 {
		return fmt.Errorf("Unable to find LBU: %s", *elbName)
	}

	lb := (lbs)[0]

	v := (*lb.Listeners)[0]

	l := make(map[string]interface{})
	l["backend_port"] = v.BackendPort
	l["backend_protocol"] = v.BackendProtocol
	l["load_balancer_port"] = v.LoadBalancerPort
	l["load_balancer_protocol"] = v.LoadBalancerProtocol
	l["server_certificate_id"] = v.ServerCertificateId

	if err := d.Set("listener", l); err != nil {
		return err
	}

	d.Set("request_id", resp.ResponseContext.RequestId)
	d.SetId(resource.UniqueId())

	return d.Set("policy_name", flattenStringList(v.PolicyNames))
}
