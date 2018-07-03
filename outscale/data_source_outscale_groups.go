package outscale

import (
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-outscale/osc/eim"
)

func dataSourceOutscaleGroups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOutscaleGroupsRead,
		Schema: map[string]*schema.Schema{
			"groups": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"arn": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"group_id": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"group_name": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"user_name": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"path": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"path_prefix": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"request_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceOutscaleGroupsRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).EIM
	request := &eim.ListGroupsInput{
		PathPrefix: aws.String(d.Get("path_prefix").(string)),
	}

	var getResp *eim.ListGroupsOutput
	var err error

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		getResp, err = conn.API.ListGroups(request)
		if err != nil {
			if strings.Contains(err.Error(), "Throttling:") {
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
		return fmt.Errorf("Error reading IAM Group %s: %s", d.Id(), err)
	}

	if len(getResp.Groups) < 1 {
		return fmt.Errorf("No results found")
	}

	grps := make([]map[string]interface{}, len(getResp.Groups))
	for k, v := range getResp.Groups {
		grp := make(map[string]interface{})
		grp["arn"] = aws.StringValue(v.Arn)
		grp["group_id"] = aws.StringValue(v.GroupId)
		grp["group_name"] = aws.StringValue(v.GroupName)
		grp["user_name"] = aws.StringValue(v.UserName)
		grp["path"] = aws.StringValue(v.Path)
		grps[k] = grp
	}

	d.SetId(resource.UniqueId())

	return d.Set("groups", grps)
}
