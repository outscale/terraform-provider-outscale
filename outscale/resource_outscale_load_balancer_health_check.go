package outscale

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/antihax/optional"
	oscgo "github.com/marinsalinas/osc-sdk-go"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceOutscaleOAPILoadBalancerHealthCheck() *schema.Resource {
	return &schema.Resource{
		Read:   resourceOutscaleOAPILoadBalancerHealthCheckRead,
		Create: resourceOutscaleOAPILoadBalancerHealthCheckCreate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Delete: resourceOutscaleOAPILoadBalancerHealthCheckDelete,

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
						"path": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},
						"check_interval": {
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

func isLoadBalancerNotFound(err error) bool {
	return strings.Contains(fmt.Sprint(err), "LoadBalancerNotFound")
}

func resourceOutscaleOAPILoadBalancerHealthCheckCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

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
	i, ierr := strconv.Atoi(check["check_interval"].(string))
	t, terr := strconv.Atoi(check["timeout"].(string))
	p, perr := strconv.Atoi(check["port"].(string))

	if hterr != nil {
		return fmt.Errorf("please provide an number in health_check.healthy_threshold argument")
	}

	if uterr != nil {
		return fmt.Errorf("please provide an number in health_check.unhealthy_threshold argument")
	}

	if ierr != nil {
		return fmt.Errorf("please provide an number in health_check.check_interval argument")
	}

	if terr != nil {
		return fmt.Errorf("please provide an number in health_check.timeout argument")
	}

	if perr != nil {
		return fmt.Errorf("please provide an number in health_check.port argument")
	}

	req := oscgo.UpdateLoadBalancerRequest{
		LoadBalancerName: ename.(string),
		HealthCheck: &oscgo.HealthCheck{
			HealthyThreshold:   int64(ht),
			UnhealthyThreshold: int64(ut),
			CheckInterval:      int64(i),
			Protocol:           check["protocol"].(string),
			Port:               int64(p),
			Timeout:            int64(t),
		},
	}
	if check["path"] != nil {
		req.HealthCheck.Path = check["path"].(string)
	}

	configureHealthCheckOpts := oscgo.UpdateLoadBalancerOpts{
		optional.NewInterface(req),
	}

	var err error

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, _, err = conn.LoadBalancerApi.UpdateLoadBalancer(
			context.Background(), &configureHealthCheckOpts)

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

	return resourceOutscaleOAPILoadBalancerHealthCheckRead(d, meta)
}

func resourceOutscaleOAPILoadBalancerHealthCheckRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI
	ename, ok := d.GetOk("load_balancer_name")

	if !ok {
		return fmt.Errorf("please provide the name of the load balancer")
	}

	elbName := ename.(string)

	// Retrieve the ELB properties for updating the state
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

	if resp.LoadBalancers == nil {
		return fmt.Errorf("NO ELB FOUND")
	}

	if len(*resp.LoadBalancers) != 1 {
		return fmt.Errorf("Unable to find ELB: %#v", resp.LoadBalancers)
	}

	lb := (*resp.LoadBalancers)[0]

	h := ""
	i := ""
	t := ""
	ti := ""
	u := ""

	healthCheck := make(map[string]interface{})

	if lb.HealthCheck.Path != "" {
		h = strconv.FormatInt(lb.HealthCheck.HealthyThreshold, 10)
		i = strconv.FormatInt(lb.HealthCheck.CheckInterval, 10)
		t = lb.HealthCheck.Path
		ti = strconv.FormatInt(lb.HealthCheck.Timeout, 10)
		u = strconv.FormatInt(lb.HealthCheck.UnhealthyThreshold, 10)
	}

	healthCheck["healthy_threshold"] = h
	healthCheck["check_interval"] = i
	healthCheck["path"] = t
	healthCheck["timeout"] = ti
	healthCheck["unhealthy_threshold"] = u

	d.Set("health_check", healthCheck)
	d.Set("load_balancer_name", *lb.LoadBalancerName)

	reqID := ""
	if resp.ResponseContext != nil {
		reqID = *resp.ResponseContext.RequestId
	}

	return d.Set("request_id", reqID)
}

func resourceOutscaleOAPILoadBalancerHealthCheckDelete(d *schema.ResourceData, meta interface{}) error {
	d.SetId("")

	return nil
}
