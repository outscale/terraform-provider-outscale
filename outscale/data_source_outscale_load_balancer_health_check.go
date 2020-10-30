package outscale

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceOutscaleOAPILoadBalancerHealthCheck() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOutscaleOAPILoadBalancerHealthCheckRead,
		Schema: getDataSourceSchemas(
			map[string]*schema.Schema{
				"load_balancer_name": {
					Type:     schema.TypeString,
					Required: true,
				},
				"healthy_threshold": {
					Type:     schema.TypeInt,
					Computed: true,
				},
				"unhealthy_threshold": {
					Type:     schema.TypeInt,
					Computed: true,
				},
				"path": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"check_interval": {
					Type:     schema.TypeInt,
					Computed: true,
				},
				"timeout": {
					Type:     schema.TypeInt,
					Computed: true,
				},
				"request_id": {
					Type:     schema.TypeString,
					Computed: true,
				},
			}),
	}
}

func dataSourceOutscaleOAPILoadBalancerHealthCheckRead(d *schema.ResourceData, meta interface{}) error {
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

	if lb.AccessLog == nil {
		return fmt.Errorf("NO Attributes FOUND")
	}

	h := int32(0)
	i := int32(0)
	t := ""
	ti := int32(0)
	u := int32(0)

	h = lb.HealthCheck.HealthyThreshold
	i = lb.HealthCheck.CheckInterval
	t = *lb.HealthCheck.Path
	ti = lb.HealthCheck.Timeout
	u = lb.HealthCheck.UnhealthyThreshold

	d.Set("healthy_threshold", h)
	d.Set("check_interval", i)
	d.Set("path", t)
	d.Set("timeout", ti)
	d.Set("unhealthy_threshold", u)

	d.Set("request_id", resp.ResponseContext.RequestId)
	d.SetId(*lb.LoadBalancerName)

	return nil
}
