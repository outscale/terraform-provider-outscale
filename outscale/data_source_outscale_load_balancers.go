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

func dataSourceOutscaleLoadBalancers() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOutscaleLoadBalancersRead,

		Schema: map[string]*schema.Schema{
			"load_balancer_name": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"load_balancer_descriptions": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"load_balancer_name": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"availability_zones": &schema.Schema{
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"dns_name": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"health_check": &schema.Schema{
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"healthy_threshold": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"unhealthy_threshold": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"target": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"interval": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"timeout": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"instances": &schema.Schema{
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
						"listener_descriptions": &schema.Schema{
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
									"policy_names": &schema.Schema{
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
									"app_cookie_stickiness_policies": &schema.Schema{
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
									"lb_cookie_stickiness_policies": &schema.Schema{
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
									"other_policies": &schema.Schema{
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
						"security_groups": &schema.Schema{
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
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
						"subnets": &schema.Schema{
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

func dataSourceOutscaleLoadBalancersRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).LBU

	elbName, ok := d.GetOk("load_balancer_name")

	if !ok {
		return fmt.Errorf("load_balancer_name(s) must be provided")
	}

	describeElbOpts := &lbu.DescribeLoadBalancersInput{
		LoadBalancerNames: expandStringList(elbName.([]interface{})),
	}
	var resp *lbu.DescribeLoadBalancersOutput
	var describeResp *lbu.DescribeLoadBalancersResult
	var err error
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, err = conn.API.DescribeLoadBalancers(describeElbOpts)
		if err != nil {
			if strings.Contains(err.Error(), "Throttling:") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		describeResp = resp.DescribeLoadBalancersResult
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

		l["availability_zones"] = flattenStringList(v.AvailabilityZones)
		l["dns_name"] = aws.StringValue(v.DNSName)
		l["health_check"] = flattenHealthCheck(v.HealthCheck)

		l["instances"] = flattenInstances(v.Instances)
		l["listener_descriptions"] = flattenListeners(v.ListenerDescriptions)
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
			policies["app_cookie_stickiness_policies"] = app
			vc := make([]map[string]interface{}, len(v.Policies.LBCookieStickinessPolicies))
			for k, v := range v.Policies.LBCookieStickinessPolicies {
				a := make(map[string]interface{})
				a["policy_name"] = aws.StringValue(v.PolicyName)
				vc[k] = a
			}
			policies["lb_cookie_stickiness_policies"] = vc
			policies["other_policies"] = flattenStringList(v.Policies.OtherPolicies)
		}

		pl[0] = policies
		l["policies"] = pl
		l["scheme"] = aws.StringValue(v.Scheme)
		l["security_groups"] = flattenStringList(v.SecurityGroups)
		ssg := make(map[string]string)
		if v.SourceSecurityGroup != nil {
			ssg["group_name"] = aws.StringValue(v.SourceSecurityGroup.GroupName)
			ssg["owner_alias"] = aws.StringValue(v.SourceSecurityGroup.OwnerAlias)
		}
		l["source_security_group"] = ssg
		l["subnets"] = flattenStringList(v.Subnets)
		l["vpc_id"] = aws.StringValue(v.VPCId)

		lbs[k] = l
	}

	d.Set("request_id", resp.ResponseMetadata.RequestID)
	d.SetId(resource.UniqueId())

	return d.Set("load_balancer_descriptions", lbs)
}
