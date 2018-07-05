package outscale

import (
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/terraform-providers/terraform-provider-outscale/osc/icu"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceOutscaleAccountConsumption() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOutscaleAccountConsumptionRead,

		Schema: map[string]*schema.Schema{
			"from_date": {
				Type:     schema.TypeString,
				Required: true,
			},
			"to_date": {
				Type:     schema.TypeString,
				Required: true,
			},
			"entries": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"category": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"operation": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"service": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"title": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"value": &schema.Schema{
							Type:     schema.TypeInt,
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

func dataSourceOutscaleAccountConsumptionRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).ICU

	request := &icu.ReadConsumptionAccountInput{
		FromDate: aws.String(d.Get("from_date").(string)),
		ToDate:   aws.String(d.Get("to_date").(string)),
	}

	var getResp *icu.ReadConsumptionAccountOutput
	var err error
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		getResp, err = conn.API.ReadConsumptionAccount(request)

		if err != nil {
			if strings.Contains(fmt.Sprint(err), "RequestLimitExceeded:") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		if strings.Contains(fmt.Sprint(err), "NoSuchEntity") {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error reading account consuption: %s", err)
	}

	entries := make([]map[string]interface{}, len(getResp.Entries))

	for k, v := range getResp.Entries {
		entry := make(map[string]interface{})
		entry["category"] = aws.StringValue(v.Category)
		entry["operation"] = aws.StringValue(v.Operation)
		entry["service"] = aws.StringValue(v.Service)
		entry["title"] = aws.StringValue(v.Title)
		entry["type"] = aws.StringValue(v.Type)
		entry["value"] = aws.Float64Value(v.Value)

		entries[k] = entry
	}

	d.SetId(resource.UniqueId())
	d.Set("entries", entries)

	return d.Set("request_id", getResp.ResponseMetadata.RequestID)
}
