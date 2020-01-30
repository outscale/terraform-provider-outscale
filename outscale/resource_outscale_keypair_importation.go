package outscale

import (
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceOutscaleOAPIKeyPairImportation() *schema.Resource {
	return &schema.Resource{
		Create: resourceOAPIKeyPairImportationCreate,
		Read:   resourceOAPIKeyPairImportationRead,
		Delete: resourceOAPIKeyPairImportationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Read:   schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: getOAPIKeyPairImportationSchema(),
	}
}

func resourceOAPIKeyPairImportationCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	var keyName string
	if v, ok := d.GetOk("public_key_material"); ok {
		keyName = v.(string)
	} else {
		keyName = resource.UniqueId()
		d.Set("public_key_material", keyName)
	}
	if publicKey, ok := d.GetOk("public_key_material"); ok {
		req := &fcu.ImportKeyPairInput{
			KeyName:           aws.String(keyName),
			PublicKeyMaterial: []byte(publicKey.(string)),
		}

		var resp *fcu.ImportKeyPairOutput
		err := resource.Retry(120*time.Second, func() *resource.RetryError {
			var err error
			resp, err = conn.VM.ImportKeyPair(req)

			if err != nil {
				if strings.Contains(err.Error(), "RequestLimitExceeded:") {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return resource.RetryableError(err)
		})

		if err != nil {
			return fmt.Errorf("Error import KeyPair: %s", err)
		}
		d.SetId(*resp.KeyName)

	} else {
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
			return fmt.Errorf("Error creating OAPIKeyPairImportation: %s", err)
		}
		d.SetId(*resp.KeyName)
		d.Set("public_key_material", *resp.KeyMaterial)
	}
	return nil
}

func resourceOAPIKeyPairImportationRead(d *schema.ResourceData, meta interface{}) error {
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

	for _, keyPair := range resp.KeyPairs {
		if *keyPair.KeyName == d.Id() {
			d.Set("public_key_material", keyPair.KeyName)
			d.Set("fingerprint", keyPair.KeyFingerprint)
			return nil
		}
	}

	return fmt.Errorf("Unable to find key pair within: %#v", resp.KeyPairs)
}

func resourceOAPIKeyPairImportationDelete(d *schema.ResourceData, meta interface{}) error {
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

func getOAPIKeyPairImportationSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		// Attributes
		"key_fingerprint": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"public_key_material": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"key_name": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"request_id": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
	}
}
