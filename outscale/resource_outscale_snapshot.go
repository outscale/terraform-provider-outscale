package outscale

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/outscale/terraform-provider-outscale/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceOutscaleSnapshot() *schema.Resource {
	return &schema.Resource{
		Create: ResourceOutscaleSnapshotCreate,
		Read:   ResourceOutscaleSnapshotRead,
		Update: ResourceOutscaleSnapshotUpdate,
		Delete: ResourceOutscaleSnapshotDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(40 * time.Minute),
			Update: schema.DefaultTimeout(40 * time.Minute),
			Delete: schema.DefaultTimeout(40 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"description": {
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
			"creation_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"permissions_to_create_volume": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"global_permission": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"account_id": {
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

func ResourceOutscaleSnapshotCreate(d *schema.ResourceData, meta interface{}) error {
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
	err := retry.RetryContext(context.Background(), d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		var err error
		rp, httpResp, err := conn.SnapshotApi.CreateSnapshot(context.Background()).CreateSnapshotRequest(request).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		resp = rp
		return nil
	})
	if err != nil {
		return err
	}

	log.Printf("Waiting for Snapshot %s to become available...", resp.Snapshot.GetSnapshotId())

	stateConf := &retry.StateChangeConf{
		Pending:    []string{"pending", "in-queue", "queued", "importing"},
		Target:     []string{"completed"},
		Refresh:    SnapshotOAPIStateRefreshFunc(conn, resp.Snapshot.GetSnapshotId(), d.Timeout(schema.TimeoutCreate)),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      OutscaleImageRetryDelay,
		MinTimeout: OutscaleImageRetryMinTimeout,
	}

	if _, err = stateConf.WaitForStateContext(context.Background()); err != nil {
		return fmt.Errorf("Error waiting for Snapshot (%s) to be ready: %w", resp.Snapshot.GetSnapshotId(), err)
	}

	if tags, ok := d.GetOk("tags"); ok {
		err := assignTags(tags.(*schema.Set), resp.Snapshot.GetSnapshotId(), conn)
		if err != nil {
			return err
		}
	}

	d.SetId(resp.Snapshot.GetSnapshotId())

	return ResourceOutscaleSnapshotRead(d, meta)
}

func ResourceOutscaleSnapshotRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	req := oscgo.ReadSnapshotsRequest{
		Filters: &oscgo.FiltersSnapshot{SnapshotIds: &[]string{d.Id()}},
	}

	var resp oscgo.ReadSnapshotsResponse
	err := retry.RetryContext(context.Background(), d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		var err error
		rp, httpResp, err := conn.SnapshotApi.ReadSnapshots(context.Background()).ReadSnapshotsRequest(req).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		resp = rp
		return nil
	})
	if err != nil {
		return fmt.Errorf("Error reading the snapshot: %w", err)
	}
	if utils.IsResponseEmpty(len(resp.GetSnapshots()), "Snapshot", d.Id()) {
		d.SetId("")
		return nil
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
		if err := set("creation_date", snapshot.GetCreationDate()); err != nil {
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
		return nil
	})
}

func ResourceOutscaleSnapshotUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	if err := setOSCAPITags(conn, d); err != nil {
		return err
	}
	return ResourceOutscaleSnapshotRead(d, meta)
}

func ResourceOutscaleSnapshotDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	return retry.RetryContext(context.Background(), d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		request := oscgo.DeleteSnapshotRequest{SnapshotId: d.Id()}
		_, httpResp, err := conn.SnapshotApi.DeleteSnapshot(context.Background()).DeleteSnapshotRequest(request).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		return nil
	})
}

// SnapshotOAPIStateRefreshFunc ...
func SnapshotOAPIStateRefreshFunc(client *oscgo.APIClient, id string, timeOut time.Duration) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		emptyResp := oscgo.ReadSnapshotsResponse{}

		var resp oscgo.ReadSnapshotsResponse
		var statusCode int
		err := retry.RetryContext(context.Background(), timeOut, func() *retry.RetryError {
			var err error
			rp, httpResp, err := client.SnapshotApi.ReadSnapshots(context.Background()).ReadSnapshotsRequest(oscgo.ReadSnapshotsRequest{
				Filters: &oscgo.FiltersSnapshot{SnapshotIds: &[]string{id}},
			}).Execute()
			if err != nil {
				return utils.CheckThrottling(httpResp, err)
			}
			resp = rp
			statusCode = httpResp.StatusCode
			return nil
		})

		if err != nil {
			if statusCode == http.StatusNotFound {
				log.Printf("[INFO] OMI %s state %s", id, "destroyed")
				return emptyResp, "destroyed", nil

			} else if len(resp.GetSnapshots()) == 0 {
				log.Printf("[INFO] OMI %s state %s", id, "destroyed")
				return emptyResp, "destroyed", nil
			} else {
				return emptyResp, "", fmt.Errorf("Error on refresh: %w", err)
			}
		}

		if resp.GetSnapshots() == nil || len(resp.GetSnapshots()) == 0 {
			return emptyResp, "destroyed", nil
		}

		// OMI is valid, so return it's state
		return resp.GetSnapshots()[0], resp.GetSnapshots()[0].GetState(), nil
	}
}
