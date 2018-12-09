package outscale

import (
	"fmt"
	"log"
	"strings"
	"time"

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
			"volume_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"account_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"progress": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"account_alias": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"snapshot_id": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"volume_size": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"permissions_to_create_volume": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"account_ids": &schema.Schema{
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"global_permission": &schema.Schema{
							Type:     schema.TypeBool,
							Optional: true,
						},
					},
				},
			},
			"tags": tagsSchema(),
			"request_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceOutscaleOAPISnapshotCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OAPI

	v, ok := d.GetOk("volume_id")
	de, dok := d.GetOk("description")

	if !ok {
		return fmt.Errorf("please provide the volume_id required attribute")
	}

	request := oapi.CreateSnapshotRequest{
		VolumeId: v.(string),
	}

	if dok {
		request.Description = de.(string)
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

	d.SetId(res.OK.Snapshot.SnapshotId)
	d.Set("snapshot_id", res.OK.Snapshot.SnapshotId)

	/*
		if d.IsNewResource() {
			if err := setTags(conn, d); err != nil {
				return err
			}
			d.SetPartial("tag")
		}
	*/

	err = resourceOutscaleOAPISnapshotWaitForAvailable(d.Id(), conn)
	if err != nil {
		return err
	}

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

	if ec2err, ok := err.(awserr.Error); ok && ec2err.Code() == "InvalidSnapshotID.NotFound" {
		log.Printf("Snapshot %q Not found - removing from state", d.Id())
		d.SetId("")
		return nil
	}

	snapshot := res.OK.Snapshots[0]

	d.Set("description", snapshot.Description)
	d.Set("account_id", snapshot.AccountId)
	d.Set("progress", snapshot.Progress)
	d.Set("snapshot_id", snapshot.SnapshotId)
	d.Set("account_alias", snapshot.AccountAlias)
	d.Set("volume_id", snapshot.VolumeId)
	d.Set("state", snapshot.State)
	d.Set("volume_id", snapshot.VolumeId)
	d.Set("volume_size", snapshot.VolumeSize)
	d.Set("volume_size", snapshot.VolumeSize)
	d.Set("tags", tagsOAPIToMap(snapshot.Tags))
	d.Set("request_id", res.OK.ResponseContext.RequestId)

	permsMap := make(map[string]interface{})
	permsMap["account_ids"] = snapshot.PermissionsToCreateVolume.AccountIds
	permsMap["global_permission"] = snapshot.PermissionsToCreateVolume.GlobalPermission

	d.Set("permissions_to_create_volume", permsMap)

	return nil
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

func resourceOutscaleOAPISnapshotWaitForAvailable(id string, conn *oapi.Client) error {
	log.Printf("Waiting for Snapshot %s to become available...", id)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"pending", "pending/queued", "queued"},
		Target:     []string{"completed"},
		Refresh:    SnapshotOAPIStateRefreshFunc(conn, id),
		Timeout:    OutscaleImageRetryTimeout,
		Delay:      OutscaleImageRetryDelay,
		MinTimeout: OutscaleImageRetryMinTimeout,
	}

	_, err := stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Error waiting for Snapshot (%s) to be ready: %s", id, err)
	}
	return nil
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
