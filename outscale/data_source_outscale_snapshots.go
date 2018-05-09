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

func dataSourceOutscaleSnapshots() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOutscaleSnapshotsRead,

		Schema: map[string]*schema.Schema{
			//selection criteria
			"filter": dataSourceFiltersSchema(),
			"owner": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"snapshot_id": {
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
			"snapshot_set": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"progress": {
							Type:     schema.TypeString,
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
						"volume_size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"tag_set": tagsSchemaComputed(),
					},
				},
			},
		},
	}
}

func dataSourceOutscaleSnapshotsRead(d *schema.ResourceData, meta interface{}) error {
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
		params.SnapshotIds = expandStringList(snapshotIds.([]interface{}))
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

	if len(resp.Snapshots) < 1 {
		return fmt.Errorf("Your query returned no results. Please change your search criteria and try again.")
	}

	snapshots := make([]map[string]interface{}, len(resp.Snapshots))
	for k, v := range resp.Snapshots {
		snapshot := make(map[string]interface{})

		snapshot["description"] = aws.StringValue(v.Description)
		snapshot["owner_alias"] = aws.StringValue(v.OwnerAlias)
		snapshot["owner_id"] = aws.StringValue(v.OwnerId)
		snapshot["progress"] = aws.StringValue(v.Progress)
		snapshot["snapshot_id"] = aws.StringValue(v.SnapshotId)
		snapshot["status"] = aws.StringValue(v.State)
		snapshot["status_message"] = aws.StringValue(v.StateMessage)
		snapshot["volume_id"] = aws.StringValue(v.VolumeId)
		snapshot["volume_size"] = aws.Int64Value(v.VolumeSize)
		snapshot["tag_set"] = tagsToMap(v.Tags)

		snapshots[k] = snapshot
	}

	d.SetId(resource.UniqueId())
	//Single Snapshot found so set to state
	return d.Set("snapshot_set", snapshots)
}
