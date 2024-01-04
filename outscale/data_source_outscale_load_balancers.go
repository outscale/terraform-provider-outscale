package outscale

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/terraform-providers/terraform-provider-outscale/utils"
)

func attrLBSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"load_balancer_name": {
			Type:     schema.TypeList,
			Optional: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
		"load_balancer": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"load_balancer_name": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"subregion_names": {
						Type:     schema.TypeList,
						Computed: true,
						Elem:     &schema.Schema{Type: schema.TypeString},
					},
					"dns_name": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"access_log": {
						Type:     schema.TypeSet,
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
					"backend_vm_ids": {
						Type:     schema.TypeList,
						Computed: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"vm_id": {
									Type:     schema.TypeString,
									Computed: true,
								},
							},
						},
					},
					"listeners": {
						Type:     schema.TypeList,
						Required: true,
						Elem: &schema.Resource{
							Schema: lb_listener_schema(true),
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
					"load_balancer_type": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"tags": tagsListOAPISchema2(true),
					"security_groups": {
						Type:     schema.TypeList,
						Computed: true,
						Elem:     &schema.Schema{Type: schema.TypeString},
					},
					"source_security_group": {
						Type:     schema.TypeSet,
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
					"subnet_id": {
						Type:     schema.TypeList,
						Computed: true,
						Elem:     &schema.Schema{Type: schema.TypeString},
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
				},
			},
		},
		"request_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}

func dataSourceOutscaleOAPILoadBalancers() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOutscaleOAPILoadBalancersRead,

		Schema: getDataSourceSchemas(attrLBSchema()),
	}
}

func dataSourceOutscaleOAPILoadBalancersRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	resp, _, err := readLbs_(conn, d, schema.TypeList)
	if err != nil {
		return err
	}

	lbs_len := len(*resp.LoadBalancers)
	lbs_ret := make([]map[string]interface{}, lbs_len)

	lbs := *resp.LoadBalancers

	for k, v := range lbs {
		l := make(map[string]interface{})

		l["subregion_names"] = v.SubregionNames
		l["dns_name"] = *v.DnsName
		l["access_log"] = flattenOAPIAccessLog(v.AccessLog)
		l["health_check"] = flattenOAPIHealthCheck(v.HealthCheck)
		l["backend_vm_ids"] = utils.StringSlicePtrToInterfaceSlice(v.BackendVmIds)
		if v.Listeners != nil {
			l["listeners"] = flattenOAPIListeners(v.Listeners)
		} else {
			l["listeners"] = make([]interface{}, 0)
		}
		l["load_balancer_name"] = v.LoadBalancerName

		if v.ApplicationStickyCookiePolicies != nil {
			app := make([]map[string]interface{}, len(*v.ApplicationStickyCookiePolicies))
			for k, v := range *v.ApplicationStickyCookiePolicies {
				a := make(map[string]interface{})
				a["cookie_name"] = v.CookieName
				a["policy_name"] = v.PolicyName
				app[k] = a
			}
			l["application_sticky_cookie_policies"] = app
		} else {
			l["application_sticky_cookie_policies"] =
				make([]map[string]interface{}, 0)
		}

		if v.LoadBalancerStickyCookiePolicies != nil {
			vc := make([]map[string]interface{},
				len(*v.LoadBalancerStickyCookiePolicies))
			for k, v := range *v.LoadBalancerStickyCookiePolicies {
				a := make(map[string]interface{})
				a["policy_name"] = v.PolicyName
				vc[k] = a
			}
			l["load_balancer_sticky_cookie_policies"] = vc
		} else {
			l["load_balancer_sticky_cookie_policies"] =
				make([]map[string]interface{}, 0)
		}
		if v.Tags != nil {
			ta := make([]map[string]interface{}, len(*v.Tags))
			for k1, v1 := range *v.Tags {
				t := make(map[string]interface{})
				t["key"] = v1.Key
				t["value"] = v1.Value
				ta[k1] = t
			}
			l["tags"] = ta
		}

		l["load_balancer_type"] = v.LoadBalancerType
		l["security_groups"] = utils.StringSlicePtrToInterfaceSlice(v.SecurityGroups)

		ssg := make([]map[string]interface{}, 0)
		if v.SourceSecurityGroup != nil {
			l["source_security_group"] = flattenSource_sg(v.SourceSecurityGroup)
		} else {
			l["source_security_group"] = ssg

		}
		l["subnet_id"] = utils.StringSlicePtrToInterfaceSlice(v.Subnets)
		l["public_ip"] = v.PublicIp
		l["secured_cookies"] = v.SecuredCookies
		l["net_id"] = v.NetId

		lbs_ret[k] = l
	}

	err = d.Set("load_balancer", lbs_ret)
	if err != nil {
		return err
	}
	d.SetId(resource.UniqueId())

	return nil
}
