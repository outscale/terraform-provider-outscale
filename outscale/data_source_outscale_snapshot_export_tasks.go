package outscale

import (
	"context"
	"fmt"
	"time"

	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/outscale/terraform-provider-outscale/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceOutscaleSnapshotExportTasks() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOAPISnapshotExportTasksRead,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
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
	conn := meta.(*OutscaleClient).OSCAPI

	filters, filtersOk := d.GetOk("filter")

	var err error
	var filtersReq *oscgo.FiltersExportTask
	if filtersOk {
		filtersReq, err = buildOutscaleOSCAPIDataSourceSnapshotExportTaskFilters(filters.(*schema.Set))
		if err != nil {
			return err
		}
	}

	var resp oscgo.ReadSnapshotExportTasksResponse
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		rp, httpResp, err := conn.SnapshotApi.ReadSnapshotExportTasks(context.Background()).
			ReadSnapshotExportTasksRequest(oscgo.ReadSnapshotExportTasksRequest{
				Filters: filtersReq,
			}).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		resp = rp
		return nil
	})

	if err != nil {
		return fmt.Errorf("Error reading task image %s", err)
	}

	if len(resp.GetSnapshotExportTasks()) == 0 {
		return fmt.Errorf("your query returned no results, please change your search criteria and try again")
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

		snapshot["tags"] = flattenOAPITagsSDK(v.GetTags())

		snapshots[k] = snapshot
	}

	d.SetId(id.UniqueId())

	return d.Set("snapshot_export_tasks", snapshots)
}
