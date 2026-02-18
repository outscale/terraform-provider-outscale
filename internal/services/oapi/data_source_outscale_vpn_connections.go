package oapi

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/outscale/osc-sdk-go/v3/pkg/osc"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/utils"
)

func DataSourceOutscaleVPNConnections() *schema.Resource {
	return &schema.Resource{
		Read: DataSourceOutscaleVPNclientectionsRead,

		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			"vpn_clientection_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"vpn_clientections": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"vpn_clientection_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"client_gateway_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"virtual_gateway_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"clientection_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"static_routes_only": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"client_gateway_configuration": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"state": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"routes": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"destination_ip_range": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"route_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"state": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"tags": TagsSchemaComputedSDK(),
						"vgw_telemetries": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"accepted_route_count": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"last_state_change_date": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"outside_ip_address": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"state": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"state_description": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
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

func DataSourceOutscaleVPNclientectionsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*client.OutscaleClient).OSC

	filters, filtersOk := d.GetOk("filter")
	vpnclientectionIDs, vpnclientectionOk := d.GetOk("vpn_clientection_ids")

	if !filtersOk && !vpnclientectionOk {
		return fmt.Errorf("one of filters, or vpn_clientection_ids must be assigned")
	}

	log.Printf("vpnclientectionIDs: %#+v\n", vpnclientectionIDs)
	params := osc.ReadVpnclientectionsRequest{}

	if vpnclientectionOk {
		params.Filters = &osc.FiltersVpnclientection{
			VpnclientectionIds: utils.InterfaceSliceToStringSlicePtr(vpnclientectionIDs.([]interface{})),
		}
	}

	var err error
	if filtersOk {
		params.Filters, err = buildOutscaleDataSourceVPNclientectionFilters(filters.(*schema.Set))
		if err != nil {
			return err
		}
	}

	var resp osc.ReadVpnclientectionsResponse
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		rp, httpResp, err := client.VpnclientectionApi.ReadVpnclientections(ctx).ReadVpnclientectionsRequest(params).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		resp = rp
		return nil
	})
	if err != nil {
		return err
	}

	if len(resp.GetVpnclientections()) == 0 {
		return ErrNoResults
	}
	if err := d.Set("vpn_clientections", flattenVPNclientections(resp.GetVpnclientections())); err != nil {
		return err
	}

	d.SetId(id.UniqueId())
	return nil
}

func flattenVPNclientections(vpnclientections []osc.Vpnclientection) []map[string]interface{} {
	vpnclientectionsMap := make([]map[string]interface{}, len(vpnclientections))

	for i, vpnclientection := range vpnclientections {
		vpnclientectionsMap[i] = map[string]interface{}{
			"vpn_clientection_id":          vpnclientection.GetVpnclientectionId(),
			"client_gateway_id":            vpnclientection.GetClientGatewayId(),
			"virtual_gateway_id":           vpnclientection.GetVirtualGatewayId(),
			"clientection_type":            vpnclientection.GetclientectionType(),
			"static_routes_only":           vpnclientection.GetStaticRoutesOnly(),
			"client_gateway_configuration": vpnclientection.GetClientGatewayConfiguration(),
			"state":                        vpnclientection.GetState(),
			"routes":                       flattenVPNclientection(vpnclientection.GetRoutes()),
			"tags":                         FlattenOAPITagsSDK(vpnclientection.Tags),
			"vgw_telemetries":              flattenVgwTelemetries(vpnclientection.GetVgwTelemetries()),
		}
	}
	return vpnclientectionsMap
}
