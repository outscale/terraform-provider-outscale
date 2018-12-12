package outscale

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/terraform-providers/terraform-provider-outscale/osc/oapi"
	"github.com/terraform-providers/terraform-provider-outscale/utils"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceOutscaleOAPISnapshot() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOutscaleOAPISnapshotRead,

		Schema: map[string]*schema.Schema{
			//selection criteria
			"filter": dataSourceFiltersSchema(),
			"permissions_to_create_volume": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"account_ids": &schema.Schema{
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"global_permission": &schema.Schema{
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},

			//Computed values returned
			"progress": {
				Type:     schema.TypeInt,
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
			"tags": dataSourceTagsSchema(),
		},
	}
}

func dataSourceOutscaleOAPISnapshotRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OAPI

	restorableUsers, restorableUsersOk := d.GetOk("permission_to_create_volume")
	filters, filtersOk := d.GetOk("filter")
	snapshotIds, snapshotIdsOk := d.GetOk("snapshot_id")
	owners, ownersOk := d.GetOk("account_id")

	if restorableUsers == false && filtersOk == false && snapshotIds == false && ownersOk == false {
		return fmt.Errorf("One of snapshot_ids, filters, restorable_by_user_ids, or owners must be assigned")
	}

	params := oapi.ReadSnapshotsRequest{}
	if restorableUsersOk {
		params.Filters.PermissionsToCreateVolumeAccountIds = oapiExpandStringList(restorableUsers.([]interface{}))
	}
	if filtersOk {
		buildOutscaleOapiSnapshootDataSourceFilters(filters.(*schema.Set), &params.Filters)
	}
	if ownersOk {
		params.Filters.AccountIds = []string{owners.(string)}
	}
	if snapshotIdsOk {
		params.Filters.SnapshotIds = []string{snapshotIds.(string)}
	}

	var resp *oapi.POST_ReadSnapshotsResponses
	var err error
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, err = conn.POST_ReadSnapshots(params)

		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}

		return nil
	})
	if err != nil {
		return err
	}

	var snapshot oapi.Snapshot
	if len(resp.OK.Snapshots) < 1 {
		return fmt.Errorf("your query returned no results, please change your search criteria and try again")
	}
	if len(resp.OK.Snapshots) > 1 {
		return fmt.Errorf("your query returned more than one result, please try a more specific search criteria")
	}

	snapshot = resp.OK.Snapshots[0]

	//Single Snapshot found so set to state
	return snapshotOAPIDescriptionAttributes(d, &snapshot)
}

func snapshotOAPIDescriptionAttributes(d *schema.ResourceData, snapshot *oapi.Snapshot) error {
	d.SetId(snapshot.SnapshotId)
	d.Set("description", snapshot.Description)
	d.Set("account_alias", snapshot.AccountAlias)
	d.Set("account_id", snapshot.AccountId)
	d.Set("progress", snapshot.Progress)
	d.Set("snapshot_id", snapshot.SnapshotId)
	d.Set("state", snapshot.State)
	d.Set("volume_id", snapshot.VolumeId)
	d.Set("volume_size", snapshot.VolumeSize)

	lp := make([]map[string]interface{}, 1)
	lp[0] = make(map[string]interface{})
	lp[0]["global_permission"] = snapshot.PermissionsToCreateVolume.GlobalPermission
	lp[0]["account_ids"] = snapshot.PermissionsToCreateVolume.AccountIds

	if err := d.Set("permissions_to_create_volume", lp); err != nil {
		return err
	}

	return d.Set("tags", tagsOAPIToMap(snapshot.Tags))
}

func buildOutscaleOapiSnapshootDataSourceFilters(set *schema.Set, filter *oapi.FiltersSnapshot) *oapi.FiltersSnapshot {

	for _, v := range set.List() {
		m := v.(map[string]interface{})
		var values []string

		for _, e := range m["values"].([]interface{}) {
			values = append(values, e.(string))
		}

		switch name := m["name"].(string); name {
		case "account_aliases":
			filter.AccountAliases = values

		case "account_ids":
			filter.AccountIds = values

		case "descriptions":
			filter.Descriptions = values

		case "permissions_to_create_volume_account_ids":
			filter.PermissionsToCreateVolumeAccountIds = values

		case "permissions_to_create_volume_global_permission":
			filter.PermissionsToCreateVolumeGlobalPermission, _ = strconv.ParseBool(values[0])

		case "progresses":
			filter.Progresses = utils.StringSliceToInt64Slice(values)

		case "snapshot_ids":
			filter.SnapshotIds = values

		case "states":
			filter.States = values

		case "tag_keys":
			filter.TagKeys = values

		case "tag_values":
			filter.TagValues = values

		case "tags":
			filter.Tags = values

		case "volume_ids":
			filter.VolumeIds = values

		case "volume_sizes":
			filter.VolumeSizes = utils.StringSliceToInt64Slice(values)

		default:
			log.Printf("[Debug] Unknown Filter Name: %s.", name)
		}
	}
	return filter
}

func oapiExpandStringList(configured []interface{}) []string {
	vs := make([]string, 0, len(configured))
	for _, v := range configured {
		val, ok := v.(string)
		if ok && val != "" {
			vs = append(vs, v.(string))
		}
	}
	return vs
}
