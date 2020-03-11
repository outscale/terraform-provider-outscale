package outscale

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/antihax/optional"
	oscgo "github.com/marinsalinas/osc-sdk-go"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceOutscaleClientGateway() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOutscaleClientGatewayRead,
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
			"tags": tagsOAPIListSchemaComputed(),
			"request_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceOutscaleClientGatewayRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	filters, filtersOk := d.GetOk("filter")
	clientGatewayID, clientGatewayOk := d.GetOk("client_gateway_id")

	if !filtersOk && !clientGatewayOk {
		return fmt.Errorf("One of filters, or client_gateway_id must be assigned")
	}

	params := oscgo.ReadClientGatewaysRequest{}

	if clientGatewayOk {
		params.Filters = &oscgo.FiltersClientGateway{
			ClientGatewayIds: &[]string{clientGatewayID.(string)},
		}
	}

	if filtersOk {
		params.Filters = buildOutscaleDataSourceClientGatewayFilters(filters.(*schema.Set))
	}

	var resp oscgo.ReadClientGatewaysResponse
	var err error
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, _, err = conn.ClientGatewayApi.ReadClientGateways(context.Background(), &oscgo.ReadClientGatewaysOpts{
			ReadClientGatewaysRequest: optional.NewInterface(params),
		})
		if err != nil {
			if strings.Contains(fmt.Sprint(err), "RequestLimitExceeded:") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return err
	}

	if len(resp.GetClientGateways()) == 0 {
		return fmt.Errorf("Unable to find Client Gateway")
	}

	if len(resp.GetClientGateways()) > 1 {
		return fmt.Errorf("multiple results returned, please use a more specific criteria in your query")
	}

	clientGateway := resp.GetClientGateways()[0]

	if err := d.Set("bgp_asn", clientGateway.GetBgpAsn()); err != nil {
		return err
	}
	if err := d.Set("client_gateway_id", clientGateway.GetClientGatewayId()); err != nil {
		return err
	}
	if err := d.Set("connection_type", clientGateway.GetConnectionType()); err != nil {
		return err
	}
	if err := d.Set("public_ip", clientGateway.GetPublicIp()); err != nil {
		return err
	}
	if err := d.Set("state", clientGateway.GetState()); err != nil {
		return err
	}
	if err := d.Set("tags", tagsOSCAPIToMap(clientGateway.GetTags())); err != nil {
		return err
	}
	if err := d.Set("request_id", resp.ResponseContext.GetRequestId()); err != nil {
		return err
	}

	d.SetId(clientGateway.GetClientGatewayId())

	return nil
}

func buildOutscaleDataSourceClientGatewayFilters(set *schema.Set) *oscgo.FiltersClientGateway {
	var filters oscgo.FiltersClientGateway
	for _, v := range set.List() {
		m := v.(map[string]interface{})
		var filterValues []string
		for _, e := range m["values"].([]interface{}) {
			filterValues = append(filterValues, e.(string))
		}

		switch name := m["name"].(string); name {
		case "client_gateway_ids":
			filters.SetClientGatewayIds(filterValues)
		case "connection_types":
			filters.SetConnectionTypes(filterValues)
		case "public_ips":
			filters.SetPublicIps(filterValues)
		case "states":
			filters.SetStates(filterValues)
		case "tag_keys":
			filters.SetTagKeys(filterValues)
		case "tag_values":
			filters.SetTagKeys(filterValues)
		case "tags":
			filters.SetTags(filterValues)
		default:
			log.Printf("[Debug] Unknown Filter Name: %s.", name)
		}
	}
	return &filters
}
