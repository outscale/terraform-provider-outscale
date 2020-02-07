package outscale

import (
	"context"
	"fmt"
	"github.com/antihax/optional"
	oscgo "github.com/marinsalinas/osc-sdk-go"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourcedOutscaleOAPISnapshotAttributes() *schema.Resource {
	return &schema.Resource{
		Create: resourcedOutscaleOAPISnapshotAttributesCreate,
		Read:   resourcedOutscaleOAPISnapshotAttributesRead,
		Delete: resourcedOutscaleOAPISnapshotAttributesDelete,

		Schema: map[string]*schema.Schema{
			"permissions_to_create_volume_additions": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"account_ids": &schema.Schema{
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"global_permission": &schema.Schema{
							Type:     schema.TypeBool,
							Optional: true,
						},
					},
				},
			},
			"permissions_to_create_volume_removals": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"account_ids": &schema.Schema{
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"global_permission": &schema.Schema{
							Type:     schema.TypeBool,
							Optional: true,
						},
					},
				},
			},
			"snapshot_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"account_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"request_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func expandAccountIds(param interface{}) []string {
	var values []string
	for _, v := range param.([]interface{}) {
		values = append(values, v.(string))
	}
	return values
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
		_, _, err = conn.SnapshotApi.UpdateSnapshot(context.Background(), &oscgo.UpdateSnapshotOpts{UpdateSnapshotRequest: optional.NewInterface(req)})
		if err != nil {
			if strings.Contains(fmt.Sprint(err), "RequestLimitExceeded") {
				log.Printf("[DEBUG] Error: %q", err)
				return resource.RetryableError(err)
			}

			return resource.NonRetryableError(err)
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
	var err error
	err = resource.Retry(2*time.Minute, func() *resource.RetryError {
		resp, _, err = conn.SnapshotApi.ReadSnapshots(context.Background(), &oscgo.ReadSnapshotsOpts{ReadSnapshotsRequest: optional.NewInterface(oscgo.ReadSnapshotsRequest{
			Filters: &oscgo.FiltersSnapshot{
				SnapshotIds: &[]string{d.Id()},
			},
		})})
		if err != nil {
			if strings.Contains(fmt.Sprint(err), "RequestLimitExceeded") {
				log.Printf("[DEBUG] Error: %q", err)
				return resource.RetryableError(err)
			}

			return resource.NonRetryableError(err)
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("Error refreshing snapshot createVolumePermission state: %s", err)
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
	if err := d.Set("request_id", resp.ResponseContext.GetRequestId()); err != nil {
		return err
	}

	return nil
}

func resourcedOutscaleOAPISnapshotAttributesDelete(d *schema.ResourceData, meta interface{}) error {
	d.SetId("")

	return nil
}
