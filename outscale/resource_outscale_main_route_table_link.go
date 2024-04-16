package outscale

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/outscale/terraform-provider-outscale/utils"
)

func resourceLinkMainRouteTable() *schema.Resource {
	return &schema.Resource{
		Create: resourceLinkMainRouteTableCreate,
		Read:   resourceLinkMainRouteTableRead,
		Delete: resourceLinkMainRouteTableDelete,
		Schema: map[string]*schema.Schema{
			"net_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"route_table_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"link_route_table_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"default_route_table_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"subnet_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"main": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"request_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceLinkMainRouteTableCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI
	netID := d.Get("net_id").(string)

	routeTable, err := readMainLinkRouteTable(meta.(*OutscaleClient), netID)
	if err != nil {
		return err
	}
	linkRouteTable := routeTable.GetLinkRouteTables()
	oldLinkRouteTableId := linkRouteTable[0].GetLinkRouteTableId()
	defaultRouteTableId := linkRouteTable[0].GetRouteTableId()

	updateRequest := oscgo.UpdateRouteTableLinkRequest{
		RouteTableId: d.Get("route_table_id").(string),
	}
	updateRequest.SetLinkRouteTableId(oldLinkRouteTableId)
	resp := oscgo.UpdateRouteTableLinkResponse{}
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		rp, httpResp, err := conn.RouteTableApi.UpdateRouteTableLink(
			context.Background()).UpdateRouteTableLinkRequest(updateRequest).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		resp = rp
		return nil
	})
	if err != nil {
		return err
	}

	if err := d.Set("default_route_table_id", defaultRouteTableId); err != nil {
		return err
	}
	d.SetId(resp.GetLinkRouteTableId())

	return resourceOutscaleOAPILinkRouteTableRead(d, meta)
}

func resourceLinkMainRouteTableRead(d *schema.ResourceData, meta interface{}) error {
	netID := d.Get("net_id").(string)
	routeTable, err := readMainLinkRouteTable(meta.(*OutscaleClient), netID)
	if err != nil {
		return err
	}
	linkRTable := routeTable.GetLinkRouteTables()
	if linkRTable == nil {
		utils.LogManuallyDeleted("RouteTableLink", d.Id())
		d.SetId("")
		return nil
	}

	if err := d.Set("net_id", linkRTable[0].GetNetId()); err != nil {
		return err
	}
	if linkRTable[0].GetSubnetId() != "" {
		if err := d.Set("subnet_id", linkRTable[0].GetSubnetId()); err != nil {
			return err
		}
	}
	if err := d.Set("link_route_table_id", linkRTable[0].GetLinkRouteTableId()); err != nil {
		return err
	}
	if err := d.Set("main", linkRTable[0].GetMain()); err != nil {
		return err
	}
	if err := d.Set("route_table_id", linkRTable[0].GetRouteTableId()); err != nil {
		return err
	}

	return nil
}

func resourceLinkMainRouteTableDelete(d *schema.ResourceData, meta interface{}) error {
	var err error
	conn := meta.(*OutscaleClient).OSCAPI

	updateRequest := oscgo.UpdateRouteTableLinkRequest{
		LinkRouteTableId: d.Get("link_route_table_id").(string),
		RouteTableId:     d.Get("default_route_table_id").(string),
	}

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, httpResp, err := conn.RouteTableApi.UpdateRouteTableLink(
			context.Background()).UpdateRouteTableLinkRequest(updateRequest).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("Error deleting link route table: %s", err)
	}

	return nil
}

func readMainLinkRouteTable(meta *OutscaleClient, netID string) (oscgo.RouteTable, error) {
	conn := meta.OSCAPI

	var resp oscgo.ReadRouteTablesResponse
	var err error
	var routeTable oscgo.RouteTable

	rtbRequest := oscgo.ReadRouteTablesRequest{}
	rtbRequest.Filters = &oscgo.FiltersRouteTable{
		NetIds:             &[]string{netID},
		LinkRouteTableMain: &[]bool{true}[0],
	}
	err = resource.Retry(15*time.Minute, func() *resource.RetryError {
		var err error
		rp, httpResp, err := conn.RouteTableApi.ReadRouteTables(
			context.Background()).ReadRouteTablesRequest(rtbRequest).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		resp = rp
		return nil
	})

	if err != nil {
		return routeTable, err
	}
	if len(resp.GetRouteTables()) == 0 {
		return routeTable, nil
	}

	return resp.GetRouteTables()[0], nil
}
