package oapi

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/outscale/goutils/sdk/ptr"
	"github.com/outscale/osc-sdk-go/v3/pkg/options"
	"github.com/outscale/osc-sdk-go/v3/pkg/osc"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/framework/fwhelpers/from"
	"github.com/outscale/terraform-provider-outscale/internal/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceOutscaleSnapshot() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceOutscaleSnapshotCreate,
		ReadContext:   ResourceOutscaleSnapshotRead,
		UpdateContext: ResourceOutscaleSnapshotUpdate,
		DeleteContext: ResourceOutscaleSnapshotDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(40 * time.Minute),
			Read:   schema.DefaultTimeout(ReadDefaultTimeout),
			Update: schema.DefaultTimeout(UpdateDefaultTimeout),
			Delete: schema.DefaultTimeout(DeleteDefaultTimeout),
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
			"tags": TagsSchemaSDK(),
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

func ResourceOutscaleSnapshotCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.OutscaleClient).OSC

	timeout := d.Timeout(schema.TimeoutCreate)

	v, ok := d.GetOk("volume_id")
	snp, sok := d.GetOk("snapshot_size")
	source, sourceok := d.GetOk("source_snapshot_id")

	if !ok && !sok && !sourceok {
		return diag.Errorf("please provide the source_snapshot_id, volume_id or snapshot_size argument")
	}

	description := d.Get("description").(string)
	fileLocation := d.Get("file_location").(string)
	sourceRegionName := d.Get("source_region_name").(string)

	request := osc.CreateSnapshotRequest{
		Description:  &description,
		FileLocation: &fileLocation,
	}

	if ok {
		request.VolumeId = new(v.(string))
	}

	if sok && snp.(int) > 0 {
		log.Printf("[DEBUG] Snapshot Size %d", snp.(int))

		request.SnapshotSize = new(int64(snp.(int)))
	}

	if sourceok {
		request.SourceSnapshotId = new(source.(string))
	}
	if sourceRegionName != "" {
		request.SourceRegionName = new(sourceRegionName)
	}

	resp, err := client.CreateSnapshot(ctx, request, options.WithRetryTimeout(timeout))
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("Waiting for Snapshot %s to become available...", resp.Snapshot.SnapshotId)

	stateConf := &retry.StateChangeConf{
		Pending: []string{"pending", "in-queue", "queued", "importing"},
		Target:  []string{"completed"},
		Timeout: timeout,
		Refresh: SnapshotOAPIStateRefreshFunc(ctx, client, resp.Snapshot.SnapshotId, timeout),
	}

	if _, err = stateConf.WaitForStateContext(ctx); err != nil {
		return diag.Errorf("error waiting for snapshot (%s) to be ready: %v", resp.Snapshot.SnapshotId, err)
	}
	d.SetId(resp.Snapshot.SnapshotId)

	err = createOAPITagsSDK(ctx, client, timeout, d)
	if err != nil {
		return diag.FromErr(err)
	}

	return ResourceOutscaleSnapshotRead(ctx, d, meta)
}

func ResourceOutscaleSnapshotRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.OutscaleClient).OSC

	timeout := d.Timeout(schema.TimeoutRead)

	req := osc.ReadSnapshotsRequest{
		Filters: &osc.FiltersSnapshot{SnapshotIds: &[]string{d.Id()}},
	}

	resp, err := client.ReadSnapshots(ctx, req, options.WithRetryTimeout(timeout))
	if err != nil {
		return diag.Errorf("error reading the snapshot: %v", err)
	}
	if resp.Snapshots == nil || utils.IsResponseEmpty(len(*resp.Snapshots), "Snapshot", d.Id()) {
		d.SetId("")
		return nil
	}

	snapshot := (*resp.Snapshots)[0]
	return diag.FromErr(resourceDataAttrSetter(d, func(set AttributeSetter) error {
		permisions := snapshot.PermissionsToCreateVolume
		if err := set("description", snapshot.Description); err != nil {
			return err
		}
		if err := set("volume_id", snapshot.VolumeId); err != nil {
			return err
		}
		if err := set("account_alias", snapshot.AccountAlias); err != nil {
			return err
		}
		if err := set("account_id", snapshot.AccountId); err != nil {
			return err
		}
		if err := set("creation_date", from.ISO8601(snapshot.CreationDate)); err != nil {
			return err
		}
		if err := set("permissions_to_create_volume", omiOAPIPermissionToLuch(permisions)); err != nil {
			return err
		}
		if err := set("progress", snapshot.Progress); err != nil {
			return err
		}
		if err := set("snapshot_id", snapshot.SnapshotId); err != nil {
			return err
		}
		if err := set("state", snapshot.State); err != nil {
			return err
		}
		if err := set("tags", FlattenOAPITagsSDK(ptr.From(snapshot.Tags))); err != nil {
			return err
		}
		if err := set("volume_size", snapshot.VolumeSize); err != nil {
			return err
		}
		return nil
	}))
}

func ResourceOutscaleSnapshotUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.OutscaleClient).OSC
	timeout := d.Timeout(schema.TimeoutUpdate)

	if err := updateOAPITagsSDK(ctx, client, timeout, d); err != nil {
		return diag.FromErr(err)
	}
	return ResourceOutscaleSnapshotRead(ctx, d, meta)
}

func ResourceOutscaleSnapshotDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.OutscaleClient).OSC
	timeout := d.Timeout(schema.TimeoutDelete)

	request := osc.DeleteSnapshotRequest{SnapshotId: d.Id()}
	_, err := client.DeleteSnapshot(ctx, request, options.WithRetryTimeout(timeout))

	return diag.FromErr(err)
}

// SnapshotOAPIStateRefreshFunc ...
func SnapshotOAPIStateRefreshFunc(ctx context.Context, client *osc.Client, id string, to time.Duration) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		emptyResp := osc.ReadSnapshotsResponse{}

		resp, err := client.ReadSnapshots(ctx, osc.ReadSnapshotsRequest{
			Filters: &osc.FiltersSnapshot{SnapshotIds: &[]string{id}},
		}, options.WithRetryTimeout(to))
		if err != nil {
			return emptyResp, "", fmt.Errorf("error on refresh: %w", err)
		}

		if resp.Snapshots == nil || len(*resp.Snapshots) == 0 {
			return emptyResp, "destroyed", nil
		}

		// OMI is valid, so return it's state
		return (*resp.Snapshots)[0], string((*resp.Snapshots)[0].State), nil
	}
}
