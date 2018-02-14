package outscale

import (
	"errors"
	"fmt"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-outscale/osc/fcu"
)

func datasourceOutscaleVMS() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOutscaleVMSRead,

		Schema: datasourceOutscaleVMSSchema(),
	}
}

func datasourceOutscaleVMSSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"filter": dataSourceFiltersSchema(),
		"reservation_set": &schema.Schema{
			Type:     schema.TypeSet,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"group_set": &schema.Schema{
						Type:     schema.TypeSet,
						Computed: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"group_id": &schema.Schema{
									Type:     schema.TypeString,
									Computed: true,
								},
								"group_name": &schema.Schema{
									Type:     schema.TypeString,
									Computed: true,
								},
							},
						},
					},
					"instances_set": {
						Type:     schema.TypeSet,
						Computed: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"ami_launch_index": {
									Type:     schema.TypeInt,
									Computed: true,
								},
								"architecture": {
									Type:     schema.TypeString,
									Computed: true,
								},
								"block_device_mapping": {
									Type: schema.TypeSet,
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"device_name": {
												Type:     schema.TypeString,
												Computed: true,
											},
											"ebs": {
												Type: schema.TypeSet,
												Elem: &schema.Resource{
													Schema: map[string]*schema.Schema{
														"delete_on_termination": {
															Type:     schema.TypeBool,
															Computed: true,
														},
														"status": {
															Type:     schema.TypeString,
															Computed: true,
														},
														"volume_id": {
															Type:     schema.TypeString,
															Computed: true,
														},
													},
												},
												Computed: true,
											},
										},
									},
									Computed: true,
								},
								"client_token": {
									Type:     schema.TypeString,
									Computed: true,
								},
								"dns_name": {
									Type:     schema.TypeString,
									Computed: true,
								},
								"ebs_optimised": {
									Type:     schema.TypeBool,
									Computed: true,
								},
								"group_set": {
									Type:     schema.TypeSet,
									Computed: true,
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"group_id": {
												Type:     schema.TypeInt,
												Computed: true,
											},
											"group_name": {
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
								"iam_instance_profile": {
									Type: schema.TypeSet,
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"arn": {
												Type:     schema.TypeString,
												Computed: true,
											},
											"id": {
												Type:     schema.TypeString,
												Computed: true,
											},
										},
									},
									Computed: true,
								},
								"image_id": {
									Type:     schema.TypeInt,
									Computed: true,
								},
								"instance_id": {
									Type:     schema.TypeString,
									Computed: true,
								},
								"instance_lifecycle": {
									Type:     schema.TypeString,
									Computed: true,
								},
								"instance_state": {
									Type: schema.TypeSet,
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"code": {
												Type:     schema.TypeString,
												Computed: true,
											},
											"name": {
												Type:     schema.TypeString,
												Computed: true,
											},
										},
									},
									Computed: true,
								},
								"instance_type": {
									Type:     schema.TypeString,
									Computed: true,
								},
								"ip_address": {
									Type:     schema.TypeString,
									Computed: true,
								},
								"kernel_id": {
									Type:     schema.TypeString,
									Computed: true,
								},
								"key_name": {
									Type:     schema.TypeString,
									Computed: true,
								},
								"monitoring": {
									Type: schema.TypeSet,
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
								"network_interface_set": {
									Type:     schema.TypeSet,
									Computed: true,
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"association": {
												Type:     schema.TypeSet,
												Computed: true,
												Elem: &schema.Resource{
													Schema: map[string]*schema.Schema{
														"ip_owner_id": {
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
											"attachment": {
												Type:     schema.TypeSet,
												Computed: true,
												Elem: &schema.Resource{
													Schema: map[string]*schema.Schema{
														"attachement_id": {
															Type:     schema.TypeString,
															Computed: true,
														},
														"delete_on_termination": {
															Type:     schema.TypeBool,
															Computed: true,
														},
														"device_index": {
															Type:     schema.TypeInt,
															Computed: true,
														},
														"status": {
															Type:     schema.TypeString,
															Computed: true,
														},
													},
												},
											},
											"description": {
												Type:     schema.TypeString,
												Computed: true,
											},
											"group_set": {
												Type:     schema.TypeSet,
												Computed: true,
												Elem: &schema.Resource{
													Schema: map[string]*schema.Schema{
														"group_id": {
															Type:     schema.TypeString,
															Computed: true,
														},
														"group_name": {
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
											"network_interface_id": {
												Type:     schema.TypeString,
												Computed: true,
											},
											"owner_id": {
												Type:     schema.TypeString,
												Computed: true,
											},
											"private_dns_name": {
												Type:     schema.TypeString,
												Computed: true,
											},
											"private_ip_address": {
												Type:     schema.TypeString,
												Computed: true,
											},
											"private_ip_addresses_set": {
												Type:     schema.TypeSet,
												Computed: true,
												Elem: &schema.Resource{
													Schema: map[string]*schema.Schema{
														"association": {
															Type:     schema.TypeSet,
															Computed: true,
															Elem: &schema.Resource{
																Schema: map[string]*schema.Schema{
																	"ip_owner_id": {
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
														"primary": {
															Type:     schema.TypeString,
															Computed: true,
														},
														"private_dns_name": {
															Type:     schema.TypeString,
															Computed: true,
														},
														"private_ip_address": {
															Type:     schema.TypeString,
															Computed: true,
														},
													},
												},
											},
											"source_dest_check": {
												Type:     schema.TypeBool,
												Computed: true,
											},
											"status": {
												Type:     schema.TypeString,
												Computed: true,
											},
											"subnet_id": {
												Type:     schema.TypeString,
												Computed: true,
											},
											"vpc_id": {
												Type:     schema.TypeInt,
												Computed: true,
											},
										},
									},
								},
								"placement": {
									Type:     schema.TypeSet,
									Computed: true,
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"affinity": {
												Type:     schema.TypeString,
												Computed: true,
											},
											"availability_zone": {
												Type:     schema.TypeString,
												Computed: true,
											},
											"group_name": {
												Type:     schema.TypeString,
												Computed: true,
											},
											"host_id": {
												Type:     schema.TypeString,
												Computed: true,
											},
											"tenancy": {
												Type:     schema.TypeString,
												Computed: true,
											},
										},
									},
								},
								"platform": {
									Type:     schema.TypeString,
									Computed: true,
								},
								"private_dns_name": {
									Type:     schema.TypeString,
									Computed: true,
								},
								"private_ip_address": {
									Type:     schema.TypeString,
									Computed: true,
								},
								"product_codes": {
									Type:     schema.TypeSet,
									Computed: true,
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"product_code": {
												Type:     schema.TypeString,
												Computed: true,
											},
											"type": {
												Type:     schema.TypeString,
												Computed: true,
											},
										},
									},
								},
								"ramdisk_id": {
									Type:     schema.TypeString,
									Computed: true,
								},
								"reason": {
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
								"source_dest_check": {
									Type:     schema.TypeString,
									Computed: true,
								},
								"sopt_instance_request_id": {
									Type:     schema.TypeString,
									Computed: true,
								},
								"sriov_net_support": {
									Type:     schema.TypeString,
									Computed: true,
								},
								"state_reason": {
									Type: schema.TypeSet,
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"code": {
												Type:     schema.TypeInt,
												Computed: true,
											},
											"message": {
												Type:     schema.TypeString,
												Computed: true,
											},
										},
									},
									Computed: true,
								},
								"subnet_id": {
									Type:     schema.TypeString,
									Computed: true,
								},
								"tag_set": {
									Type: schema.TypeSet,
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
								"virtualization_type": {
									Type:     schema.TypeString,
									Computed: true,
								},
								"vpc_id": {
									Type:     schema.TypeString,
									Computed: true,
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceOutscaleVMSRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*OutscaleClient).FCU.VM

	filters, filtersOk := d.GetOk("filter")

	if filtersOk == false {
		return fmt.Errorf("One of filters must be assigned")
	}

	// Build up search parameters
	params := &fcu.DescribeInstancesInput{}
	if filtersOk {
		params.Filters = buildOutscaleDataSourceFilters(filters.(*schema.Set))
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

	if len(filteredInstances) < 1 {
		return errors.New("Your query returned no results. Please change your search criteria and try again")
	}

	return vmsDescriptionAttributes(d, filteredInstances, client)
}

// Populate instance attribute fields with the returned instance
func vmsDescriptionAttributes(d *schema.ResourceData, instances []*fcu.Instance, conn fcu.VMService) error {
	d.Set("instances_set", dataSourceInstance(instances))
	return nil
}

func dataSourceInstance(i []*fcu.Instance) *schema.Set {
	s := &schema.Set{}
	for _, v := range i {
		instance := map[string]interface{}{
			"ami_launch_index":         v.AmiLaunchIndex,
			"architecture":             v.Architecture,
			"blocking_device_mapping":  v.BlockDeviceMappings,
			"client_token":             v.ClientToken,
			"dns_name":                 v.DnsName,
			"ebs_optimized":            v.EbsOptimized,
			"group_set":                v.GroupSet,
			"hypervisor":               v.Hypervisor,
			"iam_instance_profile":     iamInstanceProfileArnToName(v.IamInstanceProfile),
			"image_id":                 v.ImageId,
			"instance_id":              v.InstanceId,
			"instance_lifecycle":       v.InstanceLifecycle,
			"instance_state":           v.InstanceState,
			"ip_address":               v.IpAddress,
			"kernel_id":                v.KernelId,
			"key_name":                 v.KeyName,
			"monitoring":               v.Monitoring,
			"network_interfaces":       v.NetworkInterfaces,
			"placement":                v.Placement,
			"platform":                 v.Platform,
			"private_dns":              v.PrivateDnsName,
			"private_ip_address":       v.PrivateIpAddress,
			"product_codes":            v.ProductCodes,
			"ramdisk_id":               v.RamdiskId,
			"reason":                   v.Reason,
			"root_device_type":         v.RootDeviceType,
			"source_dest_check":        v.SourceDestCheck,
			"spot_instance_request_id": v.SpotInstanceRequestId,
			"sriov_net_support":        v.SriovNetSupport,
			"state":                    v.State,
			"state_reason":             v.StateReason,
			"subnet_id":                v.SubnetId,
			"tags":                     v.Tags,
			"virtualization_type":      v.VirtualizationType,
			"vpc_id":                   v.VpcId,
		}
		s.Add(instance)
	}
	return s
}

func dataSourceFiltersSchema() *schema.Schema {
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
