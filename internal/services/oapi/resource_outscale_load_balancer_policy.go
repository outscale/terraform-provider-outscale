package oapi

import (
	"context"
	"fmt"
	"regexp"

	"github.com/outscale/goutils/sdk/ptr"
	"github.com/outscale/osc-sdk-go/v3/pkg/options"
	"github.com/outscale/osc-sdk-go/v3/pkg/osc"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/utils"
	"github.com/samber/lo"
	"github.com/spf13/cast"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceOutscaleAppCookieStickinessPolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceOutscaleAppCookieStickinessPolicyCreate,
		ReadContext:   ResourceOutscaleAppCookieStickinessPolicyRead,
		DeleteContext: ResourceOutscaleAppCookieStickinessPolicyDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(CreateDefaultTimeout),
			Read:   schema.DefaultTimeout(ReadDefaultTimeout),
			Delete: schema.DefaultTimeout(DeleteDefaultTimeout),
		},

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
			"tags": TagsSchemaComputedSDK(),

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

func ResourceOutscaleAppCookieStickinessPolicyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.OutscaleClient).OSC

	timeout := d.Timeout(schema.TimeoutCreate)

	l := d.Get("load_balancer_name")
	pn := d.Get("policy_name")
	pt := d.Get("policy_type")

	cep, cepok := d.GetOk("cookie_expiration_period")
	v, cnok := d.GetOk("cookie_name")

	if cepok && pt.(string) == "app" {
		return diag.Errorf("if you want define \"cookie_expiration_period\", use policy_type = \"load_balancer\"")
	}
	if cnok && pt.(string) == "load_balancer" {
		return diag.Errorf("if you want define \"cookie_name\", use policy_type = \"app\"")
	}

	vs := v.(string)
	req := osc.CreateLoadBalancerPolicyRequest{
		LoadBalancerName: l.(string),
		PolicyName:       pn.(string),
		PolicyType:       pt.(string),
	}
	if cnok {
		req.CookieName = &vs
	}
	if cepok {
		req.CookieExpirationPeriod = new(cast.ToInt(cep))
	}
	_, err := client.CreateLoadBalancerPolicy(ctx, req, options.WithRetryTimeout(timeout))
	if err != nil {
		return diag.Errorf("error creating appcookiestickinesspolicy: %s", err)
	}
	d.SetId(id.UniqueId())
	d.Set("load_balancer_name", l.(string))
	d.Set("policy_name", pn.(string))
	d.Set("policy_type", pt.(string))

	return ResourceOutscaleAppCookieStickinessPolicyRead(ctx, d, meta)
}

func ResourceOutscaleAppCookieStickinessPolicyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.OutscaleClient).OSC
	timeout := d.Timeout(schema.TimeoutRead)

	lbuName := d.Get("load_balancer_name").(string)
	policyName := d.Get("policy_name").(string)
	lb, _, err := readResourceLb(ctx, client, lbuName, timeout)
	if err != nil {
		return diag.FromErr(err)
	}
	if lb == nil || (lb.ApplicationStickyCookiePolicies == nil && lb.LoadBalancerStickyCookiePolicies == nil) {
		d.SetId("")
		return nil
	}
	_, foundAppPolicy := lo.Find(lb.ApplicationStickyCookiePolicies, func(v osc.ApplicationStickyCookiePolicy) bool {
		return ptr.From(v.PolicyName) == policyName
	})
	_, foundLbuPolicy := lo.Find(lb.LoadBalancerStickyCookiePolicies, func(v osc.LoadBalancerStickyCookiePolicy) bool {
		return ptr.From(v.PolicyName) == policyName
	})
	if !foundAppPolicy && !foundLbuPolicy {
		d.SetId("")
		return nil
	}

	d.Set("access_log", flattenOAPIAccessLog(&lb.AccessLog))
	d.Set("listeners", flattenOAPIListeners(&lb.Listeners))
	d.Set("subregion_names", utils.StringSlicePtrToInterfaceSlice(&lb.SubregionNames))
	d.Set("load_balancer_type", lb.LoadBalancerType)
	if lb.SecurityGroups != nil {
		d.Set("security_groups", utils.StringSlicePtrToInterfaceSlice(&lb.SecurityGroups))
	} else {
		d.Set("security_groups", make([]map[string]interface{}, 0))
	}
	d.Set("dns_name", lb.DnsName)
	d.Set("subnets", utils.StringSlicePtrToInterfaceSlice(&lb.Subnets))
	d.Set("health_check", flattenOAPIHealthCheck(&lb.HealthCheck))
	d.Set("backend_vm_ids", utils.StringSlicePtrToInterfaceSlice(&lb.BackendVmIds))
	if lb.Tags != nil {
		ta := make([]map[string]interface{}, len(lb.Tags))
		for k1, v1 := range lb.Tags {
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
			len(lb.ApplicationStickyCookiePolicies))
		for k, v := range lb.ApplicationStickyCookiePolicies {
			a := make(map[string]interface{})
			a["cookie_name"] = v.CookieName
			a["policy_name"] = v.PolicyName
			if ptr.From(v.PolicyName) == policyName {
				d.Set("cookie_name", ptr.From(v.CookieName))
			}
			app[k] = a
		}
		d.Set("application_sticky_cookie_policies", app)
	}
	if lb.LoadBalancerStickyCookiePolicies != nil {
		lbc := make([]map[string]interface{},
			len(lb.LoadBalancerStickyCookiePolicies))
		for k, v := range lb.LoadBalancerStickyCookiePolicies {
			a := make(map[string]interface{})
			a["policy_name"] = v.PolicyName
			if ptr.From(v.PolicyName) == policyName {
				d.Set("cookie_expiration_period", cast.ToInt32(v.CookieExpirationPeriod))
			}
			lbc[k] = a
		}
		d.Set("load_balancer_sticky_cookie_policies", lbc)
	}

	d.Set("source_security_group", flattenSource_sg(&lb.SourceSecurityGroup))
	d.Set("public_ip", ptr.From(lb.PublicIp))
	d.Set("secured_cookies", lb.SecuredCookies)
	d.Set("net_id", ptr.From(lb.NetId))

	return nil
}

func ResourceOutscaleAppCookieStickinessPolicyDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	elbclient := meta.(*client.OutscaleClient).OSC

	timeout := d.Timeout(schema.TimeoutDelete)

	l := d.Get("load_balancer_name").(string)
	p := d.Get("policy_name").(string)

	request := osc.DeleteLoadBalancerPolicyRequest{
		LoadBalancerName: l,
		PolicyName:       p,
	}

	_, err := elbclient.DeleteLoadBalancerPolicy(ctx, request, options.WithRetryTimeout(timeout))
	if err != nil {
		return diag.Errorf("error deleting app stickiness policy %s: %s", d.Id(), err)
	}
	return nil
}
