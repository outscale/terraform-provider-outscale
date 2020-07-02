package outscale

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/antihax/optional"
	oscgo "github.com/marinsalinas/osc-sdk-go"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceOutscaleOAPILoadBalancerHealthCheck() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOutscaleOAPILoadBalancerHealthCheckRead,

		Schema: map[string]*schema.Schema{
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
		},
	}
}

func dataSourceOutscaleOAPILoadBalancerHealthCheckRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI
	ename, ok := d.GetOk("load_balancer_name")

	if !ok {
		return fmt.Errorf("please provide the name of the load balancer")
	}

	elbName := ename.(string)

	filter := &oscgo.FiltersLoadBalancer{
		LoadBalancerNames: &[]string{elbName},
	}

	req := oscgo.ReadLoadBalancersRequest{
		Filters: filter,
	}

	describeElbOpts := &oscgo.ReadLoadBalancersOpts{
		ReadLoadBalancersRequest: optional.NewInterface(req),
	}

	var resp oscgo.ReadLoadBalancersResponse
	var err error
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, _, err = conn.LoadBalancerApi.ReadLoadBalancers(
			context.Background(),
			describeElbOpts)

		if err != nil {
			if strings.Contains(fmt.Sprint(err), "Throttling:") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		if isLoadBalancerNotFound(err) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error retrieving ELB: %s", err)
	}

	lbs := *resp.LoadBalancers
	if len(lbs) != 1 {
		return fmt.Errorf("Unable to find LBU: %s", elbName)
	}

	lb := (lbs)[0]

	if lb.AccessLog == nil {
		return fmt.Errorf("NO Attributes FOUND")
	}

	h := int64(0)
	i := int64(0)
	t := ""
	ti := int64(0)
	u := int64(0)

	if lb.HealthCheck.Path != "" {
		h = lb.HealthCheck.HealthyThreshold
		i = lb.HealthCheck.CheckInterval
		t = lb.HealthCheck.Path
		ti = lb.HealthCheck.Timeout
		u = lb.HealthCheck.UnhealthyThreshold
	}

	d.Set("healthy_threshold", h)
	d.Set("check_interval", i)
	d.Set("path", t)
	d.Set("timeout", ti)
	d.Set("unhealthy_threshold", u)

	d.Set("request_id", resp.ResponseContext.RequestId)
	d.SetId(*lb.LoadBalancerName)

	return nil
}
