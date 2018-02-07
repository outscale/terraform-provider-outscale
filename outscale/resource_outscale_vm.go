package outscale

import (
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceOutscaleVM() *schema.Resource {
	return &schema.Resource{
		Create: resourceVMCreate,
		Read:   resourceVMRead,
		Update: resourceVMUpdate,
		Delete: resourceVMDelete,

		Schema: getVMSchema(),
	}
}

func resourceVMCreate(d *schema.ResourceData, meta interface{}) error {
	return nil
}
func resourceVMRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}
func resourceVMUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}
func resourceVMDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func getVMSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		// Attributes
		"block_device_mapping": {
			Type: schema.TypeSet,
			Elem: schema.Resource{
				Schema: map[string]*schema.Schema{
					"device_name": {
						Type: schema.TypeString,
					},
					"ebs": {
						Type: schema.TypeSet,
						Elem: schema.Resource{
							Schema: map[string]*schema.Schema{
								"delete_on_termination": {
									Type: schema.TypeBool,
								},
								"iops": {
									Type: schema.TypeString,
								},
								"snapshot_id": {
									Type:     schema.TypeInt,
									Required: true,
								},
								"volume_size": {
									Type: schema.TypeFloat,
								},
								"volume_type": {
									Type: schema.TypeString,
								},
							},
						},
					},
					"no_device": {
						Type: schema.TypeBool,
					},
					"virtual_name": {
						Type: schema.TypeString,
					},
				},
			},
			Optional: true,
		},

		"client_token": {
			Type:     schema.TypeString,
			Required: true,
		},
		"disable_api_termination": {
			Type:     schema.TypeBool,
			Required: true,
		},
		"ebs_optimized": {
			Type:     schema.TypeBool,
			Required: true,
		},
		"image_id": {
			Type:     schema.TypeInt,
			Required: true,
		},
		"instance_initiated_shutdown_behavior": {
			Type:     schema.TypeBool,
			Required: true,
		},
		"instance_type": {
			Type:     schema.TypeString,
			Required: true,
		},
		"key_name": {
			Type:     schema.TypeString,
			Required: true,
		},
		"max_count": {
			Type:     schema.TypeInt,
			Optional: true,
		},
		"min_count": {
			Type:     schema.TypeInt,
			Required: true,
		},
		"network_interface": {

			Type:     schema.TypeSet,
			Optional: true,
			Elem: schema.Resource{
				Schema: map[string]*schema.Schema{
					"delete_on_termination": {
						Type: schema.TypeBool,
					},
					"description": {
						Type: schema.TypeString,
					},
					"device_index": {
						Type:     schema.TypeInt,
						Required: true,
					},
					"network_interface_id": {
						Type:     schema.TypeInt,
						Required: true,
					},
					"private_ip_address": {
						Type:     schema.TypeString,
						Optional: true,
					},
					"private_ip_addresses_set": {
						Type: schema.TypeSet,
						Elem: schema.Resource{
							Schema: map[string]*schema.Schema{
								"primary": {
									Type: schema.TypeString,
								},
								"private_ip_address": {
									Type: schema.TypeString,
								},
							},
						},
					},
					"secondary_private_ip_address_count": {
						Type: schema.TypeString,
					},
					"security_group_id": {
						Type:     schema.TypeSet,
						Optional: true,
						Elem:     &schema.Schema{Type: schema.TypeString},
					},
					"subnet_id": {
						Type:     schema.TypeString,
						Required: true,
					},
				},
			},
		},
		"placement": {
			Type:     schema.TypeSet,
			Optional: true,
			Elem: schema.Resource{
				Schema: map[string]*schema.Schema{
					"affinity": {
						Type: schema.TypeString,
					},
					"availability_zone": {
						Type: schema.TypeString,
					},
					"group_name": {
						Type: schema.TypeString,
					},
					"host_id": {
						Type: schema.TypeInt,
					},
					"tenancy": {
						Type: schema.TypeString,
					},
				},
			},
		},
		"private_ip_address": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"private_ip_addresses": {
			Type:     schema.TypeSet,
			Optional: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"security_group": {
			Type:     schema.TypeSet,
			Optional: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"security_group_id": {

			Type:     schema.TypeSet,
			Optional: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
		"subnet_id": {
			Type:     schema.TypeString,
			Required: true,
		},
		"user_data": {
			Type:     schema.TypeString,
			Optional: true,
		},
		//Attributes reference:
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
			Computed: true,
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
												Required: true,
											},
										},
									},
									Computed: true,
								},
							},
						},
						Computed: true,
						Required: true,
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
						Required: true,
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
						Required: true,
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
						Required: true,
					},
					"instance_id": {
						Type:     schema.TypeString,
						Computed: true,
						Required: true,
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
									Required: true,
								},
							},
						},
						Computed: true,
					},
					"instance_type": {
						Type:     schema.TypeString,
						Computed: true,
						Required: true,
					},
					"ip_address": {
						Type:     schema.TypeString,
						Computed: true,
						Required: true,
					},
					"kernel_id": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"key_name": {
						Type:     schema.TypeString,
						Required: true,
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
												Required: true,
											},
											"delete_on_termination": {
												Type:     schema.TypeBool,
												Computed: true,
											},
											"device_index": {
												Type:     schema.TypeInt,
												Computed: true,
												Required: true,
											},
											"status": {
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
									Required: true,
								},
								"mac_address": {
									Type:     schema.TypeString,
									Computed: true,
								},
								"network_interface_id": {
									Type:     schema.TypeString,
									Computed: true,
									Required: true,
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
									Required: true,
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
									Computed: true,
								},
								"source_dest_check": {
									Type:     schema.TypeBool,
									Computed: true,
									Required: true,
								},
								"status": {
									Type:     schema.TypeString,
									Computed: true,
								},
								"subnet_id": {
									Type:     schema.TypeString,
									Computed: true,
									Required: true,
								},
								"vpc_id": {
									Type:     schema.TypeInt,
									Computed: true,
									Required: true,
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
						Computed: true,
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
						Type: schema.TypeString,

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
									Required: true,
								},
							},
						},
						Computed: true,
					},
					"ramdisk_id": {
						Type:     schema.TypeString,
						Computed: true,
						Required: true,
					},
					"reason": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"root_device_name": {
						Type:     schema.TypeString,
						Computed: true,
						Required: true,
					},
					"source_dest_check": {
						Type:     schema.TypeString,
						Computed: true,
						Required: true,
					},
					"spot_instance_request_id": {
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
						Required: true,
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
						Required: true,
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
		//instance set is closed here
	}
}
