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

func dataSourceOutscaleDirectLink() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOutscaleDirectLinkRead,

		Schema: map[string]*schema.Schema{
			"connection_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"connections": {
				Type:     schema.TypeString,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"bandwidth": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"connection_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"connection_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"connection_state": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"location": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"owner_account": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"region": {
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

func dataSourceOutscaleDirectLinkRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).DL

	var resp *dl.Connections
	var err error
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, err = conn.API.DescribeConnections(&dl.DescribeConnectionsInput{
			ConnectionID: aws.String(d.Get("connection_id").(string)),
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
		con["connection_id"] = aws.StringValue(connection.ConnectionID)
		con["connection_name"] = aws.StringValue(connection.ConnectionName)
		con["connection_state"] = aws.StringValue(connection.ConnectionState)
		con["location"] = aws.StringValue(connection.Location)
		con["owner_account"] = aws.StringValue(connection.OwnerAccount)
		con["region"] = aws.StringValue(connection.Region)
		conections[k] = con
	}

	d.Set("connections", conections)

	d.SetId(resource.UniqueId())

	return d.Set("request_id", resp.RequestID)
}
