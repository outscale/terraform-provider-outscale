package oapi

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/oapi-codegen/runtime/types"
	"github.com/outscale/goutils/sdk/ptr"
	"github.com/outscale/osc-sdk-go/v3/pkg/options"
	"github.com/outscale/osc-sdk-go/v3/pkg/osc"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/framework/fwhelpers/from"
	"github.com/outscale/terraform-provider-outscale/internal/services/oapi/oapihelpers"
	"github.com/outscale/terraform-provider-outscale/internal/utils"
	"github.com/samber/lo"
	"github.com/spf13/cast"
)

func DataSourceOutscaleVM() *schema.Resource {
	return &schema.Resource{
		ReadContext: DataSourceOutscaleVMRead,
		Schema:      getDataSourceOAPIVMSchemas(),
	}
}

func DataSourceOutscaleVMRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.OutscaleClient).OSC
	timeout := d.Timeout(schema.TimeoutRead)

	filters, filtersOk := d.GetOk("filter")
	instanceID, instanceIDOk := d.GetOk("vm_id")
	var err error
	if !filtersOk && !instanceIDOk {
		return diag.Errorf("one of filters, or instance_id must be assigned")
	}
	// Build up search parameters
	params := osc.ReadVmsRequest{}
	if filtersOk {
		params.Filters, err = buildOutscaleDataSourceVMFilters(filters.(*schema.Set))
		if err != nil {
			return diag.FromErr(err)
		}
	}
	if instanceIDOk {
		params.Filters.VmIds = &[]string{instanceID.(string)}
	}

	log.Printf("[DEBUG] ReadVmsRequest -> %+v\n", params)

	resp, err := client.ReadVms(ctx, params, options.WithRetryTimeout(timeout))
	if err != nil {
		return diag.Errorf("error reading the vm %s", err)
	}

	if resp.Vms == nil {
		return diag.FromErr(ErrNoResults)
	}

	var filteredVms []osc.Vm

	// loop through reservations, and remove terminated instances, populate vm slice
	for _, res := range ptr.From(resp.Vms) {
		if res.State != "terminated" {
			filteredVms = append(filteredVms, res)
		}
	}

	var vm osc.Vm
	if len(filteredVms) == 0 {
		return diag.FromErr(ErrNoResults)
	}

	if len(filteredVms) > 1 {
		return diag.FromErr(ErrMultipleResults)
	}

	vm = filteredVms[0]

	// Populate vm attribute fields with the returned vm
	return diag.FromErr(resourceDataAttrSetter(d, func(set AttributeSetter) error {
		d.SetId(vm.VmId)

		booTags, errTags := oapihelpers.GetBsuTagsMaps(ctx, client, timeout, vm)
		if errTags != nil {
			return errTags
		}
		if err := d.Set("block_device_mappings_created", getOscAPIVMBlockDeviceMapping(
			booTags, vm.BlockDeviceMappings)); err != nil {
			return err
		}

		return oapiVMDescriptionAttributes(set, &vm)
	}))
}

// Populate instance attribute fields with the returned instance
func oapiVMDescriptionAttributes(set AttributeSetter, vm *osc.Vm) error {
	if err := set("actions_on_next_boot", getActionsOnNextBoot(vm.ActionsOnNextBoot)); err != nil {
		return err
	}
	if err := set("boot_mode", vm.BootMode); err != nil {
		return err
	}
	if err := set("tpm_enabled", vm.TpmEnabled); err != nil {
		return err
	}
	if err := set("architecture", vm.Architecture); err != nil {
		return err
	}
	if err := set("bsu_optimized", vm.BsuOptimized); err != nil {
		return err
	}
	if err := set("client_token", vm.ClientToken); err != nil {
		return err
	}
	if err := set("creation_date", from.ISO8601(vm.CreationDate)); err != nil {
		return err
	}
	if err := set("deletion_protection", vm.DeletionProtection); err != nil {
		return err
	}
	if err := set("hypervisor", vm.Hypervisor); err != nil {
		return err
	}
	if err := set("image_id", vm.ImageId); err != nil {
		return err
	}
	if err := set("is_source_dest_checked", vm.IsSourceDestChecked); err != nil {
		return err
	}
	if err := set("keypair_name", vm.KeypairName); err != nil {
		return err
	}
	if err := set("launch_number", vm.LaunchNumber); err != nil {
		return err
	}
	if err := set("net_id", vm.NetId); err != nil {
		return err
	}
	if err := set("nested_virtualization", vm.NestedVirtualization); err != nil {
		return err
	}
	prNic, secNic := oapihelpers.GetOAPIVMNetworkInterfaceLightSet(vm.Nics)
	if err := set("primary_nic", prNic); err != nil {
		return err
	}
	if err := set("nics", secNic); err != nil {
		return err
	}
	if err := set("os_family", vm.OsFamily); err != nil {
		return err
	}
	if err := set("performance", vm.Performance); err != nil {
		return err
	}
	if err := set("placement_subregion_name", vm.Placement.SubregionName); err != nil {
		return err
	}
	if err := set("placement_tenancy", vm.Placement.Tenancy); err != nil {
		return err
	}
	if err := set("private_dns_name", vm.PrivateDnsName); err != nil {
		return err
	}
	if err := set("private_ip", vm.PrivateIp); err != nil {
		return err
	}
	if err := set("product_codes", vm.ProductCodes); err != nil {
		return err
	}
	if err := set("public_dns_name", vm.PublicDnsName); err != nil {
		return err
	}
	if err := set("public_ip", vm.PublicIp); err != nil {
		return err
	}
	if err := set("reservation_id", vm.ReservationId); err != nil {
		return err
	}
	if err := set("root_device_name", vm.RootDeviceName); err != nil {
		return err
	}
	if err := set("root_device_type", vm.RootDeviceType); err != nil {
		return err
	}
	if err := set("security_groups", getSecurityGroups(vm.SecurityGroups)); err != nil {
		return err
	}
	if err := set("state", vm.State); err != nil {
		return err
	}
	if err := set("state_reason", vm.StateReason); err != nil {
		return err
	}
	if err := set("subnet_id", vm.SubnetId); err != nil {
		return err
	}
	if err := set("user_data", vm.UserData); err != nil {
		return err
	}
	if err := set("vm_id", vm.VmId); err != nil {
		return err
	}
	if err := set("vm_initiated_shutdown_behavior", vm.VmInitiatedShutdownBehavior); err != nil {
		return err
	}
	if err := set("tags", FlattenOAPITagsSDK(vm.Tags)); err != nil {
		return err
	}
	return set("vm_type", vm.VmType)
}

func getOscAPIVMBlockDeviceMapping(busTagsMaps map[string]interface{}, blockDeviceMappings []osc.BlockDeviceMappingCreated) (blockDeviceMapping []map[string]interface{}) {
	for _, v := range blockDeviceMappings {
		blockDevice := map[string]interface{}{
			"device_name": v.DeviceName,
			"bsu":         getbusToSet(v.Bsu, busTagsMaps, v.DeviceName),
		}
		blockDeviceMapping = append(blockDeviceMapping, blockDevice)
	}
	return
}

func getbusToSet(bsu osc.BsuCreated, busTagsMaps map[string]interface{}, deviceName string) (res []map[string]interface{}) {
	res = append(res, map[string]interface{}{
		"delete_on_vm_deletion": bsu.DeleteOnVmDeletion,
		"volume_id":             bsu.VolumeId,
		"state":                 bsu.State,
		"link_date":             from.ISO8601(bsu.LinkDate),
		"tags":                  FlattenOAPITagsSDK(busTagsMaps[deviceName].([]osc.ResourceTag)),
	})
	return
}

func getSecurityGroups(groupSet []osc.SecurityGroupLight) []map[string]interface{} {
	res := []map[string]interface{}{}
	for _, g := range groupSet {
		r := map[string]interface{}{
			"security_group_id":   g.SecurityGroupId,
			"security_group_name": g.SecurityGroupName,
		}
		res = append(res, r)
	}

	return res
}

func getActionsOnNextBoot(actionsOnNextBoot osc.ActionsOnNextBoot) []map[string]interface{} {
	return []map[string]interface{}{
		{
			"secure_boot": new(osc.SecureBootAction(ptr.From(actionsOnNextBoot.SecureBoot))),
		},
	}
}

func getSecurityGroupIds(sgIds []osc.SecurityGroupLight) []string {
	res := make([]string, len(sgIds))
	for k, ids := range sgIds {
		res[k] = ids.SecurityGroupId
	}
	return res
}

func getDataSourceOAPIVMSchemas() map[string]*schema.Schema {
	wholeSchema := map[string]*schema.Schema{
		"filter": dataSourceFiltersSchema(),
	}

	attrsSchema := getOApiVMAttributesSchema()

	for k, v := range attrsSchema {
		wholeSchema[k] = v
	}

	wholeSchema["request_id"] = &schema.Schema{
		Type:     schema.TypeString,
		Computed: true,
	}

	return wholeSchema
}

func buildOutscaleDataSourceVMFilters(set *schema.Set) (*osc.FiltersVm, error) {
	filters := new(osc.FiltersVm)

	for _, v := range set.List() {
		m := v.(map[string]interface{})
		filterValues := make([]string, 0)
		for _, e := range m["values"].([]interface{}) {
			filterValues = append(filterValues, e.(string))
		}

		switch name := m["name"].(string); name {
		case "architectures":
			filters.Architectures = &filterValues
		case "Block_device_mapping_delete_on_vm_deletion":
			filters.BlockDeviceMappingDeleteOnVmDeletion = new(cast.ToBool(filterValues[0]))
		case "block_device_mapping_device_names":
			filters.BlockDeviceMappingDeviceNames = &filterValues
		case "block_device_mapping_states":
			filters.BlockDeviceMappingStates = &filterValues
		case "block_device_mapping_link_dates":
			linkDates, err := utils.StringSliceToTimeSlice(
				filterValues, "block_device_mapping_link_dates")
			if err != nil {
				return filters, err
			}
			filters.BlockDeviceMappingLinkDates = new(lo.Map(linkDates, func(t time.Time, _ int) types.Date { return types.Date{Time: t} }))
		case "block_device_mapping_volume_ids":
			filters.BlockDeviceMappingVolumeIds = &filterValues
		case "boot_modes":
			filters.BootModes = new(lo.Map(filterValues, func(s string, _ int) osc.BootMode { return (osc.BootMode)(s) }))
		case "ClientTokens":
			filters.ClientTokens = &filterValues
		case "creation_dates":
			creationDates, err := utils.StringSliceToTimeSlice(
				filterValues, "creation_dates")
			if err != nil {
				return filters, err
			}
			filters.CreationDates = new(lo.Map(creationDates, func(t time.Time, _ int) types.Date { return types.Date{Time: t} }))
		case "image_ids":
			filters.ImageIds = &filterValues
		case "is_source_dest_checked":
			filters.IsSourceDestChecked = new(cast.ToBool(filterValues[0]))
		case "keypair_names":
			filters.KeypairNames = &filterValues
		case "launch_numbers":
			filters.LaunchNumbers = new(utils.StringSliceToIntSlice(filterValues))
		case "lifecycles":
			filters.Lifecycles = &filterValues
		case "net_ids":
			filters.NetIds = &filterValues
		case "nic_account_ids":
			filters.NicAccountIds = &filterValues
		case "nic_descriptions":
			filters.NicDescriptions = &filterValues
		case "nic_is_source_dest_checked":
			filters.NicIsSourceDestChecked = new(cast.ToBool(filterValues[0]))
		case "nic_link_nic_delete_on_vm_deletion":
			filters.NicLinkNicDeleteOnVmDeletion = new(cast.ToBool(filterValues[0]))
		case "nic_link_nic_device_numbers":
			filters.NicLinkNicDeviceNumbers = new(utils.StringSliceToIntSlice(filterValues))
		case "nic_link_nic_link_nic_dates":
			linkDates, err := utils.StringSliceToTimeSlice(
				filterValues, "nic_link_nic_link_nic_dates")
			if err != nil {
				return filters, err
			}
			filters.NicLinkNicLinkNicDates = new(lo.Map(linkDates, func(t time.Time, _ int) types.Date { return types.Date{Time: t} }))
		case "nic_link_nic_link_nic_ids":
			filters.NicLinkNicLinkNicIds = &filterValues
		case "nic_link_nic_states":
			filters.NicLinkNicStates = &filterValues
		case "nic_link_nic_vm_account_ids":
			filters.NicLinkNicVmAccountIds = &filterValues
		case "nic_link_nic_vm_ids":
			filters.NicLinkNicVmIds = &filterValues
		case "nic_link_public_ip_account_ids":
			filters.NicLinkPublicIpAccountIds = &filterValues
		case "nic_link_public_ip_link_public_ip_ids":
			filters.NicLinkPublicIpLinkPublicIpIds = &filterValues
		case "nic_link_public_ip_public_ip_ids":
			filters.NicLinkPublicIpPublicIpIds = &filterValues
		case "nic_link_public_Ip_public_ips":
			filters.NicLinkPublicIpPublicIps = &filterValues
		case "nic_mac_addresses":
			filters.NicMacAddresses = &filterValues
		case "nic_net_ids":
			filters.NicNetIds = &filterValues
		case "nic_nic_ids":
			filters.NicNicIds = &filterValues
		case "nic_private_ips_link_public_ip_account_ids":
			filters.NicPrivateIpsLinkPublicIpAccountIds = &filterValues
		case "nic_private_ips_primary_ip":
			filters.NicPrivateIpsPrimaryIp = new(cast.ToBool(filterValues[0]))
		case "nic_private_ips_private_ips":
			filters.NicPrivateIpsPrivateIps = &filterValues
		case "nic_security_group_ids":
			filters.NicSecurityGroupIds = &filterValues
		case "nic_security_group_names":
			filters.NicSecurityGroupNames = &filterValues
		case "nic_states":
			filters.NicStates = &filterValues
		case "nic_subnet_ids":
			filters.NicSubnetIds = &filterValues
		case "nic_subregion_names":
			filters.NicSubregionNames = &filterValues
		case "platforms":
			filters.Platforms = &filterValues
		case "private_ips":
			filters.PrivateIps = &filterValues
		case "product_codes":
			filters.ProductCodes = &filterValues
		case "public_ips":
			filters.PublicIps = &filterValues
		case "reservation_ids":
			filters.ReservationIds = &filterValues
		case "root_device_names":
			filters.RootDeviceNames = &filterValues
		case "root_tevice_types":
			filters.RootDeviceTypes = &filterValues
		case "security_group_ids":
			filters.SecurityGroupIds = &filterValues
		case "security_group_names":
			filters.SecurityGroupNames = &filterValues
		case "state_reason_codes":
			filters.StateReasonCodes = new(utils.StringSliceToIntSlice(filterValues))
		case "state_reason_messages":
			filters.StateReasonMessages = &filterValues
		case "state_reasons":
			filters.StateReasons = &filterValues
		case "subnet_ids":
			filters.SubnetIds = &filterValues
		case "subregion_names":
			filters.SubregionNames = &filterValues
		case "tag_keys":
			filters.TagKeys = &filterValues
		case "tag_values":
			filters.TagValues = &filterValues
		case "tags":
			filters.Tags = &filterValues
		case "vm_ids":
			filters.VmIds = &filterValues
		case "tenancies":
			filters.Tenancies = &filterValues
		case "vm_security_group_ids":
			filters.VmSecurityGroupIds = &filterValues
		case "vm_security_group_names":
			filters.VmSecurityGroupNames = &filterValues
		case "vm_state_codes":
			filters.VmStateCodes = new(utils.StringSliceToIntSlice(filterValues))
		case "vm_state_names":
			filters.VmStateNames = new(lo.Map(filterValues, func(s string, _ int) osc.VmState { return osc.VmState(s) }))
		case "VmTypes":
			filters.VmTypes = &filterValues
		default:
			return nil, utils.UnknownDataSourceFilterError(name)
		}
	}
	return filters, nil
}

func getOApiVMAttributesSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		// Attributes
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
						Type:     schema.TypeList,
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
		"boot_mode": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"tpm_enabled": {
			Type:     schema.TypeBool,
			Computed: true,
		},
		"bsu_optimized": {
			Type:     schema.TypeBool,
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
			Computed: true,
		},
		"hypervisor": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"image_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"is_source_dest_checked": {
			Type:     schema.TypeBool,
			Computed: true,
		},
		"keypair_name": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"security_group_ids": {
			Type:     schema.TypeList,
			Computed: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
		"security_group_names": {
			Type:     schema.TypeList,
			Computed: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
		"launch_number": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"nested_virtualization": {
			Type:     schema.TypeBool,
			Computed: true,
		},
		"net_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"primary_nic": {
			Type:     schema.TypeSet,
			Computed: true,
			Set: func(v interface{}) int {
				return v.(map[string]interface{})["device_number"].(int)
			},
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"delete_on_vm_deletion": {
						Type:     schema.TypeBool,
						Computed: true,
					},
					"description": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"device_number": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"nic_id": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"private_ips": {
						Type:     schema.TypeSet,
						Computed: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"is_primary": {
									Type:     schema.TypeBool,
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
									Computed: true,
								},
							},
						},
					},
					"secondary_private_ip_count": {
						Type:     schema.TypeInt,
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
		"nics": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"delete_on_vm_deletion": {
						Type:     schema.TypeBool,
						Computed: true,
					},
					"description": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"device_number": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"nic_id": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"private_ips": {
						Type:     schema.TypeSet,
						Computed: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"is_primary": {
									Type:     schema.TypeBool,
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
									Computed: true,
								},
							},
						},
					},
					"secondary_private_ip_count": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"security_group_ids": {
						Type:     schema.TypeList,
						Computed: true,
						Elem:     &schema.Schema{Type: schema.TypeString},
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
					"security_groups_names": {
						Type:     schema.TypeList,
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
		"os_family": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"performance": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"placement_subregion_name": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"placement_tenancy": {
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
		"state_reason": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"subnet_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"tags": {
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
		"user_data": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"vm_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"vm_initiated_shutdown_behavior": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"vm_type": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"request_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"private_ips": {
			Type:     schema.TypeList,
			Computed: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
	}
}
