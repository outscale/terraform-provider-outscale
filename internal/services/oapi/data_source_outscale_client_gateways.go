package oapi

import (
	"context"
	"fmt"
	"time"

	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceOutscaleClientGateways() *schema.Resource {
	return &schema.Resource{
		Read: DataSourceOutscaleClientGatewaysRead,
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

func DataSourceOutscaleClientGatewaysRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*client.OutscaleClient).OSCAPI

	filters, filtersOk := d.GetOk("filter")
	clientGatewayIDs, clientGatewayOk := d.GetOk("client_gateway_ids")

	if !filtersOk && !clientGatewayOk {
		return fmt.Errorf("one of filters, or client_gateway_id must be assigned")
	}

	params := oscgo.ReadClientGatewaysRequest{}

	if clientGatewayOk {
		params.Filters = &oscgo.FiltersClientGateway{
			ClientGatewayIds: utils.InterfaceSliceToStringList(clientGatewayIDs.([]interface{})),
		}
	}

	if filtersOk {
		var err error
		params.Filters, err = buildOutscaleDataSourceClientGatewayFilters(filters.(*schema.Set))
		if err != nil {
			return err
		}
	}

	var resp oscgo.ReadClientGatewaysResponse
	err := retry.Retry(5*time.Minute, func() *retry.RetryError {
		rp, httpResp, err := conn.ClientGatewayApi.ReadClientGateways(context.Background()).ReadClientGatewaysRequest(params).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		resp = rp
		return nil
	})
	if err != nil {
		return err
	}

	if len(resp.GetClientGateways()) == 0 {
		return ErrNoResults
	}

	if err := d.Set("client_gateways", flattenClientGateways(resp.GetClientGateways())); err != nil {
		return err
	}

	d.SetId(id.UniqueId())
	return nil
}

func flattenClientGateways(clientGateways []oscgo.ClientGateway) []map[string]interface{} {
	clientGatewaysMap := make([]map[string]interface{}, len(clientGateways))

	for i, clientGateway := range clientGateways {
		clientGatewaysMap[i] = map[string]interface{}{
			"bgp_asn":           clientGateway.GetBgpAsn(),
			"client_gateway_id": clientGateway.GetClientGatewayId(),
			"connection_type":   clientGateway.GetConnectionType(),
			"public_ip":         clientGateway.GetPublicIp(),
			"state":             clientGateway.GetState(),
			"tags":              FlattenOAPITagsSDK(clientGateway.GetTags()),
		}
	}
	return clientGatewaysMap
}
