package outscale

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-outscale/osc/icu"
)

func resourceOutscaleAccessKey() *schema.Resource {
	return &schema.Resource{
		Create: resourceOutscaleAccessKeyCreate,
		Read:   resourceOutscaleAccessKeyRead,
		Delete: resourceOutscaleAccessKeyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Update: resourceOutscaleAccessKeyUpdate,

		Schema: map[string]*schema.Schema{
			"access_key_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"secret_access_key": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"tag": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:     schema.TypeString,
							Computed: true,
							Optional: true,
						},
						"value": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"owner_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"request_id": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
		},
	}
}

//Create AccessKey
func resourceOutscaleAccessKeyCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClientICU).ICU

	request := &icu.CreateAccessKeyInput{
		AccessKeyId:     aws.String(d.Get("access_key_id").(string)),
		SecretAccessKey: aws.String(d.Get("secret_access_key").(string)),
		//aqui faltan las tags pero quiero ver como se declaran dentro del struct
	}
	var err error
	createResp, err := conn.ICU_VM.CreateAccessKey(request)

	if err != nil {
		return fmt.Errorf(
			"Error creating access key for user %s: %s",
			*request.AccessKeyId,
			err,
		)
	}

	d.SetId(*createResp.AccessKey)

	if createResp.AccessKey == nil || createResp.ResponseMetadata == nil {
		return fmt.Errorf("[ERR] CreateAccessKey response did not contain a Secret Access Key as expected")
	}
	return resourceOutscaleAccessKeyRead(d, meta)
}
func resourceOutscaleAccessKeyRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClientICU).ICU

	request := &icu.DescribeAccessKeyInput{
		AccessKeyId:     aws.String(d.Get("access_key_id").(string)),
		SecretAccessKey: aws.String(d.Get("secret_access_key").(string)),
	}

	_, err := conn.ICU_VM.DescribeAccessKey(request)
	if err != nil {
		// the user does not exist, so the key can't exist.
		d.SetId("")

		return fmt.Errorf("Error reading IAM acces key: %s", err)
	}

	// Guess the key isn't around anymore.
	d.SetId("")
	return nil
}

func resourceOutscaleAccessKeyDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClientICU).ICU

	request := &icu.DeleteAccessKeyInput{
		AccessKeyId: aws.String(d.Id()),
	}

	if _, err := conn.ICU_VM.DeleteAccessKey(request); err != nil {
		fmt.Errorf("Error deleting access key %s: %s", d.Id(), err)
	}

	return nil
}

func resourceOutscaleAccessKeyUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClientICU).ICU

	request := &icu.UpdateAccessKeyInput{
		AccessKeyId: aws.String(d.Get("access_key_id").(string)),
		Status:      aws.String(d.Get("status").(string)),
	}
	var err error
	err = resource.Retry(120*time.Second, func() *resource.RetryError {
		_, err = conn.ICU_VM.UpdateAccessKey(request)

		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded:") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		log.Printf("[DEBUG] Error Updating AccessKey(%s)", err)
		return err
	}
	return nil
}
