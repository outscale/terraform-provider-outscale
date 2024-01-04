package outscale

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/terraform-providers/terraform-provider-outscale/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
	conn := meta.(*OutscaleClient).OSCAPI

	var keyName string
	if v, ok := d.GetOk("keypair_name"); ok {
		keyName = v.(string)
	} else {
		keyName = resource.UniqueId()
		if err := d.Set("keypair_name", keyName); err != nil {
			return err
		}
	}

	req := oscgo.CreateKeypairRequest{
		KeypairName: keyName,
	}

	//Accept public key as argument
	if v, ok := d.GetOk("public_key"); ok {
		req.SetPublicKey(v.(string))
	}

	var resp oscgo.CreateKeypairResponse
	var err error
	err = resource.Retry(120*time.Second, func() *resource.RetryError {
		rp, httpResp, err := conn.KeypairApi.CreateKeypair(context.Background()).CreateKeypairRequest(req).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		resp = rp
		return nil
	})

	var errString string
	log.Printf("[DEBUG] resp keypair: %+v", resp)
	if err != nil {
		errString = err.Error()
		return fmt.Errorf("Error creating OAPIKeyPair: %s", errString)
	}

	d.SetId(resp.Keypair.GetKeypairName())
	if err := d.Set("keypair_fingerprint", resp.Keypair.GetKeypairFingerprint()); err != nil {
		return err
	}

	//Set private key in creation
	if resp.Keypair.GetPrivateKey() != "" {
		if err := d.Set("private_key", resp.Keypair.GetPrivateKey()); err != nil {
			return err
		}
	}

	return resourceOAPIKeyPairRead(d, meta)
}

func resourceOAPIKeyPairRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI
	req := oscgo.ReadKeypairsRequest{
		Filters: &oscgo.FiltersKeypair{KeypairNames: &[]string{d.Id()}},
	}

	var resp oscgo.ReadKeypairsResponse
	err := resource.Retry(120*time.Second, func() *resource.RetryError {
		var err error
		rp, httpResp, err := conn.KeypairApi.ReadKeypairs(context.Background()).ReadKeypairsRequest(req).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		resp = rp
		return nil
	})

	var errString string

	if err != nil {
		if strings.Contains(fmt.Sprint(err), "InvalidOAPIKeyPair.NotFound") {
			d.SetId("")
			return nil
		}
		errString = err.Error()

		return fmt.Errorf("Error retrieving OAPIKeyPair: %s", errString)
	}
	for _, keyPair := range resp.GetKeypairs() {
		if keyPair.GetKeypairName() == d.Id() {
			if err := d.Set("keypair_name", keyPair.GetKeypairName()); err != nil {
				return err
			}
			if err := d.Set("keypair_fingerprint", keyPair.GetKeypairFingerprint()); err != nil {
				return err
			}
			return nil
		}
	}
	utils.LogManuallyDeleted("Keypair", d.Id())
	d.SetId("")
	return nil
}

func resourceOAPIKeyPairDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		request := oscgo.DeleteKeypairRequest{
			KeypairName: d.Id(),
		}

		var err error
		_, httpResp, err := conn.KeypairApi.DeleteKeypair(context.Background()).DeleteKeypairRequest(request).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
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
			Computed: true,
		},
		"private_key": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"keypair_name": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"public_key": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"request_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}
