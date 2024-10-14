package outscale

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/outscale/terraform-provider-outscale/utils"
)

func ResourcePolicyVersion() *schema.Resource {
	return &schema.Resource{
		Create: ResourcePolicyVersionCreate,
		Read:   ResourcePolicyVersionRead,
		Delete: ResourcePolicyVersionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"policy_orn": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"document": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"set_as_default": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"version_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"default_version": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"creation_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"body": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func ResourcePolicyVersionCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI
	pvDocument := d.Get("document").(string)
	req := oscgo.NewCreatePolicyVersionRequest(pvDocument, d.Get("policy_orn").(string))

	if asDefault, ok := d.GetOk("set_as_default"); ok {
		req.SetSetAsDefault(asDefault.(bool))
	}

	var resp oscgo.CreatePolicyVersionResponse
	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		rp, httpResp, err := conn.PolicyApi.CreatePolicyVersion(context.Background()).CreatePolicyVersionRequest(*req).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		resp = rp
		return nil
	})
	if err != nil {
		return err
	}

	d.SetId(resource.UniqueId())
	pVersion := resp.GetPolicyVersion()
	if err := d.Set("version_id", pVersion.GetVersionId()); err != nil {
		return err
	}

	// Remove d.Set when read user_group return user_group_id
	return ResourcePolicyVersionRead(d, meta)
}

func ResourcePolicyVersionRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI
	req := oscgo.NewReadPolicyVersionRequest(d.Get("policy_orn").(string), d.Get("version_id").(string))

	var resp oscgo.ReadPolicyVersionResponse
	err := resource.Retry(2*time.Minute, func() *resource.RetryError {
		rp, httpResp, err := conn.PolicyApi.ReadPolicyVersion(context.Background()).ReadPolicyVersionRequest(*req).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		resp = rp
		return nil
	})

	if err != nil {
		return err
	}

	if _, ok := resp.GetPolicyVersionOk(); !ok {
		d.SetId("")
		return nil
	}
	pVersion := resp.GetPolicyVersion()
	if err := d.Set("default_version", pVersion.GetDefaultVersion()); err != nil {
		return err
	}
	if err := d.Set("creation_date", (pVersion.GetCreationDate())); err != nil {
		return err
	}
	if err := d.Set("body", pVersion.GetBody()); err != nil {
		return err
	}
	/*
		usrs := resp.GetUsers()
		users := make([]map[string]interface{}, len(usrs))
		if len(usrs) > 0 {
			usrs := resp.GetUsers()
			for i, v := range usrs {
				user := make(map[string]interface{})
				user["user_id"] = v.GetUserId()
				user["user_name"] = v.GetUserName()
				user["path"] = v.GetPath()
				user["creation_date"] = v.GetCreationDate()
				user["last_modification_date"] = v.GetLastModificationDate()
				users[i] = user
			}
		}
		if err := d.Set("users", users); err != nil {
			return err
		}*/
	return nil
}

func ResourcePolicyVersionDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI
	req := oscgo.NewDeletePolicyVersionRequest(d.Get("policy_orn").(string), d.Get("version_id").(string))
	var err error
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, httpResp, err := conn.PolicyApi.DeletePolicyVersion(context.Background()).DeletePolicyVersionRequest(*req).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("Error deleting Outscale Policy version %s: %s", d.Id(), err)
	}

	return nil
}
