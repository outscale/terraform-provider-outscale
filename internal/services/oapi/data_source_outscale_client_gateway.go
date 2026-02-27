package oapi

import (
	"context"
	"log"
	"time"

	"github.com/outscale/goutils/sdk/ptr"
	"github.com/outscale/osc-sdk-go/v3/pkg/options"
	"github.com/outscale/osc-sdk-go/v3/pkg/osc"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/utils"
)

func DataSourceOutscaleClientGateway() *schema.Resource {
	return &schema.Resource{
		ReadContext: DataSourceOutscaleClientGatewayRead,
		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
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
			"request_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func DataSourceOutscaleClientGatewayRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.OutscaleClient).OSC

	filters, filtersOk := d.GetOk("filter")
	clientGatewayID, clientGatewayOk := d.GetOk("client_gateway_id")

	if !filtersOk && !clientGatewayOk {
		return diag.Errorf("one of filters, or client_gateway_id must be assigned")
	}

	params := osc.ReadClientGatewaysRequest{}

	if clientGatewayOk {
		params.Filters = &osc.FiltersClientGateway{
			ClientGatewayIds: &[]string{clientGatewayID.(string)},
		}
	}

	if filtersOk {
		filterParams, err := buildOutscaleDataSourceClientGatewayFilters(filters.(*schema.Set))
		if err != nil {
			return diag.FromErr(err)
		}
		params.Filters = filterParams
	}

	resp, err := client.ReadClientGateways(ctx, params, options.WithRetryTimeout(5*time.Minute))
	if err != nil {
		return diag.FromErr(err)
	}

	if resp.ClientGateways == nil || len(*resp.ClientGateways) == 0 {
		return diag.FromErr(ErrNoResults)
	}

	if len(*resp.ClientGateways) > 1 {
		return diag.FromErr(ErrMultipleResults)
	}

	clientGateway := (*resp.ClientGateways)[0]

	if err := d.Set("bgp_asn", ptr.From(clientGateway.BgpAsn)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("client_gateway_id", ptr.From(clientGateway.ClientGatewayId)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("connection_type", ptr.From(clientGateway.ConnectionType)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("public_ip", ptr.From(clientGateway.PublicIp)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("state", ptr.From(clientGateway.State)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("tags", FlattenOAPITagsSDK(ptr.From(clientGateway.Tags))); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(ptr.From(clientGateway.ClientGatewayId))

	return nil
}

func buildOutscaleDataSourceClientGatewayFilters(set *schema.Set) (*osc.FiltersClientGateway, error) {
	var filters osc.FiltersClientGateway
	for _, v := range set.List() {
		log.Printf("[DEBUG] gateway filters %+v", v)
		m := v.(map[string]interface{})
		var filterValues []string
		for _, e := range m["values"].([]interface{}) {
			filterValues = append(filterValues, e.(string))
		}

		switch name := m["name"].(string); name {
		case "bgp_asns":
			filters.BgpAsns = new(utils.StringSliceToIntSlice(filterValues))
		case "client_gateway_ids":
			filters.ClientGatewayIds = &filterValues
		case "connection_types":
			filters.ConnectionTypes = &filterValues
		case "public_ips":
			filters.PublicIps = &filterValues
		case "states":
			filters.States = &filterValues
		case "tag_keys":
			filters.TagKeys = &filterValues
		case "tag_values":
			filters.TagValues = &filterValues
		case "tags":
			filters.Tags = &filterValues
		default:
			return nil, utils.UnknownDataSourceFilterError(name)
		}
	}
	return &filters, nil
}
