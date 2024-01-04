package outscale

import (
	"context"
	"fmt"
	"time"

	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/terraform-providers/terraform-provider-outscale/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
						"tags": dataSourceTagsSchema(),
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

func dataSourceOutscaleOAPISnapshotsRead(d *schema.ResourceData, meta interface{}) error {
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
		if err := buildOutscaleOapiSnapshootDataSourceFilters(filters.(*schema.Set), params.Filters); err != nil {
			return err
		}
	}

	if ownersOk {
		filter.SetAccountIds(utils.InterfaceSliceToStringSlice(owners.([]interface{})))
		params.SetFilters(filter)
	}
	if snapshotIdsOk {
		filter.SetSnapshotIds(utils.InterfaceSliceToStringSlice(snapshotIds.([]interface{})))
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

	snapshots := make([]map[string]interface{}, len(resp.GetSnapshots()))
	for k, v := range resp.GetSnapshots() {
		snapshot := make(map[string]interface{})

		snapshot["description"] = v.GetDescription()
		snapshot["account_alias"] = v.GetAccountAlias()
		snapshot["account_id"] = v.GetAccountId()
		snapshot["creation_date"] = v.GetCreationDate()
		snapshot["progress"] = v.GetProgress()
		snapshot["snapshot_id"] = v.GetSnapshotId()
		snapshot["state"] = v.GetState()
		snapshot["volume_id"] = v.GetVolumeId()
		snapshot["volume_size"] = v.GetVolumeSize()
		snapshot["tags"] = tagsOSCAPIToMap(v.GetTags())

		lp := make([]map[string]interface{}, 1)
		lp[0] = make(map[string]interface{})
		lp[0]["global_permission"] = v.PermissionsToCreateVolume.GetGlobalPermission()
		lp[0]["account_ids"] = v.PermissionsToCreateVolume.GetAccountIds()

		snapshot["permissions_to_create_volume"] = lp

		snapshots[k] = snapshot
	}

	d.SetId(resource.UniqueId())
	//Single Snapshot found so set to state
	return d.Set("snapshots", snapshots)
}
