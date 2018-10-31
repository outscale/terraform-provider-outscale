package outscale

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-outscale/osc/oapi"
)

func dataSourceOutscaleOAPIVpcs() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOutscaleOAPIVpcsRead,

		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			"net_id": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"net": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
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
						"tags": tagsOAPISchemaComputed(),
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

func dataSourceOutscaleOAPIVpcsRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OAPI

	req := oapi.ReadNetsRequest{}

	filters, filtersOk := d.GetOk("filter")
	netIds, netIdsOk := d.GetOk("net_id")

	if filtersOk == false && netIdsOk == false {
		return fmt.Errorf("filters or net_id(s) must be provided")
	}

	if filtersOk {
		req.Filters = buildOutscaleOAPIDataSourceNetFilters(filters.(*schema.Set))
	}

	if netIdsOk {
		ids := make([]string, len(netIds.([]interface{})))

		for k, v := range netIds.([]interface{}) {
			ids[k] = v.(string)
		}

		req.Filters.NetIds = ids
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
		return fmt.Errorf("no matching VPC found")
	}

	d.SetId(resource.UniqueId())

	nets := make([]map[string]interface{}, len(resp.OK.Nets))

	for i, v := range resp.OK.Nets {
		net := make(map[string]interface{})

		net["net_id"] = v.NetId
		net["ip_range"] = v.IpRange
		net["dhcp_options_set_id"] = v.DhcpOptionsSetId
		net["tenancy"] = v.Tenancy
		net["state"] = v.State
		if v.Tags != nil {
			net["tags"] = tagsOAPIToMapString(v.Tags)
		}

		nets[i] = net
	}

	d.Set("net", nets)
	d.Set("request_id", resp.OK.ResponseContext.RequestId)

	return nil
}
