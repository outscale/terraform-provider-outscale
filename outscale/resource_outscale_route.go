package outscale

import (
	"context"
	"errors"
	"fmt"
	"log"
	"reflect"
	"strings"
	"time"

	"github.com/openlyinc/pointy"
	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/terraform-providers/terraform-provider-outscale/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var errOAPIRoute = errors.New("Error: more than 1 target specified. Only 1 of gateway_id, " +
	"nat_service_id, vm_id, nic_id or net_peering_id is allowed.")

var allowedTargets = []string{
	"gateway_id",
	"nat_service_id",
	"vm_id",
	"nic_id",
	"net_peering_id",
}

func resourceOutscaleOAPIRoute() *schema.Resource {
	return &schema.Resource{
		Create: resourceOutscaleOAPIRouteCreate,
		Read:   resourceOutscaleOAPIRouteRead,
		Update: resourceOutscaleOAPIRouteUpdate,
		Delete: resourceOutscaleOAPIRouteDelete,
		Exists: resourceOutscaleOAPIRouteExists,
		Importer: &schema.ResourceImporter{
			State: resourceOutscaleOAPIRouteImportState,
		},
		Schema: map[string]*schema.Schema{
			"creation_method": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"destination_ip_range": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"destination_service_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"gateway_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: allowedTargets,
			},
			"nat_service_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: allowedTargets,
			},
			"nat_access_point": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"net_peering_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: allowedTargets,
			},
			"nic_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true, // Computed because if vm_id is set, and the nic is attached to a VM, it will be set
				ExactlyOneOf: allowedTargets,
			},
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"await_active_state": {
				Type:     schema.TypeBool,
				Default:  true,
				Optional: true,
			},
			"vm_account_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vm_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true, // Computed because if nic_id is set, and the nic is attached to a VM, it will be set
				ExactlyOneOf: allowedTargets,
			},
			"route_table_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"request_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceOutscaleOAPIRouteCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI
	numTargets, target := getTarget(d)
	awaitActiveState := d.Get("await_active_state").(bool)

	if numTargets > 1 {
		return errOAPIRoute
	}

	createOpts := oscgo.CreateRouteRequest{
		RouteTableId:       d.Get("route_table_id").(string),
		DestinationIpRange: d.Get("destination_ip_range").(string),
	}
	switch target {
	case "gateway_id":
		createOpts.SetGatewayId(d.Get("gateway_id").(string))
	case "nat_service_id":
		createOpts.SetNatServiceId(d.Get("nat_service_id").(string))
	case "vm_id":
		createOpts.SetVmId(d.Get("vm_id").(string))
	case "nic_id":
		createOpts.SetNicId(d.Get("nic_id").(string))
	case "net_peering_id":
		createOpts.SetNetPeeringId(d.Get("net_peering_id").(string))
	default:
		return fmt.Errorf("An invalid target type specified: %s", target)
	}
	log.Printf("[DEBUG] Route create config: %+v", createOpts)

	var err error
	err = resource.Retry(2*time.Minute, func() *resource.RetryError {
		_, httpResp, err := conn.RouteApi.CreateRoute(context.Background()).CreateRouteRequest(createOpts).Execute()
		if err != nil {
			if strings.Contains(fmt.Sprint(err), utils.InvalidState) {
				log.Printf("[OKHT] === ERROR: %v ====\n", err)
				log.Printf("[DEBUG] Trying to create route again: %q", err)
				return resource.RetryableError(err)
			}
			return utils.CheckThrottling(httpResp, err)
		}
		return nil
	})

	var errString string

	if err != nil {
		errString = err.Error()
		return fmt.Errorf("Error creating route: %s", errString)
	}

	var route *oscgo.Route
	var requestID string

	if v, ok := d.GetOk("destination_ip_range"); ok {
		err = resource.Retry(2*time.Minute, func() *resource.RetryError {
			route, requestID, err = findResourceOAPIRoute(conn, d.Get("route_table_id").(string), v.(string))
			if awaitActiveState && err == nil {
				if route.GetState() != "active" {
					return resource.RetryableError(fmt.Errorf("still await route to be active"))
				}
			}
			if err != nil {
				return resource.NonRetryableError(err)
			}
			return nil
		})
		if err != nil {
			return fmt.Errorf("Error finding route after creating it: %s", err)
		}
	}
	d.SetId(d.Get("route_table_id").(string) + "_" + d.Get("destination_ip_range").(string))
	return resourceOutscaleOAPIRouteSetResourceData(d, route, requestID)
}

func resourceOutscaleOAPIRouteRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI
	routeTableID := d.Get("route_table_id").(string)

	destinationIPRange := d.Get("destination_ip_range").(string)
	var requestID string

	route, requestID, err := findResourceOAPIRoute(conn, routeTableID, destinationIPRange)
	if err != nil {
		return err
	}
	if route == nil {
		utils.LogManuallyDeleted("Route", d.Id())
		d.SetId("")
		return nil
	}
	return resourceOutscaleOAPIRouteSetResourceData(d, route, requestID)
}

func resourceOutscaleOAPIRouteSetResourceData(d *schema.ResourceData, route *oscgo.Route, requestID string) error {
	if err := d.Set("destination_service_id", route.GetDestinationServiceId()); err != nil {
		return err
	}
	if err := d.Set("gateway_id", route.GetGatewayId()); err != nil {
		return err
	}
	if err := d.Set("vm_id", route.GetVmId()); err != nil {
		return err
	}
	if err := d.Set("nat_access_point", route.GetNetAccessPointId()); err != nil {
		return err
	}
	if err := d.Set("nat_service_id", route.GetNatServiceId()); err != nil {
		return err
	}
	if err := d.Set("nic_id", route.GetNicId()); err != nil {
		return err
	}
	if err := d.Set("net_peering_id", route.GetNetPeeringId()); err != nil {
		return err
	}
	if err := d.Set("vm_account_id", route.GetVmAccountId()); err != nil {
		return err
	}
	if err := d.Set("creation_method", route.GetCreationMethod()); err != nil {
		return err
	}
	if err := d.Set("state", route.GetState()); err != nil {
		return err
	}
	if err := d.Set("request_id", requestID); err != nil {
		return err
	}
	return nil
}

func getTarget(d *schema.ResourceData) (n int, target string) {
	for _, allowedTarget := range allowedTargets {
		if allowed := d.Get(allowedTarget); allowed != nil {
			if len(allowed.(string)) > 0 {
				n++
				target = allowedTarget
			}
		}
	}
	return
}

func resourceOutscaleOAPIRouteUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI
	nothingToDo := true
	o, n := d.GetChange("")
	os := o.(map[string]interface{})
	ns := n.(map[string]interface{})

	for k := range os {
		if d.HasChange(k) && k != "await_active_state" {
			nothingToDo = false
		}
	}

	for k := range ns {
		if d.HasChange(k) && k != "await_active_state" {
			nothingToDo = false
		}
	}

	if nothingToDo == true {
		return nil
	}

	// Check for the new target
	// With ExacltyOneOf, we know that it will only be one target new o none for nic_id

	var target string
	for _, allowedTarget := range allowedTargets {
		old_value := os[allowedTarget]
		new_value := ns[allowedTarget]
		if new_value != "" && old_value != new_value {
			target = allowedTarget
			log.Printf("Possible new target is %v\n", target)
		}
	}

	if target == "" {
		return errors.New("no target found for the update")
	}

	replaceOpts := oscgo.UpdateRouteRequest{}

	switch target {
	case "gateway_id":
		replaceOpts = oscgo.UpdateRouteRequest{
			RouteTableId:       d.Get("route_table_id").(string),
			DestinationIpRange: d.Get("destination_ip_range").(string),
			GatewayId:          pointy.String(d.Get("gateway_id").(string)),
		}
	case "nat_service_id":
		replaceOpts = oscgo.UpdateRouteRequest{
			RouteTableId:       d.Get("route_table_id").(string),
			DestinationIpRange: d.Get("destination_ip_range").(string),
			NatServiceId:       pointy.String(d.Get("nat_service_id").(string)),
		}
	case "vm_id":
		replaceOpts = oscgo.UpdateRouteRequest{
			RouteTableId:       d.Get("route_table_id").(string),
			DestinationIpRange: d.Get("destination_ip_range").(string),
			VmId:               pointy.String(d.Get("vm_id").(string)),
		}
	case "nic_id":
		replaceOpts = oscgo.UpdateRouteRequest{
			RouteTableId:       d.Get("route_table_id").(string),
			DestinationIpRange: d.Get("destination_ip_range").(string),
			NicId:              pointy.String(d.Get("nic_id").(string)),
		}
	case "net_peering_id":
		replaceOpts = oscgo.UpdateRouteRequest{
			RouteTableId:       d.Get("route_table_id").(string),
			DestinationIpRange: d.Get("destination_ip_range").(string),
			NetPeeringId:       pointy.String(d.Get("net_peering_id").(string)),
		}
	default:
		return fmt.Errorf("An invalid target type specified: %s", target)
	}
	log.Printf("[DEBUG] Route replace config: %+v", replaceOpts)

	var err error
	err = resource.Retry(2*time.Minute, func() *resource.RetryError {
		_, httpResp, err := conn.RouteApi.UpdateRoute(context.Background()).UpdateRouteRequest(replaceOpts).Execute()
		if err != nil {
			if strings.Contains(fmt.Sprint(err), utils.InvalidState) {
				log.Printf("[DEBUG] Trying to create route again: %q", err)
				return resource.RetryableError(err)
			}
			return utils.CheckThrottling(httpResp, err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("error updating route: %s", utils.GetErrorResponse(err))
	}

	return resourceOutscaleOAPIRouteRead(d, meta)
}

func resourceOutscaleOAPIRouteDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	deleteOpts := oscgo.DeleteRouteRequest{
		RouteTableId: d.Get("route_table_id").(string),
	}
	if v, ok := d.GetOk("destination_ip_range"); ok {
		deleteOpts.SetDestinationIpRange(v.(string))
	}

	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		log.Printf("[DEBUG] Trying to delete route with opts %+v", deleteOpts)
		resp, httpResp, err := conn.RouteApi.DeleteRoute(context.Background()).DeleteRouteRequest(deleteOpts).Execute()
		log.Printf("[DEBUG] Route delete result: %+v", resp)
		if err != nil {
			if strings.Contains(fmt.Sprint(err), utils.InvalidState) {
				log.Printf("[DEBUG] Trying to delete route again: %q", fmt.Sprint(err))
				return resource.RetryableError(err)
			}
			return utils.CheckThrottling(httpResp, err)
		}
		return nil
	})

	if err != nil {
		return err
	}

	d.SetId("")
	return nil
}

func resourceOutscaleOAPIRouteExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	conn := meta.(*OutscaleClient).OSCAPI
	routeTableID := d.Get("route_table_id").(string)

	findOpts := oscgo.ReadRouteTablesRequest{
		Filters: &oscgo.FiltersRouteTable{RouteTableIds: &[]string{routeTableID}},
	}

	var resp oscgo.ReadRouteTablesResponse
	var err error
	err = resource.Retry(2*time.Minute, func() *resource.RetryError {
		rp, httpResp, err := conn.RouteTableApi.ReadRouteTables(context.Background()).ReadRouteTablesRequest(findOpts).Execute()
		if err != nil {
			if strings.Contains(fmt.Sprint(err), utils.InvalidState) {
				log.Printf("[DEBUG] Trying to create route again: %q", err)
				return resource.RetryableError(err)
			}
			return utils.CheckThrottling(httpResp, err)
		}
		resp = rp
		return nil
	})

	var errString string

	if err != nil {
		if strings.Contains(fmt.Sprint(err), "InvalidRouteTableID.NotFound") {
			log.Printf("[WARN] Route Table %q could not be found.", routeTableID)
			return false, nil
		}
		errString = err.Error()

		return false, fmt.Errorf("Error creating route: %s", errString)
	}

	if len(resp.GetRouteTables()) < 1 || reflect.DeepEqual(resp.GetRouteTables()[0], oscgo.RouteTable{}) {
		log.Printf("[WARN] Route Table %q is gone, or route does not exist.", routeTableID)
		return false, nil
	}

	if v, ok := d.GetOk("destination_ip_range"); ok {
		for _, route := range resp.GetRouteTables()[0].GetRoutes() {
			if route.GetDestinationIpRange() != "" && route.GetDestinationIpRange() == v.(string) {
				return true, nil
			}
		}
	}

	return false, nil
}

func findResourceOAPIRoute(conn *oscgo.APIClient, rtbid string, cidr string) (*oscgo.Route, string, error) {
	routeTableID := rtbid

	findOpts := oscgo.ReadRouteTablesRequest{}
	findOpts.Filters = &oscgo.FiltersRouteTable{
		RouteTableIds: &[]string{routeTableID},
	}

	var resp oscgo.ReadRouteTablesResponse
	var err error
	err = resource.Retry(2*time.Minute, func() *resource.RetryError {
		rp, httpResp, err := conn.RouteTableApi.ReadRouteTables(context.Background()).ReadRouteTablesRequest(findOpts).Execute()
		if err != nil {
			if strings.Contains(fmt.Sprint(err), utils.InvalidState) {
				log.Printf("[DEBUG] Trying to create route again: %q", err)
				return resource.RetryableError(err)
			}
			return utils.CheckThrottling(httpResp, err)
		}
		resp = rp
		return nil
	})

	var errString string

	if err != nil {
		errString = err.Error()
		return nil, "", fmt.Errorf("Error finding route resource: %s", errString)
	}
	requestID := resp.ResponseContext.GetRequestId()

	if len(resp.GetRouteTables()) < 1 || reflect.DeepEqual(resp.GetRouteTables()[0], oscgo.RouteTable{}) {
		return nil, requestID, nil
	}

	if cidr != "" {
		for _, route := range (resp.GetRouteTables()[0]).GetRoutes() {
			if route.GetDestinationIpRange() != "" && route.GetDestinationIpRange() == cidr {
				return &route, requestID, nil
			}
		}

		return nil, requestID, fmt.Errorf("Unable to find matching route for Route Table (%s) "+
			"and destination CIDR block (%s).", rtbid, cidr)
	}

	return nil, requestID, fmt.Errorf("When trying to find a matching route for Route Table %q "+
		"you need to specify a CIDR block", rtbid)

}

func resourceOutscaleOAPIRouteImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	conn := meta.(*OutscaleClient).OSCAPI

	parts := strings.SplitN(d.Id(), "_", 2)
	if len(parts) != 2 {
		return nil, errors.New("import format error: to import a Outscale Route, use the format {route_table_id}_{destination_ip_range}")
	}

	routeTableID := parts[0]
	destinationIPRange := parts[1]

	route, _, err := findResourceOAPIRoute(conn, routeTableID, destinationIPRange)
	if err != nil {
		return nil, err
	}
	if route == nil {
		d.SetId("")
		return nil, fmt.Errorf("Route Table %q is gone, or route does not exist", routeTableID)
	}
	if err := d.Set("route_table_id", routeTableID); err != nil {
		return nil, fmt.Errorf("error setting `%s` for Outscale Route(%s): %s", "route_table_id", routeTableID, err)
	}
	if err := d.Set("destination_ip_range", destinationIPRange); err != nil {
		return nil, fmt.Errorf("error setting `%s` for Outscale Route(%s): %s", "destination_ip_range", destinationIPRange, err)
	}

	d.SetId(d.Get("route_table_id").(string) + "_" + d.Get("destination_ip_range").(string))

	return []*schema.ResourceData{d}, nil
}
