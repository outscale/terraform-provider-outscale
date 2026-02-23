package oapi

import (
	"context"
	"time"

	"github.com/outscale/goutils/sdk/ptr"
	"github.com/outscale/osc-sdk-go/v3/pkg/options"
	"github.com/outscale/osc-sdk-go/v3/pkg/osc"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceOutscaleDHCPOptions() *schema.Resource {
	return &schema.Resource{
		ReadContext: DataSourceOutscaleDHCPOptionsRead,

		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			"dhcp_options_set_ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"dhcp_options": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"domain_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"dhcp_options_set_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"domain_name_servers": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"log_servers": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"ntp_servers": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"default": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"tags": TagsSchemaComputedSDK(),
					},
				},
			},
			"request_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func DataSourceOutscaleDHCPOptionsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.OutscaleClient).OSC

	filters, filtersOk := d.GetOk("filter")
	dhcpIDs, dhcpIDOk := d.GetOk("dhcp_options_set_ids")
	if !dhcpIDOk && !filtersOk {
		return diag.Errorf("one of filters, or dhcp_options_set_id must be provided")
	}

	var err error
	params := osc.ReadDhcpOptionsRequest{}
	if dhcpIDOk {
		params.Filters = &osc.FiltersDhcpOptions{
			DhcpOptionsSetIds: utils.InterfaceSliceToStringList(dhcpIDs.([]interface{})),
		}
	}
	if filtersOk {
		params.Filters, err = buildOutscaleDataSourceDHCPOptionFilters(filters.(*schema.Set))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	resp, err := client.ReadDhcpOptions(ctx, params, options.WithRetryTimeout(120*time.Second))
	if err != nil {
		return diag.FromErr(err)
	}
	if resp.DhcpOptionsSets == nil || len(*resp.DhcpOptionsSets) == 0 {
		return diag.FromErr(ErrNoResults)
	}

	if err := d.Set("dhcp_options", flattenDHCPOption(*resp.DhcpOptionsSets)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(id.UniqueId())

	return nil
}

func flattenDHCPOption(dhcpOptions []osc.DhcpOptionsSet) []map[string]interface{} {
	dhcpOptionsMap := make([]map[string]interface{}, len(dhcpOptions))

	for i, dhcpOption := range dhcpOptions {
		dhcpOptionsMap[i] = map[string]interface{}{
			"domain_name":         dhcpOption.DomainName,
			"domain_name_servers": dhcpOption.DomainNameServers,
			"log_servers":         dhcpOption.LogServers,
			"ntp_servers":         dhcpOption.NtpServers,
			"default":             dhcpOption.Default,
			"dhcp_options_set_id": dhcpOption.DhcpOptionsSetId,
			"tags":                FlattenOAPITagsSDK(ptr.From(dhcpOption.Tags)),
		}
	}
	return dhcpOptionsMap
}
