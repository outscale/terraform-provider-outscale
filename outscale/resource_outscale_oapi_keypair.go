package outscale

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-outscale/osc/oapi"
	"github.com/terraform-providers/terraform-provider-outscale/utils"
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
	conn := meta.(*OutscaleClient).OAPI

	var keyName string
	if v, ok := d.GetOk("keypair_name"); ok {
		keyName = v.(string)
	} else {
		keyName = resource.UniqueId()
		d.Set("keypair_name", keyName)
	}

	req := &oapi.CreateKeypairRequest{
		KeypairName: keyName,
	}

	var result *oapi.CreateKeypairResponse
	var resp *oapi.POST_CreateKeypairResponses
	var err error
	err = resource.Retry(120*time.Second, func() *resource.RetryError {
		var err error
		resp, err = conn.POST_CreateKeypair(*req)

		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded:") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	var errString string

	if err != nil || resp.OK == nil {
		if err != nil {
			errString = err.Error()
		} else if resp.Code401 != nil {
			errString = fmt.Sprintf("Status Code: 401, %s", utils.ToJSONString(resp.Code401))
		} else if resp.Code400 != nil {
			errString = fmt.Sprintf("Status Code: 400, %s", utils.ToJSONString(resp.Code400))
		} else if resp.Code500 != nil {
			errString = fmt.Sprintf("Status: 500, %s", utils.ToJSONString(resp.Code500))
		}
		return fmt.Errorf("Error creating OAPIKeyPair: %s", errString)
	}

	result = resp.OK

	d.SetId(result.Keypair.KeypairName)
	d.Set("keypair_fingerprint", result.Keypair.KeypairFingerprint)

	return resourceOAPIKeyPairRead(d, meta)
}

func resourceOAPIKeyPairRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OAPI
	req := &oapi.ReadKeypairsRequest{
		Filters: oapi.FiltersKeypair{KeypairNames: []string{d.Id()}},
	}

	var response *oapi.ReadKeypairsResponse
	var resp *oapi.POST_ReadKeypairsResponses
	err := resource.Retry(120*time.Second, func() *resource.RetryError {
		var err error
		resp, err = conn.POST_ReadKeypairs(*req)

		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded:") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	var errString string

	if err != nil || resp.OK == nil {
		if err != nil {
			if strings.Contains(fmt.Sprint(err), "InvalidOAPIKeyPair.NotFound") {
				d.SetId("")
				return nil
			}
			errString = err.Error()
		} else if resp.Code401 != nil {
			errString = fmt.Sprintf("ErrorCode: 401, %s", utils.ToJSONString(resp.Code401))
		} else if resp.Code400 != nil {
			errString = fmt.Sprintf("ErrorCode: 400, %s", utils.ToJSONString(resp.Code400))
		} else if resp.Code500 != nil {
			errString = fmt.Sprintf("ErrorCode: 500, %s", utils.ToJSONString(resp.Code500))
		}

		return fmt.Errorf("Error retrieving OAPIKeyPair: %s", errString)
	}

	response = resp.OK

	for _, keyPair := range response.Keypairs {
		if keyPair.KeypairName == d.Id() {
			d.Set("keypair_name", keyPair.KeypairName)
			d.Set("keypair_fingerprint", keyPair.KeypairFingerprint)
			return nil
		}
	}

	return fmt.Errorf("Unable to find key pair within: %#v", response.Keypairs)
}

func resourceOAPIKeyPairDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OAPI

	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		request := &oapi.DeleteKeypairRequest{
			KeypairName: d.Id(),
		}

		var err error
		_, err = conn.POST_DeleteKeypair(*request)
		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded:") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return err
	}

	return err
}

func getOAPIKeyPairSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		// Attributes
		"keypair_fingerprint": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"private_key": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"keypair_name": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
	}
}
