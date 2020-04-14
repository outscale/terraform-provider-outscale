package outscale

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/antihax/optional"
	oscgo "github.com/marinsalinas/osc-sdk-go"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceOutscaleAppCookieStickinessPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceOutscaleAppCookieStickinessPolicyCreate,
		Read:   resourceOutscaleAppCookieStickinessPolicyRead,
		Delete: resourceOutscaleAppCookieStickinessPolicyDelete,

		Schema: map[string]*schema.Schema{
			"policy_name": {
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

			"load_balancer_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"cookie_name": {
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

func resourceOutscaleAppCookieStickinessPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	l, lok := d.GetOk("load_balancer_name")
	p, pok := d.GetOk("policy_name")
	v, ok := d.GetOk("cookie_name")

	if !lok && !pok && !ok {
		return fmt.Errorf("please provide the required attributes load_balancer_name, policy_name and cookie_name")
	}

	vs := v.(string)
	req := oscgo.CreateLoadBalancerPolicyRequest{
		LoadBalancerName: l.(string),
		PolicyName:       p.(string),
		CookieName:       &vs,
	}
	acspOpts := oscgo.CreateLoadBalancerPolicyOpts{
		optional.NewInterface(req),
	}

	var err error
	var resp oscgo.CreateLoadBalancerPolicyResponse
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, _, err = conn.LoadBalancerPolicyApi.
			CreateLoadBalancerPolicy(
				context.Background(),
				&acspOpts)

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
	if resp.ResponseContext != nil {
		reqID = *resp.ResponseContext.RequestId
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
	elbconn := meta.(*OutscaleClient).OSCAPI

	l := d.Get("load_balancer_name").(string)
	p := d.Get("policy_name").(string)

	request := oscgo.DeleteLoadBalancerPolicyRequest{
		LoadBalancerName: l,
		PolicyName:       p,
	}

	opts := &oscgo.DeleteLoadBalancerPolicyOpts{
		optional.NewInterface(request),
	}

	var err error
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, _, err = elbconn.LoadBalancerPolicyApi.
			DeleteLoadBalancerPolicy(
				context.Background(), opts)

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
