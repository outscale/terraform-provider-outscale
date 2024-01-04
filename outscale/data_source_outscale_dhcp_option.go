package outscale

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/spf13/cast"
	"github.com/terraform-providers/terraform-provider-outscale/utils"
)

func dataSourceOutscaleDHCPOption() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOutscaleDHCPOptionRead,

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
			"tags": dataSourceTagsSchema(),
			"request_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceOutscaleDHCPOptionRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	filters, filtersOk := d.GetOk("filter")
	dhcpID, dhcpIDOk := d.GetOk("dhcp_options_set_id")
	if !dhcpIDOk && !filtersOk {
		return fmt.Errorf("One of filters, or dhcp_options_set_id must be provided")
	}

	params := oscgo.ReadDhcpOptionsRequest{}
	if dhcpIDOk {
		params.Filters = &oscgo.FiltersDhcpOptions{
			DhcpOptionsSetIds: &[]string{dhcpID.(string)},
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

	if len(resp.GetDhcpOptionsSets()) > 1 {
		return fmt.Errorf("multiple results returned, please use a more specific criteria in your query")
	}

	dhcpOption := resp.GetDhcpOptionsSets()[0]

	if err := d.Set("domain_name", dhcpOption.GetDomainName()); err != nil {
		return err
	}
	if err := d.Set("domain_name_servers", dhcpOption.GetDomainNameServers()); err != nil {
		return err
	}
	if err := d.Set("log_servers", dhcpOption.GetLogServers()); err != nil {
		return err
	}
	if err := d.Set("ntp_servers", dhcpOption.GetNtpServers()); err != nil {
		return err
	}
	if err := d.Set("default", dhcpOption.GetDefault()); err != nil {
		return err
	}
	if err := d.Set("dhcp_options_set_id", dhcpOption.GetDhcpOptionsSetId()); err != nil {
		return err
	}
	if err := d.Set("tags", tagsOSCAPIToMap(dhcpOption.GetTags())); err != nil {
		return err
	}

	d.SetId(dhcpOption.GetDhcpOptionsSetId())

	return nil
}

func buildOutscaleDataSourceDHCPOptionFilters(set *schema.Set) *oscgo.FiltersDhcpOptions {
	var filters oscgo.FiltersDhcpOptions
	for _, v := range set.List() {
		m := v.(map[string]interface{})
		var filterValues []string
		for _, e := range m["values"].([]interface{}) {
			filterValues = append(filterValues, e.(string))
		}

		switch name := m["name"].(string); name {
		case "dhcp_options_set_ids":
			filters.SetDhcpOptionsSetIds(filterValues)
		case "dhcp_options_set_id":
			filters.SetDhcpOptionsSetIds(filterValues)
		case "domain_name_servers":
			filters.SetDomainNameServers(filterValues)
		case "domain_names":
			filters.SetDomainNames(filterValues)
		case "log_servers":
			filters.SetLogServers(filterValues)
		case "ntp_servers":
			filters.SetNtpServers(filterValues)
		case "tag_keys":
			filters.SetTagKeys(filterValues)
		case "tag_values":
			filters.SetTagValues(filterValues)
		case "tags":
			filters.SetTags(filterValues)
		case "default":
			filters.SetDefault(cast.ToBool(filterValues[0]))
		default:
			log.Printf("[Debug] Unknown Filter Name: %s.", name)
		}
	}
	return &filters
}
