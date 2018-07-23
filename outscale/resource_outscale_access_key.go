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

func resourceOutscaleIamAccessKey() *schema.Resource {
	return &schema.Resource{
		Create: resourceOutscaleIamAccessKeyCreate,
		Read:   resourceOutscaleIamAccessKeyRead,
		Delete: resourceOutscaleIamAccessKeyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

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
			"tag":     tagsSchema(),
			"tag_set": tagsSchemaComputed(),
			"request_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceOutscaleIamAccessKeyCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient)
	iamconn := conn.ICU

	request := &icu.CreateAccessKeyInput{}

	if v, ok := d.GetOk("access_key_id"); ok {
		request.AccessKeyID = aws.String(v.(string))
	}
	if v, ok := d.GetOk("secret_access_key"); ok {
		request.SecretAccessKey = aws.String(v.(string))
	}
	if _, ok := d.GetOk("tag"); ok {
		oraw, nraw := d.GetChange("tag")
		o := oraw.(map[string]interface{})
		n := nraw.(map[string]interface{})
		create, _ := diffTagsCommon(tagsFromMapCommon(o), tagsFromMapCommon(n))

		request.Tags = create
	}

	var createResp *icu.CreateAccessKeyOutput
	var err error
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		createResp, err = iamconn.API.CreateAccessKey(request)

		if err != nil {
			if strings.Contains(fmt.Sprint(err), "RequestLimitExceeded:") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("ERROR Creating key, %s", err)
	}

	d.SetId(*createResp.AccessKey.AccessKeyId)

	d.Set("tag_set", make([]map[string]interface{}, 0))

	return resourceOutscaleIamAccessKeyRead(d, meta)
}

func resourceOutscaleIamAccessKeyRead(d *schema.ResourceData, meta interface{}) error {
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

	for _, key := range getResp.AccessKeyMetadata {
		if key.AccessKeyID != nil && *key.AccessKeyID == d.Id() {
			d.Set("access_key_id", key.AccessKeyID)
			d.Set("secret_access_key", key.SecretAccessKey)
			d.Set("owner_id", key.OwnerID)
			d.Set("status", key.Status)
			d.Set("tag_set", tagsToMapI(key.Tags))
			return d.Set("request_id", getResp.ResponseMetadata.RequestID)
		}
	}
	d.SetId("")
	return fmt.Errorf("AccessKey not found")
}

func resourceOutscaleIamAccessKeyDelete(d *schema.ResourceData, meta interface{}) error {
	iamconn := meta.(*OutscaleClient).ICU

	request := &icu.DeleteAccessKeyInput{
		AccessKeyId: aws.String(d.Id()),
	}

	var err error
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, err = iamconn.API.DeleteAccessKey(request)

		if err != nil {
			if strings.Contains(fmt.Sprint(err), "RequestLimitExceeded:") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("Error deleting access key %s: %s", d.Id(), err)
	}
	return nil
}
