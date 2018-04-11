package outscale

import (
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-outscale/osc/fcu"
)

func resourceOutscaleKeyPairImportation() *schema.Resource {
	return &schema.Resource{
		Create: resourceKeyPairImportationCreate,
		Read:   resourceKeyPairImportationRead,
		Delete: resourceKeyPairImportationDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Read:   schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: getKeyPairImportationSchema(),
	}
}

func resourceKeyPairImportationCreate(d *schema.ResourceData, meta interface{}) error {
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
			return fmt.Errorf("Error creating KeyPairImportation: %s", err)
		}
		d.SetId(*resp.KeyName)
		d.Set("public_key_material", *resp.KeyMaterial)
	}
	return nil
}

func resourceKeyPairImportationRead(d *schema.ResourceData, meta interface{}) error {
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
		if strings.Contains(fmt.Sprint(err), "InvalidKeyPair.NotFound") {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error retrieving KeyPair: %s", err)
	}

	d.Set("public_key_material", resp.KeyPairs[0].KeyName)
	d.Set("key_fingerprint", resp.KeyPairs[0].KeyFingerprint)
	d.Set("request_id", resp.RequesterId)

	return nil
}

func resourceKeyPairImportationDelete(d *schema.ResourceData, meta interface{}) error {
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

	return nil
}

func getKeyPairImportationSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		// Attributes
		"key_fingerprint": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"public_key_material": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"key_name": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"request_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}
