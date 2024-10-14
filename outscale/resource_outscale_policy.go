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

func ResourceOutscalePolicy() *schema.Resource {
	return &schema.Resource{
		Create: ResourceOutscalePolicyCreate,
		Read:   ResourceOutscalePolicyRead,
		Delete: ResourceOutscalePolicyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"policy_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"document": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"path": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"policy_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"resources_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"policy_default_version_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_linkable": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"orn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"creation_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"last_modification_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func ResourceOutscalePolicyCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI
	polDocument := d.Get("document").(string)
	req := oscgo.NewCreatePolicyRequest(polDocument, d.Get("policy_name").(string))
	if polPath := d.Get("path").(string); polPath != "" {
		req.SetPath(polPath)
	}
	if polDescription := d.Get("description").(string); polDescription != "" {
		req.SetDescription(polDescription)
	}

	var resp oscgo.CreatePolicyResponse
	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		rp, httpResp, err := conn.PolicyApi.CreatePolicy(context.Background()).CreatePolicyRequest(*req).Execute()
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
	ply := resp.GetPolicy()
	if err := d.Set("orn", ply.GetOrn()); err != nil {
		return err
	}

	// Remove d.Set when read user_group return user_group_id
	return ResourceOutscalePolicyRead(d, meta)
}

func ResourceOutscalePolicyRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI
	req := oscgo.NewReadPolicyRequest(d.Get("orn").(string))

	var resp oscgo.ReadPolicyResponse
	err := resource.Retry(2*time.Minute, func() *resource.RetryError {
		rp, httpResp, err := conn.PolicyApi.ReadPolicy(context.Background()).ReadPolicyRequest(*req).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		resp = rp
		return nil
	})

	if err != nil {
		return err
	}

	if _, ok := resp.GetPolicyOk(); !ok {
		d.SetId("")
		return nil
	}
	policy := resp.GetPolicy()
	if err := d.Set("policy_name", policy.GetPolicyName()); err != nil {
		return err
	}
	if err := d.Set("policy_id", policy.GetPolicyId()); err != nil {
		return err
	}
	if err := d.Set("path", policy.GetPath()); err != nil {
		return err
	}
	if err := d.Set("orn", policy.GetOrn()); err != nil {
		return err
	}
	if err := d.Set("resources_count", policy.GetResourcesCount()); err != nil {
		return err
	}
	if err := d.Set("is_linkable", policy.GetIsLinkable()); err != nil {
		return err
	}
	if err := d.Set("policy_default_version_id", policy.GetPolicyDefaultVersionId()); err != nil {
		return err
	}
	if err := d.Set("description", policy.GetDescription()); err != nil {
		return err
	}

	if err := d.Set("creation_date", (policy.GetCreationDate())); err != nil {
		return err
	}

	if err := d.Set("last_modification_date", (policy.GetLastModificationDate())); err != nil {
		return err
	}
	return nil
}

func ResourceOutscalePolicyDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI
	req := oscgo.NewDeletePolicyRequest(d.Get("orn").(string))
	var err error
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, httpResp, err := conn.PolicyApi.DeletePolicy(context.Background()).DeletePolicyRequest(*req).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("Error deleting Outscale Policy %s: %s", d.Id(), err)
	}

	return nil
}
