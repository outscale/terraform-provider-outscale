package outscale

import (
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-outscale/osc/fcu"
)

func resourceOutscaleRouteTable() *schema.Resource {
	return &schema.Resource{
		Create: resourceOutscaleRouteTableCreate,
		Read:   resourceOutscaleRouteTableRead,
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

			"propagating_vgw_set": {
				Type:     schema.TypeList,
				Computed: true,
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

	createOpts := &fcu.CreateRouteTableInput{
		VpcId: aws.String(d.Get("vpc_id").(string)),
	}

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

	rt := resp.RouteTable
	d.SetId(*rt.RouteTableId)

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

	if d.IsNewResource() {
		if err := setTags(conn, d); err != nil {
			return err
		}
		d.SetPartial("tag_set")
	}

	d.Set("tag_set", make([]interface{}, 0))
	d.Set("route_set", make([]interface{}, 0))
	d.Set("association_set", make([]interface{}, 0))

	return resourceOutscaleRouteTableRead(d, meta)
}

func resourceOutscaleRouteTableRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	var resp *fcu.DescribeRouteTablesOutput
	var err error
	err = resource.Retry(15*time.Minute, func() *resource.RetryError {
		resp, err = conn.VM.DescribeRouteTables(&fcu.DescribeRouteTablesInput{
			RouteTableIds: []*string{aws.String(d.Id())},
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
			return err
		}
	}

	if resp == nil {
		d.SetId("")
		return nil
	}

	rt := resp.RouteTables[0]

	propagatingVGWs := make([]string, 0, len(rt.PropagatingVgws))
	for _, vgw := range rt.PropagatingVgws {
		propagatingVGWs = append(propagatingVGWs, *vgw.GatewayId)
	}

	d.Set("propagating_vgw_set", propagatingVGWs)
	d.Set("tag_set", tagsToMap(rt.Tags))
	d.Set("request_id", resp.RequestId)
	d.Set("route_table_id", rt.RouteTableId)
	d.Set("vpc_id", rt.VpcId)

	if err := d.Set("route_set", setRouteSet(rt.Routes)); err != nil {
		return err
	}

	return d.Set("association_set", setAssociactionSet(rt.Associations))
}

func resourceOutscaleRouteTableDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	rtRaw, _, err := resourceOutscaleRouteTableStateRefreshFunc(conn, d.Id())()
	if err != nil {
		return err
	}
	if rtRaw == nil {
		return nil
	}
	rt := rtRaw.(*fcu.RouteTable)

	for _, a := range rt.Associations {

		var err error
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			_, err := conn.VM.DisassociateRouteTable(&fcu.DisassociateRouteTableInput{
				AssociationId: a.RouteTableAssociationId,
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
				err = nil
			}
			return err
		}
	}

	err = resource.Retry(15*time.Minute, func() *resource.RetryError {
		_, err = conn.VM.DeleteRouteTable(&fcu.DeleteRouteTableInput{
			RouteTableId: aws.String(d.Id()),
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
			return nil
		}

		return fmt.Errorf("Error deleting route table: %s", err)
	}

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

func resourceOutscaleRouteTableStateRefreshFunc(conn *fcu.Client, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {

		var resp *fcu.DescribeRouteTablesOutput
		var err error
		err = resource.Retry(15*time.Minute, func() *resource.RetryError {
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

func setRouteSet(rt []*fcu.Route) []map[string]interface{} {

	route := make([]map[string]interface{}, len(rt))

	for k, r := range rt {
		m := make(map[string]interface{})

		m["destination_cidr_block"] = aws.StringValue(r.DestinationCidrBlock)
		m["destination_prefix_list_id"] = aws.StringValue(r.DestinationPrefixListId)
		m["gateway_id"] = aws.StringValue(r.GatewayId)
		m["instance_id"] = aws.StringValue(r.InstanceId)
		m["instance_owner_id"] = aws.StringValue(r.InstanceOwnerId)
		m["vpc_peering_connection_id"] = aws.StringValue(r.VpcPeeringConnectionId)
		m["network_interface_id"] = aws.StringValue(r.NetworkInterfaceId)
		m["origin"] = aws.StringValue(r.Origin)
		m["state"] = aws.StringValue(r.State)

		route[k] = m
	}

	return route
}

func setAssociactionSet(rt []*fcu.RouteTableAssociation) []map[string]interface{} {
	association := make([]map[string]interface{}, len(rt))

	for k, r := range rt {
		m := make(map[string]interface{})
		m["main"] = aws.BoolValue(r.Main)
		m["route_table_association_id"] = aws.StringValue(r.RouteTableAssociationId)
		m["route_table_id"] = aws.StringValue(r.RouteTableId)
		m["subnet_id"] = aws.StringValue(r.SubnetId)

		association[k] = m
	}

	return association
}
