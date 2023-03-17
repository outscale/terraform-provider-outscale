package outscale

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/terraform-providers/terraform-provider-outscale/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceOutscaleOAPIRouteTable() *schema.Resource {
	return &schema.Resource{
		Create: resourceOutscaleOAPIRouteTableCreate,
		Read:   resourceOutscaleOAPIRouteTableRead,
		Update: resourceOutscaleOAPIRouteTableUpdate,
		Delete: resourceOutscaleOAPIRouteTableDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"net_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"route_table_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"request_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"tags": tagsListOAPISchema(),

			"route_propagating_virtual_gateways": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"virtual_gateway_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
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
						"destination_service_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"gateway_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vm_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vm_account_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"nic_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"creation_method": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"state": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"net_access_point_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"net_peering_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"nat_service_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"link_route_tables": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"main": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"route_table_to_subnet_link_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"link_route_table_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"route_table_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"subnet_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func resourceOutscaleOAPIRouteTableCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	createOpts := oscgo.CreateRouteTableRequest{
		NetId: d.Get("net_id").(string),
	}
	log.Printf("[DEBUG] RouteTable create config: %#v", createOpts)

	var resp oscgo.CreateRouteTableResponse
	var err error
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		rp, httpResp, err := conn.RouteTableApi.CreateRouteTable(context.Background()).CreateRouteTableRequest(createOpts).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		resp = rp
		return nil
	})
	var errString string
	if err != nil {
		errString = err.Error()

		return fmt.Errorf("Error creating route table: %s", errString)
	}

	d.SetId(resp.RouteTable.GetRouteTableId())
	log.Printf("[INFO] Route Table ID: %s", d.Id())

	log.Printf("[DEBUG] Waiting for route table (%s) to become available", d.Id())
	stateConf := &resource.StateChangeConf{
		Pending: []string{"pending"},
		Target:  []string{"ready"},
		Refresh: resourceOutscaleOAPIRouteTableStateRefreshFunc(conn, d.Id()),
		Timeout: 5 * time.Minute,
	}
	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf(
			"Error waiting for route table (%s) to become available: %s",
			d.Id(), err)
	}

	if d.IsNewResource() {
		if err := setOSCAPITags(conn, d); err != nil {
			return err
		}
	}

	a := make([]interface{}, 0)

	//d.Set("tags", a)
	if err := d.Set("routes", a); err != nil {
		return err
	}
	if err := d.Set("link_route_tables", a); err != nil {
		return err
	}

	return resourceOutscaleOAPIRouteTableRead(d, meta)
}

func resourceOutscaleOAPIRouteTableRead(d *schema.ResourceData, meta interface{}) error {
	rtRaw, _, err := readOAPIRouteTable(meta.(*OutscaleClient).OSCAPI, d.Id())
	if err != nil {
		return err
	}
	if rtRaw == nil {
		utils.LogManuallyDeleted("RouteTable", d.Id())
		d.SetId("")
		return nil
	}

	rt := rtRaw.(oscgo.RouteTable)
	if err := d.Set("route_table_id", rt.GetRouteTableId()); err != nil {
		return err
	}
	if err := d.Set("net_id", rt.GetNetId()); err != nil {
		return err
	}
	if err := d.Set("route_propagating_virtual_gateways", setOSCAPIPropagatingVirtualGateways(rt.GetRoutePropagatingVirtualGateways())); err != nil {
		return err
	}
	if err := d.Set("routes", setOSCAPIRoutes(rt.GetRoutes())); err != nil {
		return err
	}
	if err := d.Set("link_route_tables", setOSCAPILinkRouteTables(rt.GetLinkRouteTables())); err != nil {
		return err
	}
	if err := d.Set("tags", tagsOSCAPIToMap(rt.GetTags())); err != nil {
		return err
	}

	return nil
}

func resourceOutscaleOAPIRouteTableUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	if err := setOSCAPITags(conn, d); err != nil {
		return err
	}
	return resourceOutscaleOAPIRouteTableRead(d, meta)
}

func resourceOutscaleOAPIRouteTableDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	rtRaw, _, err := readOAPIRouteTable(meta.(*OutscaleClient).OSCAPI, d.Id())
	if err != nil {
		return err
	}
	if rtRaw == nil {
		return nil
	}
	rt := rtRaw.(oscgo.RouteTable)

	for _, a := range rt.GetLinkRouteTables() {
		if !a.GetMain() {
			log.Printf("[INFO] Unlinking LinkRouteTable: %s", a.GetLinkRouteTableId())

			var err error
			var statusCode int
			err = resource.Retry(5*time.Minute, func() *resource.RetryError {
				_, httpResp, err := conn.RouteTableApi.UnlinkRouteTable(context.Background()).UnlinkRouteTableRequest(oscgo.UnlinkRouteTableRequest{
					LinkRouteTableId: a.GetLinkRouteTableId(),
				}).Execute()
				if err != nil {
					return utils.CheckThrottling(httpResp, err)
				}
				statusCode = httpResp.StatusCode
				return nil
			})

			if err != nil {
				if statusCode == http.StatusNotFound {
					return nil
				}
				return err
			}
		}
	}

	log.Printf("[INFO] Deleting Route Table: %s", d.Id())
	var statusCode int
	err = resource.Retry(15*time.Minute, func() *resource.RetryError {
		_, httpResp, err := conn.RouteTableApi.DeleteRouteTable(context.Background()).DeleteRouteTableRequest(oscgo.DeleteRouteTableRequest{
			RouteTableId: d.Id(),
		}).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		statusCode = httpResp.StatusCode
		return nil
	})

	if err != nil {
		if statusCode == http.StatusNotFound {
			return nil
		}

		return fmt.Errorf("Error deleting route table: %s", err)
	}

	log.Printf("[DEBUG] Waiting for route table (%s) to become destroyed", d.Id())

	stateConf := &resource.StateChangeConf{
		Pending: []string{"ready"},
		Target:  []string{},
		Refresh: resourceOutscaleOAPIRouteTableStateRefreshFunc(conn, d.Id()),
		Timeout: 5 * time.Minute,
	}
	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("Error waiting for route table (%s) to become destroyed: %s", d.Id(), err)
	}

	return nil
}

func readOAPIRouteTable(conn *oscgo.APIClient, routeTableID string, linkIds ...string) (interface{}, string, error) {
	log.Printf("[DEBUG] Looking for RouteTable with: id %v and link_ids %v", routeTableID, linkIds)
	var resp oscgo.ReadRouteTablesResponse
	var err error
	routeTableRequest := oscgo.ReadRouteTablesRequest{}
	routeTableRequest.Filters = &oscgo.FiltersRouteTable{RouteTableIds: &[]string{routeTableID}}

	err = resource.Retry(15*time.Minute, func() *resource.RetryError {
		var err error
		rp, httpResp, err := conn.RouteTableApi.ReadRouteTables(context.Background()).ReadRouteTablesRequest(routeTableRequest).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		resp = rp
		return nil
	})

	var errString string
	var requestID string
	if err != nil {
		errString = err.Error()
		return nil, requestID, fmt.Errorf("Error getting route table: %s", errString)
	}
	if len(resp.GetRouteTables()) <= 0 {
		return nil, resp.ResponseContext.GetRequestId(), err
	}

	//Fix for OAPI issue when passing routeTableIds and routeTableLinkIds
	rts := resp.GetRouteTables()[0].GetLinkRouteTables()

	if len(linkIds) > 0 {
		for _, linkID := range linkIds {
			for _, rt := range rts {
				if rt.GetLinkRouteTableId() == linkID {
					return resp.GetRouteTables()[0], resp.ResponseContext.GetRequestId(), err
				}
			}
		}
		return nil, resp.ResponseContext.GetRequestId(), fmt.Errorf("Error getting route table: LinkRouteTables didn't match with provided (%+v)", linkIds)
	}
	return resp.GetRouteTables()[0], resp.ResponseContext.GetRequestId(), err
}

func resourceOutscaleOAPIRouteTableStateRefreshFunc(conn *oscgo.APIClient, routeTableID string, linkIds ...string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		rtRaw, _, err := readOAPIRouteTable(conn, routeTableID, linkIds...)
		if rtRaw == nil {
			return nil, "", err
		}
		return rtRaw.(oscgo.RouteTable), "ready", err
	}
}

func setOSCAPIRoutes(rt []oscgo.Route) []map[string]interface{} {
	route := make([]map[string]interface{}, len(rt))
	if len(rt) > 0 {
		for k, r := range rt {
			m := make(map[string]interface{})

			if r.GetNatServiceId() != "" {
				m["nat_service_id"] = r.GetNatServiceId()
			}

			if r.GetCreationMethod() != "" {
				m["creation_method"] = r.GetCreationMethod()
			}
			if r.GetDestinationIpRange() != "" {
				m["destination_ip_range"] = r.GetDestinationIpRange()
			}
			if r.GetDestinationServiceId() != "" {
				m["destination_service_id"] = r.GetDestinationServiceId()
			}
			if r.GetGatewayId() != "" {
				m["gateway_id"] = r.GetGatewayId()
			}
			if r.GetNetAccessPointId() != "" {
				m["net_access_point_id"] = r.GetNetAccessPointId()
			}
			if r.GetNetPeeringId() != "" {
				m["net_peering_id"] = r.GetNetPeeringId()
			}
			if r.GetVmId() != "" {
				m["vm_id"] = r.GetVmId()
			}
			if r.GetNicId() != "" {
				m["nic_id"] = r.GetNicId()
			}
			if r.GetState() != "" {
				m["state"] = r.GetState()
			}
			if r.GetVmAccountId() != "" {
				m["vm_account_id"] = r.GetVmAccountId()
			}
			route[k] = m
		}
	}

	return route
}

func setOSCAPILinkRouteTables(rt []oscgo.LinkRouteTable) []map[string]interface{} {
	linkRouteTables := make([]map[string]interface{}, len(rt))
	log.Printf("[DEBUG] LinkRouteTable: %#v", rt)
	if len(rt) > 0 {
		for k, r := range rt {
			m := make(map[string]interface{})
			if r.GetMain() {
				m["main"] = r.GetMain()
			}
			if r.GetRouteTableId() != "" {
				m["route_table_id"] = r.GetRouteTableId()
			}
			if r.GetLinkRouteTableId() != "" {
				m["link_route_table_id"] = r.GetLinkRouteTableId()
			}
			if r.GetSubnetId() != "" {
				m["subnet_id"] = r.GetSubnetId()
			}
			linkRouteTables[k] = m
		}
	}

	return linkRouteTables
}

func setOSCAPIPropagatingVirtualGateways(vg []oscgo.RoutePropagatingVirtualGateway) (propagatingVGWs []map[string]interface{}) {
	propagatingVGWs = make([]map[string]interface{}, len(vg))

	if len(vg) > 0 {
		for k, vgw := range vg {
			m := make(map[string]interface{})
			if vgw.GetVirtualGatewayId() != "" {
				m["virtual_gateway_id"] = vgw.GetVirtualGatewayId()
			}
			propagatingVGWs[k] = m
		}
	}
	return propagatingVGWs
}
