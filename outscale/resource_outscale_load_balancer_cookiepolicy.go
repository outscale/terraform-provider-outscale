package outscale

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-outscale/osc/lbu"
)

func resourceOutscaleAppCookieStickinessPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceOutscaleAppCookieStickinessPolicyCreate,
		Read:   resourceOutscaleAppCookieStickinessPolicyRead,
		Delete: resourceOutscaleAppCookieStickinessPolicyDelete,

		Schema: map[string]*schema.Schema{
			"policy_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, es []error) {
					value := v.(string)
					if !regexp.MustCompile(`^[0-9A-Za-z-]+$`).MatchString(value) {
						es = append(es, fmt.Errorf(
							"only alphanumeric characters and hyphens allowed in %q", k))
					}
					return
				},
			},

			"load_balancer_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"cookie_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"request_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceOutscaleAppCookieStickinessPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).LBU

	l, lok := d.GetOk("load_balancer_name")
	p, pok := d.GetOk("policy_name")
	v, ok := d.GetOk("cookie_name")

	if !lok && !pok && !ok {
		return fmt.Errorf("please provide the required attributes load_balancer_name, policy_name and cookie_name")
	}

	acspOpts := &lbu.CreateAppCookieStickinessPolicyInput{
		LoadBalancerName: aws.String(l.(string)),
		PolicyName:       aws.String(p.(string)),
		CookieName:       aws.String(v.(string)),
	}

	var err error
	var resp *lbu.CreateAppCookieStickinessPolicyOutput
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, err = conn.API.CreateAppCookieStickinessPolicy(acspOpts)

		if err != nil {
			if strings.Contains(fmt.Sprint(err), "Throttling") {
				return resource.RetryableError(
					fmt.Errorf("[WARN] Error creating AppCookieStickinessPolicy, retrying: %s", err))
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("Error creating AppCookieStickinessPolicy: %s", err)
	}

	//utils.PrintToJSON(resp, "RESPONSECookie")

	reqID := ""
	if resp.ResponseMetadata != nil {
		reqID = aws.StringValue(resp.ResponseMetadata.RequestID)
	}
	d.Set("request_id", reqID)
	d.SetId(resource.UniqueId())
	d.Set("load_balancer_name", l.(string))
	d.Set("policy_name", p.(string))
	d.Set("cookie_name", v.(string))

	return resourceOutscaleAppCookieStickinessPolicyRead(d, meta)
}

func resourceOutscaleAppCookieStickinessPolicyRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceOutscaleAppCookieStickinessPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	elbconn := meta.(*OutscaleClient).LBU

	l := d.Get("load_balancer_name").(string)
	p := d.Get("policy_name").(string)

	request := &lbu.DeleteLoadBalancerPolicyInput{
		LoadBalancerName: aws.String(l),
		PolicyName:       aws.String(p),
	}

	var err error
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, err = elbconn.API.DeleteLoadBalancerPolicy(request)

		if err != nil {
			if strings.Contains(fmt.Sprint(err), "Throttling") {
				return resource.RetryableError(
					fmt.Errorf("[WARN] Error deleting App stickiness policy, retrying: %s", err))
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("Error deleting App stickiness policy %s: %s", d.Id(), err)
	}
	return nil
}
