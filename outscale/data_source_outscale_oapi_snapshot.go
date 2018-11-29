package outscale

import (
	"fmt"
	"log"
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
			"permission_to_create_volume": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			//Computed values returned
			"completion": {
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
			"comment": {
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
			"tag": dataSourceTagsSchema(),
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
		params.Filters.AccountIds = oapiExpandStringList(owners.([]interface{}))
	}
	if snapshotIdsOk {
		params.Filters.SnapshotIds = oapiExpandStringList(snapshotIds.([]interface{}))
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
	//d.Set("account_alias", snapshot.OwnerAlias)
	d.Set("account_id", snapshot.AccountId)
	d.Set("completion", snapshot.Progress)
	d.Set("snapshot_id", snapshot.SnapshotId)
	d.Set("state", snapshot.State)
	//d.Set("comment", snapshot.StateMessage)
	d.Set("volume_id", snapshot.VolumeId)
	d.Set("volume_size", snapshot.VolumeSize)

	return setSnapshotArgTags("tag", d, &snapshot.Tags)
}

func buildOutscaleOapiSnapshootDataSourceFilters(set *schema.Set, filter *oapi.FiltersSnapshot) *oapi.FiltersSnapshot {

	for _, v := range set.List() {
		m := v.(map[string]interface{})
		var values []string

		for _, e := range m["values"].([]interface{}) {
			values = append(values, e.(string))
		}

		switch name := m["name"].(string); name {
		case "description":
			filter.Descriptions = values

		case "owner-alias":
			filter.AccountAliases = values

		case "owner-id":
			filter.AccountIds = values

		case "progress":
			filter.Progresses = utils.StringSliceToInt64Slice(values)

		case "snapshot-id":
			filter.SnapshotIds = values

		case "status":
			filter.States = values

		case "volume-id":
			filter.VolumeIds = values

		case "volume-size":
			filter.VolumeSizes = utils.StringSliceToInt64Slice(values)

		case "tag":
			filter.Tags = values

		case "tag-key":
			filter.TagKeys = values

		case "tag-value":
			filter.TagValues = values

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

func setSnapshotArgTags(attrName string, d *schema.ResourceData, tags *[]oapi.ResourceTag) error {
	if *tags != nil {
		if err := d.Set(attrName, tagsOAPIToMap(*tags)); err != nil {
			return err
		}
	} else {
		if err := d.Set(attrName, []map[string]string{
			map[string]string{
				"key":   "",
				"value": "",
			},
		}); err != nil {
			return err
		}
	}
	return nil
}
