package outscale

import (
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceOutscaleVM() *schema.Resource {
	return &schema.Resource{
		Read:   dataSourceOutscaleVMRead,
		Schema: getDataSourceVMSchemas(),
	}
}
func dataSourceOutscaleVMRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func getDataSourceVMSchemas() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		//Attributes
		"instance_id": {
			Type:     schema.TypeString,
			Required: true,
		},
		"group_set": {

			Type:     schema.TypeSet,
			Optional: true,
			Elem: schema.Resource{
				Schema: map[string]*schema.Schema{
					"group_id": {
						Type:     schema.TypeString,
						Required: true,
					},
					"group_name": {
						Type:     schema.TypeString,
						Required: true,
					},
				},
			},
		},
		"instance_set": {
			Type:     schema.TypeSet,
			Optional: true,
			Elem: schema.Resource{
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
						Elem: schema.Resource{
							Schema: map[string]*schema.Schema{
								"device_name": {
									Type:     schema.TypeString,
									Computed: true,
								},
								"ebs": {
									Type: schema.TypeSet,
									Elem: schema.Resource{
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
						Type: schema.TypeSet,
						Elem: schema.Resource{
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
						Elem: schema.Resource{
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
						Elem: schema.Resource{
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
						Elem: schema.Resource{
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
						Elem: schema.Resource{
							Schema: map[string]*schema.Schema{
								"association": {
									Type: schema.TypeSet,
									Elem: schema.Resource{
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
									Type: schema.TypeSet,
									Elem: schema.Resource{
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
									Type: schema.TypeSet,
									Elem: schema.Resource{
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
									Computed: true,
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
									Type: schema.TypeSet,
									Elem: schema.Resource{
										Schema: map[string]*schema.Schema{
											"association": {
												Type: schema.TypeSet,
												Elem: schema.Resource{
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
						Type: schema.TypeSet,
						Elem: schema.Resource{
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
						Type: schema.TypeSet,
						Elem: schema.Resource{
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
						Elem: schema.Resource{
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
						Elem: schema.Resource{
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
		"owner_id": {
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
		"password_data": {
			Type:     schema.TypeString,
			Computed: true,
		},
		//End of Attributes
	}
}
