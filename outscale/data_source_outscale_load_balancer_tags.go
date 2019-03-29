package outscale

import (
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/terraform-providers/terraform-provider-outscale/osc/lbu"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceOutscaleLBUTags() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOutscaleLBUTagsRead,

		Schema: getDSLBUTagsSchema(),
	}
}

func dataSourceOutscaleLBUTagsRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).LBU

	lbus := d.Get("load_balancer_names")

	params := &lbu.DescribeTagsInput{
		LoadBalancerNames: expandStringList(lbus.([]interface{})),
	}

	var resp *lbu.DescribeTagsResult
	var rs *lbu.DescribeTagsOutput
	var err error

	err = resource.Retry(60*time.Second, func() *resource.RetryError {
		rs, err = conn.API.DescribeTags(params)
		if rs != nil {
			resp = rs.DescribeTagsResult
		}
		return resource.RetryableError(err)
	})

	if err != nil {
		return err
	}

	td := make([]map[string]interface{}, len(resp.TagDescriptions))

	for k, v := range resp.TagDescriptions {
		t := make(map[string]interface{})
		t["load_balancer_name"] = aws.StringValue(v.LoadBalancerName)

		ta := make([]map[string]interface{}, len(v.Tags))
		for k1, v1 := range v.Tags {
			t := make(map[string]interface{})
			t["key"] = aws.StringValue(v1.Key)
			t["value"] = aws.StringValue(v1.Key)
			ta[k1] = t
		}

		t["tags"] = ta

		td[k] = t
	}

	d.SetId(resource.UniqueId())
	d.Set("request_id", rs.ResponseMetadata.RequestID)

	return d.Set("tag_descriptions", td)
}

func getDSLBUTagsSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"load_balancer_names": {
			Type:     schema.TypeList,
			Optional: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
		"tag_descriptions": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"load_balancer_name": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"tags": {
						Type:     schema.TypeList,
						Computed: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"key": {
									Type:     schema.TypeString,
									Computed: true,
								},
								"value": {
									Type:     schema.TypeString,
									Computed: true,
								},
							},
						},
					},
				},
			},
		},
		"request_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}
