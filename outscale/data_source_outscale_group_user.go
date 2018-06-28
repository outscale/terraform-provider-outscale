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

func dataSourceOutscaleGroupUser() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOutscaleGroupUserRead,
		Schema: map[string]*schema.Schema{
			"user_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
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
			"request_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceOutscaleGroupUserRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).EIM
	user := d.Get("user_name").(string)

	var err error
	var resp *eim.ListGroupsForUserOutput
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, err = conn.API.ListGroupsForUser(&eim.ListGroupsForUserInput{
			UserName: aws.String(user),
		})
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
		return err
	}

	if len(resp.Groups) < 1 {
		return fmt.Errorf("no results")
	}

	uls := make([]map[string]interface{}, len(resp.Groups))

	for k, v := range resp.Groups {
		ul := make(map[string]interface{})
		ul["arn"] = aws.StringValue(v.Arn)
		ul["group_id"] = aws.StringValue(v.GroupId)
		ul["group_name"] = aws.StringValue(v.GroupName)
		ul["user_name"] = aws.StringValue(v.UserName)
		ul["path"] = aws.StringValue(v.Path)
		uls[k] = ul
	}

	d.SetId(resource.UniqueId())

	return d.Set("groups", uls)
}
