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

func resourceOutscaleOAPIUserAPIKey() *schema.Resource {
	return &schema.Resource{
		Read:   resourceOutscaleOAPIUserAPIKeyRead,
		Create: resourceOutscaleOAPIUserAPIKeyCreate,
		Update: resourceOutscaleOAPIUserAPIKeyUpdate,
		Delete: resourceOutscaleOAPIUserAPIKeyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"user_name": {
				Type:     schema.TypeString,
				Computed: true,
				ForceNew: true,
				Optional: true,
			},
			"state": { //Only works on update
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"api_key_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"secret_key": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"request_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			//missing account_id //string
			//tags //key-value (argument)
		},
	}
}

func resourceOutscaleOAPIUserAPIKeyCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).EIM
	name := aws.String(d.Get("user_name").(string))

	request := &eim.CreateAccessKeyInput{
		UserName: name,
	}

	var err error
	var resp *eim.CreateAccessKeyOutput
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, err = conn.API.CreateAccessKey(request)

		if err != nil {
			if strings.Contains(err.Error(), "Throttling:") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("Failure creating access key for EIM: %s", err)
	}

	if resp.CreateAccessKeyResult == nil {
		return fmt.Errorf("Cannot unmarshal result of AccessKeys")
	}

	if resp.CreateAccessKeyResult.AccessKey == nil || resp.CreateAccessKeyResult.AccessKey.SecretAccessKey == nil {
		return fmt.Errorf("[ERR] CreateAccessKey response did not contain a Secret Access Key as expected")
	}

	if err := d.Set("secret_key", resp.CreateAccessKeyResult.AccessKey.SecretAccessKey); err != nil {
		return err
	}

	d.SetId(aws.StringValue(resp.CreateAccessKeyResult.AccessKey.AccessKeyID))
	return resourceOutscaleOAPIUserAPIKeyRead(d, meta)
}

func resourceOutscaleOAPIUserAPIKeyRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).EIM

	request := &eim.ListAccessKeysInput{
		UserName: aws.String(d.Get("user_name").(string)),
	}

	var err error
	var resp *eim.ListAccessKeysOutput
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, err = conn.API.ListAccessKeys(request)

		if err != nil {
			if strings.Contains(err.Error(), "Throttling:") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("Failure get access key for EIM: %s", err)
	}

	if resp.ListAccessKeysResult == nil {
		return fmt.Errorf("Cannot unmarshal result of AccessKeys")
	}

	for _, key := range resp.ListAccessKeysResult.AccessKeyMetadata {
		if key.AccessKeyID != nil && *key.AccessKeyID == d.Id() {
			d.Set("request_id", aws.StringValue(resp.ResponseMetadata.RequestID))
			return resourceOutsclaeOAPIAccessKeyReadResult(d, key)
		}
	}

	d.SetId("")
	return nil
}

func resourceOutscaleOAPIUserAPIKeyUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).EIM

	if !d.HasChange("state") {
		return nil
	}

	state, ok := d.GetOk("state")

	if !ok {
		return fmt.Errorf("You must set `Active` or `Inactive` value")
	}

	if state.(string) != "Active" && state.(string) != "Inactive" {
		return fmt.Errorf("You must set `Active` or `Inactive` value")
	}

	request := &eim.UpdateAccessKeyInput{
		AccessKeyID: aws.String(d.Id()),
		UserName:    aws.String(d.Get("user_name").(string)),
		Status:      aws.String(state.(string)),
	}

	var err error
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, err = conn.API.UpdateAccessKey(request)

		if err != nil {
			if strings.Contains(err.Error(), "Throttling:") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("Failure get access key for EIM: %s", err)
	}

	return resourceOutscaleOAPIUserAPIKeyRead(d, meta)
}
func resourceOutscaleOAPIUserAPIKeyDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).EIM

	request := &eim.DeleteAccessKeyInput{
		AccessKeyID: aws.String(d.Id()),
		UserName:    aws.String(d.Get("user_name").(string)),
	}

	var err error
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, err = conn.API.DeleteAccessKey(request)

		if err != nil {
			if strings.Contains(err.Error(), "Throttling:") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("Failure delete access key for EIM: %s", err)
	}

	d.SetId("")
	return nil
}

func resourceOutsclaeOAPIAccessKeyReadResult(d *schema.ResourceData, key *eim.AccessKeyMetadata) error {
	d.SetId(*key.AccessKeyID)

	if err := d.Set("api_key_id", key.AccessKeyID); err != nil {
		return err
	}

	if err := d.Set("state", key.Status); err != nil {
		return err
	}
	return nil
}
