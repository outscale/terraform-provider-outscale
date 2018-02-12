package outscale

import (
	"github.com/hashicorp/terraform/helper/schema"
)

func datasourceOutscaleOApiVMS() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOutscaleVMSRead,

		Schema: datasourceOutscaleOApiVMSSchema(),
	}
}

func datasourceOutscaleOApiVMSSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"filter": dataSourceFiltersSchema(),
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
	return nil
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
