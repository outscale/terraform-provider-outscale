package outscale

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-outscale/osc/fcu"
)

func resourceOutscaleOAPIRouteTable() *schema.Resource {
	return &schema.Resource{
		Create: resourceOutscaleOAPIRouteTableCreate,
		Read:   resourceOutscaleOAPIRouteTableRead,
		Delete: resourceOutscaleOAPIRouteTableDelete,
		Importer: &schema.ResourceImporter{
			State: resourceOutscaleRouteTableImportState,
		},

		Schema: map[string]*schema.Schema{
			"lin_id": {
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

			"tag": tagsSchema(),

			"route_propagating_vpn_gateway": {
				Type:     schema.TypeList,
				ForceNew: true,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"route": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"destination_ip_range": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"destinaton_prefix_list_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vpn_gateway_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vm_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vm_account_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"nic_id": {
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
						"lin_peering_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"link": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"main": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"route_table_to_subnet_link_id": {
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

func resourceOutscaleOAPIRouteTableCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	createOpts := &fcu.CreateRouteTableInput{
		VpcId: aws.String(d.Get("lin_id").(string)),
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

	rt := resp.RouteTable
	d.SetId(*rt.RouteTableId)
	log.Printf("[INFO] Route Table ID: %s", d.Id())

	log.Printf(
		"[DEBUG] Waiting for route table (%s) to become available",
		d.Id())
	stateConf := &resource.StateChangeConf{
		Pending: []string{"pending"},
		Target:  []string{"ready"},
		Refresh: resourceOutscaleOAPIRouteTableStateRefreshFunc(conn, d.Id()),
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
		d.SetPartial("tag")
	}

	a := make([]interface{}, 0)

	d.Set("tag", a)
	d.Set("route", a)
	d.Set("link", a)

	return resourceOutscaleOAPIRouteTableRead(d, meta)
}

func resourceOutscaleOAPIRouteTableRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	rtRaw, _, err := resourceOutscaleOAPIRouteTableStateRefreshFunc(conn, d.Id())()
	if err != nil {
		return err
	}
	if rtRaw == nil {
		d.SetId("")
		return nil
	}

	rt := rtRaw.(*fcu.RouteTable)
	d.Set("lin_id", rt.VpcId)

	propagatingVGWs := make([]string, 0, len(rt.PropagatingVgws))
	for _, vgw := range rt.PropagatingVgws {
		propagatingVGWs = append(propagatingVGWs, *vgw.GatewayId)
	}
	d.Set("route_propagating_vpn_gateway", propagatingVGWs)

	d.Set("route", setOAPIRouteSet(rt.Routes))

	d.Set("link", setOAPIAssociactionSet(rt.Associations))

	d.Set("tag", tagsToMap(rt.Tags))

	return nil
}

func resourceOutscaleOAPIRouteTableDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	rtRaw, _, err := resourceOutscaleOAPIRouteTableStateRefreshFunc(conn, d.Id())()
	if err != nil {
		return err
	}
	if rtRaw == nil {
		return nil
	}
	rt := rtRaw.(*fcu.RouteTable)

	for _, a := range rt.Associations {
		log.Printf("[INFO] Disassociating association: %s", *a.RouteTableAssociationId)

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

	log.Printf("[INFO] Deleting Route Table: %s", d.Id())

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

	log.Printf(
		"[DEBUG] Waiting for route table (%s) to become destroyed",
		d.Id())

	stateConf := &resource.StateChangeConf{
		Pending: []string{"ready"},
		Target:  []string{},
		Refresh: resourceOutscaleOAPIRouteTableStateRefreshFunc(conn, d.Id()),
		Timeout: 5 * time.Minute,
	}
	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf(
			"Error waiting for route table (%s) to become destroyed: %s",
			d.Id(), err)
	}

	return nil
}

func resourceOutscaleOAPIRouteTableStateRefreshFunc(conn *fcu.Client, id string) resource.StateRefreshFunc {
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

func setOAPIRouteSet(rt []*fcu.Route) []map[string]interface{} {

	route := make([]map[string]interface{}, len(rt))

	if len(rt) > 0 {
		for k, r := range rt {
			if r.GatewayId != nil && *r.GatewayId == "local" {
				continue
			}

			if r.Origin != nil && *r.Origin == "EnableVgwRoutePropagation" {
				continue
			}

			if r.DestinationPrefixListId != nil {
				continue
			}

			m := make(map[string]interface{})

			if r.DestinationCidrBlock != nil {
				m["destination_ip_range"] = *r.DestinationCidrBlock
			}
			if r.DestinationPrefixListId != nil {
				m["destinaton_prefix_list_id"] = *r.DestinationPrefixListId
			}
			if r.GatewayId != nil {
				m["vpn_gateway_id"] = *r.GatewayId
			}
			if r.NatGatewayId != nil {
				m["nat_gateway_id"] = *r.NatGatewayId
			}
			if r.InstanceId != nil {
				m["vm_id"] = *r.InstanceId
			}
			if r.InstanceOwnerId != nil {
				m["vm_account_id"] = *r.InstanceOwnerId
			}
			if r.VpcPeeringConnectionId != nil {
				m["lin_peering_id"] = *r.VpcPeeringConnectionId
			}
			if r.NetworkInterfaceId != nil {
				m["nic_id"] = *r.NetworkInterfaceId
			}
			if r.Origin != nil {
				m["creation_method"] = *r.Origin
			}
			if r.State != nil {
				m["state"] = *r.State
			}

			route[k] = m
		}
	}

	return route
}

func setOAPIAssociactionSet(rt []*fcu.RouteTableAssociation) []map[string]interface{} {
	association := make([]map[string]interface{}, len(rt))

	if len(rt) > 0 {
		for k, r := range rt {
			m := make(map[string]interface{})

			if r.Main != nil {
				m["main"] = *r.Main
			}
			if r.RouteTableAssociationId != nil {
				m["route_table_to_subnet_link_id"] = *r.RouteTableAssociationId
			}
			if r.RouteTableId != nil {
				m["route_table_id"] = *r.RouteTableId
			}
			if r.SubnetId != nil {
				m["subnet_id"] = *r.SubnetId
			}

			association[k] = m
		}
	}

	return association
}
