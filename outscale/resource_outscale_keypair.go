package outscale

import (
	"fmt"
	"strings"
	"time"

	"github.com/terraform-providers/terraform-provider-outscale/utils"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-outscale/osc/fcu"
)

func resourceOutscaleKeyPair() *schema.Resource {
	return &schema.Resource{
		Create: resourceKeyPairCreate,
		Read:   resourceKeyPairRead,
		Delete: resourceKeyPairDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Read:   schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: getKeyPairSchema(),
	}
}

func resourceKeyPairCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	var keyName string
	if v, ok := d.GetOk("key_name"); ok {
		keyName = v.(string)
	} else {
		keyName = resource.UniqueId()
		d.Set("key_name", keyName)
	}

	req := &fcu.CreateKeyPairInput{
		KeyName: aws.String(keyName),
	}

	var resp *fcu.CreateKeyPairOutput
	err := resource.Retry(120*time.Second, func() *resource.RetryError {
		var err error
		resp, err = conn.VM.CreateKeyPair(req)

		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded:") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return resource.RetryableError(err)
	})
	if err != nil {
		return fmt.Errorf("Error creating KeyPair: %s", err)
	}
	d.SetId(*resp.KeyName)
	d.Set("key_material", *resp.KeyMaterial)
	return nil
}

func resourceKeyPairRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU
	req := &fcu.DescribeKeyPairsInput{
		KeyNames: []*string{aws.String(d.Id())},
	}

	var resp *fcu.DescribeKeyPairsOutput
	err := resource.Retry(120*time.Second, func() *resource.RetryError {
		var err error
		resp, err = conn.VM.DescribeKeyPairs(req)

		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded:") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return resource.RetryableError(err)
	})

	if err != nil {
		return err
	}

	if err != nil {
		awsErr, ok := err.(awserr.Error)
		if ok && awsErr.Code() == "InvalidKeyPair.NotFound" {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error retrieving KeyPair: %s", err)
	}

	fmt.Println("\n\n[DEBUG] RESP")
	utils.PrintToJSON(resp, "KEY_PAIR")

	d.Set("key_name", resp.KeyPairs[0].KeyName)
	d.Set("key_fingerprint", resp.KeyPairs[0].KeyFingerprint)
	d.Set("request_id", resp.RequestId)

	return nil
}

func resourceKeyPairDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		var err error
		_, err = conn.VM.DeleteKeyPairs(&fcu.DeleteKeyPairInput{
			KeyName: aws.String(d.Id()),
		})
		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded:") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return resource.NonRetryableError(err)
	})

	if err != nil {
		return err
	}

	return err
}

func getKeyPairSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		// Attributes
		"key_fingerprint": {
			Type:     schema.TypeString,
			Optional: true,
			ForceNew: true,
			Computed: true,
		},
		"key_material": {
			Type:     schema.TypeString,
			Optional: true,
			ForceNew: true,
			Computed: true,
		},
		"key_name": {
			Type:     schema.TypeString,
			Optional: true,
			ForceNew: true,
			Computed: true,
		},
		"request_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}
