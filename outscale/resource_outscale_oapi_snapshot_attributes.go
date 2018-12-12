package outscale

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-outscale/osc/oapi"
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
			"permissions_to_create_volume": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"account_ids": &schema.Schema{
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"global_permission": &schema.Schema{
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
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
	conn := meta.(*OutscaleClient).OAPI

	snapshotID := d.Get("snapshot_id").(string)

	req := oapi.UpdateSnapshotRequest{
		SnapshotId: snapshotID,
	}

	perms := oapi.PermissionsOnResourceCreation{}

	if addPermsParam, ok := d.GetOk("permissions_to_create_volume_additions"); ok {
		AddPerms := addPermsParam.([]interface{})

		if len(AddPerms) > 0 {

			perms.Additions = oapi.PermissionsOnResource{}

			addMap := AddPerms[0].(map[string]interface{})
			if addMap["account_ids"] != nil {
				paramIds := addMap["account_ids"].([]interface{})
				accountIds := make([]string, len(paramIds))
				for i, v := range paramIds {
					accountIds[i] = v.(string)
				}
				perms.Additions.AccountIds = accountIds
			}
			if addMap["global_permission"] != nil {
				globalPermission := addMap["global_permission"].(bool)
				perms.Additions.GlobalPermission = globalPermission
			}
		}
	}

	if removalPermsParam, ok := d.GetOk("permissions_to_create_volume_removals"); ok {
		removalPerms := removalPermsParam.([]interface{})

		if len(removalPerms) > 0 {

			perms.Removals = oapi.PermissionsOnResource{}

			removalMap := removalPerms[0].(map[string]interface{})
			if removalMap["account_ids"] != nil {
				paramIds := removalMap["account_ids"].([]interface{})
				accountIds := make([]string, len(paramIds))
				for i, v := range paramIds {
					accountIds[i] = v.(string)
				}
				perms.Removals.AccountIds = accountIds
			}
			if removalMap["global_permission"] != nil {
				globalPermission := removalMap["global_permission"].(bool)
				perms.Additions.GlobalPermission = globalPermission
			}
		}
	}

	req.PermissionsToCreateVolume = perms

	var err error
	err = resource.Retry(2*time.Minute, func() *resource.RetryError {
		_, err = conn.POST_UpdateSnapshot(req)
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
	conn := meta.(*OutscaleClient).OAPI

	var attrs *oapi.POST_ReadSnapshotsResponses
	var err error
	err = resource.Retry(2*time.Minute, func() *resource.RetryError {
		attrs, err = conn.POST_ReadSnapshots(oapi.ReadSnapshotsRequest{
			Filters: oapi.FiltersSnapshot{
				SnapshotIds: []string{d.Id()},
			},
		})
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
	lp[0]["global_permission"] = attrs.OK.Snapshots[0].PermissionsToCreateVolume.GlobalPermission
	lp[0]["account_ids"] = attrs.OK.Snapshots[0].PermissionsToCreateVolume.AccountIds

	if err := d.Set("permissions_to_create_volume", lp); err != nil {
		return err
	}

	d.Set("request_id", attrs.OK.ResponseContext.RequestId)

	return nil
}

func resourcedOutscaleOAPISnapshotAttributesDelete(d *schema.ResourceData, meta interface{}) error {
	d.SetId("")

	return nil
}
