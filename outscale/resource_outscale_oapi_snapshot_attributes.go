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
			"permissions_to_create_volume": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"additions": &schema.Schema{
							Type:     schema.TypeList,
							Optional: true,
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
						"removals": &schema.Schema{
							Type:     schema.TypeList,
							Optional: true,
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
					},
				},
			},
			"snapshot_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
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

	if permsParam, ok := d.GetOk("permissions_to_create_volume"); ok {
		perms := oapi.PermissionsOnResourceCreation{}

		if additions, additionsOk := permsParam.(map[string]interface{})["additions"]; additionsOk {
			perms.Additions = oapi.PermissionsOnResource{}
			if accountIdsParam, accountIdsOk := additions.(map[string]interface{})["account_ids"]; accountIdsOk {
				perms.Additions.AccountIds = expandAccountIds(accountIdsParam)
			}
			if globalPermsParam, globalPermsOk := additions.(map[string]interface{})["global_permission"]; globalPermsOk {
				perms.Additions.GlobalPermission = globalPermsParam.(bool)
			}
		}

		if removals, removalsOk := permsParam.(map[string]interface{})["removals"]; removalsOk {
			perms.Removals = oapi.PermissionsOnResource{}
			if accountIdsParam, accountIdsOk := removals.(map[string]interface{})["account_ids"]; accountIdsOk {
				perms.Removals.AccountIds = expandAccountIds(accountIdsParam)
			}
			if globalPermsParam, globalPermsOk := removals.(map[string]interface{})["global_permission"]; globalPermsOk {
				perms.Removals.GlobalPermission = globalPermsParam.(bool)
			}
		}

		req.PermissionsToCreateVolume = perms
	}

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

	permsMap := make(map[string]interface{})
	permsMap["account_ids"] = attrs.OK.Snapshots[0].PermissionsToCreateVolume.AccountIds
	permsMap["global_permission"] = attrs.OK.Snapshots[0].PermissionsToCreateVolume.GlobalPermission

	d.Set("request_id", attrs.OK.ResponseContext.RequestId)

	return d.Set("permissions_to_create_volume", permsMap)
}

func resourcedOutscaleOAPISnapshotAttributesDelete(d *schema.ResourceData, meta interface{}) error {
	d.SetId("")

	return nil
}
