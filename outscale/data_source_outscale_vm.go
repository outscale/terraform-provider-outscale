package outscale

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/antihax/optional"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	oscgo "github.com/marinsalinas/osc-sdk-go"
	"github.com/terraform-providers/terraform-provider-outscale/osc/fcu"
)

func buildOutscaleDataSourceFilters(set *schema.Set) []*fcu.Filter {
	var filters []*fcu.Filter
	for _, v := range set.List() {
		m := v.(map[string]interface{})
		var filterValues []*string
		for _, e := range m["values"].([]interface{}) {
			filterValues = append(filterValues, aws.String(e.(string)))
		}
		filters = append(filters, &fcu.Filter{
			Name:   aws.String(m["name"].(string)),
			Values: filterValues,
		})
	}
	return filters
}

func dataSourceOutscaleOAPIVM() *schema.Resource {
	return &schema.Resource{
		Read:   dataSourceOutscaleOAPIVMRead,
		Schema: getDataSourceOAPIVMSchemas(),
	}
}
func dataSourceOutscaleOAPIVMRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*OutscaleClient).OSCAPI

	filters, filtersOk := d.GetOk("filter")
	instanceID, instanceIDOk := d.GetOk("vm_id")

	if filtersOk == false && instanceIDOk == false {
		return fmt.Errorf("One of filters, or instance_id must be assigned")
	}

	// Build up search parameters
	params := oscgo.ReadVmsRequest{}
	if filtersOk {
		params.Filters = buildOutscaleOAPIDataSourceVMFilters(filters.(*schema.Set))
	}
	if instanceIDOk {
		params.Filters.VmIds = &[]string{instanceID.(string)}
	}

	log.Printf("[DEBUG] ReadVmsRequest -> %+v\n", params)

	var resp oscgo.ReadVmsResponse
	err := resource.Retry(30*time.Second, func() *resource.RetryError {
		r, _, err := client.VmApi.ReadVms(context.Background(), &oscgo.ReadVmsOpts{
			ReadVmsRequest: optional.NewInterface(params),
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

	d.Set("request_id", resp.GetResponseContext().RequestId)

	// Populate vm attribute fields with the returned vm
	return resourceDataAttrSetter(d, func(set AttributeSetter) error {
		d.SetId(vm.GetVmId())
		set("tags", getOscAPITagSet(vm.GetTags()))
		return oapiVMDescriptionAttributes(set, &vm)
	})
}

// Populate instance attribute fields with the returned instance
func oapiVMDescriptionAttributes(set AttributeSetter, vm *oscgo.Vm) error {

	set("architecture", vm.GetArchitecture())
	if err := set("block_device_mappings_created", getOscAPIVMBlockDeviceMapping(vm.GetBlockDeviceMappings())); err != nil {
		log.Printf("[DEBUG] BLOCKING DEVICE MAPPING ERR %+v", err)
		return err
	}
	set("bsu_optimized", vm.GetBsuOptimized())
	set("client_token", vm.GetClientToken())
	set("deletion_protection", vm.GetDeletionProtection())
	set("hypervisor", vm.GetHypervisor())
	set("image_id", vm.GetImageId())
	set("is_source_dest_checked", vm.GetIsSourceDestChecked())
	set("keypair_name", vm.GetKeypairName())
	set("launch_number", vm.GetLaunchNumber())
	set("net_id", vm.GetNetId())

	if err := set("nics", getOAPIVMNetworkInterfaceLightSet(vm.GetNics())); err != nil {
		log.Printf("[DEBUG] NICS ERR %+v", err)
		return err
	}
	set("os_family", vm.GetOsFamily())
	set("placement_subregion_name", aws.StringValue(vm.GetPlacement().SubregionName))
	set("placement_tenancy", aws.StringValue(vm.GetPlacement().Tenancy))
	set("private_dns_name", vm.GetPrivateDnsName())
	set("private_ip", vm.GetPrivateIp())
	set("product_codes", vm.GetProductCodes())
	set("public_dns_name", vm.GetPublicDnsName())
	set("public_ip", vm.GetPublicIp())
	set("reservation_id", vm.GetReservationId())
	set("root_device_name", vm.GetRootDeviceName())
	set("root_device_type", vm.GetRootDeviceType())
	if err := set("security_groups", getOAPIVMSecurityGroups(vm.GetSecurityGroups())); err != nil {
		log.Printf("[DEBUG] SECURITY GROUPS ERR %+v", err)
		return err
	}
	set("state", vm.GetState())
	set("state_reason", vm.GetStateReason())
	set("subnet_id", vm.GetSubnetId())
	set("user_data", vm.GetUserData())
	set("vm_id", vm.GetVmId())
	set("vm_initiated_shutdown_behavior", vm.GetVmInitiatedShutdownBehavior())

	return set("vm_type", vm.GetVmType())
}

func getOscAPIVMBlockDeviceMapping(blockDeviceMappings []oscgo.BlockDeviceMappingCreated) []map[string]interface{} {
	blockDeviceMapping := make([]map[string]interface{}, len(blockDeviceMappings))

	for k, v := range blockDeviceMappings {
		blockDeviceMapping[k] = map[string]interface{}{
			"device_name": aws.StringValue(v.DeviceName),
			"bsu": map[string]interface{}{
				"delete_on_vm_deletion": fmt.Sprintf("%t", aws.BoolValue(v.GetBsu().DeleteOnVmDeletion)),
				"volume_id":             aws.StringValue(v.GetBsu().VolumeId),
				"state":                 aws.StringValue(v.GetBsu().State),
				"link_date":             aws.StringValue(v.GetBsu().LinkDate),
			},
		}
	}
	return blockDeviceMapping
}

func getOAPIVMSecurityGroups(groupSet []oscgo.SecurityGroupLight) []map[string]interface{} {
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

func buildOutscaleOAPIDataSourceVMFilters(set *schema.Set) *oscgo.FiltersVm {
	filters := new(oscgo.FiltersVm)

	for _, v := range set.List() {
		m := v.(map[string]interface{})
		filterValues := make([]string, 0)
		for _, e := range m["values"].([]interface{}) {
			filterValues = append(filterValues, e.(string))
		}

		switch name := m["name"].(string); name {
		case "tag_keys":
			filters.TagKeys = &filterValues
		case "tag_values":
			filters.TagValues = &filterValues
		case "tags":
			filters.Tags = &filterValues
		case "vm_ids":
			filters.VmIds = &filterValues
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

func getOApiVMAttributesSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		// Attributes
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
