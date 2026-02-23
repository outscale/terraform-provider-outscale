package oapi

import (
	"context"
	"time"

	"github.com/outscale/osc-sdk-go/v3/pkg/options"
	"github.com/outscale/osc-sdk-go/v3/pkg/osc"
	"github.com/samber/lo"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/services/oapi/oapihelpers"
	"github.com/outscale/terraform-provider-outscale/internal/utils"
	"github.com/spf13/cast"
)

// Creates a network interface in the specified subnet
func DataSourceOutscaleNic() *schema.Resource {
	return &schema.Resource{
		ReadContext: DataSourceOutscaleNicRead,

		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			// This is attribute part for schema Nic
			// Argument
			"nic_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			// Attributes
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"private_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"security_group_id": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"subnet_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			// Attributes
			"link_public_ip": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"public_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"public_ip_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"link_public_ip_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"public_ip_account_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"public_dns_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"link_nic": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"link_nic_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"delete_on_vm_deletion": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"device_number": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"vm_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vm_account_id": {
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
			"subregion_name": {
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

			"mac_address": {
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

			"private_ips": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"link_public_ip": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"public_ip_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"link_public_ip_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
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
						"private_dns_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"private_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"is_primary": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},
			"request_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"requester_managed": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"is_source_dest_checked": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": TagsSchemaComputedSDK(),
			"net_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

// Read Nic
func DataSourceOutscaleNicRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.OutscaleClient).OSC

	filters, okFilters := d.GetOk("filter")

	if !okFilters {
		return diag.Errorf("filters must be assigned")
	}

	var err error
	dnri := osc.ReadNicsRequest{}
	if okFilters {
		dnri.Filters, err = buildOutscaleDataSourceNicFilters(filters.(*schema.Set))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	resp, err := client.ReadNics(ctx, dnri, options.WithRetryTimeout(5*time.Minute))
	if err != nil {
		return diag.Errorf("error describing network interfaces: %s", err)
	}
	if resp.Nics == nil || len(*resp.Nics) == 0 {
		return diag.FromErr(ErrNoResults)
	}
	if len(*resp.Nics) > 1 {
		return diag.FromErr(ErrMultipleResults)
	}

	eni := (*resp.Nics)[0]

	if err := d.Set("description", eni.Description); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("nic_id", eni.NicId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("subregion_name", eni.SubregionName); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("subnet_id", eni.SubnetId); err != nil {
		return diag.FromErr(err)
	}
	if eni.LinkPublicIp != nil {
		if err := d.Set("link_public_ip", flattenLinkPublicIp(eni.LinkPublicIp)); err != nil {
			return diag.FromErr(err)
		}
	}

	if eni.LinkNic != nil {
		if err := d.Set("link_nic", flattenLinkNic(eni.LinkNic)); err != nil {
			return diag.FromErr(err)
		}
	}

	if err := d.Set("subregion_name", eni.SubregionName); err != nil {
		return diag.FromErr(err)
	}

	x := make([]map[string]interface{}, len(eni.SecurityGroups))
	for k, v := range eni.SecurityGroups {
		b := make(map[string]interface{})
		b["security_group_id"] = v.SecurityGroupId
		b["security_group_name"] = v.SecurityGroupName
		x[k] = b
	}
	if err := d.Set("security_groups", x); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("mac_address", eni.MacAddress); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("nic_id", eni.NicId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("account_id", eni.AccountId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("private_dns_name", eni.PrivateDnsName); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("private_ips", oapihelpers.GetPrivateIPsForNic(eni.PrivateIps)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("is_source_dest_checked", eni.IsSourceDestChecked); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("state", eni.State); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("tags", FlattenOAPITagsSDK(eni.Tags)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("net_id", eni.NetId); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(eni.NicId)
	return nil
}

func buildOutscaleDataSourceNicFilters(set *schema.Set) (*osc.FiltersNic, error) {
	var filters osc.FiltersNic
	for _, v := range set.List() {
		m := v.(map[string]interface{})
		var filterValues []string
		for _, e := range m["values"].([]interface{}) {
			filterValues = append(filterValues, e.(string))
		}

		switch name := m["name"].(string); name {
		case "descriptions":
			filters.Descriptions = &filterValues
		case "is_source_dest_check":
			filters.IsSourceDestCheck = new(cast.ToBool(filterValues[0]))
		case "link_nic_delete_on_vm_deletion":
			filters.LinkNicDeleteOnVmDeletion = new(cast.ToBool(filterValues[0]))
		case "link_nic_device_numbers":
			filters.LinkNicDeviceNumbers = new(utils.StringSliceToIntSlice(filterValues))
		case "link_nic_link_nic_ids":
			filters.LinkNicLinkNicIds = &filterValues
		case "link_nic_states":
			filters.LinkNicStates = &filterValues
		case "link_nic_vm_account_ids":
			filters.LinkNicVmAccountIds = &filterValues
		case "link_nic_vm_ids":
			filters.LinkNicVmIds = &filterValues
		case "link_public_ip_account_ids":
			filters.LinkPublicIpAccountIds = &filterValues
		case "link_public_ip_link_public_ip_ids":
			filters.LinkPublicIpLinkPublicIpIds = &filterValues
		case "link_public_ip_public_ip_ids":
			filters.LinkPublicIpPublicIpIds = &filterValues
		case "link_public_ip_public_ips":
			filters.LinkPublicIpPublicIps = &filterValues
		case "mac_addresses":
			filters.MacAddresses = &filterValues
		case "private_ips_primary_ip":
			filters.PrivateIpsPrimaryIp = new(cast.ToBool(filterValues[0]))
		case "tag_keys":
			filters.TagKeys = &filterValues
		case "tag_values":
			filters.TagValues = &filterValues
		case "tags":
			filters.Tags = &filterValues
		case "net_ids":
			filters.NetIds = &filterValues
		case "nic_ids":
			filters.NicIds = &filterValues
		case "private_dns_names":
			filters.PrivateDnsNames = &filterValues
		case "private_ips_link_public_ip_account_ids":
			filters.PrivateIpsLinkPublicIpAccountIds = &filterValues
		case "private_ips_link_public_ip_public_ips":
			filters.PrivateIpsLinkPublicIpPublicIps = &filterValues
		case "private_ips_private_ips":
			filters.PrivateIpsPrivateIps = &filterValues
		case "security_group_ids":
			filters.SecurityGroupIds = &filterValues
		case "security_group_names":
			filters.SecurityGroupNames = &filterValues
		case "states":
			filters.States = new(lo.Map(filterValues, func(s string, _ int) osc.NicState { return osc.NicState(s) }))
		case "subnet_ids":
			filters.SubnetIds = &filterValues
		case "subregion_names":
			filters.SubregionNames = &filterValues
		default:
			return nil, utils.UnknownDataSourceFilterError(name)
		}
	}
	return &filters, nil
}
