package outscale

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-outscale/osc/oapi"
)

func dataSourceOutscaleOAPISnapshots() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOutscaleOAPISnapshotsRead,

		Schema: map[string]*schema.Schema{
			//selection criteria
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
			//Computed values returned
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
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"volume_size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
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
						"tags": tagsSchemaComputed(),
					},
				},
			},
		},
	}
}

func dataSourceOutscaleOAPISnapshotsRead(d *schema.ResourceData, meta interface{}) error {
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

	if len(resp.OK.Snapshots) < 1 {
		return fmt.Errorf("your query returned no results, please change your search criteria and try again")
	}

	snapshots := make([]map[string]interface{}, len(resp.OK.Snapshots))
	for k, v := range resp.OK.Snapshots {
		snapshot := make(map[string]interface{})

		snapshot["description"] = v.Description
		snapshot["account_alias"] = v.AccountAlias
		snapshot["account_id"] = v.AccountId
		snapshot["progress"] = v.Progress
		snapshot["snapshot_id"] = v.SnapshotId
		snapshot["state"] = v.State
		snapshot["volume_id"] = v.VolumeId
		snapshot["volume_size"] = v.VolumeSize
		snapshot["tags"] = tagsOAPIToMap(v.Tags)

		lp := make([]map[string]interface{}, 1)
		lp[0] = make(map[string]interface{})
		lp[0]["global_permission"] = v.PermissionsToCreateVolume.GlobalPermission
		lp[0]["account_ids"] = v.PermissionsToCreateVolume.AccountIds

		snapshot["permissions_to_create_volume"] = lp

		snapshots[k] = snapshot
	}

	d.SetId(resource.UniqueId())
	//Single Snapshot found so set to state
	return d.Set("snapshots", snapshots)
}
