package outscale

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/outscale/terraform-provider-outscale/utils"
	"github.com/spf13/cast"
)

func DataSourceOutscaleVM() *schema.Resource {
	return &schema.Resource{
		Read:   DataSourceOutscaleVMRead,
		Schema: getDataSourceOAPIVMSchemas(),
	}
}
func DataSourceOutscaleVMRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*OutscaleClient).OSCAPI

	filters, filtersOk := d.GetOk("filter")
	instanceID, instanceIDOk := d.GetOk("vm_id")
	var err error
	if !filtersOk && !instanceIDOk {
		return fmt.Errorf("one of filters, or instance_id must be assigned")
	}
	// Build up search parameters
	params := oscgo.ReadVmsRequest{}
	if filtersOk {
		params.Filters, err = buildOutscaleDataSourceVMFilters(filters.(*schema.Set))
		if err != nil {
			return err
		}
	}
	if instanceIDOk {
		params.Filters.VmIds = &[]string{instanceID.(string)}
	}

	log.Printf("[DEBUG] ReadVmsRequest -> %+v\n", params)

	var resp oscgo.ReadVmsResponse
	err = resource.Retry(30*time.Second, func() *resource.RetryError {
		rp, httpResp, err := client.VmApi.ReadVms(context.Background()).ReadVmsRequest(params).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		resp = rp
		return nil
	})

	if err != nil {
		return fmt.Errorf("error reading the VM %s", err)
	}

	if !resp.HasVms() {
		return fmt.Errorf("Your query returned no results. Please change your search criteria and try again")
	}

	var filteredVms []oscgo.Vm

	// loop through reservations, and remove terminated instances, populate vm slice
	for _, res := range resp.GetVms() {
		if res.GetState() != "terminated" {
			filteredVms = append(filteredVms, res)
		}
	}

	var vm oscgo.Vm
	if len(filteredVms) < 1 {
		return errors.New("Your query returned no results. Please change your search criteria and try again")
	}

	if len(filteredVms) > 1 {
		return errors.New("Your query returned more than one result. Please try a more " +
			"specific search criteria")
	}

	vm = filteredVms[0]

	// Populate vm attribute fields with the returned vm
	return resourceDataAttrSetter(d, func(set AttributeSetter) error {
		d.SetId(vm.GetVmId())

		booTags, errTags := utils.GetBsuTagsMaps(vm, client)
		if errTags != nil {
			return errTags
		}
		if err := d.Set("block_device_mappings_created", getOscAPIVMBlockDeviceMapping(
			booTags, vm.GetBlockDeviceMappings())); err != nil {
			return err
		}

		return oapiVMDescriptionAttributes(set, &vm)
	})
}

// Populate instance attribute fields with the returned instance
func oapiVMDescriptionAttributes(set AttributeSetter, vm *oscgo.Vm) error {
	if err := set("actions_on_next_boot", getActionsOnNextBoot(vm.GetActionsOnNextBoot())); err != nil {
		return err
	}
	if err := set("boot_mode", vm.GetBootMode()); err != nil {
		return err
	}
	if err := set("architecture", vm.GetArchitecture()); err != nil {
		return err
	}
	if err := set("bsu_optimized", vm.GetBsuOptimized()); err != nil {
		return err
	}
	if err := set("client_token", vm.GetClientToken()); err != nil {
		return err
	}
	if err := set("creation_date", vm.GetCreationDate()); err != nil {
		return err
	}
	if err := set("deletion_protection", vm.GetDeletionProtection()); err != nil {
		return err
	}
	if err := set("hypervisor", vm.GetHypervisor()); err != nil {
		return err
	}
	if err := set("image_id", vm.GetImageId()); err != nil {
		return err
	}
	if err := set("is_source_dest_checked", vm.GetIsSourceDestChecked()); err != nil {
		return err
	}
	if err := set("keypair_name", vm.GetKeypairName()); err != nil {
		return err
	}
	if err := set("launch_number", vm.GetLaunchNumber()); err != nil {
		return err
	}
	if err := set("net_id", vm.GetNetId()); err != nil {
		return err
	}
	if err := set("nested_virtualization", vm.GetNestedVirtualization()); err != nil {
		return err
	}
	prNic, secNic := getOAPIVMNetworkInterfaceLightSet(vm.GetNics())
	if err := set("primary_nic", prNic); err != nil {
		return err
	}
	if err := set("nics", secNic); err != nil {
		return err
	}
	if err := set("os_family", vm.GetOsFamily()); err != nil {
		return err
	}
	if err := set("performance", vm.GetPerformance()); err != nil {
		return err
	}
	if err := set("placement_subregion_name", aws.StringValue(vm.GetPlacement().SubregionName)); err != nil {
		return err
	}
	if err := set("placement_tenancy", aws.StringValue(vm.GetPlacement().Tenancy)); err != nil {
		return err
	}
	if err := set("private_dns_name", vm.GetPrivateDnsName()); err != nil {
		return err
	}
	if err := set("private_ip", vm.GetPrivateIp()); err != nil {
		return err
	}
	if err := set("product_codes", vm.GetProductCodes()); err != nil {
		return err
	}
	if err := set("public_dns_name", vm.GetPublicDnsName()); err != nil {
		return err
	}
	if err := set("public_ip", vm.GetPublicIp()); err != nil {
		return err
	}
	if err := set("reservation_id", vm.GetReservationId()); err != nil {
		return err
	}
	if err := set("root_device_name", vm.GetRootDeviceName()); err != nil {
		return err
	}
	if err := set("root_device_type", vm.GetRootDeviceType()); err != nil {
		return err
	}
	if err := set("security_groups", getSecurityGroups(vm.GetSecurityGroups())); err != nil {
		return err
	}
	if err := set("state", vm.GetState()); err != nil {
		return err
	}
	if err := set("state_reason", vm.GetStateReason()); err != nil {
		return err
	}
	if err := set("subnet_id", vm.GetSubnetId()); err != nil {
		return err
	}
	if err := set("user_data", vm.GetUserData()); err != nil {
		return err
	}
	if err := set("vm_id", vm.GetVmId()); err != nil {
		return err
	}
	if err := set("vm_initiated_shutdown_behavior", vm.GetVmInitiatedShutdownBehavior()); err != nil {
		return err
	}
	if err := set("tags", flattenOAPITagsSDK(vm.GetTags())); err != nil {
		return err
	}
	return set("vm_type", vm.GetVmType())
}

func getOscAPIVMBlockDeviceMapping(busTagsMaps map[string]interface{}, blockDeviceMappings []oscgo.BlockDeviceMappingCreated) (blockDeviceMapping []map[string]interface{}) {
	for _, v := range blockDeviceMappings {
		blockDevice := map[string]interface{}{
			"device_name": v.GetDeviceName(),
			"bsu":         getbusToSet(v.GetBsu(), busTagsMaps, *v.DeviceName),
		}
		blockDeviceMapping = append(blockDeviceMapping, blockDevice)
	}
	return
}

func getbusToSet(bsu oscgo.BsuCreated, busTagsMaps map[string]interface{}, deviceName string) (res []map[string]interface{}) {
	res = append(res, map[string]interface{}{
		"delete_on_vm_deletion": bsu.GetDeleteOnVmDeletion(),
		"volume_id":             bsu.GetVolumeId(),
		"state":                 bsu.GetState(),
		"link_date":             bsu.GetLinkDate(),
		"tags":                  flattenOAPITagsSDK(busTagsMaps[deviceName].([]oscgo.ResourceTag)),
	})
	return
}

func getSecurityGroups(groupSet []oscgo.SecurityGroupLight) []map[string]interface{} {
	res := []map[string]interface{}{}
	for _, g := range groupSet {
		r := map[string]interface{}{
			"security_group_id":   g.GetSecurityGroupId(),
			"security_group_name": g.GetSecurityGroupName(),
		}
		res = append(res, r)
	}

	return res
}

func getActionsOnNextBoot(actionsOnNextBoot oscgo.ActionsOnNextBoot) []map[string]interface{} {
	return []map[string]interface{}{
		{
			"secure_boot": string(actionsOnNextBoot.GetSecureBoot()),
		},
	}
}

func getSecurityGroupIds(sgIds []oscgo.SecurityGroupLight) []string {
	res := make([]string, len(sgIds))
	for k, ids := range sgIds {
		res[k] = ids.GetSecurityGroupId()
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

func buildOutscaleDataSourceVMFilters(set *schema.Set) (*oscgo.FiltersVm, error) {
	filters := new(oscgo.FiltersVm)

	for _, v := range set.List() {
		m := v.(map[string]interface{})
		filterValues := make([]string, 0)
		for _, e := range m["values"].([]interface{}) {
			filterValues = append(filterValues, e.(string))
		}

		switch name := m["name"].(string); name {
		case "architectures":
			filters.SetArchitectures(filterValues)
		case "Block_device_mapping_delete_on_vm_deletion":
			filters.SetBlockDeviceMappingDeleteOnVmDeletion(cast.ToBool(filterValues[0]))
		case "block_device_mapping_device_names":
			filters.SetBlockDeviceMappingDeviceNames(filterValues)
		case "block_device_mapping_states":
			filters.SetBlockDeviceMappingStates(filterValues)
		case "block_device_mapping_link_dates":
			linkDates, err := utils.FiltersTimesToStringSlice(
				filterValues, "block_device_mapping_link_dates")
			if err != nil {
				return filters, err
			}
			filters.SetBlockDeviceMappingLinkDates(linkDates)
		case "block_device_mapping_volume_ids":
			filters.SetBlockDeviceMappingVolumeIds(filterValues)
		case "boot_modes":
			filters.SetBootModes(utils.Map(filterValues, func(s string) oscgo.BootMode { return (oscgo.BootMode)(s) }))
		case "ClientTokens":
			filters.SetClientTokens(filterValues)
		case "creation_dates":
			creationDates, err := utils.FiltersTimesToStringSlice(
				filterValues, "creation_dates")
			if err != nil {
				return filters, err
			}
			filters.SetCreationDates(creationDates)
		case "image_ids":
			filters.SetImageIds(filterValues)
		case "is_source_dest_checked":
			filters.SetIsSourceDestChecked(cast.ToBool(filterValues[0]))
		case "keypair_names":
			filters.SetKeypairNames(filterValues)
		case "launch_numbers":
			filters.SetLaunchNumbers(utils.StringSliceToInt32Slice(filterValues))
		case "lifecycles":
			filters.SetLifecycles(filterValues)
		case "net_ids":
			filters.SetNetIds(filterValues)
		case "nic_account_ids":
			filters.SetNicAccountIds(filterValues)
		case "nic_descriptions":
			filters.SetNicDescriptions(filterValues)
		case "nic_is_source_dest_checked":
			filters.SetNicIsSourceDestChecked(cast.ToBool(filterValues[0]))
		case "nic_link_nic_delete_on_vm_deletion":
			filters.SetNicLinkNicDeleteOnVmDeletion(cast.ToBool(filterValues[0]))
		case "nic_link_nic_device_numbers":
			filters.SetNicLinkNicDeviceNumbers(
				utils.StringSliceToInt32Slice(filterValues))
		case "nic_link_nic_link_nic_dates":
			linkDates, err := utils.FiltersTimesToStringSlice(
				filterValues, "nic_link_nic_link_nic_dates")
			if err != nil {
				return filters, err
			}
			filters.SetNicLinkNicLinkNicDates(linkDates)
		case "nic_link_nic_link_nic_ids":
			filters.SetNicLinkNicLinkNicIds(filterValues)
		case "nic_link_nic_states":
			filters.SetNicLinkNicStates(filterValues)
		case "nic_link_nic_vm_account_ids":
			filters.SetNicLinkNicVmAccountIds(filterValues)
		case "nic_link_nic_vm_ids":
			filters.SetNicLinkNicVmIds(filterValues)
		case "nic_link_public_ip_account_ids":
			filters.SetNicLinkPublicIpAccountIds(filterValues)
		case "nic_link_public_ip_link_public_ip_ids":
			filters.SetNicLinkPublicIpLinkPublicIpIds(filterValues)
		case "nic_link_public_ip_public_ip_ids":
			filters.SetNicLinkPublicIpPublicIpIds(filterValues)
		case "nic_link_public_Ip_public_ips":
			filters.SetNicLinkPublicIpPublicIps(filterValues)
		case "nic_mac_addresses":
			filters.SetNicMacAddresses(filterValues)
		case "nic_net_ids":
			filters.SetNicNetIds(filterValues)
		case "nic_nic_ids":
			filters.SetNicNicIds(filterValues)
		case "nic_private_ips_link_public_ip_account_ids":
			filters.SetNicPrivateIpsLinkPublicIpAccountIds(filterValues)
		case "nic_private_ips_primary_ip":
			filters.SetNicPrivateIpsPrimaryIp(cast.ToBool(filterValues[0]))
		case "nic_private_ips_private_ips":
			filters.SetNicPrivateIpsPrivateIps(filterValues)
		case "nic_security_group_ids":
			filters.SetNicSecurityGroupIds(filterValues)
		case "nic_security_group_names":
			filters.SetNicSecurityGroupNames(filterValues)
		case "nic_states":
			filters.SetNicStates(filterValues)
		case "nic_subnet_ids":
			filters.SetNicSubnetIds(filterValues)
		case "nic_subregion_names":
			filters.SetNicSubregionNames(filterValues)
		case "platforms":
			filters.SetPlatforms(filterValues)
		case "private_ips":
			filters.SetPrivateIps(filterValues)
		case "product_codes":
			filters.SetProductCodes(filterValues)
		case "public_ips":
			filters.SetPublicIps(filterValues)
		case "reservation_ids":
			filters.SetReservationIds(filterValues)
		case "root_device_names":
			filters.SetRootDeviceNames(filterValues)
		case "root_tevice_types":
			filters.SetRootDeviceTypes(filterValues)
		case "security_group_ids":
			filters.SetSecurityGroupIds(filterValues)
		case "security_group_names":
			filters.SetSecurityGroupNames(filterValues)
		case "state_reason_codes":
			filters.SetStateReasonCodes(
				utils.StringSliceToInt32Slice(filterValues))
		case "state_reason_messages":
			filters.SetStateReasonMessages(filterValues)
		case "state_reasons":
			filters.SetStateReasons(filterValues)
		case "subnet_ids":
			filters.SetSubnetIds(filterValues)
		case "subregion_names":
			filters.SetSubregionNames(filterValues)
		case "tag_keys":
			filters.SetTagKeys(filterValues)
		case "tag_values":
			filters.SetTagValues(filterValues)
		case "tags":
			filters.SetTags(filterValues)
		case "vm_ids":
			filters.SetVmIds(filterValues)
		case "tenancies":
			filters.SetTenancies(filterValues)
		case "vm_security_group_ids":
			filters.SetVmSecurityGroupIds(filterValues)
		case "vm_security_group_names":
			filters.SetVmSecurityGroupNames(filterValues)
		case "vm_state_codes":
			filters.SetVmStateCodes(
				utils.StringSliceToInt32Slice(filterValues))
		case "vm_state_names":
			filters.SetVmStateNames(filterValues)
		case "VmTypes":
			filters.SetVmTypes(filterValues)
		default:
			return nil, utils.UnknownDataSourceFilterError(context.Background(), name)
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
