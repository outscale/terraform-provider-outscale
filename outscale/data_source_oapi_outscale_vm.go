package outscale

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-outscale/osc/fcu"
)

func dataSourceOApiOutscaleVM() *schema.Resource {
	return &schema.Resource{
		Read:   dataSourceOApiOutscaleVMRead,
		Schema: getDataSourceOApiVMSchemas(),
	}
}
func dataSourceOApiOutscaleVMRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*OutscaleClient).FCU.VM

	filters, filtersOk := d.GetOk("filter")
	instanceID, instanceIDOk := d.GetOk("instance_id")

	if filtersOk == false && instanceIDOk == false {
		return fmt.Errorf("One of filters, or instance_id must be assigned")
	}

	// Build up search parameters
	params := &fcu.DescribeInstancesInput{}
	if filtersOk {
		params.Filters = buildOApiOutscaleDataSourceFilters(filters.(*schema.Set))
	}
	if instanceIDOk {
		params.InstanceIds = []*string{aws.String(instanceID.(string))}
	}

	// Perform the lookup
	resp, err := client.DescribeInstances(params)
	if err != nil {
		return err
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

	log.Printf("[DEBUG] aws_instance - Single Instance ID found: %s", *instance.InstanceId)

	return instanceDescriptionOApiAttributes(d, instance, client)
}

// Populate instance attribute fields with the returned instance
func instanceDescriptionOApiAttributes(d *schema.ResourceData, instance *fcu.Instance, conn fcu.VMService) error {
	d.SetId(*instance.InstanceId)
	// Set the easy attributes
	d.Set("instance_state", instance.State.Name)
	if instance.Placement != nil {
		d.Set("availability_zone", instance.Placement.AvailabilityZone)
	}
	if instance.Placement.Tenancy != nil {
		d.Set("tenancy", instance.Placement.Tenancy)
	}
	d.Set("ami", instance.ImageId)
	d.Set("instance_type", instance.InstanceType)
	d.Set("key_name", instance.KeyName)
	d.Set("private_dns", instance.PrivateDnsName)
	d.Set("private_ip", instance.PrivateIpAddress)
	d.Set("iam_instance_profile", iamInstanceOApiProfileArnToName(instance.IamInstanceProfile))

	// iterate through network interfaces, and set subnet, network_interface, public_addr
	if len(instance.NetworkInterfaces) > 0 {
		for _, ni := range instance.NetworkInterfaces {
			if *ni.Attachment.DeviceIndex == 0 {
				d.Set("subnet_id", ni.SubnetId)
				d.Set("network_interface_id", ni.NetworkInterfaceId)
				d.Set("associate_public_ip_address", ni.Association != nil)
			}
		}
	} else {
		d.Set("subnet_id", instance.SubnetId)
		d.Set("network_interface_id", "")
	}

	d.Set("ebs_optimized", instance.EbsOptimized)
	if instance.SubnetId != nil && *instance.SubnetId != "" {
		d.Set("source_dest_check", instance.SourceDestCheck)
	}

	if instance.Monitoring != nil && instance.Monitoring.State != nil {
		monitoringState := *instance.Monitoring.State
		d.Set("monitoring", monitoringState == "enabled" || monitoringState == "pending")
	}

	return nil
}

func buildOApiOutscaleDataSourceFilters(set *schema.Set) []*fcu.Filter {
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

func getDataSourceOApiVMSchemas() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		//Attributes
		"instance_id": {
			Type:     schema.TypeString,
			Required: true,
		},
		"firewall_rules_sets": {
			Type: schema.TypeSet,
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
		"vm": {
			Type:     schema.TypeSet,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"launch_sort_number": {
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
								"bsu": {
									Type: schema.TypeSet,
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"delete_on_vm_deleti": {
												Type:     schema.TypeBool,
												Computed: true,
											},
											"on": {
												Type:     schema.TypeString,
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
					"token": {
						Type:     schema.TypeString,
						Optional: true,
					},
					"public_dns_name": {
						Type:     schema.TypeString,
						Optional: true,
					},
					"bsu_optimized": {
						Type:     schema.TypeBool,
						Computed: true,
					},
					"firewall_rules_set": {
						Type: schema.TypeSet,
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
						Required: true,
					},
					"hypervisor": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"vm_profile": {
						Type: schema.TypeSet,
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
					"vm_id": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"spot_vm": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"state": {
						Type: schema.TypeSet,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"state_code": {
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
					"nic": {
						Type:     schema.TypeSet,
						Computed: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"public_ip_link": {
									Type:     schema.TypeSet,
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
									Type: schema.TypeSet,
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
							},
						},
					},
					"description": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"firewall_rules_sets": {
						Type: schema.TypeSet,
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
						Type:     schema.TypeSet,
						Computed: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"public_ip": {
									Type:     schema.TypeSet,
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
					"activated_check": {
						Type:     schema.TypeBool,
						Computed: true,
					},
					// "state": {
					// 	Type:     schema.TypeString,
					// 	Computed: true,
					// },
					"subnet_id": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"lin_id": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"placement": {
						Type: schema.TypeSet,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"affinity": {
									Type:     schema.TypeString,
									Computed: true,
								},
								"sub_regio_name": {
									Type:     schema.TypeString,
									Computed: true,
								},
								"firewall_rules_set_name": {
									Type:     schema.TypeString,
									Computed: true,
								},
								"dedicated_host_id": {
									Type:     schema.TypeString,
									Computed: true,
								},
								"tenancy": {
									Type:     schema.TypeString,
									Computed: true,
								},
							},
						},
						Computed: true,
					},
					"system": {
						Type:     schema.TypeString,
						Computed: true,
					},
					// "private_dns_name": {
					// 	Type:     schema.TypeString,
					// 	Computed: true,
					// },
					// "private_ip": {
					// 	Type:     schema.TypeString,
					// 	Computed: true,
					// },
					"product_codes": {
						Type: schema.TypeSet,
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
						Computed: true,
					},
					"ramdisk_id": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"comment": {
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
					// "activated_check": {
					// 	Type:     schema.TypeString,
					// 	Computed: true,
					// },
					"spot_vm_request_id": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"sriov_net_support": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"comments": {
						Type: schema.TypeSet,
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
						Computed: true,
					},
					// "subnet_id": {
					// 	Type:     schema.TypeString,
					// 	Computed: true,
					// },
					"tags": {
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
					//					"lin_id": {
					//						Type:     schema.TypeString,
					//						Computed: true,
					//					},
				},
			},
		},
		"account_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"requester_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"reservation_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"admin_password": {
			Type:     schema.TypeString,
			Computed: true,
		},
		//End of Attributes
	}
}

func iamInstanceOApiProfileArnToName(ip *fcu.IamInstanceProfile) string {
	if ip == nil || ip.Arn == nil {
		return ""
	}
	parts := strings.Split(*ip.Arn, "/")
	return parts[len(parts)-1]
}
