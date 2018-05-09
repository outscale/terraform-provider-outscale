package outscale

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-outscale/osc/lbu"
)

func resourceOutscaleOAPILoadBalancer() *schema.Resource {
	return &schema.Resource{
		Create: resourceOutscaleOAPILoadBalancerCreate,
		Read:   resourceOutscaleOAPILoadBalancerRead,
		Update: resourceOutscaleOAPILoadBalancerUpdate,
		Delete: resourceOutscaleOAPILoadBalancerDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"sub_region_name": &schema.Schema{
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"listener": &schema.Schema{
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"backend_port": &schema.Schema{
							Type:     schema.TypeInt,
							Required: true,
						},

						"backend_protocol": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},

						"load_balancer_port": &schema.Schema{
							Type:     schema.TypeInt,
							Required: true,
						},

						"load_balancer_protocol": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"server_certificate_id": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"load_balancer_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"load_balancer_type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"firewall_rules_set_name": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"subnet_id": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"tag": tagsSchema(),

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
				Optional: true,
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
			"source_firewall_rules_set": &schema.Schema{
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

func resourceOutscaleOAPILoadBalancerCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).LBU

	elbOpts := &lbu.CreateLoadBalancerInput{}

	listeners, err := expandOAPIListeners(d.Get("listener").([]interface{}))
	if err != nil {
		return err
	}

	elbOpts.Listeners = listeners

	if v, ok := d.GetOk("load_balancer_name"); ok {
		elbOpts.LoadBalancerName = aws.String(v.(string))
	}

	if v, ok := d.GetOk("tag"); ok {
		elbOpts.Tags = tagsFromMapC(v.(map[string]interface{}))

	}

	if v, ok := d.GetOk("load_balancer_type"); ok {
		elbOpts.Scheme = aws.String(v.(string))
	}

	if v, ok := d.GetOk("sub_region_name"); ok {
		elbOpts.AvailabilityZones = expandStringList(v.([]interface{}))
	}

	if v, ok := d.GetOk("firewall_rules_set_name"); ok {
		elbOpts.SecurityGroups = expandStringList(v.([]interface{}))
	}

	if v, ok := d.GetOk("subnet_id"); ok {
		elbOpts.Subnets = expandStringList(v.([]interface{}))
	}

	log.Printf("[DEBUG] ELB create configuration: %#v", elbOpts)
	var resp *lbu.CreateLoadBalancerOutput
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, err = conn.API.CreateLoadBalancer(elbOpts)

		if err != nil {
			if strings.Contains(fmt.Sprint(err), "CertificateNotFound") {
				return resource.RetryableError(
					fmt.Errorf("[WARN] Error creating ELB Listener with SSL Cert, retrying: %s", err))
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return err
	}

	// Assign the lbu's unique identifier for use later
	d.SetId(*elbOpts.LoadBalancerName)
	log.Printf("[INFO] ELB ID: %s", d.Id())

	// Enable partial mode and record what we set
	d.Partial(true)
	d.SetPartial("load_balancer_name")
	d.SetPartial("load_balancer_type")
	d.SetPartial("sub_region_name")
	d.SetPartial("listener_member")
	d.SetPartial("firewall_rules_set_name")
	d.SetPartial("subnet_id")

	return resourceOutscaleOAPILoadBalancerUpdate(d, meta)
}

func resourceOutscaleOAPILoadBalancerRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).LBU
	elbName := d.Id()

	// Retrieve the ELB properties for updating the state
	describeElbOpts := &lbu.DescribeLoadBalancersInput{
		LoadBalancerNames: []*string{aws.String(elbName)},
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
	if len(describeResp.LoadBalancerDescriptions) != 1 {
		return fmt.Errorf("Unable to find ELB: %#v", describeResp.LoadBalancerDescriptions)
	}

	describeAttrsOpts := &lbu.DescribeLoadBalancerAttributesInput{
		LoadBalancerName: aws.String(elbName),
	}

	var describeAttrsResp *lbu.DescribeLoadBalancerAttributesOutput
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		describeAttrsResp, err = conn.API.DescribeLoadBalancerAttributes(describeAttrsOpts)

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

	// lbAttrs := describeAttrsResp.LoadBalancerAttributes

	lb := describeResp.LoadBalancerDescriptions[0]

	d.Set("availability_zones", flattenOAPIStringList(lb.AvailabilityZones))
	d.Set("public_dns_name", aws.StringValue(lb.DNSName))
	if *lb.HealthCheck.Target != "" {
		d.Set("health_check", flattenOAPIHealthCheck(lb.HealthCheck))
	} else {
		d.Set("health_check", make(map[string]interface{}))
	}
	d.Set("backend_vm_id", flattenOAPIInstances(lb.Instances))
	d.Set("listeners", flattenOAPIListeners(lb.ListenerDescriptions))
	d.Set("load_balancer_name", lb.LoadBalancerName)

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
		policies["other_policy"] = flattenOAPIStringList(lb.Policies.OtherPolicies)
	}
	d.Set("policies", policies)
	d.Set("load_balancer_type", aws.StringValue(lb.Scheme))
	d.Set("firewall_rules_set_name", flattenOAPIStringList(lb.SecurityGroups))
	ssg := make(map[string]string)
	if lb.SourceSecurityGroup != nil {
		ssg["firewall_rules_set_name"] = aws.StringValue(lb.SourceSecurityGroup.GroupName)
		ssg["account_alias"] = aws.StringValue(lb.SourceSecurityGroup.OwnerAlias)
	}
	d.Set("source_firewall_rules_set", ssg)
	d.Set("subnet_id", flattenOAPIStringList(lb.Subnets))
	d.Set("lin_id", lb.VPCId)

	return nil
}

func resourceOutscaleOAPILoadBalancerUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).LBU

	d.Partial(true)

	if d.HasChange("listener") {
		o, n := d.GetChange("listener")
		os := o.(*schema.Set)
		ns := n.(*schema.Set)

		remove, _ := expandOAPIListeners(os.Difference(ns).List())
		add, _ := expandOAPIListeners(ns.Difference(os).List())

		if len(remove) > 0 {
			ports := make([]*int64, 0, len(remove))
			for _, listener := range remove {
				ports = append(ports, listener.LoadBalancerPort)
			}

			deleteListenersOpts := &lbu.DeleteLoadBalancerListenersInput{
				LoadBalancerName:  aws.String(d.Id()),
				LoadBalancerPorts: ports,
			}

			log.Printf("[DEBUG] ELB Delete Listeners opts: %s", deleteListenersOpts)

			var err error
			err = resource.Retry(5*time.Minute, func() *resource.RetryError {
				_, err = conn.API.DeleteLoadBalancerListeners(deleteListenersOpts)

				if err != nil {
					if strings.Contains(err.Error(), "RequestLimitExceeded:") {
						return resource.RetryableError(err)
					}
					return resource.NonRetryableError(err)
				}
				return nil
			})

			if err != nil {
				return fmt.Errorf("Failure removing outdated ELB listeners: %s", err)
			}
		}

		if len(add) > 0 {
			createListenersOpts := &lbu.CreateLoadBalancerListenersInput{
				LoadBalancerName: aws.String(d.Id()),
				Listeners:        add,
			}

			// Occasionally AWS will error with a 'duplicate listener', without any
			// other listeners on the ELB. Retry here to eliminate that.
			err := resource.Retry(5*time.Minute, func() *resource.RetryError {
				log.Printf("[DEBUG] ELB Create Listeners opts: %s", createListenersOpts)
				if _, err := conn.API.CreateLoadBalancerListeners(createListenersOpts); err != nil {
					if awsErr, ok := err.(awserr.Error); ok {
						if strings.Contains(fmt.Sprint(err), "DuplicateListener") {
							log.Printf("[DEBUG] Duplicate listener found for ELB (%s), retrying", d.Id())
							return resource.RetryableError(awsErr)
						}
						if strings.Contains(fmt.Sprint(err), "CertificateNotFound") && strings.Contains(fmt.Sprint(err), "Server Certificate not found for the key: arn") {
							log.Printf("[DEBUG] SSL Cert not found for given ARN, retrying")
							return resource.RetryableError(awsErr)
						}
					}

					// Didn't recognize the error, so shouldn't retry.
					return resource.NonRetryableError(err)
				}
				// Successful creation
				return nil
			})
			if err != nil {
				return fmt.Errorf("Failure adding new or updated ELB listeners: %s", err)
			}
		}

		d.SetPartial("listener")
	}

	// If we currently have instances, or did have instances,
	// we want to figure out what to add and remove from the load
	// balancer
	if d.HasChange("backend_vm_id") {
		o, n := d.GetChange("backend_vm_id")
		os := o.(*schema.Set)
		ns := n.(*schema.Set)
		remove := expandOAPIInstanceString(os.Difference(ns).List())
		add := expandOAPIInstanceString(ns.Difference(os).List())

		if len(add) > 0 {
			registerInstancesOpts := lbu.RegisterInstancesWithLoadBalancerInput{
				LoadBalancerName: aws.String(d.Id()),
				Instances:        add,
			}

			var err error
			err = resource.Retry(5*time.Minute, func() *resource.RetryError {
				_, err = conn.API.RegisterInstancesWithLoadBalancer(&registerInstancesOpts)

				if err != nil {
					if strings.Contains(err.Error(), "RequestLimitExceeded:") {
						return resource.RetryableError(err)
					}
					return resource.NonRetryableError(err)
				}
				return nil
			})

			if err != nil {
				return fmt.Errorf("Failure registering instances with ELB: %s", err)
			}
		}
		if len(remove) > 0 {
			deRegisterInstancesOpts := lbu.DeregisterInstancesFromLoadBalancerInput{
				LoadBalancerName: aws.String(d.Id()),
				Instances:        remove,
			}

			var err error
			err = resource.Retry(5*time.Minute, func() *resource.RetryError {
				_, err = conn.API.DeregisterInstancesFromLoadBalancer(&deRegisterInstancesOpts)

				if err != nil {
					if strings.Contains(err.Error(), "RequestLimitExceeded:") {
						return resource.RetryableError(err)
					}
					return resource.NonRetryableError(err)
				}
				return nil
			})

			if err != nil {
				return fmt.Errorf("Failure deregistering instances from ELB: %s", err)
			}
		}

		d.SetPartial("backend_vm_id")
	}

	if d.HasChange("health_check") {
		hc := d.Get("health_check").([]interface{})
		if len(hc) > 0 {
			check := hc[0].(map[string]interface{})
			configureHealthCheckOpts := lbu.ConfigureHealthCheckInput{
				LoadBalancerName: aws.String(d.Id()),
				HealthCheck: &lbu.HealthCheck{
					HealthyThreshold:   aws.Int64(int64(check["healthy_threshold"].(int))),
					UnhealthyThreshold: aws.Int64(int64(check["unhealthy_threshold"].(int))),
					Interval:           aws.Int64(int64(check["check_interval"].(int))),
					Target:             aws.String(check["checked_vm"].(string)),
					Timeout:            aws.Int64(int64(check["timeout"].(int))),
				},
			}
			var err error

			err = resource.Retry(5*time.Minute, func() *resource.RetryError {
				_, err = conn.API.ConfigureHealthCheck(&configureHealthCheckOpts)

				if err != nil {
					if strings.Contains(err.Error(), "RequestLimitExceeded:") {
						return resource.RetryableError(err)
					}
					return resource.NonRetryableError(err)
				}
				return nil
			})

			if err != nil {
				return fmt.Errorf("Failure configuring health check for ELB: %s", err)
			}
			d.SetPartial("health_check")
		}
	}

	if d.HasChange("firewall_rules_set_name") {
		groups := d.Get("firewall_rules_set_name").([]interface{})

		applySecurityGroupsOpts := lbu.ApplySecurityGroupsToLoadBalancerInput{
			LoadBalancerName: aws.String(d.Id()),
			SecurityGroups:   expandStringList(groups),
		}

		var err error
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			_, err = conn.API.ApplySecurityGroupsToLoadBalancer(&applySecurityGroupsOpts)

			if err != nil {
				if strings.Contains(err.Error(), "RequestLimitExceeded:") {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})

		if err != nil {
			return fmt.Errorf("Failure applying security groups to ELB: %s", err)
		}

		d.SetPartial("firewall_rules_set_name")
	}

	if d.HasChange("sub_region_name") {
		o, n := d.GetChange("sub_region_name")
		os := o.(*schema.Set)
		ns := n.(*schema.Set)

		removed := expandStringList(os.Difference(ns).List())
		added := expandStringList(ns.Difference(os).List())

		if len(added) > 0 {
			enableOpts := &lbu.EnableAvailabilityZonesForLoadBalancerInput{
				LoadBalancerName:  aws.String(d.Id()),
				AvailabilityZones: added,
			}

			log.Printf("[DEBUG] ELB enable availability zones opts: %s", enableOpts)
			var err error

			err = resource.Retry(5*time.Minute, func() *resource.RetryError {
				_, err = conn.API.EnableAvailabilityZonesForLoadBalancer(enableOpts)

				if err != nil {
					if strings.Contains(err.Error(), "RequestLimitExceeded:") {
						return resource.RetryableError(err)
					}
					return resource.NonRetryableError(err)
				}
				return nil
			})

			if err != nil {
				return fmt.Errorf("Failure enabling ELB availability zones: %s", err)
			}
		}

		if len(removed) > 0 {
			disableOpts := &lbu.DisableAvailabilityZonesForLoadBalancerInput{
				LoadBalancerName:  aws.String(d.Id()),
				AvailabilityZones: removed,
			}

			log.Printf("[DEBUG] ELB disable availability zones opts: %s", disableOpts)
			var err error

			err = resource.Retry(5*time.Minute, func() *resource.RetryError {
				_, err = conn.API.DisableAvailabilityZonesForLoadBalancer(disableOpts)

				if err != nil {
					if strings.Contains(err.Error(), "RequestLimitExceeded:") {
						return resource.RetryableError(err)
					}
					return resource.NonRetryableError(err)
				}
				return nil
			})

			if err != nil {
				return fmt.Errorf("Failure disabling ELB availability zones: %s", err)
			}
		}

		d.SetPartial("availability_zones")
	}

	if d.HasChange("subnet_id") {
		o, n := d.GetChange("subnet_id")
		os := o.(*schema.Set)
		ns := n.(*schema.Set)

		removed := expandStringList(os.Difference(ns).List())
		added := expandStringList(ns.Difference(os).List())

		if len(removed) > 0 {
			detachOpts := &lbu.DetachLoadBalancerFromSubnetsInput{
				LoadBalancerName: aws.String(d.Id()),
				Subnets:          removed,
			}

			log.Printf("[DEBUG] ELB detach subnet_id opts: %s", detachOpts)

			var err error

			err = resource.Retry(5*time.Minute, func() *resource.RetryError {
				_, err := conn.API.DetachLoadBalancerFromSubnets(detachOpts)

				if err != nil {
					if strings.Contains(err.Error(), "RequestLimitExceeded:") {
						return resource.RetryableError(err)
					}
					return resource.NonRetryableError(err)
				}
				return nil
			})

			if err != nil {
				return fmt.Errorf("Failure removing ELB subnets: %s", err)
			}
		}

		if len(added) > 0 {
			attachOpts := &lbu.AttachLoadBalancerToSubnetsInput{
				LoadBalancerName: aws.String(d.Id()),
				Subnets:          added,
			}
			var err error

			log.Printf("[DEBUG] ELB attach subnets opts: %s", attachOpts)
			err = resource.Retry(5*time.Minute, func() *resource.RetryError {
				_, err = conn.API.AttachLoadBalancerToSubnets(attachOpts)
				if err != nil {
					if awsErr, ok := err.(awserr.Error); ok {
						// eventually consistent issue with removing a subnet in AZ1 and
						// immediately adding a new one in the same AZ
						if awsErr.Code() == "InvalidConfigurationRequest" && strings.Contains(awsErr.Message(), "cannot be attached to multiple subnets in the same AZ") {
							log.Printf("[DEBUG] retrying az association")
							return resource.RetryableError(awsErr)
						}
					}
					return resource.NonRetryableError(err)
				}
				return nil
			})
			if err != nil {
				return fmt.Errorf("Failure adding ELB subnets: %s", err)
			}
		}

		d.SetPartial("subnet_id")
	}

	d.Partial(false)

	return resourceOutscaleOAPILoadBalancerRead(d, meta)
}

func resourceOutscaleOAPILoadBalancerDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).LBU

	log.Printf("[INFO] Deleting ELB: %s", d.Id())

	// Destroy the load balancer
	deleteElbOpts := lbu.DeleteLoadBalancerInput{
		LoadBalancerName: aws.String(d.Id()),
	}
	var err error

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, err = conn.API.DeleteLoadBalancer(&deleteElbOpts)
		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded:") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("Error deleting ELB: %s", err)
	}

	return nil
}

func expandOAPIListeners(configured []interface{}) ([]*lbu.Listener, error) {
	listeners := make([]*lbu.Listener, 0, len(configured))

	for _, lRaw := range configured {
		data := lRaw.(map[string]interface{})

		ip := int64(data["backend_port"].(int))
		lp := int64(data["load_balancer_port"].(int))
		l := &lbu.Listener{
			InstancePort:     &ip,
			InstanceProtocol: aws.String(data["backend_protocol"].(string)),
			LoadBalancerPort: &lp,
			Protocol:         aws.String(data["load_balancer_protocol"].(string)),
		}

		if v, ok := data["server_certificate_id"]; ok {
			l.SSLCertificateId = aws.String(v.(string))
		}

		var valid bool
		if l.SSLCertificateId != nil && *l.SSLCertificateId != "" {
			// validate the load_balancer_protocol is correct
			for _, p := range []string{"https", "ssl"} {
				if (strings.ToLower(*l.InstanceProtocol) == p) || (strings.ToLower(*l.Protocol) == p) {
					valid = true
				}
			}
		} else {
			valid = true
		}

		if valid {
			listeners = append(listeners, l)
		} else {
			return nil, fmt.Errorf("[ERR] ELB Listener: server_certificate_id may be set only when load_balancer_protocol is 'https' or 'ssl'")
		}
	}

	return listeners, nil
}

func flattenOAPIStringList(list []*string) []interface{} {
	vs := make([]interface{}, 0, len(list))
	for _, v := range list {
		vs = append(vs, *v)
	}
	return vs
}

func flattenOAPIInstances(list []*lbu.Instance) []map[string]string {
	result := make([]map[string]string, len(list))
	for _, i := range list {
		result = append(result, map[string]string{"vm_id": *i.InstanceId})
	}
	return result
}

// Expands an array of String Instance IDs into a []Instances
func expandOAPIInstanceString(list []interface{}) []*lbu.Instance {
	result := make([]*lbu.Instance, 0, len(list))
	for _, i := range list {
		result = append(result, &lbu.Instance{InstanceId: aws.String(i.(string))})
	}
	return result
}

// Flattens an array of Listeners into a []map[string]interface{}
func flattenOAPIListeners(list []*lbu.ListenerDescription) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(list))

	for _, i := range list {
		l := make(map[string]interface{})
		listener := map[string]interface{}{
			"backend_port":           aws.Int64Value(i.Listener.InstancePort),
			"backend_protocol":       strings.ToLower(aws.StringValue(i.Listener.InstanceProtocol)),
			"load_balancer_port":     aws.Int64Value(i.Listener.LoadBalancerPort),
			"load_balancer_protocol": strings.ToLower(aws.StringValue(i.Listener.Protocol)),
			"server_certificate_id":  aws.StringValue(i.Listener.SSLCertificateId),
		}
		l["listener"] = listener
		l["policy_name"] = flattenOAPIStringList(i.PolicyNames)
		result = append(result, l)
	}
	return result
}

func flattenOAPIHealthCheck(check *lbu.HealthCheck) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, 1)

	chk := make(map[string]interface{})
	chk["unhealthy_threshold"] = *check.UnhealthyThreshold
	chk["healthy_threshold"] = *check.HealthyThreshold
	chk["checked_vm"] = *check.Target
	chk["timeout"] = *check.Timeout
	chk["check_interval"] = *check.Interval

	result = append(result, chk)

	return result
}
