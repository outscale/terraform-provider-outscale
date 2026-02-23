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

func DataSourceOutscaleClientGateways() *schema.Resource {
	return &schema.Resource{
		ReadContext: DataSourceOutscaleClientGatewaysRead,
		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			"client_gateway_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"client_gateways": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"bgp_asn": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"client_gateway_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"connection_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"public_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"state": {
							Type:     schema.TypeString,
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

func DataSourceOutscaleClientGatewaysRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.OutscaleClient).OSC

	filters, filtersOk := d.GetOk("filter")
	clientGatewayIDs, clientGatewayOk := d.GetOk("client_gateway_ids")

	if !filtersOk && !clientGatewayOk {
		return diag.Errorf("one of filters, or client_gateway_id must be assigned")
	}

	params := osc.ReadClientGatewaysRequest{}

	if clientGatewayOk {
		params.Filters = &osc.FiltersClientGateway{
			ClientGatewayIds: utils.InterfaceSliceToStringList(clientGatewayIDs.([]interface{})),
		}
	}

	if filtersOk {
		var err error
		params.Filters, err = buildOutscaleDataSourceClientGatewayFilters(filters.(*schema.Set))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	resp, err := client.ReadClientGateways(ctx, params, options.WithRetryTimeout(5*time.Minute))
	if err != nil {
		return diag.FromErr(err)
	}

	if resp.ClientGateways == nil || len(*resp.ClientGateways) == 0 {
		return diag.FromErr(ErrNoResults)
	}

	if err := d.Set("client_gateways", flattenClientGateways(*resp.ClientGateways)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(id.UniqueId())
	return nil
}

func flattenClientGateways(clientGateways []osc.ClientGateway) []map[string]interface{} {
	clientGatewaysMap := make([]map[string]interface{}, len(clientGateways))

	for i, clientGateway := range clientGateways {
		clientGatewaysMap[i] = map[string]interface{}{
			"bgp_asn":           clientGateway.BgpAsn,
			"client_gateway_id": clientGateway.ClientGatewayId,
			"connection_type":   clientGateway.ConnectionType,
			"public_ip":         clientGateway.PublicIp,
			"state":             clientGateway.State,
			"tags":              FlattenOAPITagsSDK(ptr.From(clientGateway.Tags)),
		}
	}
	return clientGatewaysMap
}
