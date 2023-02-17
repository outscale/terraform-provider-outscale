package outscale

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/spf13/cast"
	"github.com/terraform-providers/terraform-provider-outscale/utils"
)

func resourceOutscaleOApiVM() *schema.Resource {
	return &schema.Resource{
		Create: resourceOAPIVMCreate,
		Read:   resourceOAPIVMRead,
		Update: resourceOAPIVMUpdate,
		Delete: resourceOAPIVMDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"block_device_mappings": {
				Type:     schema.TypeList,
				Optional: true,
				//ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"bsu": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"delete_on_vm_deletion": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"iops": {
										Type:     schema.TypeInt,
										Optional: true,
										ForceNew: true,
										ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
											iopsVal := val.(int)
											if iopsVal < utils.MinIops || iopsVal > utils.MaxIops {
												errs = append(errs, fmt.Errorf("%q must be between %d and %d inclusive, got: %d", key, utils.MinIops, utils.MaxIops, iopsVal))
											}
											return
										},
									},
									"snapshot_id": {
										Type:     schema.TypeString,
										Optional: true,
										ForceNew: true,
									},
									"volume_size": {
										Type:     schema.TypeInt,
										Optional: true,
										ForceNew: true,
										ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
											vSize := val.(int)
											if vSize < 1 || vSize > utils.MaxSize {
												errs = append(errs, fmt.Errorf("%q must be between 1 and %d gibibytes inclusive, got: %d", key, utils.MaxSize, vSize))
											}
											return
										},
									},
									"volume_type": {
										Type:     schema.TypeString,
										Optional: true,
										ForceNew: true,
									},
								},
							},
						},
						"device_name": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"no_device": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"virtual_device_name": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
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
			"creation_date": {
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
				Required: true,
				ForceNew: true,
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
				ForceNew: true,
				Set: func(v interface{}) int {
					return v.(map[string]interface{})["device_number"].(int)
				},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"delete_on_vm_deletion": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
							Optional: true,
							ForceNew: true,
						},
						"device_number": {
							Type:     schema.TypeInt,
							Required: true,
							ForceNew: true,
						},
						"nic_id": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},
						"private_ips": {
							Type:     schema.TypeSet,
							Optional: true,
							Computed: true,
							ForceNew: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"is_primary": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
										ForceNew: true,
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
										ForceNew: true,
									},
								},
							},
						},
						"secondary_private_ip_count": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
							ForceNew: true,
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
							ForceNew: true,
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
							ForceNew: true,
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
			"security_group_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"security_group_names": {
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
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"bsu": {
							Type:     schema.TypeMap,
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
							Computed: true,
						},
						"boot_disk_tags": {
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
			"nested_virtualization": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
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
			"product_codes": {
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
				Optional: true,
				Default:  "running",
			},
			"state_reason": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"user_data": {
				Type:     schema.TypeString,
				Optional: true,
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
			"get_admin_password": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"boot_disk_tags": tagsListOAPISchema(),
			"tags":           tagsListOAPISchema(),
		},
	}
}

func resourceOAPIVMCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	vmOpts, err := buildCreateVmsRequest(d, meta)
	if err != nil {
		return err
	}

	vState := d.Get("state").(string)
	if vState != "stopped" && vState != "running" {
		return fmt.Errorf("Error: state should be `stopped or running`")
	}
	vmStateTarget := []string{"running"}
	if vState == "stopped" {
		vmStateTarget[0] = "stopped"
		vmOpts.BootOnCreation = oscgo.PtrBool(false)
	}

	// Create the vm
	var resp oscgo.CreateVmsResponse
	err = resource.Retry(120*time.Second, func() *resource.RetryError {
		rp, httpResp, err := conn.VmApi.CreateVms(context.Background()).CreateVmsRequest(vmOpts).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		resp = rp
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

	if get_psswd := d.Get("get_admin_password").(bool); get_psswd {
		psswd_err := resource.Retry(2500*time.Second, func() *resource.RetryError {
			psswd, err := getOAPIVMAdminPassword(vm.GetVmId(), conn)
			if err != nil || len(psswd) < 1 {
				return resource.RetryableError(errors.New("timeout awaiting windows password"))
			}
			if err != nil {
				return resource.NonRetryableError(err)
			}
			return nil
		})
		if psswd_err != nil {
			return psswd_err
		}
	}

	if tags, ok := d.GetOk("tags"); ok {
		err := assignTags(tags.(*schema.Set), vm.GetVmId(), conn)
		if err != nil {
			return err
		}
	}
	if tags, ok := d.GetOk("boot_disk_tags"); ok {
		err := assignTags(tags.(*schema.Set), utils.GetBootDiskId(vm), conn)
		if err != nil {
			return err
		}
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"pending", "ending/wait"},
		Target:     vmStateTarget,
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

	if v, exist := d.GetOkExists("is_source_dest_checked"); exist {
		opts := oscgo.UpdateVmRequest{
			VmId: vm.GetVmId(),
		}
		opts.SetIsSourceDestChecked(v.(bool))
		if err := updateVmAttr(conn, opts); err != nil {
			return err
		}
	}
	return resourceOAPIVMRead(d, meta)
}

func resourceOAPIVMRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	var resp oscgo.ReadVmsResponse
	err := resource.Retry(60*time.Second, func() *resource.RetryError {
		rp, httpResp, err := conn.VmApi.ReadVms(context.Background()).ReadVmsRequest(oscgo.ReadVmsRequest{
			Filters: &oscgo.FiltersVm{
				VmIds: &[]string{d.Id()},
			},
		}).Execute()

		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		resp = rp
		return nil
	})

	if err != nil {
		return fmt.Errorf("error reading the VM (%s): %s", d.Id(), err)
	}
	if utils.IsResponseEmpty(len(resp.GetVms()), "Snapshot", d.Id()) {
		d.SetId("")
		return nil
	}

	vm := resp.GetVms()[0]
	if vm.GetState() == "terminated" {
		utils.LogManuallyDeleted("Vm", d.Id())
		d.SetId("")
		return nil
	}
	adminPassword, err := getOAPIVMAdminPassword(vm.GetVmId(), conn)
	if err != nil {
		return err
	}
	bsu := d.Get("bsu_optimized")
	if err := resourceDataAttrSetter(d, func(set AttributeSetter) error {
		if err := d.Set("admin_password", adminPassword); err != nil {
			return err
		}
		d.SetId(vm.GetVmId())

		booTags, errTags := utils.GetBootDiskTags(utils.GetBootDiskId(vm), conn)
		if errTags != nil {
			return errTags
		}

		if err := d.Set("block_device_mappings_created", getOscAPIVMBlockDeviceMapping(
			booTags, vm.GetBlockDeviceMappings())); err != nil {
			return err
		}

		return oapiVMDescriptionAttributes(set, &vm)
	}); err != nil {
		return err
	}
	return d.Set("bsu_optimized", bsu)
}

func getOAPIVMAdminPassword(VMID string, conn *oscgo.APIClient) (string, error) {
	var resp oscgo.ReadAdminPasswordResponse
	err := resource.Retry(60*time.Second, func() *resource.RetryError {
		rp, httpResp, err := conn.VmApi.ReadAdminPassword(context.Background()).ReadAdminPasswordRequest(oscgo.ReadAdminPasswordRequest{VmId: VMID}).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		resp = rp
		return nil
	})

	if err != nil {
		return "", fmt.Errorf("error reading the VM's password %s", err)
	}
	return resp.GetAdminPassword(), nil
}

func resourceOAPIVMUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI
	id := d.Get("vm_id").(string)

	nothingToDo := true
	onlyTags := d.HasChange("tags")
	o, n := d.GetChange("")
	os := o.(map[string]interface{})
	ns := n.(map[string]interface{})

	for k := range os {
		if d.HasChange(k) && k != "get_admin_password" {
			nothingToDo = false
		}
		if d.HasChange(k) && k != "tags" {
			onlyTags = false
		}
	}

	for k := range ns {
		if d.HasChange(k) && k != "get_admin_password" {
			nothingToDo = false
		}
		if d.HasChange(k) && k != "tags" {
			onlyTags = false
		}
	}

	if nothingToDo == true {
		return nil
	}

	if !d.IsNewResource() &&
		(d.HasChange("vm_type") || d.HasChange("user_data") ||
			d.HasChange("performance") || d.HasChange("nested_virtualization")) {
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

	if d.HasChange("performance") && !d.IsNewResource() {
		opts := oscgo.UpdateVmRequest{VmId: id}
		opts.SetPerformance(d.Get("performance").(string))

		if err := updateVmAttr(conn, opts); err != nil {
			return err
		}
	}

	if d.HasChange("nested_virtualization") && !d.IsNewResource() {
		opts := oscgo.UpdateVmRequest{VmId: id}
		opts.SetNestedVirtualization(d.Get("nested_virtualization").(bool))

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

		opts.SetSecurityGroupIds(utils.InterfaceSliceToStringSlice(d.Get("security_group_ids").([]interface{})))
		if err := updateVmAttr(conn, opts); err != nil {
			return err
		}
	}

	if d.HasChange("security_group_names") && !d.IsNewResource() {
		opts := oscgo.UpdateVmRequest{VmId: id}
		opts.SetSecurityGroupIds(utils.InterfaceSliceToStringSlice(d.Get("security_group_names").([]interface{})))
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

	if d.HasChange("block_device_mappings") && !d.IsNewResource() {
		maps := d.Get("block_device_mappings").([]interface{})
		mappings := []oscgo.BlockDeviceMappingVmUpdate{}

		for _, m := range maps {
			f := m.(map[string]interface{})
			mapping := oscgo.BlockDeviceMappingVmUpdate{}

			if v, ok := f["device_name"]; ok && v.(string) != "" {
				mapping.SetDeviceName(v.(string))
			}

			if v, ok := f["no_device"]; ok && v.(string) != "" {
				mapping.SetNoDevice(v.(string))
			}

			if v, ok := f["virtual_device_name"]; ok && v.(string) != "" {
				mapping.SetVirtualDeviceName(v.(string))
			}

			if bsuList, ok := f["bsu"].([]interface{}); ok && len(bsuList) > 0 {
				bsu := oscgo.BsuToUpdateVm{}

				if e, ok1 := bsuList[0].(map[string]interface{}); ok1 {
					bsu.SetDeleteOnVmDeletion(cast.ToBool(e["delete_on_vm_deletion"]))

					if v, ok := e["volume_id"]; ok {
						bsu.SetVolumeId(v.(string))
					}
					mapping.SetBsu(bsu)
				}
			}

			mappings = append(mappings, mapping)
		}

		opts := oscgo.UpdateVmRequest{VmId: id}

		opts.SetBlockDeviceMappings(mappings)

		if err := updateVmAttr(conn, opts); err != nil {
			return utils.GetErrorResponse(err)
		}
	}

	if err := setOSCAPITags(conn, d, "tags"); err != nil {
		return err
	}

	if !d.IsNewResource() && d.HasChange("boot_disk_tags") {
		onlyTags = true
		if err := setOSCAPITags(conn, d, "boot_disk_tags"); err != nil {
			return err
		}
	}

	if onlyTags {
		goto out
	}

	if d.HasChange("state") && !d.IsNewResource() {
		upState := d.Get("state").(string)
		if upState != "stopped" && upState != "running" {
			return fmt.Errorf("Error: state should be `stopped or running`")
		}
		if upState == "stopped" {
			if err := stopVM(id, conn); err != nil {
				return err
			}
		} else {
			if err := startVM(id, conn); err != nil {
				return err
			}
		}
	} else {
		if err := startVM(id, conn); err != nil {
			return err
		}
	}

out:
	return resourceOAPIVMRead(d, meta)
}

func resourceOAPIVMDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	id := d.Id()

	log.Printf("[INFO] Terminating VM: %s", id)

	var err error
	err = resource.Retry(30*time.Second, func() *resource.RetryError {
		_, httpResp, err := conn.VmApi.DeleteVms(context.Background()).DeleteVmsRequest(oscgo.DeleteVmsRequest{
			VmIds: []string{id},
		}).Execute()

		if err != nil {
			return utils.CheckThrottling(httpResp, err)
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
		BootOnCreation:     oscgo.PtrBool(true),
		MaxVmsCount:        oscgo.PtrInt32(1),
		MinVmsCount:        oscgo.PtrInt32(1),
		ImageId:            d.Get("image_id").(string),
	}

	placement, err := expandPlacement(d)
	if err != nil {
		return request, err
	} else if placement != nil {
		request.SetPlacement(*placement)
	}

	subNet := d.Get("subnet_id").(string)
	if subNet != "" {
		request.SetSubnetId(subNet)
	}
	blockDevices, err := expandBlockDeviceOApiMappings(d)
	if err != nil {
		return request, err
	}
	if len(blockDevices) > 0 {
		request.SetBlockDeviceMappings(blockDevices)
	}

	if nics := buildNetworkOApiInterfaceOpts(d); len(nics) > 0 {
		if subNet != "" || placement != nil {
			return request, errors.New("If you specify nics parameter, you must not specify subnet_id and placement parameters.")
		}
		request.SetNics(nics)
	}

	if privateIPs := utils.InterfaceSliceToStringSlice(d.Get("private_ips").([]interface{})); len(privateIPs) > 0 {
		request.SetPrivateIps(privateIPs)
	}

	if sgIDs := utils.InterfaceSliceToStringSlice(d.Get("security_group_ids").([]interface{})); len(sgIDs) > 0 {
		request.SetSecurityGroupIds(sgIDs)
	}

	if sgNames := utils.InterfaceSliceToStringSlice(d.Get("security_group_names").([]interface{})); len(sgNames) > 0 {
		request.SetSecurityGroups(sgNames)
	}

	nestedVirtualization := d.Get("nested_virtualization").(bool)
	if tenacy := d.Get("placement_tenancy").(string); nestedVirtualization && tenacy != "dedicated" {
		return request, errors.New("The field nested_virtualization can be true, only if placement_tenancy is \"dedicated\".")
	}
	request.SetNestedVirtualization(nestedVirtualization)

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

func expandBlockDeviceOApiMappings(d *schema.ResourceData) ([]oscgo.BlockDeviceMappingVmCreation, error) {
	var blockDevices []oscgo.BlockDeviceMappingVmCreation
	block := d.Get("block_device_mappings").([]interface{})

	for _, v := range block {
		blockDevice := oscgo.BlockDeviceMappingVmCreation{}
		value := v.(map[string]interface{})

		if bsu, ok := value["bsu"].([]interface{}); ok && len(bsu) > 0 {
			expandBSU, err := expandBlockDeviceBSU(bsu[0].(map[string]interface{}))
			if err != nil {
				return nil, err
			}
			blockDevice.SetBsu(expandBSU)
		}
		if deviceName, ok := value["device_name"]; ok && deviceName != "" {
			blockDevice.SetDeviceName(cast.ToString(deviceName))
		}
		if noDevice, ok := value["no_device"]; ok && noDevice != "" {
			blockDevice.SetNoDevice(cast.ToString(noDevice))
		}
		if virtualDeviceName, ok := value["virtual_device_name"]; ok && virtualDeviceName != "" {
			blockDevice.SetVirtualDeviceName(cast.ToString(virtualDeviceName))
		}
		blockDevices = append(blockDevices, blockDevice)
	}
	return blockDevices, nil
}

func expandBlockDeviceBSU(bsu map[string]interface{}) (oscgo.BsuToCreate, error) {
	bsuToCreate := oscgo.BsuToCreate{}
	snapshotID := bsu["snapshot_id"].(string)
	volumeType := bsu["volume_type"].(string)
	volumeSize := int32(bsu["volume_size"].(int))

	if snapshotID == "" && volumeSize == 0 {
		return bsuToCreate, fmt.Errorf("Error: 'volume_size' parameter is required if the volume is not created from a snapshot (SnapshotId unspecified)")
	}
	if iops, _ := bsu["iops"]; iops.(int) > 0 {
		if volumeType != "io1" {
			return bsuToCreate, fmt.Errorf("Error: %s", utils.VolumeIOPSError)
		}
		bsuToCreate.SetIops(int32(iops.(int)))
	} else {
		delete(bsu, "iops")
	}
	if snapshotID != "" {
		bsuToCreate.SetSnapshotId(snapshotID)
	}
	if volumeSize > 0 {
		bsuToCreate.SetVolumeSize(volumeSize)
	}
	if volumeType != "" {
		bsuToCreate.SetVolumeType(volumeType)
	}
	if deleteOnVMDeletion, ok := bsu["delete_on_vm_deletion"]; ok && deleteOnVMDeletion != "" {
		bsuToCreate.SetDeleteOnVmDeletion(cast.ToBool(deleteOnVMDeletion))
	}
	return bsuToCreate, nil
}

func buildNetworkOApiInterfaceOpts(d *schema.ResourceData) []oscgo.NicForVmCreation {
	nics := d.Get("nics").(*schema.Set).List()
	networkInterfaces := []oscgo.NicForVmCreation{}

	for i, v := range nics {
		nic := v.(map[string]interface{})
		ni := oscgo.NicForVmCreation{
			DeviceNumber: oscgo.PtrInt32(int32(nic["device_number"].(int))),
		}

		if v := nic["nic_id"].(string); v != "" {
			ni.SetNicId(v)
		}
		if v := nic["secondary_private_ip_count"].(int); v > 0 {
			ni.SetSecondaryPrivateIpCount(int32(v))
		}
		if delete, deleteOK := d.GetOk(fmt.Sprintf("nics.%d.delete_on_vm_deletion", i)); deleteOK {
			log.Printf("[DEBUG] delete=%+v, deleteOK=%+v", delete, deleteOK)
			ni.SetDeleteOnVmDeletion(delete.(bool))
		}

		ni.SetDescription(nic["description"].(string))

		ni.SetPrivateIps(expandPrivatePublicIps(nic["private_ips"].(*schema.Set)))
		ni.SetSubnetId(nic["subnet_id"].(string))

		if sg := utils.InterfaceSliceToStringSlice(nic["security_group_ids"].([]interface{})); len(sg) > 0 {
			ni.SetSecurityGroupIds(sg)
		}

		if v, ok := d.GetOk("private_ip"); ok {
			ni.SetPrivateIps([]oscgo.PrivateIpLight{{
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

func expandPlacement(d *schema.ResourceData) (*oscgo.Placement, error) {
	placement := &oscgo.Placement{}

	subregionName, sOK := d.GetOk("placement_subregion_name")
	tenancy, tOK := d.GetOk("placement_tenancy")

	if sOK {
		placement.SetSubregionName(subregionName.(string))
	}
	if tOK {
		if v := tenancy.(string); v == "default" || v == "dedicated" {
			placement.SetTenancy(v)
		} else {
			return nil, errors.New("The value of field placement_tenancy can be only \"default\" or \"dedicated\"")
		}
	}
	if sOK || tOK {
		return placement, nil
	} else {
		return nil, nil
	}
}

func vmStateRefreshFunc(conn *oscgo.APIClient, instanceID, failState string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		var resp oscgo.ReadVmsResponse
		err := resource.Retry(30*time.Second, func() *resource.RetryError {
			rp, httpResp, err := conn.VmApi.ReadVms(context.Background()).ReadVmsRequest(oscgo.ReadVmsRequest{
				Filters: &oscgo.FiltersVm{
					VmIds: &[]string{instanceID},
				},
			}).Execute()
			if err != nil {
				return utils.CheckThrottling(httpResp, err)
			}
			resp = rp
			return nil
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
			return vm, state, fmt.Errorf("Failed to reach target state. Reason: %v", *vm.State)

		}

		return vm, state, nil
	}
}

func stopVM(vmID string, conn *oscgo.APIClient) error {
	vmResp, _, err := readVM(vmID, conn)
	if err != nil {
		return err
	}
	shutdownBehaviorOriginal := ""
	if len(vmResp.GetVms()) > 0 {
		if vmResp.GetVms()[0].GetVmInitiatedShutdownBehavior() != "stop" {
			shutdownBehaviorOriginal = vmResp.GetVms()[0].GetVmInitiatedShutdownBehavior()
			opts := oscgo.UpdateVmRequest{VmId: vmID}
			opts.SetVmInitiatedShutdownBehavior("stop")
			if err = updateVmAttr(conn, opts); err != nil {
				return err
			}
		}
	}

	err = resource.Retry(50*time.Second, func() *resource.RetryError {
		_, httpResp, err := conn.VmApi.StopVms(context.Background()).StopVmsRequest(oscgo.StopVmsRequest{
			VmIds: []string{vmID},
		}).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		return nil
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

	if shutdownBehaviorOriginal != "" {
		opts := oscgo.UpdateVmRequest{VmId: vmID}
		opts.SetVmInitiatedShutdownBehavior(shutdownBehaviorOriginal)
		if err = updateVmAttr(conn, opts); err != nil {
			return err
		}
	}

	return nil
}

func startVM(vmID string, conn *oscgo.APIClient) error {
	err := resource.Retry(50*time.Second, func() *resource.RetryError {
		_, httpResp, err := conn.VmApi.StartVms(context.Background()).StartVmsRequest(oscgo.StartVmsRequest{
			VmIds: []string{vmID},
		}).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		return nil
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
	err := resource.Retry(50*time.Second, func() *resource.RetryError {
		_, httpResp, err := conn.VmApi.UpdateVm(context.Background()).UpdateVmRequest(instanceAttrOpts).Execute()
		if err != nil {
			_, errBody := ioutil.ReadAll(httpResp.Body)
			if errBody != nil {
				fmt.Println(errBody)
			}
			return utils.CheckThrottling(httpResp, err)
		}
		return nil
	})

	if err != nil {
		return err
	}
	return nil
}

func readVM(vmID string, conn *oscgo.APIClient) (oscgo.ReadVmsResponse, *http.Response, error) {
	var resp oscgo.ReadVmsResponse
	var httpResult *http.Response
	err := resource.Retry(50*time.Second, func() *resource.RetryError {
		rp, httpResp, err := conn.VmApi.ReadVms(context.Background()).ReadVmsRequest(oscgo.ReadVmsRequest{
			Filters: &oscgo.FiltersVm{
				VmIds: &[]string{vmID},
			},
		}).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		resp = rp
		httpResult = httpResp
		return nil
	})
	return resp, httpResult, err
}

// AttributeSetter you can use this function to set the attributes
type AttributeSetter func(key string, value interface{}) error

func resourceDataAttrSetter(d *schema.ResourceData, callback func(AttributeSetter) error) error {
	setterFunc := func(key string, value interface{}) error {
		return d.Set(key, value)
	}
	return callback(setterFunc)
}
