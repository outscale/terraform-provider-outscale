package outscale

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-outscale/osc/fcu"
)

// Route table import also imports all the rules
func resourceOutscaleRouteTableImportState(
	d *schema.ResourceData,
	meta interface{}) ([]*schema.ResourceData, error) {
	conn := meta.(*OutscaleClient).FCU

	// First query the resource itself
	id := d.Id()
	resp, err := conn.VM.DescribeRouteTables(&fcu.DescribeRouteTablesInput{
		RouteTableIds: []*string{&id},
	})
	if err != nil {
		return nil, err
	}
	if len(resp.RouteTables) < 1 || resp.RouteTables[0] == nil {
		return nil, fmt.Errorf("route table %s is not found", id)
	}
	table := resp.RouteTables[0]

	// Start building our results
	results := make([]*schema.ResourceData, 1,
		2+len(table.Associations)+len(table.Routes))
	results[0] = d

	{
		// Construct the routes
		subResource := resourceOutscaleRoute()
		for _, route := range table.Routes {
			// Ignore the local/default route
			if route.GatewayId != nil && *route.GatewayId == "local" {
				continue
			}

			if route.DestinationPrefixListId != nil {
				// Skipping because VPC endpoint routes are handled separately
				// See aws_vpc_endpoint
				continue
			}

			// Minimal data for route
			d := subResource.Data(nil)
			d.SetType("outscale_route")
			d.Set("route_table_id", id)
			d.Set("destination_cidr_block", route.DestinationCidrBlock)
			d.SetId(routeIDHash(d, route))
			results = append(results, d)
		}
	}

	{
		// Construct the associations
		subResource := resourceOutscaleRouteTableAssociation()
		for _, assoc := range table.Associations {
			if *assoc.Main {
				// Ignore
				continue
			}

			// Minimal data for route
			d := subResource.Data(nil)
			d.SetType("outscale_route_table_link")
			d.Set("route_table_id", assoc.RouteTableId)
			d.SetId(*assoc.RouteTableAssociationId)
			results = append(results, d)
		}
	}

	return results, nil
}
