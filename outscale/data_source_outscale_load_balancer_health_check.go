package outscale

import (
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-outscale/osc/lbu"
)

func dataSourceOutscaleLoadBalancerHealthCheck() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOutscaleLoadBalancerHealthCheckRead,

		Schema: map[string]*schema.Schema{
			"load_balancer_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"healthy_threshold": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"unhealthy_threshold": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"target": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"interval": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"timeout": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"request_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceOutscaleLoadBalancerHealthCheckRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).LBU
	ename, ok := d.GetOk("load_balancer_name")

	if !ok {
		return fmt.Errorf("please provide the name of the load balancer")
	}

	elbName := ename.(string)

	// Retrieve the ELB properties for updating the state
	describeElbOpts := &lbu.DescribeLoadBalancersInput{
		LoadBalancerNames: []*string{aws.String(elbName)},
	}

	var describeResp *lbu.DescribeLoadBalancersOutput
	var err error
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		describeResp, err = conn.API.DescribeLoadBalancers(describeElbOpts)

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

	if describeResp.LoadBalancerDescriptions == nil {
		return fmt.Errorf("NO ELB FOUND")
	}

	if len(describeResp.LoadBalancerDescriptions) != 1 {
		return fmt.Errorf("Unable to find ELB: %#v", describeResp.LoadBalancerDescriptions)
	}

	lb := describeResp.LoadBalancerDescriptions[0]

	h := int64(0)
	i := int64(0)
	t := ""
	ti := int64(0)
	u := int64(0)

	if *lb.HealthCheck.Target != "" {
		h = aws.Int64Value(lb.HealthCheck.HealthyThreshold)
		i = aws.Int64Value(lb.HealthCheck.Interval)
		t = aws.StringValue(lb.HealthCheck.Target)
		ti = aws.Int64Value(lb.HealthCheck.Timeout)
		u = aws.Int64Value(lb.HealthCheck.UnhealthyThreshold)
	}

	d.Set("healthy_threshold", h)
	d.Set("interval", i)
	d.Set("target", t)
	d.Set("timeout", ti)
	d.Set("unhealthy_threshold", u)

	// d.Set("request_id", resp.ResponseMetadata.RequestID)
	d.SetId(*lb.LoadBalancerName)

	return nil
}
