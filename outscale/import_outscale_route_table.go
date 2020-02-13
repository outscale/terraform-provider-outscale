package outscale

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/hashcode"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	oscgo "github.com/marinsalinas/osc-sdk-go"
)

func routeIDHash(d *schema.ResourceData, r *oscgo.Route) string {
	return fmt.Sprintf("r-%s%d", d.Get("route_table_id").(string),
		hashcode.String(r.GetDestinationIpRange()))
}

// Route table import also imports all the rules
func resourceOutscaleRouteTableImportState(
	d *schema.ResourceData,
	meta interface{}) ([]*schema.ResourceData, error) {
	conn := meta.(*OutscaleClient).OSCAPI

	// First query the resource itself
	id := d.Id()
	tableRaw, _, _ := readOAPIRouteTable(conn, id)

	table := tableRaw.(oscgo.RouteTable)
	// Start building our results
	results := make([]*schema.ResourceData, 1,
		2+len(table.GetLinkRouteTables())+len(table.GetRoutes()))
	results[0] = d

	{
		// Construct the routes
		subResource := resourceOutscaleOAPIRoute()
		for _, route := range table.GetRoutes() {
			// Ignore the local/default route
			if route.GatewayId != nil && *route.GatewayId == "local" {
				continue
			}

			if route.DestinationServiceId != nil {
				// Skipping because VPC endpoint routes are handled separately
				// See aws_vpc_endpoint
				continue
			}

			// Minimal data for route
			d := subResource.Data(nil)
			d.SetType("outscale_route")
			d.Set("route_table_id", id)
			d.Set("destination_cidr_block", route.DestinationIpRange)
			d.SetId(routeIDHash(d, &route))
			results = append(results, d)
		}
	}

	{
		// Construct the associations
		subResource := resourceOutscaleOAPILinkRouteTable()
		for _, assoc := range table.GetLinkRouteTables() {
			if *assoc.Main {
				// Ignore
				continue
			}

			// Minimal data for route
			d := subResource.Data(nil)
			d.SetType("outscale_route_table_link")
			d.Set("route_table_id", assoc.RouteTableId)
			d.SetId(*assoc.LinkRouteTableId)
			results = append(results, d)
		}
	}

	return results, nil
}
