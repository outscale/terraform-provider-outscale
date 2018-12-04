package outscale

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-outscale/osc/fcu"
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

	var resp *oapi.POST_ReadVmsResponses
	var err error

	err = resource.Retry(30*time.Second, func() *resource.RetryError {
		resp, err = client.POST_ReadVms(params)
		return resource.RetryableError(err)
	})

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

	return oapiVMDescriptionAttributes(d, &instance, client)
}

// Populate instance attribute fields with the returned instance
func oapiVMDescriptionAttributes(d *schema.ResourceData, instance *oapi.Vm, conn *oapi.Client) error {
	d.SetId(instance.VmId)
	// Set the easy attributes

	d.Set("launch_sort_number", instance.LaunchNumber)
	d.Set("architecture", instance.Architecture)
	// d.Set("block_device_mapping", getOAPIVMBlockDeviceMapping(instance.BlockDeviceMappings))
	d.Set("token", instance.ClientToken)
	d.Set("public_dns_name", instance.PublicDnsName)
	d.Set("bsu_optimised", instance.BsuOptimized)

	// TODO: add to struct for OAPI
	// d.Set("firewall_rules_set", instance.FirewallRulesSets)

	d.Set("hypervisor", instance.Hypervisor)

	// what field to map?
	//d.Set("vm_profile", map[string]string{
	//	"resource_id":   "",
	//	"vm_profile_id": iamInstanceProfileArnToName(instance.IamInstanceProfile),
	//})
	d.Set("image_id", instance.ImageId)
	d.Set("vm_id", instance.VmId)
	// what field to map?
	//d.Set("spot_vm", instance.SpotInstanceRequestId)
	d.Set("state", map[string]interface{}{
		"state_code": instance.State,
		"state_name": instance.State,
	})
	d.Set("type", instance.VmType)
	d.Set("public_ip", instance.PublicIp)
	// what field to map?
	//d.Set("kernel_id", instance.KernelId)
	d.Set("keypair_name", instance.KeypairName)

	// what field to map?
	//if instance.Monitoring != nil && instance.Monitoring.State != nil {
	//	monitoringState := *instance.Monitoring.State
	//	d.Set("monitoring", map[string]interface{}{
	//		"state": monitoringState == "enabled" || monitoringState == "pending",
	//	})
	//}

	d.Set("private_dns_name", instance.PrivateDnsName)
	d.Set("private_ip", instance.PrivateIp)
	//TODO:OAPI d.Set("nics", getOAPIVMNetworkInterfaceSet(instance.NetworkInterfaces))

	d.Set("placement", map[string]interface{}{
		// How to map this field?
		//"affinity":        instance.Placement.Affinity,
		"sub_region_name": instance.Placement.SubregionName,
		// TODO: Add to struct for OAPI
		// "firewall_rules_set_name": instance.Placement.FirewallRulesSetName,
		// How to map these fields?
		//"dedicated_host_id": instance.Placement.HostId,
		//"tenancy":           instance.Placement.Tenancy,
	})

	// TODO: Add to struct for OAPI
	// d.Set("system", instance.System)

	d.Set("product_codes", getOAPIVMProductCodes(instance.ProductCodes))
	// How to map this field?
	//d.Set("ramdisk_id", map[string]interface{}{
	//	"comment": instance.RamdiskId,
	//})
	d.Set("root_device_name", instance.RootDeviceName)
	d.Set("root_device_type", instance.RootDeviceType)
	d.Set("nat_check", instance.IsSourceDestChecked)
	// How to map these fields?
	//d.Set("spot_vm_request_id", instance.SpotInstanceRequestId)
	//d.Set("sriov_net_support", instance.SriovNetSupport)
	d.Set("comment", map[string]interface{}{
		// "state_code": instance.StateReason.Code,
		// "message": instance.StateReason.Message,
	})
	d.Set("subnet_id", instance.SubnetId)
	d.Set("tag_set", getOapiTagSet(instance.Tags))
	// How to map this field?
	//d.Set("virtualization_type", instance.VirtualizationType)
	d.Set("lin_id", instance.NetId)

	return nil
}

//Missing on Swagger spec
// func getOAPIVMBlockDeviceMapping(blockDeviceMappings []oapi.BlockDeviceMappingVmCreation) []map[string]interface{} {
// 	var blockDeviceMapping []map[string]interface{}

// 	if len(blockDeviceMappings) > 0 {
// 		blockDeviceMapping = make([]map[string]interface{}, len(blockDeviceMappings))
// 		for _, mapping := range blockDeviceMappings {
// 			r := map[string]interface{}{}
// 			r["device_name"] = mapping.DeviceName

// 			bsu := map[string]interface{}{}
// 			bsu["delete_on_vm_deletion"] = mapping.Bsu.DeleteOnVmDeletion
// 			bsu["state"] = mapping.Bsu.State
// 			bsu["volume_id"] = mapping.Bsu.VolumeId
// 			r["bsu"] = bsu

// 			blockDeviceMapping = append(blockDeviceMapping, r)
// 		}
// 	} else {
// 		blockDeviceMapping = make([]map[string]interface{}, 0)
// 	}
// 	return blockDeviceMapping
// }

func getOAPIVMGroupSet(groupSet []*fcu.GroupIdentifier) []map[string]interface{} {
	res := []map[string]interface{}{}
	for _, g := range groupSet {

		r := map[string]interface{}{
			"group_id":   g.GroupId,
			"group_name": g.GroupName,
		}
		res = append(res, r)
	}

	return res
}

func getOAPIVMTagSet(tagSet []*fcu.Tag) []map[string]interface{} {
	res := []map[string]interface{}{}
	for _, t := range tagSet {

		r := map[string]interface{}{
			"key":   t.Key,
			"value": t.Value,
		}
		res = append(res, r)
	}

	return res
}

func getOAPIVMProductCodes(productCode []string) []map[string]interface{} {
	res := []map[string]interface{}{}
	for _, p := range productCode {

		r := map[string]interface{}{
			"product_code": p,
		}
		res = append(res, r)
	}

	return res
}

func getOAPIVMPrivateIPAddressSet(privateIPs []*fcu.InstancePrivateIpAddress) []map[string]interface{} {
	res := []map[string]interface{}{}
	if privateIPs != nil {
		for _, p := range privateIPs {
			var inter map[string]interface{}

			assoc := map[string]interface{}{}
			assoc["ip_owner_id"] = p.Association.IpOwnerId
			assoc["public_dns_name"] = p.Association.PublicDnsName
			assoc["public_ip"] = p.Association.PublicIp

			inter["association"] = assoc
			inter["private_dns_name"] = p.Primary
			inter["private_ip_address"] = p.PrivateIpAddress

			res = append(res, inter)
		}
	}
	return res
}

func getDataSourceOAPIVMSchemas() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		//Attributes
		"filter": dataSourceFiltersSchema(),
		"vm_id": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"launch_sort_number": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"architecture": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"block_device_mapping": {
			Type:     schema.TypeMap,
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
		"bsu_optimised": {
			Type:     schema.TypeBool,
			Computed: true,
		},
		"firewall_rules_set": {
			Type: schema.TypeMap,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"firewall_rules_set_id": {
						Type:     schema.TypeInt,
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
					"state_name": {
						Type:     schema.TypeString,
						Computed: true,
					},
				},
			},
			Computed: true,
		},
		"type": {
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
			Type:     schema.TypeMap,
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
					"description": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"firewall_rules_set": {
						Type: schema.TypeMap,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"firewall_rules_set_id": {
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
						Type:     schema.TypeMap,
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
					"nat_check": {
						Type:     schema.TypeBool,
						Computed: true,
					},
					"state": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"subnet_id": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"lin_id": {
						Type:     schema.TypeInt,
						Computed: true,
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
		"system": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"private_dns_name": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"private_ip": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"product_codes": {
			Type:     schema.TypeMap,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"product_code": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"product_type": {
						Type:     schema.TypeString,
						Computed: true,
					},
				},
			},
		},
		"ramdisk_id": {
			Type:     schema.TypeMap,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"comment": {
						Type:     schema.TypeString,
						Computed: true,
					},
				},
			},
		},
		"root_device_name": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"root_device_type": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"nat_check": {
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
		"comment": {
			Type:     schema.TypeMap,
			Computed: true,
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
		},
		"subnet_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"tag_set": {
			Type:     schema.TypeMap,
			Computed: true,
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
		},
		"virtualization_type": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"lin_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
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
		case "account-id":
			filters.AccountIds = filterValues
		case "activated-check":
			filters.ActivatedCheck, _ = strconv.ParseBool(filterValues[0])
		case "architecture":
			filters.Architectures = filterValues
		case "block-device-mapping-delete-on-vm-deletion":
			filters.BlockDeviceMappingDeleteOnVmDeletion, _ = strconv.ParseBool(filterValues[0])
		case "block-device-mapping-device-name":
			filters.BlockDeviceMappingDeviceNames = filterValues
		case "block-device-mapping-link-date":
			filters.BlockDeviceMappingLinkDates = filterValues
		case "block-device-mapping-state":
			filters.BlockDeviceMappingStates = filterValues
		case "block-device-mapping-volume-id":
			filters.BlockDeviceMappingVolumeIds = filterValues
		case "comment":
			filters.Comments = filterValues
		case "creation-date":
			filters.CreationDates = filterValues
		case "dns-name":
			filters.DnsNames = filterValues
		case "hypervisor":
			filters.Hypervisors = filterValues
		case "image-id":
			filters.ImageIds = filterValues
		case "kernel-id":
			filters.KernelIds = filterValues
		case "keypair-name":
			filters.KeypairNames = filterValues
		case "launch-sort-number":
			filters.LaunchSortNumbers, _ = sliceAtoi(filterValues)
		case "link-nic-delete-on-vm-deletion":
			filters.LinkNicDeleteOnVmDeletion, _ = strconv.ParseBool(filterValues[0])
		case "link-nic-link-date":
			filters.LinkNicLinkDates = filterValues
		case "link-nic-link-nic-id":
			filters.LinkNicLinkNicIds = filterValues
		case "link-nic-link-public-ip-id":
			filters.LinkNicLinkPublicIpIds = filterValues
		case "link-nic-nic-id":
			filters.LinkNicNicIds = filterValues
		case "link-nic-nic-sort-number":
			filters.LinkNicNicSortNumbers, _ = sliceAtoi(filterValues)
		case "link-nic-public-ip-account-id":
			filters.LinkNicPublicIpAccountIds = filterValues
		case "link-nic-public-ip-id":
			filters.LinkNicPublicIpIds = filterValues
		case "link-nic-public-ip":
			filters.LinkNicPublicIps = filterValues
		case "link-nic-state":
			filters.LinkNicStates = filterValues
		case "link-nic-vm-account-id":
			filters.LinkNicVmAccountIds = filterValues
		case "link-nic-vm-id":
			filters.LinkNicVmIds = filterValues
		case "monitoring-state":
			filters.MonitoringStates = filterValues
		case "net-id":
			filters.NetIds = filterValues
		case "nic-account-id":
			filters.NicAccountIds = filterValues
		case "nic-activated-check":
			filters.NicActivatedCheck, _ = strconv.ParseBool(filterValues[0])
		case "nic-description":
			filters.NicDescriptions = filterValues
		case "nic-mac-address":
			filters.NicMacAddresses = filterValues
		case "nic-net-id":
			filters.NicNetIds = filterValues
		case "nic-nic-id":
			filters.NicNicIds = filterValues
		case "nic-private-dns-name":
			filters.NicPrivateDnsNames = filterValues
		case "nic-security-group-id":
			filters.NicSecurityGroupIds = filterValues
		case "nic-security-group-name":
			filters.NicSecurityGroupNames = filterValues
		case "nic-state":
			filters.NicStates = filterValues
		case "nic-subnet-id":
			filters.NicSubnetIds = filterValues
		case "nic-subregion-name":
			filters.NicSubregionNames = filterValues
		case "placement-group":
			filters.PlacementGroups = filterValues
		case "private-dns-name":
			filters.PrivateDnsNames = filterValues
		case "private-ip-link-private-ip-account-id":
			filters.PrivateIpLinkPrivateIpAccountIds = filterValues
		case "private-ip-link-public-ip":
			filters.PrivateIpLinkPublicIps = filterValues
		case "private-ip-primary-ip":
			filters.PrivateIpPrimaryIps = filterValues
		case "private-ip-private-ip":
			filters.PrivateIpPrivateIps = filterValues
		case "private-ip":
			filters.PrivateIps = filterValues
		case "product-code":
			filters.ProductCodes = filterValues
		case "public-ip":
			filters.PublicIps = filterValues
		case "ram-disk-id":
			filters.RamDiskIds = filterValues
		case "root-device-name":
			filters.RootDeviceNames = filterValues
		case "root-device-type":
			filters.RootDeviceTypes = filterValues
		case "security-group-id":
			filters.SecurityGroupIds = filterValues
		case "security-group-name":
			filters.SecurityGroupNames = filterValues
		case "spot-vm-request-id":
			filters.SpotVmRequestIds = filterValues
		case "spot-vm":
			filters.SpotVms = filterValues
		case "state-comment":
			filters.StateComments = filterValues
		case "subnet-id":
			filters.SubnetIds = filterValues
		case "subregion-name":
			filters.SubregionNames = filterValues
		case "system":
			filters.Systems = filterValues
		case "tag-key":
			filters.TagKeys = filterValues
		case "tag-value":
			filters.TagValues = filterValues
		case "tag":
			filters.Tags = filterValues
		case "tenancy":
			filters.Tenancies = filterValues
		case "token":
			filters.Tokens = filterValues
		case "virtualization-type":
			filters.VirtualizationTypes = filterValues
		case "vm-id":
			filters.VmIds = filterValues
		case "vm-state":
			filters.VmStates = filterValues
		case "vm-type":
			filters.VmTypes = filterValues
		case "vms-security-group-id":
			filters.VmsSecurityGroupIds = filterValues
		case "vms-security-group-name":
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
