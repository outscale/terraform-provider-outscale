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

func dataSourceOutscaleOAPIGroupUser() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOutscaleOAPIGroupUserRead,
		Schema: map[string]*schema.Schema{
			"user_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
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
			"request_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceOutscaleOAPIGroupUserRead(d *schema.ResourceData, meta interface{}) error {
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

	d.Set("group_id", aws.StringValue(resp.Groups[0].GroupId))
	d.Set("group_name", aws.StringValue(resp.Groups[0].GroupName))
	d.Set("path", aws.StringValue(resp.Groups[0].Path))
	d.SetId(resource.UniqueId())

	return nil
}
