package outscale

import (
	"context"
	"fmt"
	"github.com/antihax/optional"
	oscgo "github.com/marinsalinas/osc-sdk-go"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceOutscaleOAPIVpc() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOutscaleOAPIVpcRead,

		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			"ip_range": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"dhcp_options_set_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"net_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"tenancy": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"request_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"tags": tagsOAPIListSchemaComputed(),
		},
	}
}

func dataSourceOutscaleOAPIVpcRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	req := oscgo.ReadNetsRequest{}

	if v, ok := d.GetOk("filter"); ok {
		req.SetFilters(buildOutscaleOAPIDataSourceNetFilters(v.(*schema.Set)))
	}

	if id := d.Get("net_id"); id != "" {
		req.Filters.SetNetIds([]string{id.(string)})
	}

	var err error
	var resp oscgo.ReadNetsResponse
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, _, err = conn.NetApi.ReadNets(context.Background(), &oscgo.ReadNetsOpts{ReadNetsRequest: optional.NewInterface(req)})
		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded:") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return err
	}
	if len(resp.GetNets()) == 0 {
		return fmt.Errorf("No matching Net found")
	}
	if len(resp.GetNets()) > 1 {
		return fmt.Errorf("Multiple Nets matched; use additional constraints to reduce matches to a single Net")
	}

	net := resp.GetNets()[0]

	d.SetId(net.GetNetId())
	d.Set("net_id", net.GetNetId())
	d.Set("ip_range", net.GetIpRange())
	d.Set("dhcp_options_set_id", net.GetDhcpOptionsSetId())
	d.Set("tenancy", net.GetTenancy())
	d.Set("state", net.GetState())
	d.Set("request_id", resp.ResponseContext.GetRequestId())

	return d.Set("tags", tagsOSCAPIToMap(net.GetTags()))
}

func buildOutscaleOAPIDataSourceNetFilters(set *schema.Set) oscgo.FiltersNet {
	var filters oscgo.FiltersNet
	for _, v := range set.List() {
		m := v.(map[string]interface{})
		var filterValues []string
		for _, e := range m["values"].([]interface{}) {
			filterValues = append(filterValues, e.(string))
		}

		switch name := m["name"].(string); name {
		case "net_ids":
			filters.SetNetIds(filterValues)
		case "ip_range":
			filters.SetIpRanges(filterValues)
		case "dhcp_options_set_id":
			filters.SetDhcpOptionsSetIds(filterValues)
		case "is_default":
			//bool
			//filters.IsDefault = filterValues
		case "state":
			filters.SetStates(filterValues)
		case "tag_key":
			filters.SetTagKeys(filterValues)
		case "tag_value":
			filters.SetTagValues(filterValues)
		default:
			log.Printf("[Debug] Unknown Filter Name: %s.", name)
		}
	}
	return filters
}
