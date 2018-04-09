package outscale

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-outscale/osc/fcu"
)

func dataSourceOutscaleOAPIVM() *schema.Resource {
	return &schema.Resource{
		Read:   dataSourceOutscaleOAPIVMRead,
		Schema: getDataSourceOAPIVMSchemas(),
	}
}
func dataSourceOutscaleOAPIVMRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*OutscaleClient).FCU.VM

	filters, filtersOk := d.GetOk("filter")
	instanceID, instanceIDOk := d.GetOk("vm_id")

	if filtersOk == false && instanceIDOk == false {
		return fmt.Errorf("One of filters, or instance_id must be assigned")
	}

	// Build up search parameters
	params := &fcu.DescribeInstancesInput{}
	if filtersOk {
		params.Filters = buildOutscaleDataSourceFilters(filters.(*schema.Set))
	}
	if instanceIDOk {
		params.InstanceIds = []*string{aws.String(instanceID.(string))}
	}

	var resp *fcu.DescribeInstancesOutput
	var err error

	err = resource.Retry(30*time.Second, func() *resource.RetryError {
		resp, err = client.DescribeInstances(params)
		return resource.RetryableError(err)
	})

	if resp.Reservations == nil {
		return fmt.Errorf("Your query returned no results. Please change your search criteria and try again")
	}

	// If no instances were returned, return
	if len(resp.Reservations) == 0 {
		return fmt.Errorf("Your query returned no results. Please change your search criteria and try again")
	}

	var filteredInstances []*fcu.Instance

	// loop through reservations, and remove terminated instances, populate instance slice
	for _, res := range resp.Reservations {
		for _, instance := range res.Instances {
			if instance.State != nil && *instance.State.Name != "terminated" {
				filteredInstances = append(filteredInstances, instance)
			}
		}
	}

	var instance *fcu.Instance
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

	log.Printf("[DEBUG] outscale_vm - Single VM ID found: %s", *instance.InstanceId)

	return oapiVMDescriptionAttributes(d, instance, client)
}

// Populate instance attribute fields with the returned instance
func oapiVMDescriptionAttributes(d *schema.ResourceData, instance *fcu.Instance, conn fcu.VMService) error {
	d.SetId(*instance.InstanceId)
	// Set the easy attributes

	d.Set("launch_sort_number", instance.AmiLaunchIndex)
	d.Set("architecture", instance.Architecture)
	// d.Set("block_device_mapping", getOAPIVMBlockDeviceMapping(instance.BlockDeviceMappings))
	d.Set("token", instance.ClientToken)
	d.Set("public_dns_name", instance.DnsName)
	d.Set("bsu_optimised", instance.EbsOptimized)

	// TODO: add to struct for OAPI
	// d.Set("firewall_rules_set", instance.FirewallRulesSets)

	d.Set("hypervisor", instance.Hypervisor)
	d.Set("vm_profile", map[string]string{
		"resource_id":   "",
		"vm_profile_id": iamInstanceProfileArnToName(instance.IamInstanceProfile),
	})
	d.Set("image_id", instance.ImageId)
	d.Set("vm_id", instance.InstanceId)
	d.Set("spot_vm", instance.SpotInstanceRequestId)
	d.Set("state", map[string]interface{}{
		"state_code": instance.State.Code,
		"state_name": instance.State.Name,
	})
	d.Set("type", instance.InstanceType)
	d.Set("public_ip", instance.IpAddress)
	d.Set("kernel_id", instance.KernelId)
	d.Set("keypair_name", instance.KeyName)

	if instance.Monitoring != nil && instance.Monitoring.State != nil {
		monitoringState := *instance.Monitoring.State
		d.Set("monitoring", map[string]interface{}{
			"state": monitoringState == "enabled" || monitoringState == "pending",
		})
	}

	d.Set("private_dns_name", instance.PrivateDnsName)
	d.Set("private_ip", instance.PrivateIpAddress)
	d.Set("nics", getOAPIVMNetworkInterfaceSet(instance.NetworkInterfaces))

	d.Set("placement", map[string]interface{}{
		"affinity":        instance.Placement.Affinity,
		"sub_region_name": instance.Placement.GroupName,
		// TODO: Add to struct for OAPI
		// "firewall_rules_set_name": instance.Placement.FirewallRulesSetName,
		"dedicated_host_id": instance.Placement.HostId,
		"tenancy":           instance.Placement.Tenancy,
	})

	// TODO: Add to struct for OAPI
	// d.Set("system", instance.System)

	d.Set("product_codes", getOAPIVMProductCodes(instance.ProductCodes))
	d.Set("ramdisk_id", map[string]interface{}{
		"comment": instance.RamdiskId,
	})
	d.Set("root_device_name", instance.RootDeviceName)
	d.Set("root_device_type", instance.RootDeviceType)
	d.Set("nat_check", instance.SourceDestCheck)
	d.Set("spot_vm_request_id", instance.SpotInstanceRequestId)
	d.Set("sriov_net_support", instance.SriovNetSupport)
	d.Set("comment", map[string]interface{}{
		// "state_code": instance.StateReason.Code,
		// "message": instance.StateReason.Message,
	})
	d.Set("subnet_id", instance.SubnetId)
	d.Set("tags", getOAPIVMTagSet(instance.Tags))
	d.Set("virtualization_type", instance.VirtualizationType)
	d.Set("lin_id", instance.VpcId)

	return nil
}

func getOAPIVMBlockDeviceMapping(blockDeviceMappings []*fcu.InstanceBlockDeviceMapping) []map[string]interface{} {
	s := []map[string]interface{}{}
	for _, mapping := range blockDeviceMappings {
		r := map[string]interface{}{
			"device_name": mapping.DeviceName,
			"bsu": map[string]interface{}{
				"delete_on_vm_deletion": mapping.Ebs.DeleteOnTermination,
				"state":                 mapping.Ebs.Status,
				"volume_id":             mapping.Ebs.VolumeId,
			},
		}
		s = append(s, r)
	}
	return s
}

func getOAPIVMNetworkInterfaceSet(interfaces []*fcu.InstanceNetworkInterface) []map[string]interface{} {
	res := []map[string]interface{}{}

	if interfaces != nil {
		for _, i := range interfaces {
			assoc := map[string]interface{}{}

			assoc["public_ip_link"] = map[string]interface{}{
				"public_ip_account_id": i.Association.IpOwnerId,
				"public_dns_name":      i.Association.PublicDnsName,
				"public_ip":            i.Association.PublicIp,
			}

			// TODO: add to struct for OAPI
			assoc["nic_link"] = map[string]interface{}{
				"nic_link_id":              i.Attachment.AttachmentId,
				"delete_on_vm_termination": i.Attachment.DeleteOnTermination,
				"nic_sort_number":          i.Attachment.DeviceIndex,
				"state":                    i.Attachment.Status,
			}

			assoc["description"] = *i.Description

			// TODO: add to struct for OAPI
			// firewall := []map[string]string{}
			// for _, f := range i.FirewallRulesSets {
			// 	rule := map[string]string{
			// 		"firewall_rules_set_id": "",
			// 		"firewall_rules_name":   "",
			// 	}
			// 	firewall = append(firewall, rule)
			// }
			// assoc["firewall_rules_sets"] = firewall

			assoc["mac_address"] = i.MacAddress
			assoc["nic_id"] = i.Attachment.AttachmentId

			assoc["account_id"] = i.OwnerId

			assoc["private_dns_name"] = i.PrivateDnsName
			assoc["private_ip"] = i.PrivateIpAddress

			ips := []map[string]interface{}{}

			for _, p := range i.PrivateIpAddresses {
				ip := map[string]interface{}{
					"public_ip_link": map[string]interface{}{
						"public_ip_account_id": p.Association.IpOwnerId,
						"public_dns_name":      p.Association.PublicDnsName,
						"public_ip":            p.Association.PublicIp,
					},
					"primary_ip":       p.Primary,
					"private_dns_name": p.PrivateDnsName,
					"private_ip":       p.PrivateIpAddress,
				}
				ips = append(ips, ip)
			}
			assoc["private_ips"] = ips
			assoc["nat_check"] = i.SourceDestCheck
			assoc["state"] = i.Status
			assoc["subnet_id"] = i.SubnetId
			assoc["lin_id"] = i.VpcId

			res = append(res, assoc)
		}
	}

	return res
}

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

func getOAPIVMProductCodes(productCode []*fcu.ProductCode) []map[string]interface{} {
	res := []map[string]interface{}{}
	for _, p := range productCode {

		r := map[string]interface{}{
			"product_code": p.ProductCode,
			"product_type": p.Type,
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
		"tags": {
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
