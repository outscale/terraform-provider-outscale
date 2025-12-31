package oapi

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	oscgo "github.com/outscale/osc-sdk-go/v2"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/services/oapi/oapihelpers"
	"github.com/outscale/terraform-provider-outscale/internal/utils"
	"github.com/spf13/cast"
)

// Creates a network interface in the specified subnet
func DataSourceOutscaleNic() *schema.Resource {
	return &schema.Resource{
		Read: DataSourceOutscaleNicRead,

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
func DataSourceOutscaleNicRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*client.OutscaleClient).OSCAPI

	filters, okFilters := d.GetOk("filter")

	if !okFilters {
		return errors.New("filters must be assigned")
	}

	var err error
	dnri := oscgo.ReadNicsRequest{}
	if okFilters {
		dnri.Filters, err = buildOutscaleDataSourceNicFilters(filters.(*schema.Set))
		if err != nil {
			return err
		}
	}

	var resp oscgo.ReadNicsResponse
	var statusCode int
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		rp, httpResp, err := conn.NicApi.ReadNics(context.Background()).ReadNicsRequest(dnri).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		resp = rp
		statusCode = httpResp.StatusCode
		return nil
	})
	if err != nil {
		if statusCode == http.StatusNotFound {
			// The ENI is gone now, so just remove it from the state
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error describing Network Interfaces: %s", err)
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
	if linkIp, ok := eni.GetLinkPublicIpOk(); ok {
		if err := d.Set("link_public_ip", flattenLinkPublicIp(linkIp)); err != nil {
			return err
		}
	}

	if linkNic, ok := eni.GetLinkNicOk(); ok {
		if err := d.Set("link_nic", flattenLinkNic(linkNic)); err != nil {
			return err
		}
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
	if privIps, ok := eni.GetPrivateIpsOk(); ok {
		if err := d.Set("private_ips", oapihelpers.GetPrivateIPsForNic(*privIps)); err != nil {
			return err
		}
	}
	if err := d.Set("is_source_dest_checked", eni.GetIsSourceDestChecked()); err != nil {
		return err
	}
	if err := d.Set("state", eni.GetState()); err != nil {
		return err
	}
	if err := d.Set("tags", FlattenOAPITagsSDK(eni.GetTags())); err != nil {
		return err
	}
	if err := d.Set("net_id", eni.GetNetId()); err != nil {
		return err
	}

	d.SetId(eni.GetNicId())
	return nil
}

func buildOutscaleDataSourceNicFilters(set *schema.Set) (*oscgo.FiltersNic, error) {
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
			return nil, utils.UnknownDataSourceFilterError(context.Background(), name)
		}
	}
	return &filters, nil
}
