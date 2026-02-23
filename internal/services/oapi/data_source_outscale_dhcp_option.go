package oapi

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/outscale/goutils/sdk/ptr"
	"github.com/outscale/osc-sdk-go/v3/pkg/options"
	"github.com/outscale/osc-sdk-go/v3/pkg/osc"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/utils"
	"github.com/spf13/cast"
)

func DataSourceOutscaleDHCPOption() *schema.Resource {
	return &schema.Resource{
		ReadContext: DataSourceOutscaleDHCPOptionRead,

		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			"dhcp_options_set_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"domain_name": {
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
			"request_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func DataSourceOutscaleDHCPOptionRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.OutscaleClient).OSC

	filters, filtersOk := d.GetOk("filter")
	dhcpID, dhcpIDOk := d.GetOk("dhcp_options_set_id")
	if !dhcpIDOk && !filtersOk {
		return diag.Errorf("one of filters, or dhcp_options_set_id must be provided")
	}

	params := osc.ReadDhcpOptionsRequest{}
	if dhcpIDOk {
		params.Filters = &osc.FiltersDhcpOptions{
			DhcpOptionsSetIds: &[]string{dhcpID.(string)},
		}
	}
	if filtersOk {
		filterParams, err := buildOutscaleDataSourceDHCPOptionFilters(filters.(*schema.Set))
		if err != nil {
			return diag.FromErr(err)
		}
		params.Filters = filterParams
	}

	resp, err := client.ReadDhcpOptions(ctx, params, options.WithRetryTimeout(120*time.Second))
	if err != nil {
		return diag.FromErr(err)
	}

	if resp.DhcpOptionsSets == nil || len(*resp.DhcpOptionsSets) == 0 {
		return diag.FromErr(ErrNoResults)
	}

	if len(*resp.DhcpOptionsSets) > 1 {
		return diag.FromErr(ErrMultipleResults)
	}

	dhcpOption := (*resp.DhcpOptionsSets)[0]

	if err := d.Set("domain_name", ptr.From(dhcpOption.DomainName)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("domain_name_servers", ptr.From(dhcpOption.DomainNameServers)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("log_servers", ptr.From(dhcpOption.LogServers)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("ntp_servers", ptr.From(dhcpOption.NtpServers)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("default", ptr.From(dhcpOption.Default)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("dhcp_options_set_id", ptr.From(dhcpOption.DhcpOptionsSetId)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("tags", FlattenOAPITagsSDK(ptr.From(dhcpOption.Tags))); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(ptr.From(dhcpOption.DhcpOptionsSetId))

	return nil
}

func buildOutscaleDataSourceDHCPOptionFilters(set *schema.Set) (*osc.FiltersDhcpOptions, error) {
	var filters osc.FiltersDhcpOptions
	for _, v := range set.List() {
		m := v.(map[string]interface{})
		var filterValues []string
		for _, e := range m["values"].([]interface{}) {
			filterValues = append(filterValues, e.(string))
		}

		switch name := m["name"].(string); name {
		case "dhcp_options_set_ids":
			filters.DhcpOptionsSetIds = &filterValues
		case "dhcp_options_set_id":
			filters.DhcpOptionsSetIds = &filterValues
		case "domain_name_servers":
			filters.DomainNameServers = &filterValues
		case "domain_names":
			filters.DomainNames = &filterValues
		case "log_servers":
			filters.LogServers = &filterValues
		case "ntp_servers":
			filters.NtpServers = &filterValues
		case "tag_keys":
			filters.TagKeys = &filterValues
		case "tag_values":
			filters.TagValues = &filterValues
		case "tags":
			filters.Tags = &filterValues
		case "default":
			filters.Default = new(cast.ToBool(filterValues[0]))
		default:
			return nil, utils.UnknownDataSourceFilterError(name)
		}
	}
	return &filters, nil
}
