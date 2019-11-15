package outscale

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/outscale/osc-go/oapi"
	"github.com/terraform-providers/terraform-provider-outscale/utils"
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
	conn := meta.(*OutscaleClient).OAPI
	subnetID := d.Get("subnet_id").(string)
	routeTableID := d.Get("route_table_id").(string)
	log.Printf("[INFO] Creating route table link: %s => %s", subnetID, routeTableID)
	linkRouteTableOpts := oapi.LinkRouteTableRequest{
		RouteTableId: routeTableID,
		SubnetId:     subnetID,
	}

	var resp *oapi.POST_LinkRouteTableResponses
	var err error
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, err = conn.POST_LinkRouteTable(linkRouteTableOpts)
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
	if err != nil || resp.OK == nil {
		if err != nil {
			errString = err.Error()
		} else if resp.Code401 != nil {
			errString = fmt.Sprintf("ErrorCode: 401, %s", utils.ToJSONString(resp.Code401))
		} else if resp.Code400 != nil {
			errString = fmt.Sprintf("ErrorCode: 400, %s", utils.ToJSONString(resp.Code400))
		} else if resp.Code500 != nil {
			errString = fmt.Sprintf("ErrorCode: 500, %s", utils.ToJSONString(resp.Code500))
		}

		return fmt.Errorf("Error creating route table link: %s", errString)
	}

	d.SetId(resp.OK.LinkRouteTableId)
	d.Set("link_route_table_id", d.Id())
	d.Set("request_id", resp.OK.ResponseContext.RequestId)
	log.Printf("[INFO] LinkRouteTable ID: %s", d.Id())

	return nil
}

func resourceOutscaleOAPILinkRouteTableRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OAPI

	rtRaw, _, err := resourceOutscaleOAPIRouteTableStateRefreshFunc(
		conn, d.Get("route_table_id").(string), d.Get("link_route_table_id").(string))()
	if err != nil {
		return err
	}
	if rtRaw == nil {
		return nil
	}
	rt := rtRaw.(oapi.RouteTable)
	log.Printf("[DEBUG] LinkRouteTables: %v and %v", rt.LinkRouteTables, d.Get("link_route_table_id"))

	found := false
	for _, a := range rt.LinkRouteTables {
		if a.LinkRouteTableId == d.Id() {
			found = true
			d.Set("subnet_id", a.SubnetId)
			d.Set("main", a.Main)
			break
		}
	}

	if !found {
		d.SetId("")
	}

	return nil
}

func resourceOutscaleOAPILinkRouteTableDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OAPI

	log.Printf("[INFO] Deleting link route table: %s", d.Id())

	var err error
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, err = conn.POST_UnlinkRouteTable(oapi.UnlinkRouteTableRequest{
			LinkRouteTableId: d.Id(),
		})
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
