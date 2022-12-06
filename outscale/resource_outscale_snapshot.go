package outscale

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	oscgo "github.com/outscale/osc-sdk-go/v2"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-outscale/utils"
)

func resourceOutscaleOAPISnapshot() *schema.Resource {
	return &schema.Resource{
		Create: resourceOutscaleOAPISnapshotCreate,
		Read:   resourceOutscaleOAPISnapshotRead,
		Update: resourceOutscaleOAPISnapshotUpdate,
		Delete: resourceOutscaleOAPISnapshotDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"snapshot_size": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"file_location": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"source_region_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"source_snapshot_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"volume_id": {
				Type:     schema.TypeString,
				Optional: true,
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
			"permissions_to_create_volume_global_permission": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"permissions_to_create_volume_account_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
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
	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
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
	d.SetId(resp.Snapshot.GetSnapshotId())
	log.Printf("Waiting for Snapshot %s to become available...", d.Id())

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"pending", "pending/queued", "queued", "importing"},
		Target:     []string{"completed"},
		Refresh:    SnapshotOAPIStateRefreshFunc(conn, d.Id()),
		Timeout:    OutscaleImageRetryTimeout,
		Delay:      OutscaleImageRetryDelay,
		MinTimeout: OutscaleImageRetryMinTimeout,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Error waiting for Snapshot (%s) to be ready: %s", d.Id(), err)
	}

	if tags, ok := d.GetOk("tags"); ok {
		err := assignTags(tags.(*schema.Set), d.Id(), conn)
		if err != nil {
			return err
		}
	}

	if _, ok := d.GetOk("permissions_to_create_volume_accounts_ids"); ok || d.Get("permissions_to_create_volume_global_permission").(bool) {
		if err := UpdateSnapshot(d, meta); err != nil {
			return err
		}
	}

	return resourceOutscaleOAPISnapshotRead(d, meta)
}

func resourceOutscaleOAPISnapshotRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	req := oscgo.ReadSnapshotsRequest{
		Filters: &oscgo.FiltersSnapshot{SnapshotIds: &[]string{d.Id()}},
	}

	var resp oscgo.ReadSnapshotsResponse
	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		var err error
		rp, httpResp, err := conn.SnapshotApi.ReadSnapshots(context.Background()).ReadSnapshotsRequest(req).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		resp = rp
		return nil
	})
	if err != nil {
		return fmt.Errorf("Error reading the snapshot %snapshot", err)
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
		if err := set("permissions_to_create_volume_account_ids", permisions.GetAccountIds()); err != nil {
			return err
		}
		if err := set("permissions_to_create_volume_global_permission", permisions.GetGlobalPermission()); err != nil {
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

func UpdateSnapshot(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	snapshotID := d.Id()

	req := oscgo.UpdateSnapshotRequest{
		SnapshotId: snapshotID,
	}

	oldAccount, newAccount := d.GetChange("permissions_to_create_volume_account_ids")
	inter := oldAccount.(*schema.Set).Intersection(newAccount.(*schema.Set))
	added := newAccount.(*schema.Set).Difference(inter).List()
	removed := oldAccount.(*schema.Set).Difference(inter).List()

	globalPermission := d.Get("permissions_to_create_volume_global_permission").(bool)

	if len(added) > 0 || globalPermission {
		perms := oscgo.PermissionsOnResourceCreation{}
		addition := oscgo.PermissionsOnResource{}
		if len(added) > 0 {
			addition.SetAccountIds(utils.InterfaceSliceToStringSlice(added))
		}
		if globalPermission {
			addition.SetGlobalPermission(true)
		}
		perms.SetAdditions(addition)
		req.SetPermissionsToCreateVolume(perms)
		var err error
		err = resource.Retry(2*time.Minute, func() *resource.RetryError {
			_, httpResp, err := conn.SnapshotApi.UpdateSnapshot(context.Background()).UpdateSnapshotRequest(req).Execute()
			if err != nil {
				return utils.CheckThrottling(httpResp, err)
			}
			return nil
		})
		if err != nil {
			return fmt.Errorf("Error updating snapshot: %s", err)
		}
	}
	if len(removed) > 0 || !globalPermission {
		perms := oscgo.PermissionsOnResourceCreation{}
		removal := oscgo.PermissionsOnResource{}
		if len(removed) > 0 {
			removal.SetAccountIds(utils.InterfaceSliceToStringSlice(removed))
		}
		if !globalPermission {
			removal.SetGlobalPermission(true)
		}
		perms.SetRemovals(removal)
		req.SetPermissionsToCreateVolume(perms)

		var err error
		err = resource.Retry(2*time.Minute, func() *resource.RetryError {
			_, httpResp, err := conn.SnapshotApi.UpdateSnapshot(context.Background()).UpdateSnapshotRequest(req).Execute()
			if err != nil {
				return utils.CheckThrottling(httpResp, err)
			}
			return nil
		})
		if err != nil {
			return fmt.Errorf("Error updating snapshot: %s", err)
		}
	}
	return nil
}

func resourceOutscaleOAPISnapshotUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	if err := setOSCAPITags(conn, d); err != nil {
		return err
	}
	if d.HasChanges("permissions_to_create_volume_account_ids", "permissions_to_create_volume_global_permission") {
		if err := UpdateSnapshot(d, meta); err != nil {
			return err
		}
	}
	return resourceOutscaleOAPISnapshotRead(d, meta)
}

func resourceOutscaleOAPISnapshotDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		request := oscgo.DeleteSnapshotRequest{SnapshotId: d.Id()}
		_, httpResp, err := conn.SnapshotApi.DeleteSnapshot(context.Background()).DeleteSnapshotRequest(request).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		return nil
	})
}

// SnapshotOAPIStateRefreshFunc ...
func SnapshotOAPIStateRefreshFunc(client *oscgo.APIClient, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		emptyResp := oscgo.ReadSnapshotsResponse{}

		var resp oscgo.ReadSnapshotsResponse
		var statusCode int
		err := resource.Retry(5*time.Minute, func() *resource.RetryError {
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
