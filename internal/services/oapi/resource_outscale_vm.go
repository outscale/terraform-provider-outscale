package oapi

import (
	"context"
	"errors"
	"fmt"
	"log"
	"reflect"
	"time"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/outscale/goutils/sdk/ptr"
	"github.com/outscale/osc-sdk-go/v3/pkg/options"
	"github.com/outscale/osc-sdk-go/v3/pkg/osc"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/services/oapi/oapihelpers"
	"github.com/outscale/terraform-provider-outscale/internal/utils"
	"github.com/spf13/cast"
)

func ResourceOutscaleVM() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceOAPIVMCreate,
		ReadContext:   resourceOAPIVMRead,
		UpdateContext: resourceOAPIVMUpdate,
		DeleteContext: resourceOAPIVMDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(15 * time.Minute),
			Read:   schema.DefaultTimeout(ReadDefaultTimeout),
			Update: schema.DefaultTimeout(UpdateDefaultTimeout),
			Delete: schema.DefaultTimeout(DeleteDefaultTimeout),
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
			"tpm_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
				Optional: true,
				ForceNew: true,
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

func resourceOAPIVMCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.OutscaleClient).OSC

	vmOpts, bsuMapsTags, err := buildCreateVmsRequest(d)
	if err != nil {
		return diag.FromErr(err)
	}
	vState := d.Get("state").(string)
	if vState != "stopped" && vState != "running" {
		return diag.FromErr(errors.New("error: state must be `stopped or running`"))
	}
	vmStateTarget := []string{"running"}
	if vState == "stopped" {
		vmStateTarget[0] = "stopped"
		vmOpts.BootOnCreation = new(false)
	}

	timeout := d.Timeout(schema.TimeoutCreate)

	// Create the vm
	resp, err := client.CreateVms(ctx, vmOpts, options.WithRetryTimeout(timeout))
	if err != nil {
		return diag.Errorf("error launching source vm: %v", err)
	}

	if resp.Vms == nil || len(*resp.Vms) == 0 {
		return diag.Errorf("error launching source VM: no VMs returned in response")
	}

	vm := (*resp.Vms)[0]

	d.SetId(vm.VmId)

	if get_psswd := d.Get("get_admin_password").(bool); get_psswd {
		psswd, err := getOAPIVMAdminPassword(ctx, client, vm.VmId, timeout)
		if err != nil || len(psswd) < 1 {
			return diag.Errorf("timeout awaiting windows password")
		}
	}

	if bsuMapsTags != nil {
		err := createBsuTags(ctx, client, timeout, vm, bsuMapsTags)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	err = createOAPITagsSDK(ctx, client, timeout, d)
	if err != nil {
		return diag.FromErr(err)
	}

	stateConf := &retry.StateChangeConf{
		Pending: []string{"pending", "ending/wait"},
		Target:  vmStateTarget,
		Timeout: timeout,
		Refresh: vmStateRefreshFunc(ctx, client, vm.VmId, "terminated", timeout),
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf(
			"error waiting for instance (%s) to become created: %s", d.Id(), err)
	}

	// Initialize the connection info
	if vm.PublicIp != nil {
		d.SetConnInfo(map[string]string{
			"type": "ssh",
			"host": ptr.From(vm.PublicIp),
		})
	} else {
		d.SetConnInfo(map[string]string{
			"type": "ssh",
			"host": vm.PrivateIp,
		})
	}

	if v, exist := d.GetOkExists("is_source_dest_checked"); exist {
		opts := osc.UpdateVmRequest{
			VmId: vm.VmId,
		}
		opts.IsSourceDestChecked = new(v.(bool))
		if err := updateVmAttr(ctx, client, timeout, opts); err != nil {
			return diag.FromErr(err)
		}
	}
	return resourceOAPIVMRead(ctx, d, meta)
}

func createBsuTags(ctx context.Context, client *osc.Client, timeout time.Duration, vm osc.Vm, bsuMapsTags []map[string]interface{}) error {
	for _, tMaps := range bsuMapsTags {
		for dName, tagsSchema := range tMaps {
			set := tagsSchema.(*schema.Set)
			tags := expandOAPITagsSDK(set)
			id := oapihelpers.GetBsuId(vm, dName)

			err := createOAPITags(ctx, client, timeout, tags, id)
			if err != nil {
				return fmt.Errorf("unable to create tags: %s", err)
			}
		}
	}
	return nil
}

func updateBsuTags(ctx context.Context, client *osc.Client, timeout time.Duration, d *schema.ResourceData, addTags map[string]interface{}, delTags map[string]interface{}) error {
	resp, err := client.ReadVms(ctx, osc.ReadVmsRequest{
		Filters: &osc.FiltersVm{
			VmIds: &[]string{d.Id()},
		},
	}, options.WithRetryTimeout(timeout))
	if err != nil {
		return err
	}

	vms := ptr.From(resp.Vms)
	var empty []osc.ResourceTag
	for dName := range delTags {
		id := oapihelpers.GetBsuId(vms[0], dName)
		toRemove := expandOAPITagsSDK(delTags[dName].(*schema.Set))
		err := updateOAPITags(ctx, client, timeout, empty, toRemove, id)
		if err != nil {
			return err
		}
	}
	for dName := range addTags {
		id := oapihelpers.GetBsuId(vms[0], dName)
		toAdd := expandOAPITagsSDK(addTags[dName].(*schema.Set))
		err := updateOAPITags(ctx, client, timeout, toAdd, empty, id)
		if err != nil {
			return err
		}
	}
	return nil
}

func resourceOAPIVMRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.OutscaleClient).OSC

	timeout := d.Timeout(schema.TimeoutRead)

	resp, err := client.ReadVms(ctx, osc.ReadVmsRequest{
		Filters: &osc.FiltersVm{
			VmIds: &[]string{d.Id()},
		},
	}, options.WithRetryTimeout(timeout))
	if err != nil {
		return diag.Errorf("error reading the vm (%s): %s", d.Id(), err)
	}
	if resp.Vms == nil || utils.IsResponseEmpty(len(*resp.Vms), "Snapshot", d.Id()) {
		d.SetId("")
		return nil
	}

	vm := (*resp.Vms)[0]
	if vm.State == "terminated" {
		utils.LogManuallyDeleted("Vm", d.Id())
		d.SetId("")
		return nil
	}
	adminPassword, err := getOAPIVMAdminPassword(ctx, client, vm.VmId, timeout)
	if err != nil {
		return diag.FromErr(err)
	}
	bsu := d.Get("bsu_optimized")
	if err := resourceDataAttrSetter(d, func(set AttributeSetter) error {
		if err := d.Set("admin_password", adminPassword); err != nil {
			return err
		}
		if nics := buildNetworkOApiInterfaceOpts(d); len(nics) == 0 {
			if err := set("security_group_ids", getSecurityGroupIds(vm.SecurityGroups)); err != nil {
				return err
			}
		}
		d.SetId(vm.VmId)

		bsuTagsMaps, errTags := oapihelpers.GetBsuTagsMaps(ctx, client, timeout, vm)
		if errTags != nil {
			return errTags
		}

		if err := d.Set("block_device_mappings_created", getOscAPIVMBlockDeviceMapping(
			bsuTagsMaps, vm.BlockDeviceMappings)); err != nil {
			return err
		}

		return oapiVMDescriptionAttributes(set, &vm)
	}); err != nil {
		return diag.FromErr(err)
	}
	return diag.FromErr(d.Set("bsu_optimized", bsu))
}

func getOAPIVMAdminPassword(ctx context.Context, client *osc.Client, VMID string, timeout time.Duration) (string, error) {
	resp, err := client.ReadAdminPassword(ctx, osc.ReadAdminPasswordRequest{VmId: VMID}, options.WithRetryTimeout(timeout))
	if err != nil {
		return "", fmt.Errorf("error reading the vm's password %s", err)
	}
	return ptr.From(resp.AdminPassword), nil
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

func resourceOAPIVMUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.OutscaleClient).OSC

	timeout := d.Timeout(schema.TimeoutUpdate)
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

	updateRequest := osc.UpdateVmRequest{}
	mustStartVM := false
	if !d.IsNewResource() &&
		(d.HasChange("vm_type") || d.HasChange("user_data") ||
			d.HasChange("performance") || d.HasChange("nested_virtualization")) {

		if err := stopVM(ctx, client, timeout, id); err != nil {
			return diag.FromErr(err)
		}
		mustStartVM = true

		if d.HasChange("vm_type") {
			updateRequest.VmType = new(d.Get("vm_type").(string))
		}

		if d.HasChange("user_data") {
			updateRequest.UserData = new(d.Get("user_data").(string))
		}

		if d.HasChange("performance") {
			updateRequest.Performance = new(osc.UpdateVmRequestPerformance(d.Get("performance").(string)))
		}

		if d.HasChange("nested_virtualization") {
			updateRequest.NestedVirtualization = new(d.Get("nested_virtualization").(bool))
		}
	}

	if d.HasChange("deletion_protection") && !d.IsNewResource() {
		updateRequest.DeletionProtection = new(d.Get("deletion_protection").(bool))
	}

	if d.HasChange("keypair_name") && !d.IsNewResource() {
		updateRequest.KeypairName = new(d.Get("keypair_name").(string))
	}

	if d.HasChange("security_group_ids") && !d.IsNewResource() {
		updateRequest.SecurityGroupIds = utils.SetToStringSlice(d.Get("security_group_ids").(*schema.Set))
	}

	if d.HasChange("security_group_names") && !d.IsNewResource() {
		updateRequest.SecurityGroupIds = utils.SetToStringSlice(d.Get("security_group_names").(*schema.Set))
	}

	if d.HasChange("vm_initiated_shutdown_behavior") && !d.IsNewResource() {
		updateRequest.VmInitiatedShutdownBehavior = new(d.Get("vm_initiated_shutdown_behavior").(string))
	}

	if d.HasChange("is_source_dest_checked") && !d.IsNewResource() {
		updateRequest.IsSourceDestChecked = new(d.Get("is_source_dest_checked").(bool))
	}

	var updateBSUVolumeReqs []osc.UpdateVolumeRequest
	if d.HasChange("block_device_mappings") && !d.IsNewResource() {
		oldT, newT := d.GetChange("block_device_mappings")
		oldMapsTags, newMapsTags := getChangeTags(oldT, newT)
		if oldMapsTags != nil || newMapsTags != nil {
			if err := updateBsuTags(ctx, client, timeout, d, oldMapsTags, newMapsTags); err != nil {
				return diag.FromErr(err)
			}
		}

		oldMappings := oldT.([]any)
		newMappings := newT.([]any)
		var mappingsReqs []osc.BlockDeviceMappingVmUpdate
		for i, newMapping := range newMappings {
			oldMapping := oldMappings[i].(map[string]any)
			newMappingMap := newMapping.(map[string]any)

			hasMappingChanges := false
			updateMappingReq := osc.BlockDeviceMappingVmUpdate{}

			deviceName, okName := newMappingMap["device_name"]
			if okName && deviceName.(string) != "" {
				updateMappingReq.DeviceName = new(deviceName.(string))
			}

			if v, ok := newMappingMap["no_device"]; ok && v.(string) != "" {
				updateMappingReq.NoDevice = new(v.(string))
			}

			if v, ok := newMappingMap["virtual_device_name"]; ok && v.(string) != "" {
				updateMappingReq.VirtualDeviceName = new(v.(string))
			}

			if newBsu, ok := newMappingMap["bsu"].([]any); ok && len(newBsu) > 0 {
				newBsu := newBsu[0].(map[string]any)
				oldBsu := oldMapping["bsu"].([]any)[0].(map[string]any)

				updateBsuReq := osc.BsuToUpdateVm{}
				if deletion, ok := newBsu["delete_on_vm_deletion"].(bool); ok && oldBsu["delete_on_vm_deletion"].(bool) != deletion {
					updateBsuReq.DeleteOnVmDeletion = deletion
					updateMappingReq.Bsu = &updateBsuReq
					hasMappingChanges = true
				}

				hasVolumeChanges := false
				updateVolumeReq := osc.UpdateVolumeRequest{}

				if size, ok := newBsu["volume_size"]; ok && oldBsu["volume_size"].(int) != size.(int) && size.(int) > 0 {
					updateVolumeReq.Size = new(size.(int))
					hasVolumeChanges = true
				}
				if iops, ok := newBsu["iops"]; ok && oldBsu["iops"].(int) != iops.(int) && iops.(int) > 0 {
					updateVolumeReq.Iops = new(iops.(int))
					hasVolumeChanges = true
				}
				if volType, ok := newBsu["volume_type"]; ok && oldBsu["volume_type"].(string) != volType.(string) && volType.(string) != "" {
					updateVolumeReq.VolumeType = new(osc.VolumeType(volType.(string)))
					hasVolumeChanges = true
				}

				if hasVolumeChanges && okName {
					id, err := findVolumeIdByDeviceName(d, deviceName.(string))
					if err != nil {
						return diag.FromErr(err)
					}
					updateVolumeReq.VolumeId = id

					updateBSUVolumeReqs = append(updateBSUVolumeReqs, updateVolumeReq)
				}
			}
			if hasMappingChanges {
				mappingsReqs = append(mappingsReqs, updateMappingReq)
			}
		}
		if len(mappingsReqs) > 0 {
			updateRequest.BlockDeviceMappings = mappingsReqs
		}
	}

	if d.HasChange("secure_boot_action") && !d.IsNewResource() {
		if action := d.Get("secure_boot_action").(string); action != "" {
			bootAction := osc.SecureBootAction(action)
			updateRequest.ActionsOnNextBoot = &osc.ActionsOnNextBoot{
				SecureBoot: &bootAction,
			}
		}
	}

	if err := updateOAPITagsSDK(ctx, client, timeout, d); err != nil {
		return diag.FromErr(err)
	}

	if !reflect.ValueOf(updateRequest).IsZero() {
		updateRequest.VmId = id
		if err := updateVmAttr(ctx, client, timeout, updateRequest); err != nil {
			return diag.FromErr(err)
		}
	}

	if !onlyTags {
		if d.HasChange("state") && !d.IsNewResource() {
			upState := d.Get("state").(string)
			if upState != "stopped" && upState != "running" {
				return diag.Errorf("error: state should be `stopped or running`")
			}
			mustStartVM = false
			if upState == "stopped" {
				if err := stopVM(ctx, client, timeout, id); err != nil {
					return diag.FromErr(err)
				}
			} else {
				if err := startVM(ctx, client, timeout, id); err != nil {
					return diag.FromErr(err)
				}
			}
		}
		if mustStartVM {
			if err := startVM(ctx, client, timeout, id); err != nil {
				return diag.FromErr(err)
			}
		}

		var tasksIds []string
		for _, volumeReq := range updateBSUVolumeReqs {
			rp, err := client.UpdateVolume(ctx, volumeReq, options.WithRetryTimeout(timeout))
			if err != nil {
				return diag.FromErr(err)
			}
			if rp.Volume != nil {
				if ptr.From(rp.Volume.TaskId) != "" {
					tasksIds = append(tasksIds, ptr.From(rp.Volume.TaskId))
				}
			}
		}
		if len(tasksIds) > 0 {
			err := WaitForVolumeTasks(ctx, timeout, tasksIds, client)
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	return resourceOAPIVMRead(ctx, d, meta)
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

func resourceOAPIVMDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.OutscaleClient).OSC

	timeout := d.Timeout(schema.TimeoutDelete)
	id := d.Id()

	_, err := client.StopVms(ctx, osc.StopVmsRequest{
		VmIds:     []string{id},
		ForceStop: new(true),
	}, options.WithRetryTimeout(timeout))
	if err != nil {
		return diag.Errorf("error force stopping vms before destroy %s", err)
	}

	_, err = client.DeleteVms(ctx, osc.DeleteVmsRequest{
		VmIds: []string{id},
	}, options.WithRetryTimeout(timeout))
	if err != nil {
		return diag.Errorf("error deleting the vm")
	}

	log.Printf("[DEBUG] Waiting for VM (%s) to become terminated", id)

	stateConf := &retry.StateChangeConf{
		Pending: []string{"pending", "running", "shutting-down", "stopped", "stopping"},
		Target:  []string{"terminated"},
		Timeout: timeout,
		Refresh: vmStateRefreshFunc(ctx, client, id, "", timeout),
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf(
			"error waiting for instance (%s) to terminate: %s", id, err)
	}

	return nil
}

func buildCreateVmsRequest(d *schema.ResourceData) (osc.CreateVmsRequest, []map[string]interface{}, error) {
	request := osc.CreateVmsRequest{
		DeletionProtection: new(d.Get("deletion_protection").(bool)),
		BootOnCreation:     new(true),
		MaxVmsCount:        new(1),
		MinVmsCount:        new(1),
		ImageId:            d.Get("image_id").(string),
	}

	placement, err := expandPlacement(d)
	if err != nil {
		return request, nil, err
	} else if placement != nil {
		request.Placement = placement
	}

	subNet := d.Get("subnet_id").(string)
	if subNet != "" {
		request.SubnetId = &subNet
	}
	blockDevices, bsuMapsTags, err := expandBlockDeviceOApiMappings(d)
	if err != nil {
		return request, nil, err
	}
	if len(blockDevices) > 0 {
		request.BlockDeviceMappings = blockDevices
	}

	if nics := buildNetworkOApiInterfaceOpts(d); len(nics) > 0 {
		if subNet != "" || placement != nil {
			return request, nil, errors.New("if you specify nics parameter, you must not specify subnet_id and placement parameters")
		}
		request.Nics = nics
	}

	if privateIPs := utils.InterfaceSliceToStringSlice(d.Get("private_ips").([]interface{})); len(privateIPs) > 0 {
		request.PrivateIps = privateIPs
	}

	if sgIDs := utils.SetToStringSlice(d.Get("security_group_ids").(*schema.Set)); len(sgIDs) > 0 {
		request.SecurityGroupIds = sgIDs
	}

	if sgNames := utils.SetToStringSlice(d.Get("security_group_names").(*schema.Set)); len(sgNames) > 0 {
		request.SecurityGroups = sgNames
	}

	nestedVirtualization := d.Get("nested_virtualization").(bool)
	if tenacy := d.Get("placement_tenancy").(string); nestedVirtualization && tenacy != "dedicated" {
		return request, nil, errors.New("the field nested_virtualization can be true only if placement_tenancy is \"dedicated\"")
	}
	request.NestedVirtualization = &nestedVirtualization

	if v := d.Get("user_data").(string); v != "" {
		request.UserData = &v
	}

	if v := d.Get("vm_type").(string); v != "" {
		request.VmType = &v
	}

	if v := d.Get("client_token").(string); v != "" {
		request.ClientToken = &v
	}

	if v := d.Get("keypair_name").(string); v != "" {
		request.KeypairName = &v
	}
	if v, ok := d.GetOk("vm_initiated_shutdown_behavior"); ok && v != "" {
		request.VmInitiatedShutdownBehavior = new(v.(string))
	}

	if v := d.Get("performance").(string); v != "" {
		request.Performance = new(osc.CreateVmsRequestPerformance(v))
	}

	if v := d.Get("boot_mode").(string); v != "" {
		action := (osc.BootMode)(d.Get("boot_mode").(string))
		request.BootMode = &action
	}
	if v := d.Get("secure_boot_action").(string); v != "" {
		action := (osc.SecureBootAction)(d.Get("secure_boot_action").(string))
		request.ActionsOnNextBoot = &osc.ActionsOnNextBoot{SecureBoot: &action}
	}
	tpmEnabled := d.GetRawConfig().GetAttr("tpm_enabled")
	if !tpmEnabled.IsNull() {
		request.TpmEnabled = new(tpmEnabled.True())
	}

	kpName, diags := d.GetRawConfigAt(cty.GetAttrPath("keypair_name_wo"))
	if diags.HasError() {
		return request, bsuMapsTags, fmt.Errorf("error retrieving write-only argument: keypair_name_wo: %v", diags)
	}
	if !kpName.Type().Equals(cty.String) {
		return request, bsuMapsTags, errors.New("error retrieving write-only argument: password_wo, retrieved config value is not a string")
	}
	if !kpName.IsNull() {
		request.KeypairName = new(kpName.AsString())
	}

	return request, bsuMapsTags, nil
}

func expandBlockDeviceOApiMappings(d *schema.ResourceData) ([]osc.BlockDeviceMappingVmCreation, []map[string]interface{}, error) {
	var blockDevices []osc.BlockDeviceMappingVmCreation
	block := d.Get("block_device_mappings").([]any)
	bsuMapsTags := make([]map[string]any, len(block))
	for k, v := range block {
		blockDevice := osc.BlockDeviceMappingVmCreation{}
		value := v.(map[string]any)
		if bsu := value["bsu"].([]any); len(bsu) > 0 {
			expandBSU, mapsTags, err := expandBlockDeviceBSU(bsu[0].(map[string]any), value["device_name"].(string))
			bsuMapsTags[k] = mapsTags
			if err != nil {
				return nil, nil, err
			}
			blockDevice.Bsu = &expandBSU
		}
		if deviceName, ok := value["device_name"]; ok && deviceName != "" {
			blockDevice.DeviceName = new(cast.ToString(deviceName))
		}
		if noDevice, ok := value["no_device"]; ok && noDevice != "" {
			blockDevice.NoDevice = new(cast.ToString(noDevice))
		}
		if virtualDeviceName, ok := value["virtual_device_name"]; ok && virtualDeviceName != "" {
			blockDevice.VirtualDeviceName = new(cast.ToString(virtualDeviceName))
		}
		blockDevices = append(blockDevices, blockDevice)
	}
	return blockDevices, bsuMapsTags, nil
}

func expandBlockDeviceBSU(bsu map[string]interface{}, deviceName string) (osc.BsuToCreate, map[string]interface{}, error) {
	bsuMapsTags := make(map[string]interface{})
	bsuToCreate := osc.BsuToCreate{}
	snapshotID := bsu["snapshot_id"].(string)
	volumeType := bsu["volume_type"].(string)
	volumeSize := bsu["volume_size"].(int)

	if snapshotID == "" && volumeSize == 0 {
		return bsuToCreate, nil, fmt.Errorf("error: 'volume_size' parameter is required if the volume is not created from a snapshot (snapshotid unspecified)")
	}
	if iops := bsu["iops"]; iops.(int) > 0 {
		if volumeType != "io1" {
			return bsuToCreate, nil, ErrResourceInvalidIOPS
		}
		bsuToCreate.Iops = new(iops.(int))
	} else {
		delete(bsu, "iops")
	}
	if snapshotID != "" {
		bsuToCreate.SnapshotId = &snapshotID
	}
	if volumeSize > 0 {
		bsuToCreate.VolumeSize = &volumeSize
	}
	if volumeType != "" {
		bsuToCreate.VolumeType = new(osc.VolumeType(volumeType))
	}
	if deleteOnVMDeletion, ok := bsu["delete_on_vm_deletion"]; ok && deleteOnVMDeletion != "" {
		bsuToCreate.DeleteOnVmDeletion = new(cast.ToBool(deleteOnVMDeletion))
	}
	if bsu_tags := bsu["tags"]; len(bsu_tags.(*schema.Set).List()) != 0 {
		bsuMapsTags[deviceName] = bsu_tags
	}

	return bsuToCreate, bsuMapsTags, nil
}

func buildNetworkOApiInterfaceOpts(d *schema.ResourceData) []osc.NicForVmCreation {
	networkInterfaces := []osc.NicForVmCreation{}
	if nics := d.Get("primary_nic").(*schema.Set).List(); len(nics) > 0 {
		buildNicForVmCreation(nics, &networkInterfaces)
	}
	if nics := d.Get("nics").(*schema.Set).List(); len(nics) > 0 {
		buildNicForVmCreation(nics, &networkInterfaces)
	}
	return networkInterfaces
}

func buildNicForVmCreation(nics []interface{}, listNics *[]osc.NicForVmCreation) {
	for _, v := range nics {
		nic := v.(map[string]interface{})
		ni := osc.NicForVmCreation{
			DeviceNumber: new(nic["device_number"].(int)),
		}

		if v := nic["nic_id"].(string); v != "" {
			ni.NicId = &v
		}
		if v := nic["secondary_private_ip_count"].(int); v > 0 {
			ni.SecondaryPrivateIpCount = &v
		}
		if v := nic["delete_on_vm_deletion"]; v != nil {
			ni.DeleteOnVmDeletion = new(v.(bool))
		}
		ni.Description = new(nic["description"].(string))
		ni.PrivateIps = new(expandPrivatePublicIps(nic["private_ips"].(*schema.Set)))
		ni.SubnetId = new(nic["subnet_id"].(string))

		if sg := utils.InterfaceSliceToStringSlice(nic["security_group_ids"].([]interface{})); len(sg) > 0 {
			ni.SecurityGroupIds = &sg
		}
		*listNics = append(*listNics, ni)
	}
}

func expandPrivatePublicIps(p *schema.Set) []osc.PrivateIpLight {
	privatePublicIPS := make([]osc.PrivateIpLight, len(p.List()))

	for i, v := range p.List() {
		value := v.(map[string]interface{})
		privatePublicIPS[i].IsPrimary = value["is_primary"].(bool)
		privatePublicIPS[i].PrivateIp = value["private_ip"].(string)
	}
	return privatePublicIPS
}

func expandPlacement(d *schema.ResourceData) (*osc.Placement, error) {
	placement := &osc.Placement{}

	subregionName, sOK := d.GetOk("placement_subregion_name")
	tenancy, tOK := d.GetOk("placement_tenancy")

	if sOK {
		placement.SubregionName = subregionName.(string)
	}
	if tOK {
		if v := tenancy.(string); v == "default" || v == "dedicated" {
			placement.Tenancy = v
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

func vmStateRefreshFunc(ctx context.Context, client *osc.Client, instanceID, failState string, timeout time.Duration) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := client.ReadVms(ctx, osc.ReadVmsRequest{
			Filters: &osc.FiltersVm{
				VmIds: &[]string{instanceID},
			},
		}, options.WithRetryTimeout(timeout))
		if err != nil {
			log.Printf("[ERROR] error on InstanceStateRefresh: %s", err)
			return nil, "", err
		}

		if resp.Vms == nil {
			return nil, "", nil
		}

		vm := (*resp.Vms)[0]
		state := string(vm.State)

		if state == failState {
			return vm, state, fmt.Errorf("failed to reach target state: %v", state)
		}

		return vm, state, nil
	}
}

func stopVM(ctx context.Context, client *osc.Client, timeout time.Duration, id string) error {
	resp, err := readVM(ctx, client, timeout, id)
	if err != nil {
		return err
	}
	if resp.Vms == nil || len(*resp.Vms) == 0 {
		return fmt.Errorf("no VM found with ID %s", id)
	}

	vm := (*resp.Vms)[0]
	shutdownBehaviorOriginal := ""
	if vm.VmInitiatedShutdownBehavior != "stop" {
		shutdownBehaviorOriginal = vm.VmInitiatedShutdownBehavior
		opts := osc.UpdateVmRequest{VmId: id}
		opts.VmInitiatedShutdownBehavior = new("stop")
		if err = updateVmAttr(ctx, client, timeout, opts); err != nil {
			return err
		}
	}

	_, err = client.StopVms(ctx, osc.StopVmsRequest{
		VmIds: []string{id},
	}, options.WithRetryTimeout(timeout))
	if err != nil {
		return fmt.Errorf("error stopping vms %s", err)
	}

	stateConf := &retry.StateChangeConf{
		Pending: []string{"pending", "running", "shutting-down", "stopped", "stopping"},
		Target:  []string{"stopped"},
		Timeout: timeout,
		Refresh: vmStateRefreshFunc(ctx, client, id, "", timeout),
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("error waiting for instance (%s) to stop: %s", id, err)
	}

	if shutdownBehaviorOriginal != "" {
		opts := osc.UpdateVmRequest{VmId: id}
		opts.VmInitiatedShutdownBehavior = &shutdownBehaviorOriginal
		if err = updateVmAttr(ctx, client, timeout, opts); err != nil {
			return err
		}
	}

	return nil
}

func startVM(ctx context.Context, client *osc.Client, timeOut time.Duration, id string) error {
	_, err := client.StartVms(ctx, osc.StartVmsRequest{
		VmIds: []string{id},
	}, options.WithRetryTimeout(timeOut))
	if err != nil {
		return fmt.Errorf("error starting vm %s", err)
	}

	stateConf := &retry.StateChangeConf{
		Pending: []string{"pending", "stopped"},
		Target:  []string{"running"},
		Timeout: timeOut,
		Refresh: vmStateRefreshFunc(ctx, client, id, "", timeOut),
	}

	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("error waiting for instance (%s) to become ready: %s", id, err)
	}

	return nil
}

func updateVmAttr(ctx context.Context, client *osc.Client, timeout time.Duration, instanceAttrOpts osc.UpdateVmRequest) error {
	_, err := client.UpdateVm(ctx, instanceAttrOpts, options.WithRetryTimeout(timeout))
	if err != nil {
		return err
	}
	return nil
}

func readVM(ctx context.Context, client *osc.Client, timeOut time.Duration, id string) (*osc.ReadVmsResponse, error) {
	resp, err := client.ReadVms(ctx, osc.ReadVmsRequest{
		Filters: &osc.FiltersVm{
			VmIds: &[]string{id},
		},
	}, options.WithRetryTimeout(timeOut))
	return resp, err
}

// AttributeSetter you can use this function to set the attributes
type AttributeSetter func(key string, value interface{}) error

func resourceDataAttrSetter(d *schema.ResourceData, callback func(AttributeSetter) error) error {
	setterFunc := func(key string, value interface{}) error {
		return d.Set(key, value)
	}
	return callback(setterFunc)
}
