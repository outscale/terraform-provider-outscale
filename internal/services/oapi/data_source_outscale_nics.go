package oapi

import (
	"context"
	"fmt"
	"time"

	"github.com/outscale/osc-sdk-go/v2"
	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/services/oapi/oapihelpers"
	"github.com/outscale/terraform-provider-outscale/internal/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Creates a network interface in the specified subnet
func DataSourceOutscaleNics() *schema.Resource {
	return &schema.Resource{
		Read:   DataSourceOutscaleNicsRead,
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
						Type:     schema.TypeSet,
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
								"link_nic_id": {
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
						Type:     schema.TypeSet,
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
									Type:     schema.TypeSet,
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
						Computed: true,
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
					},
				},
			},
		},
		"request_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}

// Read Nic
func DataSourceOutscaleNicsRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*client.OutscaleClient).OSCAPI

	filters, filtersOk := d.GetOk("filter")

	var err error
	params := oscgo.ReadNicsRequest{}
	if filtersOk {
		params.Filters, err = buildOutscaleDataSourceNicFilters(filters.(*schema.Set))
		if err != nil {
			return err
		}
	}

	var resp oscgo.ReadNicsResponse
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		rp, httpResp, err := conn.NicApi.ReadNics(context.Background()).ReadNicsRequest(params).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		resp = rp
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
		d.SetId(id.UniqueId())

		if err := set("nics", getVMNetworkInterfaceSet(nics)); err != nil {
			return err
		}
		return nil
	})
}

func getVMNetworkInterfaceSet(nics []osc.Nic) (res []map[string]interface{}) {
	for _, nic := range nics {
		securityGroups, _ := oapihelpers.GetSecurityGroups(*nic.SecurityGroups)
		r := map[string]interface{}{
			"account_id":             nic.GetAccountId(),
			"description":            nic.GetDescription(),
			"is_source_dest_checked": nic.GetIsSourceDestChecked(),
			"mac_address":            nic.GetMacAddress(),
			"net_id":                 nic.GetNetId(),
			"nic_id":                 nic.GetNicId(),
			"private_dns_name":       nic.GetPrivateDnsName(),
			"private_ips":            oapihelpers.GetPrivateIPsForNic(nic.GetPrivateIps()),
			"security_groups":        securityGroups,
			"state":                  nic.GetState(),
			"subnet_id":              nic.GetSubnetId(),
			"subregion_name":         nic.GetSubregionName(),
			"tags":                   FlattenOAPITagsSDK(nic.GetTags()),
		}
		if _, ok := nic.GetLinkNicOk(); ok {
			r["link_nic"] = oapihelpers.GetOAPILinkNic(nic.GetLinkNic())
		}
		if _, ok := nic.GetLinkPublicIpOk(); ok {
			r["link_public_ip"] = oapihelpers.GetOAPILinkPublicIPsForNic(nic.GetLinkPublicIp())
		}
		res = append(res, r)
	}

	return
}
