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
	"github.com/outscale/terraform-provider-outscale/utils"
)

func ResourceOutscaleAccessKey() *schema.Resource {
	return &schema.Resource{
		Create: ResourceOutscaleAccessKeyCreate,
		Read:   ResourceOutscaleAccessKeyRead,
		Update: ResourceOutscaleAccessKeyUpdate,
		Delete: ResourceOutscaleAccessKeyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"user_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
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

func ResourceOutscaleAccessKeyCreate(d *schema.ResourceData, meta interface{}) error {
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
	if userName := d.Get("user_name").(string); userName != "" {
		req.SetUserName(userName)
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
	if err := d.Set("secret_key", *resp.GetAccessKey().SecretKey); err != nil {
		return err
	}
	if d.Get("state").(string) != "ACTIVE" {
		if err := inactiveAccessKey(d, conn, "INACTIVE"); err != nil {
			return err
		}
	}
	return ResourceOutscaleAccessKeyRead(d, meta)
}

func ResourceOutscaleAccessKeyRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	filter := oscgo.FiltersAccessKeys{
		AccessKeyIds: &[]string{d.Id()},
	}

	req := oscgo.ReadAccessKeysRequest{
		Filters: &filter,
	}
	if userName := d.Get("user_name").(string); userName != "" {
		req.SetUserName(userName)
	}
	var resp oscgo.ReadAccessKeysResponse
	err := resource.Retry(2*time.Minute, func() *resource.RetryError {
		rp, httpResp, err := conn.AccessKeyApi.ReadAccessKeys(context.Background()).ReadAccessKeysRequest(req).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		resp = rp
		return nil
	})

	if err != nil {
		return err
	}

	accessKey := resp.GetAccessKeys()
	if len(accessKey) == 0 {
		d.SetId("")
		return nil
	}
	if userName := d.Get("user_name").(string); userName != "" {
		if err := d.Set("user_name", userName); err != nil {
			return err
		}
	}
	if err := d.Set("access_key_id", accessKey[0].GetAccessKeyId()); err != nil {
		return err
	}
	if err := d.Set("creation_date", accessKey[0].GetCreationDate()); err != nil {
		return err
	}
	if err := d.Set("expiration_date", accessKey[0].GetExpirationDate()); err != nil {
		return err
	}
	if err := d.Set("last_modification_date", accessKey[0].GetLastModificationDate()); err != nil {
		return err
	}
	if err := d.Set("state", accessKey[0].GetState()); err != nil {
		return err
	}

	return nil
}

func ResourceOutscaleAccessKeyUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI
	req := oscgo.UpdateAccessKeyRequest{AccessKeyId: d.Id()}

	if expirDate, newdate := d.GetChange("expiration_date"); newdate.(string) != "" {
		req.SetExpirationDate(expirDate.(string))
	}
	if userName := d.Get("user_name").(string); userName != "" {
		req.SetUserName(userName)
	}
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
			req.SetExpirationDate(newExpirDate)
		} else {
			if !req.HasUserName() {
				if userName := d.Get("user_name").(string); userName != "" {
					req.SetUserName(userName)
				}
			}
		}
		if req.State == "" {
			req.SetState(state)
		}
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
	return ResourceOutscaleAccessKeyRead(d, meta)
}

func ResourceOutscaleAccessKeyDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	req := oscgo.DeleteAccessKeyRequest{
		AccessKeyId: d.Id(),
	}
	if userName := d.Get("user_name").(string); userName != "" {
		req.SetUserName(userName)
		if err := inactiveAccessKey(d, conn, "INACTIVE"); err != nil {
			return err
		}
	}
	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, httpResp, err := conn.AccessKeyApi.DeleteAccessKey(context.Background()).DeleteAccessKeyRequest(req).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf(" Error deleting Outscale Access Key %s: %s", d.Id(), err)
	}

	return nil
}

func inactiveAccessKey(d *schema.ResourceData, conn *oscgo.APIClient, state string) error {
	req := oscgo.UpdateAccessKeyRequest{
		AccessKeyId: d.Id(),
		State:       state,
	}
	if userName := d.Get("user_name").(string); userName != "" {
		req.SetUserName(userName)
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
		return fmt.Errorf(" Expiration Date should be 'ISO 8601' format ('2017-06-14' or '2017-06-14T00:00:00Z, ...) %s", err)
	}
	if currentDate.After(settingDate) {
		return fmt.Errorf(" Expiration date: '%s' should be after current date '%s'", settingDate, currentDate)
	}
	return nil
}
