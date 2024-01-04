package outscale

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/nav-inc/datetime"
	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/terraform-providers/terraform-provider-outscale/utils"
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
			"access_key_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"creation_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"expiration_date": {
				Type:     schema.TypeString,
				Optional: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					date1, _ := datetime.Parse(new, time.UTC)
					date2, _ := datetime.Parse(old, time.UTC)
					return date1.Equal(date2)
				},
			},
			"last_modification_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"secret_key": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"state": {
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

	var err error
	req := oscgo.CreateAccessKeyRequest{}

	expirDate := d.Get("expiration_date").(string)
	if expirDate != "" {
		if err = checkDateFormat(expirDate); err != nil {
			return err
		}
		req.ExpirationDate = &expirDate
	}

	var resp oscgo.CreateAccessKeyResponse
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		rp, httpResp, err := conn.AccessKeyApi.CreateAccessKey(context.Background()).CreateAccessKeyRequest(req).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		resp = rp
		return nil
	})
	if err != nil {
		return err
	}

	d.SetId(*resp.GetAccessKey().AccessKeyId)

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
	var resp oscgo.ReadSecretAccessKeyResponse
	err := resource.Retry(2*time.Minute, func() *resource.RetryError {
		rp, httpResp, err := conn.AccessKeyApi.ReadSecretAccessKey(context.Background()).ReadSecretAccessKeyRequest(filter).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		resp = rp
		return nil
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
	if err := d.Set("expiration_date", accessKey.GetExpirationDate()); err != nil {
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

	return nil
}

func resourceOutscaleAccessKeyUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI
	req := oscgo.UpdateAccessKeyRequest{AccessKeyId: d.Id()}

	if d.HasChange("state") {
		req.State = d.Get("state").(string)
	}
	if d.HasChange("expiration_date") {
		newExpirDate := d.Get("expiration_date").(string)
		state := d.Get("state").(string)
		if newExpirDate != "" {
			if err := checkDateFormat(newExpirDate); err != nil {
				return err
			}
			req.ExpirationDate = &newExpirDate
		}
		req.State = state
	}

	err := resource.Retry(2*time.Minute, func() *resource.RetryError {
		_, httpResp, err := conn.AccessKeyApi.UpdateAccessKey(context.Background()).UpdateAccessKeyRequest(req).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		return nil
	})
	if err != nil {
		return err
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
		_, httpResp, err := conn.AccessKeyApi.DeleteAccessKey(context.Background()).DeleteAccessKeyRequest(req).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
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
	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, httpResp, err := conn.AccessKeyApi.UpdateAccessKey(context.Background()).UpdateAccessKeyRequest(req).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		return nil
	})

	if err != nil {
		return err
	}
	return nil
}

func checkDateFormat(dateFormat string) error {
	var err error
	var settingDate time.Time
	currentDate := time.Now()

	if settingDate, err = datetime.Parse(dateFormat, time.UTC); err != nil {
		return fmt.Errorf("Expiration Date should be 'ISO 8601' format ('2017-06-14' or '2017-06-14T00:00:00Z, ...) %s", err)
	}
	if currentDate.After(settingDate) {
		return fmt.Errorf(" Expiration date: '%s' should be after current date '%s'", settingDate, currentDate)
	}
	return nil
}
