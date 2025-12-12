package outscale

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/outscale/terraform-provider-outscale/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	errorLinkRouteTableSetting = "error setting `%s` for Link Route Table (%s): %s"
)

func ResourceOutscaleLinkRouteTable() *schema.Resource {
	return &schema.Resource{
		Create: ResourceOutscaleLinkRouteTableCreate,
		Read:   ResourceOutscaleLinkRouteTableRead,
		Delete: ResourceOutscaleLinkRouteTableDelete,
		Importer: &schema.ResourceImporter{
			State: ResourceOutscaleLinkRouteTableImportState,
		},
		Schema: map[string]*schema.Schema{
			"subnet_id": {
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

func ResourceOutscaleLinkRouteTableCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI
	var resp oscgo.LinkRouteTableResponse
	var err error

	linkRouteTableOpts := oscgo.LinkRouteTableRequest{
		RouteTableId: d.Get("route_table_id").(string),
		SubnetId:     d.Get("subnet_id").(string),
	}
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		var err error
		rp, httpResp, err := conn.RouteTableApi.LinkRouteTable(context.Background()).LinkRouteTableRequest(linkRouteTableOpts).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		resp = rp
		return nil
	})
	// Set the ID and return
	var errString string
	if err != nil {
		errString = err.Error()

		return fmt.Errorf("Error creating route table link: %s", errString)
	}

	d.SetId(resp.GetLinkRouteTableId())

	return ResourceOutscaleLinkRouteTableRead(d, meta)
}

func ResourceOutscaleLinkRouteTableRead(d *schema.ResourceData, meta interface{}) error {
	routeTableID := d.Get("route_table_id").(string)
	linkRTable, err := readOutscaleLinkRouteTable(meta.(*OutscaleClient), routeTableID, d.Id())
	if err != nil {
		return err
	}
	if linkRTable == nil {
		utils.LogManuallyDeleted("RouteTableLink", d.Id())
		d.SetId("")
		return nil
	}
	if err := d.Set("link_route_table_id", linkRTable.GetLinkRouteTableId()); err != nil {
		return fmt.Errorf(errorLinkRouteTableSetting, "link_route_table_id", linkRTable.GetLinkRouteTableId(), err)
	}
	if err := d.Set("main", linkRTable.GetMain()); err != nil {
		return fmt.Errorf(errorLinkRouteTableSetting, "main", linkRTable.GetLinkRouteTableId(), err)
	}
	if err := d.Set("subnet_id", linkRTable.GetSubnetId()); err != nil {
		return fmt.Errorf(errorLinkRouteTableSetting, "subnet_id", linkRTable.GetLinkRouteTableId(), err)
	}

	return nil
}

func ResourceOutscaleLinkRouteTableDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	err := retry.Retry(5*time.Minute, func() *retry.RetryError {
		_, httpResp, err := conn.RouteTableApi.UnlinkRouteTable(context.Background()).UnlinkRouteTableRequest(oscgo.UnlinkRouteTableRequest{
			LinkRouteTableId: d.Id(),
		}).Execute()
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

func ResourceOutscaleLinkRouteTableImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	parts := strings.SplitN(d.Id(), "_", 2)
	if len(parts) != 2 {
		return nil, errors.New("import format error: to import a Link Route Table, use the format {route_table_id}_{link_route_table_id}")
	}

	routeTableID := parts[0]
	linkRouteTableID := parts[1]

	linkRTable, err := readOutscaleLinkRouteTable(meta.(*OutscaleClient), routeTableID, linkRouteTableID)
	if err != nil {
		return nil, fmt.Errorf("couldn't import Link Route Table(%s), error: %s", linkRouteTableID, err)
	}
	if linkRTable == nil {
		return nil, fmt.Errorf("oAPI route tables for get link table not found")
	}
	if err := d.Set("route_table_id", linkRTable.GetRouteTableId()); err != nil {
		return nil, fmt.Errorf(errorLinkRouteTableSetting, "route_table_id", linkRTable.GetLinkRouteTableId(), err)
	}
	if err := d.Set("main", linkRTable.GetMain()); err != nil {
		return nil, fmt.Errorf(errorLinkRouteTableSetting, "main", linkRTable.GetLinkRouteTableId(), err)
	}
	if err := d.Set("subnet_id", linkRTable.GetSubnetId()); err != nil {
		return nil, fmt.Errorf(errorLinkRouteTableSetting, "subnet_id", linkRTable.GetLinkRouteTableId(), err)
	}

	d.SetId(linkRouteTableID)

	return []*schema.ResourceData{d}, nil
}

func readOutscaleLinkRouteTable(meta *OutscaleClient, routeTableID, linkRouteTableID string) (*oscgo.LinkRouteTable, error) {
	conn := meta.OSCAPI

	var resp oscgo.ReadRouteTablesResponse
	var err error
	routeTableRequest := oscgo.ReadRouteTablesRequest{}
	routeTableRequest.Filters = &oscgo.FiltersRouteTable{RouteTableIds: &[]string{routeTableID}}

	err = retry.Retry(15*time.Minute, func() *retry.RetryError {
		var err error
		rp, httpResp, err := conn.RouteTableApi.ReadRouteTables(context.Background()).ReadRouteTablesRequest(routeTableRequest).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		resp = rp
		return nil
	})

	if err != nil {
		return nil, err
	}

	if len(resp.GetRouteTables()) == 0 {
		return nil, nil
	}

	var linkRTable oscgo.LinkRouteTable
	for _, lTable := range resp.GetRouteTables()[0].GetLinkRouteTables() {
		if lTable.GetLinkRouteTableId() == linkRouteTableID {
			linkRTable = lTable
		}
	}
	return &linkRTable, nil
}
