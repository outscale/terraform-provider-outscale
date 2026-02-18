package oapi

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/outscale/osc-sdk-go/v3/pkg/osc"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/utils"
)

func DataSourceOutscaleClientGateway() *schema.Resource {
	return &schema.Resource{
		Read: DataSourceOutscaleClientGatewayRead,
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
			"clientection_type": {
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

func DataSourceOutscaleClientGatewayRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*client.OutscaleClient).OSC

	filters, filtersOk := d.GetOk("filter")
	clientGatewayID, clientGatewayOk := d.GetOk("client_gateway_id")

	if !filtersOk && !clientGatewayOk {
		return fmt.Errorf("one of filters, or client_gateway_id must be assigned")
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
			return err
		}
		params.Filters = filterParams
	}

	var resp osc.ReadClientGatewaysResponse
	err := retry.Retry(5*time.Minute, func() *retry.RetryError {
		rp, httpResp, err := client.ClientGatewayApi.ReadClientGateways(ctx).ReadClientGatewaysRequest(params).Execute()
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

	if len(resp.GetClientGateways()) > 1 {
		return ErrMultipleResults
	}

	clientGateway := resp.GetClientGateways()[0]

	if err := d.Set("bgp_asn", clientGateway.GetBgpAsn()); err != nil {
		return err
	}
	if err := d.Set("client_gateway_id", clientGateway.GetClientGatewayId()); err != nil {
		return err
	}
	if err := d.Set("clientection_type", clientGateway.GetclientectionType()); err != nil {
		return err
	}
	if err := d.Set("public_ip", clientGateway.GetPublicIp()); err != nil {
		return err
	}
	if err := d.Set("state", clientGateway.GetState()); err != nil {
		return err
	}
	if err := d.Set("tags", FlattenOAPITagsSDK(clientGateway.Tags)); err != nil {
		return err
	}

	d.SetId(clientGateway.GetClientGatewayId())

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
			filters.SetBgpAsns(utils.StringSliceToInt32Slice(filterValues))
		case "client_gateway_ids":
			filters.SetClientGatewayIds(filterValues)
		case "clientection_types":
			filters.SetclientectionTypes(filterValues)
		case "public_ips":
			filters.SetPublicIps(filterValues)
		case "states":
			filters.SetStates(filterValues)
		case "tag_keys":
			filters.SetTagKeys(filterValues)
		case "tag_values":
			filters.SetTagValues(filterValues)
		case "tags":
			filters.SetTags(filterValues)
		default:
			return nil, utils.UnknownDataSourceFilterError(ctx, name)
		}
	}
	return &filters, nil
}
