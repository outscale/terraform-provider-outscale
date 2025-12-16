package outscale

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
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
				Required: true,
				ForceNew: true,
			},
			"path": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
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
	err := retry.Retry(2*time.Minute, func() *retry.RetryError {
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

	d.SetId(*resp.GetPolicy().Orn)
	return ResourceOutscalePolicyRead(d, meta)
}

func ResourceOutscalePolicyRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI
	req := oscgo.NewReadPolicyRequest(d.Id())

	var resp oscgo.ReadPolicyResponse
	err := retry.Retry(2*time.Minute, func() *retry.RetryError {
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
	policyDocument, err := getPolicyDocument(conn, policy.GetOrn(), policy.GetPolicyDefaultVersionId())
	if err != nil {
		return err
	}
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
	if err := d.Set("document", policyDocument); err != nil {
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
	policyOrn := d.Get("orn").(string)
	if err := unlinkEntitiesToPolicy(conn, policyOrn); err != nil {
		return err
	}

	req := oscgo.NewDeletePolicyRequest(policyOrn)
	err := retry.Retry(2*time.Minute, func() *retry.RetryError {
		_, httpResp, err := conn.PolicyApi.DeletePolicy(context.Background()).DeletePolicyRequest(*req).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("error deleting Outscale Policy %s: %s", d.Id(), err)
	}
	return nil
}

func unlinkEntitiesToPolicy(conn *oscgo.APIClient, policyOrn string) error {

	req := oscgo.ReadEntitiesLinkedToPolicyRequest{PolicyOrn: policyOrn}
	var users, groups []oscgo.MinimalPolicy
	err := retry.Retry(2*time.Minute, func() *retry.RetryError {
		resp, httpResp, err := conn.PolicyApi.ReadEntitiesLinkedToPolicy(context.Background()).ReadEntitiesLinkedToPolicyRequest(req).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		if resp.HasPolicyEntities() {
			users = *resp.GetPolicyEntities().Users
			groups = *resp.GetPolicyEntities().Groups
		}
		return nil
	})
	if err != nil {
		return err
	}
	if len(users) > 0 {
		unlinkReq := oscgo.UnlinkPolicyRequest{}
		unlinkReq.SetPolicyOrn(policyOrn)

		for _, user := range users {
			unlinkReq.SetUserName(*user.Name)
			err := retry.Retry(2*time.Minute, func() *retry.RetryError {
				_, httpResp, err := conn.PolicyApi.UnlinkPolicy(context.Background()).UnlinkPolicyRequest(unlinkReq).Execute()
				if err != nil {
					return utils.CheckThrottling(httpResp, err)
				}
				return nil
			})
			if err != nil {
				return err
			}
		}
	}
	if len(groups) > 0 {
		unlinkReq := oscgo.UnlinkManagedPolicyFromUserGroupRequest{PolicyOrn: policyOrn}
		for _, group := range groups {
			unlinkReq.SetUserGroupName(*group.Name)
			err := retry.Retry(2*time.Minute, func() *retry.RetryError {
				_, httpResp, err := conn.PolicyApi.UnlinkManagedPolicyFromUserGroup(context.Background()).UnlinkManagedPolicyFromUserGroupRequest(unlinkReq).Execute()
				if err != nil {
					return utils.CheckThrottling(httpResp, err)
				}
				return nil
			})
			if err != nil {
				return err
			}
		}
	}

	return err
}

func getPolicyDocument(conn *oscgo.APIClient, policyOrn, policyVersionId string) (string, error) {

	req := oscgo.NewReadPolicyVersionRequest(policyOrn, policyVersionId)
	var resp oscgo.ReadPolicyVersionResponse
	err := retry.Retry(2*time.Minute, func() *retry.RetryError {
		rp, httpResp, err := conn.PolicyApi.ReadPolicyVersion(context.Background()).ReadPolicyVersionRequest(*req).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		resp = rp
		return nil
	})
	if err != nil {
		return "", err
	}
	if _, ok := resp.GetPolicyVersionOk(); !ok {
		return "", fmt.Errorf("cannot find Policy version: %v", policyVersionId)
	}

	return *resp.GetPolicyVersion().Body, err
}
