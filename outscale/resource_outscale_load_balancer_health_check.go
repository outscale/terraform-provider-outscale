package outscale

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-outscale/osc/lbu"
)

func resourceOutscaleLoadBalancerHealthCheck() *schema.Resource {
	return &schema.Resource{
		Read:   resourceOutscaleLoadBalancerHealthCheckRead,
		Create: resourceOutscaleLoadBalancerHealthCheckCreate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Delete: resourceOutscaleLoadBalancerHealthCheckDelete,

		Schema: map[string]*schema.Schema{
			"health_check": {
				Type:     schema.TypeMap,
				Required: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"healthy_threshold": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},
						"unhealthy_threshold": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},
						"target": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},
						"interval": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},
						"timeout": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},
					},
				},
			},
			"load_balancer_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"request_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceOutscaleLoadBalancerHealthCheckCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).LBU

	ename, ok := d.GetOk("load_balancer_name")
	hc, hok := d.GetOk("health_check")

	if !ok {
		return fmt.Errorf("please provide the name of the load balancer")
	}

	if !hok {
		return fmt.Errorf("please provide health check values")
	}

	check := hc.(map[string]interface{})

	ht, hterr := strconv.Atoi(check["healthy_threshold"].(string))
	ut, uterr := strconv.Atoi(check["unhealthy_threshold"].(string))
	i, ierr := strconv.Atoi(check["interval"].(string))
	t, terr := strconv.Atoi(check["timeout"].(string))

	if hterr != nil {
		return fmt.Errorf("please provider an number in health_check.healthy_threshold argument")
	}

	if uterr != nil {
		return fmt.Errorf("please provider an number in health_check.unhealthy_threshold argument")
	}

	if ierr != nil {
		return fmt.Errorf("please provider an number in health_check.interval argument")
	}

	if terr != nil {
		return fmt.Errorf("please provider an number in health_check.timeout argument")
	}

	configureHealthCheckOpts := lbu.ConfigureHealthCheckInput{
		LoadBalancerName: aws.String(ename.(string)),
		HealthCheck: &lbu.HealthCheck{
			HealthyThreshold:   aws.Int64(int64(ht)),
			UnhealthyThreshold: aws.Int64(int64(ut)),
			Interval:           aws.Int64(int64(i)),
			Target:             aws.String(check["target"].(string)),
			Timeout:            aws.Int64(int64(t)),
		},
	}
	var err error

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, err = conn.API.ConfigureHealthCheck(&configureHealthCheckOpts)

		if err != nil {
			if strings.Contains(err.Error(), "Throttling:") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("Failure configuring health check for ELB: %s", err)
	}

	d.SetId(ename.(string))

	return resourceOutscaleLoadBalancerHealthCheckRead(d, meta)
}

func resourceOutscaleLoadBalancerHealthCheckRead(d *schema.ResourceData, meta interface{}) error {
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

	h := ""
	i := ""
	t := ""
	ti := ""
	u := ""

	healthCheck := make(map[string]interface{})

	if *lb.HealthCheck.Target != "" {
		h = strconv.FormatInt(aws.Int64Value(lb.HealthCheck.HealthyThreshold), 10)
		i = strconv.FormatInt(aws.Int64Value(lb.HealthCheck.Interval), 10)
		t = aws.StringValue(lb.HealthCheck.Target)
		ti = strconv.FormatInt(aws.Int64Value(lb.HealthCheck.Timeout), 10)
		u = strconv.FormatInt(aws.Int64Value(lb.HealthCheck.UnhealthyThreshold), 10)
	}

	healthCheck["healthy_threshold"] = h
	healthCheck["interval"] = i
	healthCheck["target"] = t
	healthCheck["timeout"] = ti
	healthCheck["unhealthy_threshold"] = u

	d.Set("health_check", healthCheck)
	d.Set("load_balancer_name", *lb.LoadBalancerName)

	reqID := ""
	if describeResp.ResponseMetadata != nil {
		reqID = aws.StringValue(describeResp.ResponseMetadata.RequestID)
	}

	return d.Set("request_id", reqID)
}

func resourceOutscaleLoadBalancerHealthCheckDelete(d *schema.ResourceData, meta interface{}) error {
	d.SetId("")

	return nil
}
