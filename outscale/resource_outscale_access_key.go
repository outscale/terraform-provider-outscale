package outscale

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-outscale/osc/fcu"
)

func resourceOutscaleAccessKey() *schema.Resource {
	return &schema.Resource{
		Create: resourceOutscaleAccessKeyCreate,
		Read: 	resourceOutscaleAccessKeyRead,
		Delete: resourceOutscaleAccessKeyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Update: resourceOutscaleAccessKeyUpdate,
		Schema: getAccessKeySchema(),
	}
}

//Create AccessKey
func resourceOutscaleAccessKeyCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	createOpts := &fcu.CreateAccessKeyInput{
		AccessKeyId: 			aws.String(d.Get("access_key_id").(string)),
		SecretAccessKey: 	aws.String(d.Get("secret_access_key").(string)),
	}

	var res *fcu.CreateAccessKeyOutput
	var err error
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		res, err = conn.VM.CreateAccessKey(createOpts)

		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("Error creating AccessKey: %s", err)
	}

	// Get the ID and store it

	d.SetId(*res.AccessKey)
	log.Printf("[DEBUG] Waiting for AccessKey creation ")

	return resourceOutscaleAccessKeyRead(d, meta)
}


func resourceOutscaleAccessKeyRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	var resp *fcu.DescribeAccessKeyOutput
	var respi *fcu.DescribeAccessKeyInput
	var err error
	err = resource.Retry(120*time.Second, func() *resource.RetryError {
		resp, err = conn.VM.DescribeAccessKey(&fcu.DescribeAccessKeyInput{
			AccessKeyId: 			aws.String(d.Get("access_key_id").(string)),
		})

		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded:") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		if strings.Contains(err.Error(), "InvalidAccessKey.NotFound") {
			// Update state to indicate the AccessKey no longer exists.
			d.SetId("")
			return nil
		}
		return err
	}
	if resp == nil {
		return nil
	}

	accesskey := respi

	d.Set("access_key_id", accesskey.AccessKeyId)
	d.Set("secret_access_key", accesskey.SecretAccessKey)


	if err := d.Set("tag_set", tagsToMap(accesskey.Tags)); err != nil {
		return err
	}

	return nil
}

func resourceOutscaleAccessKeyDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	id := d.Id()
	log.Printf("[DEBUG] Deleting AccessKey (%s)", id)

	req := &fcu.DeleteAccessKeyInput{
		AccessKeyId: &id,
	}

	var err error
	err = resource.Retry(120*time.Second, func() *resource.RetryError {
		_, err = conn.VM.DeleteAccessKey(req)

		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded:") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		log.Printf("[DEBUG] Error deleting AccessKey(%s)", err)
		return err
	}

	return nil
}

func resourceOutscaleAccessKeyUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU
	id := d.Id()
	log.Printf("Updating AccessKey")
	req :=&fcu.UpdateAccessKeyInput{
		AccessKeyId: &id,
	}
	var err error
	err = resource.Retry(120*time.Second, func() *resource.RetryError {
		_, err = conn.VM.UpdateAccessKey(req)

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

func getAccessKeySchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"access_key_id": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"secret_access_key": &schema.Schema{
			Type: schema.TypeString,
			Optional: true,
		},
		"tag_set": dataSourceTagsSchema(),
		"owner_id": &schema.Schema{
			Type: schema.TypeString,
			Optional: true,
		},
		"request_id": &schema.Schema{
			Type: schema.TypeString,
			Optional: true,
		},
	}
}
