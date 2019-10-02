package outscale

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-outscale/osc/oapi"
)

func dataSourceOutscaleOAPIVM() *schema.Resource {
	return &schema.Resource{
		Read:   dataSourceOutscaleOAPIVMRead,
		Schema: getDataSourceOAPIVMSchemas(),
	}
}
func dataSourceOutscaleOAPIVMRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*OutscaleClient).OAPI

	filters, filtersOk := d.GetOk("filter")
	instanceID, instanceIDOk := d.GetOk("vm_id")

	if filtersOk == false && instanceIDOk == false {
		return fmt.Errorf("One of filters, or instance_id must be assigned")
	}

	// Build up search parameters
	params := oapi.ReadVmsRequest{}
	if filtersOk {
		params.Filters = buildOutscaleOAPIDataSourceVmFilters(filters.(*schema.Set))
	}
	if instanceIDOk {
		params.Filters.VmIds = []string{instanceID.(string)}
	}

	fmt.Printf("ReadVmsRequest -> %+v\n", params)

	var resp *oapi.POST_ReadVmsResponses
	var err error

	err = resource.Retry(30*time.Second, func() *resource.RetryError {
		resp, err = client.POST_ReadVms(params)
		return resource.RetryableError(err)
	})

	if err != nil {
		return fmt.Errorf("Error reading the VM %s", err)
	}

	if resp.OK.Vms == nil {
		return fmt.Errorf("Your query returned no results. Please change your search criteria and try again")
	}

	// If no instances were returned, return
	if len(resp.OK.Vms) == 0 {
		return fmt.Errorf("Your query returned no results. Please change your search criteria and try again")
	}

	var filteredInstances []oapi.Vm

	// loop through reservations, and remove terminated instances, populate instance slice
	for _, res := range resp.OK.Vms {
		if res.State != "terminated" {
			filteredInstances = append(filteredInstances, res)
		}
	}

	var instance oapi.Vm
	if len(filteredInstances) < 1 {
		return errors.New("Your query returned no results. Please change your search criteria and try again")
	}

	// (TODO: Support a list of instances to be returned)
	// Possibly with a different data source that returns a list of individual instance data sources
	if len(filteredInstances) > 1 {
		return errors.New("Your query returned more than one result. Please try a more " +
			"specific search criteria")
	}

	instance = filteredInstances[0]

	log.Printf("[DEBUG] outscale_vm - Single VM ID found: %s", instance.VmId)
	d.Set("request_id", resp.OK.ResponseContext.RequestId)

	// Populate instance attribute fields with the returned instance
	return resourceDataAttrSetter(d, func(set AttributeSetter) error {
		d.SetId(instance.VmId)
		return oapiVMDescriptionAttributes(set, &instance)
	})
}

// Populate instance attribute fields with the returned instance
func oapiVMDescriptionAttributes(set AttributeSetter, instance *oapi.Vm) error {

	set("architecture", instance.Architecture)
	if err := set("block_device_mappings", getOAPIVMBlockDeviceMapping(instance.BlockDeviceMappings)); err != nil {
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
	set("tags", getOapiTagSet(instance.Tags))
	set("user_data", instance.UserData)
	set("vm_id", instance.VmId)
	set("vm_initiated_shutdown_behavior", instance.VmInitiatedShutdownBehavior)

	return set("vm_type", instance.VmType)
}

func getOAPIVMBlockDeviceMapping(b []oapi.BlockDeviceMappingCreated) (blockDeviceMapping []map[string]interface{}) {
	for i := 0; i < len(b); i++ {
		blockDeviceMapping = append(blockDeviceMapping, map[string]interface{}{
			"device_name": b[i].DeviceName,
			"bsu": map[string]interface{}{
				"delete_on_vm_deletion": fmt.Sprintf("%t", b[i].Bsu.DeleteOnVmDeletion),
				"volume_id":             b[i].Bsu.VolumeId,
				"state":                 b[i].Bsu.State,
				"link_date":             b[i].Bsu.LinkDate,
			},
		})
	}
	return
}

func getOAPIVMSecurityGroups(groupSet []oapi.SecurityGroupLight) []map[string]interface{} {
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

func buildOutscaleOAPIDataSourceVmFilters(set *schema.Set) oapi.FiltersVm {
	var filters oapi.FiltersVm
	for _, v := range set.List() {
		m := v.(map[string]interface{})
		var filterValues []string
		for _, e := range m["values"].([]interface{}) {
			filterValues = append(filterValues, e.(string))
		}

		switch name := m["name"].(string); name {
		case "account_ids":
			filters.AccountIds = filterValues
		case "activated_check":
			filters.ActivatedCheck, _ = strconv.ParseBool(filterValues[0])
		case "architectures":
			filters.Architectures = filterValues
		case "block_device_mapping_delete_on_vm_deletion":
			filters.BlockDeviceMappingDeleteOnVmDeletion, _ = strconv.ParseBool(filterValues[0])
		case "block_device_mapping_device_names":
			filters.BlockDeviceMappingDeviceNames = filterValues
		case "block_device_mapping_link_dates":
			filters.BlockDeviceMappingLinkDates = filterValues
		case "block_device_mapping_states":
			filters.BlockDeviceMappingStates = filterValues
		case "block_device_mapping_volume_ids":
			filters.BlockDeviceMappingVolumeIds = filterValues
		case "comments":
			filters.Comments = filterValues
		case "creation_dates":
			filters.CreationDates = filterValues
		case "dns_names":
			filters.DnsNames = filterValues
		case "hypervisors":
			filters.Hypervisors = filterValues
		case "image_ids":
			filters.ImageIds = filterValues
		case "kernel_ids":
			filters.KernelIds = filterValues
		case "keypair_names":
			filters.KeypairNames = filterValues
		case "launch_sort_numbers":
			filters.LaunchSortNumbers, _ = sliceAtoi(filterValues)
		case "link_nic_delete_on_vm_deletion":
			filters.LinkNicDeleteOnVmDeletion, _ = strconv.ParseBool(filterValues[0])
		case "link_nic_link_dates":
			filters.LinkNicLinkDates = filterValues
		case "link_nic_link_nic_ids":
			filters.LinkNicLinkNicIds = filterValues
		case "link_nic_link_public_ip_ids":
			filters.LinkNicLinkPublicIpIds = filterValues
		case "link_nic_nic_ids":
			filters.LinkNicNicIds = filterValues
		case "link_nic_nic_sort_numbers":
			filters.LinkNicNicSortNumbers, _ = sliceAtoi(filterValues)
		case "link_nic_public_ip_account_ids":
			filters.LinkNicPublicIpAccountIds = filterValues
		case "link_nic_public_ip_ids":
			filters.LinkNicPublicIpIds = filterValues
		case "link_nic_public_ips":
			filters.LinkNicPublicIps = filterValues
		case "link_nic_states":
			filters.LinkNicStates = filterValues
		case "link_nic_vm_account_ids":
			filters.LinkNicVmAccountIds = filterValues
		case "link_nic_vm_ids":
			filters.LinkNicVmIds = filterValues
		case "monitoring_states":
			filters.MonitoringStates = filterValues
		case "net_ids":
			filters.NetIds = filterValues
		case "nic_account_ids":
			filters.NicAccountIds = filterValues
		case "nic_activated_check":
			filters.NicActivatedCheck, _ = strconv.ParseBool(filterValues[0])
		case "nic_descriptions":
			filters.NicDescriptions = filterValues
		case "nic_mac_addresses":
			filters.NicMacAddresses = filterValues
		case "nic_net_ids":
			filters.NicNetIds = filterValues
		case "nic_nic_ids":
			filters.NicNicIds = filterValues
		case "nic_private_dns_names":
			filters.NicPrivateDnsNames = filterValues
		case "nic_security_group_ids":
			filters.NicSecurityGroupIds = filterValues
		case "nic_security_group_names":
			filters.NicSecurityGroupNames = filterValues
		case "nic_states":
			filters.NicStates = filterValues
		case "nic_subnet_ids":
			filters.NicSubnetIds = filterValues
		case "nic_subregion_names":
			filters.NicSubregionNames = filterValues
		case "placement_groups":
			filters.PlacementGroups = filterValues
		case "private_dns_names":
			filters.PrivateDnsNames = filterValues
		case "private_ip_link_private_ip_account_ids":
			filters.PrivateIpLinkPrivateIpAccountIds = filterValues
		case "private_ip_link_public_ips":
			filters.PrivateIpLinkPublicIps = filterValues
		case "private_ip_primary_ips":
			filters.PrivateIpPrimaryIps = filterValues
		case "private_ip_private_ips":
			filters.PrivateIpPrivateIps = filterValues
		case "private_ips":
			filters.PrivateIps = filterValues
		case "product_codes":
			filters.ProductCodes = filterValues
		case "public_ips":
			filters.PublicIps = filterValues
		case "ram_disk_ids":
			filters.RamDiskIds = filterValues
		case "root_device_names":
			filters.RootDeviceNames = filterValues
		case "root_device_types":
			filters.RootDeviceTypes = filterValues
		case "security_group_ids":
			filters.SecurityGroupIds = filterValues
		case "security_group_names":
			filters.SecurityGroupNames = filterValues
		case "spot_vm_request_ids":
			filters.SpotVmRequestIds = filterValues
		case "spot_vms":
			filters.SpotVms = filterValues
		case "state_comments":
			filters.StateComments = filterValues
		case "subnet_ids":
			filters.SubnetIds = filterValues
		case "subregion_names":
			filters.SubregionNames = filterValues
		case "systems":
			filters.Systems = filterValues
		case "tag_keys":
			filters.TagKeys = filterValues
		case "tag_values":
			filters.TagValues = filterValues
		case "tags":
			filters.Tags = filterValues
		case "tenancies":
			filters.Tenancies = filterValues
		case "tokens":
			filters.Tokens = filterValues
		case "virtualization_types":
			filters.VirtualizationTypes = filterValues
		case "vm_ids":
			filters.VmIds = filterValues
		case "vm_states":
			filters.VmStates = filterValues
		case "vm_types":
			filters.VmTypes = filterValues
		case "vms_security_group_ids":
			filters.VmsSecurityGroupIds = filterValues
		case "vms_security_group_names":
			filters.VmsSecurityGroupNames = filterValues

		default:
			log.Printf("[Debug] Unknown Filter Name: %s.", name)
		}
	}
	return filters
}

func sliceAtoi(sa []string) ([]int64, error) {
	si := make([]int64, 0, len(sa))
	for _, a := range sa {
		i, err := strconv.Atoi(a)
		if err != nil {
			return si, err
		}
		si = append(si, int64(i))
	}
	return si, nil
}

func getOapiTagSet(tags []oapi.ResourceTag) []map[string]interface{} {
	res := []map[string]interface{}{}

	if tags != nil {
		for _, t := range tags {
			tag := map[string]interface{}{}

			tag["key"] = t.Key
			tag["value"] = t.Value

			res = append(res, tag)
		}
	}

	return res
}

func getOApiVMAttributesSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		// Attributes
		"architecture": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"block_device_mappings": {
			Type:     schema.TypeList,
			Optional: true,
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
					"device_name": {
						Type:     schema.TypeString,
						Computed: true,
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
		"hypervisor": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"image_id": {
			Type:     schema.TypeString,
			ForceNew: true,
			Optional: true,
			Computed: true,
		},
		"is_source_dest_checked": {
			Type:     schema.TypeBool,
			Optional: true,
			Computed: true,
		},
		"keypair_name": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
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
		"launch_number": {
			Type:     schema.TypeInt,
			Computed: true,
		},

		"net_id": {
			Type:     schema.TypeString,
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
					"security_group_ids": {
						Type:     schema.TypeList,
						Optional: true,
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
		"os_family": {
			Type:     schema.TypeString,
			Computed: true,
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
			ForceNew: true,
			Optional: true,
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
		"private_ips": {
			Type:     schema.TypeList,
			Optional: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
	}
}
