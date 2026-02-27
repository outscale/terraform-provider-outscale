package oapi

import (
	"context"
	"strconv"
	"time"

	"github.com/outscale/goutils/sdk/ptr"
	"github.com/outscale/osc-sdk-go/v3/pkg/options"
	"github.com/outscale/osc-sdk-go/v3/pkg/osc"
	"github.com/samber/lo"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/framework/fwhelpers/from"
	"github.com/outscale/terraform-provider-outscale/internal/framework/fwhelpers/to"
	"github.com/outscale/terraform-provider-outscale/internal/utils"
)

func DataSourceOutscaleSnapshot() *schema.Resource {
	return &schema.Resource{
		ReadContext: DataSourceOutscaleSnapshotRead,

		Schema: map[string]*schema.Schema{
			// selection criteria
			"filter": dataSourceFiltersSchema(),
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

			// Computed values returned
			"progress": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"creation_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"snapshot_id": {
				Type:     schema.TypeString,
				Optional: true,
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
				Optional: true,
			},
			"account_alias": {
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
			"tags": TagsSchemaComputedSDK(),
			"request_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func DataSourceOutscaleSnapshotRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	var err error
	filter := osc.FiltersSnapshot{}
	if restorableUsersOk {
		filter.PermissionsToCreateVolumeAccountIds = utils.InterfaceSliceToStringSlicePtr(restorableUsers.([]interface{}))
		params.Filters = &filter
	}
	if filtersOk {
		params.Filters, err = buildOutscaleOapiSnapshootDataSourceFilters(filters.(*schema.Set))
		if err != nil {
			return diag.FromErr(err)
		}
	}
	if ownersOk {
		params.Filters.AccountIds = &[]string{owners.(string)}
	}
	if snapshotIdsOk {
		params.Filters.SnapshotIds = &[]string{snapshotIds.(string)}
	}

	resp, err := client.ReadSnapshots(ctx, params, options.WithRetryTimeout(5*time.Minute))
	if err != nil {
		return diag.FromErr(err)
	}

	if resp.Snapshots == nil || len(*resp.Snapshots) < 1 {
		return diag.FromErr(ErrNoResults)
	}
	if len(*resp.Snapshots) > 1 {
		return diag.FromErr(ErrMultipleResults)
	}

	snapshot := (*resp.Snapshots)[0]

	// Single Snapshot found so set to state
	return diag.FromErr(snapshotOAPIDescriptionAttributes(d, &snapshot))
}

func snapshotOAPIDescriptionAttributes(d *schema.ResourceData, snapshot *osc.Snapshot) error {
	d.SetId(snapshot.SnapshotId)
	if err := d.Set("description", ptr.From(snapshot.Description)); err != nil {
		return err
	}
	if err := d.Set("account_alias", ptr.From(snapshot.AccountAlias)); err != nil {
		return err
	}
	if err := d.Set("account_id", snapshot.AccountId); err != nil {
		return err
	}
	if err := d.Set("creation_date", from.ISO8601(snapshot.CreationDate)); err != nil {
		return err
	}
	if err := d.Set("progress", ptr.From(snapshot.Progress)); err != nil {
		return err
	}
	if err := d.Set("snapshot_id", snapshot.SnapshotId); err != nil {
		return err
	}
	if err := d.Set("state", snapshot.State); err != nil {
		return err
	}
	if err := d.Set("volume_id", snapshot.VolumeId); err != nil {
		return err
	}
	if err := d.Set("volume_size", snapshot.VolumeSize); err != nil {
		return err
	}

	lp := make([]map[string]interface{}, 1)
	lp[0] = make(map[string]interface{})
	perm := ptr.From(snapshot.PermissionsToCreateVolume)
	lp[0]["global_permission"] = perm.GlobalPermission
	lp[0]["account_ids"] = perm.AccountIds

	if err := d.Set("permissions_to_create_volume", lp); err != nil {
		return err
	}

	return d.Set("tags", FlattenOAPITagsSDK(ptr.From(snapshot.Tags)))
}

func buildOutscaleOapiSnapshootDataSourceFilters(set *schema.Set) (*osc.FiltersSnapshot, error) {
	var filter osc.FiltersSnapshot
	for _, v := range set.List() {
		m := v.(map[string]interface{})
		var values []string

		for _, e := range m["values"].([]interface{}) {
			values = append(values, e.(string))
		}

		switch name := m["name"].(string); name {
		case "account_aliases":
			filter.AccountAliases = &values

		case "account_ids":
			filter.AccountIds = &values

		case "descriptions":
			filter.Descriptions = &values
		case "to_creation_date":
			if values[0] != "" {
				valDate, err := to.ISO8601ToDate(values[0])
				if err != nil {
					return nil, err
				}
				filter.ToCreationDate = &valDate
			}
		case "from_creation_date":
			if values[0] != "" {
				valDate, err := to.ISO8601FromDate(values[0])
				if err != nil {
					return nil, err
				}
				filter.FromCreationDate = &valDate
			}
		case "permissions_to_create_volume_account_ids":
			filter.PermissionsToCreateVolumeAccountIds = &values

		case "permissions_to_create_volume_global_permission":
			boolean, err := strconv.ParseBool(values[0])
			if err != nil {
				return nil, err
			}
			filter.PermissionsToCreateVolumeGlobalPermission = &boolean

		case "progresses":
			filter.Progresses = new(utils.StringSliceToIntSlice(values))

		case "snapshot_ids":
			filter.SnapshotIds = &values

		case "states":
			filter.States = new(lo.Map(values, func(s string, _ int) osc.SnapshotState { return osc.SnapshotState(s) }))

		case "tag_keys":
			filter.TagKeys = &values

		case "tag_values":
			filter.TagValues = &values

		case "tags":
			filter.Tags = &values

		case "volume_ids":
			filter.VolumeIds = &values

		case "volume_sizes":
			filter.VolumeSizes = new(utils.StringSliceToIntSlice(values))

		default:
			return nil, utils.UnknownDataSourceFilterError(name)
		}
	}
	return &filter, nil
}
