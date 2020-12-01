package outscale

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	oscgo "github.com/outscale/osc-sdk-go/v2"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const (
	errorLinkRouteTableSetting = "error setting `%s` for Link Route Table (%s): %s"
)

func resourceOutscaleOAPILinkRouteTable() *schema.Resource {
	return &schema.Resource{
		Create: resourceOutscaleOAPILinkRouteTableCreate,
		Read:   resourceOutscaleOAPILinkRouteTableRead,
		Delete: resourceOutscaleOAPILinkRouteTableDelete,
		Importer: &schema.ResourceImporter{
			State: resourceOutscaleOAPILinkRouteTableImportState,
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

func resourceOutscaleOAPILinkRouteTableCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI
	subnetID := d.Get("subnet_id").(string)
	routeTableID := d.Get("route_table_id").(string)
	log.Printf("[INFO] Creating route table link: %s => %s", subnetID, routeTableID)
	linkRouteTableOpts := oscgo.LinkRouteTableRequest{
		RouteTableId: routeTableID,
		SubnetId:     subnetID,
	}

	var resp oscgo.LinkRouteTableResponse
	var err error
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, _, err = conn.RouteTableApi.LinkRouteTable(context.Background()).LinkRouteTableRequest(linkRouteTableOpts).Execute()
		if err != nil {
			if strings.Contains(fmt.Sprint(err), "InvalidRouteTableID.NotFound") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return err
	}

	// Set the ID and return
	var errString string
	if err != nil {
		errString = err.Error()

		return fmt.Errorf("Error creating route table link: %s", errString)
	}

	d.SetId(resp.GetLinkRouteTableId())

	return resourceOutscaleOAPILinkRouteTableRead(d, meta)
}

func resourceOutscaleOAPILinkRouteTableRead(d *schema.ResourceData, meta interface{}) error {
	routeTable, requestID, err := readOutscaleLinkRouteTable(meta.(*OutscaleClient), d.Get("route_table_id").(string), d.Id())
	if err != nil {
		return err
	}
	if routeTable == nil {
		d.SetId("")
	}

	if err := d.Set("link_route_table_id", routeTable.GetLinkRouteTableId()); err != nil {
		return fmt.Errorf(errorLinkRouteTableSetting, "link_route_table_id", routeTable.GetLinkRouteTableId(), err)
	}
	if err := d.Set("main", routeTable.GetMain()); err != nil {
		return fmt.Errorf(errorLinkRouteTableSetting, "main", routeTable.GetLinkRouteTableId(), err)
	}
	if err := d.Set("request_id", requestID); err != nil {
		return fmt.Errorf(errorLinkRouteTableSetting, "request_id", routeTable.GetLinkRouteTableId(), err)
	}

	return nil
}

func resourceOutscaleOAPILinkRouteTableDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	log.Printf("[INFO] Deleting link route table: %s", d.Id())

	var err error
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, _, err = conn.RouteTableApi.UnlinkRouteTable(context.Background()).UnlinkRouteTableRequest(oscgo.UnlinkRouteTableRequest{
			LinkRouteTableId: d.Id(),
		}).Execute()
		if err != nil {
			if strings.Contains(fmt.Sprint(err), "RequestLimitExceeded") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		if strings.Contains(fmt.Sprint(err), "InvalidAssociationID.NotFound") {
			return nil
		}
		return fmt.Errorf("Error deleting link route table: %s", err)
	}

	return nil
}

func resourceOutscaleOAPILinkRouteTableImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	parts := strings.SplitN(d.Id(), "_", 2)
	if len(parts) != 2 {
		return nil, errors.New("import format error: to import a Link Route Table, use the format {route_table_id}_{link_route_table_id}")
	}

	routeTableID := parts[0]
	linkRouteTableID := parts[1]

	routeTable, _, err := readOutscaleLinkRouteTable(meta.(*OutscaleClient), routeTableID, linkRouteTableID)
	if err != nil {
		return nil, fmt.Errorf("couldn't import Link Route Table(%s), error: %s", linkRouteTableID, err)
	}

	if err := d.Set("route_table_id", routeTable.GetRouteTableId()); err != nil {
		return nil, fmt.Errorf(errorLinkRouteTableSetting, "route_table_id", routeTable.GetLinkRouteTableId(), err)
	}
	if err := d.Set("subnet_id", routeTable.GetSubnetId()); err != nil {
		return nil, fmt.Errorf(errorLinkRouteTableSetting, "subnet_id", routeTable.GetLinkRouteTableId(), err)
	}

	d.SetId(linkRouteTableID)

	return []*schema.ResourceData{d}, nil
}

func readOutscaleLinkRouteTable(meta *OutscaleClient, routeTableID, linkRouteTableID string) (*oscgo.LinkRouteTable, string, error) {
	conn := meta.OSCAPI

	var rt oscgo.ReadRouteTablesResponse
	var err error

	err = resource.Retry(15*time.Minute, func() *resource.RetryError {
		rt, _, err = conn.RouteTableApi.ReadRouteTables(context.Background()).ReadRouteTablesRequest(oscgo.ReadRouteTablesRequest{
			Filters: &oscgo.FiltersRouteTable{RouteTableIds: &[]string{routeTableID}},
		}).Execute()
		if err != nil {
			if strings.Contains(fmt.Sprint(err), "RequestLimitExceeded") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return nil, rt.ResponseContext.GetRequestId(), err
	}

	return getLinkRouteTable(linkRouteTableID, rt.GetRouteTables()[0].GetLinkRouteTables()), rt.ResponseContext.GetRequestId(), nil
}

func getLinkRouteTable(id string, routeTables []oscgo.LinkRouteTable) (routeTable *oscgo.LinkRouteTable) {
	for _, rt := range routeTables {
		if rt.GetLinkRouteTableId() == id {
			routeTable = &rt
		}
	}
	return
}
