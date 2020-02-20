package outscale

import (
	"context"
	"errors"
	"fmt"
	"log"
	"reflect"
	"strings"
	"time"

	"github.com/antihax/optional"
	oscgo "github.com/marinsalinas/osc-sdk-go"

	"github.com/hashicorp/terraform-plugin-sdk/helper/hashcode"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"creation_method": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"destination_ip_range": {
				Type:     schema.TypeString,
				Required: true,
			},
			"destination_service_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"gateway_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"nat_service_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"nat_access_point": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"net_peering_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"nic_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vm_account_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vm_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
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

	if numTargets > 1 {
		return errOAPIRoute
	}

	createOpts := oscgo.CreateRouteRequest{}
	switch target {
	case "gateway_id":
		createOpts = oscgo.CreateRouteRequest{
			RouteTableId:       d.Get("route_table_id").(string),
			DestinationIpRange: d.Get("destination_ip_range").(string),
		}
		createOpts.SetGatewayId(d.Get("gateway_id").(string))
	case "nat_service_id":
		createOpts = oscgo.CreateRouteRequest{
			RouteTableId:       d.Get("route_table_id").(string),
			DestinationIpRange: d.Get("destination_ip_range").(string),
		}
		createOpts.SetNatServiceId(d.Get("nat_service_id").(string))
	case "vm_id":
		createOpts = oscgo.CreateRouteRequest{
			RouteTableId:       d.Get("route_table_id").(string),
			DestinationIpRange: d.Get("destination_ip_range").(string),
		}
		createOpts.SetVmId(d.Get("vm_id").(string))
	case "nic_id":
		createOpts = oscgo.CreateRouteRequest{
			RouteTableId:       d.Get("route_table_id").(string),
			DestinationIpRange: d.Get("destination_ip_range").(string),
		}
		createOpts.SetNicId(d.Get("nic_id").(string))
	case "net_peering_id":
		createOpts = oscgo.CreateRouteRequest{
			RouteTableId:       d.Get("route_table_id").(string),
			DestinationIpRange: d.Get("destination_ip_range").(string),
		}
		createOpts.SetNetPeeringId(d.Get("net_peering_id").(string))
	default:
		return fmt.Errorf("An invalid target type specified: %s", target)
	}
	log.Printf("[DEBUG] Route create config: %+v", createOpts)

	var err error
	err = resource.Retry(2*time.Minute, func() *resource.RetryError {
		_, _, err = conn.RouteApi.CreateRoute(context.Background(), &oscgo.CreateRouteOpts{CreateRouteRequest: optional.NewInterface(createOpts)})

		if err != nil {
			if strings.Contains(fmt.Sprint(err), "InvalidParameterException") {
				log.Printf("[DEBUG] Trying to create route again: %q", err)
				return resource.RetryableError(err)
			}

			return resource.NonRetryableError(err)
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
			return resource.RetryableError(err)
		})
		if err != nil {
			return fmt.Errorf("Error finding route after creating it: %s", err)
		}
	}

	d.SetId(routeOAPIIDHash(d, route))
	resourceOutscaleOAPIRouteSetResourceData(d, route, requestID)
	return nil
}

func resourceOutscaleOAPIRouteRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI
	routeTableID := d.Get("route_table_id").(string)

	destinationIPRange := d.Get("destination_ip_range").(string)
	var requestID string

	route, requestID, err := findResourceOAPIRoute(conn, routeTableID, destinationIPRange)
	if err != nil {
		if strings.Contains(fmt.Sprint(err), "InvalidRouteTableID.NotFound") {
			log.Printf("[WARN] Route Table %q could not be found. Removing Route from state.", routeTableID)
			d.SetId("")
			return nil
		}
		return err
	}
	resourceOutscaleOAPIRouteSetResourceData(d, route, requestID)
	return nil
}

func resourceOutscaleOAPIRouteSetResourceData(d *schema.ResourceData, route *oscgo.Route, requestID string) {
	d.Set("destination_service_id", route.GetDestinationServiceId())
	d.Set("gateway_id", route.GetGatewayId())
	d.Set("vm_id", route.GetVmId())
	d.Set("nat_access_point", route.GetNetAccessPointId())
	d.Set("nic_id", route.NicId)
	d.Set("net_peering_id", route.GetNetPeeringId())
	d.Set("vm_account_id", route.GetVmAccountId())
	d.Set("creation_method", route.GetCreationMethod())
	d.Set("state", route.GetState())
	d.Set("request_id", requestID)
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
	numTargets, target := getTarget(d)

	replaceOpts := &oscgo.UpdateRouteRequest{}

	switch target {
	case "vm_id":
		if numTargets > 2 || (numTargets == 2 && len(d.Get("nic_id").(string)) == 0) {
			return errOAPIRoute
		}
	default:
		if numTargets > 1 {
			return errOAPIRoute
		}
	}

	switch target {
	case "gateway_id":
		replaceOpts = &oscgo.UpdateRouteRequest{
			RouteTableId:       d.Get("route_table_id").(string),
			DestinationIpRange: d.Get("destination_ip_range").(string),
			GatewayId:          d.Get("gateway_id").(*string),
		}
	case "nat_service_id":
		replaceOpts = &oscgo.UpdateRouteRequest{
			RouteTableId:       d.Get("route_table_id").(string),
			DestinationIpRange: d.Get("destination_ip_range").(string),
			GatewayId:          d.Get("nat_service_id").(*string),
		}
	case "vm_id":
		replaceOpts = &oscgo.UpdateRouteRequest{
			RouteTableId:       d.Get("route_table_id").(string),
			DestinationIpRange: d.Get("destination_ip_range").(string),
			VmId:               d.Get("vm_id").(*string),
		}
	case "nic_id":
		replaceOpts = &oscgo.UpdateRouteRequest{
			RouteTableId:       d.Get("route_table_id").(string),
			DestinationIpRange: d.Get("destination_ip_range").(string),
			NicId:              d.Get("nic_id").(*string),
		}
	case "net_peering_id":
		replaceOpts = &oscgo.UpdateRouteRequest{
			RouteTableId:       d.Get("route_table_id").(string),
			DestinationIpRange: d.Get("destination_ip_range").(string),
			NetPeeringId:       d.Get("net_peering_id").(*string),
		}
	default:
		return fmt.Errorf("An invalid target type specified: %s", target)
	}
	log.Printf("[DEBUG] Route replace config: %+v", replaceOpts)

	var err error
	err = resource.Retry(2*time.Minute, func() *resource.RetryError {
		_, _, err = conn.RouteApi.UpdateRoute(context.Background(), &oscgo.UpdateRouteOpts{UpdateRouteRequest: optional.NewInterface(replaceOpts)})

		if err != nil {
			if strings.Contains(fmt.Sprint(err), "InvalidParameterException") {
				log.Printf("[DEBUG] Trying to create route again: %q", err)
				return resource.RetryableError(err)
			}

			return resource.NonRetryableError(err)
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
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
		resp, _, err := conn.RouteApi.DeleteRoute(context.Background(), &oscgo.DeleteRouteOpts{DeleteRouteRequest: optional.NewInterface(deleteOpts)})
		log.Printf("[DEBUG] Route delete result: %+v", resp)

		if err == nil {
			return nil
		}

		if strings.Contains(fmt.Sprint(err), "InvalidParameterException") {
			log.Printf("[DEBUG] Trying to delete route again: %q", fmt.Sprint(err))
			return resource.RetryableError(err)
		}

		return resource.NonRetryableError(err)
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
		resp, _, err = conn.RouteTableApi.ReadRouteTables(context.Background(), &oscgo.ReadRouteTablesOpts{ReadRouteTablesRequest: optional.NewInterface(findOpts)})

		if err != nil {
			if strings.Contains(fmt.Sprint(err), "InvalidParameterException") || strings.Contains(fmt.Sprint(err), "RequestLimitExceeded") {
				log.Printf("[DEBUG] Trying to create route again: %q", err)
				return resource.RetryableError(err)
			}

			return resource.NonRetryableError(err)
		}

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

func routeOAPIIDHash(d *schema.ResourceData, r *oscgo.Route) string {
	return fmt.Sprintf("r-%s%d", d.Get("route_table_id").(string), hashcode.String(r.GetDestinationIpRange()))
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
		resp, _, err = conn.RouteTableApi.ReadRouteTables(context.Background(), &oscgo.ReadRouteTablesOpts{ReadRouteTablesRequest: optional.NewInterface(findOpts)})

		if err != nil {
			if strings.Contains(fmt.Sprint(err), "InvalidParameterException") || strings.Contains(fmt.Sprint(err), "RequestLimitExceeded") {
				log.Printf("[DEBUG] Trying to create route again: %q", err)
				return resource.RetryableError(err)
			}

			return resource.NonRetryableError(err)
		}

		return nil
	})

	var errString string

	if err != nil {
		errString = err.Error()
		return nil, "", fmt.Errorf("Error finding route resource: %s", errString)
	}
	requestID := resp.ResponseContext.GetRequestId()

	if len(resp.GetRouteTables()) < 1 || reflect.DeepEqual(resp.GetRouteTables()[0], oscgo.RouteTable{}) {
		return nil, requestID, fmt.Errorf("Route Table %q is gone, or route does not exist", routeTableID)
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
