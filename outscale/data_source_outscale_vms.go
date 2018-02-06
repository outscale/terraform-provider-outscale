package outscale

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/hashicorp/terraform/helper/schema"
)

func datasourceOutscaleInstance() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOutscaleInstanceRead,

		Schema: datasourceOutscaleInstanceSchema(),
	}
}

func datasourceOutscaleInstanceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{

		// Optional attributes
		"filter": dataSourceFiltersSchema(),
		"ami_launch_index": &schema.Schema{
			Type: schema.TypeInt,

			Computed: true,
		},
		"architecture": &schema.Schema{
			Type: schema.TypeString,

			Computed: true,
		},
		"block_device_mapping": &schema.Schema{
			Type: schema.TypeList,

			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"device_name": &schema.Schema{
						Type:     schema.TypeString,
						Computed: true,
					},
					"ebs": &schema.Schema{
						Type:     schema.TypeSet,
						Computed: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"delete_on_termination": &schema.Schema{
									Type: schema.TypeBool,

									Computed: true,
								},
								"status": &schema.Schema{
									Type: schema.TypeString,

									Computed: true,
								},
								"volume_id": &schema.Schema{
									Type: schema.TypeString,

									Computed: true,
								},
							},
						},
					},
				},
			},
		},
		"client_token": &schema.Schema{
			Type: schema.TypeString,

			Computed: true,
		},
		"dns_name": &schema.Schema{
			Type: schema.TypeString,

			Computed: true,
		},
		"ebs_optimized": &schema.Schema{
			Type: schema.TypeBool,

			Computed: true,
		},
		"group_set": &schema.Schema{
			Type: schema.TypeList,

			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"arn": &schema.Schema{
						Type: schema.TypeString,

						Computed: true,
					},
					"id": &schema.Schema{
						Type: schema.TypeString,

						Computed: true,
					},
				},
			},
		},
		"hypervisor": &schema.Schema{
			Type: schema.TypeString,

			Computed: true,
		},
		"iam_instance_profile": &schema.Schema{
			Type: schema.TypeSet,

			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"arn": &schema.Schema{
						Type: schema.TypeString,

						Computed: true,
					},
					"id": &schema.Schema{
						Type: schema.TypeString,

						Computed: true,
					},
				},
			},
		},
		"image_id": &schema.Schema{
			Type: schema.TypeString,

			Computed: true,
		},
		"instance_id": &schema.Schema{
			Type:     schema.TypeSet,
			Optional: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"instance_lifecycle": &schema.Schema{
			Type: schema.TypeString,

			Computed: true,
		},
		"instance_state": &schema.Schema{
			Type: schema.TypeSet,

			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"code": &schema.Schema{
						Type: schema.TypeInt,

						Computed: true,
					},
					"name": &schema.Schema{
						Type: schema.TypeString,

						Computed: true,
					},
				},
			},
		},
		"instance_type": &schema.Schema{
			Type: schema.TypeString,

			Computed: true,
		},
		"ip_address": &schema.Schema{
			Type: schema.TypeString,

			Computed: true,
		},
		"kernel_id": &schema.Schema{
			Type: schema.TypeString,

			Computed: true,
		},
		"key_name": &schema.Schema{
			Type: schema.TypeString,

			Computed: true,
		},
		"monitoring": &schema.Schema{
			Type: schema.TypeSet,

			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"state": &schema.Schema{
						Type: schema.TypeString,

						Computed: true,
					},
				},
			},
		},
		"network_interface_set": &schema.Schema{
			Type: schema.TypeList,

			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"association": &schema.Schema{
						Type:     schema.TypeSet,
						Computed: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"ip_owner_id": &schema.Schema{
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
					"attachment": &schema.Schema{
						Type: schema.TypeSet,

						Computed: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"attachment_id": &schema.Schema{
									Type: schema.TypeString,

									Computed: true,
								},
								"delete_on_termination": &schema.Schema{
									Type: schema.TypeBool,

									Computed: true,
								},
								"device_index": &schema.Schema{
									Type: schema.TypeInt,

									Computed: true,
								},
								"status": &schema.Schema{
									Type: schema.TypeString,

									Computed: true,
								},
							},
						},
					},
					"description": &schema.Schema{
						Type: schema.TypeString,

						Computed: true,
					},
					"group_set": &schema.Schema{
						Type: schema.TypeList,

						Computed: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"group_id": &schema.Schema{
									Type: schema.TypeString,

									Computed: true,
								},
								"group_name": &schema.Schema{
									Type: schema.TypeString,

									Computed: true,
								},
							},
						},
					},
					"mac_address": &schema.Schema{
						Type: schema.TypeString,

						Computed: true,
					},
					"network_interface_id": &schema.Schema{
						Type: schema.TypeString,

						Computed: true,
					},
					"owner_id": &schema.Schema{
						Type: schema.TypeString,

						Computed: true,
					},
					"private_dns_name": &schema.Schema{
						Type: schema.TypeString,

						Computed: true,
					},
					"private_ip_address": &schema.Schema{
						Type: schema.TypeString,

						Computed: true,
					},
					"private_ip_address_set": &schema.Schema{
						Type: schema.TypeSet,

						Computed: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"association": &schema.Schema{
									Type:     schema.TypeSet,
									Computed: true,
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"ip_owner_id": &schema.Schema{
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
								"primary": &schema.Schema{
									Type: schema.TypeBool,

									Computed: true,
								},
								"private_dns": &schema.Schema{
									Type: schema.TypeString,

									Computed: true,
								},
								"private_ip_address": &schema.Schema{
									Type: schema.TypeString,

									Computed: true,
								},
							},
						},
					},
					"source_dest_check": &schema.Schema{
						Type: schema.TypeBool,

						Computed: true,
					},
					"status": &schema.Schema{
						Type: schema.TypeString,

						Computed: true,
					},
					"subnet_id": &schema.Schema{
						Type: schema.TypeString,

						Computed: true,
					},
					"vpc_id": &schema.Schema{
						Type: schema.TypeString,

						Computed: true,
					},
				},
			},
		},
		"placement": &schema.Schema{
			Type: schema.TypeSet,

			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"affinity": &schema.Schema{
						Type: schema.TypeString,

						Computed: true,
					},
					"availability_zone": &schema.Schema{
						Type: schema.TypeString,

						Computed: true,
					},
					"group_name": &schema.Schema{
						Type: schema.TypeString,

						Computed: true,
					},
					"host_id": &schema.Schema{
						Type: schema.TypeString,

						Computed: true,
					},
					"tenancy": &schema.Schema{
						Type: schema.TypeString,

						Computed: true,
					},
				},
			},
		},
		"platform": &schema.Schema{
			Type: schema.TypeString,

			Computed: true,
		},
		"private_dns_name": &schema.Schema{
			Type: schema.TypeString,

			Computed: true,
		},
		"product_codes": &schema.Schema{
			Type: schema.TypeList,

			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"product_code": &schema.Schema{
						Type: schema.TypeString,

						Computed: true,
					},
					"type": &schema.Schema{
						Type: schema.TypeString,

						Computed: true,
					},
				},
			},
		},
		"rundisk_id": &schema.Schema{
			Type: schema.TypeString,

			Computed: true,
		},
		"reason": &schema.Schema{
			Type: schema.TypeString,

			Computed: true,
		},
		"root_device_name": &schema.Schema{
			Type: schema.TypeString,

			Computed: true,
		},
		"root_device_type": &schema.Schema{
			Type: schema.TypeString,

			Computed: true,
		},
		"source_dest_check": &schema.Schema{
			Type: schema.TypeBool,

			Computed: true,
		},
		"spot_instance_request_id": &schema.Schema{
			Type: schema.TypeString,

			Computed: true,
		},
		"sriov_net_support": &schema.Schema{
			Type: schema.TypeString,

			Computed: true,
		},
		"state_reason": &schema.Schema{
			Type: schema.TypeSet,

			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"code": &schema.Schema{
						Type: schema.TypeString,

						Computed: true,
					},
					"message": &schema.Schema{
						Type: schema.TypeString,

						Computed: true,
					},
				},
			},
		},
		"subnet_id": &schema.Schema{
			Type: schema.TypeString,

			Computed: true,
		},
		"tag_set": &schema.Schema{
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"key": &schema.Schema{
						Type: schema.TypeString,

						Computed: true,
					},
					"value": &schema.Schema{
						Type: schema.TypeString,

						Computed: true,
					},
				},
			},
		},
		"virtualization_type": &schema.Schema{
			Type: schema.TypeString,

			Computed: true,
		},
		"vpc_id": &schema.Schema{
			Type: schema.TypeString,

			Computed: true,
		},
		// Computed attributes

		"allocation_id": &schema.Schema{
			Type: schema.TypeString,

			Computed: true,
		},
		"association_id": &schema.Schema{
			Type: schema.TypeString,

			Computed: true,
		},
		"network_interface_id": &schema.Schema{
			Type: schema.TypeString,

			Computed: true,
		},
		"network_interface_owner_id": &schema.Schema{
			Type: schema.TypeString,

			Computed: true,
		},
		"private_ip_address": &schema.Schema{
			Type: schema.TypeString,

			Computed: true,
		},
		"public_ip": &schema.Schema{
			Type: schema.TypeString,

			Computed: true,
		},
	}
}

func dataSourceOutscaleInstanceRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func buildOutscaleDataSourceFilters(set *schema.Set) []*Filter {
	var filters []*Filter
	for _, v := range set.List() {
		m := v.(map[string]interface{})
		var filterValues []*string
		for _, e := range m["values"].([]interface{}) {
			filterValues = append(filterValues, aws.String(e.(string)))
		}
		filters = append(filters, &Filter{
			Name:   String(m["name"].(string)),
			Values: filterValues,
		})
	}
	return filters
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
