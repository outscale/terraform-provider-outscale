package outscale

import (
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/terraform-providers/terraform-provider-outscale/osc/eim"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceOutscaleOAPIUsers() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOutscaleOAPIUsersRead,

		Schema: map[string]*schema.Schema{
			"users": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
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
			"request_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"path": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func dataSourceOutscaleOAPIUsersRead(d *schema.ResourceData, meta interface{}) error {
	iamconn := meta.(*OutscaleClient).EIM

	request := &eim.ListUsersInput{
		PathPrefix: aws.String(d.Get("path").(string)),
	}

	var err error
	var getResp *eim.ListUsersOutput
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		getResp, err = iamconn.API.ListUsers(request)
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
		return fmt.Errorf("Error reading IAM User %s: %s", d.Id(), err)
	}

	if len(getResp.Users) < 1 {
		return fmt.Errorf("No results")
	}

	users := make([]map[string]interface{}, len(getResp.Users))
	for k, v := range getResp.Users {
		user := make(map[string]interface{})
		user["path"] = aws.StringValue(v.Path)
		user["user_id"] = aws.StringValue(v.UserId)
		user["user_name"] = aws.StringValue(v.UserName)
		users[k] = user
	}

	d.SetId(resource.UniqueId())

	return d.Set("users", users)
}
