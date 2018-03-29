package outscale

import (
	"bytes"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-outscale/osc/fcu"
)

func resourceOutscaleRouteTable() *schema.Resource {
	return &schema.Resource{
		Create: resourceOutscaleRouteTableCreate,
		Read:   resourceOutscaleRouteTableRead,
		Update: resourceOutscaleRouteTableUpdate,
		Delete: resourceOutscaleRouteTableDelete,
		Importer: &schema.ResourceImporter{
			State: resourceOutscaleRouteTableImportState,
		},

		Schema: map[string]*schema.Schema{
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"route_table_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"request_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"tag":     tagsSchema(),
			"tag_set": tagsSchemaComputed(),

			"propagating_vgws": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"route_set": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"destination_cidr_block": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"destination_prefix_list_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"gateway_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_owner_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"network_interface_id": {
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
						"vpc_peering_connection_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"association_set": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"main": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"route_table_association_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"route_table_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"subnet_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func resourceOutscaleRouteTableCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	// Create the routing table
	createOpts := &fcu.CreateRouteTableInput{
		VpcId: aws.String(d.Get("vpc_id").(string)),
	}
	log.Printf("[DEBUG] RouteTable create config: %#v", createOpts)

	var resp *fcu.CreateRouteTableOutput
	var err error
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, err = conn.VM.CreateRouteTable(createOpts)
		if err != nil {
			if strings.Contains(fmt.Sprint(err), "RequestLimitExceeded") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("Error creating route table: %s", err)
	}

	// Get the ID and store it
	rt := resp.RouteTable
	d.SetId(*rt.RouteTableId)
	log.Printf("[INFO] Route Table ID: %s", d.Id())

	// Wait for the route table to become available
	log.Printf(
		"[DEBUG] Waiting for route table (%s) to become available",
		d.Id())
	stateConf := &resource.StateChangeConf{
		Pending: []string{"pending"},
		Target:  []string{"ready"},
		Refresh: resourceOutscaleRouteTableStateRefreshFunc(conn, d.Id()),
		Timeout: 5 * time.Minute,
	}
	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf(
			"Error waiting for route table (%s) to become available: %s",
			d.Id(), err)
	}

	return resourceOutscaleRouteTableUpdate(d, meta)
}

func resourceOutscaleRouteTableRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	rtRaw, _, err := resourceOutscaleRouteTableStateRefreshFunc(conn, d.Id())()
	if err != nil {
		return err
	}
	if rtRaw == nil {
		d.SetId("")
		return nil
	}

	rt := rtRaw.(*fcu.RouteTable)
	d.Set("vpc_id", rt.VpcId)

	propagatingVGWs := make([]string, 0, len(rt.PropagatingVgws))
	for _, vgw := range rt.PropagatingVgws {
		propagatingVGWs = append(propagatingVGWs, *vgw.GatewayId)
	}
	d.Set("propagating_vgws", propagatingVGWs)

	// Create an empty schema.Set to hold all routes
	route := &schema.Set{F: resourceOutscaleRouteTableHash}

	// Loop through the routes and add them to the set
	for _, r := range rt.Routes {
		if r.GatewayId != nil && *r.GatewayId == "local" {
			continue
		}

		if r.Origin != nil && *r.Origin == "EnableVgwRoutePropagation" {
			continue
		}

		if r.DestinationPrefixListId != nil {
			// Skipping because VPC endpoint routes are handled separately
			// See aws_vpc_endpoint
			continue
		}

		m := make(map[string]interface{})

		if r.DestinationCidrBlock != nil {
			m["cidr_block"] = *r.DestinationCidrBlock
		}
		if r.GatewayId != nil {
			m["gateway_id"] = *r.GatewayId
		}
		if r.NatGatewayId != nil {
			m["nat_gateway_id"] = *r.NatGatewayId
		}
		if r.InstanceId != nil {
			m["instance_id"] = *r.InstanceId
		}
		if r.VpcPeeringConnectionId != nil {
			m["vpc_peering_connection_id"] = *r.VpcPeeringConnectionId
		}
		if r.NetworkInterfaceId != nil {
			m["network_interface_id"] = *r.NetworkInterfaceId
		}

		route.Add(m)
	}
	d.Set("route", route)

	// Tags
	d.Set("tags", tagsToMap(rt.Tags))

	return nil
}

func resourceOutscaleRouteTableUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	if d.HasChange("propagating_vgws") {
		o, n := d.GetChange("propagating_vgws")
		os := o.(*schema.Set)
		ns := n.(*schema.Set)
		remove := os.Difference(ns).List()
		add := ns.Difference(os).List()

		// Now first loop through all the old propagations and disable any obsolete ones
		for _, vgw := range remove {
			id := vgw.(string)

			// Disable the propagation as it no longer exists in the config
			log.Printf(
				"[INFO] Deleting VGW propagation from %s: %s",
				d.Id(), id)
			_, err := conn.VM.DisableVgwRoutePropagation(&fcu.DisableVgwRoutePropagationInput{
				RouteTableId: aws.String(d.Id()),
				GatewayId:    aws.String(id),
			})
			if err != nil {
				return err
			}
		}

		// Make sure we save the state of the currently configured rules
		propagatingVGWs := os.Intersection(ns)
		d.Set("propagating_vgws", propagatingVGWs)

		// Then loop through all the newly configured propagations and enable them
		for _, vgw := range add {
			id := vgw.(string)

			var err error
			for i := 0; i < 5; i++ {
				log.Printf("[INFO] Enabling VGW propagation for %s: %s", d.Id(), id)
				_, err = conn.VM.EnableVgwRoutePropagation(&fcu.EnableVgwRoutePropagationInput{
					RouteTableId: aws.String(d.Id()),
					GatewayId:    aws.String(id),
				})
				if err == nil {
					break
				}

				// If we get a Gateway.NotAttached, it is usually some
				// eventually consistency stuff. So we have to just wait a
				// bit...
				ec2err, ok := err.(awserr.Error)
				if ok && ec2err.Code() == "Gateway.NotAttached" {
					time.Sleep(20 * time.Second)
					continue
				}
			}
			if err != nil {
				return err
			}

			propagatingVGWs.Add(vgw)
			d.Set("propagating_vgws", propagatingVGWs)
		}
	}

	// Check if the route set as a whole has changed
	if d.HasChange("route") {
		o, n := d.GetChange("route")
		ors := o.(*schema.Set).Difference(n.(*schema.Set))
		nrs := n.(*schema.Set).Difference(o.(*schema.Set))

		// Now first loop through all the old routes and delete any obsolete ones
		for _, route := range ors.List() {
			m := route.(map[string]interface{})

			deleteOpts := &fcu.DeleteRouteInput{
				RouteTableId: aws.String(d.Id()),
			}

			if s := m["ipv6_cidr_block"].(string); s != "" {
				deleteOpts.DestinationIpv6CidrBlock = aws.String(s)

				log.Printf(
					"[INFO] Deleting route from %s: %s",
					d.Id(), m["ipv6_cidr_block"].(string))
			}

			if s := m["cidr_block"].(string); s != "" {
				deleteOpts.DestinationCidrBlock = aws.String(s)

				log.Printf(
					"[INFO] Deleting route from %s: %s",
					d.Id(), m["cidr_block"].(string))
			}

			_, err := conn.VM.DeleteRoute(deleteOpts)
			if err != nil {
				return err
			}
		}

		// Make sure we save the state of the currently configured rules
		routes := o.(*schema.Set).Intersection(n.(*schema.Set))
		d.Set("route", routes)

		// Then loop through all the newly configured routes and create them
		for _, route := range nrs.List() {
			m := route.(map[string]interface{})

			opts := fcu.CreateRouteInput{
				RouteTableId: aws.String(d.Id()),
			}

			if s := m["vpc_peering_connection_id"].(string); s != "" {
				opts.VpcPeeringConnectionId = aws.String(s)
			}

			if s := m["network_interface_id"].(string); s != "" {
				opts.NetworkInterfaceId = aws.String(s)
			}

			if s := m["instance_id"].(string); s != "" {
				opts.InstanceId = aws.String(s)
			}

			if s := m["ipv6_cidr_block"].(string); s != "" {
				opts.DestinationIpv6CidrBlock = aws.String(s)
			}

			if s := m["cidr_block"].(string); s != "" {
				opts.DestinationCidrBlock = aws.String(s)
			}

			if s := m["gateway_id"].(string); s != "" {
				opts.GatewayId = aws.String(s)
			}

			if s := m["nat_gateway_id"].(string); s != "" {
				opts.NatGatewayId = aws.String(s)
			}

			log.Printf("[INFO] Creating route for %s: %#v", d.Id(), opts)
			err := resource.Retry(5*time.Minute, func() *resource.RetryError {
				_, err := conn.VM.CreateRoute(&opts)
				if err != nil {
					if awsErr, ok := err.(awserr.Error); ok {
						if awsErr.Code() == "InvalidRouteTableID.NotFound" {
							return resource.RetryableError(awsErr)
						}
					}
					return resource.NonRetryableError(err)
				}
				return nil
			})
			if err != nil {
				return err
			}

			routes.Add(route)
			d.Set("route", routes)
		}
	}

	if err := setTags(conn, d); err != nil {
		return err
	} else {
		d.SetPartial("tags")
	}

	return resourceOutscaleRouteTableRead(d, meta)
}

func resourceOutscaleRouteTableDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	// First request the routing table since we'll have to disassociate
	// all the subnets first.
	rtRaw, _, err := resourceOutscaleRouteTableStateRefreshFunc(conn, d.Id())()
	if err != nil {
		return err
	}
	if rtRaw == nil {
		return nil
	}
	rt := rtRaw.(*fcu.RouteTable)

	// Do all the disassociations
	for _, a := range rt.Associations {
		log.Printf("[INFO] Disassociating association: %s", *a.RouteTableAssociationId)
		_, err := conn.VM.DisassociateRouteTable(&fcu.DisassociateRouteTableInput{
			AssociationId: a.RouteTableAssociationId,
		})
		if err != nil {
			// First check if the association ID is not found. If this
			// is the case, then it was already disassociated somehow,
			// and that is okay.
			if ec2err, ok := err.(awserr.Error); ok && ec2err.Code() == "InvalidAssociationID.NotFound" {
				err = nil
			}
		}
		if err != nil {
			return err
		}
	}

	// Delete the route table
	log.Printf("[INFO] Deleting Route Table: %s", d.Id())
	_, err = conn.VM.DeleteRouteTable(&fcu.DeleteRouteTableInput{
		RouteTableId: aws.String(d.Id()),
	})
	if err != nil {
		ec2err, ok := err.(awserr.Error)
		if ok && ec2err.Code() == "InvalidRouteTableID.NotFound" {
			return nil
		}

		return fmt.Errorf("Error deleting route table: %s", err)
	}

	// Wait for the route table to really destroy
	log.Printf(
		"[DEBUG] Waiting for route table (%s) to become destroyed",
		d.Id())

	stateConf := &resource.StateChangeConf{
		Pending: []string{"ready"},
		Target:  []string{},
		Refresh: resourceOutscaleRouteTableStateRefreshFunc(conn, d.Id()),
		Timeout: 5 * time.Minute,
	}
	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf(
			"Error waiting for route table (%s) to become destroyed: %s",
			d.Id(), err)
	}

	return nil
}

func resourceOutscaleRouteTableHash(v interface{}) int {
	var buf bytes.Buffer
	m, castOk := v.(map[string]interface{})
	if !castOk {
		return 0
	}

	if v, ok := m["cidr_block"]; ok {
		buf.WriteString(fmt.Sprintf("%s-", v.(string)))
	}

	if v, ok := m["gateway_id"]; ok {
		buf.WriteString(fmt.Sprintf("%s-", v.(string)))
	}

	natGatewaySet := false
	if v, ok := m["nat_gateway_id"]; ok {
		natGatewaySet = v.(string) != ""
		buf.WriteString(fmt.Sprintf("%s-", v.(string)))
	}

	instanceSet := false
	if v, ok := m["instance_id"]; ok {
		instanceSet = v.(string) != ""
		buf.WriteString(fmt.Sprintf("%s-", v.(string)))
	}

	if v, ok := m["vpc_peering_connection_id"]; ok {
		buf.WriteString(fmt.Sprintf("%s-", v.(string)))
	}

	if v, ok := m["network_interface_id"]; ok && !(instanceSet || natGatewaySet) {
		buf.WriteString(fmt.Sprintf("%s-", v.(string)))
	}

	return hashcode.String(buf.String())
}

func resourceOutscaleRouteTableStateRefreshFunc(conn *fcu.Client, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {

		var resp *fcu.DescribeRouteTablesOutput
		var err error
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			resp, err = conn.VM.DescribeRouteTables(&fcu.DescribeRouteTablesInput{
				RouteTableIds: []*string{aws.String(id)},
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
			if strings.Contains(fmt.Sprint(err), "InvalidRouteTableID.NotFound") {
				resp = nil
			} else {
				log.Printf("Error on RouteTableStateRefresh: %s", err)
				return nil, "", err
			}
		}

		if resp == nil {
			return nil, "", nil
		}

		rt := resp.RouteTables[0]
		return rt, "ready", nil
	}
}
