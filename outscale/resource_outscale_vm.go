package outscale

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/spf13/cast"

	"github.com/antihax/optional"
	oscgo "github.com/marinsalinas/osc-sdk-go"
	"github.com/terraform-providers/terraform-provider-outscale/utils"
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
				Type:     schema.TypeList,
				Optional: true,
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
										Optional: true,
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
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Set: func(v interface{}) int {
					return v.(map[string]interface{})["device_number"].(int)
				},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"delete_on_vm_deletion": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
							Optional: true,
						},
						"device_number": {
							Type:     schema.TypeInt,
							Required: true,
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
							Type:     schema.TypeList,
							MaxItems: 1,
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
			"block_device_mappings_created": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
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
									},
									"link_date": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"state": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"volume_id": {
										Type:     schema.TypeFloat,
										Computed: true,
									},
								},
							},
						},
						"device_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
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
			"performance": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
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
			"tags": tagsListOAPISchema(),
		},
	}
}

func resourceOAPIVMCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	vmOpts, err := buildCreateVmsRequest(d, meta)
	if err != nil {
		return err
	}

	// Create the vm
	var resp oscgo.CreateVmsResponse
	err = resource.Retry(30*time.Second, func() *resource.RetryError {
		var err error
		resp, _, err = conn.VmApi.CreateVms(context.Background(), &oscgo.CreateVmsOpts{
			CreateVmsRequest: optional.NewInterface(vmOpts),
		})

		if err != nil {
			if strings.Contains(fmt.Sprint(err), "Throttling") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("Error launching source VM: %s", utils.GetErrorResponse(err))
	}

	if !resp.HasVms() || len(resp.GetVms()) == 0 {
		return errors.New("Error launching source VM: no VMs returned in response")
	}

	vm := resp.GetVms()[0]

	d.SetId(vm.GetVmId())

	if tags, ok := d.GetOk("tags"); ok {
		err := assignTags(tags.([]interface{}), vm.GetVmId(), conn)
		if err != nil {
			return err
		}
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"pending", "ending/wait"},
		Target:     []string{"running"},
		Refresh:    vmStateRefreshFunc(conn, vm.GetVmId(), "terminated"),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf(
			"Error waiting for instance (%s) to become created: %s", d.Id(), err)
	}

	// Initialize the connection info
	if vm.HasPublicIp() {
		d.SetConnInfo(map[string]string{
			"type": "ssh",
			"host": vm.GetPublicIp(),
		})
	} else if vm.HasPrivateIp() {
		d.SetConnInfo(map[string]string{
			"type": "ssh",
			"host": vm.GetPrivateIp(),
		})
	}

	//Check if source dest check is enabled.
	if v, ok := d.GetOk("is_source_dest_checked"); ok {
		opts := oscgo.UpdateVmRequest{
			VmId: vm.GetVmId(),
		}

		opts.SetIsSourceDestChecked(v.(bool))

		log.Printf("[DEBUG] is_source_dest_checked argument is not in CreateVms, we have to update the vm (%s)", vm.GetVmId())
		if err := updateVmAttr(conn, opts); err != nil {
			return err
		}
	}

	return resourceOAPIVMRead(d, meta)
}

func resourceOAPIVMRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	var resp oscgo.ReadVmsResponse
	err := resource.Retry(30*time.Second, func() *resource.RetryError {
		r, _, err := conn.VmApi.ReadVms(context.Background(), &oscgo.ReadVmsOpts{
			ReadVmsRequest: optional.NewInterface(oscgo.ReadVmsRequest{
				Filters: &oscgo.FiltersVm{
					VmIds: &[]string{d.Id()},
				},
			}),
		})

		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded:") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}

		resp = r
		return nil
	})

	if err != nil {
		return fmt.Errorf("error reading the VM (%s): %s", d.Id(), err)
	}

	// If nothing was found, then return no state
	if !resp.HasVms() {
		d.SetId("")
		return nil
	}

	vm := resp.GetVms()[0]

	// Get the admin password from the server to save in the state
	adminPassword, err := getOAPIVMAdminPassword(vm.GetVmId(), conn)
	if err != nil {
		return err
	}

	return resourceDataAttrSetter(d, func(set AttributeSetter) error {
		if err := d.Set("request_id", resp.ResponseContext.RequestId); err != nil {
			return err
		}
		if err := d.Set("admin_password", adminPassword); err != nil {
			return err
		}
		d.SetId(vm.GetVmId())
		return oapiVMDescriptionAttributes(set, &vm)
	})
}

func getOAPIVMAdminPassword(VMID string, conn *oscgo.APIClient) (string, error) {
	resp, _, err := conn.VmApi.ReadAdminPassword(context.Background(),
		&oscgo.ReadAdminPasswordOpts{
			ReadAdminPasswordRequest: optional.NewInterface(oscgo.ReadAdminPasswordRequest{VmId: VMID}),
		},
	)

	if err != nil {
		return "", fmt.Errorf("error reading the VM's password %s", err)
	}
	return resp.GetAdminPassword(), nil
}

func resourceOAPIVMUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	d.Partial(true)

	id := d.Get("vm_id").(string)

	if d.HasChange("vm_type") && !d.IsNewResource() ||
		d.HasChange("user_data") && !d.IsNewResource() ||
		d.HasChange("bsu_optimized") && !d.IsNewResource() {
		if err := stopVM(id, conn); err != nil {
			return err
		}
	}

	if d.HasChange("vm_type") && !d.IsNewResource() {
		opts := oscgo.UpdateVmRequest{VmId: id}
		opts.SetVmType(d.Get("vm_type").(string))

		if err := updateVmAttr(conn, opts); err != nil {
			return err
		}
	}

	if d.HasChange("user_data") && !d.IsNewResource() {
		opts := oscgo.UpdateVmRequest{VmId: id}
		opts.SetUserData(d.Get("user_data").(string))

		if err := updateVmAttr(conn, opts); err != nil {
			return err
		}
	}

	if d.HasChange("bsu_optimized") && !d.IsNewResource() {
		opts := oscgo.UpdateVmRequest{VmId: id}
		opts.SetBsuOptimized(d.Get("bsu_optimized").(bool))

		if err := updateVmAttr(conn, opts); err != nil {
			return err
		}
	}

	if d.HasChange("deletion_protection") && !d.IsNewResource() {
		opts := oscgo.UpdateVmRequest{VmId: id}
		opts.SetDeletionProtection(d.Get("deletion_protection").(bool))

		if err := updateVmAttr(conn, opts); err != nil {
			return err
		}
	}

	if d.HasChange("keypair_name") && !d.IsNewResource() {
		opts := oscgo.UpdateVmRequest{VmId: id}
		opts.SetKeypairName(d.Get("keypair_name").(string))
		if err := updateVmAttr(conn, opts); err != nil {
			return err
		}
	}

	if d.HasChange("security_group_ids") && !d.IsNewResource() {
		opts := oscgo.UpdateVmRequest{VmId: id}

		opts.SetSecurityGroupIds(expandStringValueList(d.Get("security_group_ids").([]interface{})))
		if err := updateVmAttr(conn, opts); err != nil {
			return err
		}
	}

	if d.HasChange("security_group_names") && !d.IsNewResource() {
		opts := oscgo.UpdateVmRequest{VmId: id}
		opts.SetSecurityGroupIds(expandStringValueList(d.Get("security_group_names").([]interface{})))
		if err := updateVmAttr(conn, opts); err != nil {
			return err
		}
	}

	if d.HasChange("vm_initiated_shutdown_behavior") && !d.IsNewResource() {
		opts := oscgo.UpdateVmRequest{VmId: id}
		opts.SetVmInitiatedShutdownBehavior(d.Get("vm_initiated_shutdown_behavior").(string))

		if err := updateVmAttr(conn, opts); err != nil {
			return err
		}
	}

	if d.HasChange("is_source_dest_checked") && !d.IsNewResource() {
		opts := oscgo.UpdateVmRequest{VmId: id}
		opts.SetIsSourceDestChecked(d.Get("is_source_dest_checked").(bool))

		if err := updateVmAttr(conn, opts); err != nil {
			return err
		}
	}

	if d.HasChange("performance") && !d.IsNewResource() {
		opts := oscgo.UpdateVmRequest{VmId: id}
		opts.SetPerformance(d.Get("performance").(string))

		if err := updateVmAttr(conn, opts); err != nil {
			return err
		}
	}

	if d.HasChange("block_device_mappings") && !d.IsNewResource() {
		maps := d.Get("block_device_mappings").(*schema.Set).List()
		mappings := []oscgo.BlockDeviceMappingVmUpdate{}

		for _, m := range maps {
			f := m.(map[string]interface{})
			mapping := oscgo.BlockDeviceMappingVmUpdate{}
			mapping.SetDeviceName(f["device_name"].(string))
			mapping.SetNoDevice(f["no_device"].(string))
			mapping.SetVirtualDeviceName(f["virtual_device_name"].(string))

			e := f["bsu"].(map[string]interface{})
			bsu := oscgo.BsuToUpdateVm{}

			bsu.SetDeleteOnVmDeletion(e["delete_on_vm_deletion"].(bool))
			bsu.SetVolumeId(e["volume_id"].(string))

			mapping.SetBsu(bsu)

			mappings = append(mappings, mapping)
		}

		opts := oscgo.UpdateVmRequest{VmId: id}

		opts.SetBlockDeviceMappings(mappings)

		if err := updateVmAttr(conn, opts); err != nil {
			return err
		}
	}

	if err := setOSCAPITags(conn, d); err != nil {
		return err
	}

	d.SetPartial("tags")

	d.Partial(false)

	if err := startVM(id, conn); err != nil {
		return err
	}

	return resourceOAPIVMRead(d, meta)
}

func resourceOAPIVMDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	id := d.Id()

	log.Printf("[INFO] Terminating VM: %s", id)

	var err error
	err = resource.Retry(30*time.Second, func() *resource.RetryError {
		_, _, err = conn.VmApi.DeleteVms(context.Background(), &oscgo.DeleteVmsOpts{
			DeleteVmsRequest: optional.NewInterface(oscgo.DeleteVmsRequest{
				VmIds: []string{id},
			}),
		})

		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded") {
				fmt.Printf("[INFO] Request limit exceeded")
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("Error deleting the VM")
	}

	log.Printf("[DEBUG] Waiting for VM (%s) to become terminated", id)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"pending", "running", "shutting-down", "stopped", "stopping"},
		Target:     []string{"terminated"},
		Refresh:    vmStateRefreshFunc(conn, id, ""),
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

func buildCreateVmsRequest(d *schema.ResourceData, meta interface{}) (oscgo.CreateVmsRequest, error) {
	request := oscgo.CreateVmsRequest{
		DeletionProtection: oscgo.PtrBool(d.Get("deletion_protection").(bool)),
		BsuOptimized:       oscgo.PtrBool(d.Get("bsu_optimized").(bool)),
		MaxVmsCount:        oscgo.PtrInt64(1),
		MinVmsCount:        oscgo.PtrInt64(1),
		ImageId:            d.Get("image_id").(string),
		Placement:          expandPlacement(d),
	}

	if nics := buildNetworkOApiInterfaceOpts(d); len(nics) > 0 {
		request.SetNics(nics)
	}

	if blockDevices := expandBlockDeviceOApiMappings(d); len(blockDevices) > 0 {
		request.SetBlockDeviceMappings(blockDevices)
	}

	if privateIPs := expandStringValueList(d.Get("private_ips").([]interface{})); len(privateIPs) > 0 {
		request.SetPrivateIps(privateIPs)
	}

	if sgIDs := expandStringValueList(d.Get("security_group_ids").([]interface{})); len(sgIDs) > 0 {
		request.SetSecurityGroupIds(sgIDs)
	}

	if sgNames := expandStringValueList(d.Get("security_group_names").([]interface{})); len(sgNames) > 0 {
		request.SetSecurityGroups(sgNames)
	}

	if v := d.Get("subnet_id").(string); v != "" {
		request.SetSubnetId(v)
	}

	if v := d.Get("user_data").(string); v != "" {
		request.SetUserData(v)
	}

	if v := d.Get("vm_type").(string); v != "" {
		request.SetVmType(v)
	}

	if v := d.Get("client_token").(string); v != "" {
		request.SetClientToken(v)
	}

	if v := d.Get("keypair_name").(string); v != "" {
		request.SetKeypairName(v)
	}

	if v, ok := d.GetOk("vm_initiated_shutdown_behavior"); ok && v != "" {
		request.SetVmInitiatedShutdownBehavior(v.(string))
	}

	if v := d.Get("performance").(string); v != "" {
		request.SetPerformance(v)
	}

	return request, nil
}

func expandBlockDeviceOApiMappings(d *schema.ResourceData) []oscgo.BlockDeviceMappingVmCreation {
	var blockDevices []oscgo.BlockDeviceMappingVmCreation

	block := d.Get("block_device_mappings").([]interface{})

	for _, v := range block {
		blockDevice := oscgo.BlockDeviceMappingVmCreation{}

		value := v.(map[string]interface{})
		bsu := value["bsu"].(map[string]interface{})

		blockDevice.SetBsu(expandBlockDeviceBSU(bsu))

		if deviceName, ok := value["device_name"]; ok {
			blockDevice.SetDeviceName(cast.ToString(deviceName))
		}
		if noDevice, ok := value["no_device"]; ok {
			blockDevice.SetNoDevice(cast.ToString(noDevice))
		}
		if virtualDeviceName, ok := value["virtual_device_name"]; ok {
			blockDevice.SetVirtualDeviceName(cast.ToString(virtualDeviceName))
		}

		blockDevices = append(blockDevices, blockDevice)
	}
	return blockDevices
}

func expandBlockDeviceBSU(bsu map[string]interface{}) oscgo.BsuToCreate {
	bsuToCreate := oscgo.BsuToCreate{}

	if deleteOnVMDeletion, ok := bsu["delete_on_vm_deletion"]; ok {
		bsuToCreate.SetDeleteOnVmDeletion(cast.ToBool(deleteOnVMDeletion))
	}

	if iops, ok := bsu["iops"]; ok {
		bsuToCreate.SetIops(cast.ToInt64(iops))
	}
	if snapshotID, ok := bsu["snapshot_id"]; ok {
		bsuToCreate.SetSnapshotId(cast.ToString(snapshotID))
	}
	if volumeSize, ok := bsu["volume_size"]; ok {
		bsuToCreate.SetVolumeSize(cast.ToInt64(volumeSize))
	}
	if volumeType, ok := bsu["volume_type"]; ok {
		bsuToCreate.SetVolumeType(cast.ToString(volumeType))
	}

	return bsuToCreate
}

func buildNetworkOApiInterfaceOpts(d *schema.ResourceData) []oscgo.NicForVmCreation {

	nics := d.Get("nics").(*schema.Set).List()
	networkInterfaces := []oscgo.NicForVmCreation{}

	for i, v := range nics {
		nic := v.(map[string]interface{})

		ni := oscgo.NicForVmCreation{
			DeviceNumber: oscgo.PtrInt64(int64(nic["device_number"].(int))),
		}

		if v := nic["nic_id"].(string); v != "" {
			ni.SetNicId(v)
		}

		if v := nic["secondary_private_ip_count"].(int); v > 0 {
			ni.SetSecondaryPrivateIpCount(int64(v))
		}

		if delete, deleteOK := d.GetOk(fmt.Sprintf("nics.%d.delete_on_vm_deletion", i)); deleteOK {
			log.Printf("[DEBUG] delete=%+v, deleteOK=%+v", delete, deleteOK)
			ni.SetDeleteOnVmDeletion(delete.(bool))
		}

		ni.SetDescription(nic["description"].(string))

		ni.SetPrivateIps(expandPrivatePublicIps(nic["private_ips"].(*schema.Set)))
		ni.SetSubnetId(nic["subnet_id"].(string))

		if sg := expandStringValueList(nic["security_group_ids"].([]interface{})); len(sg) > 0 {
			ni.SetSecurityGroupIds(sg)
		}

		if v, ok := d.GetOk("private_ip"); ok {
			ni.SetPrivateIps([]oscgo.PrivateIpLight{oscgo.PrivateIpLight{
				PrivateIp: aws.String(v.(string)),
			}})
		}
		networkInterfaces = append(networkInterfaces, ni)
	}

	return networkInterfaces
}

func expandPrivatePublicIps(p *schema.Set) []oscgo.PrivateIpLight {
	privatePublicIPS := make([]oscgo.PrivateIpLight, len(p.List()))

	for i, v := range p.List() {
		value := v.(map[string]interface{})
		privatePublicIPS[i].SetIsPrimary(value["is_primary"].(bool))
		privatePublicIPS[i].SetPrivateIp(value["private_ip"].(string))
	}
	return privatePublicIPS
}

func expandPlacement(d *schema.ResourceData) *oscgo.Placement {
	var placement *oscgo.Placement

	subregionName, sOK := d.GetOk("placement_subregion_name")
	tenancy, tOK := d.GetOk("placement_tenancy")

	if sOK || tOK {
		placement = &oscgo.Placement{
			SubregionName: oscgo.PtrString(subregionName.(string)),
		}

		placement.Tenancy = oscgo.PtrString(tenancy.(string))
	}
	return placement
}

func vmStateRefreshFunc(conn *oscgo.APIClient, instanceID, failState string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, _, err := conn.VmApi.ReadVms(context.Background(), &oscgo.ReadVmsOpts{
			ReadVmsRequest: optional.NewInterface(oscgo.ReadVmsRequest{
				Filters: &oscgo.FiltersVm{
					VmIds: &[]string{instanceID},
				},
			}),
		})

		if err != nil {
			log.Printf("[ERROR] error on InstanceStateRefresh: %s", err)
			return nil, "", err
		}

		if !resp.HasVms() {
			return nil, "", nil
		}

		vm := resp.GetVms()[0]
		state := vm.GetState()

		if state == failState {
			return vm, state, fmt.Errorf("Failed to reach target state. Reason: %v", vm.State)

		}

		return vm, state, nil
	}
}

func stopVM(vmID string, conn *oscgo.APIClient) error {
	_, _, err := conn.VmApi.StopVms(context.Background(), &oscgo.StopVmsOpts{
		StopVmsRequest: optional.NewInterface(oscgo.StopVmsRequest{
			VmIds: []string{vmID},
		}),
	})

	if err != nil {
		return fmt.Errorf("error stopping vms %s", err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"pending", "running", "shutting-down", "stopped", "stopping"},
		Target:     []string{"stopped"},
		Refresh:    vmStateRefreshFunc(conn, vmID, ""),
		Timeout:    10 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Error waiting for instance (%s) to stop: %s", vmID, err)
	}

	return nil
}

func startVM(vmID string, conn *oscgo.APIClient) error {
	_, _, err := conn.VmApi.StartVms(context.Background(), &oscgo.StartVmsOpts{
		StartVmsRequest: optional.NewInterface(oscgo.StartVmsRequest{
			VmIds: []string{vmID},
		}),
	})

	if err != nil {
		return fmt.Errorf("error starting vm %s", err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"pending", "stopped"},
		Target:     []string{"running"},
		Refresh:    vmStateRefreshFunc(conn, vmID, ""),
		Timeout:    10 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("Error waiting for instance (%s) to become ready: %s", vmID, err)
	}

	return nil
}

func updateVmAttr(conn *oscgo.APIClient, instanceAttrOpts oscgo.UpdateVmRequest) error {
	if _, _, err := conn.VmApi.UpdateVm(context.Background(), &oscgo.UpdateVmOpts{
		UpdateVmRequest: optional.NewInterface(instanceAttrOpts),
	}); err != nil {
		return err
	}
	return nil
}

// AttributeSetter you can use this function to set the attributes
type AttributeSetter func(key string, value interface{}) error

func resourceDataAttrSetter(d *schema.ResourceData, callback func(AttributeSetter) error) error {
	setterFunc := func(key string, value interface{}) error {
		return d.Set(key, value)
	}
	return callback(setterFunc)
}
