package oapi

import (
	"context"
	"time"

	"github.com/outscale/goutils/sdk/ptr"
	"github.com/outscale/osc-sdk-go/v3/pkg/options"
	"github.com/outscale/osc-sdk-go/v3/pkg/osc"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/framework/fwhelpers/from"
	"github.com/outscale/terraform-provider-outscale/internal/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceOutscaleSnapshots() *schema.Resource {
	return &schema.Resource{
		ReadContext: DataSourceOutscaleSnapshotsRead,

		Schema: map[string]*schema.Schema{
			// selection criteria
			"filter": dataSourceFiltersSchema(),
			"account_id": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"snapshot_id": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"permission_to_create_volume": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			// Computed values returned
			"snapshots": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"progress": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"snapshot_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"volume_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"state": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"account_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"account_alias": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"creation_date": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"volume_size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"permissions_to_create_volume": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"account_ids": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"global_permission": {
										Type:     schema.TypeBool,
										Computed: true,
									},
								},
							},
						},
						"tags": TagsSchemaComputedSDK(),
					},
				},
			},
			"request_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func DataSourceOutscaleSnapshotsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.OutscaleClient).OSC

	restorableUsers, restorableUsersOk := d.GetOk("permission_to_create_volume")
	filters, filtersOk := d.GetOk("filter")
	snapshotIds, snapshotIdsOk := d.GetOk("snapshot_id")
	owners, ownersOk := d.GetOk("account_id")

	if restorableUsers == false && !filtersOk && snapshotIds == false && !ownersOk {
		return diag.Errorf("one of snapshot_ids, filters, restorable_by_user_ids, or owners must be assigned")
	}

	params := osc.ReadSnapshotsRequest{
		Filters: &osc.FiltersSnapshot{},
	}
	filter := osc.FiltersSnapshot{}
	if restorableUsersOk {
		filter.PermissionsToCreateVolumeAccountIds = new(utils.InterfaceSliceToStringSlice(restorableUsers.([]interface{})))
		params.Filters = &filter
	}
	if ownersOk {
		filter.AccountIds = new(utils.InterfaceSliceToStringSlice(owners.([]interface{})))
		params.Filters = &filter
	}
	if snapshotIdsOk {
		filter.SnapshotIds = new(utils.InterfaceSliceToStringSlice(snapshotIds.([]interface{})))
		params.Filters = &filter
	}

	var err error
	if filtersOk {
		params.Filters, err = buildOutscaleOapiSnapshootDataSourceFilters(filters.(*schema.Set))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	resp, err := client.ReadSnapshots(ctx, params, options.WithRetryTimeout(5*time.Minute))
	if err != nil {
		return diag.FromErr(err)
	}

	if resp.Snapshots == nil || len(*resp.Snapshots) < 1 {
		return diag.FromErr(ErrNoResults)
	}

	snapshots := make([]map[string]interface{}, len(*resp.Snapshots))
	for k, v := range *resp.Snapshots {
		snapshot := make(map[string]interface{})

		snapshot["description"] = v.Description
		snapshot["account_alias"] = v.AccountAlias
		snapshot["account_id"] = v.AccountId
		snapshot["creation_date"] = from.ISO8601(v.CreationDate)
		snapshot["progress"] = v.Progress
		snapshot["snapshot_id"] = v.SnapshotId
		snapshot["state"] = v.State
		snapshot["volume_id"] = v.VolumeId
		snapshot["volume_size"] = v.VolumeSize
		snapshot["tags"] = FlattenOAPITagsSDK(ptr.From(v.Tags))

		lp := make([]map[string]interface{}, 1)
		lp[0] = make(map[string]interface{})
		lp[0]["global_permission"] = ptr.From(v.PermissionsToCreateVolume).GlobalPermission
		lp[0]["account_ids"] = ptr.From(v.PermissionsToCreateVolume).AccountIds

		snapshot["permissions_to_create_volume"] = lp

		snapshots[k] = snapshot
	}

	d.SetId(id.UniqueId())
	// Single Snapshot found so set to state
	return diag.FromErr(d.Set("snapshots", snapshots))
}
