package outscale

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-outscale/osc/dl"
)

func dataSourceOutscaleOAPIDirectLink() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOutscaleOAPIDirectLinkRead,

		Schema: map[string]*schema.Schema{
			"directlink_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"directlinks": {
				Type:     schema.TypeString,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"bandwidth": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"directlink_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"directlink_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"state": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"site": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"account_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"region_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
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

func dataSourceOutscaleOAPIDirectLinkRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).DL

	var resp *dl.Connections
	var err error
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, err = conn.API.DescribeConnections(&dl.DescribeConnectionsInput{
			ConnectionID: aws.String(d.Get("directlink_id").(string)),
		})

		if err != nil {
			if strings.Contains(fmt.Sprint(err), "RequestLimitExceeded:") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		if isNoSuchDxConnectionErr(err) {
			log.Printf("[WARN] Direct Connect connection (%s) not found, removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return err
	}

	if len(resp.Connections) < 1 {
		log.Printf("[WARN] Direct Connect connection (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	conections := make([]map[string]interface{}, len(resp.Connections))
	for k, connection := range resp.Connections {
		con := make(map[string]interface{})
		con["bandwidth"] = aws.StringValue(connection.Bandwidth)
		con["directlink_id"] = aws.StringValue(connection.ConnectionID)
		con["directlink_name"] = aws.StringValue(connection.ConnectionName)
		con["state"] = aws.StringValue(connection.ConnectionState)
		con["site"] = aws.StringValue(connection.Location)
		con["account_id"] = aws.StringValue(connection.OwnerAccount)
		con["region_name"] = aws.StringValue(connection.Region)
		conections[k] = con
	}

	d.Set("directlinks", conections)

	d.SetId(resource.UniqueId())

	return d.Set("request_id", resp.RequestID)
}
