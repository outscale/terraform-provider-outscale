package outscale

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	oscgo "github.com/outscale/osc-sdk-go/v2"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/spf13/cast"
	"github.com/terraform-providers/terraform-provider-outscale/utils"
)

// Creates a network interface in the specified subnet
func dataSourceOutscaleOAPINic() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOutscaleOAPINicRead,

		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			// This is attribute part for schema Nic
			// Argument
			"nic_id": {
				Type:     schema.TypeString,
				Optional: true,
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
				Type:     schema.TypeMap,
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
					},
				},
			},
			"link_nic": {
				Type:     schema.TypeMap,
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
							Type:     schema.TypeMap,
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
			"tags": dataSourceTagsSchema(),
			"net_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

// Read Nic
func dataSourceOutscaleOAPINicRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	nicID, okID := d.GetOk("nic_id")
	filters, okFilters := d.GetOk("filter")

	if okID && okFilters {
		return errors.New("nic_id and filter set")
	}

	dnri := oscgo.ReadNicsRequest{}

	if okID {
		dnri.Filters = &oscgo.FiltersNic{
			NicIds: &[]string{nicID.(string)},
		}
	}

	if okFilters {
		dnri.SetFilters(buildOutscaleOAPIDataSourceNicFilters(filters.(*schema.Set)))
	}

	var resp oscgo.ReadNicsResponse
	var err error
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, _, err = conn.NicApi.ReadNics(context.Background()).ReadNicsRequest(dnri).Execute()
		if err != nil {
			return utils.CheckThrottling(err)
		}
		return nil
	})

	if err != nil {

		return fmt.Errorf("Error describing Network Interfaces : %s", err)
	}

	if err != nil {
		if strings.Contains(err.Error(), utils.ResourceNotFound) {
			// The ENI is gone now, so just remove it from the state
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error retrieving ENI: %s", err)
	}

	if err := utils.IsResponseEmptyOrMutiple(len(resp.GetNics()), "Nic"); err != nil {
		return err
	}

	eni := resp.GetNics()[0]

	if err := d.Set("description", eni.GetDescription()); err != nil {
		return err
	}
	if err := d.Set("nic_id", eni.GetNicId()); err != nil {
		return err
	}
	if err := d.Set("subregion_name", eni.GetSubregionName()); err != nil {
		return err
	}
	if err := d.Set("subnet_id", eni.GetSubnetId()); err != nil {
		return err
	}

	b := make(map[string]interface{})

	link := eni.GetLinkPublicIp()
	b["public_ip_id"] = link.GetPublicIpId()
	b["link_public_ip_id"] = link.GetLinkPublicIpId()
	b["public_ip_account_id"] = link.GetPublicIpAccountId()
	b["public_dns_name"] = link.GetPublicDnsName()
	b["public_ip"] = link.GetPublicIp()

	if err := d.Set("link_public_ip", b); err != nil {
		return err
	}

	bb := make(map[string]interface{})

	linkNic := eni.GetLinkNic()

	bb["link_nic_id"] = linkNic.GetLinkNicId()
	bb["delete_on_vm_deletion"] = fmt.Sprintf("%t", linkNic.GetDeleteOnVmDeletion())
	bb["device_number"] = fmt.Sprintf("%d", linkNic.GetDeviceNumber())
	bb["vm_id"] = linkNic.GetVmId()
	bb["vm_account_id"] = linkNic.GetVmAccountId()
	bb["state"] = linkNic.GetState()

	if err := d.Set("link_nic", bb); err != nil {
		return err
	}

	if err := d.Set("subregion_name", eni.GetSubregionName()); err != nil {
		return err
	}

	x := make([]map[string]interface{}, len(eni.GetSecurityGroups()))
	for k, v := range eni.GetSecurityGroups() {
		b := make(map[string]interface{})
		b["security_group_id"] = v.GetSecurityGroupId()
		b["security_group_name"] = v.GetSecurityGroupName()
		x[k] = b
	}

	if err := d.Set("security_groups", x); err != nil {
		return err
	}

	if err := d.Set("mac_address", eni.GetMacAddress()); err != nil {
		return err
	}
	if err := d.Set("nic_id", eni.GetNicId()); err != nil {
		return err
	}
	if err := d.Set("account_id", eni.GetAccountId()); err != nil {
		return err
	}
	if err := d.Set("private_dns_name", eni.GetPrivateDnsName()); err != nil {
		return err
	}

	y := make([]map[string]interface{}, len(eni.GetPrivateIps()))
	if eni.PrivateIps != nil {
		for k, v := range eni.GetPrivateIps() {
			b := make(map[string]interface{})

			d := make(map[string]interface{})
			assoc := v.GetLinkPublicIp()
			d["public_ip_id"] = assoc.GetPublicIpId()
			d["link_public_ip_id"] = assoc.GetLinkPublicIpId()
			d["public_ip_account_id"] = assoc.GetPublicIpAccountId()
			d["public_dns_name"] = assoc.GetPublicDnsName()
			d["public_ip"] = assoc.GetPublicIp()

			b["link_public_ip"] = d
			b["private_dns_name"] = v.GetPrivateDnsName()
			b["private_ip"] = v.GetPrivateIp()
			b["is_primary"] = v.GetIsPrimary()

			y[k] = b
		}
	}
	if err := d.Set("private_ips", y); err != nil {
		return err
	}

	if err := d.Set("is_source_dest_checked", eni.GetIsSourceDestChecked()); err != nil {
		return err
	}
	if err := d.Set("state", eni.GetState()); err != nil {
		return err
	}
	if err := d.Set("tags", tagsOSCAPIToMap(eni.GetTags())); err != nil {
		return err
	}
	if err := d.Set("net_id", eni.GetNetId()); err != nil {
		return err
	}

	d.SetId(eni.GetNicId())
	return nil
}

func buildOutscaleOAPIDataSourceNicFilters(set *schema.Set) oscgo.FiltersNic {
	var filters oscgo.FiltersNic
	for _, v := range set.List() {
		m := v.(map[string]interface{})
		var filterValues []string
		for _, e := range m["values"].([]interface{}) {
			filterValues = append(filterValues, e.(string))
		}

		switch name := m["name"].(string); name {
		case "descriptions":
			filters.SetDescriptions(filterValues)
		case "is_source_dest_check":
			filters.SetIsSourceDestCheck(cast.ToBool(filterValues[0]))
		case "link_nic_delete_on_vm_deletion":
			filters.SetLinkNicDeleteOnVmDeletion(cast.ToBool(filterValues[0]))
		case "link_nic_device_numbers":
			filters.SetLinkNicDeviceNumbers(utils.StringSliceToInt32Slice(filterValues))
		case "link_nic_link_nic_ids":
			filters.SetLinkNicLinkNicIds(filterValues)
		case "link_nic_states":
			filters.SetLinkNicStates(filterValues)
		case "link_nic_vm_account_ids":
			filters.SetLinkNicVmAccountIds(filterValues)
		case "link_nic_vm_ids":
			filters.SetLinkNicVmIds(filterValues)
		case "link_public_ip_account_ids":
			filters.SetLinkPublicIpAccountIds(filterValues)
		case "link_public_ip_link_public_ip_ids":
			filters.SetLinkPublicIpLinkPublicIpIds(filterValues)
		case "link_public_ip_public_ip_ids":
			filters.SetLinkPublicIpPublicIpIds(filterValues)
		case "link_public_ip_public_ips":
			filters.SetLinkPublicIpPublicIps(filterValues)
		case "mac_addresses":
			filters.SetMacAddresses(filterValues)
		case "private_ips_primary_ip":
			filters.SetPrivateIpsPrimaryIp(cast.ToBool(filterValues[0]))
		case "tag_keys":
			filters.SetTagKeys(filterValues)
		case "tag_values":
			filters.SetTagValues(filterValues)
		case "tags":
			filters.SetTags(filterValues)
		case "net_ids":
			filters.SetNetIds(filterValues)
		case "nic_ids":
			filters.SetNicIds(filterValues)
		case "private_dns_names":
			filters.SetPrivateDnsNames(filterValues)
		case "private_ips_link_public_ip_account_ids":
			filters.SetPrivateIpsLinkPublicIpAccountIds(filterValues)
		case "private_ips_link_public_ip_public_ips":
			filters.SetPrivateIpsLinkPublicIpPublicIps(filterValues)
		case "private_ips_private_ips":
			filters.SetPrivateIpsPrivateIps(filterValues)
		case "security_group_ids":
			filters.SetSecurityGroupIds(filterValues)
		case "security_group_names":
			filters.SetSecurityGroupNames(filterValues)
		case "states":
			filters.SetStates(filterValues)
		case "subnet_ids":
			filters.SetSubnetIds(filterValues)
		case "subregion_names":
			filters.SetSubregionNames(filterValues)
		default:
			log.Printf("[Debug] Unknown Filter Name: %s.", name)
		}
	}
	return filters
}
