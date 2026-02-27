package oapi

import (
	"context"
	"time"

	"github.com/outscale/osc-sdk-go/v3/pkg/options"
	"github.com/outscale/osc-sdk-go/v3/pkg/osc"
	"github.com/outscale/terraform-provider-outscale/internal/client"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceOutscaleSnapshotExportTasks() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceOAPISnapshotExportTasksRead,
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

func dataSourceOAPISnapshotExportTasksRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.OutscaleClient).OSC

	filters, filtersOk := d.GetOk("filter")

	var err error
	var filtersReq *osc.FiltersExportTask
	if filtersOk {
		filtersReq, err = buildOutscaleOSCAPIDataSourceSnapshotExportTaskFilters(filters.(*schema.Set))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	resp, err := client.ReadSnapshotExportTasks(ctx, osc.ReadSnapshotExportTasksRequest{
		Filters: filtersReq,
	}, options.WithRetryTimeout(5*time.Minute))
	if err != nil {
		return diag.Errorf("error reading task image %s", err)
	}

	if resp.SnapshotExportTasks == nil || len(*resp.SnapshotExportTasks) == 0 {
		return diag.FromErr(ErrNoResults)
	}

	snapshots := make([]map[string]interface{}, len(*resp.SnapshotExportTasks))
	for k, v := range *resp.SnapshotExportTasks {
		snapshot := make(map[string]interface{})

		snapshot["progress"] = v.Progress
		snapshot["task_id"] = v.TaskId
		snapshot["state"] = v.State
		snapshot["comment"] = v.Comment

		exp := make([]map[string]interface{}, 1)
		exportToOsu := make(map[string]interface{})
		exportToOsu["disk_image_format"] = v.OsuExport.DiskImageFormat
		exportToOsu["osu_bucket"] = v.OsuExport.OsuBucket
		exportToOsu["osu_prefix"] = v.OsuExport.OsuPrefix

		exp[0] = exportToOsu

		snapshot["snapshot_id"] = v.SnapshotId
		snapshot["osu_export"] = exp

		snapshot["tags"] = FlattenOAPITagsSDK(v.Tags)

		snapshots[k] = snapshot
	}

	d.SetId(id.UniqueId())

	return diag.FromErr(d.Set("snapshot_export_tasks", snapshots))
}
