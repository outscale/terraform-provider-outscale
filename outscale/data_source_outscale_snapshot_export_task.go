package outscale

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/terraform-providers/terraform-provider-outscale/utils"
)

func dataSourceOutscaleOAPISnapshotExportTask() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOAPISnapshotExportTaskRead,
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
			"request_id": {
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
			"tags": dataSourceTagsSchema(),
		},
	}
}

func dataSourceOAPISnapshotExportTaskRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	filters, filtersOk := d.GetOk("filter")

	filtersReq := &oscgo.FiltersExportTask{}
	if filtersOk {
		filtersReq = buildOutscaleOSCAPIDataSourceSnapshotExportTaskFilters(filters.(*schema.Set))
	}

	var resp oscgo.ReadSnapshotExportTasksResponse
	var err error
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
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
		return fmt.Errorf("Error reading task snapshot export %s", err)
	}

	if len(resp.GetSnapshotExportTasks()) == 0 {
		return fmt.Errorf("your query returned no results, please change your search criteria and try again")
	}
	v := resp.GetSnapshotExportTasks()[0]

	if err = d.Set("progress", v.GetProgress()); err != nil {
		return err
	}
	if err = d.Set("task_id", v.GetTaskId()); err != nil {
		return err
	}
	if err = d.Set("state", v.GetState()); err != nil {
		return err
	}
	if err = d.Set("comment", v.GetComment()); err != nil {
		return err
	}

	exp := make([]map[string]interface{}, 1)
	exportToOsu := make(map[string]interface{})
	exportToOsu["disk_image_format"] = v.OsuExport.GetDiskImageFormat()
	exportToOsu["osu_bucket"] = v.OsuExport.GetOsuBucket()
	exportToOsu["osu_prefix"] = v.OsuExport.GetOsuPrefix()

	exp[0] = exportToOsu

	if err = d.Set("snapshot_id", v.GetSnapshotId()); err != nil {
		return err
	}
	if err = d.Set("osu_export", exp); err != nil {
		return err
	}
	if err = d.Set("tags", tagsOSCAPIToMap(v.GetTags())); err != nil {
		return err
	}

	d.SetId(v.GetTaskId())

	return nil
}

func buildOutscaleOSCAPIDataSourceSnapshotExportTaskFilters(set *schema.Set) *oscgo.FiltersExportTask {
	var filters oscgo.FiltersExportTask
	for _, v := range set.List() {
		m := v.(map[string]interface{})
		var filterValues []string
		for _, e := range m["values"].([]interface{}) {
			filterValues = append(filterValues, e.(string))
		}

		switch name := m["name"].(string); name {
		case "task_ids":
			filters.TaskIds = &filterValues
		default:
			log.Printf("[Debug] Unknown Filter Name: %s.", name)
		}
	}
	return &filters
}
