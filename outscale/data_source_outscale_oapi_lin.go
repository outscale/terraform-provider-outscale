package outscale

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-outscale/osc/oapi"
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
	conn := meta.(*OutscaleClient).OAPI

	req := oapi.ReadNetsRequest{}

	if v, ok := d.GetOk("filters"); ok {
		req.Filters = buildOutscaleOAPIDataSourceNetFilters(v.(*schema.Set))
	}

	if id := d.Get("net_id"); id != "" {
		req.Filters.NetIds = []string{id.(string)}
	}

	var err error
	var resp *oapi.POST_ReadNetsResponses
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		var err error

		resp, err = conn.POST_ReadNets(req)
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
	if resp == nil || len(resp.OK.Nets) == 0 {
		return fmt.Errorf("No matching Net found")
	}
	if len(resp.OK.Nets) > 1 {
		return fmt.Errorf("Multiple Nets matched; use additional constraints to reduce matches to a single Net")
	}

	net := resp.OK.Nets[0]

	d.SetId(net.NetId)
	d.Set("net_id", net.NetId)
	d.Set("ip_range", net.IpRange)
	d.Set("dhcp_options_set_id", net.DhcpOptionsSetId)
	d.Set("tenancy", net.Tenancy)
	d.Set("state", net.State)
	d.Set("request_id", resp.OK.ResponseContext.RequestId)

	return d.Set("tags", tagsOAPIToMap(net.Tags))
}

func buildOutscaleOAPIDataSourceNetFilters(set *schema.Set) oapi.FiltersNet {
	var filters oapi.FiltersNet
	for _, v := range set.List() {
		m := v.(map[string]interface{})
		var filterValues []string
		for _, e := range m["values"].([]interface{}) {
			filterValues = append(filterValues, e.(string))
		}

		switch name := m["name"].(string); name {
		case "ip-range":
			filters.IpRanges = filterValues
		case "dhcp-options-set-id":
			filters.DhcpOptionsSetIds = filterValues
		case "is-default":
			//bool
			//filters.IsDefault = filterValues
		case "state":
			filters.States = filterValues
		case "tag-key":
			filters.TagKeys = filterValues
		case "tag-value":
			filters.TagValues = filterValues
		default:
			log.Printf("[Debug] Unknown Filter Name: %s.", name)
		}
	}
	return filters
}
