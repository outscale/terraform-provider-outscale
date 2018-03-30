package outscale

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-outscale/osc/fcu"
)

var routeTargetValidationError = errors.New("Error: more than 1 target specified. Only 1 of gateway_id, " +
	"egress_only_gateway_id, nat_gateway_id, instance_id, network_interface_id, route_table_id or " +
	"vpc_peering_connection_id is allowed.")

func resourceOutscaleRoute() *schema.Resource {
	return &schema.Resource{
		Create: resourceOutscaleRouteCreate,
		Read:   resourceOutscaleRouteRead,
		Update: resourceOutscaleRouteUpdate,
		Delete: resourceOutscaleRouteDelete,
		Exists: resourceOutscaleRouteExists,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"destination_cidr_block": {
				Type:     schema.TypeString,
				Required: true,
			},

			"gateway_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"instance_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"nat_gateway_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"network_interface_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"route_table_id": {
				Type:     schema.TypeString,
				Required: true,
			},

			"vpc_peering_connection_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"destination_prefix_list_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"instance_owner_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"origin": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceOutscaleRouteCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU
	var numTargets int
	var setTarget string
	allowedTargets := []string{
		"gateway_id",
		"nat_gateway_id",
		"instance_id",
		"network_interface_id",
		"vpc_peering_connection_id",
	}

	for _, target := range allowedTargets {
		if len(d.Get(target).(string)) > 0 {
			numTargets++
			setTarget = target
		}
	}

	if numTargets > 1 {
		return routeTargetValidationError
	}

	createOpts := &fcu.CreateRouteInput{}
	switch setTarget {
	case "gateway_id":
		createOpts = &fcu.CreateRouteInput{
			RouteTableId:         aws.String(d.Get("route_table_id").(string)),
			DestinationCidrBlock: aws.String(d.Get("destination_cidr_block").(string)),
			GatewayId:            aws.String(d.Get("gateway_id").(string)),
		}
	case "nat_gateway_id":
		createOpts = &fcu.CreateRouteInput{
			RouteTableId:         aws.String(d.Get("route_table_id").(string)),
			DestinationCidrBlock: aws.String(d.Get("destination_cidr_block").(string)),
			NatGatewayId:         aws.String(d.Get("nat_gateway_id").(string)),
		}
	case "instance_id":
		createOpts = &fcu.CreateRouteInput{
			RouteTableId:         aws.String(d.Get("route_table_id").(string)),
			DestinationCidrBlock: aws.String(d.Get("destination_cidr_block").(string)),
			InstanceId:           aws.String(d.Get("instance_id").(string)),
		}
	case "network_interface_id":
		createOpts = &fcu.CreateRouteInput{
			RouteTableId:         aws.String(d.Get("route_table_id").(string)),
			DestinationCidrBlock: aws.String(d.Get("destination_cidr_block").(string)),
			NetworkInterfaceId:   aws.String(d.Get("network_interface_id").(string)),
		}
	case "vpc_peering_connection_id":
		createOpts = &fcu.CreateRouteInput{
			RouteTableId:           aws.String(d.Get("route_table_id").(string)),
			DestinationCidrBlock:   aws.String(d.Get("destination_cidr_block").(string)),
			VpcPeeringConnectionId: aws.String(d.Get("vpc_peering_connection_id").(string)),
		}
	default:
		return fmt.Errorf("An invalid target type specified: %s", setTarget)
	}
	log.Printf("[DEBUG] Route create config: %s", createOpts)

	var err error
	err = resource.Retry(2*time.Minute, func() *resource.RetryError {
		_, err = conn.VM.CreateRoute(createOpts)

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
		return fmt.Errorf("Error creating route: %s", err)
	}

	var route *fcu.Route

	if v, ok := d.GetOk("destination_cidr_block"); ok {
		err = resource.Retry(2*time.Minute, func() *resource.RetryError {
			route, err = findResourceRoute(conn, d.Get("route_table_id").(string), v.(string))
			return resource.RetryableError(err)
		})
		if err != nil {
			return fmt.Errorf("Error finding route after creating it: %s", err)
		}
	}

	d.SetId(routeIDHash(d, route))
	resourceOutscaleRouteSetResourceData(d, route)
	return nil
}

func resourceOutscaleRouteRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU
	routeTableId := d.Get("route_table_id").(string)

	destinationCidrBlock := d.Get("destination_cidr_block").(string)

	route, err := findResourceRoute(conn, routeTableId, destinationCidrBlock)
	if err != nil {
		if strings.Contains(fmt.Sprint(err), "InvalidRouteTableID.NotFound") {
			log.Printf("[WARN] Route Table %q could not be found. Removing Route from state.", routeTableId)
			d.SetId("")
			return nil
		}
		return err
	}
	resourceOutscaleRouteSetResourceData(d, route)
	return nil
}

func resourceOutscaleRouteSetResourceData(d *schema.ResourceData, route *fcu.Route) {
	d.Set("destination_prefix_list_id", route.DestinationPrefixListId)
	d.Set("gateway_id", route.GatewayId)
	d.Set("instance_id", route.InstanceId)
	d.Set("nat_gateway_id", route.NatGatewayId)
	d.Set("network_interface_id", route.NetworkInterfaceId)
	d.Set("vpc_peering_connection_id", route.VpcPeeringConnectionId)
	d.Set("instance_owner_id", route.InstanceOwnerId)
	d.Set("origin", route.Origin)
	d.Set("state", route.State)
}

func resourceOutscaleRouteUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU
	var numTargets int
	var setTarget string

	allowedTargets := []string{
		"gateway_id",
		"nat_gateway_id",
		"network_interface_id",
		"instance_id",
		"vpc_peering_connection_id",
	}
	replaceOpts := &fcu.ReplaceRouteInput{}

	for _, target := range allowedTargets {
		if len(d.Get(target).(string)) > 0 {
			numTargets++
			setTarget = target
		}
	}

	switch setTarget {
	case "instance_id":
		if numTargets > 2 || (numTargets == 2 && len(d.Get("network_interface_id").(string)) == 0) {
			return routeTargetValidationError
		}
	default:
		if numTargets > 1 {
			return routeTargetValidationError
		}
	}

	switch setTarget {
	case "gateway_id":
		replaceOpts = &fcu.ReplaceRouteInput{
			RouteTableId:         aws.String(d.Get("route_table_id").(string)),
			DestinationCidrBlock: aws.String(d.Get("destination_cidr_block").(string)),
			GatewayId:            aws.String(d.Get("gateway_id").(string)),
		}
	case "nat_gateway_id":
		replaceOpts = &fcu.ReplaceRouteInput{
			RouteTableId:         aws.String(d.Get("route_table_id").(string)),
			DestinationCidrBlock: aws.String(d.Get("destination_cidr_block").(string)),
			NatGatewayId:         aws.String(d.Get("nat_gateway_id").(string)),
		}
	case "instance_id":
		replaceOpts = &fcu.ReplaceRouteInput{
			RouteTableId:         aws.String(d.Get("route_table_id").(string)),
			DestinationCidrBlock: aws.String(d.Get("destination_cidr_block").(string)),
			InstanceId:           aws.String(d.Get("instance_id").(string)),
		}
	case "network_interface_id":
		replaceOpts = &fcu.ReplaceRouteInput{
			RouteTableId:         aws.String(d.Get("route_table_id").(string)),
			DestinationCidrBlock: aws.String(d.Get("destination_cidr_block").(string)),
			NetworkInterfaceId:   aws.String(d.Get("network_interface_id").(string)),
		}
	case "vpc_peering_connection_id":
		replaceOpts = &fcu.ReplaceRouteInput{
			RouteTableId:           aws.String(d.Get("route_table_id").(string)),
			DestinationCidrBlock:   aws.String(d.Get("destination_cidr_block").(string)),
			VpcPeeringConnectionId: aws.String(d.Get("vpc_peering_connection_id").(string)),
		}
	default:
		return fmt.Errorf("An invalid target type specified: %s", setTarget)
	}
	log.Printf("[DEBUG] Route replace config: %s", replaceOpts)

	var err error
	err = resource.Retry(2*time.Minute, func() *resource.RetryError {
		_, err = conn.VM.ReplaceRoute(replaceOpts)

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

func resourceOutscaleRouteDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	deleteOpts := &fcu.DeleteRouteInput{
		RouteTableId: aws.String(d.Get("route_table_id").(string)),
	}
	if v, ok := d.GetOk("destination_cidr_block"); ok {
		deleteOpts.DestinationCidrBlock = aws.String(v.(string))
	}
	log.Printf("[DEBUG] Route delete opts: %s", deleteOpts)

	var err error
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		log.Printf("[DEBUG] Trying to delete route with opts %s", deleteOpts)
		resp, err := conn.VM.DeleteRoute(deleteOpts)
		log.Printf("[DEBUG] Route delete result: %s", resp)

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

func resourceOutscaleRouteExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	conn := meta.(*OutscaleClient).FCU
	routeTableId := d.Get("route_table_id").(string)

	findOpts := &fcu.DescribeRouteTablesInput{
		RouteTableIds: []*string{&routeTableId},
	}

	var res *fcu.DescribeRouteTablesOutput
	var err error
	err = resource.Retry(2*time.Minute, func() *resource.RetryError {
		res, err = conn.VM.DescribeRouteTables(findOpts)

		if err != nil {
			if strings.Contains(fmt.Sprint(err), "InvalidParameterException") || strings.Contains(fmt.Sprint(err), "RequestLimitExceeded") {
				log.Printf("[DEBUG] Trying to create route again: %q", err)
				return resource.RetryableError(err)
			}

			return resource.NonRetryableError(err)
		}

		return nil
	})

	if err != nil {
		if strings.Contains(fmt.Sprint(err), "InvalidRouteTableID.NotFound") {
			log.Printf("[WARN] Route Table %q could not be found.", routeTableId)
			return false, nil
		}
		return false, fmt.Errorf("Error while checking if route exists: %s", err)
	}

	if len(res.RouteTables) < 1 || res.RouteTables[0] == nil {
		log.Printf("[WARN] Route Table %q is gone, or route does not exist.",
			routeTableId)
		return false, nil
	}

	if v, ok := d.GetOk("destination_cidr_block"); ok {
		for _, route := range (*res.RouteTables[0]).Routes {
			if route.DestinationCidrBlock != nil && *route.DestinationCidrBlock == v.(string) {
				return true, nil
			}
		}
	}

	return false, nil
}

func routeIDHash(d *schema.ResourceData, r *fcu.Route) string {
	return fmt.Sprintf("r-%s%d", d.Get("route_table_id").(string), hashcode.String(*r.DestinationCidrBlock))
}

func findResourceRoute(conn *fcu.Client, rtbid string, cidr string) (*fcu.Route, error) {
	routeTableID := rtbid

	findOpts := &fcu.DescribeRouteTablesInput{
		RouteTableIds: []*string{&routeTableID},
	}

	var resp *fcu.DescribeRouteTablesOutput
	var err error
	err = resource.Retry(2*time.Minute, func() *resource.RetryError {
		resp, err = conn.VM.DescribeRouteTables(findOpts)

		if err != nil {
			if strings.Contains(fmt.Sprint(err), "InvalidParameterException") || strings.Contains(fmt.Sprint(err), "RequestLimitExceeded") {
				log.Printf("[DEBUG] Trying to create route again: %q", err)
				return resource.RetryableError(err)
			}

			return resource.NonRetryableError(err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	if len(resp.RouteTables) < 1 || resp.RouteTables[0] == nil {
		return nil, fmt.Errorf("Route Table %q is gone, or route does not exist.",
			routeTableID)
	}

	if cidr != "" {
		for _, route := range (*resp.RouteTables[0]).Routes {
			if route.DestinationCidrBlock != nil && *route.DestinationCidrBlock == cidr {
				return route, nil
			}
		}

		return nil, fmt.Errorf("Unable to find matching route for Route Table (%s) "+
			"and destination CIDR block (%s).", rtbid, cidr)
	}

	return nil, fmt.Errorf("When trying to find a matching route for Route Table %q "+
		"you need to specify a CIDR block", rtbid)

}
