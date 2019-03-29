package outscale

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-outscale/osc/oapi"
	"github.com/terraform-providers/terraform-provider-outscale/utils"
)

func dataSourceOutscaleOAPISubnet() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOutscaleOAPISubnetRead,

		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			"subregion_name": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"ip_range": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"subnet_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"tag": dataSourceTagsSchema(),

			"net_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"request_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"available_ips_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func dataSourceOutscaleOAPISubnetRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OAPI

	req := &oapi.ReadSubnetsRequest{}

	if id := d.Get("subnet_id"); id != "" {
		req.Filters = oapi.FiltersSubnet{SubnetIds: []string{id.(string)}}
	}

	filters, filtersOk := d.GetOk("filter")

	if filtersOk {
		req.Filters = buildOutscaleOAPISubnetDataSourceFilters(filters.(*schema.Set))
	}

	var resp *oapi.POST_ReadSubnetsResponses
	var err error
	err = resource.Retry(120*time.Second, func() *resource.RetryError {
		resp, err = conn.POST_ReadSubnets(*req)

		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded:") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	var errString string

	if err != nil || resp.OK == nil {
		if err != nil {
			errString = err.Error()
		} else if resp.Code401 != nil {
			errString = fmt.Sprintf("ErrorCode: 401, %s", utils.ToJSONString(resp.Code401))
		} else if resp.Code400 != nil {
			errString = fmt.Sprintf("ErrorCode: 400, %s", utils.ToJSONString(resp.Code400))
		} else if resp.Code500 != nil {
			errString = fmt.Sprintf("ErrorCode: 500, %s", utils.ToJSONString(resp.Code500))
		}

		return fmt.Errorf("[DEBUG] Error reading Subnet (%s)", errString)
	}

	response := resp.OK

	if resp == nil || len(response.Subnets) == 0 {
		return fmt.Errorf("no matching subnet found")
	}

	if len(response.Subnets) > 1 {
		return fmt.Errorf("multiple subnets matched; use additional constraints to reduce matches to a single subnet")
	}

	subnet := response.Subnets[0]

	d.SetId(subnet.SubnetId)
	d.Set("subnet_id", subnet.SubnetId)
	d.Set("net_id", subnet.NetId)
	d.Set("subregion_name", subnet.SubregionName)
	d.Set("ip_range", subnet.IpRange)
	d.Set("state", subnet.State)
	d.Set("tag", tagsOAPIToMap(subnet.Tags))
	d.Set("available_ips_count", subnet.AvailableIpsCount)
	d.Set("request_id", response.ResponseContext.RequestId)

	return nil
}

func buildOutscaleOAPISubnetDataSourceFilters(set *schema.Set) oapi.FiltersSubnet {
	var filters oapi.FiltersSubnet
	for _, v := range set.List() {
		m := v.(map[string]interface{})
		var filterValues []string
		for _, e := range m["values"].([]interface{}) {
			filterValues = append(filterValues, e.(string))
		}

		switch name := m["name"].(string); name {
		// case "available_ips_counts":
		// 	filters.AvailableIpsCounts = filterValues
		case "ip_ranges":
			filters.IpRanges = filterValues
		case "net_ids":
			filters.NetIds = filterValues
		case "states":
			filters.States = filterValues
		case "subnet_ids":
			filters.SubnetIds = filterValues
		case "subregion_names":
			filters.SubregionNames = filterValues
		default:
			log.Printf("[Debug] Unknown Filter Name: %s.", name)
		}
	}
	return filters
}
