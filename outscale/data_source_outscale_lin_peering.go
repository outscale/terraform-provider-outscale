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

func dataSourceOutscaleLinPeeringConnection() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOutscaleLinPeeringConnectionRead,

		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			"vpc_peering_connection_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"code": {
							Type:     schema.TypeString,
							Computed: true,
							Optional: true,
						},
						"message": {
							Type:     schema.TypeString,
							Computed: true,
							Optional: true,
						},
					},
				},
			},
			"accepter_vpc_info":  vpcPeeringConnectionOptionsSchema(),
			"requester_vpc_info": vpcPeeringConnectionOptionsSchema(),
			"tag_set":            tagsSchemaComputed(),
			"request_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceOutscaleLinPeeringConnectionRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	log.Printf("[DEBUG] Reading VPC Peering Connections.")

	id, ok := d.GetOk("vpc_peering_connection_id")
	v, vok := d.GetOk("filter")

	if ok == false && vok == false {
		return fmt.Errorf("One of filters, or instance_id must be assigned")
	}

	req := &fcu.DescribeVpcPeeringConnectionsInput{}

	if ok {
		req.VpcPeeringConnectionIds = aws.StringSlice([]string{id.(string)})
	}
	if vok {
		req.Filters = buildOutscaleDataSourceFilters(v.(*schema.Set))
	}

	var resp *fcu.DescribeVpcPeeringConnectionsOutput
	var err error
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, err = conn.VM.DescribeVpcPeeringConnections(req)

		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded:") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		if strings.Contains(fmt.Sprint(err), "InvalidVpcPeeringConnectionID.NotFound") {
			resp = nil
		} else {
			log.Printf("Error reading VPC Peering Connection details: %s", err)
			return err
		}
	}

	if err != nil {
		return err
	}
	if resp == nil || len(resp.VpcPeeringConnections) == 0 {
		return fmt.Errorf("no matching VPC peering connection found")
	}
	if len(resp.VpcPeeringConnections) > 1 {
		return fmt.Errorf("multiple VPC peering connections matched; use additional constraints to reduce matches to a single VPC peering connection")
	}

	pc := resp.VpcPeeringConnections[0]

	d.SetId(aws.StringValue(pc.VpcPeeringConnectionId))

	accepter := make(map[string]interface{})
	requester := make(map[string]interface{})
	stat := make(map[string]interface{})

	if pc.AccepterVpcInfo != nil {
		accepter["cidr_block"] = aws.StringValue(pc.AccepterVpcInfo.CidrBlock)
		accepter["owner_id"] = aws.StringValue(pc.AccepterVpcInfo.OwnerId)
		accepter["vpc_id"] = aws.StringValue(pc.AccepterVpcInfo.VpcId)
	}
	if pc.RequesterVpcInfo != nil {
		requester["cidr_block"] = aws.StringValue(pc.AccepterVpcInfo.CidrBlock)
		requester["owner_id"] = aws.StringValue(pc.AccepterVpcInfo.OwnerId)
		requester["vpc_id"] = aws.StringValue(pc.AccepterVpcInfo.VpcId)
	}
	if pc.Status != nil {
		stat["code"] = aws.StringValue(pc.Status.Code)
		stat["message"] = aws.StringValue(pc.Status.Message)
	}

	if err := d.Set("accepter_vpc_info", accepter); err != nil {
		return err
	}
	if err := d.Set("requester_vpc_info", requester); err != nil {
		return err
	}
	if err := d.Set("status", stat); err != nil {
		return err
	}
	if err := d.Set("vpc_peering_connection_id", pc.VpcPeeringConnectionId); err != nil {
		return err
	}
	if err := d.Set("tag_set", tagsToMap(pc.Tags)); err != nil {
		return err
	}

	d.Set("request_id", resp.RequestId)

	return nil
}
