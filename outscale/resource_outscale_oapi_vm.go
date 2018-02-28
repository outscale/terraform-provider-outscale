package outscale

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-outscale/osc/fcu"
)

func resourceOutscaleOApiVM() *schema.Resource {
	return &schema.Resource{
		Create: resourceOAPIVMCreate,
		Read:   resourceOAPIVMRead,
		Update: resourceOAPIVMUpdate,
		Delete: resourceOAPIVMDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: getOApiVMSchema(),
	}
}

func resourceOAPIVMCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	instanceOpts, err := buildOutscaleOAPIVMOpts(d, meta)
	if err != nil {
		return err
	}

	// Build the creation struct
	runOpts := &fcu.RunInstancesInput{
		BlockDeviceMappings: instanceOpts.BlockDeviceMappings,
		// DisableApiTermination: instanceOpts.DisableAPITermination,
		EbsOptimized: instanceOpts.EBSOptimized,
		// Monitoring:            instanceOpts.Monitoring,
		// IamInstanceProfile:    instanceOpts.IAMInstanceProfile,
		ImageId: instanceOpts.ImageID,
		InstanceInitiatedShutdownBehavior: instanceOpts.InstanceInitiatedShutdownBehavior,
		InstanceType:                      instanceOpts.InstanceType,
		// Ipv6AddressCount:                  instanceOpts.Ipv6AddressCount,
		// Ipv6Addresses:                     instanceOpts.Ipv6Addresses,
		KeyName:           instanceOpts.KeyName,
		MaxCount:          aws.Int64(int64(1)),
		MinCount:          aws.Int64(int64(1)),
		NetworkInterfaces: instanceOpts.NetworkInterfaces,
		Placement:         instanceOpts.Placement,
		// PrivateIpAddress:                  instanceOpts.PrivateIPAddress,
		SecurityGroupIds: instanceOpts.SecurityGroupIDs,
		SecurityGroups:   instanceOpts.SecurityGroups,
		SubnetId:         instanceOpts.SubnetID,
		UserData:         instanceOpts.UserData,
	}

	// Create the instance
	log.Printf("[DEBUG] Run configuration: %s", runOpts)

	var runResp *fcu.Reservation
	err = resource.Retry(30*time.Second, func() *resource.RetryError {
		var err error
		runResp, err = conn.VM.RunInstance(runOpts)

		return resource.RetryableError(err)
	})

	if err != nil {
		return fmt.Errorf("Error launching source instance: %s", err)
	}
	if runResp == nil || len(runResp.Instances) == 0 {
		return errors.New("Error launching source instance: no instances returned in response")
	}

	instance := runResp.Instances[0]
	log.Printf("[INFO] Instance ID: %s", *instance.InstanceId)

	d.SetId(*instance.InstanceId)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"pending"},
		Target:     []string{"running"},
		Refresh:    InstanceStateOApiRefreshFunc(conn, *instance.InstanceId, "terminated"),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf(
			"Error waiting for instance (%s) to stop: %s", d.Id(), err)
	}

	// Initialize the connection info
	if instance.IpAddress != nil {
		d.SetConnInfo(map[string]string{
			"type": "ssh",
			"host": *instance.IpAddress,
		})
	} else if instance.PrivateIpAddress != nil {
		d.SetConnInfo(map[string]string{
			"type": "ssh",
			"host": *instance.PrivateIpAddress,
		})
	}

	return resourceOAPIVMRead(d, meta)
}

func resourceOAPIVMRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	input := &fcu.DescribeInstancesInput{
		InstanceIds: []*string{aws.String(d.Id())},
	}

	var resp *fcu.DescribeInstancesOutput
	var err error

	err = resource.Retry(30*time.Second, func() *resource.RetryError {
		resp, err = conn.VM.DescribeInstances(input)

		return resource.RetryableError(err)
	})

	if err != nil {
		return fmt.Errorf("Error deleting the instance %s", err)
	}

	if err != nil {
		// If the instance was not found, return nil so that we can show
		// that the instance is gone.
		if ec2err, ok := err.(awserr.Error); ok && ec2err.Code() == "InvalidInstanceID.NotFound" {
			d.SetId("")
			return nil
		}

		// Some other error, report it
		return err
	}

	// If nothing was found, then return no state
	if len(resp.Reservations) == 0 {
		d.SetId("")
		return nil
	}

	instance := resp.Reservations[0].Instances[0]

	d.Set("block_device_mapping", getOAPIVMBlockDeviceMapping(instance.BlockDeviceMappings))
	d.Set("token", instance.ClientToken)

	// d.Set("delete_protection", instance.DnsName)

	d.Set("z", instance.EbsOptimized)
	d.Set("image_id", instance.ImageId)
	d.Set("vm_id", instance.InstanceId)

	// d.Set("shutdown_automatic_behavior", instance.SpotInstanceRequestId)

	d.Set("type", instance.InstanceType)
	d.Set("keypair_name", instance.KeyName)

	// d.Set("max_vms_count", instance)
	// d.Set("min_vms_count", instance.KernelId)

	d.Set("nics", getOAPIVMNetworkInterfaceSet(instance.NetworkInterfaces))
	d.Set("placement", map[string]interface{}{
		"affinity":        instance.Placement.Affinity,
		"sub_region_name": instance.Placement.GroupName,
		// TODO: Add to struct for OAPI
		// "firewall_rules_set_name": instance.Placement.FirewallRulesSetName,
		"dedicated_host_id": instance.Placement.HostId,
		"tenancy":           instance.Placement.Tenancy,
	})
	d.Set("private_ip", instance.PrivateIpAddress)
	// d.Set("private_ips", ips)
	// d.Set("firewall_rules_set", ips)
	// d.Set("firewall_rules_set_id", ips)
	// d.Set("subnet_id", ips)
	// d.Set("user_data", ips)

	return nil
}

func resourceOAPIVMUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU
	log.Printf("[DEBUG] updating the instance %s", d.Id())

	if d.HasChange("key_name") {
		input := &fcu.ModifyInstanceKeyPairInput{
			InstanceId: aws.String(d.Id()),
			KeyName:    aws.String(d.Get("keypair_name").(string)),
		}

		err := conn.VM.ModifyInstanceKeyPair(input)
		if err != nil {
			return err
		}
	}
	return resourceVMRead(d, meta)
}

func resourceOAPIVMDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	id := d.Id()

	log.Printf("[INFO] Terminating instance: %s", id)
	req := &fcu.TerminateInstancesInput{
		InstanceIds: []*string{aws.String(id)},
	}

	var err error
	err = resource.Retry(30*time.Second, func() *resource.RetryError {
		_, err = conn.VM.TerminateInstances(req)

		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded") {
				log.Printf("[INFO] Request limit exceeded")
				return resource.RetryableError(err)
			}
		}

		return resource.RetryableError(err)
	})

	if err != nil {
		return fmt.Errorf("Error deleting the instance")
	}

	log.Printf("[DEBUG] Waiting for instance (%s) to become terminated", id)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"pending", "running", "shutting-down", "stopped", "stopping"},
		Target:     []string{"terminated"},
		Refresh:    InstanceStateRefreshFunc(conn, id, ""),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf(
			"Error waiting for instance (%s) to terminate: %s", id, err)
	}

	return nil
}

func getOApiVMSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		// Attributes
		"block_device_mapping": {
			Type:     schema.TypeSet,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"device_name": {
						Type:     schema.TypeString,
						Optional: true,
					},
					"bsu": {
						Type:     schema.TypeMap,
						Optional: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"delete_on_vm_deletion": {
									Type:     schema.TypeBool,
									Optional: true,
								},
								"iops": {
									Type:     schema.TypeString,
									Optional: true,
								},
								"snapshot_id": {
									Type:     schema.TypeInt,
									Optional: true,
								},
								"volume_size": {
									Type:     schema.TypeFloat,
									Optional: true,
								},
								"type": {
									Type:     schema.TypeString,
									Optional: true,
								},
							},
						},
					},
					"no_device": {
						Type:     schema.TypeBool,
						Optional: true,
					},
					"virtual_device_name": {
						Type:     schema.TypeString,
						Optional: true,
					},
				},
			},
		},
		"token": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"deletion_protection": {
			Type:     schema.TypeBool,
			Optional: true,
		},
		"bsu_optimized": {
			Type:     schema.TypeBool,
			Optional: true,
		},
		"image_id": {
			Type:     schema.TypeString,
			Required: true,
		},
		"shutdown_automatic_behavior": {
			Type:     schema.TypeBool,
			Optional: true,
		},
		"type": {
			Type:     schema.TypeString,
			Required: true,
		},
		"keypair_name": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"max_vms_count": {
			Type:     schema.TypeInt,
			Optional: true,
		},
		"min_vms_count": {
			Type:     schema.TypeInt,
			Optional: true,
		},
		"nics": {
			Type:     schema.TypeSet,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"delete_on_vm_deletion": {
						Type:     schema.TypeBool,
						Optional: true,
					},
					"description": {
						Type:     schema.TypeString,
						Optional: true,
					},
					"nic_sort_number": {
						Type:     schema.TypeInt,
						Optional: true,
					},
					"nic_id": {
						Type:     schema.TypeString,
						Optional: true,
					},
					"private_ip": {
						Type:     schema.TypeString,
						Optional: true,
					},
					"private_ips": {
						Type:     schema.TypeSet,
						Optional: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"primary_ip": {
									Type:     schema.TypeString,
									Optional: true,
								},
								"private_ip": {
									Type:     schema.TypeString,
									Optional: true,
								},
								"secondary_private_ip_count": {
									Type:     schema.TypeString,
									Optional: true,
								},
								"firewall_rules_set_id": {
									Type:     schema.TypeSet,
									Optional: true,
									Elem:     &schema.Schema{Type: schema.TypeString},
								},
								"subnet_id": {
									Type:     schema.TypeString,
									Optional: true,
								},
							},
						},
					},
				},
			},
		},
		"placement": {
			Type:     schema.TypeMap,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"affinity": {
						Type:     schema.TypeString,
						Optional: true,
					},
					"sub_region_name": {
						Type:     schema.TypeString,
						Optional: true,
					},
					"firewall_rules_set_name": {
						Type:     schema.TypeString,
						Optional: true,
					},
					"dedicated_host_id": {
						Type:     schema.TypeInt,
						Optional: true,
					},
					"tenancy": {
						Type:     schema.TypeString,
						Optional: true,
					},
				},
			},
		},
		"private_ip": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"private_ips": {
			Type:     schema.TypeSet,
			Optional: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
		"firewall_rules_set": {
			Type:     schema.TypeSet,
			Optional: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
		"firewall_rules_set_id": {
			Type:     schema.TypeString,
			Optional: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
		"subnet_id": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"user_data": {
			Type:     schema.TypeString,
			Optional: true,
		},
		//Attributes reference:
		"firewall_rules_sets": {
			Type:     schema.TypeSet,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"firewall_rules_set_id": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"firewall_rules_set_name": {
						Type:     schema.TypeString,
						Computed: true,
					},
				},
			},
		},
		"vms": {
			Type:     schema.TypeSet,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"launch_sort_number": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"architecture": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"block_device_mapping": {
						Type:     schema.TypeSet,
						Computed: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"device_name": {
									Type:     schema.TypeString,
									Computed: true,
								},
								"bsu": {
									Type:     schema.TypeMap,
									Computed: true,
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"delete_on_vm_deletion": {
												Type:     schema.TypeBool,
												Computed: true,
											},
											"state": {
												Type:     schema.TypeString,
												Computed: true,
											},
											"volume_id": {
												Type:     schema.TypeString,
												Computed: true,
											},
										},
									},
								},
							},
						},
					},
					"token": {
						Type:     schema.TypeString,
						Optional: true,
					},
					"public_dns_name": {
						Type:     schema.TypeString,
						Optional: true,
					},
					"bsu_optimized": {
						Type:     schema.TypeBool,
						Computed: true,
					},
					"firewall_rules_set": {
						Type:     schema.TypeSet,
						Computed: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"firewall_rules_set_id": {
									Type:     schema.TypeInt,
									Computed: true,
								},
								"firewall_rules_set_name": {
									Type:     schema.TypeString,
									Computed: true,
								},
							},
						},
					},
					"hypervisor": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"vm_profile": {
						Type: schema.TypeMap,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"resource_id": {
									Type:     schema.TypeString,
									Computed: true,
								},
								"vm_profile_id": {
									Type:     schema.TypeString,
									Computed: true,
								},
							},
						},
						Computed: true,
					},
					"image_id": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"vm_id": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"spot_vm": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"state": {
						Type: schema.TypeMap,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"state_code": {
									Type:     schema.TypeString,
									Computed: true,
								},
								"name": {
									Type:     schema.TypeString,
									Computed: true,
								},
							},
						},
						Computed: true,
					},
					"type": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"public_ip": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"kernel_id": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"keypair_name": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"monitoring": {
						Type: schema.TypeMap,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"state": {
									Type:     schema.TypeString,
									Computed: true,
								},
							},
						},
						Computed: true,
					},
					"nics": {
						Type:     schema.TypeSet,
						Computed: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"public_ip_link": {
									Type:     schema.TypeMap,
									Computed: true,
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"public_ip_account_id": {
												Type:     schema.TypeString,
												Computed: true,
											},
											"public_dns_name": {
												Type:     schema.TypeString,
												Computed: true,
											},
											"public_ip": {
												Type:     schema.TypeString,
												Computed: true,
											},
										},
									},
								},
								"nic_link": {
									Type: schema.TypeMap,
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"nic_link_id": {
												Type:     schema.TypeString,
												Computed: true,
											},
											"delete_on_vm_deletion": {
												Type:     schema.TypeBool,
												Computed: true,
											},
											"nic_sort_number": {
												Type:     schema.TypeInt,
												Computed: true,
											},
											"state": {
												Type:     schema.TypeString,
												Computed: true,
											},
										},
									},
									Computed: true,
								},
							},
						},
					},
					"description": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"firewall_rules_sets": {
						Type: schema.TypeSet,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"firewall_rules_set_id": {
									Type:     schema.TypeString,
									Computed: true,
								},
								"firewall_rules_set_name": {
									Type:     schema.TypeString,
									Computed: true,
								},
							},
						},
						Computed: true,
					},
					"mac_address": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"nic_id": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"account_id": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"private_dns_name": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"private_ip": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"private_ips": {
						Type:     schema.TypeSet,
						Computed: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"public_ip_link": {
									Type:     schema.TypeSet,
									Computed: true,
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"public_ip_account_id": {
												Type:     schema.TypeString,
												Computed: true,
											},
											"public_dns_name": {
												Type:     schema.TypeString,
												Computed: true,
											},
											"public_ip": {
												Type:     schema.TypeString,
												Computed: true,
											},
										},
									},
								},
								"primary_ip": {
									Type:     schema.TypeString,
									Computed: true,
								},
								"private_dns_name": {
									Type:     schema.TypeString,
									Computed: true,
								},
								"private_ip": {
									Type:     schema.TypeString,
									Computed: true,
								},
							},
						},
					},
					"activated_check": {
						Type:     schema.TypeBool,
						Computed: true,
					},
					"subnet_id": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"placement": {
						Type: schema.TypeMap,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"affinity": {
									Type:     schema.TypeString,
									Computed: true,
								},
								"sub_region_name": {
									Type:     schema.TypeString,
									Computed: true,
								},
								"firewall_rules_set_name": {
									Type:     schema.TypeString,
									Computed: true,
								},
								"dedicated_host_id": {
									Type:     schema.TypeString,
									Computed: true,
								},
								"tenancy": {
									Type:     schema.TypeString,
									Computed: true,
								},
							},
						},
						Computed: true,
					},
					"system": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"product_codes": {
						Type: schema.TypeSet,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"product_code": {
									Type:     schema.TypeString,
									Computed: true,
								},
								"type": {
									Type:     schema.TypeString,
									Computed: true,
								},
							},
						},
						Computed: true,
					},
					"ramdisk_id": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"comment": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"root_device_name": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"root_device_type": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"spot_vm_request_id": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"sriov_net_support": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"comments": {
						Type: schema.TypeMap,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"state_code": {
									Type:     schema.TypeInt,
									Computed: true,
								},
								"message": {
									Type:     schema.TypeString,
									Computed: true,
								},
							},
						},
						Computed: true,
					},
					"tags": {
						Type: schema.TypeSet,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"key": {
									Type:     schema.TypeString,
									Computed: true,
								},
								"value": {
									Type:     schema.TypeString,
									Computed: true,
								},
							},
						},
						Computed: true,
					},
					"virtualization_type": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"lin_id": {
						Type:     schema.TypeString,
						Computed: true,
					},
				},
			},
		},
		"account_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"requester_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"reservation_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"admin_password": {
			Type:     schema.TypeString,
			Computed: true,
		},
		//instance set is closed here
	}
}

type outscaleOApiInstanceOpts struct {
	BlockDeviceMappings               []*fcu.BlockDeviceMapping
	DisableAPITermination             *bool
	EBSOptimized                      *bool
	ImageID                           *string
	InstanceInitiatedShutdownBehavior *string
	InstanceType                      *string
	Ipv6AddressCount                  *int64
	KeyName                           *string
	NetworkInterfaces                 []*fcu.InstanceNetworkInterfaceSpecification
	Placement                         *fcu.Placement
	PrivateIPAddress                  *string
	SecurityGroupIDs                  []*string
	SecurityGroups                    []*string
	SubnetID                          *string
	UserData                          *string
	Monitoring                        *fcu.Monitoring
	// SpotPlacement                     *fcu.SpotPlacement
	// Ipv6Addresses                     []*fcu.InstanceIpv6Address
	// IAMInstanceProfile                *fcu.IamInstanceProfileSpecification
}

func buildOutscaleOAPIVMOpts(
	d *schema.ResourceData, meta interface{}) (*outscaleOApiInstanceOpts, error) {
	conn := meta.(*OutscaleClient).FCU

	opts := &outscaleOApiInstanceOpts{
		DisableAPITermination: aws.Bool(d.Get("deletion_protection").(bool)),
		EBSOptimized:          aws.Bool(d.Get("bsu_optimized").(bool)),
		ImageID:               aws.String(d.Get("image_id").(string)),
		InstanceType:          aws.String(d.Get("type").(string)),
	}

	if v := d.Get("shutdown_automatic_behavior").(string); v != "" {
		opts.InstanceInitiatedShutdownBehavior = aws.String(v)
	}

	userData := d.Get("user_data").(string)
	opts.UserData = &userData

	subnetID, hasSubnet := d.GetOk("subnet_id")
	if hasSubnet {
		s := subnetID.(string)
		opts.SubnetID = &s
	}

	if t, hasTenancy := d.GetOk("tenancy"); hasTenancy {
		opts.Placement.Tenancy = aws.String(t.(string))
	}

	gn, gnOk := d.GetOk("placement")

	if gnOk {
		opts.Placement = &fcu.Placement{
			GroupName: aws.String(gn.(string)),
		}
	}

	var groups []*string

	networkInterfaces, interfacesOk := d.GetOk("nics")
	if interfacesOk {
		opts.NetworkInterfaces = buildNetworkOApiInterfaceOpts(d, groups, networkInterfaces)
	}

	if v, ok := d.GetOk("private_ip"); ok {
		opts.PrivateIPAddress = aws.String(v.(string))
	}

	if v, ok := d.GetOk("keypair_name"); ok {
		opts.KeyName = aws.String(v.(string))
	}

	blockDevices, err := readBlockDeviceOApiMappingsFromConfig(d, conn)
	if err != nil {
		return nil, err
	}
	if len(blockDevices) > 0 {
		opts.BlockDeviceMappings = blockDevices
	}

	return opts, nil
}

func buildNetworkOApiInterfaceOpts(d *schema.ResourceData, groups []*string, nInterfaces interface{}) []*fcu.InstanceNetworkInterfaceSpecification {
	networkInterfaces := []*fcu.InstanceNetworkInterfaceSpecification{}
	// Get necessary items
	subnet, hasSubnet := d.GetOk("subnet_id")

	if hasSubnet {
		// If we have a non-default VPC / Subnet specified, we can flag
		// AssociatePublicIpAddress to get a Public IP assigned. By default these are not provided.
		// You cannot specify both SubnetId and the NetworkInterface.0.* parameters though, otherwise
		// you get: Network interfaces and an instance-level subnet ID may not be specified on the same request
		// You also need to attach Security Groups to the NetworkInterface instead of the instance,
		// to avoid: Network interfaces and an instance-level security groups may not be specified on
		// the same request
		ni := &fcu.InstanceNetworkInterfaceSpecification{
			DeviceIndex: aws.Int64(int64(0)),
			SubnetId:    aws.String(subnet.(string)),
			Groups:      groups,
		}

		if v, ok := d.GetOk("private_ip"); ok {
			ni.PrivateIpAddress = aws.String(v.(string))
		}

		networkInterfaces = append(networkInterfaces, ni)
	} else {
		// If we have manually specified network interfaces, build and attach those here.
		vL := nInterfaces.(*schema.Set).List()
		for _, v := range vL {
			ini := v.(map[string]interface{})
			ni := &fcu.InstanceNetworkInterfaceSpecification{
				DeviceIndex:         aws.Int64(int64(ini["nic_sort_number"].(int))),
				NetworkInterfaceId:  aws.String(ini["nic_id"].(string)),
				DeleteOnTermination: aws.Bool(ini["delete_on_vm_deletion"].(bool)),
			}
			networkInterfaces = append(networkInterfaces, ni)
		}
	}

	return networkInterfaces
}

func readBlockDeviceOApiMappingsFromConfig(
	d *schema.ResourceData, conn *fcu.Client) ([]*fcu.BlockDeviceMapping, error) {
	blockDevices := make([]*fcu.BlockDeviceMapping, 0)

	if v, ok := d.GetOk("bsu"); ok {
		vL := v.(*schema.Set).List()
		for _, v := range vL {
			bd := v.(map[string]interface{})
			ebs := &fcu.EbsBlockDevice{
				DeleteOnTermination: aws.Bool(bd["delete_on_vm_deletion"].(bool)),
			}
			if v, ok := bd["snapshot_id"].(string); ok && v != "" {
				ebs.SnapshotId = aws.String(v)
			}
			if v, ok := bd["volume_size"].(int); ok && v != 0 {
				ebs.VolumeSize = aws.Int64(int64(v))
			}
			if v, ok := bd["type"].(string); ok && v != "" {
				ebs.VolumeType = aws.String(v)
			}
			if v, ok := bd["iops"].(int); ok && v > 0 {
				ebs.Iops = aws.Int64(int64(v))
			}

			blockDevices = append(blockDevices, &fcu.BlockDeviceMapping{
				DeviceName:  aws.String(bd["device_name"].(string)),
				NoDevice:    aws.String(bd["no_device"].(string)),
				VirtualName: aws.String(bd["virtual_device_name"].(string)),
				Ebs:         ebs,
			})
		}
	}

	return blockDevices, nil
}

func InstanceStateOApiRefreshFunc(conn *fcu.Client, instanceID, failState string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		var resp *fcu.DescribeInstancesOutput
		var err error

		err = resource.Retry(30*time.Second, func() *resource.RetryError {
			resp, err = conn.VM.DescribeInstances(&fcu.DescribeInstancesInput{
				InstanceIds: []*string{aws.String(instanceID)},
			})
			return resource.RetryableError(err)
		})

		if err != nil {
			log.Printf("Error on InstanceStateRefresh: %s", err)

			return nil, "", err
		}

		if resp == nil || len(resp.Reservations) == 0 || len(resp.Reservations[0].Instances) == 0 {
			return nil, "", nil
		}

		i := resp.Reservations[0].Instances[0]
		state := *i.State.Name

		if state == failState {
			return i, state, fmt.Errorf("Failed to reach target state. Reason: %s",
				*i.StateReason)

		}

		return i, state, nil
	}
}

func InstanceOApiPa(conn *fcu.Client, instanceID, failState string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		var resp *fcu.DescribeInstancesOutput
		var err error

		err = resource.Retry(30*time.Second, func() *resource.RetryError {
			resp, err = conn.VM.DescribeInstances(&fcu.DescribeInstancesInput{
				InstanceIds: []*string{aws.String(instanceID)},
			})

			return resource.RetryableError(err)
		})

		if err != nil {
			log.Printf("Error on InstanceStateRefresh: %s", err)

			return nil, "", err
		}

		if resp == nil || len(resp.Reservations) == 0 || len(resp.Reservations[0].Instances) == 0 {
			return nil, "", nil
		}

		i := resp.Reservations[0].Instances[0]
		state := *i.State.Name

		if state == failState {
			return i, state, fmt.Errorf("Failed to reach target state. Reason: %s",
				*i.StateReason)

		}

		return i, state, nil
	}
}
