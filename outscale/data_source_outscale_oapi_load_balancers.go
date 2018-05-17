package outscale

import (
	"github.com/hashicorp/terraform/helper/schema"
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
			"load_balancer_descriptions_member": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"load_balancer_name": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"availability_zones_member": &schema.Schema{
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"dns_name": &schema.Schema{
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
									"target": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"interval": &schema.Schema{
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
						"instances_member": &schema.Schema{
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"instance_id": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"listener_descriptions_member": &schema.Schema{
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"listener": &schema.Schema{
										Type:     schema.TypeMap,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"instance_port": &schema.Schema{
													Type:     schema.TypeInt,
													Computed: true,
												},
												"instance_protocol": &schema.Schema{
													Type:     schema.TypeString,
													Computed: true,
												},
												"load_balancer_port": &schema.Schema{
													Type:     schema.TypeInt,
													Computed: true,
												},
												"protocol": &schema.Schema{
													Type:     schema.TypeString,
													Computed: true,
												},
												"ssl_certificate_id": &schema.Schema{
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
									"policy_names_member": &schema.Schema{
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
									"app_cookie_stickiness_policies_member": &schema.Schema{
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
									"lb_cookie_stickiness_policies_member": &schema.Schema{
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
									"other_policies_member": &schema.Schema{
										Type:     schema.TypeList,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"scheme": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"source_security_group": &schema.Schema{
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"group_name": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"owner_alias": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"subnets_member": &schema.Schema{
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"vpc_id": &schema.Schema{
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
	// conn := meta.(*OutscaleClient).LBU

	// elbName, ok := d.GetOk("load_balancer_name")

	// if !ok {
	// 	return fmt.Errorf("please provide the required attribute load_balancer_name")
	// }

	// describeElbOpts := &lbu.DescribeLoadBalancersInput{
	// 	LoadBalancerNames: []*string{aws.String(elbName.(string))},
	// }

	// var describeResp *lbu.DescribeLoadBalancersOutput
	// var err error
	// err = resource.Retry(5*time.Minute, func() *resource.RetryError {
	// 	describeResp, err = conn.API.DescribeLoadBalancers(describeElbOpts)

	// 	if err != nil {
	// 		if strings.Contains(err.Error(), "RequestLimitExceeded:") {
	// 			return resource.RetryableError(err)
	// 		}
	// 		return resource.NonRetryableError(err)
	// 	}
	// 	return nil
	// })

	// if err != nil {
	// 	if isLoadBalancerNotFound(err) {
	// 		d.SetId("")
	// 		return nil
	// 	}

	// 	return fmt.Errorf("Error retrieving ELB: %s", err)
	// }
	// if len(describeResp.LoadBalancerDescriptions) != 1 {
	// 	return fmt.Errorf("Unable to find ELB: %#v", describeResp.LoadBalancerDescriptions)
	// }

	// lb := describeResp.LoadBalancerDescriptions

	// lbs := make([]map[string]interface{}, len(lb))

	// for k, v := range lb {
	// 	l := make(map[string]interface{})

	// 	l["availability_zones_member"] = flattenStringList(v.AvailabilityZones)
	// 	l["dns_name"] = aws.StringValue(v.DNSName)
	// 	if *v.HealthCheck.Target != "" {
	// 		l["health_check"] = flattenHealthCheck(v.HealthCheck)
	// 	} else {
	// 		l["health_check"] = make(map[string]interface{})
	// 	}
	// 	l["instances_member"] = flattenInstances(v.Instances)
	// 	l["listener_descriptions_member"] = flattenListeners(v.ListenerDescriptions)
	// 	l["load_balancer_name"] = aws.StringValue(v.LoadBalancerName)

	// 	policies := make(map[string]interface{})
	// 	if v.Policies != nil {
	// 		app := make([]map[string]interface{}, len(v.Policies.AppCookieStickinessPolicies))
	// 		for k, v := range v.Policies.AppCookieStickinessPolicies {
	// 			a := make(map[string]interface{})
	// 			a["cookie_name"] = aws.StringValue(v.CookieName)
	// 			a["policy_name"] = aws.StringValue(v.PolicyName)
	// 			app[k] = a
	// 		}
	// 		policies["app_cookie_stickiness_policies_member"] = app
	// 		vc := make([]map[string]interface{}, len(v.Policies.LBCookieStickinessPolicies))
	// 		for k, v := range v.Policies.LBCookieStickinessPolicies {
	// 			a := make(map[string]interface{})
	// 			a["policy_name"] = aws.StringValue(v.PolicyName)
	// 			vc[k] = a
	// 		}
	// 		policies["lb_cookie_stickiness_policies_member"] = vc
	// 		policies["other_policies_member"] = flattenStringList(v.Policies.OtherPolicies)
	// 	}

	// 	l["policies"] = policies
	// 	l["scheme"] = aws.StringValue(v.Scheme)
	// 	l["security_groups_member"] = flattenStringList(v.SecurityGroups)
	// 	ssg := make(map[string]string)
	// 	if v.SourceSecurityGroup != nil {
	// 		ssg["group_name"] = aws.StringValue(v.SourceSecurityGroup.GroupName)
	// 		ssg["owner_alias"] = aws.StringValue(v.SourceSecurityGroup.OwnerAlias)
	// 	}
	// 	l["source_security_group"] = ssg
	// 	l["subnets_member"] = flattenStringList(v.Subnets)
	// 	l["vpc_id"] = aws.StringValue(v.VPCId)

	// 	lbs[k] = l
	// }

	// d.Set("load_balancer_descriptions_member", lbs)
	// d.Set("request_id", describeResp.RequestID)
	// d.SetId(resource.UniqueId())

	return nil
}
