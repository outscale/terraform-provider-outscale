package outscale

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/outscale/terraform-provider-outscale/utils"
)

func DataSourceOutscaleOAPISnapshot() *schema.Resource {
	return &schema.Resource{
		Read: DataSourceOutscaleOAPISnapshotRead,

		Schema: map[string]*schema.Schema{
			//selection criteria
			"filter": DataSourceFiltersSchema(),
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

			//Computed values returned
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
			"tags": DataSourceTagsSchema(),
			"request_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func DataSourceOutscaleOAPISnapshotRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	restorableUsers, restorableUsersOk := d.GetOk("permission_to_create_volume")
	filters, filtersOk := d.GetOk("filter")
	snapshotIds, snapshotIdsOk := d.GetOk("snapshot_id")
	owners, ownersOk := d.GetOk("account_id")

	if restorableUsers == false && !filtersOk && snapshotIds == false && !ownersOk {
		return fmt.Errorf("One of snapshot_ids, filters, restorable_by_user_ids, or owners must be assigned")
	}

	params := oscgo.ReadSnapshotsRequest{
		Filters: &oscgo.FiltersSnapshot{},
	}

	filter := oscgo.FiltersSnapshot{}
	if restorableUsersOk {
		filter.SetPermissionsToCreateVolumeAccountIds(utils.InterfaceSliceToStringSlice(restorableUsers.([]interface{})))
		params.SetFilters(filter)
	}
	if filtersOk {
		buildOutscaleOapiSnapshootDataSourceFilters(filters.(*schema.Set), params.Filters)
	}
	if ownersOk {
		filter.SetAccountIds([]string{owners.(string)})
		params.SetFilters(filter)
	}
	if snapshotIdsOk {
		filter.SetSnapshotIds([]string{snapshotIds.(string)})
		params.SetFilters(filter)
	}

	var resp oscgo.ReadSnapshotsResponse
	var err error
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		rp, httpResp, err := conn.SnapshotApi.ReadSnapshots(context.Background()).ReadSnapshotsRequest(params).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		resp = rp
		return nil
	})
	if err != nil {
		return err
	}

	if len(resp.GetSnapshots()) < 1 {
		return fmt.Errorf("your query returned no results, please change your search criteria and try again")
	}
	if len(resp.GetSnapshots()) > 1 {
		return fmt.Errorf("your query returned more than one result, please try a more specific search criteria")
	}

	snapshot := resp.GetSnapshots()[0]

	//Single Snapshot found so set to state
	return snapshotOAPIDescriptionAttributes(d, &snapshot)
}

func snapshotOAPIDescriptionAttributes(d *schema.ResourceData, snapshot *oscgo.Snapshot) error {
	d.SetId(snapshot.GetSnapshotId())
	if err := d.Set("description", snapshot.GetDescription()); err != nil {
		return err
	}
	if err := d.Set("account_alias", snapshot.GetAccountAlias()); err != nil {
		return err
	}
	if err := d.Set("account_id", snapshot.GetAccountId()); err != nil {
		return err
	}
	if err := d.Set("creation_date", snapshot.GetCreationDate()); err != nil {
		return err
	}
	if err := d.Set("progress", snapshot.GetProgress()); err != nil {
		return err
	}
	if err := d.Set("snapshot_id", snapshot.GetSnapshotId()); err != nil {
		return err
	}
	if err := d.Set("state", snapshot.GetState()); err != nil {
		return err
	}
	if err := d.Set("volume_id", snapshot.GetVolumeId()); err != nil {
		return err
	}
	if err := d.Set("volume_size", snapshot.GetVolumeSize()); err != nil {
		return err
	}

	lp := make([]map[string]interface{}, 1)
	lp[0] = make(map[string]interface{})
	lp[0]["global_permission"] = snapshot.PermissionsToCreateVolume.GetGlobalPermission()
	lp[0]["account_ids"] = snapshot.PermissionsToCreateVolume.GetAccountIds()

	if err := d.Set("permissions_to_create_volume", lp); err != nil {
		return err
	}

	return d.Set("tags", tagsOSCAPIToMap(snapshot.GetTags()))
}

func buildOutscaleOapiSnapshootDataSourceFilters(set *schema.Set, filter *oscgo.FiltersSnapshot) *oscgo.FiltersSnapshot {

	for _, v := range set.List() {
		m := v.(map[string]interface{})
		var values []string

		for _, e := range m["values"].([]interface{}) {
			values = append(values, e.(string))
		}

		switch name := m["name"].(string); name {
		case "account_aliases":
			filter.SetAccountAliases(values)

		case "account_ids":
			filter.SetAccountIds(values)

		case "descriptions":
			filter.SetDescriptions(values)

		case "permissions_to_create_volume_account_ids":
			filter.SetPermissionsToCreateVolumeAccountIds(values)

		case "permissions_to_create_volume_global_permission":
			boolean, _ := strconv.ParseBool(values[0])
			filter.SetPermissionsToCreateVolumeGlobalPermission(boolean)

		case "progresses":
			filter.SetProgresses(utils.StringSliceToInt32Slice(values))

		case "snapshot_ids":
			filter.SetSnapshotIds(values)

		case "states":
			filter.SetStates(values)

		case "tag_keys":
			filter.SetTagKeys(values)

		case "tag_values":
			filter.SetTagValues(values)

		case "tags":
			filter.SetTags(values)

		case "volume_ids":
			filter.SetVolumeIds(values)

		case "volume_sizes":
			filter.SetVolumeSizes(utils.StringSliceToInt32Slice(values))

		default:
			log.Printf("[Debug] Unknown Filter Name: %s.", name)
		}
	}
	return filter
}
