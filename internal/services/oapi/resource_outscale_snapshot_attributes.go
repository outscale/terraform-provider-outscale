package oapi

import (
	"context"

	"github.com/outscale/goutils/sdk/ptr"
	"github.com/outscale/osc-sdk-go/v3/pkg/options"
	"github.com/outscale/osc-sdk-go/v3/pkg/osc"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceOutscaleSnapshotAttributes() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceOutscaleSnapshotAttributesCreate,
		ReadContext:   ResourceOutscaleSnapshotAttributesRead,
		DeleteContext: ResourceOutscaleSnapshotAttributesDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(CreateDefaultTimeout),
			Read:   schema.DefaultTimeout(ReadDefaultTimeout),
			Delete: schema.DefaultTimeout(DeleteDefaultTimeout),
		},

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

func ResourceOutscaleSnapshotAttributesCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.OutscaleClient).OSC

	timeout := d.Timeout(schema.TimeoutCreate)

	snapshotID := d.Get("snapshot_id").(string)

	req := osc.UpdateSnapshotRequest{
		SnapshotId: snapshotID,
	}

	perms := osc.PermissionsOnResourceCreation{}

	if addPermsParam, ok := d.GetOk("permissions_to_create_volume_additions"); ok {
		AddPerms := addPermsParam.([]interface{})
		addition := osc.PermissionsOnResource{}
		if len(AddPerms) > 0 {
			perms.Additions = &addition

			addMap := AddPerms[0].(map[string]interface{})
			if addMap["account_ids"] != nil {
				paramIds := addMap["account_ids"].([]interface{})
				accountIds := make([]string, len(paramIds))
				for i, v := range paramIds {
					accountIds[i] = v.(string)
				}
				addition.AccountIds = &accountIds
				perms.Additions = &addition
			}
			if addMap["global_permission"] != nil {
				globalPermission := addMap["global_permission"].(bool)
				addition.GlobalPermission = &globalPermission
				perms.Additions = &addition
			}
		}
	}

	if removalPermsParam, ok := d.GetOk("permissions_to_create_volume_removals"); ok {
		removalPerms := removalPermsParam.([]interface{})

		if len(removalPerms) > 0 {
			removal := osc.PermissionsOnResource{}
			perms.Removals = &removal

			removalMap := removalPerms[0].(map[string]interface{})
			if removalMap["account_ids"] != nil {
				paramIds := removalMap["account_ids"].([]interface{})
				accountIds := make([]string, len(paramIds))
				for i, v := range paramIds {
					accountIds[i] = v.(string)
				}
				removal.AccountIds = &accountIds
				perms.Removals = &removal
			}
			if removalMap["global_permission"] != nil {
				globalPermission := removalMap["global_permission"].(bool)
				removal.GlobalPermission = &globalPermission
				perms.Removals = &removal
			}
		}
	}

	req.PermissionsToCreateVolume = perms

	_, err := client.UpdateSnapshot(ctx, req, options.WithRetryTimeout(timeout))
	if err != nil {
		return diag.Errorf("error createing snapshot createvolumepermission: %s", err)
	}
	d.SetId(snapshotID)

	return ResourceOutscaleSnapshotAttributesRead(ctx, d, meta)
}

func ResourceOutscaleSnapshotAttributesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.OutscaleClient).OSC

	timeout := d.Timeout(schema.TimeoutRead)

	resp, err := client.ReadSnapshots(ctx, osc.ReadSnapshotsRequest{
		Filters: &osc.FiltersSnapshot{
			SnapshotIds: &[]string{d.Id()},
		},
	}, options.WithRetryTimeout(timeout))
	if err != nil {
		return diag.Errorf("error refreshing snapshot createvolumepermission state: %s", err)
	}
	if resp.Snapshots == nil || utils.IsResponseEmpty(len(*resp.Snapshots), "SnapshotAtribute", d.Id()) {
		d.SetId("")
		return nil
	}
	lp := make([]map[string]interface{}, 1)
	lp[0] = make(map[string]interface{})
	lp[0]["global_permission"] = ptr.From((*resp.Snapshots)[0].PermissionsToCreateVolume).GlobalPermission
	lp[0]["account_ids"] = ptr.From((*resp.Snapshots)[0].PermissionsToCreateVolume).AccountIds

	if err := d.Set("permissions_to_create_volume_additions", lp); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("account_id", (*resp.Snapshots)[0].AccountId); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func ResourceOutscaleSnapshotAttributesDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	d.SetId("")

	return nil
}
