package outscale

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/antihax/optional"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	oscgo "github.com/marinsalinas/osc-sdk-go"
)

func resourceOutscaleAccessKey() *schema.Resource {
	return &schema.Resource{
		Create: resourceOutscaleAccessKeyCreate,
		Read:   resourceOutscaleAccessKeyRead,
		Update: resourceOutscaleAccessKeyUpdate,
		Delete: resourceOutscaleAccessKeyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"access_key_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"creation_date": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"last_modification_date": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"secret_key": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"state": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "ACTIVE",
				ValidateFunc: validation.StringInSlice([]string{"ACTIVE", "INACTIVE"}, false),
			},
			"request_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceOutscaleAccessKeyCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	var res oscgo.CreateAccessKeyResponse
	var err error
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		res, _, err = conn.AccessKeyApi.CreateAccessKey(context.Background())
		if err != nil {
			if strings.Contains(fmt.Sprint(err), "RequestLimitExceeded:") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return err
	}

	d.SetId(*res.GetAccessKey().AccessKeyId)

	if d.Get("state").(string) != "ACTIVE" {
		if err := updateAccessKey(conn, d.Id(), "INACTIVE"); err != nil {
			return err
		}
	}

	return resourceOutscaleAccessKeyRead(d, meta)
}

func resourceOutscaleAccessKeyRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	filter := oscgo.ReadSecretAccessKeyRequest{
		AccessKeyId: d.Id(),
	}

	resp, _, err := conn.AccessKeyApi.ReadSecretAccessKey(context.Background(), &oscgo.ReadSecretAccessKeyOpts{
		ReadSecretAccessKeyRequest: optional.NewInterface(filter),
	})
	if err != nil {
		return err
	}

	accessKey, ok := resp.GetAccessKeyOk()
	if !ok {
		d.SetId("")
		return nil
	}

	if err := d.Set("access_key_id", accessKey.GetAccessKeyId()); err != nil {
		return err
	}
	if err := d.Set("creation_date", accessKey.GetCreationDate()); err != nil {
		return err
	}
	if err := d.Set("last_modification_date", accessKey.GetLastModificationDate()); err != nil {
		return err
	}
	if err := d.Set("secret_key", accessKey.GetSecretKey()); err != nil {
		return err
	}
	if err := d.Set("state", accessKey.GetState()); err != nil {
		return err
	}
	if err := d.Set("request_id", resp.ResponseContext.GetRequestId()); err != nil {
		return err
	}

	return nil
}

func resourceOutscaleAccessKeyUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	if d.HasChange("state") {
		if err := updateAccessKey(conn, d.Id(), d.Get("state").(string)); err != nil {
			return err
		}
	}

	return resourceOutscaleAccessKeyRead(d, meta)
}

func resourceOutscaleAccessKeyDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	req := oscgo.DeleteAccessKeyRequest{
		AccessKeyId: d.Id(),
	}

	var err error
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, _, err = conn.AccessKeyApi.DeleteAccessKey(context.Background(), &oscgo.DeleteAccessKeyOpts{
			DeleteAccessKeyRequest: optional.NewInterface(req),
		})
		if err != nil {
			if strings.Contains(fmt.Sprint(err), "RequestLimitExceeded:") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("Error deleting Outscale Access Key %s: %s", d.Id(), err)
	}

	return nil
}

func updateAccessKey(conn *oscgo.APIClient, id, state string) error {
	req := oscgo.UpdateAccessKeyRequest{
		AccessKeyId: id,
		State:       state,
	}

	_, _, err := conn.AccessKeyApi.UpdateAccessKey(context.Background(), &oscgo.UpdateAccessKeyOpts{
		UpdateAccessKeyRequest: optional.NewInterface(req),
	})
	if err != nil {
		return err
	}

	return nil
}
