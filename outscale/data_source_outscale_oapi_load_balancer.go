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

func dataSourceOutscaleOAPILoadBalancer() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOutscaleOAPILoadBalancerRead,

		Schema: map[string]*schema.Schema{
			"sub_region_name": &schema.Schema{
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

			"public_dns_name": &schema.Schema{
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
						"checked_vm": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"check_interval": &schema.Schema{
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
			"lin_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
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
			"request_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceOutscaleOAPILoadBalancerRead(d *schema.ResourceData, meta interface{}) error {
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
	var resp *lbu.DescribeLoadBalancersOutput
	var describeResp *lbu.DescribeLoadBalancersResult
	var err error
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, err = conn.API.DescribeLoadBalancers(describeElbOpts)
		if err != nil {
			if strings.Contains(fmt.Sprint(err), "Throttling:") {
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

	if describeResp.LoadBalancerDescriptions == nil {
		return fmt.Errorf("NO ELB FOUND")
	}

	if len(describeResp.LoadBalancerDescriptions) != 1 {
		return fmt.Errorf("Unable to find ELB: %#v", describeResp.LoadBalancerDescriptions)
	}

	lb := describeResp.LoadBalancerDescriptions[0]

	d.Set("sub_region_name", flattenStringList(lb.AvailabilityZones))
	d.Set("public_dns_name", aws.StringValue(lb.DNSName))
	d.Set("health_check", flattenOAPIHealthCheck(lb.HealthCheck))

	if lb.Instances != nil {
		d.Set("backend_vm_id", flattenInstances(lb.Instances))
	} else {
		d.Set("backend_vm_id", make([]map[string]interface{}, 0))
	}
	if lb.ListenerDescriptions != nil {
		if err := d.Set("listeners", flattenOAPIListeners(lb.ListenerDescriptions)); err != nil {
			return err
		}
	} else {
		if err := d.Set("listeners", make([]map[string]interface{}, 0)); err != nil {
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
		policies["application_sticky_cookie_policy"] = app
		lbc := make([]map[string]interface{}, len(lb.Policies.LBCookieStickinessPolicies))
		for k, v := range lb.Policies.LBCookieStickinessPolicies {
			a := make(map[string]interface{})
			a["policy_name"] = aws.StringValue(v.PolicyName)
			lbc[k] = a
		}
		policies["load_balancer_sticky_cookie_policy"] = lbc
		policies["other_policy"] = flattenStringList(lb.Policies.OtherPolicies)
	} else {
		lbc := make([]map[string]interface{}, 0)
		policies["load_balancer_sticky_cookie_policy"] = lbc
		policies["other_policy"] = lbc
	}
	d.Set("policies", policies)
	d.Set("scheme", aws.StringValue(lb.Scheme))
	if lb.SecurityGroups != nil {
		d.Set("security_groups_member", flattenStringList(lb.SecurityGroups))
	} else {
		d.Set("security_groups_member", make([]map[string]interface{}, 0))
	}
	ssg := make(map[string]string)
	if lb.SourceSecurityGroup != nil {
		ssg["firewall_rules_set_name"] = aws.StringValue(lb.SourceSecurityGroup.GroupName)
		ssg["account_alias"] = aws.StringValue(lb.SourceSecurityGroup.OwnerAlias)
	}
	d.Set("firewall_rules_set_name", ssg)
	d.Set("subnets_member", flattenStringList(lb.Subnets))
	d.Set("lin_id", aws.StringValue(lb.VPCId))
	d.Set("request_id", resp.ResponseMetadata.RequestID)
	d.SetId(*lb.LoadBalancerName)

	return nil
}
