package outscale

import (
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/terraform-providers/terraform-provider-outscale/osc/dl"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceOutscaleSites() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOutscaleSitesRead,

		Schema: map[string]*schema.Schema{
			"locations": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"location_code": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"location_name": &schema.Schema{
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

func dataSourceOutscaleSitesRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).DL

	request := &dl.DescribeLocationsInput{}

	var getResp *dl.DescribeLocationsOutput
	var err error
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		getResp, err = conn.API.DescribeLocations(request)

		if err != nil {
			if strings.Contains(fmt.Sprint(err), "Throttling") {
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
		return fmt.Errorf("Error reading sites: %s", err)
	}

	locations := make([]map[string]interface{}, len(getResp.Locations))

	for k, v := range getResp.Locations {
		location := make(map[string]interface{})
		location["location_code"] = aws.StringValue(v.LocationCode)
		location["location_name"] = aws.StringValue(v.LocationName)

		locations[k] = location
	}

	d.SetId(resource.UniqueId())
	d.Set("locations", locations)

	return d.Set("request_id", getResp.RequestID)
}
