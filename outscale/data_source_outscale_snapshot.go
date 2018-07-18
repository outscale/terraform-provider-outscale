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

func dataSourceOutscaleSnapshot() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOutscaleSnapshotRead,

		Schema: map[string]*schema.Schema{
			//selection criteria
			"filter": dataSourceFiltersSchema(),
			"owner": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"restorable_by": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			//Computed values returned
			"progress": {
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
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status_message": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"owner_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"owner_alias": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"start_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"request_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"volume_size": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"tag_set": tagsSchemaComputed(),
		},
	}
}

func dataSourceOutscaleSnapshotRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	restorableUsers, restorableUsersOk := d.GetOk("restorable_by")
	filters, filtersOk := d.GetOk("filter")
	snapshotIds, snapshotIdsOk := d.GetOk("snapshot_id")
	owners, ownersOk := d.GetOk("owner")

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
		return fmt.Errorf("our query returned no results, please change your search criteria and try again")
	}
	if len(resp.Snapshots) > 1 {
		return fmt.Errorf("our query returned more than one result, please try a more specific search criteria")
	}

	snapshot = resp.Snapshots[0]
	d.Set("request_id", resp.RequestId)

	//Single Snapshot found so set to state
	return snapshotDescriptionAttributes(d, snapshot)
}

func snapshotDescriptionAttributes(d *schema.ResourceData, snapshot *fcu.Snapshot) error {
	d.SetId(*snapshot.SnapshotId)
	d.Set("description", snapshot.Description)
	d.Set("owner_alias", snapshot.OwnerAlias)
	d.Set("owner_id", snapshot.OwnerId)
	d.Set("progress", snapshot.Progress)
	d.Set("snapshot_id", snapshot.SnapshotId)
	d.Set("status", snapshot.State)
	d.Set("status_message", snapshot.StateMessage)
	d.Set("volume_id", snapshot.VolumeId)
	d.Set("volume_size", snapshot.VolumeSize)
	d.Set("start_time", snapshot.StartTime.String())

	return d.Set("tag_set", tagsToMap(snapshot.Tags))
}
