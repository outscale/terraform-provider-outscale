package oapi

import (
	"context"
	"time"

	"github.com/outscale/goutils/sdk/ptr"
	"github.com/outscale/osc-sdk-go/v3/pkg/options"
	"github.com/outscale/osc-sdk-go/v3/pkg/osc"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/services/oapi/oapihelpers"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Creates a network interface in the specified subnet
func DataSourceOutscaleNics() *schema.Resource {
	return &schema.Resource{
		ReadContext: DataSourceOutscaleNicsRead,
		Schema:      getDSOAPINicsSchema(),
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
func DataSourceOutscaleNicsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.OutscaleClient).OSC

	filters, filtersOk := d.GetOk("filter")

	var err error
	params := osc.ReadNicsRequest{}
	if filtersOk {
		params.Filters, err = buildOutscaleDataSourceNicFilters(filters.(*schema.Set))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	resp, err := client.ReadNics(ctx, params, options.WithRetryTimeout(5*time.Minute))
	if err != nil {
		return diag.Errorf("error reading network interface cards : %s", err)
	}

	if resp.Nics == nil {
		return diag.FromErr(ErrNoResults)
	}

	if resp.Nics == nil || len(*resp.Nics) == 0 {
		return diag.FromErr(ErrNoResults)
	}
	nics := ptr.From(resp.Nics)

	return diag.FromErr(resourceDataAttrSetter(d, func(set AttributeSetter) error {
		d.SetId(id.UniqueId())

		if err := set("nics", getVMNetworkInterfaceSet(nics)); err != nil {
			return err
		}
		return nil
	}))
}

func getVMNetworkInterfaceSet(nics []osc.Nic) (res []map[string]interface{}) {
	for _, nic := range nics {
		securityGroups, _ := oapihelpers.GetSecurityGroups(nic.SecurityGroups)
		r := map[string]interface{}{
			"account_id":             nic.AccountId,
			"description":            nic.Description,
			"is_source_dest_checked": nic.IsSourceDestChecked,
			"mac_address":            nic.MacAddress,
			"net_id":                 nic.NetId,
			"nic_id":                 nic.NicId,
			"private_dns_name":       nic.PrivateDnsName,
			"private_ips":            oapihelpers.GetPrivateIPsForNic(nic.PrivateIps),
			"security_groups":        securityGroups,
			"state":                  nic.State,
			"subnet_id":              nic.SubnetId,
			"subregion_name":         nic.SubregionName,
			"tags":                   FlattenOAPITagsSDK(nic.Tags),
		}
		if nic.LinkNic != nil {
			r["link_nic"] = oapihelpers.GetOAPILinkNic(*nic.LinkNic)
		}
		if nic.LinkPublicIp != nil {
			r["link_public_ip"] = oapihelpers.GetOAPILinkPublicIPsForNic(*nic.LinkPublicIp)
		}
		res = append(res, r)
	}

	return
}
