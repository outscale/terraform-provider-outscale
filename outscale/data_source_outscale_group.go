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

func dataSourceOutscaleGroup() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOutscaleGroupRead,
		Schema: map[string]*schema.Schema{
			"group": &schema.Schema{
				Type:     schema.TypeMap,
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
						"path": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"users": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"arn": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"user_id": &schema.Schema{
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
			"group_name": &schema.Schema{
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

func dataSourceOutscaleGroupRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).EIM
	request := &eim.GetGroupInput{
		GroupName: aws.String(d.Get("group_name").(string)),
	}

	var getResp *eim.GetGroupOutput
	var err error

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		getResp, err = conn.API.GetGroup(request)
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

	grp := make(map[string]interface{})
	grp["arn"] = aws.StringValue(getResp.Group.Arn)
	grp["group_id"] = aws.StringValue(getResp.Group.GroupId)
	grp["group_name"] = aws.StringValue(getResp.Group.GroupName)
	grp["path"] = aws.StringValue(getResp.Group.Path)

	usr := make([]map[string]interface{}, len(getResp.Users))
	for k, v := range getResp.Users {
		us := make(map[string]interface{})
		us["arn"] = aws.StringValue(v.Arn)
		us["user_id"] = aws.StringValue(v.UserId)
		us["user_name"] = aws.StringValue(v.UserName)
		us["path"] = aws.StringValue(v.Path)
		usr[k] = us
	}

	if err := d.Set("group", grp); err != nil {
		return err
	}

	d.SetId(resource.UniqueId())

	return d.Set("users", usr)
}
