package oapi

import (
	"context"
	"fmt"
	"time"

	"github.com/outscale/osc-sdk-go/v3/pkg/osc"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceOutscaleSnapshotExportTasks() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOAPISnapshotExportTasksRead,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(40 * time.Minute),
			Delete: schema.DefaultTimeout(40 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			"dry_run": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"snapshot_export_tasks": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"osu_export": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"disk_image_format": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"osu_bucket": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"osu_prefix": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"snapshot_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"progress": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"task_id": {
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
						"tags": TagsSchemaComputedSDK(),
					},
				},
			},
			"request_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceOAPISnapshotExportTasksRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*client.OutscaleClient).OSC
	

	filters, filtersOk := d.GetOk("filter")

	var err error
	var filtersReq *osc.FiltersExportTask
	if filtersOk {
		filtersReq, err = buildOutscaleOSCAPIDataSourceSnapshotExportTaskFilters(filters.(*schema.Set))
		if err != nil {
			return err
		}
	}

	var resp osc.ReadSnapshotExportTasksResponse
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		rp, httpResp, err := client.SnapshotApi.ReadSnapshotExportTasks(ctx).
			ReadSnapshotExportTasksRequest(osc.ReadSnapshotExportTasksRequest{
				Filters: filtersReq,
			}).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		resp = rp
		return nil
	})
	if err != nil {
		return fmt.Errorf("error reading task image %s", err)
	}

	if len(resp.GetSnapshotExportTasks()) == 0 {
		return ErrNoResults
	}

	snapshots := make([]map[string]interface{}, len(resp.GetSnapshotExportTasks()))
	for k, v := range resp.GetSnapshotExportTasks() {
		snapshot := make(map[string]interface{})

		snapshot["progress"] = v.GetProgress()
		snapshot["task_id"] = v.GetTaskId()
		snapshot["state"] = v.GetState()
		snapshot["comment"] = v.GetComment()

		exp := make([]map[string]interface{}, 1)
		exportToOsu := make(map[string]interface{})
		exportToOsu["disk_image_format"] = v.OsuExport.GetDiskImageFormat()
		exportToOsu["osu_bucket"] = v.OsuExport.GetOsuBucket()
		exportToOsu["osu_prefix"] = v.OsuExport.GetOsuPrefix()

		exp[0] = exportToOsu

		snapshot["snapshot_id"] = v.GetSnapshotId()
		snapshot["osu_export"] = exp

		snapshot["tags"] = FlattenOAPITagsSDK(v.Tags)

		snapshots[k] = snapshot
	}

	d.SetId(id.UniqueId())

	return d.Set("snapshot_export_tasks", snapshots)
}
