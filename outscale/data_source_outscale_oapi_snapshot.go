package outscale

import (
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-outscale/osc/fcu"
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
	conn := meta.(*OutscaleClient).FCU

	restorableUsers, restorableUsersOk := d.GetOk("permission_to_create_volume")
	filters, filtersOk := d.GetOk("filter")
	snapshotIds, snapshotIdsOk := d.GetOk("snapshot_id")
	owners, ownersOk := d.GetOk("account_id")

	if restorableUsers == false && filtersOk == false && snapshotIds == false && ownersOk == false {
		return fmt.Errorf("One of snapshot_ids, filters, restorable_by_user_ids, or owners must be assigned")
	}

	params := &fcu.DescribeSnapshotsInput{}
	if restorableUsersOk {
		params.RestorableByUserIds = expandStringList(restorableUsers.([]interface{}))
	}
	if filtersOk {
		params.Filters = buildOutscaleDataSourceFilters(filters.(*schema.Set))
	}
	if ownersOk {
		params.OwnerIds = expandStringList(owners.([]interface{}))
	}
	if snapshotIdsOk {
		params.SnapshotIds = []*string{aws.String(snapshotIds.(string))}
	}

	var resp *fcu.DescribeSnapshotsOutput
	var err error
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, err = conn.VM.DescribeSnapshots(params)

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

	var snapshot *fcu.Snapshot
	if len(resp.Snapshots) < 1 {
		return fmt.Errorf("your query returned no results, please change your search criteria and try again")
	}
	if len(resp.Snapshots) > 1 {
		return fmt.Errorf("your query returned more than one result, please try a more specific search criteria")
	}

	snapshot = resp.Snapshots[0]

	//Single Snapshot found so set to state
	return snapshotOAPIDescriptionAttributes(d, snapshot)
}

func snapshotOAPIDescriptionAttributes(d *schema.ResourceData, snapshot *fcu.Snapshot) error {
	d.SetId(*snapshot.SnapshotId)
	d.Set("description", snapshot.Description)
	d.Set("account_alias", snapshot.OwnerAlias)
	d.Set("account_id", snapshot.OwnerId)
	d.Set("completion", snapshot.Progress)
	d.Set("snapshot_id", snapshot.SnapshotId)
	d.Set("state", snapshot.State)
	d.Set("comment", snapshot.StateMessage)
	d.Set("volume_id", snapshot.VolumeId)
	d.Set("volume_size", snapshot.VolumeSize)

	return d.Set("tag", tagsToMap(snapshot.Tags))
}
