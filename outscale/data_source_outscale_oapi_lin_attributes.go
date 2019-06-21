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

func dataSourceOutscaleOAPIVpcAttr() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOutscaleOAPIVpcAttrRead,

		Schema: map[string]*schema.Schema{
			//"filter": dataSourceFiltersSchema(),
			"dhcp_options_set_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"net_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"request_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceOutscaleOAPIVpcAttrRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OAPI

	filters := oapi.FiltersNet{
		NetIds: []string{d.Get("net_id").(string)},
	}

	req := oapi.ReadNetsRequest{
		Filters: filters,
	}

	var rs *oapi.POST_ReadNetsResponses
	var resp *oapi.ReadNetsResponse
	var err error
	err = resource.Retry(120*time.Second, func() *resource.RetryError {
		rs, err = conn.POST_ReadNets(req)

		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded:") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return resource.RetryableError(err)
	})
	if err != nil {
		log.Printf("[DEBUG] Error reading lin (%s)", err)
	}

	resp = rs.OK

	if resp == nil || len(resp.Nets) == 0 {
		d.SetId("")
		return fmt.Errorf("oAPI Net not found")
	}

	d.SetId(resp.Nets[0].NetId)

	d.Set("net_id", resp.Nets[0].NetId)
	d.Set("dhcp_options_set_id", resp.Nets[0].DhcpOptionsSetId)
	d.Set("request_id", resp.ResponseContext.RequestId)

	return nil
}
