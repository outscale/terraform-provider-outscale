package outscale

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/antihax/optional"
	oscgo "github.com/marinsalinas/osc-sdk-go"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceOutscaleOAPILoadBalancerPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceOutscaleOAPILoadBalancerPolicyCreate,
		Read:   resourceOutscaleOAPILoadBalancerPolicyRead,
		Delete: resourceOutscaleOAPILoadBalancerPolicyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"policy_names": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"load_balancer_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"load_balancer_port": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"request_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func expandStringList(ifs []interface{}) *[]string {
	r := make([]string, len(ifs))

	for k, v := range ifs {
		r[k] = v.(string)
	}
	return &r
}

func resourceOutscaleOAPILoadBalancerPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	pInput := oscgo.UpdateLoadBalancerRequest{}

	pInput.PolicyNames = expandStringList(d.Get("policy_names").([]interface{}))

	if v, ok := d.GetOk("load_balancer_name"); ok {
		pInput.LoadBalancerName = v.(string)
	}

	if v, ok := d.GetOk("load_balancer_port"); ok {
		port := int64(v.(int))
		pInput.LoadBalancerPort = &port
	}

	opts := &oscgo.UpdateLoadBalancerOpts{
		optional.NewInterface(pInput),
	}

	var err error
	var resp oscgo.UpdateLoadBalancerResponse
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, _, err = conn.LoadBalancerApi.UpdateLoadBalancer(
			context.Background(), opts)

		if err != nil {
			if strings.Contains(fmt.Sprint(err), "Throttling") {
				return resource.RetryableError(
					fmt.Errorf("[WARN] Error creating ELB Policy, retrying: %s", err))
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		fmt.Printf("POLICY ERROR:%s", err)
		return err
	}

	if resp.ResponseContext != nil {
		d.Set("request_id", resp.ResponseContext.RequestId)
	}

	d.SetId(pInput.LoadBalancerName)
	log.Printf("[INFO] ELB Policies Listener ID: %s", d.Id())

	return nil
}

func resourceOutscaleOAPILoadBalancerPolicyRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceOutscaleOAPILoadBalancerPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	pInput := &oscgo.UpdateLoadBalancerRequest{}

	pols := make([]string, 0)
	pInput.PolicyNames = &pols

	if v, ok := d.GetOk("load_balancer_name"); ok {
		pInput.LoadBalancerName = v.(string)
	}

	if v, ok := d.GetOk("load_balancer_port"); ok {
		p := int64(v.(int))
		pInput.LoadBalancerPort = &p
	}

	opts := &oscgo.UpdateLoadBalancerOpts{
		optional.NewInterface(pInput),
	}

	var err error

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, _, err = conn.LoadBalancerApi.UpdateLoadBalancer(
			context.Background(), opts)

		if err != nil {
			if strings.Contains(fmt.Sprint(err), "Throttling") {
				return resource.RetryableError(
					fmt.Errorf("[WARN] Error creating ELB Policy, retrying: %s", err))
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return err
	}

	d.SetId("")

	return nil
}
