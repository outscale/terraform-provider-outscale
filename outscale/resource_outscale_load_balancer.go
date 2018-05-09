package outscale

import (
	"bytes"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/terraform-providers/terraform-provider-outscale/osc/lbu"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceOutscaleLoadBalancer() *schema.Resource {
	return &schema.Resource{
		Create: resourceOutscaleLoadBalancerCreate,
		Read:   resourceOutscaleLoadBalancerRead,
		Update: resourceOutscaleLoadBalancerUpdate,
		Delete: resourceOutscaleLoadBalancerDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"availability_zones_member": &schema.Schema{
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"listeners_member": &schema.Schema{
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_port": &schema.Schema{
							Type:         schema.TypeInt,
							Required:     true,
						},

						"instance_protocol": &schema.Schema{
							Type:         schema.TypeString,
							Required:     true,
						},

						"load_balancer_port": &schema.Schema{
							Type:         schema.TypeInt,
							Required:     true,
						},

						"protocol": &schema.Schema{
							Type:         schema.TypeString,
							Required:     true,
						},
						"ssl_certificate_id": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"load_balancer_name": &schema.Schema{
				Type:          schema.TypeString,
				Required:     true,
				ForceNew:      true,
			},
			"scheme": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Computed: true,
				ForceNew:     true,
			},
			"security_groups.member": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
"subnets_member": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"tags.member": tagsSchema(),

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
							Type:         schema.TypeInt,
							Computed:     true,
						},
						"unhealthy_threshold": &schema.Schema{
							Type:         schema.TypeInt,
							Computed:     true,
						},
						"target": &schema.Schema{
							Type:         schema.TypeString,
							Computed:     true,
						},
						"interval": &schema.Schema{
							Type:         schema.TypeInt,
							Computed:     true,
						},
						"timeout": &schema.Schema{
							Type:         schema.TypeInt,
							Computed:     true,
						},
					},
				},
			},
	"instances.member": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_id": &schema.Schema{
							Type: schema.TypeString,
							Computed: true,
						},
					},
				},
			},
"listener_descriptions.member": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"listener": &schema.Schema{
							Type:         schema.TypeMap,
							Computed:     true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"instance_port": &schema.Schema{
							Type:         schema.TypeInt,
							Computed:     true,
						},
						"instance_protocol": &schema.Schema{
							Type:         schema.TypeString,
							Computed:     true,
						},
						"load_balancer_port": &schema.Schema{
							Type:         schema.TypeInt,
							Computed:     true,
						},
						"protocol": &schema.Schema{
							Type:         schema.TypeString,
							Computed:     true,
						},
						"ssl_certificate_id": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
								},
							},
						},
						"policy_names.member": &schema.Schema{
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{Type: schema.TypeString},
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

		},
	}
}

func resourceOutscaleLoadBalancerCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).LBU

		elbOpts := &lbu.CreateLoadBalancerInput{}

	listeners, err := expandListeners(d.Get("listeners_member").([]interface{}))
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

	if v, ok := d.GetOk("scheme"); ok {
		elbOpts.Scheme = aws.String(v.(string))
	}

	if v, ok := d.GetOk("availability_zones_member"); ok {
		elbOpts.AvailabilityZones = expandStringList(v.([]interface{}))
	}

	if v, ok := d.GetOk("security_groups_member"); ok {
		elbOpts.SecurityGroups = expandStringList(v.([]interface{}))
	}

	if v, ok := d.GetOk("subnets_member"); ok {
		elbOpts.Subnets = expandStringList(v.([]interface{}))
	}

	log.Printf("[DEBUG] ELB create configuration: %#v", elbOpts)
	var resp  *lbu.CreateLoadBalancerOutput
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, err := conn.API.CreateLoadBalancer(elbOpts)

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
	d.SetPartial("scheme")
	d.SetPartial("availability_zones_member")
	d.SetPartial("listener_member")
	d.SetPartial("security_groups_member")
	d.SetPartial("subnets_member")

	return resourceOutscaleLoadBalancerUpdate(d, meta)
}

func resourceOutscaleLoadBalancerRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).LBU
	elbName := d.Id()

	// Retrieve the ELB properties for updating the state
	describeElbOpts := &lbu.DescribeLoadBalancersInput{
		LoadBalancerNames: []*string{aws.String(elbName)},
	}

	describeResp, err := conn.API.DescribeLoadBalancers(describeElbOpts)
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
	describeAttrsResp, err := conn.API.DescribeLoadBalancerAttributes(describeAttrsOpts)
	if err != nil {
		if isLoadBalancerNotFound(err) {
			// The ELB is gone now, so just remove it from the state
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error retrieving ELB: %s", err)
	}

	lbAttrs := describeAttrsResp.LoadBalancerAttributes

	lb := describeResp.LoadBalancerDescriptions[0]

	d.Set("name", lb.LoadBalancerName)
	d.Set("dns_name", lb.DNSName)
	d.Set("zone_id", lb.CanonicalHostedZoneNameID)

	var scheme bool
	if lb.Scheme != nil {
		scheme = *lb.Scheme == "internal"
	}
	d.Set("internal", scheme)
	d.Set("availability_zones", flattenStringList(lb.AvailabilityZones))
	d.Set("instances", flattenInstances(lb.Instances))
	d.Set("listener", flattenListeners(lb.ListenerDescriptions))
	d.Set("security_groups", flattenStringList(lb.SecurityGroups))
	if lb.SourceSecurityGroup != nil {
		group := lb.SourceSecurityGroup.GroupName
		if lb.SourceSecurityGroup.OwnerAlias != nil && *lb.SourceSecurityGroup.OwnerAlias != "" {
			group = aws.String(*lb.SourceSecurityGroup.OwnerAlias + "/" + *lb.SourceSecurityGroup.GroupName)
		}
		d.Set("source_security_group", group)

		// Manually look up the ELB Security Group ID, since it's not provided
		var elbVpc string
		if lb.VPCId != nil {
			elbVpc = *lb.VPCId
			sgId, err := sourceSGIdByName(meta, *lb.SourceSecurityGroup.GroupName, elbVpc)
			if err != nil {
				return fmt.Errorf("[WARN] Error looking up ELB Security Group ID: %s", err)
			} else {
				d.Set("source_security_group_id", sgId)
			}
		}
	}
	d.Set("subnets", flattenStringList(lb.Subnets))
	if lbAttrs.ConnectionSettings != nil {
		d.Set("idle_timeout", lbAttrs.ConnectionSettings.IdleTimeout)
	}
	d.Set("connection_draining", lbAttrs.ConnectionDraining.Enabled)
	d.Set("connection_draining_timeout", lbAttrs.ConnectionDraining.Timeout)
	d.Set("cross_zone_load_balancing", lbAttrs.CrossZoneLoadBalancing.Enabled)
	if lbAttrs.AccessLog != nil {
		// The AWS API does not allow users to remove access_logs, only disable them.
		// During creation of the ELB, Terraform sets the access_logs to disabled,
		// so there should not be a case where lbAttrs.AccessLog above is nil.

		// Here we do not record the remove value of access_log if:
		// - there is no access_log block in the configuration
		// - the remote access_logs are disabled
		//
		// This indicates there is no access_log in the configuration.
		// - externally added access_logs will be enabled, so we'll detect the drift
		// - locally added access_logs will be in the config, so we'll add to the
		// API/state
		// See https://github.com/hashicorp/terraform/issues/10138
		_, n := d.GetChange("access_logs")
		elbal := lbAttrs.AccessLog
		nl := n.([]interface{})
		if len(nl) == 0 && !*elbal.Enabled {
			elbal = nil
		}
		if err := d.Set("access_logs", flattenAccessLog(elbal)); err != nil {
			return err
		}
	}

	resp, err := conn.API.DescribeTags(&lbu.DescribeTagsInput{
		LoadBalancerNames: []*string{lb.LoadBalancerName},
	})

	var et []*lbu.Tag
	if len(resp.TagDescriptions) > 0 {
		et = resp.TagDescriptions[0].Tags
	}
	d.Set("tags", tagsToMapELB(et))

	// There's only one health check, so save that to state as we
	// currently can
	if *lb.HealthCheck.Target != "" {
		d.Set("health_check", flattenHealthCheck(lb.HealthCheck))
	}

	return nil
}

func resourceOutscaleLoadBalancerUpdate(d *schema.ResourceData, meta interface{}) error {
	LBU := meta.(*OutscaleClient).LBU

	d.Partial(true)

	if d.HasChange("listener") {
		o, n := d.GetChange("listener")
		os := o.(*schema.Set)
		ns := n.(*schema.Set)

		remove, _ := expandListeners(os.Difference(ns).List())
		add, _ := expandListeners(ns.Difference(os).List())

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
			_, err := LBU.DeleteLoadBalancerListeners(deleteListenersOpts)
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
				if _, err := LBU.CreateLoadBalancerListeners(createListenersOpts); err != nil {
					if awsErr, ok := err.(awserr.Error); ok {
						if awsErr.Code() == "DuplicateListener" {
							log.Printf("[DEBUG] Duplicate listener found for ELB (%s), retrying", d.Id())
							return resource.RetryableError(awsErr)
						}
						if awsErr.Code() == "CertificateNotFound" && strings.Contains(awsErr.Message(), "Server Certificate not found for the key: arn") {
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
	if d.HasChange("instances") {
		o, n := d.GetChange("instances")
		os := o.(*schema.Set)
		ns := n.(*schema.Set)
		remove := expandInstanceString(os.Difference(ns).List())
		add := expandInstanceString(ns.Difference(os).List())

		if len(add) > 0 {
			registerInstancesOpts := lbu.RegisterInstancesWithLoadBalancerInput{
				LoadBalancerName: aws.String(d.Id()),
				Instances:        add,
			}

			_, err := LBU.RegisterInstancesWithLoadBalancer(&registerInstancesOpts)
			if err != nil {
				return fmt.Errorf("Failure registering instances with ELB: %s", err)
			}
		}
		if len(remove) > 0 {
			deRegisterInstancesOpts := lbu.DeregisterInstancesFromLoadBalancerInput{
				LoadBalancerName: aws.String(d.Id()),
				Instances:        remove,
			}

			_, err := LBU.DeregisterInstancesFromLoadBalancer(&deRegisterInstancesOpts)
			if err != nil {
				return fmt.Errorf("Failure deregistering instances from ELB: %s", err)
			}
		}

		d.SetPartial("instances")
	}

	if d.HasChange("cross_zone_load_balancing") || d.HasChange("idle_timeout") || d.HasChange("access_logs") {
		attrs := lbu.ModifyLoadBalancerAttributesInput{
			LoadBalancerName: aws.String(d.Get("name").(string)),
			LoadBalancerAttributes: &lbu.LoadBalancerAttributes{
				CrossZoneLoadBalancing: &lbu.CrossZoneLoadBalancing{
					Enabled: aws.Bool(d.Get("cross_zone_load_balancing").(bool)),
				},
				ConnectionSettings: &lbu.ConnectionSettings{
					IdleTimeout: aws.Int64(int64(d.Get("idle_timeout").(int))),
				},
			},
		}

		logs := d.Get("access_logs").([]interface{})
		if len(logs) == 1 {
			l := logs[0].(map[string]interface{})
			accessLog := &lbu.AccessLog{
				Enabled:      aws.Bool(l["enabled"].(bool)),
				EmitInterval: aws.Int64(int64(l["interval"].(int))),
				S3BucketName: aws.String(l["bucket"].(string)),
			}

			if l["bucket_prefix"] != "" {
				accessLog.S3BucketPrefix = aws.String(l["bucket_prefix"].(string))
			}

			attrs.LoadBalancerAttributes.AccessLog = accessLog
		} else if len(logs) == 0 {
			// disable access logs
			attrs.LoadBalancerAttributes.AccessLog = &lbu.AccessLog{
				Enabled: aws.Bool(false),
			}
		}

		log.Printf("[DEBUG] ELB Modify Load Balancer Attributes Request: %#v", attrs)
		_, err := LBU.ModifyLoadBalancerAttributes(&attrs)
		if err != nil {
			return fmt.Errorf("Failure configuring ELB attributes: %s", err)
		}

		d.SetPartial("cross_zone_load_balancing")
		d.SetPartial("idle_timeout")
		d.SetPartial("connection_draining_timeout")
	}

	// We have to do these changes separately from everything else since
	// they have some weird undocumented rules. You can't set the timeout
	// without having connection draining to true, so we set that to true,
	// set the timeout, then reset it to false if requested.
	if d.HasChange("connection_draining") || d.HasChange("connection_draining_timeout") {
		// We do timeout changes first since they require us to set draining
		// to true for a hot second.
		if d.HasChange("connection_draining_timeout") {
			attrs := lbu.ModifyLoadBalancerAttributesInput{
				LoadBalancerName: aws.String(d.Get("name").(string)),
				LoadBalancerAttributes: &lbu.LoadBalancerAttributes{
					ConnectionDraining: &lbu.ConnectionDraining{
						Enabled: aws.Bool(true),
						Timeout: aws.Int64(int64(d.Get("connection_draining_timeout").(int))),
					},
				},
			}

			_, err := LBU.ModifyLoadBalancerAttributes(&attrs)
			if err != nil {
				return fmt.Errorf("Failure configuring ELB attributes: %s", err)
			}

			d.SetPartial("connection_draining_timeout")
		}

		// Then we always set connection draining even if there is no change.
		// This lets us reset to "false" if requested even with a timeout
		// change.
		attrs := lbu.ModifyLoadBalancerAttributesInput{
			LoadBalancerName: aws.String(d.Get("name").(string)),
			LoadBalancerAttributes: &lbu.LoadBalancerAttributes{
				ConnectionDraining: &lbu.ConnectionDraining{
					Enabled: aws.Bool(d.Get("connection_draining").(bool)),
				},
			},
		}

		_, err := LBU.ModifyLoadBalancerAttributes(&attrs)
		if err != nil {
			return fmt.Errorf("Failure configuring ELB attributes: %s", err)
		}

		d.SetPartial("connection_draining")
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
					Interval:           aws.Int64(int64(check["interval"].(int))),
					Target:             aws.String(check["target"].(string)),
					Timeout:            aws.Int64(int64(check["timeout"].(int))),
				},
			}
			_, err := LBU.ConfigureHealthCheck(&configureHealthCheckOpts)
			if err != nil {
				return fmt.Errorf("Failure configuring health check for ELB: %s", err)
			}
			d.SetPartial("health_check")
		}
	}

	if d.HasChange("security_groups") {
		groups := d.Get("security_groups").(*schema.Set).List()

		applySecurityGroupsOpts := lbu.ApplySecurityGroupsToLoadBalancerInput{
			LoadBalancerName: aws.String(d.Id()),
			SecurityGroups:   expandStringList(groups),
		}

		_, err := LBU.ApplySecurityGroupsToLoadBalancer(&applySecurityGroupsOpts)
		if err != nil {
			return fmt.Errorf("Failure applying security groups to ELB: %s", err)
		}

		d.SetPartial("security_groups")
	}

	if d.HasChange("availability_zones") {
		o, n := d.GetChange("availability_zones")
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
			_, err := LBU.EnableAvailabilityZonesForLoadBalancer(enableOpts)
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
			_, err := LBU.DisableAvailabilityZonesForLoadBalancer(disableOpts)
			if err != nil {
				return fmt.Errorf("Failure disabling ELB availability zones: %s", err)
			}
		}

		d.SetPartial("availability_zones")
	}

	if d.HasChange("subnets") {
		o, n := d.GetChange("subnets")
		os := o.(*schema.Set)
		ns := n.(*schema.Set)

		removed := expandStringList(os.Difference(ns).List())
		added := expandStringList(ns.Difference(os).List())

		if len(removed) > 0 {
			detachOpts := &lbu.DetachLoadBalancerFromSubnetsInput{
				LoadBalancerName: aws.String(d.Id()),
				Subnets:          removed,
			}

			log.Printf("[DEBUG] ELB detach subnets opts: %s", detachOpts)
			_, err := LBU.DetachLoadBalancerFromSubnets(detachOpts)
			if err != nil {
				return fmt.Errorf("Failure removing ELB subnets: %s", err)
			}
		}

		if len(added) > 0 {
			attachOpts := &lbu.AttachLoadBalancerToSubnetsInput{
				LoadBalancerName: aws.String(d.Id()),
				Subnets:          added,
			}

			log.Printf("[DEBUG] ELB attach subnets opts: %s", attachOpts)
			err := resource.Retry(5*time.Minute, func() *resource.RetryError {
				_, err := LBU.AttachLoadBalancerToSubnets(attachOpts)
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

		d.SetPartial("subnets")
	}

	if err := setTagsELB(LBU, d); err != nil {
		return err
	}

	d.SetPartial("tags")
	d.Partial(false)

	return resourceOutscaleLoadBalancerRead(d, meta)
}

func resourceOutscaleLoadBalancerDelete(d *schema.ResourceData, meta interface{}) error {
	LBU := meta.(*OutscaleClient).LBU

	log.Printf("[INFO] Deleting ELB: %s", d.Id())

	// Destroy the load balancer
	deleteElbOpts := lbu.DeleteLoadBalancerInput{
		LoadBalancerName: aws.String(d.Id()),
	}
	if _, err := LBU.DeleteLoadBalancer(&deleteElbOpts); err != nil {
		return fmt.Errorf("Error deleting ELB: %s", err)
	}

	return nil
}

func resourceOutscaleLoadBalancerListenerHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})
	buf.WriteString(fmt.Sprintf("%d-", m["instance_port"].(int)))
	buf.WriteString(fmt.Sprintf("%s-",
		strings.ToLower(m["instance_protocol"].(string))))
	buf.WriteString(fmt.Sprintf("%d-", m["lb_port"].(int)))
	buf.WriteString(fmt.Sprintf("%s-",
		strings.ToLower(m["lb_protocol"].(string))))

	if v, ok := m["ssl_certificate_id"]; ok {
		buf.WriteString(fmt.Sprintf("%s-", v.(string)))
	}

	return hashcode.String(buf.String())
}

func isLoadBalancerNotFound(err error) bool {
	elberr, ok := err.(awserr.Error)
	return ok && elberr.Code() == "LoadBalancerNotFound"
}

func sourceSGIdByName(meta interface{}, sg, vpcId string) (string, error) {
	conn := meta.(*OutscaleClient).ec2conn
	var filters []*ec2.Filter
	var sgFilterName, sgFilterVPCID *ec2.Filter
	sgFilterName = &ec2.Filter{
		Name:   aws.String("group-name"),
		Values: []*string{aws.String(sg)},
	}

	if vpcId != "" {
		sgFilterVPCID = &ec2.Filter{
			Name:   aws.String("vpc-id"),
			Values: []*string{aws.String(vpcId)},
		}
	}

	filters = append(filters, sgFilterName)

	if sgFilterVPCID != nil {
		filters = append(filters, sgFilterVPCID)
	}

	req := &ec2.DescribeSecurityGroupsInput{
		Filters: filters,
	}
	resp, err := conn.DescribeSecurityGroups(req)
	if err != nil {
		if ec2err, ok := err.(awserr.Error); ok {
			if ec2err.Code() == "InvalidSecurityGroupID.NotFound" ||
				ec2err.Code() == "InvalidGroup.NotFound" {
				resp = nil
				err = nil
			}
		}

		if err != nil {
			log.Printf("Error on ELB SG look up: %s", err)
			return "", err
		}
	}

	if resp == nil || len(resp.SecurityGroups) == 0 {
		return "", fmt.Errorf("No security groups found for name %s and vpc id %s", sg, vpcId)
	}

	group := resp.SecurityGroups[0]
	return *group.GroupId, nil
}

func validateAccessLogsInterval(v interface{}, k string) (ws []string, errors []error) {
	value := v.(int)

	// Check if the value is either 5 or 60 (minutes).
	if value != 5 && value != 60 {
		errors = append(errors, fmt.Errorf(
			"%q contains an invalid Access Logs interval \"%d\". "+
				"Valid intervals are either 5 or 60 (minutes).",
			k, value))
	}
	return
}

func validateHeathCheckTarget(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)

	// Parse the Health Check target value.
	matches := regexp.MustCompile(`\A(\w+):(\d+)(.+)?\z`).FindStringSubmatch(value)

	// Check if the value contains a valid target.
	if matches == nil || len(matches) < 1 {
		errors = append(errors, fmt.Errorf(
			"%q contains an invalid Health Check: %s",
			k, value))

		// Invalid target? Return immediately,
		// there is no need to collect other
		// errors.
		return
	}

	// Check if the value contains a valid protocol.
	if !isValidProtocol(matches[1]) {
		errors = append(errors, fmt.Errorf(
			"%q contains an invalid Health Check protocol %q. "+
				"Valid protocols are either %q, %q, %q, or %q.",
			k, matches[1], "TCP", "SSL", "HTTP", "HTTPS"))
	}

	// Check if the value contains a valid port range.
	port, _ := strconv.Atoi(matches[2])
	if port < 1 || port > 65535 {
		errors = append(errors, fmt.Errorf(
			"%q contains an invalid Health Check target port \"%d\". "+
				"Valid port is in the range from 1 to 65535 inclusive.",
			k, port))
	}

	switch strings.ToLower(matches[1]) {
	case "tcp", "ssl":
		// Check if value is in the form <PROTOCOL>:<PORT> for TCP and/or SSL.
		if matches[3] != "" {
			errors = append(errors, fmt.Errorf(
				"%q cannot contain a path in the Health Check target: %s",
				k, value))
		}
		break
	case "http", "https":
		// Check if value is in the form <PROTOCOL>:<PORT>/<PATH> for HTTP and/or HTTPS.
		if matches[3] == "" {
			errors = append(errors, fmt.Errorf(
				"%q must contain a path in the Health Check target: %s",
				k, value))
		}

		// Cannot be longer than 1024 multibyte characters.
		if len([]rune(matches[3])) > 1024 {
			errors = append(errors, fmt.Errorf("%q cannot contain a path longer "+
				"than 1024 characters in the Health Check target: %s",
				k, value))
		}
		break
	}

	return
}

func validateListenerProtocol(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)

	if !isValidProtocol(value) {
		errors = append(errors, fmt.Errorf(
			"%q contains an invalid Listener protocol %q. "+
				"Valid protocols are either %q, %q, %q, or %q.",
			k, value, "TCP", "SSL", "HTTP", "HTTPS"))
	}
	return
}

func isValidProtocol(s string) bool {
	if s == "" {
		return false
	}
	s = strings.ToLower(s)

	validProtocols := map[string]bool{
		"http":  true,
		"https": true,
		"ssl":   true,
		"tcp":   true,
	}

	if _, ok := validProtocols[s]; !ok {
		return false
	}

	return true
}

func expandListeners(configured []interface{}) ([]*lbu.Listener, error) {
	listeners := make([]*elb.Listener, 0, len(configured))

	// Loop over our configured listeners and create
	// an array of aws-sdk-go compatible objects
	for _, lRaw := range configured {
		data := lRaw.(map[string]interface{})

		ip := int64(data["instance_port"].(int))
		lp := int64(data["lb_port"].(int))
		l := &elb.Listener{
			InstancePort:     &ip,
			InstanceProtocol: aws.String(data["instance_protocol"].(string)),
			LoadBalancerPort: &lp,
			Protocol:         aws.String(data["lb_protocol"].(string)),
		}

		if v, ok := data["ssl_certificate_id"]; ok {
			l.SSLCertificateId = aws.String(v.(string))
		}

		var valid bool
		if l.SSLCertificateId != nil && *l.SSLCertificateId != "" {
			// validate the protocol is correct
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
			return nil, fmt.Errorf("[ERR] ELB Listener: ssl_certificate_id may be set only when protocol is 'https' or 'ssl'")
		}
	}

	return listeners, nil
}

func expandStringList(configured []interface{}) []*string {
	vs := make([]*string, 0, len(configured))
	for _, v := range configured {
		val, ok := v.(string)
		if ok && val != "" {
			vs = append(vs, aws.String(v.(string)))
		}
	}
	return vs
}