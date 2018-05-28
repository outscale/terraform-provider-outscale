package outscale

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-outscale/osc/fcu"
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
			"completion": {
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
			"tag": tagsSchema(),
		},
	}
}

func resourceOutscaleOAPISnapshotCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	v, ok := d.GetOk("volume_id")
	de, dok := d.GetOk("description")

	if !ok {
		return fmt.Errorf("please provide the volume_id required attribute")
	}

	request := &fcu.CreateSnapshotInput{
		VolumeId: aws.String(v.(string)),
	}

	if dok {
		request.Description = aws.String(de.(string))
	}

	var res *fcu.Snapshot
	var err error
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		res, err = conn.VM.CreateSnapshot(request)

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

	d.SetId(*res.SnapshotId)
	d.Set("snapshot_id", res.SnapshotId)

	if d.IsNewResource() {
		if err := setTags(conn, d); err != nil {
			return err
		}
		d.SetPartial("tag")
	}

	err = resourceOutscaleOAPISnapshotWaitForAvailable(d.Id(), conn)
	if err != nil {
		return err
	}

	return resourceOutscaleOAPISnapshotRead(d, meta)
}

func resourceOutscaleOAPISnapshotRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	req := &fcu.DescribeSnapshotsInput{
		SnapshotIds: []*string{aws.String(d.Id())},
	}
	var res *fcu.DescribeSnapshotsOutput
	var err error
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		res, err = conn.VM.DescribeSnapshots(req)

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

	snapshot := res.Snapshots[0]

	d.Set("description", snapshot.Description)
	d.Set("account_id", snapshot.OwnerId)
	d.Set("completion", snapshot.Progress)
	d.Set("account_alias", snapshot.OwnerAlias)
	d.Set("volume_id", snapshot.VolumeId)
	d.Set("state", snapshot.State)
	d.Set("status_message", snapshot.StateMessage)
	d.Set("volume_size", snapshot.VolumeSize)
	d.Set("tag", tagsToMap(snapshot.Tags))

	return nil
}

func resourceOutscaleOAPISnapshotDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		request := &fcu.DeleteSnapshotInput{
			SnapshotId: aws.String(d.Id()),
		}
		var err error
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			_, err := conn.VM.DeleteSnapshot(request)

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

func resourceOutscaleOAPISnapshotWaitForAvailable(id string, conn *fcu.Client) error {
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
func SnapshotOAPIStateRefreshFunc(client *fcu.Client, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		emptyResp := &fcu.DescribeSnapshotsOutput{}

		var resp *fcu.DescribeSnapshotsOutput
		var err error

		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			resp, err = client.VM.DescribeSnapshots(&fcu.DescribeSnapshotsInput{
				SnapshotIds: []*string{aws.String(id)},
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

			} else if resp != nil && len(resp.Snapshots) == 0 {
				log.Printf("[INFO] OMI %s state %s", id, "destroyed")
				return emptyResp, "destroyed", nil
			} else {
				return emptyResp, "", fmt.Errorf("Error on refresh: %+v", err)
			}
		}

		if resp == nil || resp.Snapshots == nil || len(resp.Snapshots) == 0 {
			return emptyResp, "destroyed", nil
		}

		// OMI is valid, so return it's state
		return resp.Snapshots[0], *resp.Snapshots[0].State, nil
	}
}
