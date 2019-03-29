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

func dataSourceOutscaleLinPeeringsConnection() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOutscaleLinPeeringsConnectionRead,

		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			"vpc_peering_connection_id": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"vpc_peering_connection_set": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
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
					},
				},
			},
			"request_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceOutscaleLinPeeringsConnectionRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	log.Printf("[DEBUG] Reading VPC Peering Connections.")

	id, ok := d.GetOk("vpc_peering_connection_id")
	v, vok := d.GetOk("filter")

	if ok == false && vok == false {
		return fmt.Errorf("One of filters, or instance_id must be assigned")
	}

	req := &fcu.DescribeVpcPeeringConnectionsInput{}

	if ok {
		ids := make([]*string, len(id.([]interface{})))

		for k, v := range id.([]interface{}) {
			ids[k] = aws.String(v.(string))
		}

		req.VpcPeeringConnectionIds = ids
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

	lps := make([]map[string]interface{}, len(resp.VpcPeeringConnections))
	for k, v := range resp.VpcPeeringConnections {
		lp := make(map[string]interface{})
		accepter := make(map[string]interface{})
		requester := make(map[string]interface{})
		stat := make(map[string]interface{})

		if v.AccepterVpcInfo != nil {
			accepter["cidr_block"] = aws.StringValue(v.AccepterVpcInfo.CidrBlock)
			accepter["owner_id"] = aws.StringValue(v.AccepterVpcInfo.OwnerId)
			accepter["vpc_id"] = aws.StringValue(v.AccepterVpcInfo.VpcId)
		}
		if v.RequesterVpcInfo != nil {
			requester["cidr_block"] = aws.StringValue(v.AccepterVpcInfo.CidrBlock)
			requester["owner_id"] = aws.StringValue(v.AccepterVpcInfo.OwnerId)
			requester["vpc_id"] = aws.StringValue(v.AccepterVpcInfo.VpcId)
		}
		if v.Status != nil {
			stat["code"] = aws.StringValue(v.Status.Code)
			stat["message"] = aws.StringValue(v.Status.Message)
		}

		lp["accepter_vpc_info"] = accepter
		lp["requester_vpc_info"] = requester
		lp["status"] = stat
		lp["vpc_peering_connection_id"] = *v.VpcPeeringConnectionId
		lp["tag_set"] = tagsToMap(v.Tags)

		lps[k] = lp
	}

	d.SetId(resource.UniqueId())
	if err := d.Set("vpc_peering_connection_set", lps); err != nil {
		return err
	}
	d.Set("request_id", resp.RequestId)

	return nil
}
