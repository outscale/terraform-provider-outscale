package outscale

import (
	"context"
	"fmt"
	"regexp"
	"time"

	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/spf13/cast"
	"github.com/terraform-providers/terraform-provider-outscale/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
			"access_log": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"is_enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"osu_bucket_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"osu_bucket_prefix": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"publication_interval": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
			"health_check": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
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
						"port": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"protocol": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"timeout": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
			"application_sticky_cookie_policies": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cookie_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"policy_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"load_balancer_sticky_cookie_policies": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"policy_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"listeners": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: lb_listener_schema(true),
				},
			},
			"source_security_group": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"security_group_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"security_group_account_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"public_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"secured_cookies": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"net_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"backend_vm_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				ForceNew: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"subregion_names": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"load_balancer_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"security_groups": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"subnets": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"tags": tagsListOAPISchema2(true),

			"dns_name": {
				Type:     schema.TypeString,
				Computed: true,
				ForceNew: true,
			},
			"policy_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, es []error) {
					value := v.(string)
					if !regexp.MustCompile(`^app|load_balancer$`).MatchString(value) {
						es = append(es, fmt.Errorf(
							"only \"app\" or \"load_balancer\" allowed in %q", k))
					}
					return
				},
			},
			"load_balancer_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"cookie_expiration_period": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"cookie_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
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

	l := d.Get("load_balancer_name")
	pn := d.Get("policy_name")
	pt := d.Get("policy_type")

	cep, cepok := d.GetOk("cookie_expiration_period")
	v, cnok := d.GetOk("cookie_name")

	if cepok && pt.(string) == "app" {
		return fmt.Errorf("if you want define \"cookie_expiration_period\", use policy_type = \"load_balancer\"")
	}
	if cnok && pt.(string) == "load_balancer" {
		return fmt.Errorf("if you want define \"cookie_name\", use policy_type = \"app\"")
	}

	vs := v.(string)
	req := oscgo.CreateLoadBalancerPolicyRequest{
		LoadBalancerName: l.(string),
		PolicyName:       pn.(string),
		PolicyType:       pt.(string),
	}
	if cnok {
		req.CookieName = &vs
	}
	if cepok {
		req.SetCookieExpirationPeriod(cast.ToInt32(cep))
	}
	var err error
	var resp oscgo.CreateLoadBalancerPolicyResponse
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		rp, httpResp, err := conn.LoadBalancerPolicyApi.
			CreateLoadBalancerPolicy(
				context.Background()).
			CreateLoadBalancerPolicyRequest(req).Execute()

		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		resp = rp
		return nil
	})

	if err != nil {
		return fmt.Errorf("Error creating AppCookieStickinessPolicy: %s", err)
	}

	if resp.LoadBalancer != nil {
		lb := resp.LoadBalancer
		d.Set("access_log", flattenOAPIAccessLog(lb.AccessLog))
		d.Set("listeners", flattenOAPIListeners(lb.Listeners))
		d.Set("subregion_names", utils.StringSlicePtrToInterfaceSlice(lb.SubregionNames))
		d.Set("load_balancer_type", lb.LoadBalancerType)
		if lb.SecurityGroups != nil {
			d.Set("security_groups", utils.StringSlicePtrToInterfaceSlice(lb.SecurityGroups))
		} else {
			d.Set("security_groups", make([]map[string]interface{}, 0))
		}
		d.Set("dns_name", lb.DnsName)
		d.Set("subnets", utils.StringSlicePtrToInterfaceSlice(lb.Subnets))
		d.Set("health_check", flattenOAPIHealthCheck(lb.HealthCheck))
		d.Set("backend_vm_ids", utils.StringSlicePtrToInterfaceSlice(lb.BackendVmIds))
		if lb.Tags != nil {
			ta := make([]map[string]interface{}, len(*lb.Tags))
			for k1, v1 := range *lb.Tags {
				t := make(map[string]interface{})
				t["key"] = v1.Key
				t["value"] = v1.Value
				ta[k1] = t
			}
			d.Set("tags", ta)
		} else {
			d.Set("tags", make([]map[string]interface{}, 0))
		}
		if lb.ApplicationStickyCookiePolicies != nil {
			app := make([]map[string]interface{},
				len(*lb.ApplicationStickyCookiePolicies))
			for k, v := range *lb.ApplicationStickyCookiePolicies {
				a := make(map[string]interface{})
				a["cookie_name"] = v.CookieName
				a["policy_name"] = v.PolicyName
				if v.GetPolicyName() == pn.(string) {
					d.Set("cookie_name", v.GetCookieName())
				}
				app[k] = a
			}
			d.Set("application_sticky_cookie_policies", app)
		}
		if lb.LoadBalancerStickyCookiePolicies != nil {
			lbc := make([]map[string]interface{},
				len(*lb.LoadBalancerStickyCookiePolicies))
			for k, v := range *lb.LoadBalancerStickyCookiePolicies {
				a := make(map[string]interface{})
				a["policy_name"] = v.PolicyName
				if v.GetPolicyName() == pn.(string) {
					d.Set("cookie_expiration_period", cast.ToInt32(v.CookieExpirationPeriod))
				}
				lbc[k] = a
			}
			d.Set("load_balancer_sticky_cookie_policies", lbc)
		}

		if lb.SourceSecurityGroup != nil {
			d.Set("source_security_group", flattenSource_sg(lb.SourceSecurityGroup))
		}
		d.Set("public_ip", lb.PublicIp)
		d.Set("secured_cookies", lb.SecuredCookies)
		d.Set("net_id", lb.NetId)
	}

	d.SetId(resource.UniqueId())
	d.Set("load_balancer_name", l.(string))
	d.Set("policy_name", pn.(string))
	d.Set("policy_type", pt.(string))
	return nil
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

	var err error
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, httpResp, err := elbconn.LoadBalancerPolicyApi.
			DeleteLoadBalancerPolicy(
				context.Background()).
			DeleteLoadBalancerPolicyRequest(request).Execute()

		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("Error deleting App stickiness policy %s: %s", d.Id(), err)
	}
	return nil
}
