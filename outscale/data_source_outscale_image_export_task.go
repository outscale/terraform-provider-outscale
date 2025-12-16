package outscale

import (
	"context"
	"fmt"
	"strings"
	"time"

	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/outscale/terraform-provider-outscale/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceOutscaleImageExportTask() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOAPISnapshotImageTaskRead,
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
						"osu_manifest_url": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"image_id": {
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
			"tags": TagsSchemaComputedSDK(),
		},
	}
}

func dataSourceOAPISnapshotImageTaskRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	filters, filtersOk := d.GetOk("filter")

	var err error
	filtersReq := &oscgo.FiltersExportTask{}
	if filtersOk {
		filtersReq, err = buildOutscaleOSCAPIDataSourceImageExportTaskFilters(filters.(*schema.Set))
		if err != nil {
			return err
		}
	}

	var resp oscgo.ReadImageExportTasksResponse
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		rp, httpResp, err := conn.ImageApi.ReadImageExportTasks(context.Background()).
			ReadImageExportTasksRequest(oscgo.ReadImageExportTasksRequest{
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

	if len(resp.GetImageExportTasks()) == 0 {
		return fmt.Errorf("your query returned no results, please change your search criteria and try again")
	}
	v := resp.GetImageExportTasks()[0]

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
	osuPrefix := v.OsuExport.GetOsuPrefix()
	if strings.Contains(osuPrefix, "/") {
		osuList := strings.Split(osuPrefix, "/")
		osuPrefix = osuList[0]
	}
	exportToOsu["osu_prefix"] = osuPrefix
	exportToOsu["osu_manifest_url"] = v.OsuExport.GetOsuManifestUrl()

	exp[0] = exportToOsu

	if err = d.Set("image_id", v.GetImageId()); err != nil {
		return err
	}
	if err = d.Set("osu_export", exp); err != nil {
		return err
	}
	if err = d.Set("tags", flattenOAPITagsSDK(v.GetTags())); err != nil {
		return err
	}
	d.SetId(v.GetTaskId())

	return nil
}

func buildOutscaleOSCAPIDataSourceImageExportTaskFilters(set *schema.Set) (*oscgo.FiltersExportTask, error) {
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
			return nil, utils.UnknownDataSourceFilterError(context.Background(), name)
		}
	}
	return &filters, nil
}
