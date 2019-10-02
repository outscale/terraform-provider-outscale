package outscale

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-outscale/osc/oapi"

	"github.com/spf13/cast"
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

		Schema: map[string]*schema.Schema{
			"block_device_mappings": {
				Type:     schema.TypeSet,
				Optional: true,
				Set: func(v interface{}) int {
					return hashcode.String(v.(map[string]interface{})["device_name"].(string))
				},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"bsu": {
							Type:     schema.TypeMap,
							Optional: true,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"delete_on_vm_deletion": {
										Type:     schema.TypeBool,
										Computed: true,
										Optional: true,
									},
									"link_date": {
										Type:     schema.TypeString,
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
									"iops": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"snapshot_id": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"volume_size": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"volume_type": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"device_name": {
							Type:     schema.TypeString,
							Computed: true,
							Optional: true,
						},
						"no_device": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"virtual_device_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"bsu_optimized": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"client_token": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"deletion_protection": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"image_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
				Computed: true,
			},
			"keypair_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"nics": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"delete_on_vm_deletion": {
							Type:     schema.TypeBool,
							Computed: true,
							Optional: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
							Optional: true,
						},
						"device_number": {
							Type:     schema.TypeInt,
							Computed: true,
							Optional: true,
						},
						"nic_id": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"private_ips": {
							Type:     schema.TypeSet,
							Optional: true,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"is_primary": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"link_public_ip": {
										Type:     schema.TypeSet,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"public_dns_name": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"public_ip": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"public_ip_account_id": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
									"private_dns_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"private_ip": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
								},
							},
						},
						"secondary_private_ip_count": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"account_id": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"is_source_dest_checked": {
							Type:     schema.TypeBool,
							Computed: true,
						},

						"subnet_id": {
							Type:     schema.TypeString,
							Computed: true,
							Optional: true,
						},
						"link_nic": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"delete_on_vm_deletion": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"device_number": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"link_nic_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"state": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"link_public_ip": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"public_dns_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"public_ip": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"public_ip_account_id": {
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
						"net_id": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"private_dns_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"security_group_ids": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"security_groups_names": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"security_groups": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"security_group_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"security_group_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"state": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"placement_subregion_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"placement_tenancy": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"private_ips": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"security_group_ids": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"security_group_names": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"subnet_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
				Computed: true,
			},

			"security_groups": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"security_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"security_group_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"architecture": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"hypervisor": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_source_dest_checked": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"launch_number": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"net_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"os_family": {
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
			"product_codes": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"public_dns_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"public_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"reservation_id": {
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
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"state_reason": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"user_data": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"vm_id": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"vm_initiated_shutdown_behavior": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"vm_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"request_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"admin_password": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"value": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func resourceOAPIVMCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OAPI

	instanceOpts, err := buildCreateVmsRequest(d, meta)
	if err != nil {
		return err
	}

	// Create the instance
	var runResp *oapi.CreateVmsResponse
	var resp *oapi.POST_CreateVmsResponses
	err = resource.Retry(30*time.Second, func() *resource.RetryError {
		var err error
		resp, err = conn.POST_CreateVms(*instanceOpts)

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

	if tags, ok := d.GetOk("tags"); ok {
		err := assignOapiTags(tags.([]interface{}), vm.VmId, conn)
		if err != nil {
			return err
		}
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
			"type": "ssh",
			"host": vm.PublicIp,
		})
	} else if vm.PrivateIp != "" {
		d.SetConnInfo(map[string]string{
			"type": "ssh",
			"host": vm.PrivateIp,
		})
	}

	//Check if source dest check is enabled.
	if v, ok := d.GetOk("is_source_dest_checked"); ok {

		opts := &oapi.UpdateVmRequest{
			VmId:                vm.VmId,
			IsSourceDestChecked: v.(bool),
		}
		log.Printf("[DEBGUG] is_source_dest_checked argument is not in CreateVms, we have to update the vm (%s)", vm.VmId)
		if err := oapiModifyInstanceAttr(conn, opts); err != nil {
			return err
		}
	}

	return resourceOAPIVMRead(d, meta)
}

func resourceOAPIVMRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OAPI

	var resp *oapi.ReadVmsResponse
	var rs *oapi.POST_ReadVmsResponses
	var err error

	err = resource.Retry(30*time.Second, func() *resource.RetryError {
		rs, err = conn.POST_ReadVms(*&oapi.ReadVmsRequest{
			Filters: oapi.FiltersVm{
				VmIds: []string{d.Id()},
			},
		})

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

	// Get the admin password from the server to save in the state
	adminPassword, err := getOAPIVMAdminPassword(instance.VmId, conn)
	if err != nil {
		return err
	}

	d.Set("request_id", resp.ResponseContext.RequestId)
	return resourceDataAttrSetter(d, func(set AttributeSetter) error {
		d.SetId(instance.VmId)

		set("architecture", instance.Architecture)
		// we need to get since index 1 ([1:]) all the slices because the server response always gets us another one item.
		if err := set("block_device_mappings", getOAPIVMBlockDeviceMapping(instance.BlockDeviceMappings)[1:]); err != nil {
			log.Printf("[DEBUG] BLOCKING DEVICE MAPPING ERR %+v", err)
			return err
		}
		set("bsu_optimized", instance.BsuOptimized)
		set("client_token", instance.ClientToken)
		set("deletion_protection", instance.DeletionProtection)
		set("hypervisor", instance.Hypervisor)
		set("image_id", instance.ImageId)
		set("is_source_dest_checked", instance.IsSourceDestChecked)
		set("keypair_name", instance.KeypairName)
		set("launch_number", instance.LaunchNumber)
		set("net_id", instance.NetId)
		if err := set("nics", getOAPIVMNetworkInterfaceSet(instance.Nics)); err != nil {
			log.Printf("[DEBUG] NICS ERR %+v", err)
			return err
		}
		set("os_family", instance.OsFamily)
		set("placement_subregion_name", instance.Placement.SubregionName)
		set("placement_tenancy", instance.Placement.Tenancy)
		set("private_dns_name", instance.PrivateDnsName)
		set("private_ip", instance.PrivateIp)
		set("product_codes", instance.ProductCodes)
		set("public_dns_name", instance.PublicDnsName)
		set("public_ip", instance.PublicIp)
		set("reservation_id", instance.ReservationId)
		set("root_device_name", instance.RootDeviceName)
		set("root_device_type", instance.RootDeviceType)
		if err := set("security_groups", getOAPIVMSecurityGroups(instance.SecurityGroups)); err != nil {
			log.Printf("[DEBUG] SECURITY GROUPS ERR %+v", err)
			return err
		}
		set("state", instance.State)
		set("state_reason", instance.StateReason)
		set("subnet_id", instance.SubnetId)
		set("user_data", instance.UserData)
		set("vm_id", instance.VmId)
		set("vm_initiated_shutdown_behavior", instance.VmInitiatedShutdownBehavior)
		set("admin_password", adminPassword)
		set("tags", getOapiTagSet(instance.Tags))

		return set("vm_type", instance.VmType)
	})
}

func getOAPIVMAdminPassword(VMID string, conn *oapi.Client) (string, error) {
	resp, err := conn.POST_ReadAdminPassword(oapi.ReadAdminPasswordRequest{VmId: VMID})
	if err != nil {
		return "", fmt.Errorf("Error reading the VM %s", err)
	}
	return resp.OK.AdminPassword, nil
}

func resourceOAPIVMUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OAPI

	d.Partial(true)

	id := d.Get("vm_id").(string)

	var stateConf *resource.StateChangeConf
	var err error
	if d.HasChange("vm_type") && !d.IsNewResource() ||
		d.HasChange("user_data") && !d.IsNewResource() ||
		d.HasChange("bsu_optimized") && !d.IsNewResource() {
		stateConf, err = oapiStopInstance(id, conn)
		if err != nil {
			return err
		}
	}

	if d.HasChange("vm_type") && !d.IsNewResource() {
		opts := &oapi.UpdateVmRequest{
			VmId:   id,
			VmType: d.Get("vm_type").(string),
		}
		if err := oapiModifyInstanceAttr(conn, opts); err != nil {
			return err
		}
	}

	if d.HasChange("user_data") && !d.IsNewResource() {
		opts := &oapi.UpdateVmRequest{
			VmId:     id,
			UserData: d.Get("user_data").(string),
		}
		if err := oapiModifyInstanceAttr(conn, opts); err != nil {
			return err
		}
	}

	if d.HasChange("bsu_optimized") && !d.IsNewResource() {
		opts := &oapi.UpdateVmRequest{
			VmId:         id,
			BsuOptimized: d.Get("bsu_optimized").(bool),
		}
		if err := oapiModifyInstanceAttr(conn, opts); err != nil {
			return err
		}
	}

	if d.HasChange("deletion_protection") && !d.IsNewResource() {
		deletionProtection := d.Get("deletion_protection").(bool)
		opts := &oapi.UpdateVmRequest{
			VmId:               id,
			DeletionProtection: &deletionProtection,
		}

		if err := oapiModifyInstanceAttr(conn, opts); err != nil {
			return err
		}
	}

	if d.HasChange("keypair_name") && !d.IsNewResource() {
		opts := &oapi.UpdateVmRequest{
			VmId:        id,
			KeypairName: d.Get("keypair_name").(string),
		}
		if err := oapiModifyInstanceAttr(conn, opts); err != nil {
			return err
		}
	}

	if d.HasChange("security_group_ids") && !d.IsNewResource() {
		opts := &oapi.UpdateVmRequest{
			VmId:             id,
			SecurityGroupIds: expandStringValueList(d.Get("security_group_ids").([]interface{})),
		}
		if err := oapiModifyInstanceAttr(conn, opts); err != nil {
			return err
		}
	}

	if d.HasChange("security_group_names") && !d.IsNewResource() {
		opts := &oapi.UpdateVmRequest{
			VmId:             id,
			SecurityGroupIds: expandStringValueList(d.Get("security_group_names").([]interface{})),
		}
		if err := oapiModifyInstanceAttr(conn, opts); err != nil {
			return err
		}
	}

	if d.HasChange("vm_initiated_shutdown_behavior") && !d.IsNewResource() {
		opts := &oapi.UpdateVmRequest{
			VmId:                        id,
			VmInitiatedShutdownBehavior: d.Get("vm_initiated_shutdown_behavior").(string),
		}
		if err := oapiModifyInstanceAttr(conn, opts); err != nil {
			return err
		}
	}

	if d.HasChange("is_source_dest_checked") && !d.IsNewResource() {
		opts := &oapi.UpdateVmRequest{
			VmId:                id,
			IsSourceDestChecked: d.Get("is_source_dest_checked").(bool),
		}
		if err := oapiModifyInstanceAttr(conn, opts); err != nil {
			return err
		}
	}

	if d.HasChange("block_device_mappings") && !d.IsNewResource() {
		maps := d.Get("block_device_mappings").(*schema.Set).List()
		mappings := []oapi.BlockDeviceMappingVmUpdate{}

		for _, m := range maps {
			f := m.(map[string]interface{})
			mapping := oapi.BlockDeviceMappingVmUpdate{
				DeviceName:        f["device_name"].(string),
				NoDevice:          f["no_device"].(string),
				VirtualDeviceName: f["virtual_device_name"].(string),
			}

			e := f["bsu"].(map[string]interface{})

			bsu := oapi.BsuToUpdateVm{
				DeleteOnVmDeletion: e["delete_on_vm_deletion"].(bool),
				VolumeId:           e["volume_id"].(string),
			}

			mapping.Bsu = bsu

			mappings = append(mappings, mapping)
		}

		opts := &oapi.UpdateVmRequest{
			VmId:                id,
			BlockDeviceMappings: mappings,
		}

		if err := oapiModifyInstanceAttr(conn, opts); err != nil {
			return err
		}
	}

	d.Partial(false)

	if err := oapiStartInstance(id, stateConf, conn); err != nil {
		return err
	}

	return resourceOAPIVMRead(d, meta)
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

func buildCreateVmsRequest(d *schema.ResourceData, meta interface{}) (*oapi.CreateVmsRequest, error) {
	request := &oapi.CreateVmsRequest{
		BlockDeviceMappings:         expandBlockDeviceOApiMappings(d),
		BsuOptimized:                d.Get("bsu_optimized").(bool),
		ClientToken:                 d.Get("client_token").(string),
		ImageId:                     d.Get("image_id").(string),
		KeypairName:                 d.Get("keypair_name").(string),
		MaxVmsCount:                 int64(1),
		MinVmsCount:                 int64(1),
		Nics:                        buildNetworkOApiInterfaceOpts(d),
		PrivateIps:                  expandStringValueList(d.Get("private_ips").([]interface{})),
		SecurityGroupIds:            expandStringValueList(d.Get("security_group_ids").([]interface{})),
		SecurityGroups:              expandStringValueList(d.Get("security_group_names").([]interface{})),
		SubnetId:                    d.Get("subnet_id").(string),
		UserData:                    d.Get("user_data").(string),
		VmInitiatedShutdownBehavior: d.Get("vm_initiated_shutdown_behavior").(string),
		VmType:                      d.Get("vm_type").(string),
	}

	deletionProtection := cast.ToBool(d.Get("deletion_protection")) == true
	request.DeletionProtection = &deletionProtection

	request.Placement = expandOAPIPlacement(d)
	return request, nil
}

func expandBlockDeviceOApiMappings(d *schema.ResourceData) []oapi.BlockDeviceMappingVmCreation {

	block := d.Get("block_device_mappings").(*schema.Set).List()
	blockDevices := make([]oapi.BlockDeviceMappingVmCreation, len(block))

	for i, v := range block {
		value := v.(map[string]interface{})
		bsu := value["bsu"].(map[string]interface{})

		deleteOnVMDeletion := (cast.ToBool(bsu["delete_on_vm_deletion"]) == true)
		blockDevices[i].Bsu.DeleteOnVmDeletion = &deleteOnVMDeletion

		if iops, ok := bsu["iops"]; ok {
			blockDevices[i].Bsu.Iops = cast.ToInt64(iops)
		}
		if snapshotID, ok := bsu["snapshot_id"]; ok {
			blockDevices[i].Bsu.SnapshotId = cast.ToString(snapshotID)
		}
		if volumeSize, ok := bsu["volume_size"]; ok {
			blockDevices[i].Bsu.VolumeSize = cast.ToInt64(volumeSize)
		}
		if volumeType, ok := bsu["volume_type"]; ok {
			blockDevices[i].Bsu.VolumeType = cast.ToString(volumeType)
		}
		if deviceName, ok := value["device_name"]; ok {
			blockDevices[i].DeviceName = cast.ToString(deviceName)
		}
		if noDevice, ok := value["no_device"]; ok {
			blockDevices[i].NoDevice = cast.ToString(noDevice)
		}
		if virtualDeviceName, ok := value["virtual_device_name"]; ok {
			blockDevices[i].VirtualDeviceName = cast.ToString(virtualDeviceName)
		}
	}
	return blockDevices
}

func buildNetworkOApiInterfaceOpts(d *schema.ResourceData) []oapi.NicForVmCreation {

	nics := d.Get("nics").([]interface{})
	networkInterfaces := []oapi.NicForVmCreation{}

	for _, v := range nics {
		nic := v.(map[string]interface{})

		ni := oapi.NicForVmCreation{
			DeleteOnVmDeletion: nic["delete_on_vm_deletion"].(bool),
			Description:        nic["description"].(string),
			DeviceNumber:       int64(nic["device_number"].(int)),
		}

		ni.PrivateIps = expandPrivatePublicIps(nic["private_ips"].(*schema.Set))
		ni.SubnetId = nic["subnet_id"].(string)
		ni.SecurityGroupIds = expandStringValueList(nic["security_group_ids"].([]interface{}))
		ni.SecondaryPrivateIpCount = int64(nic["secondary_private_ip_count"].(int))
		ni.NicId = nic["nic_id"].(string)

		if v, ok := d.GetOk("private_ip"); ok {
			ni.PrivateIps = []oapi.PrivateIpLight{oapi.PrivateIpLight{
				PrivateIp: v.(string),
			}}
		}
		networkInterfaces = append(networkInterfaces, ni)
	}

	return networkInterfaces
}

func expandPrivatePublicIps(p *schema.Set) []oapi.PrivateIpLight {
	privatePublicIPS := make([]oapi.PrivateIpLight, len(p.List()))

	for i, v := range p.List() {
		value := v.(map[string]interface{})
		privatePublicIPS[i].IsPrimary = value["is_primary"].(bool)
		privatePublicIPS[i].PrivateIp = value["private_ip"].(string)
	}
	return privatePublicIPS
}

func expandOAPIPlacement(d *schema.ResourceData) oapi.Placement {
	return oapi.Placement{
		SubregionName: d.Get("placement_subregion_name").(string),
		Tenancy:       d.Get("placement_tenancy").(string),
	}
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

// AttributeSetter you can use this function to set the attributes
type AttributeSetter func(key string, value interface{}) error

func resourceDataAttrSetter(d *schema.ResourceData, callback func(AttributeSetter) error) error {
	setterFunc := func(key string, value interface{}) error {
		return d.Set(key, value)
	}
	return callback(setterFunc)
}
