package outscale

import (
<<<<<<< HEAD
	"context"
=======
>>>>>>> fbd6a594... chore: updated vendor files
	"fmt"
	"strings"
	"time"

<<<<<<< HEAD
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
=======
	"github.com/aws/aws-sdk-go/aws"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceOutscaleIamAccessKey() *schema.Resource {
	return &schema.Resource{
		Create: resourceOutscaleIamAccessKeyCreate,
		Read:   resourceOutscaleIamAccessKeyRead,
		Delete: resourceOutscaleIamAccessKeyDelete,
>>>>>>> fbd6a594... chore: updated vendor files
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
<<<<<<< HEAD
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "ACTIVE",
				ValidateFunc: validation.StringInSlice([]string{"ACTIVE", "INACTIVE"}, false),
=======
				Type:     schema.TypeString,
				Computed: true,
>>>>>>> fbd6a594... chore: updated vendor files
			},
			"request_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

<<<<<<< HEAD
func resourceOutscaleAccessKeyCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	var res oscgo.CreateAccessKeyResponse
	var err error
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		res, _, err = conn.AccessKeyApi.CreateAccessKey(context.Background())
=======
func resourceOutscaleIamAccessKeyCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	return resourceOutscaleIamAccessKeyRead(d, meta)
}

func resourceOutscaleIamAccessKeyRead(d *schema.ResourceData, meta interface{}) error {
	iamconn := meta.(*OutscaleClient).ICU

	request := &icu.ListAccessKeysInput{}

	var getResp *icu.ListAccessKeysOutput
	var err error
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		getResp, err = iamconn.API.ListAccessKeys(request)

>>>>>>> fbd6a594... chore: updated vendor files
		if err != nil {
			if strings.Contains(fmt.Sprint(err), "RequestLimitExceeded:") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
<<<<<<< HEAD
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
=======
		if strings.Contains(fmt.Sprint(err), "NoSuchEntity") {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error reading acces key: %s", err)
	}

	for _, key := range getResp.AccessKeyMetadata {
		if key.AccessKeyID != nil && *key.AccessKeyID == d.Id() {
			d.Set("access_key_id", key.AccessKeyID)
			d.Set("secret_access_key", key.SecretAccessKey)
			d.Set("owner_id", key.OwnerID)
			d.Set("status", key.Status)
			d.Set("tag_set", tagsToMapI(key.Tags))
			return d.Set("request_id", getResp.ResponseMetadata.RequestID)
		}
	}
	d.SetId("")
	return fmt.Errorf("AccessKey not found")
}

func resourceOutscaleIamAccessKeyDelete(d *schema.ResourceData, meta interface{}) error {
	iamconn := meta.(*OutscaleClient).ICU

	request := &icu.DeleteAccessKeyInput{
		AccessKeyId: aws.String(d.Id()),
>>>>>>> fbd6a594... chore: updated vendor files
	}

	var err error
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
<<<<<<< HEAD
		_, _, err = conn.AccessKeyApi.DeleteAccessKey(context.Background(), &oscgo.DeleteAccessKeyOpts{
			DeleteAccessKeyRequest: optional.NewInterface(req),
		})
=======
		_, err = iamconn.API.DeleteAccessKey(request)

>>>>>>> fbd6a594... chore: updated vendor files
		if err != nil {
			if strings.Contains(fmt.Sprint(err), "RequestLimitExceeded:") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
<<<<<<< HEAD
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

=======

	if err != nil {
		return fmt.Errorf("Error deleting access key %s: %s", d.Id(), err)
	}
>>>>>>> fbd6a594... chore: updated vendor files
	return nil
}
