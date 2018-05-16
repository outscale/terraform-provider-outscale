package outscale

import (
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
			"destination_ip_range": {
				Type:     schema.TypeString,
				Required: true,
			},

			"vpn_gateway_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"vm_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"nat_service_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"nic_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"route_table_id": {
				Type:     schema.TypeString,
				Required: true,
			},

			"lin_peering_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"destinaton_prefix_list_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"vm_account_id": {
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
		},
	}
}

func resourceOutscaleOAPIRouteCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU
	var numTargets int
	var setTarget string
	allowedTargets := []string{
		"vpn_gateway_id",
		"nat_service_id",
		"vm_id",
		"nic_id",
		"lin_peering_id",
	}

	for _, target := range allowedTargets {
		if len(d.Get(target).(string)) > 0 {
			numTargets++
			setTarget = target
		}
	}

	if numTargets > 1 {
		return errRoute
	}

	createOpts := &fcu.CreateRouteInput{}
	switch setTarget {
	case "vpn_gateway_id":
		createOpts = &fcu.CreateRouteInput{
			RouteTableId:         aws.String(d.Get("route_table_id").(string)),
			DestinationCidrBlock: aws.String(d.Get("destination_ip_range").(string)),
			GatewayId:            aws.String(d.Get("vpn_gateway_id").(string)),
		}
	case "nat_service_id":
		createOpts = &fcu.CreateRouteInput{
			RouteTableId:         aws.String(d.Get("route_table_id").(string)),
			DestinationCidrBlock: aws.String(d.Get("destination_ip_range").(string)),
			NatGatewayId:         aws.String(d.Get("nat_service_id").(string)),
		}
	case "vm_id":
		createOpts = &fcu.CreateRouteInput{
			RouteTableId:         aws.String(d.Get("route_table_id").(string)),
			DestinationCidrBlock: aws.String(d.Get("destination_ip_range").(string)),
			InstanceId:           aws.String(d.Get("vm_id").(string)),
		}
	case "nic_id":
		createOpts = &fcu.CreateRouteInput{
			RouteTableId:         aws.String(d.Get("route_table_id").(string)),
			DestinationCidrBlock: aws.String(d.Get("destination_ip_range").(string)),
			NetworkInterfaceId:   aws.String(d.Get("nic_id").(string)),
		}
	case "lin_peering_id":
		createOpts = &fcu.CreateRouteInput{
			RouteTableId:           aws.String(d.Get("route_table_id").(string)),
			DestinationCidrBlock:   aws.String(d.Get("destination_ip_range").(string)),
			VpcPeeringConnectionId: aws.String(d.Get("lin_peering_id").(string)),
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

	if v, ok := d.GetOk("destination_ip_range"); ok {
		err = resource.Retry(2*time.Minute, func() *resource.RetryError {
			route, err = findResourceOAPIRoute(conn, d.Get("route_table_id").(string), v.(string))
			return resource.RetryableError(err)
		})
		if err != nil {
			return fmt.Errorf("Error finding route after creating it: %s", err)
		}
	}

	d.SetId(routeOAPIIDHash(d, route))
	resourceOutscaleOAPIRouteSetResourceData(d, route)
	return nil
}

func resourceOutscaleOAPIRouteRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU
	routeTableId := d.Get("route_table_id").(string)

	destinationCidrBlock := d.Get("destination_ip_range").(string)

	route, err := findResourceOAPIRoute(conn, routeTableId, destinationCidrBlock)
	if err != nil {
		if strings.Contains(fmt.Sprint(err), "InvalidRouteTableID.NotFound") {
			log.Printf("[WARN] Route Table %q could not be found. Removing Route from state.", routeTableId)
			d.SetId("")
			return nil
		}
		return err
	}
	resourceOutscaleOAPIRouteSetResourceData(d, route)
	return nil
}

func resourceOutscaleOAPIRouteSetResourceData(d *schema.ResourceData, route *fcu.Route) {
	d.Set("destinaton_prefix_list_id", route.DestinationPrefixListId)
	d.Set("vpn_gateway_id", route.GatewayId)
	d.Set("vm_id", route.InstanceId)
	d.Set("nat_service_id", route.NatGatewayId)
	d.Set("nic_id", route.NetworkInterfaceId)
	d.Set("lin_peering_id", route.VpcPeeringConnectionId)
	d.Set("vm_account_id", route.InstanceOwnerId)
	d.Set("creation_method", route.Origin)
	d.Set("state", route.State)
}

func resourceOutscaleOAPIRouteUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU
	var numTargets int
	var setTarget string

	allowedTargets := []string{
		"vpn_gateway_id",
		"nat_service_id",
		"nic_id",
		"vm_id",
		"lin_peering_id",
	}
	replaceOpts := &fcu.ReplaceRouteInput{}

	for _, target := range allowedTargets {
		if len(d.Get(target).(string)) > 0 {
			numTargets++
			setTarget = target
		}
	}

	switch setTarget {
	case "vm_id":
		if numTargets > 2 || (numTargets == 2 && len(d.Get("nic_id").(string)) == 0) {
			return errRoute
		}
	default:
		if numTargets > 1 {
			return errRoute
		}
	}

	switch setTarget {
	case "vpn_gateway_id":
		replaceOpts = &fcu.ReplaceRouteInput{
			RouteTableId:         aws.String(d.Get("route_table_id").(string)),
			DestinationCidrBlock: aws.String(d.Get("destination_ip_range").(string)),
			GatewayId:            aws.String(d.Get("vpn_gateway_id").(string)),
		}
	case "nat_service_id":
		replaceOpts = &fcu.ReplaceRouteInput{
			RouteTableId:         aws.String(d.Get("route_table_id").(string)),
			DestinationCidrBlock: aws.String(d.Get("destination_ip_range").(string)),
			NatGatewayId:         aws.String(d.Get("nat_service_id").(string)),
		}
	case "vm_id":
		replaceOpts = &fcu.ReplaceRouteInput{
			RouteTableId:         aws.String(d.Get("route_table_id").(string)),
			DestinationCidrBlock: aws.String(d.Get("destination_ip_range").(string)),
			InstanceId:           aws.String(d.Get("vm_id").(string)),
		}
	case "nic_id":
		replaceOpts = &fcu.ReplaceRouteInput{
			RouteTableId:         aws.String(d.Get("route_table_id").(string)),
			DestinationCidrBlock: aws.String(d.Get("destination_ip_range").(string)),
			NetworkInterfaceId:   aws.String(d.Get("nic_id").(string)),
		}
	case "lin_peering_id":
		replaceOpts = &fcu.ReplaceRouteInput{
			RouteTableId:           aws.String(d.Get("route_table_id").(string)),
			DestinationCidrBlock:   aws.String(d.Get("destination_ip_range").(string)),
			VpcPeeringConnectionId: aws.String(d.Get("lin_peering_id").(string)),
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

func resourceOutscaleOAPIRouteDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	deleteOpts := &fcu.DeleteRouteInput{
		RouteTableId: aws.String(d.Get("route_table_id").(string)),
	}
	if v, ok := d.GetOk("destination_ip_range"); ok {
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

func resourceOutscaleOAPIRouteExists(d *schema.ResourceData, meta interface{}) (bool, error) {
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

	if v, ok := d.GetOk("destination_ip_range"); ok {
		for _, route := range (*res.RouteTables[0]).Routes {
			if route.DestinationCidrBlock != nil && *route.DestinationCidrBlock == v.(string) {
				return true, nil
			}
		}
	}

	return false, nil
}

func routeOAPIIDHash(d *schema.ResourceData, r *fcu.Route) string {
	return fmt.Sprintf("r-%s%d", d.Get("route_table_id").(string), hashcode.String(*r.DestinationCidrBlock))
}

func findResourceOAPIRoute(conn *fcu.Client, rtbid string, cidr string) (*fcu.Route, error) {
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
