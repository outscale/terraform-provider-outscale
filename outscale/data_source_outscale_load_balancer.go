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

func dataSourceOutscaleLoadBalancer() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOutscaleLoadBalancerRead,

		Schema: map[string]*schema.Schema{
			"availability_zones_member": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"load_balancer_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"scheme": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"security_groups_member": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"subnets_member": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"tag": tagsSchema(),

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
			"vpc_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
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
			"request_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceOutscaleLoadBalancerRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).LBU
	ename, ok := d.GetOk("load_balancer_name")

	if !ok {
		return fmt.Errorf("please provide the name of the load balancer")
	}

	elbName := ename.(string)

	// Retrieve the ELB properties for updating the state
	describeElbOpts := &lbu.DescribeLoadBalancersInput{
		LoadBalancerNames: []*string{aws.String(elbName)},
	}

	var describeResp *lbu.DescribeLoadBalancersOutput
	var err error
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		describeResp, err = conn.API.DescribeLoadBalancers(describeElbOpts)

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

	if describeResp.LoadBalancerDescriptions == nil {
		return fmt.Errorf("NO ELB FOUND")
	}

	if len(describeResp.LoadBalancerDescriptions) != 1 {
		return fmt.Errorf("Unable to find ELB: %#v", describeResp.LoadBalancerDescriptions)
	}

	lb := describeResp.LoadBalancerDescriptions[0]

	d.Set("availability_zones_member", flattenStringList(lb.AvailabilityZones))
	d.Set("dns_name", aws.StringValue(lb.DNSName))
	if *lb.HealthCheck.Target != "" {
		d.Set("health_check", flattenHealthCheck(lb.HealthCheck))
	} else {
		d.Set("health_check", make(map[string]interface{}))
	}
	if lb.Instances != nil {
		d.Set("instances_member", flattenInstances(lb.Instances))
	} else {
		d.Set("instances_member", make([]map[string]interface{}, 0))
	}
	if lb.ListenerDescriptions != nil {
		if err := d.Set("listener_descriptions_member", flattenListeners(lb.ListenerDescriptions)); err != nil {
			return err
		}
	} else {
		if err := d.Set("listener_descriptions_member", make([]map[string]interface{}, 0)); err != nil {
			return err
		}
	}
	d.Set("load_balancer_name", aws.StringValue(lb.LoadBalancerName))

	policies := make(map[string]interface{})
	if lb.Policies != nil {
		app := make([]map[string]interface{}, len(lb.Policies.AppCookieStickinessPolicies))
		for k, v := range lb.Policies.AppCookieStickinessPolicies {
			a := make(map[string]interface{})
			a["cookie_name"] = aws.StringValue(v.CookieName)
			a["policy_name"] = aws.StringValue(v.PolicyName)
			app[k] = a
		}
		policies["app_cookie_stickiness_policies_member"] = app
		lbc := make([]map[string]interface{}, len(lb.Policies.LBCookieStickinessPolicies))
		for k, v := range lb.Policies.LBCookieStickinessPolicies {
			a := make(map[string]interface{})
			a["policy_name"] = aws.StringValue(v.PolicyName)
			lbc[k] = a
		}
		policies["lb_cookie_stickiness_policies_member"] = lbc
		policies["other_policies_member"] = flattenStringList(lb.Policies.OtherPolicies)
	} else {
		lbc := make([]map[string]interface{}, 0)
		policies["lb_cookie_stickiness_policies_member"] = lbc
		policies["other_policies_member"] = lbc
	}
	pl := make([]map[string]interface{}, 1)
	pl[0] = policies
	d.Set("policies", pl)
	d.Set("scheme", aws.StringValue(lb.Scheme))
	if lb.SecurityGroups != nil {
		d.Set("security_groups_member", flattenStringList(lb.SecurityGroups))
	} else {
		d.Set("security_groups_member", make([]map[string]interface{}, 0))
	}
	ssg := make(map[string]string)
	if lb.SourceSecurityGroup != nil {
		ssg["group_name"] = aws.StringValue(lb.SourceSecurityGroup.GroupName)
		ssg["owner_alias"] = aws.StringValue(lb.SourceSecurityGroup.OwnerAlias)
	}
	d.Set("source_security_group", ssg)
	d.Set("subnets_member", flattenStringList(lb.Subnets))
	d.Set("vpc_id", aws.StringValue(lb.VPCId))
	// d.Set("request_id", resp.ResponseMetadata.RequestID)
	d.SetId(*lb.LoadBalancerName)

	return nil
}
