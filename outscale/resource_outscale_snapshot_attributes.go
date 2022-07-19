package outscale

import (
	"context"
	"fmt"
	"time"

	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/terraform-providers/terraform-provider-outscale/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourcedOutscaleOAPISnapshotAttributes() *schema.Resource {
	return &schema.Resource{
		Create: resourcedOutscaleOAPISnapshotAttributesCreate,
		Read:   resourcedOutscaleOAPISnapshotAttributesRead,
		Delete: resourcedOutscaleOAPISnapshotAttributesDelete,

		Schema: map[string]*schema.Schema{
			"permissions_to_create_volume_additions": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"account_ids": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"global_permission": {
							Type:     schema.TypeBool,
							Optional: true,
						},
					},
				},
			},
			"permissions_to_create_volume_removals": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"account_ids": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"global_permission": {
							Type:     schema.TypeBool,
							Optional: true,
						},
					},
				},
			},
			"snapshot_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"account_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"request_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourcedOutscaleOAPISnapshotAttributesCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	snapshotID := d.Get("snapshot_id").(string)

	req := oscgo.UpdateSnapshotRequest{
		SnapshotId: snapshotID,
	}

	perms := oscgo.PermissionsOnResourceCreation{}

	if addPermsParam, ok := d.GetOk("permissions_to_create_volume_additions"); ok {
		AddPerms := addPermsParam.([]interface{})
		addition := oscgo.PermissionsOnResource{}
		if len(AddPerms) > 0 {
			perms.SetAdditions(addition)

			addMap := AddPerms[0].(map[string]interface{})
			if addMap["account_ids"] != nil {
				paramIds := addMap["account_ids"].([]interface{})
				accountIds := make([]string, len(paramIds))
				for i, v := range paramIds {
					accountIds[i] = v.(string)
				}
				addition.SetAccountIds(accountIds)
				perms.SetAdditions(addition)
			}
			if addMap["global_permission"] != nil {
				globalPermission := addMap["global_permission"].(bool)
				addition.SetGlobalPermission(globalPermission)
				perms.SetAdditions(addition)
			}
		}
	}

	if removalPermsParam, ok := d.GetOk("permissions_to_create_volume_removals"); ok {
		removalPerms := removalPermsParam.([]interface{})

		if len(removalPerms) > 0 {
			removal := oscgo.PermissionsOnResource{}
			perms.SetRemovals(removal)

			removalMap := removalPerms[0].(map[string]interface{})
			if removalMap["account_ids"] != nil {
				paramIds := removalMap["account_ids"].([]interface{})
				accountIds := make([]string, len(paramIds))
				for i, v := range paramIds {
					accountIds[i] = v.(string)
				}
				removal.SetAccountIds(accountIds)
				perms.SetRemovals(removal)
			}
			if removalMap["global_permission"] != nil {
				globalPermission := removalMap["global_permission"].(bool)
				removal.SetGlobalPermission(globalPermission)
				perms.SetRemovals(removal)
			}
		}
	}

	req.SetPermissionsToCreateVolume(perms)

	var err error
	err = resource.Retry(2*time.Minute, func() *resource.RetryError {
		_, httpResp, err := conn.SnapshotApi.UpdateSnapshot(context.Background()).UpdateSnapshotRequest(req).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("Error createing snapshot createVolumePermission: %s", err)
	}
	d.SetId(snapshotID)

	return resourcedOutscaleOAPISnapshotAttributesRead(d, meta)
}

func resourcedOutscaleOAPISnapshotAttributesRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	var resp oscgo.ReadSnapshotsResponse
	err := resource.Retry(2*time.Minute, func() *resource.RetryError {
		var err error
		rp, httpResp, err := conn.SnapshotApi.ReadSnapshots(context.Background()).ReadSnapshotsRequest(oscgo.ReadSnapshotsRequest{
			Filters: &oscgo.FiltersSnapshot{
				SnapshotIds: &[]string{d.Id()},
			},
		}).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		resp = rp
		return nil
	})
	if err != nil {
		return fmt.Errorf("Error refreshing snapshot createVolumePermission state: %s", err)
	}
	if utils.IsResponseEmpty(len(resp.GetSnapshots()), "SnapshotAtribute", d.Id()) {
		d.SetId("")
		return nil
	}
	lp := make([]map[string]interface{}, 1)
	lp[0] = make(map[string]interface{})
	lp[0]["global_permission"] = resp.GetSnapshots()[0].PermissionsToCreateVolume.GetGlobalPermission()
	lp[0]["account_ids"] = resp.GetSnapshots()[0].PermissionsToCreateVolume.GetAccountIds()

	if err := d.Set("permissions_to_create_volume_additions", lp); err != nil {
		return err
	}
	if err := d.Set("account_id", resp.GetSnapshots()[0].GetAccountId()); err != nil {
		return err
	}

	return nil
}

func resourcedOutscaleOAPISnapshotAttributesDelete(d *schema.ResourceData, meta interface{}) error {
	d.SetId("")

	return nil
}
