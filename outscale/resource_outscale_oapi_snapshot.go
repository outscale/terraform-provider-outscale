package outscale

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/spf13/cast"

	"github.com/terraform-providers/terraform-provider-outscale/osc/oapi"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceOutscaleOAPISnapshot() *schema.Resource {
	return &schema.Resource{
		Create: resourceOutscaleOAPISnapshotCreate,
		Read:   resourceOutscaleOAPISnapshotRead,
		Delete: resourceOutscaleOAPISnapshotDelete,

		Schema: map[string]*schema.Schema{
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"file_location": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"snapshot_size": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"source_region_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"source_snapshot_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"volume_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"account_alias": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"account_id": {
				Type:     schema.TypeString,
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
			"progress": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"snapshot_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": tagsOAPIListSchemaComputed(),
			"volume_size": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"request_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceOutscaleOAPISnapshotCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OAPI

	volumeID, ok := d.GetOk("volume_id")
	if !ok {
		return fmt.Errorf("please provide the volume_id required attribute")
	}

	request := oapi.CreateSnapshotRequest{
		Description:      d.Get("description").(string),
		FileLocation:     d.Get("file_location").(string),
		SnapshotSize:     cast.ToInt64(d.Get("snapshot_size")),
		SourceRegionName: d.Get("source_region_name").(string),
		SourceSnapshotId: d.Get("source_snapshot_id").(string),
		VolumeId:         volumeID.(string),
	}

	var res *oapi.POST_CreateSnapshotResponses
	var err error
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		res, err = conn.POST_CreateSnapshot(request)
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

	log.Printf("Waiting for Snapshot %s to become available...", res.OK.Snapshot.SnapshotId)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"pending", "pending/queued", "queued"},
		Target:     []string{"completed"},
		Refresh:    SnapshotOAPIStateRefreshFunc(conn, res.OK.Snapshot.SnapshotId),
		Timeout:    OutscaleImageRetryTimeout,
		Delay:      OutscaleImageRetryDelay,
		MinTimeout: OutscaleImageRetryMinTimeout,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Error waiting for Snapshot (%s) to be ready: %s", res.OK.Snapshot.SnapshotId, err)
	}

	d.SetId(res.OK.Snapshot.SnapshotId)

	return resourceOutscaleOAPISnapshotRead(d, meta)
}

func resourceOutscaleOAPISnapshotRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OAPI

	req := oapi.ReadSnapshotsRequest{
		Filters: oapi.FiltersSnapshot{SnapshotIds: []string{d.Id()}},
	}

	var res *oapi.POST_ReadSnapshotsResponses
	var err error
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		res, err = conn.POST_ReadSnapshots(req)
		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("Error reading the snapshot %s", err)
	}

	s := res.OK.Snapshots[0]

	return resourceDataAttrSetter(d, func(set AttributeSetter) error {

		set("description", s.Description)
		set("volume_id", s.VolumeId)
		set("account_alias", s.AccountAlias)
		set("account_id", s.AccountId)
		set("permissions_to_create_volume", omiOAPIPermissionToLuch(s.PermissionsToCreateVolume))
		set("progress", s.Progress)
		set("snapshot_id", s.SnapshotId)
		set("state", s.State)
		set("tags", tagsOAPIToMap(s.Tags))
		set("volume_size", s.VolumeSize)
		return set("request_id", res.OK.ResponseContext.RequestId)
	})
}

func resourceOutscaleOAPISnapshotDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OAPI

	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		request := oapi.DeleteSnapshotRequest{
			SnapshotId: d.Id(),
		}
		var err error
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			_, err := conn.POST_DeleteSnapshot(request)

			if err != nil {
				if strings.Contains(err.Error(), "RequestLimitExceeded") {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}

			return nil
		})
		if err == nil {
			return nil
		}

		ebsErr, ok := err.(awserr.Error)
		if ebsErr.Code() == "SnapshotInUse" {
			return resource.RetryableError(fmt.Errorf("EBS SnapshotInUse - trying again while it detaches"))
		}

		if !ok {
			return resource.NonRetryableError(err)
		}

		return resource.NonRetryableError(err)
	})
}

// SnapshotOAPIStateRefreshFunc ...
func SnapshotOAPIStateRefreshFunc(client *oapi.Client, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		emptyResp := &oapi.ReadSnapshotsResponse{}

		var resp *oapi.POST_ReadSnapshotsResponses
		var err error

		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			resp, err = client.POST_ReadSnapshots(oapi.ReadSnapshotsRequest{
				Filters: oapi.FiltersSnapshot{SnapshotIds: []string{id}},
			})
			if err != nil {
				if strings.Contains(err.Error(), "RequestLimitExceeded:") {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})

		if err != nil {
			if e := fmt.Sprint(err); strings.Contains(e, "InvalidAMIID.NotFound") {
				log.Printf("[INFO] OMI %s state %s", id, "destroyed")
				return emptyResp, "destroyed", nil

			} else if resp != nil && len(resp.OK.Snapshots) == 0 {
				log.Printf("[INFO] OMI %s state %s", id, "destroyed")
				return emptyResp, "destroyed", nil
			} else {
				return emptyResp, "", fmt.Errorf("Error on refresh: %+v", err)
			}
		}

		if resp == nil || resp.OK.Snapshots == nil || len(resp.OK.Snapshots) == 0 {
			return emptyResp, "destroyed", nil
		}

		// OMI is valid, so return it's state
		return resp.OK.Snapshots[0], resp.OK.Snapshots[0].State, nil
	}
}
