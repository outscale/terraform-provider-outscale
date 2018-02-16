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

func resourceOutscaleVM() *schema.Resource {
	return &schema.Resource{
		Create: resourceVMCreate,
		Read:   resourceVMRead,
		Update: resourceVMUpdate,
		Delete: resourceVMDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: getVMSchema(),
	}
}

func resourceVMCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	instanceOpts, err := buildOutscaleVMOpts(d, meta)
	if err != nil {
		return err
	}

	// Build the creation struct
	runOpts := &fcu.RunInstancesInput{
		BlockDeviceMappings:   instanceOpts.BlockDeviceMappings,
		DisableApiTermination: instanceOpts.DisableApiTermination,
		EbsOptimized:          instanceOpts.EBSOptimized,
		DryRun:                instanceOpts.DryRun,
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
		// RamdiskId:        instanceOpts.RamdiskId,
		SecurityGroupIds: instanceOpts.SecurityGroupIDs,
		SecurityGroups:   instanceOpts.SecurityGroups,
		SubnetId:         instanceOpts.SubnetID,
		UserData:         instanceOpts.UserData,
	}

	// Create the instance
	log.Printf("[DEBUG] Run configuration: %+v", runOpts)

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
		Refresh:    InstanceStateRefreshFunc(conn, *instance.InstanceId, "terminated"),
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

	return resourceVMRead(d, meta)
}

func resourceVMRead(d *schema.ResourceData, meta interface{}) error {
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

	d.Set("block_device_mapping", getBlockDeviceMapping(instance.BlockDeviceMappings))

	d.Set("client_token", instance.ClientToken)

	// d.Set("disable_api_termination", instance.DisableApiTermination)

	// d.Set("dry_run", instance.DryRun)

	d.Set("ebs_optimized", instance.EbsOptimized)

	d.Set("image_id", instance.ImageId)

	// d.Set("instance_initiated_shutdown_behavior", instance.InstanceInitiatedShutdownBehavior)

	d.Set("instance_type", instance.InstanceType)

	// d.Set("instance_name", instance.InstanceName)

	d.Set("key_name", instance.KeyName)

	// d.Set("max_count", instance.MaxCount)

	// d.Set("min_count", instance.MinCount)

	d.Set("network_interface", getNetworkInterfaceSet(instance.NetworkInterfaces))

	// d.Set("placement", getPlacement(instance.Placement))

	d.Set("private_ip", instance.PrivateIpAddress)

	// d.Set("private_ips", getPrivateIPAddressSet(instance.PrivateIpAddresses))

	d.Set("ramdisk_id", instance.RamdiskId)

	// d.Set("security_group", )
	// d.Set("security_group_id", )

	d.Set("subnet_id", instance.SubnetId)

	// d.Set("user_date",instance.UserData)

	d.Set("group_set", getGroupSet(resp.GroupSet))

	d.Set("owner_id", resp.OwnerId)

	d.Set("requester_id", resp.RequesterId)

	d.Set("reservation_id", resp.ReservationId)

	err = d.Set("instances_set", flattenedInstanceSet([]*fcu.Instance{instance}))
	if err != nil {
		return err
	}

	if instance.Platform != nil && *instance.Platform == "windows" {
		var passRes *fcu.GetPasswordDataOutput
		err = resource.Retry(500*time.Second, func() *resource.RetryError {
			var err error
			passRes, err = conn.VM.GetPasswordData(&fcu.GetPasswordDataInput{
				InstanceId: instance.InstanceId,
			})

			if err != nil {
				if strings.Contains(fmt.Sprint(err), "RequestLimitExceeded") {
					return resource.RetryableError(fmt.Errorf("Got empty password for instance (%s)", d.Id()))
				}
			}

			if passRes.PasswordData == nil || *passRes.PasswordData == "" {
				return resource.RetryableError(fmt.Errorf("Got empty password for instance (%s)", d.Id()))
			}

			return resource.NonRetryableError(err)
		})

		if passRes == nil {
			return fmt.Errorf("Error launching source instance: (%s)", d.Id())
		}

		if err != nil {
			return err
		}

		d.Set("password_data", passRes.PasswordData)

		return nil
	}

	return nil
}

func resourceVMUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU
	log.Printf("[DEBUG] updating the instance %s", d.Id())

	if d.HasChange("key_name") {
		input := &fcu.ModifyInstanceKeyPairInput{
			InstanceId: aws.String(d.Id()),
			KeyName:    aws.String(d.Get("key_name").(string)),
		}

		err := conn.VM.ModifyInstanceKeyPair(input)
		if err != nil {
			return err
		}
	}

	return resourceVMRead(d, meta)
}

func resourceVMDelete(d *schema.ResourceData, meta interface{}) error {
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

func getVMSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		// Attributes
		"block_device_mapping": {
			Type: schema.TypeSet,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"device_name": {
						Type:     schema.TypeString,
						Optional: true,
					},
					"ebs": {
						Type:     schema.TypeMap,
						Optional: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"delete_on_termination": {
									Type:     schema.TypeString,
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
								"volume_type": {
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
					"virtual_name": {
						Type:     schema.TypeString,
						Optional: true,
					},
				},
			},
			Optional: true,
		},

		"client_token": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"disable_api_termination": {
			Type:     schema.TypeBool,
			Optional: true,
		},
		"dry_run": {
			Type:     schema.TypeBool,
			Optional: true,
		},
		"ebs_optimized": {
			Type:     schema.TypeBool,
			Optional: true,
		},
		"image_id": {
			Type:     schema.TypeString,
			ForceNew: true,
			Required: true,
		},
		"instance_initiated_shutdown_behavior": {
			Type:     schema.TypeBool,
			Optional: true,
		},
		"instance_type": {
			Type:     schema.TypeString,
			ForceNew: true,
			Required: true,
		},
		"instance_name": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"key_name": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"max_count": {
			Type:     schema.TypeInt,
			Optional: true,
		},
		"min_count": {
			Type:     schema.TypeInt,
			Optional: true,
		},
		"network_interface": {

			Type:     schema.TypeSet,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"delete_on_termination": {
						Type:     schema.TypeBool,
						Optional: true,
					},
					"description": {
						Type:     schema.TypeString,
						Optional: true,
					},
					"device_index": {
						Type:     schema.TypeInt,
						Optional: true,
					},
					"network_interface_id": {
						Type:     schema.TypeInt,
						Optional: true,
					},
					"private_ip_address": {
						Type:     schema.TypeString,
						Optional: true,
					},
					"private_ip_addresses_set": {
						Type:     schema.TypeSet,
						Optional: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"primary": {
									Type:     schema.TypeString,
									Optional: true,
								},
								"private_ip_address": {
									Type:     schema.TypeString,
									Optional: true,
								},
							},
						},
					},
					"secondary_private_ip_address_count": {
						Type:     schema.TypeString,
						Optional: true,
					},
					"security_group_id": {
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
		"placement": {
			Type:     schema.TypeMap,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"affinity": {
						Type:     schema.TypeString,
						Optional: true,
					},
					"availability_zone": {
						Type:     schema.TypeString,
						Optional: true,
					},
					"group_name": {
						Type:     schema.TypeString,
						Optional: true,
					},
					"host_id": {
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
		"private_ip_address": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"private_ip_addresses": {
			Type:     schema.TypeSet,
			Optional: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
		"ramdisk_id": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"security_group": {
			Type:     schema.TypeSet,
			Optional: true,
			Computed: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
		"security_group_id": {
			Type:     schema.TypeSet,
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
		"group_set": {
			Type:     schema.TypeSet,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"group_id": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"group_name": {
						Type:     schema.TypeString,
						Computed: true,
					},
				},
			},
		},
		"instances_set": {
			Type:     schema.TypeSet,
			Computed: true,
			// Set:      resourceInstancSetHash,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"ami_launch_index": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"architecture": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"block_device_mapping": {
						Type:     schema.TypeList,
						Computed: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"device_name": {
									Type:     schema.TypeString,
									Computed: true,
								},
								"ebs": {
									Type:     schema.TypeMap,
									Computed: true,
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"delete_on_termination": {
												Type:     schema.TypeBool,
												Computed: true,
											},
											"status": {
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
					"client_token": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"dns_name": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"ebs_optimised": {
						Type:     schema.TypeBool,
						Computed: true,
					},
					"group_set": {
						Type:     schema.TypeList,
						Computed: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"group_id": {
									Type:     schema.TypeString,
									Computed: true,
								},
								"group_name": {
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
					"iam_instance_profile": {
						Type: schema.TypeMap,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"arn": {
									Type:     schema.TypeString,
									Computed: true,
								},
								"id": {
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
					"instance_id": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"instance_state": {
						Type: schema.TypeMap,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"code": {
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
					"instance_type": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"ip_address": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"kernel_id": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"key_name": {
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
					"network_interface_set": {
						Type:     schema.TypeList,
						Computed: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"association": {
									Type:     schema.TypeMap,
									Computed: true,
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"ip_owner_id": {
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
								"attachment": {
									Type: schema.TypeMap,
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"attachement_id": {
												Type:     schema.TypeString,
												Computed: true,
											},
											"delete_on_termination": {
												Type:     schema.TypeBool,
												Computed: true,
											},
											"device_index": {
												Type:     schema.TypeInt,
												Computed: true,
											},
											"status": {
												Type:     schema.TypeString,
												Computed: true,
											},
										},
									},
									Computed: true,
								},
								"description": {
									Type:     schema.TypeString,
									Computed: true,
								},
								"group_set": {
									Type:     schema.TypeList,
									Computed: true,
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"group_id": {
												Type:     schema.TypeString,
												Computed: true,
											},
											"group_name": {
												Type:     schema.TypeString,
												Computed: true,
											},
										},
									},
								},
								"mac_address": {
									Type:     schema.TypeString,
									Computed: true,
								},
								"network_interface_id": {
									Type:     schema.TypeString,
									Computed: true,
								},
								"owner_id": {
									Type:     schema.TypeString,
									Computed: true,
								},
								"private_dns_name": {
									Type:     schema.TypeString,
									Computed: true,
								},
								"private_ip_address": {
									Type:     schema.TypeString,
									Computed: true,
								},
								"private_ip_addresses_set": {
									Type:     schema.TypeList,
									Computed: true,
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"association": {
												Type:     schema.TypeMap,
												Computed: true,
												Elem: &schema.Resource{
													Schema: map[string]*schema.Schema{
														"ip_owner_id": {
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
											"primary": {
												Type:     schema.TypeBool,
												Computed: true,
											},
											"private_dns_name": {
												Type:     schema.TypeString,
												Computed: true,
											},
											"private_ip_address": {
												Type:     schema.TypeString,
												Computed: true,
											},
										},
									},
								},
								"source_dest_check": {
									Type:     schema.TypeBool,
									Computed: true,
								},
								"status": {
									Type:     schema.TypeString,
									Computed: true,
								},
								"subnet_id": {
									Type:     schema.TypeString,
									Computed: true,
								},
								"vpc_id": {
									Type:     schema.TypeInt,
									Computed: true,
								},
							},
						},
					},
					"placement": {
						Type: schema.TypeMap,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"affinity": {
									Type:     schema.TypeString,
									Computed: true,
								},
								"availability_zone": {
									Type:     schema.TypeString,
									Computed: true,
								},
								"group_name": {
									Type:     schema.TypeString,
									Computed: true,
								},
								"host_id": {
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
					"platform": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"private_dns_name": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"private_ip_address": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"product_codes": {
						Type: schema.TypeList,
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
					"reason": {
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
					"source_dest_check": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"spot_instance_request_id": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"sriov_net_support": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"state_reason": {
						Type: schema.TypeMap,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"code": {
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
					"subnet_id": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"tag_set": {
						Type: schema.TypeList,
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
					"vpc_id": {
						Type:     schema.TypeString,
						Computed: true,
					},
				},
			},
		},
		"owner_id": {
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
		"password_data": {
			Type:     schema.TypeString,
			Computed: true,
		},
		//instance set is closed here
	}
}

type outscaleInstanceOpts struct {
	BlockDeviceMappings               []*fcu.BlockDeviceMapping
	DisableApiTermination             *bool
	EBSOptimized                      *bool
	DryRun                            *bool
	ImageID                           *string
	InstanceInitiatedShutdownBehavior *bool
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
	RamdiskId                         *string
	RequesterId                       *string
	ReservationId                     *string
	PasswordData                      *string
	OwnerId                           *string
	// Monitoring                        *fcu.RunInstancesMonitoringEnabled
	// SpotPlacement                     *fcu.SpotPlacement
	// Ipv6Addresses                     []*fcu.InstanceIpv6Address
	// IAMInstanceProfile                *fcu.IamInstanceProfileSpecification
}

func buildOutscaleVMOpts(
	d *schema.ResourceData, meta interface{}) (*outscaleInstanceOpts, error) {
	conn := meta.(*OutscaleClient).FCU

	opts := &outscaleInstanceOpts{
		DisableApiTermination: aws.Bool(d.Get("disable_api_termination").(bool)),
		EBSOptimized:          aws.Bool(d.Get("ebs_optimized").(bool)),
		ImageID:               aws.String(d.Get("image_id").(string)),
		InstanceType:          aws.String(d.Get("instance_type").(string)),
	}

	if v := d.Get("instance_initiated_shutdown_behavior").(bool); v {
		opts.InstanceInitiatedShutdownBehavior = aws.Bool(v)
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

	az, azOk := d.GetOk("availability_zone")
	gn, gnOk := d.GetOk("placement_group")

	if azOk && gnOk {
		opts.Placement = &fcu.Placement{
			AvailabilityZone: aws.String(az.(string)),
			GroupName:        aws.String(gn.(string)),
		}
	}

	var groups []*string
	if v := d.Get("security_group"); v != nil {

		sgs := v.(*schema.Set).List()
		for _, v := range sgs {
			str := v.(string)
			groups = append(groups, aws.String(str))
		}
	}

	opts.SecurityGroups = groups

	networkInterfaces, interfacesOk := d.GetOk("network_interface")
	if interfacesOk {
		opts.NetworkInterfaces = buildNetworkInterfaceOpts(d, groups, networkInterfaces)

	}

	if v, ok := d.GetOk("private_ip"); ok {
		opts.PrivateIPAddress = aws.String(v.(string))
	}

	if v, ok := d.GetOk("vpc_security_group_ids"); ok && v.(*schema.Set).Len() > 0 {
		for _, v1 := range v.(*schema.Set).List() {
			opts.SecurityGroupIDs = append(opts.SecurityGroupIDs, aws.String(v1.(string)))
		}
	}

	if v, ok := d.GetOk("ipv6_address_count"); ok {
		opts.Ipv6AddressCount = aws.Int64(int64(v.(int)))
	}

	if v, ok := d.GetOk("key_name"); ok {
		opts.KeyName = aws.String(v.(string))
	}

	blockDevices, err := readBlockDeviceMappingsFromConfig(d, conn)
	if err != nil {
		return nil, err
	}
	if len(blockDevices) > 0 {
		opts.BlockDeviceMappings = blockDevices
	}

	if dryRun, ok := d.GetOk("dry_run"); ok {
		opts.DryRun = aws.Bool(dryRun.(bool))
	}

	ramdiskID := d.Get("ramdisk_id").(string)
	opts.RamdiskId = &ramdiskID

	ownerID := d.Get("owner_id").(string)
	opts.OwnerId = &ownerID

	requesterID := d.Get("requester_id").(string)
	opts.RequesterId = &requesterID

	reservationID := d.Get("reservation_id").(string)
	opts.ReservationId = &reservationID

	if p := d.Get("password_data"); p != nil {
		opts.PasswordData = aws.String(p.(string))
	} else {
		opts.PasswordData = aws.String("pending")
	}

	return opts, nil
}

func buildNetworkInterfaceOpts(d *schema.ResourceData, groups []*string, nInterfaces interface{}) []*fcu.InstanceNetworkInterfaceSpecification {
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

		if v, ok := d.GetOkExists("associate_public_ip_address"); ok {
			ni.AssociatePublicIpAddress = aws.Bool(v.(bool))
		}

		if v, ok := d.GetOk("private_ip"); ok {
			ni.PrivateIpAddress = aws.String(v.(string))
		}

		if v, ok := d.GetOk("ipv6_address_count"); ok {
			ni.Ipv6AddressCount = aws.Int64(int64(v.(int)))
		}

		if v := d.Get("vpc_security_group_ids").(*schema.Set); v.Len() > 0 {
			for _, v := range v.List() {
				ni.Groups = append(ni.Groups, aws.String(v.(string)))
			}
		}

		networkInterfaces = append(networkInterfaces, ni)
	} else {
		// If we have manually specified network interfaces, build and attach those here.
		vL := nInterfaces.(*schema.Set).List()
		for _, v := range vL {
			ini := v.(map[string]interface{})
			ni := &fcu.InstanceNetworkInterfaceSpecification{
				DeviceIndex:         aws.Int64(int64(ini["device_index"].(int))),
				NetworkInterfaceId:  aws.String(ini["network_interface_id"].(string)),
				DeleteOnTermination: aws.Bool(ini["delete_on_termination"].(bool)),
			}
			networkInterfaces = append(networkInterfaces, ni)
		}
	}

	return networkInterfaces
}

func readBlockDeviceMappingsFromConfig(
	d *schema.ResourceData, conn *fcu.Client) ([]*fcu.BlockDeviceMapping, error) {
	blockDevices := make([]*fcu.BlockDeviceMapping, 0)

	if v, ok := d.GetOk("ebs_block_device"); ok {
		vL := v.(*schema.Set).List()
		for _, v := range vL {
			bd := v.(map[string]interface{})
			ebs := &fcu.EbsBlockDevice{
				DeleteOnTermination: aws.Bool(bd["delete_on_termination"].(bool)),
			}

			if v, ok := bd["snapshot_id"].(string); ok && v != "" {
				ebs.SnapshotId = aws.String(v)
			}

			if v, ok := bd["encrypted"].(bool); ok && v {
				ebs.Encrypted = aws.Bool(v)
			}

			if v, ok := bd["volume_size"].(int); ok && v != 0 {
				ebs.VolumeSize = aws.Int64(int64(v))
			}

			if v, ok := bd["volume_type"].(string); ok && v != "" {
				ebs.VolumeType = aws.String(v)

				if v, ok := bd["iops"].(int); ok && v > 0 {
					ebs.Iops = aws.Int64(int64(v))

				}

			}

			blockDevices = append(blockDevices, &fcu.BlockDeviceMapping{
				DeviceName: aws.String(bd["device_name"].(string)),
				Ebs:        ebs,
			})
		}
	}

	if v, ok := d.GetOk("ephemeral_block_device"); ok {
		vL := v.(*schema.Set).List()
		for _, v := range vL {
			bd := v.(map[string]interface{})
			bdm := &fcu.BlockDeviceMapping{
				DeviceName:  aws.String(bd["device_name"].(string)),
				VirtualName: aws.String(bd["virtual_name"].(string)),
			}
			if v, ok := bd["no_device"].(bool); ok && v {
				bdm.NoDevice = aws.String("")
				// When NoDevice is true, just ignore VirtualName since it's not needed
				bdm.VirtualName = nil
			}

			if bdm.NoDevice == nil && aws.StringValue(bdm.VirtualName) == "" {
				return nil, errors.New("virtual_name cannot be empty when no_device is false or undefined.")
			}

			blockDevices = append(blockDevices, bdm)
		}
	}

	if v, ok := d.GetOk("root_block_device"); ok {
		vL := v.([]interface{})
		if len(vL) > 1 {
			return nil, errors.New("Cannot specify more than one root_block_device.")
		}
		for _, v := range vL {
			bd := v.(map[string]interface{})
			ebs := &fcu.EbsBlockDevice{
				DeleteOnTermination: aws.Bool(bd["delete_on_termination"].(bool)),
			}

			if v, ok := bd["volume_size"].(int); ok && v != 0 {
				ebs.VolumeSize = aws.Int64(int64(v))
			}

			if v, ok := bd["volume_type"].(string); ok && v != "" {
				ebs.VolumeType = aws.String(v)
			}

			if v, ok := bd["iops"].(int); ok && v > 0 && *ebs.VolumeType == "io1" {
				// Only set the iops attribute if the volume type is io1. Setting otherwise
				// can trigger a refresh/plan loop based on the computed value that is given
				// from AWS, and prevent us from specifying 0 as a valid iops.
				//   See https://github.com/hashicorp/terraform/pull/4146
				//   See https://github.com/hashicorp/terraform/issues/7765
				ebs.Iops = aws.Int64(int64(v))
			} else if v, ok := bd["iops"].(int); ok && v > 0 && *ebs.VolumeType != "io1" {
				// Message user about incompatibility
				log.Print("[WARN] IOPs is only valid for storate type io1 for EBS Volumes")
			}
		}
	}

	return blockDevices, nil
}

func InstanceStateRefreshFunc(conn *fcu.Client, instanceID, failState string) resource.StateRefreshFunc {
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

// GetInstanceGetPasswordData func
func GetInstanceGetPasswordData(conn *fcu.Client, instanceID, failState string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		var resp *fcu.GetPasswordDataOutput
		var err error

		err = resource.Retry(30*time.Second, func() *resource.RetryError {
			resp, err = conn.VM.GetPasswordData(&fcu.GetPasswordDataInput{
				InstanceId: aws.String(instanceID),
			})
			return resource.RetryableError(err)
		})

		if err != nil {
			log.Printf("Error on InstanceStateRefresh: %s", err)

			return nil, "", err
		}

		if resp == nil {
			return nil, "", nil
		}

		i := resp.PasswordData

		if len(*i) < 0 {
			return nil, "running", nil
		} else {
			return nil, "terminated", nil
		}

		return nil, "", nil
	}
}

func InstancePa(conn *fcu.Client, instanceID, failState string) resource.StateRefreshFunc {
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

func getInstanceSet(instance *fcu.Instance) *schema.Set {

	instanceSet := map[string]interface{}{}
	s := schema.NewSet(nil, []interface{}{})

	instanceSet["ami_launch_index"] = *instance.AmiLaunchIndex
	instanceSet["ebs_optimised"] = *instance.EbsOptimized
	instanceSet["architecture"] = *instance.Architecture
	instanceSet["client_token"] = *instance.ClientToken
	instanceSet["hypervisor"] = *instance.Hypervisor
	instanceSet["image_id"] = *instance.ImageId
	instanceSet["instance_id"] = *instance.InstanceId
	instanceSet["instance_type"] = *instance.InstanceType
	instanceSet["kernel_id"] = *instance.KernelId
	instanceSet["key_name"] = *instance.KeyName
	instanceSet["private_dns_name"] = *instance.PrivateDnsName
	instanceSet["private_ip_address"] = *instance.PrivateIpAddress
	instanceSet["root_device_name"] = *instance.RootDeviceName

	if instance.DnsName != nil {
		instanceSet["dns_name"] = *instance.DnsName
	}
	if instance.IpAddress != nil {
		instanceSet["ip_address"] = *instance.IpAddress
	}
	if instance.Platform != nil {
		instanceSet["platform"] = *instance.Platform
	}
	if instance.RamdiskId != nil {
		instanceSet["ramdisk_id"] = *instance.RamdiskId
	}
	if instance.Reason != nil {
		instanceSet["reason"] = *instance.Reason
	}
	if instance.SourceDestCheck != nil {
		instanceSet["source_dest_check"] = *instance.SourceDestCheck
	}
	if instance.SpotInstanceRequestId != nil {
		instanceSet["spot_instance_request_id"] = *instance.SpotInstanceRequestId
	}
	if instance.SriovNetSupport != nil {
		instanceSet["sriov_net_support"] = *instance.SriovNetSupport
	}
	if instance.SubnetId != nil {
		instanceSet["subnet_id"] = *instance.SubnetId
	}
	if instance.VirtualizationType != nil {
		instanceSet["virtualization_type"] = *instance.VirtualizationType
	}
	if instance.VpcId != nil {
		instanceSet["vpc_id"] = *instance.VpcId
	}

	s.Add(instanceSet)

	instanceSet["block_device_mapping"] = getBlockDeviceMapping(instance.BlockDeviceMappings)
	// instanceSet["group_set"] = getGroupSet(instance.GroupSet)
	// instanceSet["iam_instance_profile"] = getIAMInstanceProfile(instance.IamInstanceProfile)
	// instanceSet["instance_state"] = getInstanceState(instance.State)
	// instanceSet["monitoring"] = getMonitoring(instance.Monitoring)
	// instanceSet["network_interface_set"] = getNetworkInterfaceSet(instance.NetworkInterfaces)
	// instanceSet["placement"] = getPlacement(instance.Placement)
	// instanceSet["state_reason"] = getStateReason(instance.StateReason)
	// instanceSet["product_codes"] = getProductCodes(instance.ProductCodes)
	// instanceSet["tag_set"] = getTagSet(instance.Tags)

	return s
}
