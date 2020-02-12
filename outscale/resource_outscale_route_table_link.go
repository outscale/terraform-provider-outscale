package outscale

import (
	"context"
	"fmt"
	"github.com/antihax/optional"
	oscgo "github.com/marinsalinas/osc-sdk-go"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceOutscaleOAPILinkRouteTable() *schema.Resource {
	return &schema.Resource{
		Create: resourceOutscaleOAPILinkRouteTableCreate,
		Read:   resourceOutscaleOAPILinkRouteTableRead,
		Delete: resourceOutscaleOAPILinkRouteTableDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"subnet_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"route_table_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"link_route_table_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"main": {
				Type:     schema.TypeString,
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
		resp, _, err = conn.RouteTableApi.LinkRouteTable(context.Background(), &oscgo.LinkRouteTableOpts{LinkRouteTableRequest: optional.NewInterface(linkRouteTableOpts)})
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
	d.Set("link_route_table_id", d.Id())
	d.Set("request_id", resp.ResponseContext.GetRequestId())
	log.Printf("[INFO] LinkRouteTable ID: %s", d.Id())

	return nil
}

func resourceOutscaleOAPILinkRouteTableRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	rtRaw, _, err := resourceOutscaleOAPIRouteTableStateRefreshFunc(
		conn, d.Get("route_table_id").(string), d.Get("link_route_table_id").(string))()
	if err != nil {
		return err
	}
	if rtRaw == nil {
		return nil
	}
	rt := rtRaw.(oscgo.RouteTable)
	log.Printf("[DEBUG] LinkRouteTables: %v and %v", rt.LinkRouteTables, d.Get("link_route_table_id"))

	found := false
	for _, a := range rt.GetLinkRouteTables() {
		if a.GetLinkRouteTableId() == d.Id() {
			found = true
			d.Set("subnet_id", a.GetSubnetId())
			d.Set("main", a.GetMain())
			break
		}
	}

	if !found {
		d.SetId("")
	}

	return nil
}

func resourceOutscaleOAPILinkRouteTableDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	log.Printf("[INFO] Deleting link route table: %s", d.Id())

	var err error
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, _, err = conn.RouteTableApi.UnlinkRouteTable(context.Background(), &oscgo.UnlinkRouteTableOpts{UnlinkRouteTableRequest: optional.NewInterface(oscgo.UnlinkRouteTableRequest{
			LinkRouteTableId: d.Id(),
		})})
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
