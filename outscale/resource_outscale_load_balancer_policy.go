package outscale

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-outscale/osc/lbu"
)

func resourceOutscaleLoadBalancerPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceOutscaleLoadBalancerPolicyCreate,
		Read:   resourceOutscaleLoadBalancerPolicyRead,
		Delete: resourceOutscaleLoadBalancerPolicyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"policy_names": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"load_balancer_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"load_balancer_port": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"request_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceOutscaleLoadBalancerPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).LBU

	pInput := &lbu.SetLoadBalancerPoliciesOfListenerInput{}

	pInput.PolicyNames = expandStringList(d.Get("policy_names").([]interface{}))

	if v, ok := d.GetOk("load_balancer_name"); ok {
		pInput.LoadBalancerName = aws.String(v.(string))
	}

	if v, ok := d.GetOk("load_balancer_port"); ok {
		pInput.LoadBalancerPort = aws.Int64(int64(v.(int)))
	}

	var err error
	var resp *lbu.SetLoadBalancerPoliciesOfListenerOutput
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, err = conn.API.SetLoadBalancerPoliciesOfListener(pInput)

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

	if resp.ResponseMatadata != nil {
		d.Set("request_id", resp.ResponseMatadata.RequestID)
	}

	d.SetId(*pInput.LoadBalancerName)
	log.Printf("[INFO] ELB Policies Listener ID: %s", d.Id())

	return nil
}

func resourceOutscaleLoadBalancerPolicyRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceOutscaleLoadBalancerPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).LBU

	pInput := &lbu.SetLoadBalancerPoliciesOfListenerInput{}

	pInput.PolicyNames = make([]*string, 0)

	if v, ok := d.GetOk("load_balancer_name"); ok {
		pInput.LoadBalancerName = aws.String(v.(string))
	}

	if v, ok := d.GetOk("load_balancer_port"); ok {
		pInput.LoadBalancerPort = aws.Int64(int64(v.(int)))
	}

	var err error

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, err = conn.API.SetLoadBalancerPoliciesOfListener(pInput)

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
