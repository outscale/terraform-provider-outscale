package outscale

import (
	"context"
	"fmt"
	"time"

	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/terraform-providers/terraform-provider-outscale/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceOutscaleDHCPOptions() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOutscaleDHCPOptionsRead,

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
						"tags": dataSourceTagsSchema(),
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

func dataSourceOutscaleDHCPOptionsRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	filters, filtersOk := d.GetOk("filter")
	dhcpIDs, dhcpIDOk := d.GetOk("dhcp_options_set_ids")
	if !dhcpIDOk && !filtersOk {
		return fmt.Errorf("One of filters, or dhcp_options_set_id must be provided")
	}

	params := oscgo.ReadDhcpOptionsRequest{}
	if dhcpIDOk {
		params.Filters = &oscgo.FiltersDhcpOptions{
			DhcpOptionsSetIds: utils.InterfaceSliceToStringList(dhcpIDs.([]interface{})),
		}
	}
	if filtersOk {
		params.Filters = buildOutscaleDataSourceDHCPOptionFilters(filters.(*schema.Set))
	}

	var resp oscgo.ReadDhcpOptionsResponse
	var err error
	err = resource.Retry(120*time.Second, func() *resource.RetryError {
		rp, httpResp, err := conn.DhcpOptionApi.ReadDhcpOptions(context.Background()).ReadDhcpOptionsRequest(params).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		resp = rp
		return nil
	})
	if err != nil {
		return err
	}
	if len(resp.GetDhcpOptionsSets()) == 0 {
		return fmt.Errorf("Unable to find DHCP Option")
	}

	if err := d.Set("dhcp_options", flattenDHCPOption(resp.GetDhcpOptionsSets())); err != nil {
		return err
	}

	d.SetId(resource.UniqueId())

	return nil
}

func flattenDHCPOption(dhcpOptions []oscgo.DhcpOptionsSet) []map[string]interface{} {
	dhcpOptionsMap := make([]map[string]interface{}, len(dhcpOptions))

	for i, dhcpOption := range dhcpOptions {
		dhcpOptionsMap[i] = map[string]interface{}{
			"domain_name":         dhcpOption.GetDomainName(),
			"domain_name_servers": dhcpOption.GetDomainNameServers(),
			"log_servers":         dhcpOption.GetLogServers(),
			"ntp_servers":         dhcpOption.GetNtpServers(),
			"default":             dhcpOption.GetDefault(),
			"dhcp_options_set_id": dhcpOption.GetDhcpOptionsSetId(),
			"tags":                tagsOSCAPIToMap(dhcpOption.GetTags()),
		}
	}
	return dhcpOptionsMap
}
