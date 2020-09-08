package outscale

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/antihax/optional"
	oscgo "github.com/marinsalinas/osc-sdk-go"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceOutscaleOAPILoadBalancers() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOutscaleOAPILoadBalancersRead,

		Schema: map[string]*schema.Schema{
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
						"health_check": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"healthy_threshold": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"unhealthy_threshold": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"checked_vm": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"check_interval": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"timeout": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"backend_vm_id": {
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
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"listener": {
										Type:     schema.TypeMap,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"backend_port": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"backend_protocol": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"load_balancer_port": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"load_balancer_protocol": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"server_certificate_id": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
									"policy_name": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"policies": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"application_sticky_cookie_policy": {
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
									"load_balancer_sticky_cookie_policy": {
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
								},
							},
						},
						"load_balancer_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"security_groups_member": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"firewall_rules_set_name": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"firewall_rules_set_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"account_alias": {
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
		},
	}
}

func dataSourceOutscaleOAPILoadBalancersRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	eName, ok := d.GetOk("load_balancer_name")

	if !ok {
		return fmt.Errorf("load_balancer_name(s) must be provided")
	}

	elbName := eName.(string)

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

	lbs_len := len(*resp.LoadBalancers)
	lbs_ret := make([]map[string]interface{}, lbs_len)

	lbs := *resp.LoadBalancers
	if len(lbs) < 1 {
		return fmt.Errorf("Unable to find LBU: %s", elbName)
	}

	for k, v := range lbs {
		l := make(map[string]interface{})

		l["subregion_names"] = v.SubregionNames
		l["public_dns_name"] = v.DnsName
		l["health_check"] = v.HealthCheck
		l["backend_vm_id"] = v.BackendVmIds
		l["listeners"] = flattenOAPIListeners(v.Listeners)
		l["load_balancer_name"] = elbName

		policies := make(map[string]interface{})
		pl := make([]map[string]interface{}, 1)
		if v.ApplicationStickyCookiePolicies != nil {
			app := make([]map[string]interface{}, len(*v.ApplicationStickyCookiePolicies))
			for k, v := range *v.ApplicationStickyCookiePolicies {
				a := make(map[string]interface{})
				a["cookie_name"] = v.CookieName
				a["policy_name"] = v.PolicyName
				app[k] = a
			}
			policies["application_sticky_cookie_policy"] = app
			vc := make([]map[string]interface{},
				len(*v.LoadBalancerStickyCookiePolicies))
			for k, v := range *v.LoadBalancerStickyCookiePolicies {
				a := make(map[string]interface{})
				a["policy_name"] = v.PolicyName
				vc[k] = a
			}
			policies["load_balancer_sticky_cookie_policy"] = vc
		}

		pl[0] = policies
		l["policies"] = pl
		l["load_balancer_type"] = v.LoadBalancerType
		l["security_groups_member"] = flattenStringList(v.SecurityGroups)
		ssg := make(map[string]string)
		if v.SourceSecurityGroup != nil {
			ssg["security_group_account_id"] = *v.SourceSecurityGroup.SecurityGroupAccountId
			ssg["security_group_name"] = *v.SourceSecurityGroup.SecurityGroupName
		}
		l["firewall_rules_set_name"] = ssg
		l["subnet_id"] = flattenStringList(v.Subnets)
		l["net_id"] = v.NetId

		lbs_ret[k] = l
	}

	d.Set("request_id", resp.ResponseContext.RequestId)
	d.SetId(resource.UniqueId())

	return d.Set("load_balancer", lbs)
}
