package oapi

import (
	"context"
	"log"
	"time"

	"github.com/outscale/osc-sdk-go/v3/pkg/options"
	"github.com/outscale/osc-sdk-go/v3/pkg/osc"
	"github.com/outscale/terraform-provider-outscale/internal/client"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceOutscaleVpcAttr() *schema.Resource {
	return &schema.Resource{
		ReadContext: DataSourceOutscaleVpcAttrRead,

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
			"ip_range": {
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
			"tags": TagsSchemaComputedSDK(),
		},
	}
}

func DataSourceOutscaleVpcAttrRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.OutscaleClient).OSC

	filters := osc.FiltersNet{
		NetIds: &[]string{d.Get("net_id").(string)},
	}

	req := osc.ReadNetsRequest{
		Filters: &filters,
	}

	resp, err := client.ReadNets(ctx, req, options.WithRetryTimeout(120*time.Second))
	if err != nil {
		log.Printf("[DEBUG] Error reading lin (%s)", err)
	}

	if resp.Nets == nil || len(*resp.Nets) == 0 {
		d.SetId("")
		return diag.FromErr(ErrNoResults)
	}
	net := (*resp.Nets)[0]

	d.SetId(net.NetId)

	d.Set("ip_range", net.IpRange)
	d.Set("tenancy", net.Tenancy)
	d.Set("dhcp_options_set_id", net.DhcpOptionsSetId)
	d.Set("net_id", net.NetId)
	d.Set("state", net.State)

	return diag.FromErr(d.Set("tags", FlattenOAPITagsSDK(net.Tags)))
}
