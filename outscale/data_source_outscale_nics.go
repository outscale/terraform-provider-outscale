package outscale

import (
	"context"
	"fmt"
	"time"

	oscgo "github.com/outscale/osc-sdk-go/v2"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-outscale/utils"
)

// Creates a network interface in the specified subnet
func dataSourceOutscaleOAPINics() *schema.Resource {
	return &schema.Resource{
		Read:   dataSourceOutscaleOAPINicsRead,
		Schema: getDSOAPINicsSchema(),
	}
}

func getDSOAPINicsSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		//  This is attribute part for schema Nic
		"filter": dataSourceFiltersSchema(),
		"nics": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"account_id": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"description": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"is_source_dest_checked": {
						Type:     schema.TypeBool,
						Computed: true,
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
									Type:     schema.TypeInt,
									Computed: true,
								},
								"nic_link_id": {
									Type:     schema.TypeString,
									Computed: true,
								},
								"state": {
									Type:     schema.TypeString,
									Computed: true,
								},
								"vm_account_id": {
									Type:     schema.TypeString,
									Computed: true,
								},
								"vm_id": {
									Type:     schema.TypeString,
									Computed: true,
								},
							},
						},
					},
					"link_public_ip": {
						Type:     schema.TypeMap,
						Computed: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"link_public_ip_id": {
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
								"public_ip_account_id": {
									Type:     schema.TypeString,
									Computed: true,
								},
								"public_ip_id": {
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
					"nic_id": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"private_dns_name": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"private_ips": {
						Type:     schema.TypeList,
						Computed: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"is_primary": {
									Type:     schema.TypeBool,
									Computed: true,
								},
								"link_public_ip": {
									Type:     schema.TypeMap,
									Computed: true,
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"link_public_ip_id": {
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
											"public_ip_account_id": {
												Type:     schema.TypeString,
												Computed: true,
											},
											"public_ip_id": {
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
									Computed: true,
								},
							},
						},
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
					"subnet_id": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"subregion_name": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"tags": {
						Type:     schema.TypeList,
						Optional: true,
						Computed: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"key": {
									Type:     schema.TypeString,
									Optional: true,
									Computed: true,
								},
								"value": {
									Type:     schema.TypeString,
									Optional: true,
									Computed: true,
								},
							},
						},
					}},
			},
		},
		"request_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}

// Read Nic
func dataSourceOutscaleOAPINicsRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	filters, filtersOk := d.GetOk("filter")
	if !filtersOk {
		return fmt.Errorf("filters, or owner must be assigned, or nic_id must be provided")
	}

	params := oscgo.ReadNicsRequest{}
	if filtersOk {
		params.SetFilters(buildOutscaleOAPIDataSourceNicFilters(filters.(*schema.Set)))
	}

	var resp oscgo.ReadNicsResponse
	var err error

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, _, err = conn.NicApi.ReadNics(context.Background()).ReadNicsRequest(params).Execute()
		if err != nil {
			return utils.CheckThrottling(err)
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("Error reading Network Interface Cards : %s", err)
	}

	if resp.GetNics() == nil {
		return fmt.Errorf("Your query returned no results. Please change your search criteria and try again")
	}

	if len(resp.GetNics()) == 0 {
		return fmt.Errorf("Your query returned no results. Please change your search criteria and try again")
	}

	nics := resp.GetNics()

	return resourceDataAttrSetter(d, func(set AttributeSetter) error {
		d.SetId(resource.UniqueId())

		if err := set("nics", getOAPIVMNetworkInterfaceSet(nics)); err != nil {
			return err
		}

		return nil
	})
}
