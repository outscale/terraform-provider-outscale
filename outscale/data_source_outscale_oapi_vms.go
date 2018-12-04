package outscale

import (
	"errors"
	"fmt"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-outscale/osc/oapi"
)

func datasourceOutscaleOApiVMS() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOutscaleOApiVMSRead,

		Schema: datasourceOutscaleOApiVMSSchema(),
	}
}

func datasourceOutscaleOApiVMSSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"filter": dataSourceFiltersSchema(),
		"vm_id": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"reservation_set": &schema.Schema{
			Type:     schema.TypeSet,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"firewall_rules_set": &schema.Schema{
						Type:     schema.TypeSet,
						Computed: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"firewall_rules_set_id": &schema.Schema{
									Type:     schema.TypeString,
									Computed: true,
								},
								"firewall_rules_set_name": &schema.Schema{
									Type:     schema.TypeString,
									Computed: true,
								},
							},
						},
					},
					"vm": &schema.Schema{
						Type:     schema.TypeSet,
						Computed: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"architecture": &schema.Schema{
									Type:     schema.TypeString,
									Computed: true,
								},
								"block_device_mapping": &schema.Schema{
									Type:     schema.TypeList,
									Computed: true,
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"device_name": &schema.Schema{
												Type:     schema.TypeString,
												Computed: true,
											},
											"bsu": &schema.Schema{
												Type:     schema.TypeSet,
												Computed: true,
												Elem: &schema.Resource{
													Schema: map[string]*schema.Schema{
														"delete_on_vm_deletion": &schema.Schema{
															Type:     schema.TypeBool,
															Computed: true,
														},
														"state": &schema.Schema{
															Type:     schema.TypeString,
															Computed: true,
														},
														"volume_id": &schema.Schema{
															Type:     schema.TypeString,
															Computed: true,
														},
													},
												},
											},
										},
									},
								},
								"token": &schema.Schema{
									Type:     schema.TypeString,
									Computed: true,
								},
								"public_dns_name": &schema.Schema{
									Type:     schema.TypeString,
									Computed: true,
								},
								"bsu_optimised": &schema.Schema{
									Type:     schema.TypeBool,
									Computed: true,
								},
								"firewall_rules_set": &schema.Schema{
									Type:     schema.TypeSet,
									Computed: true,
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"firewall_rules_set_id": &schema.Schema{
												Type:     schema.TypeString,
												Computed: true,
											},
											"firewall_rules_set_name": &schema.Schema{
												Type:     schema.TypeString,
												Computed: true,
											},
										},
									},
								},
								"hypervisor": &schema.Schema{
									Type:     schema.TypeString,
									Computed: true,
								},
								"vm_profile": &schema.Schema{
									Type:     schema.TypeSet,
									Computed: true,
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"resource_id": &schema.Schema{
												Type:     schema.TypeString,
												Computed: true,
											},
											"vm_profile_id": &schema.Schema{
												Type:     schema.TypeString,
												Computed: true,
											},
										},
									},
								},
								"image_id": &schema.Schema{
									Type:     schema.TypeString,
									Computed: true,
								},
								"vm_id": &schema.Schema{
									Type:     schema.TypeSet,
									Optional: true,
									Elem: &schema.Schema{
										Type: schema.TypeString,
									},
								},
								"spot_vm": &schema.Schema{
									Type:     schema.TypeString,
									Computed: true,
								},
								"state": &schema.Schema{
									Type:     schema.TypeSet,
									Computed: true,
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"state_code": &schema.Schema{
												Type:     schema.TypeInt,
												Computed: true,
											},
											"state_name": &schema.Schema{
												Type:     schema.TypeString,
												Computed: true,
											},
										},
									},
								},
								"type": &schema.Schema{
									Type:     schema.TypeString,
									Computed: true,
								},
								"public_ip": &schema.Schema{
									Type:     schema.TypeString,
									Computed: true,
								},
								"kernel_id": &schema.Schema{
									Type: schema.TypeString,

									Computed: true,
								},
								"keypair_name": &schema.Schema{
									Type:     schema.TypeString,
									Computed: true,
								},
								"monitoring": &schema.Schema{
									Type: schema.TypeSet,

									Computed: true,
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"state": &schema.Schema{
												Type:     schema.TypeString,
												Computed: true,
											},
										},
									},
								},
								"nic": &schema.Schema{
									Type: schema.TypeList,

									Computed: true,
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"public_ip_link": &schema.Schema{
												Type:     schema.TypeSet,
												Computed: true,
												Elem: &schema.Resource{
													Schema: map[string]*schema.Schema{
														"public_ip_account_id": &schema.Schema{
															Type: schema.TypeString,

															Computed: true,
														},
														"public_dns_name": &schema.Schema{
															Type: schema.TypeString,

															Computed: true,
														},
														"public_ip": &schema.Schema{
															Type: schema.TypeString,

															Computed: true,
														},
													},
												},
											},
											"nic_link": &schema.Schema{
												Type:     schema.TypeSet,
												Computed: true,
												Elem: &schema.Resource{
													Schema: map[string]*schema.Schema{
														"nic_link_id": &schema.Schema{
															Type:     schema.TypeString,
															Computed: true,
														},
														"delete_on_vm_deletion": &schema.Schema{
															Type:     schema.TypeBool,
															Computed: true,
														},
														"nic_sort_number": &schema.Schema{
															Type:     schema.TypeInt,
															Computed: true,
														},
														"state": &schema.Schema{
															Type:     schema.TypeString,
															Computed: true,
														},
													},
												},
											},
											"description": &schema.Schema{
												Type:     schema.TypeString,
												Computed: true,
											},
											"firewall_rules_set": &schema.Schema{
												Type:     schema.TypeList,
												Computed: true,
												Elem: &schema.Resource{
													Schema: map[string]*schema.Schema{
														"firewall_rules_set_id": &schema.Schema{
															Type:     schema.TypeString,
															Computed: true,
														},
														"firewall_rules_set_name": &schema.Schema{
															Type:     schema.TypeString,
															Computed: true,
														},
													},
												},
											},
											"mac_address": &schema.Schema{
												Type:     schema.TypeString,
												Computed: true,
											},
											"nic_id": &schema.Schema{
												Type:     schema.TypeString,
												Computed: true,
											},
											"account_id": &schema.Schema{
												Type:     schema.TypeString,
												Computed: true,
											},
											"private_dns_name": &schema.Schema{
												Type:     schema.TypeString,
												Computed: true,
											},
											"private_ip": &schema.Schema{
												Type:     schema.TypeString,
												Computed: true,
											},
											"private_ip_set": &schema.Schema{
												Type:     schema.TypeSet,
												Computed: true,
												Elem: &schema.Resource{
													Schema: map[string]*schema.Schema{
														"public_ip_link": &schema.Schema{
															Type:     schema.TypeSet,
															Computed: true,
															Elem: &schema.Resource{
																Schema: map[string]*schema.Schema{
																	"public_ip_account_id": &schema.Schema{
																		Type:     schema.TypeString,
																		Computed: true,
																	},
																	"public_dns_name": &schema.Schema{
																		Type:     schema.TypeString,
																		Computed: true,
																	},
																	"public_ip": &schema.Schema{
																		Type:     schema.TypeString,
																		Computed: true,
																	},
																},
															},
														},
														"primary_ip": &schema.Schema{
															Type:     schema.TypeBool,
															Computed: true,
														},
														"private_dns_name": &schema.Schema{
															Type:     schema.TypeString,
															Computed: true,
														},
														"private_ip": &schema.Schema{
															Type:     schema.TypeString,
															Computed: true,
														},
													},
												},
											},
											"nat_check": &schema.Schema{
												Type:     schema.TypeBool,
												Computed: true,
											},
											"state": &schema.Schema{
												Type:     schema.TypeString,
												Computed: true,
											},
											"subnet_id": &schema.Schema{
												Type:     schema.TypeString,
												Computed: true,
											},
											"lin_id": &schema.Schema{
												Type:     schema.TypeString,
												Computed: true,
											},
										},
									},
								},
								"placement": &schema.Schema{
									Type:     schema.TypeSet,
									Computed: true,
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"affinity": &schema.Schema{
												Type:     schema.TypeString,
												Computed: true,
											},
											"sub_region_name": &schema.Schema{
												Type:     schema.TypeString,
												Computed: true,
											},
											"firewall_rules_set_name": &schema.Schema{
												Type:     schema.TypeString,
												Computed: true,
											},
											"dedicated_host_id": &schema.Schema{
												Type:     schema.TypeString,
												Computed: true,
											},
											"tenancy": &schema.Schema{
												Type:     schema.TypeString,
												Computed: true,
											},
										},
									},
								},

								"system": &schema.Schema{
									Type:     schema.TypeString,
									Computed: true,
								},
								"private_dns_name": &schema.Schema{
									Type:     schema.TypeString,
									Computed: true,
								},
								"private_ip": &schema.Schema{
									Type:     schema.TypeString,
									Computed: true,
								},
								"product_codes": &schema.Schema{
									Type:     schema.TypeList,
									Computed: true,
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"product_code": &schema.Schema{
												Type:     schema.TypeString,
												Computed: true,
											},
											"product_type": &schema.Schema{
												Type:     schema.TypeString,
												Computed: true,
											},
										},
									},
								},

								"ramdisk_id": &schema.Schema{
									Type:     schema.TypeList,
									Computed: true,
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"comment": &schema.Schema{
												Type:     schema.TypeString,
												Computed: true,
											},
										},
									},
								},
								"root_device_name": &schema.Schema{
									Type:     schema.TypeString,
									Computed: true,
								},
								"root_device_type": &schema.Schema{
									Type:     schema.TypeString,
									Computed: true,
								},
								"nat_check": &schema.Schema{
									Type:     schema.TypeBool,
									Computed: true,
								},
								"spot_vm_request_id": &schema.Schema{
									Type:     schema.TypeString,
									Computed: true,
								},
								"sriov_net_support": &schema.Schema{
									Type:     schema.TypeString,
									Computed: true,
								},
								"comment": &schema.Schema{
									Type:     schema.TypeSet,
									Computed: true,
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"state_code": &schema.Schema{
												Type:     schema.TypeString,
												Computed: true,
											},
											"comment": &schema.Schema{
												Type:     schema.TypeString,
												Computed: true,
											},
										},
									},
								},
								"subnet_id": &schema.Schema{
									Type:     schema.TypeString,
									Computed: true,
								},
								"tag_set": &schema.Schema{
									Type:     schema.TypeList,
									Computed: true,
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"key": &schema.Schema{
												Type:     schema.TypeString,
												Computed: true,
											},
											"value": &schema.Schema{
												Type:     schema.TypeString,
												Computed: true,
											},
										},
									},
								},
								"virtualization_type": &schema.Schema{
									Type:     schema.TypeString,
									Computed: true,
								},
								"lin_id": &schema.Schema{
									Type:     schema.TypeString,
									Computed: true,
								},
							},
						},
					},
					"account_id": &schema.Schema{
						Type:     schema.TypeString,
						Computed: true,
					},
					"requester_id": &schema.Schema{
						Type:     schema.TypeString,
						Computed: true,
					},
					"reservation_id": &schema.Schema{
						Type:     schema.TypeString,
						Computed: true,
					},
					"admin_password__": &schema.Schema{
						Type:     schema.TypeString,
						Computed: true,
					},
				},
			},
		},
	}
}

func dataSourceOutscaleOApiVMSRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*OutscaleClient).OAPI

	filters, filtersOk := d.GetOk("filter")
	vmID, vmIDOk := d.GetOk("vm_id")

	if filtersOk == false && vmIDOk == false {
		return fmt.Errorf("One of filters, and vm ID must be assigned")
	}

	// Build up search parameters
	params := oapi.ReadVmsRequest{}
	if filtersOk {
		params.Filters = buildOutscaleOAPIDataSourceVmFilters(filters.(*schema.Set))
	}
	if vmIDOk {
		params.Filters.VmIds = []string{vmID.(string)}
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

	// TODO: add Firewall struct
	// var firewallRules []*fcu.FirewallRules

	// loop through reservations, and remove terminated instances, populate instance slice
	for _, res := range resp.OK.Vms {
		if res.State != "terminated" {
			filteredInstances = append(filteredInstances, res)
		}

		d.Set("requester_id", resp.OK.ResponseContext.RequestId)
		d.Set("reservation_id", resp.OK.Vms[0].ReservationId)
		// TODO: add the following in the struct
		// account_id & admin_password__
	}

	if len(filteredInstances) < 1 {
		return errors.New("Your query returned no results. Please change your search criteria and try again")
	}

	return vmsOAPIDescriptionAttributes(d, filteredInstances, client)
}

// Populate instance attribute fields with the returned instance
func vmsOAPIDescriptionAttributes(d *schema.ResourceData, instances []oapi.Vm, conn *oapi.Client) error {
	d.Set("vm", dataSourceOAPIVMS(instances))
	return nil
}

func dataSourceOAPIVMS(i []oapi.Vm) *schema.Set {
	s := &schema.Set{}
	for _, v := range i {
		instance := map[string]interface{}{
			"launch_sort_number": v.LaunchNumber,
			"architecture":       v.Architecture,

			// TODO: has different struct for OAPI
			// "blocking_device_mapping": v.BlockDeviceMappings,

			"token":           v.ClientToken,
			"public_dns_name": v.PublicDnsName,
			"bsu_optimized":   v.BsuOptimized,

			// TODO: has different struct for OAPI
			// "group_set":                v.GroupSet,

			"hypervisor": v.Hypervisor,

			// "iam_instance_profile":     iamInstanceProfileArnToName(v.IamInstanceProfile),

			"image_id": v.ImageId,
			"vm_id":    v.VmId,

			// "instance_lifecycle":       v.InstanceLifecycle,

			// TODO: has different struct for OAPI
			// "instance_state":           v.InstanceState,

			"type":      v.VmType,
			"public_ip": v.PublicIp,
			// how to map?
			//"kernel_id":    v.KernelId,
			"keypair_name": v.KeypairName,
			// how to map?
			//"monitoring":   v.Monitoring,

			// TODO: has different struct for OAPI
			// "network_interfaces":       v.NetworkInterfaces,

			"placement": v.Placement,
			// how to map?
			//"system":           v.Platform,
			"private_dns_name": v.PrivateDnsName,
			"private_ip":       v.PrivateIp,
			"product_codes":    v.ProductCodes,

			// TODO: has different struct for OAPI
			// "ramdisk_id":               v.RamdiskId,

			// "reason":                   v.Reason,

			// TODO: Missing in struct for OAPI
			// "root_device_type":         v.RootDeviceType,

			"root_device_type": v.RootDeviceType,
			"nat_check":        v.IsSourceDestChecked,
			// how to map?
			//"spot_vm_request_id": v.SpotInstanceRequestId,
			//"sriov_net_support":  v.SriovNetSupport,

			// TODO: Missing in struct for OAPI
			// "state":               v.State,

			// "state_reason":        v.StateReason,

			"subnet_id": v.SubnetId,
			"tags":      v.Tags,
			// how to map?
			//"virtualization_type": v.VirtualizationType,
			"lin_id": v.NetId,
		}
		s.Add(instance)
	}
	return s
}

func dataSourceFiltersOApiSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		ForceNew: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name": {
					Type:     schema.TypeString,
					Required: true,
				},

				"values": {
					Type:     schema.TypeList,
					Required: true,
					Elem:     &schema.Schema{Type: schema.TypeString},
				},
			},
		},
	}
}
