package outscale

import (
	"fmt"
	"strings"
	"time"

	"github.com/terraform-providers/terraform-provider-outscale/osc/icu"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceOutscaleIamAccessKey() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOutscaleIamAccessKeyRead,

		Schema: map[string]*schema.Schema{
			"access_key": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"access_key_id": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
							Optional: true,
							ForceNew: true,
						},
						"owner_id": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"secret_access_key": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
							Optional: true,
							ForceNew: true,
						},
						"status": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"tag_set": tagsSchemaComputed(),
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

func dataSourceOutscaleIamAccessKeyRead(d *schema.ResourceData, meta interface{}) error {
	iamconn := meta.(*OutscaleClient).ICU

	request := &icu.ListAccessKeysInput{}

	var getResp *icu.ListAccessKeysOutput
	var err error
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		getResp, err = iamconn.API.ListAccessKeys(request)

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
		return fmt.Errorf("Error reading acces key: %s", err)
	}

	acc := make([]map[string]interface{}, len(getResp.AccessKeyMetadata))

	for k, v := range getResp.AccessKeyMetadata {
		ac := make(map[string]interface{})
		ac["access_key_id"] = *v.AccessKeyID
		ac["owner_id"] = *v.OwnerID
		ac["secret_access_key"] = *v.SecretAccessKey
		ac["status"] = *v.Status
		ac["tag_set"] = tagsToMapI(v.Tags)
		acc[k] = ac
	}

	d.SetId(resource.UniqueId())
	d.Set("access_key", acc)
	d.Set("request_id", getResp.ResponseMetadata.RequestID)

	return nil
}
