package outscale

import (
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-outscale/osc/lbu"
)

func dataSourceOutscaleOAPILoadBalancers() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOutscaleOAPILoadBalancersRead,

		Schema: map[string]*schema.Schema{
			"load_balancer_name": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"load_balancer": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"load_balancer_name": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"sub_region_name": &schema.Schema{
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"public_dns_name": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"health_check": &schema.Schema{
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"healthy_threshold": &schema.Schema{
										Type:     schema.TypeInt,
										Computed: true,
									},
									"unhealthy_threshold": &schema.Schema{
										Type:     schema.TypeInt,
										Computed: true,
									},
									"checked_vm": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"check_interval": &schema.Schema{
										Type:     schema.TypeInt,
										Computed: true,
									},
									"timeout": &schema.Schema{
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},
						"backend_vm_id": &schema.Schema{
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"vm_id": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"listeners": &schema.Schema{
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"listener": &schema.Schema{
										Type:     schema.TypeMap,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"backend_port": &schema.Schema{
													Type:     schema.TypeInt,
													Computed: true,
												},
												"backend_protocol": &schema.Schema{
													Type:     schema.TypeString,
													Computed: true,
												},
												"load_balancer_port": &schema.Schema{
													Type:     schema.TypeInt,
													Computed: true,
												},
												"load_balancer_protocol": &schema.Schema{
													Type:     schema.TypeString,
													Computed: true,
												},
												"server_certificate_id": &schema.Schema{
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
									"policy_name": &schema.Schema{
										Type:     schema.TypeList,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"policies": &schema.Schema{
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"application_sticky_cookie_policy": &schema.Schema{
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"cookie_name": &schema.Schema{
													Type:     schema.TypeString,
													Computed: true,
												},
												"policy_name": &schema.Schema{
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
									"load_balancer_sticky_cookie_policy": &schema.Schema{
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"policy_name": &schema.Schema{
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
									"other_policy": &schema.Schema{
										Type:     schema.TypeList,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"load_balancer_type": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"security_groups_member": &schema.Schema{
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"firewall_rules_set_name": &schema.Schema{
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"firewall_rules_set_name": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"account_alias": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"subnet_id": &schema.Schema{
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"lin_id": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"request_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceOutscaleOAPILoadBalancersRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).LBU

	elbName, ok := d.GetOk("load_balancer_name")

	if !ok {
		return fmt.Errorf("load_balancer_name(s) must be provided")
	}

	describeElbOpts := &lbu.DescribeLoadBalancersInput{
		LoadBalancerNames: expandStringList(elbName.([]interface{})),
	}

	var describeResp *lbu.DescribeLoadBalancersOutput
	var err error
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		describeResp, err = conn.API.DescribeLoadBalancers(describeElbOpts)

		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded:") {
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
	if len(describeResp.LoadBalancerDescriptions) < 1 {
		return fmt.Errorf("Unable to find ELB: %#v", describeResp.LoadBalancerDescriptions)
	}

	lbs := make([]map[string]interface{}, len(describeResp.LoadBalancerDescriptions))

	for k, v := range describeResp.LoadBalancerDescriptions {
		l := make(map[string]interface{})

		l["sub_region_name"] = flattenStringList(v.AvailabilityZones)
		l["public_dns_name"] = aws.StringValue(v.DNSName)
		if *v.HealthCheck.Target != "" {
			l["health_check"] = flattenHealthCheck(v.HealthCheck)
		} else {
			l["health_check"] = make(map[string]interface{})
		}
		l["backend_vm_id"] = flattenOAPIInstances(v.Instances)
		l["listeners"] = flattenOAPIListeners(v.ListenerDescriptions)
		l["load_balancer_name"] = aws.StringValue(v.LoadBalancerName)

		policies := make(map[string]interface{})
		pl := make([]map[string]interface{}, 1)
		if v.Policies != nil {
			app := make([]map[string]interface{}, len(v.Policies.AppCookieStickinessPolicies))
			for k, v := range v.Policies.AppCookieStickinessPolicies {
				a := make(map[string]interface{})
				a["cookie_name"] = aws.StringValue(v.CookieName)
				a["policy_name"] = aws.StringValue(v.PolicyName)
				app[k] = a
			}
			policies["application_sticky_cookie_policy"] = app
			vc := make([]map[string]interface{}, len(v.Policies.LBCookieStickinessPolicies))
			for k, v := range v.Policies.LBCookieStickinessPolicies {
				a := make(map[string]interface{})
				a["policy_name"] = aws.StringValue(v.PolicyName)
				vc[k] = a
			}
			policies["load_balancer_sticky_cookie_policy"] = vc
			policies["other_policy"] = flattenStringList(v.Policies.OtherPolicies)
		}

		pl[0] = policies
		l["policies"] = pl
		l["load_balancer_type"] = aws.StringValue(v.Scheme)
		l["security_groups_member"] = flattenStringList(v.SecurityGroups)
		ssg := make(map[string]string)
		if v.SourceSecurityGroup != nil {
			ssg["firewall_rules_set_name"] = aws.StringValue(v.SourceSecurityGroup.GroupName)
			ssg["account_alias"] = aws.StringValue(v.SourceSecurityGroup.OwnerAlias)
		}
		l["firewall_rules_set_name"] = ssg
		l["subnet_id"] = flattenStringList(v.Subnets)
		l["lin_id"] = aws.StringValue(v.VPCId)

		lbs[k] = l
	}

	// d.Set("request_id", describeResp.RequestID)
	d.SetId(resource.UniqueId())

	return d.Set("load_balancer", lbs)
}
