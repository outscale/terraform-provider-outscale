package outscale

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/antihax/optional"
	oscgo "github.com/marinsalinas/osc-sdk-go"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceOutscaleOAPISnapshot() *schema.Resource {
	return &schema.Resource{
		Create: resourceOutscaleOAPISnapshotCreate,
		Read:   resourceOutscaleOAPISnapshotRead,
		Update: resourceOutscaleOAPISnapshotUpdate,
		Delete: resourceOutscaleOAPISnapshotDelete,

		Schema: map[string]*schema.Schema{
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"snapshot_size": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"file_location": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"source_region_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"source_snapshot_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
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
						"global_permission": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"account_id": &schema.Schema{
							Type:     schema.TypeString,
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
			"tags": tagsListOAPISchema(),
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
	conn := meta.(*OutscaleClient).OSCAPI

	v, ok := d.GetOk("volume_id")
	snp, sok := d.GetOk("snapshot_size")
	source, sourceok := d.GetOk("source_snapshot_id")

	if !ok && !sok && !sourceok {
		return fmt.Errorf("please provide the source_snapshot_id, volume_id or snapshot_size argument")
	}

	description := d.Get("description").(string)
	fileLocation := d.Get("file_location").(string)
	sourceRegionName := d.Get("source_region_name").(string)

	request := oscgo.CreateSnapshotRequest{
		Description:  &description,
		FileLocation: &fileLocation,
	}

	if ok {
		request.SetVolumeId(v.(string))
	}

	if sok && snp.(int) > 0 {
		log.Printf("[DEBUG] Snapshot Size %d", snp.(int))

		request.SetSnapshotSize(int64(snp.(int)))
	}

	if sourceok {
		request.SetSourceSnapshotId(source.(string))
	}
	if sourceRegionName != "" {
		request.SetSourceRegionName(sourceRegionName)
	}

	var resp oscgo.CreateSnapshotResponse
	var err error
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, _, err = conn.SnapshotApi.CreateSnapshot(context.Background(), &oscgo.CreateSnapshotOpts{CreateSnapshotRequest: optional.NewInterface(request)})
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

	log.Printf("Waiting for Snapshot %s to become available...", resp.Snapshot.GetSnapshotId())

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"pending", "pending/queued", "queued", "importing"},
		Target:     []string{"completed"},
		Refresh:    SnapshotOAPIStateRefreshFunc(conn, resp.Snapshot.GetSnapshotId()),
		Timeout:    OutscaleImageRetryTimeout,
		Delay:      OutscaleImageRetryDelay,
		MinTimeout: OutscaleImageRetryMinTimeout,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Error waiting for Snapshot (%s) to be ready: %s", resp.Snapshot.GetSnapshotId(), err)
	}

	if tags, ok := d.GetOk("tags"); ok {
               err := assignTags(tags.(*schema.Set), resp.Snapshot.GetSnapshotId(), conn)
		if err != nil {
			return err
		}
	}

	d.SetId(resp.Snapshot.GetSnapshotId())

	return resourceOutscaleOAPISnapshotRead(d, meta)
}

func resourceOutscaleOAPISnapshotRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	req := oscgo.ReadSnapshotsRequest{
		Filters: &oscgo.FiltersSnapshot{SnapshotIds: &[]string{d.Id()}},
	}

	var resp oscgo.ReadSnapshotsResponse
	var err error
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, _, err = conn.SnapshotApi.ReadSnapshots(context.Background(), &oscgo.ReadSnapshotsOpts{ReadSnapshotsRequest: optional.NewInterface(req)})
		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("Error reading the snapshot %snapshot", err)
	}

	snapshot := resp.GetSnapshots()[0]

	return resourceDataAttrSetter(d, func(set AttributeSetter) error {
		permisions := snapshot.GetPermissionsToCreateVolume()
		if err := set("description", snapshot.GetDescription()); err != nil {
			return err
		}
		if err := set("volume_id", snapshot.GetVolumeId()); err != nil {
			return err
		}
		if err := set("account_alias", snapshot.GetAccountAlias()); err != nil {
			return err
		}
		if err := set("account_id", snapshot.GetAccountId()); err != nil {
			return err
		}
		if err := set("permissions_to_create_volume", omiOAPIPermissionToLuch(&permisions)); err != nil {
			return err
		}
		if err := set("progress", snapshot.GetProgress()); err != nil {
			return err
		}
		if err := set("snapshot_id", snapshot.GetSnapshotId()); err != nil {
			return err
		}
		if err := set("state", snapshot.GetState()); err != nil {
			return err
		}
		if err := set("tags", tagsOSCAPIToMap(snapshot.GetTags())); err != nil {
			return err
		}
		if err := set("volume_size", snapshot.GetVolumeSize()); err != nil {
			return err
		}
		return set("request_id", resp.ResponseContext.GetRequestId())
	})
}

func resourceOutscaleOAPISnapshotUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	d.Partial(true)

	if err := setOSCAPITags(conn, d); err != nil {
		return err
	}

	d.SetPartial("tags")

	d.Partial(false)
	return resourceOutscaleOAPISnapshotRead(d, meta)
}

func resourceOutscaleOAPISnapshotDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		request := oscgo.DeleteSnapshotRequest{SnapshotId: d.Id()}

		err := resource.Retry(5*time.Minute, func() *resource.RetryError {
			_, _, err := conn.SnapshotApi.DeleteSnapshot(context.Background(), &oscgo.DeleteSnapshotOpts{DeleteSnapshotRequest: optional.NewInterface(request)})

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
func SnapshotOAPIStateRefreshFunc(client *oscgo.APIClient, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		emptyResp := oscgo.ReadSnapshotsResponse{}

		var resp oscgo.ReadSnapshotsResponse
		var err error

		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			resp, _, err = client.SnapshotApi.ReadSnapshots(context.Background(), &oscgo.ReadSnapshotsOpts{ReadSnapshotsRequest: optional.NewInterface(oscgo.ReadSnapshotsRequest{
				Filters: &oscgo.FiltersSnapshot{SnapshotIds: &[]string{id}},
			})})
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

			} else if len(resp.GetSnapshots()) == 0 {
				log.Printf("[INFO] OMI %s state %s", id, "destroyed")
				return emptyResp, "destroyed", nil
			} else {
				return emptyResp, "", fmt.Errorf("Error on refresh: %+v", err)
			}
		}

		if resp.GetSnapshots() == nil || len(resp.GetSnapshots()) == 0 {
			return emptyResp, "destroyed", nil
		}

		// OMI is valid, so return it's state
		return resp.GetSnapshots()[0], resp.GetSnapshots()[0].GetState(), nil
	}
}
