package outscale

import (
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-outscale/osc/fcu"
)

func resourceOutscaleOAPIKeyPair() *schema.Resource {
	return &schema.Resource{
		Create: resourceOAPIKeyPairCreate,
		Read:   resourceOAPIKeyPairRead,
		Delete: resourceOAPIKeyPairDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Read:   schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: getOAPIKeyPairSchema(),
	}
}

func resourceOAPIKeyPairCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	var keyName string
	if v, ok := d.GetOk("key_name"); ok {
		keyName = v.(string)
	} else {
		keyName = resource.UniqueId()
		d.Set("key_name", keyName)
	}
	if publicKey, ok := d.GetOk("key_material"); ok {
		req := &fcu.ImportOAPIKeyPairInput{
			KeyName:           aws.String(keyName),
			PublicKeyMaterial: []byte(publicKey.(string)),
		}

		var resp *fcu.ImportOAPIKeyPairOutput
		err := resource.Retry(120*time.Second, func() *resource.RetryError {
			var err error
			resp, err = conn.VM.ImportOAPIKeyPair(req)

			if err != nil {
				if strings.Contains(err.Error(), "RequestLimitExceeded:") {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return resource.RetryableError(err)
		})

		if err != nil {
			return fmt.Errorf("Error import OAPIKeyPair: %s", err)
		}
		d.SetId(*resp.KeyName)

	} else {
		req := &fcu.CreateOAPIKeyPairInput{
			KeyName: aws.String(keyName),
		}

		var resp *fcu.CreateOAPIKeyPairOutput
		err := resource.Retry(120*time.Second, func() *resource.RetryError {
			var err error
			resp, err = conn.VM.CreateOAPIKeyPair(req)

			if err != nil {
				if strings.Contains(err.Error(), "RequestLimitExceeded:") {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return resource.RetryableError(err)
		})
		if err != nil {
			return fmt.Errorf("Error creating OAPIKeyPair: %s", err)
		}
		d.SetId(*resp.KeyName)
		d.Set("key_material", *resp.KeyMaterial)
	}
	return nil
}

func resourceOAPIKeyPairRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU
	req := &fcu.DescribeOAPIKeyPairsInput{
		KeyNames: []*string{aws.String(d.Id())},
	}

	var resp *fcu.DescribeOAPIKeyPairsOutput
	err := resource.Retry(120*time.Second, func() *resource.RetryError {
		var err error
		resp, err = conn.VM.DescribeOAPIKeyPairs(req)

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
		if ok && awsErr.Code() == "InvalidOAPIKeyPair.NotFound" {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error retrieving OAPIKeyPair: %s", err)
	}

	for _, keyPair := range resp.OAPIKeyPairs {
		if *keyPair.KeyName == d.Id() {
			d.Set("key_name", keyPair.KeyName)
			d.Set("fingerprint", keyPair.KeyFingerprint)
			return nil
		}
	}

	return fmt.Errorf("Unable to find key pair within: %#v", resp.OAPIKeyPairs)
}

func resourceOAPIKeyPairDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		var err error
		_, err = conn.VM.DeleteOAPIKeyPairs(&fcu.DeleteOAPIKeyPairInput{
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

func getOAPIKeyPairSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		// Attributes
		"key_fingerprint": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"key_material": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"key_name": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
	}
}
