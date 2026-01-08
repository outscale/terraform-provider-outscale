package oapi

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"reflect"
	"time"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/services/oapi/oapihelpers"
	"github.com/outscale/terraform-provider-outscale/internal/utils"
	"github.com/spf13/cast"
)

func ResourceOutscaleVM() *schema.Resource {
	return &schema.Resource{
		Create: resourceOAPIVMCreate,
		Read:   resourceOAPIVMRead,
		Update: resourceOAPIVMUpdate,
		Delete: resourceOAPIVMDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(12 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		// Schema: //omitted for brevity
		ValidateRawResourceConfigFuncs: []schema.ValidateRawResourceConfigFunc{
			validation.PreferWriteOnlyAttribute(cty.GetAttrPath("keypair_name"), cty.GetAttrPath("keypair_name_wo")),
		},
		Schema: map[string]*schema.Schema{
			"actions_on_next_boot": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"secure_boot": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"secure_boot_action": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"enable", "disable", "setup-mode", "none"}, false),
			},
			"block_device_mappings": {
				Type:     schema.TypeList,
				Optional: true,
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
										Computed: true,
									},
									"iops": {
										Type:     schema.TypeInt,
										Optional: true,
										ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
											iopsVal := val.(int)
											if int32(iopsVal) < MinIops || int32(iopsVal) > MaxIops {
												errs = append(errs, fmt.Errorf("%q must be between %d and %d inclusive, got: %d", key, MinIops, MaxIops, iopsVal))
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
										ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
											vSize, _ := val.(int)
											if int32(vSize) < 1 || int32(vSize) > MaxSize {
												errs = append(errs, fmt.Errorf("%q must be between 1 and %d gibibytes inclusive, got: %d", key, MaxSize, vSize))
											}
											return
										},
									},
									"volume_type": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"tags": TagsSchemaSDK(),
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
			"boot_mode": {
				Type:         schema.TypeString,
				Computed:     true,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"legacy", "uefi"}, false),
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
			"keypair_name_wo": {
				Type:      schema.TypeString,
				Optional:  true,
				WriteOnly: true,
			},
			"primary_nic": {
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
							ValidateFunc: func(number interface{}, key string) (warns []string, errs []error) {
								deviceNumber := number.(int)
								if deviceNumber != 0 {
									errs = append(errs, fmt.Errorf("%q in primary_nic must be only '0', got: %d", key, deviceNumber))
								}
								return
							},
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
							Type:     schema.TypeSet,
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
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"security_group_names": {
				Type:     schema.TypeSet,
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
							Type:     schema.TypeSet,
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
										Type:     schema.TypeString,
										Computed: true,
									},
									"volume_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"tags": TagsSchemaSDK(),
								},
							},
						},
						"device_name": {
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
			"tags": TagsSchemaSDK(),
		},
	}
}

func resourceOAPIVMCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*client.OutscaleClient).OSCAPI

	vmOpts, bsuMapsTags, err := buildCreateVmsRequest(d)
	if err != nil {
		return err
	}
	vState := d.Get("state").(string)
	if vState != "stopped" && vState != "running" {
		return errors.New("error: state must be `stopped or running`")
	}
	vmStateTarget := []string{"running"}
	if vState == "stopped" {
		vmStateTarget[0] = "stopped"
		vmOpts.BootOnCreation = oscgo.PtrBool(false)
	}

	// Create the vm
	var resp oscgo.CreateVmsResponse
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		rp, httpResp, err := conn.VmApi.CreateVms(context.Background()).CreateVmsRequest(vmOpts).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		resp = rp
		return nil
	})
	if err != nil {
		return fmt.Errorf("error launching source vm: %v", utils.GetErrorResponse(err))
	}

	if !resp.HasVms() || len(resp.GetVms()) == 0 {
		return errors.New("error launching source VM: no VMs returned in response")
	}

	vm := resp.GetVms()[0]

	d.SetId(vm.GetVmId())

	if get_psswd := d.Get("get_admin_password").(bool); get_psswd {
		psswd_err := retry.Retry(2500*time.Second, func() *retry.RetryError {
			psswd, err := getOAPIVMAdminPassword(vm.GetVmId(), conn)
			if err != nil || len(psswd) < 1 {
				return retry.RetryableError(errors.New("timeout awaiting windows password"))
			}
			return nil
		})
		if psswd_err != nil {
			return psswd_err
		}
	}

	if bsuMapsTags != nil {
		err := createBsuTags(conn, vm, bsuMapsTags)
		if err != nil {
			return err
		}
	}

	err = createOAPITagsSDK(conn, d)
	if err != nil {
		return err
	}

	stateConf := &retry.StateChangeConf{
		Pending:    []string{"pending", "ending/wait"},
		Target:     vmStateTarget,
		Refresh:    vmStateRefreshFunc(conn, vm.GetVmId(), "terminated"),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      15 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf(
			"error waiting for instance (%s) to become created: %s", d.Id(), err)
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

func createBsuTags(client *oscgo.APIClient, vm oscgo.Vm, bsuMapsTags []map[string]interface{}) error {
	for _, tMaps := range bsuMapsTags {
		for dName, tagsSchema := range tMaps {
			set := tagsSchema.(*schema.Set)
			tags := expandOAPITagsSDK(set)
			id := oapihelpers.GetBsuId(vm, dName)

			err := createOAPITags(context.Background(), client, tags, id)
			if err != nil {
				return fmt.Errorf("unable to create tags: %s", err)
			}
		}
	}
	return nil
}

func updateBsuTags(client *oscgo.APIClient, d *schema.ResourceData, addTags map[string]interface{}, delTags map[string]interface{}) error {
	var resp oscgo.ReadVmsResponse
	err := retry.Retry(60*time.Second, func() *retry.RetryError {
		rp, httpResp, err := client.VmApi.ReadVms(context.Background()).ReadVmsRequest(oscgo.ReadVmsRequest{
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
		return err
	}

	var empty []oscgo.ResourceTag
	for dName := range delTags {
		id := oapihelpers.GetBsuId(resp.GetVms()[0], dName)
		toRemove := expandOAPITagsSDK(delTags[dName].(*schema.Set))
		err := updateOAPITags(context.Background(), client, empty, toRemove, id)
		if err != nil {
			return err
		}
	}
	for dName := range addTags {
		id := oapihelpers.GetBsuId(resp.GetVms()[0], dName)
		toAdd := expandOAPITagsSDK(addTags[dName].(*schema.Set))
		err := updateOAPITags(context.Background(), client, toAdd, empty, id)
		if err != nil {
			return err
		}
	}
	return nil
}

func resourceOAPIVMRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*client.OutscaleClient).OSCAPI

	var resp oscgo.ReadVmsResponse
	err := retry.Retry(d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
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
		return fmt.Errorf("error reading the vm (%s): %s", d.Id(), err)
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
		if nics := buildNetworkOApiInterfaceOpts(d); len(nics) == 0 {
			if err := set("security_group_ids", getSecurityGroupIds(vm.GetSecurityGroups())); err != nil {
				return err
			}
		}
		d.SetId(vm.GetVmId())

		bsuTagsMaps, errTags := oapihelpers.GetBsuTagsMaps(vm, conn)
		if errTags != nil {
			return errTags
		}

		if err := d.Set("block_device_mappings_created", getOscAPIVMBlockDeviceMapping(
			bsuTagsMaps, vm.GetBlockDeviceMappings())); err != nil {
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
	err := retry.Retry(60*time.Second, func() *retry.RetryError {
		rp, httpResp, err := conn.VmApi.ReadAdminPassword(context.Background()).ReadAdminPasswordRequest(oscgo.ReadAdminPasswordRequest{VmId: VMID}).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		resp = rp
		return nil
	})
	if err != nil {
		return "", fmt.Errorf("error reading the vm's password %s", err)
	}
	return resp.GetAdminPassword(), nil
}

func findVolumeIdByDeviceName(d *schema.ResourceData, deviceName string) (string, error) {
	mappings := d.Get("block_device_mappings_created").([]any)
	for _, mapping := range mappings {
		mapping := mapping.(map[string]any)
		currentName := mapping["device_name"].(string)

		if deviceName == currentName {
			if bsuSet, ok := mapping["bsu"].(*schema.Set); ok && bsuSet.Len() > 0 {
				bsuList := bsuSet.List()
				if e, ok := bsuList[0].(map[string]any); ok {
					return e["volume_id"].(string), nil
				}
			}
		}
	}

	return "", fmt.Errorf("no volume found for device name %s", deviceName)
}

func resourceOAPIVMUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*client.OutscaleClient).OSCAPI
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

	if nothingToDo {
		return nil
	}

	updateRequest := oscgo.UpdateVmRequest{}
	mustStartVM := false
	if !d.IsNewResource() &&
		(d.HasChange("vm_type") || d.HasChange("user_data") ||
			d.HasChange("performance") || d.HasChange("nested_virtualization")) {

		if err := stopVM(id, conn, d.Timeout(schema.TimeoutUpdate)); err != nil {
			return err
		}
		mustStartVM = true

		if d.HasChange("vm_type") {
			updateRequest.SetVmType(d.Get("vm_type").(string))
		}

		if d.HasChange("user_data") {
			updateRequest.SetUserData(d.Get("user_data").(string))
		}

		if d.HasChange("performance") {
			updateRequest.SetPerformance(d.Get("performance").(string))
		}

		if d.HasChange("nested_virtualization") {
			updateRequest.SetNestedVirtualization(d.Get("nested_virtualization").(bool))
		}
	}

	if d.HasChange("deletion_protection") && !d.IsNewResource() {
		updateRequest.SetDeletionProtection(d.Get("deletion_protection").(bool))
	}

	if d.HasChange("keypair_name") && !d.IsNewResource() {
		updateRequest.SetKeypairName(d.Get("keypair_name").(string))
	}

	if d.HasChange("security_group_ids") && !d.IsNewResource() {
		updateRequest.SetSecurityGroupIds(utils.SetToStringSlice(d.Get("security_group_ids").(*schema.Set)))
	}

	if d.HasChange("security_group_names") && !d.IsNewResource() {
		updateRequest.SetSecurityGroupIds(utils.SetToStringSlice(d.Get("security_group_names").(*schema.Set)))
	}

	if d.HasChange("vm_initiated_shutdown_behavior") && !d.IsNewResource() {
		updateRequest.SetVmInitiatedShutdownBehavior(d.Get("vm_initiated_shutdown_behavior").(string))
	}

	if d.HasChange("is_source_dest_checked") && !d.IsNewResource() {
		updateRequest.SetIsSourceDestChecked(d.Get("is_source_dest_checked").(bool))
	}

	var updateBSUVolumeReqs []oscgo.UpdateVolumeRequest
	if d.HasChange("block_device_mappings") && !d.IsNewResource() {
		oldT, newT := d.GetChange("block_device_mappings")
		oldMapsTags, newMapsTags := getChangeTags(oldT, newT)
		if oldMapsTags != nil || newMapsTags != nil {
			if err := updateBsuTags(conn, d, oldMapsTags, newMapsTags); err != nil {
				return err
			}
		}

		oldMappings := oldT.([]any)
		newMappings := newT.([]any)
		var mappingsReqs []oscgo.BlockDeviceMappingVmUpdate
		for i, newMapping := range newMappings {
			oldMapping := oldMappings[i].(map[string]any)
			newMappingMap := newMapping.(map[string]any)

			hasMappingChanges := false
			updateMappingReq := oscgo.BlockDeviceMappingVmUpdate{}

			deviceName, okName := newMappingMap["device_name"]
			if okName && deviceName.(string) != "" {
				updateMappingReq.SetDeviceName(deviceName.(string))
			}

			if v, ok := newMappingMap["no_device"]; ok && v.(string) != "" {
				updateMappingReq.SetNoDevice(v.(string))
			}

			if v, ok := newMappingMap["virtual_device_name"]; ok && v.(string) != "" {
				updateMappingReq.SetVirtualDeviceName(v.(string))
			}

			if newBsu, ok := newMappingMap["bsu"].([]any); ok && len(newBsu) > 0 {
				newBsu := newBsu[0].(map[string]any)
				oldBsu := oldMapping["bsu"].([]any)[0].(map[string]any)

				updateBsuReq := oscgo.BsuToUpdateVm{}
				if deletion, ok := newBsu["delete_on_vm_deletion"].(bool); ok && oldBsu["delete_on_vm_deletion"].(bool) != deletion {
					updateBsuReq.SetDeleteOnVmDeletion(deletion)
					updateMappingReq.SetBsu(updateBsuReq)
					hasMappingChanges = true
				}

				hasVolumeChanges := false
				updateVolumeReq := oscgo.UpdateVolumeRequest{}

				if size, ok := newBsu["volume_size"]; ok && oldBsu["volume_size"].(int) != size.(int) && size.(int) > 0 {
					updateVolumeReq.SetSize(int32(size.(int)))
					hasVolumeChanges = true
				}
				if iops, ok := newBsu["iops"]; ok && oldBsu["iops"].(int) != iops.(int) && iops.(int) > 0 {
					updateVolumeReq.SetIops(int32(iops.(int)))
					hasVolumeChanges = true
				}
				if volType, ok := newBsu["volume_type"]; ok && oldBsu["volume_type"].(string) != volType.(string) && volType.(string) != "" {
					updateVolumeReq.SetVolumeType(volType.(string))
					hasVolumeChanges = true
				}

				if hasVolumeChanges && okName {
					id, err := findVolumeIdByDeviceName(d, deviceName.(string))
					if err != nil {
						return err
					}
					updateVolumeReq.SetVolumeId(id)

					updateBSUVolumeReqs = append(updateBSUVolumeReqs, updateVolumeReq)
				}
			}
			if hasMappingChanges {
				mappingsReqs = append(mappingsReqs, updateMappingReq)
			}
		}
		if len(mappingsReqs) > 0 {
			updateRequest.SetBlockDeviceMappings(mappingsReqs)
		}
	}

	if d.HasChange("secure_boot_action") && !d.IsNewResource() {
		if action := d.Get("secure_boot_action").(string); action != "" {
			bootAction := oscgo.SecureBootAction(action)
			updateRequest.ActionsOnNextBoot = &oscgo.ActionsOnNextBoot{
				SecureBoot: &bootAction,
			}
		}
	}

	if err := updateOAPITagsSDK(conn, d); err != nil {
		return err
	}

	if !reflect.ValueOf(updateRequest).IsZero() {
		updateRequest.SetVmId(id)
		if err := updateVmAttr(conn, updateRequest); err != nil {
			return utils.GetErrorResponse(err)
		}
	}

	if !onlyTags {
		if d.HasChange("state") && !d.IsNewResource() {
			upState := d.Get("state").(string)
			if upState != "stopped" && upState != "running" {
				return fmt.Errorf("error: state should be `stopped or running`")
			}
			mustStartVM = false
			if upState == "stopped" {
				if err := stopVM(id, conn, d.Timeout(schema.TimeoutUpdate)); err != nil {
					return err
				}
			} else {
				if err := startVM(id, conn, d.Timeout(schema.TimeoutUpdate)); err != nil {
					return err
				}
			}
		}
		if mustStartVM {
			if err := startVM(id, conn, d.Timeout(schema.TimeoutUpdate)); err != nil {
				return err
			}
		}

		var tasksIds []string
		for _, volumeReq := range updateBSUVolumeReqs {
			err := retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
				rp, httpResp, err := conn.VolumeApi.UpdateVolume(context.Background()).UpdateVolumeRequest(volumeReq).Execute()
				if err != nil {
					return utils.CheckThrottling(httpResp, err)
				}
				if vol, ok := rp.GetVolumeOk(); ok {
					if vol.GetTaskId() != "" {
						tasksIds = append(tasksIds, vol.GetTaskId())
					}
				}
				return nil
			})
			if err != nil {
				return err
			}
		}
		if len(tasksIds) > 0 {
			err := WaitForVolumeTasks(context.Background(), d.Timeout(schema.TimeoutUpdate), tasksIds, conn)
			if err != nil {
				return err
			}
		}
	}

	return resourceOAPIVMRead(d, meta)
}

func getChangeTags(oldCh interface{}, newCh interface{}) (map[string]interface{}, map[string]interface{}) {
	oldMapsTags := getbsuMapsTags(oldCh.([]interface{}))
	newMapsTags := getbsuMapsTags(newCh.([]interface{}))
	addMapsTags := make(map[string]interface{})
	delMapsTags := make(map[string]interface{})

	for v := range oldMapsTags {
		inter := oldMapsTags[v].(*schema.Set).Intersection(newMapsTags[v].(*schema.Set))
		if add := oldMapsTags[v].(*schema.Set).Difference(inter); len(add.List()) > 0 {
			addMapsTags[v] = add
		}
		if del := newMapsTags[v].(*schema.Set).Difference(inter); len(del.List()) > 0 {
			delMapsTags[v] = del
		}
	}
	return delMapsTags, addMapsTags
}

func getbsuMapsTags(changeMaps []interface{}) map[string]interface{} {
	mapsTags := make(map[string]interface{})

	for _, value := range changeMaps {
		val := value.(map[string]interface{})
		bsuMaps := val["bsu"].([]interface{})
		for _, v := range bsuMaps {
			bsu := v.(map[string]interface{})
			bsu_tags := bsu["tags"]
			mapsTags[val["device_name"].(string)] = bsu_tags
		}
	}
	return mapsTags
}

func resourceOAPIVMDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*client.OutscaleClient).OSCAPI
	id := d.Id()
	var err error

	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		_, httpResp, err := conn.VmApi.StopVms(context.Background()).StopVmsRequest(oscgo.StopVmsRequest{
			VmIds:     []string{id},
			ForceStop: oscgo.PtrBool(true),
		}).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("error force stopping vms before destroy %s", err)
	}

	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		_, httpResp, err := conn.VmApi.DeleteVms(context.Background()).DeleteVmsRequest(oscgo.DeleteVmsRequest{
			VmIds: []string{id},
		}).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("error deleting the vm")
	}

	log.Printf("[DEBUG] Waiting for VM (%s) to become terminated", id)

	stateConf := &retry.StateChangeConf{
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
			"error waiting for instance (%s) to terminate: %s", id, err)
	}

	return nil
}

func buildCreateVmsRequest(d *schema.ResourceData) (oscgo.CreateVmsRequest, []map[string]interface{}, error) {
	request := oscgo.CreateVmsRequest{
		DeletionProtection: oscgo.PtrBool(d.Get("deletion_protection").(bool)),
		BootOnCreation:     oscgo.PtrBool(true),
		MaxVmsCount:        oscgo.PtrInt32(1),
		MinVmsCount:        oscgo.PtrInt32(1),
		ImageId:            d.Get("image_id").(string),
	}

	placement, err := expandPlacement(d)
	if err != nil {
		return request, nil, err
	} else if placement != nil {
		request.SetPlacement(*placement)
	}

	subNet := d.Get("subnet_id").(string)
	if subNet != "" {
		request.SetSubnetId(subNet)
	}
	blockDevices, bsuMapsTags, err := expandBlockDeviceOApiMappings(d)
	if err != nil {
		return request, nil, err
	}
	if len(blockDevices) > 0 {
		request.SetBlockDeviceMappings(blockDevices)
	}

	if nics := buildNetworkOApiInterfaceOpts(d); len(nics) > 0 {
		if subNet != "" || placement != nil {
			return request, nil, errors.New("if you specify nics parameter, you must not specify subnet_id and placement parameters")
		}
		request.SetNics(nics)
	}

	if privateIPs := utils.InterfaceSliceToStringSlice(d.Get("private_ips").([]interface{})); len(privateIPs) > 0 {
		request.SetPrivateIps(privateIPs)
	}

	if sgIDs := utils.SetToStringSlice(d.Get("security_group_ids").(*schema.Set)); len(sgIDs) > 0 {
		request.SetSecurityGroupIds(sgIDs)
	}

	if sgNames := utils.SetToStringSlice(d.Get("security_group_names").(*schema.Set)); len(sgNames) > 0 {
		request.SetSecurityGroups(sgNames)
	}

	nestedVirtualization := d.Get("nested_virtualization").(bool)
	if tenacy := d.Get("placement_tenancy").(string); nestedVirtualization && tenacy != "dedicated" {
		return request, nil, errors.New("the field nested_virtualization can be true only if placement_tenancy is \"dedicated\"")
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

	if v := d.Get("boot_mode").(string); v != "" {
		action := (oscgo.BootMode)(d.Get("boot_mode").(string))
		request.SetBootMode(action)
	}
	if v := d.Get("secure_boot_action").(string); v != "" {
		action := (oscgo.SecureBootAction)(d.Get("secure_boot_action").(string))
		request.SetActionsOnNextBoot(oscgo.ActionsOnNextBoot{SecureBoot: &action})
	}

	kpName, diags := d.GetRawConfigAt(cty.GetAttrPath("keypair_name_wo"))
	if diags.HasError() {
		return request, bsuMapsTags, fmt.Errorf("error retrieving write-only argument: keypair_name_wo: %v", diags)
	}
	if !kpName.Type().Equals(cty.String) {
		return request, bsuMapsTags, errors.New("error retrieving write-only argument: password_wo, retrieved config value is not a string")
	}
	if !kpName.IsNull() {
		request.SetKeypairName(kpName.AsString())
	}

	return request, bsuMapsTags, nil
}

func expandBlockDeviceOApiMappings(d *schema.ResourceData) ([]oscgo.BlockDeviceMappingVmCreation, []map[string]interface{}, error) {
	var blockDevices []oscgo.BlockDeviceMappingVmCreation
	block := d.Get("block_device_mappings").([]any)
	bsuMapsTags := make([]map[string]any, len(block))
	for k, v := range block {
		blockDevice := oscgo.BlockDeviceMappingVmCreation{}
		value := v.(map[string]any)
		if bsu := value["bsu"].([]any); len(bsu) > 0 {
			expandBSU, mapsTags, err := expandBlockDeviceBSU(bsu[0].(map[string]any), value["device_name"].(string))
			bsuMapsTags[k] = mapsTags
			if err != nil {
				return nil, nil, err
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
	return blockDevices, bsuMapsTags, nil
}

func expandBlockDeviceBSU(bsu map[string]interface{}, deviceName string) (oscgo.BsuToCreate, map[string]interface{}, error) {
	bsuMapsTags := make(map[string]interface{})
	bsuToCreate := oscgo.BsuToCreate{}
	snapshotID := bsu["snapshot_id"].(string)
	volumeType := bsu["volume_type"].(string)
	volumeSize := int32(bsu["volume_size"].(int))

	if snapshotID == "" && volumeSize == 0 {
		return bsuToCreate, nil, fmt.Errorf("error: 'volume_size' parameter is required if the volume is not created from a snapshot (snapshotid unspecified)")
	}
	if iops := bsu["iops"]; iops.(int) > 0 {
		if volumeType != "io1" {
			return bsuToCreate, nil, ErrResourceInvalidIOPS
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
	if bsu_tags := bsu["tags"]; len(bsu_tags.(*schema.Set).List()) != 0 {
		bsuMapsTags[deviceName] = bsu_tags
	}

	return bsuToCreate, bsuMapsTags, nil
}

func buildNetworkOApiInterfaceOpts(d *schema.ResourceData) []oscgo.NicForVmCreation {
	networkInterfaces := []oscgo.NicForVmCreation{}
	if nics := d.Get("primary_nic").(*schema.Set).List(); len(nics) > 0 {
		buildNicForVmCreation(nics, &networkInterfaces)
	}
	if nics := d.Get("nics").(*schema.Set).List(); len(nics) > 0 {
		buildNicForVmCreation(nics, &networkInterfaces)
	}
	return networkInterfaces
}

func buildNicForVmCreation(nics []interface{}, listNics *[]oscgo.NicForVmCreation) {
	for _, v := range nics {
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
		if v := nic["delete_on_vm_deletion"]; v != nil {
			ni.SetDeleteOnVmDeletion(v.(bool))
		}
		ni.SetDescription(nic["description"].(string))
		ni.SetPrivateIps(expandPrivatePublicIps(nic["private_ips"].(*schema.Set)))
		ni.SetSubnetId(nic["subnet_id"].(string))

		if sg := utils.InterfaceSliceToStringSlice(nic["security_group_ids"].([]interface{})); len(sg) > 0 {
			ni.SetSecurityGroupIds(sg)
		}
		*listNics = append(*listNics, ni)
	}
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
			return nil, errors.New("the value of field placement_tenancy can be only \"default\" or \"dedicated\"")
		}
	}
	if sOK || tOK {
		return placement, nil
	} else {
		return nil, nil
	}
}

func vmStateRefreshFunc(conn *oscgo.APIClient, instanceID, failState string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		var resp oscgo.ReadVmsResponse
		err := retry.Retry(30*time.Second, func() *retry.RetryError {
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
			return vm, state, fmt.Errorf("failed to reach target state:: %v", *vm.State)
		}

		return vm, state, nil
	}
}

func stopVM(vmID string, conn *oscgo.APIClient, timeOut time.Duration) error {
	vmResp, _, err := readVM(vmID, conn, timeOut)
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

	err = retry.Retry(timeOut, func() *retry.RetryError {
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

	stateConf := &retry.StateChangeConf{
		Pending:    []string{"pending", "running", "shutting-down", "stopped", "stopping"},
		Target:     []string{"stopped"},
		Refresh:    vmStateRefreshFunc(conn, vmID, ""),
		Timeout:    timeOut,
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("error waiting for instance (%s) to stop: %s", vmID, err)
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

func startVM(vmID string, conn *oscgo.APIClient, timeOut time.Duration) error {
	err := retry.Retry(timeOut, func() *retry.RetryError {
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

	stateConf := &retry.StateChangeConf{
		Pending:    []string{"pending", "stopped"},
		Target:     []string{"running"},
		Refresh:    vmStateRefreshFunc(conn, vmID, ""),
		Timeout:    timeOut,
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("error waiting for instance (%s) to become ready: %s", vmID, err)
	}

	return nil
}

func updateVmAttr(conn *oscgo.APIClient, instanceAttrOpts oscgo.UpdateVmRequest) error {
	err := retry.Retry(50*time.Second, func() *retry.RetryError {
		_, httpResp, err := conn.VmApi.UpdateVm(context.Background()).UpdateVmRequest(instanceAttrOpts).Execute()
		if err != nil {
			_, errBody := io.ReadAll(httpResp.Body)
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

func readVM(vmID string, conn *oscgo.APIClient, timeOut time.Duration) (oscgo.ReadVmsResponse, *http.Response, error) {
	var resp oscgo.ReadVmsResponse
	var httpResult *http.Response
	err := retry.Retry(timeOut, func() *retry.RetryError {
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
