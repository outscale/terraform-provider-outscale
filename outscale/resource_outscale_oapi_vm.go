package outscale

import (
	"errors"
	"fmt"
	"log"
	"reflect"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-outscale/osc/oapi"
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
	conn := meta.(*OutscaleClient).OAPI

	instanceOpts, err := buildOutscaleOAPIVMOpts(d, meta)
	if err != nil {
		return err
	}

	// Build the creation struct
	runOpts := &oapi.CreateVmsRequest{
		BlockDeviceMappings:         instanceOpts.BlockDeviceMappings,
		BsuOptimized:                instanceOpts.EBSOptimized,
		VmType:                      instanceOpts.InstanceType,
		Nics:                        instanceOpts.NetworkInterfaces,
		ImageId:                     instanceOpts.ImageID,
		SubnetId:                    instanceOpts.SubnetID,
		UserData:                    instanceOpts.UserData,
		Placement:                   instanceOpts.Placement,
		KeypairName:                 instanceOpts.KeyName,
		MaxVmsCount:                 int64(1),
		MinVmsCount:                 int64(1),
		SecurityGroupIds:            instanceOpts.SecurityGroupIDs,
		SecurityGroups:              instanceOpts.SecurityGroups,
		VmInitiatedShutdownBehavior: instanceOpts.InstanceInitiatedShutdownBehavior,
		DeletionProtection:          instanceOpts.DisableAPITermination,
		// Monitoring:            instanceOpts.Monitoring,
		//DisableApiTermination: instanceOpts.DisableAPITermination,

		// IamInstanceProfile:    instanceOpts.IAMInstanceProfile,
		// Ipv6AddressCount:                  instanceOpts.Ipv6AddressCount,
		// Ipv6Addresses:                     instanceOpts.Ipv6Addresses,
		// PrivateIpAddress:                  instanceOpts.PrivateIPAddress,
	}

	//Missing on Swagger Spec
	// tagsSpec := make([]*oapi.TagSpecification, 0)

	// if v, ok := d.GetOk("tags"); ok {
	// 	tags := tagsFromMap(v.(map[string]interface{}))

	// 	spec := &oapi.TagSpecification{
	// 		ResourceType: aws.String("instance"),
	// 		Tags:         tags,
	// 	}

	// 	tagsSpec = append(tagsSpec, spec)
	// }

	// if len(tagsSpec) > 0 {
	// 	runOpts.TagSpecifications = tagsSpec
	// }

	// Create the instance
	var runResp *oapi.CreateVmsResponse
	var resp *oapi.POST_CreateVmsResponses
	err = resource.Retry(30*time.Second, func() *resource.RetryError {
		var err error
		resp, err = conn.POST_CreateVms(*runOpts)

		if err != nil {
			if strings.Contains(fmt.Sprint(err), "Throttling") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("Error launching source instance: %s", err)
	}

	runResp = resp.OK

	if runResp == nil || len(runResp.Vms) == 0 {
		return errors.New("Error launching source instance: no instances returned in response")
	}

	vm := runResp.Vms[0]
	fmt.Printf("[INFO] Instance ID: %s", vm.VmId)

	d.SetId(vm.VmId)

	if d.IsNewResource() {
		if err := setOAPITags(conn, d); err != nil {
			return err
		}
		d.SetPartial("tag")
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"pending"},
		Target:     []string{"running"},
		Refresh:    InstanceStateOApiRefreshFunc(conn, vm.VmId, "terminated"),
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
	if vm.PublicIp != "" {
		d.SetConnInfo(map[string]string{
			"vm_type": "ssh",
			"host":    vm.PublicIp,
		})
	} else if vm.PrivateIp != "" {
		d.SetConnInfo(map[string]string{
			"vm_type": "ssh",
			"host":    vm.PrivateIp,
		})
	}

	return resourceOAPIVMRead(d, meta)
}

func resourceOAPIVMRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OAPI
	filters := oapi.FiltersVm{
		VmIds: []string{d.Id()},
	}

	input := &oapi.ReadVmsRequest{
		Filters: filters,
	}

	var resp *oapi.ReadVmsResponse
	var rs *oapi.POST_ReadVmsResponses
	var err error

	err = resource.Retry(30*time.Second, func() *resource.RetryError {
		rs, err = conn.POST_ReadVms(*input)

		return resource.RetryableError(err)
	})

	if err != nil {
		return fmt.Errorf("Error reading the VM %s", err)
	}

	resp = rs.OK

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
	if len(resp.Vms) == 0 {
		d.SetId("")
		return nil
	}

	instance := resp.Vms[0]

	//d.Set("block_device_mapping", getOAPIVMBlockDeviceMapping(instance.BlockDeviceMappings))
	d.Set("token", instance.ClientToken)
	d.Set("bsu_optimized", instance.BsuOptimized)
	d.Set("image_id", instance.ImageId)
	d.Set("vm_type", instance.VmType)
	d.Set("vm_id", instance.VmId)
	d.Set("keypair_name", instance.KeypairName)
	d.Set("nics", getOAPIVMNetworkInterfaceSet(instance.Nics))
	d.Set("private_ip", instance.PrivateIp)
	//ramdisk
	d.Set("subnet_id", instance.SubnetId)
	//tagSet
	//d.Set("account_id", "")
	d.Set("reservation_id", instance.ReservationId)

	if err := d.Set("group_set", getOAPISecurityGroups(instance.SecurityGroups)); err != nil {
		return err
	}

	placement := make(map[string]interface{})
	if !reflect.DeepEqual(instance.Placement, oapi.Placement{}) {
		placement["tenancy"] = instance.Placement.Tenancy
		placement["sub_region_name"] = instance.Placement.SubregionName
		//Missing on swagger spec
		//placement["affinity"] = instance.Placement.Affinity
		//placement["dedicated_host_id"] = instance.Placement.DedicatedHostId
		// "firewall_rules_set_name": instance.Placement.FirewallRulesSetName,
	}

	d.Set("placement", placement)

	//d.Set("delete_protection", instance.DeletionProtection)
	// d.Set("shutdown_automatic_behavior", instance.SpotInstanceRequestId)
	// d.Set("max_vms_count", instance)
	// d.Set("min_vms_count", instance.KernelId)
	// d.Set("private_ips", ips)
	// d.Set("security_groups", ips)
	// d.Set("security_group_ids", ips)
	// d.Set("subnet_id", ips)
	// d.Set("user_data", ips)

	return nil
}

func resourceOAPIVMUpdate(d *schema.ResourceData, meta interface{}) error {
	//conn := meta.(*OutscaleClient).OAPI
	fmt.Printf("[DEBUG] updating the instance %s", d.Id())

	d.Partial(true)

	// if d.HasChange("keypair_name") {
	// 	input := &oapi.UpdateKeypairRequest{
	// 		//VmId:        aws.String(d.Id()), Missing on Swagger Spec
	// 		KeypairName: aws.String(d.Get("keypair_name").(string)),
	// 	}

	// 	_, err := conn.POST_UpdateKeypair(*input)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	//Missing Tag_set

	// if d.HasChange("vm_type") && !d.IsNewResource() {
	// 	opts := &oapi.UpdateVmAttributeRequest{
	// 		VmId: aws.String(d.Id()),
	// 		Type: aws.String(d.Get("vm_type").(string)),
	// 	}
	// 	if err := updateVMAttr(conn, opts, "vm_type"); err != nil {
	// 		return err
	// 	}
	// }

	// if d.HasChange("user_data") && !d.IsNewResource() {
	// 	opts := &fcu.ModifyInstanceAttributeInput{
	// 		InstanceId: aws.String(d.Id()),
	// 		UserData: &fcu.BlobAttributeValue{
	// 			Value: d.Get("user_data").([]byte),
	// 		},
	// 	}
	// 	if err := modifyInstanceAttr(conn, opts, "user_data"); err != nil {
	// 		return err
	// 	}
	// }

	// if d.HasChange("ebs_optimized") && !d.IsNewResource() {
	// 	opts := &fcu.ModifyInstanceAttributeInput{
	// 		InstanceId: aws.String(d.Id()),
	// 		EbsOptimized: &fcu.AttributeBooleanValue{
	// 			Value: aws.Bool(d.Get("ebs_optimized").(bool)),
	// 		},
	// 	}
	// 	if err := modifyInstanceAttr(conn, opts, "ebs_optimized"); err != nil {
	// 		return err
	// 	}
	// }

	// if d.HasChange("delete_on_termination") && !d.IsNewResource() {
	// 	opts := &fcu.ModifyInstanceAttributeInput{
	// 		InstanceId: aws.String(d.Id()),
	// 		DeleteOnTermination: &fcu.AttributeBooleanValue{
	// 			Value: d.Get("delete_on_termination").(*bool),
	// 		},
	// 	}
	// 	if err := modifyInstanceAttr(conn, opts, "delete_on_termination"); err != nil {
	// 		return err
	// 	}
	// }

	// if d.HasChange("disable_api_termination") {
	// 	opts := &fcu.ModifyInstanceAttributeInput{
	// 		InstanceId: aws.String(d.Id()),
	// 		DisableApiTermination: &fcu.AttributeBooleanValue{
	// 			Value: aws.Bool(d.Get("disable_api_termination").(bool)),
	// 		},
	// 	}
	// 	if err := modifyInstanceAttr(conn, opts, "disable_api_termination"); err != nil {
	// 		return err
	// 	}
	// }

	// if d.HasChange("instance_initiated_shutdown_behavior") {
	// 	opts := &fcu.ModifyInstanceAttributeInput{
	// 		InstanceId: aws.String(d.Id()),
	// 		InstanceInitiatedShutdownBehavior: &fcu.AttributeValue{
	// 			Value: aws.String(d.Get("instance_initiated_shutdown_behavior").(string)),
	// 		},
	// 	}
	// 	if err := modifyInstanceAttr(conn, opts, "instance_initiated_shutdown_behavior"); err != nil {
	// 		return err
	// 	}
	// }

	// if d.HasChange("group_set") {
	// 	opts := &fcu.ModifyInstanceAttributeInput{
	// 		InstanceId: aws.String(d.Id()),
	// 		Groups:     d.Get("group_set").([]*string),
	// 	}
	// 	if err := modifyInstanceAttr(conn, opts, "group_set"); err != nil {
	// 		return err
	// 	}
	// }

	// if d.HasChange("source_dest_check") {
	// 	opts := &fcu.ModifyInstanceAttributeInput{
	// 		InstanceId: aws.String(d.Id()),
	// 		SourceDestCheck: &fcu.AttributeBooleanValue{
	// 			Value: aws.Bool(d.Get("source_dest_check").(bool)),
	// 		},
	// 	}
	// 	if err := modifyInstanceAttr(conn, opts, "source_dest_check"); err != nil {
	// 		return err
	// 	}
	// }

	// if d.HasChange("block_device_mapping") {
	// 	if err := setBlockDevice(d.Get("block_device_mapping"), conn, d.Id()); err != nil {
	// 		return err
	// 	}
	// }

	return resourceVMRead(d, meta)
}

func resourceOAPIVMDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OAPI

	id := d.Id()

	fmt.Printf("[INFO] Terminating instance: %s", id)
	req := &oapi.DeleteVmsRequest{
		VmIds: []string{id},
	}

	var err error
	err = resource.Retry(30*time.Second, func() *resource.RetryError {
		_, err = conn.POST_DeleteVms(*req)

		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded") {
				fmt.Printf("[INFO] Request limit exceeded")
				return resource.RetryableError(err)
			}
		}

		return resource.RetryableError(err)
	})

	if err != nil {
		return fmt.Errorf("Error deleting the instance")
	}

	fmt.Printf("[DEBUG] Waiting for instance (%s) to become terminated", id)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"pending", "running", "shutting-down", "stopped", "stopping"},
		Target:     []string{"terminated"},
		Refresh:    InstanceStateOApiRefreshFunc(conn, id, ""),
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
								"vm_type": {
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
			Computed: true,
			Optional: true,
		},
		"bsu_optimized": {
			Type:     schema.TypeBool,
			Optional: true,
		},
		"image_id": {
			Type:     schema.TypeString,
			ForceNew: true,
			Required: true,
		},
		"shutdown_automatic_behavior": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"vm_type": {
			Type:     schema.TypeString,
			ForceNew: true,
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
			Type: schema.TypeSet,
			//To change in for oapi attributes ConflictsWith: []string{"subnet_id", "security_group_id", "security_group"},
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
									Type:     schema.TypeBool,
									Optional: true,
								},
								"private_ip": {
									Type:     schema.TypeString,
									Optional: true,
								},
							},
						},
					},
					"secondary_private_ip_count": {
						Type:     schema.TypeString,
						Optional: true,
					},
					"security_group_ids": {
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
		"security_groups": {
			Type:     schema.TypeSet,
			Optional: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
		"security_group_ids": {
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
					"security_group_ids": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"security_group_name": {
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
					"security_groups": {
						Type:     schema.TypeSet,
						Computed: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"security_group_ids": {
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
					"vm_type": {
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
					"group_set": {
						Type: schema.TypeSet,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"security_group_ids": {
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
								"vm_type": {
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
	BlockDeviceMappings               []oapi.BlockDeviceMappingVmCreation
	DisableAPITermination             bool
	EBSOptimized                      bool
	ImageID                           string
	InstanceInitiatedShutdownBehavior string
	InstanceType                      string
	Ipv6AddressCount                  int64
	KeyName                           string
	NetworkInterfaces                 []oapi.NicForVmCreation
	Placement                         oapi.Placement
	PrivateIPAddress                  string
	SecurityGroupIDs                  []string
	SecurityGroups                    []string
	SubnetID                          string
	UserData                          string
	// Monitoring                        oapi.Monitoring
	// SpotPlacement                     oapi.SpotPlacement
	// Ipv6Addresses                     []oapi.InstanceIpv6Address
	// IAMInstanceProfile                oapi.IamInstanceProfileSpecification
}

func buildOutscaleOAPIVMOpts(
	d *schema.ResourceData, meta interface{}) (*outscaleOApiInstanceOpts, error) {
	conn := meta.(*OutscaleClient).OAPI

	opts := &outscaleOApiInstanceOpts{
		DisableAPITermination: d.Get("deletion_protection").(bool),
		EBSOptimized:          d.Get("bsu_optimized").(bool),
		ImageID:               d.Get("image_id").(string),
		InstanceType:          d.Get("vm_type").(string),
	}

	if v := d.Get("shutdown_automatic_behavior").(string); v != "" {
		opts.InstanceInitiatedShutdownBehavior = v
	}

	userData := d.Get("user_data").(string)
	opts.UserData = userData

	subnetID, hasSubnet := d.GetOk("subnet_id")

	tenancy, tenancyOK := d.GetOk("tenancy")
	az, azOk := d.GetOk("availability_zone")
	//gn, gnOk := d.GetOk("placement")

	//if gnOk && tenancyOK && azOk {
	if tenancyOK && azOk {
		opts.Placement = oapi.Placement{
			//PlacementName: gn.(string),
			SubregionName: az.(string),
			Tenancy:       tenancy.(string),
		}
	}

	groups := make([]string, 0)
	if v := d.Get("security_group"); v != nil {
		groups = expandStringValueList(v.(*schema.Set).List())
		if len(groups) > 0 && hasSubnet {
			log.Print("[WARN] Deprecated. Attempting to use 'security_group' within a VPC instance. Use 'security_group_id' instead.")
		}
	}

	networkInterfaces, interfacesOk := d.GetOk("nics")
	if hasSubnet || interfacesOk {
		opts.NetworkInterfaces = buildNetworkOApiInterfaceOpts(d, groups, networkInterfaces)
	} else {
		if hasSubnet {
			s := subnetID.(string)
			opts.SubnetID = s
		}

		if opts.SubnetID != "" {
			opts.SecurityGroupIDs = groups
		} else {
			opts.SecurityGroups = groups
		}

		var groupIDs []string
		if v := d.Get("security_group_ids"); v != nil {

			sgs := v.(*schema.Set).List()
			for _, v := range sgs {
				str := v.(string)
				groupIDs = append(groupIDs, str)
			}
		}
		opts.SecurityGroupIDs = groupIDs
	}

	if v, ok := d.GetOk("private_ip"); ok {
		opts.PrivateIPAddress = v.(string)
	}

	if v, ok := d.GetOk("keypair_name"); ok {
		opts.KeyName = v.(string)
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

func buildNetworkOApiInterfaceOpts(d *schema.ResourceData, groups []string, nInterfaces interface{}) []oapi.NicForVmCreation {
	networkInterfaces := []oapi.NicForVmCreation{}
	subnet, hasSubnet := d.GetOk("subnet_id")

	if hasSubnet {
		ni := oapi.NicForVmCreation{
			DeviceNumber:     int64(0),
			SubnetId:         subnet.(string),
			SecurityGroupIds: groups,
		}

		if v, ok := d.GetOk("private_ip"); ok {
			ni.PrivateIps = []oapi.PrivateIpLight{oapi.PrivateIpLight{
				PrivateIp: v.(string),
			}}
		}

		networkInterfaces = append(networkInterfaces, ni)
	} else {
		// If we have manually specified network interfaces, build and attach those here.
		vL := nInterfaces.(*schema.Set).List()
		for _, v := range vL {
			ini := v.(map[string]interface{})
			ni := oapi.NicForVmCreation{
				NicId:              ini["nic_id"].(string),
				DeviceNumber:       int64(ini["nic_sort_number"].(int)),
				DeleteOnVmDeletion: ini["delete_on_vm_deletion"].(bool),
			}
			networkInterfaces = append(networkInterfaces, ni)
		}
	}

	return networkInterfaces
}

func readBlockDeviceOApiMappingsFromConfig(
	d *schema.ResourceData, conn *oapi.Client) ([]oapi.BlockDeviceMappingVmCreation, error) {
	blockDevices := make([]oapi.BlockDeviceMappingVmCreation, 0)

	if v, ok := d.GetOk("bsu"); ok {
		vL := v.(*schema.Set).List()
		for _, v := range vL {
			bd := v.(map[string]interface{})
			ebs := oapi.BsuToCreate{
				DeleteOnVmDeletion: bd["delete_on_vm_deletion"].(bool),
			}

			if v, ok := bd["snapshot_id"].(string); ok && v != "" {
				ebs.SnapshotId = v
			}
			if v, ok := bd["volume_size"].(int); ok && v != 0 {
				ebs.VolumeSize = int64(v)
			}
			if v, ok := bd["vm_type"].(string); ok && v != "" {
				ebs.VolumeType = v
			}
			if v, ok := bd["iops"].(int); ok && v > 0 {
				ebs.Iops = int64(v)
			}

			blockDevices = append(blockDevices, oapi.BlockDeviceMappingVmCreation{
				Bsu:               ebs,
				DeviceName:        bd["device_name"].(string),
				NoDevice:          bd["no_device"].(string),
				VirtualDeviceName: bd["virtual_device_name"].(string),
			})
		}
	}

	return blockDevices, nil
}

// InstanceStateOApiRefreshFunc ...
func InstanceStateOApiRefreshFunc(conn *oapi.Client, instanceID, failState string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		var resp *oapi.ReadVmsResponse
		var rs *oapi.POST_ReadVmsResponses
		var err error

		err = resource.Retry(30*time.Second, func() *resource.RetryError {
			rs, err = conn.POST_ReadVms(oapi.ReadVmsRequest{
				Filters: getVMsFilterByVMID(instanceID),
			})
			return resource.RetryableError(err)
		})

		if err != nil {
			fmt.Printf("Error on InstanceStateRefresh: %s", err)

			return nil, "", err
		}

		resp = rs.OK

		if resp == nil || len(resp.Vms) == 0 {
			return nil, "", nil
		}

		i := resp.Vms[0]
		state := i.State

		if state == failState {
			return i, state, fmt.Errorf("Failed to reach target state. Reason: %v",
				i.State)

		}

		return i, state, nil
	}
}

// // InstanceOApiPa ...
// func InstanceOApiPa(conn *oapi.Client, instanceID, failState string) resource.StateRefreshFunc {
// 	return func() (interface{}, string, error) {
// 		var resp *oapi.DescribeInstancesOutput
// 		var err error

// 		err = resource.Retry(30*time.Second, func() *resource.RetryError {
// 			resp, err = conn.VM.DescribeInstances(&oapi.DescribeInstancesInput{
// 				InstanceIds: []*string{aws.String(instanceID)},
// 			})

// 			return resource.RetryableError(err)
// 		})

// 		if err != nil {
// 			fmt.Printf("Error on InstanceStateRefresh: %s", err)

// 			return nil, "", err
// 		}

// 		if resp == nil || len(resp.Reservations) == 0 || len(resp.Reservations[0].Instances) == 0 {
// 			return nil, "", nil
// 		}

// 		i := resp.Reservations[0].Instances[0]
// 		state := *i.State.Name

// 		if state == failState {
// 			return i, state, fmt.Errorf("Failed to reach target state. Reason: %v",
// 				*i.StateReason)

// 		}

// 		return i, state, nil
// 	}
// }

// func updateVMAttr(conn *oapi.Client, instanceAttrOpts *oapi.UpdateVmAttributeRequest, attr string) error {

// 	var err error
// 	var stateConf *resource.StateChangeConf

// 	switch attr {
// 	case "instance_type":
// 		fallthrough
// 	case "user_data":
// 		fallthrough
// 	case "ebs_optimized":
// 		fallthrough
// 	case "delete_on_termination":
// 		stateConf, err = stopVM(instanceAttrOpts, conn, attr)
// 	}

// 	if err != nil {
// 		return err
// 	}

// 	if _, err := conn.POST_UpdateVmAttribute(*instanceAttrOpts); err != nil {
// 		return err
// 	}

// 	switch attr {
// 	case "instance_type":
// 		fallthrough
// 	case "user_data":
// 		fallthrough
// 	case "ebs_optimized":
// 		fallthrough
// 	case "delete_on_termination":
// 		err = startVM(instanceAttrOpts, stateConf, conn, attr)
// 	}

// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

func stopVM(vmID string, conn *oapi.Client, attr string) (*resource.StateChangeConf, error) {
	_, err := conn.POST_StopVms(oapi.StopVmsRequest{
		VmIds: []string{vmID},
	})

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"pending", "running", "shutting-down", "stopped", "stopping"},
		Target:     []string{"stopped"},
		Refresh:    InstanceStateOApiRefreshFunc(conn, vmID, ""),
		Timeout:    10 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return nil, fmt.Errorf(
			"Error waiting for instance (%s) to stop: %s", vmID, err)
	}

	return stateConf, nil
}

func startVM(vmID string, stateConf *resource.StateChangeConf, conn *oapi.Client, attr string) error {
	if _, err := conn.POST_StartVms(oapi.StartVmsRequest{
		VmIds: []string{vmID},
	}); err != nil {
		return err
	}

	stateConf = &resource.StateChangeConf{
		Pending:    []string{"pending", "stopped"},
		Target:     []string{"running"},
		Refresh:    InstanceStateOApiRefreshFunc(conn, vmID, ""),
		Timeout:    10 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("Error waiting for instance (%s) to become ready: %s", vmID, err)
	}

	return nil
}

func getVMsFilterByVMID(vmID string) oapi.FiltersVm {
	return oapi.FiltersVm{
		VmIds: []string{vmID},
	}
}
