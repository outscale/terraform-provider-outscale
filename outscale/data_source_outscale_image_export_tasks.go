package outscale

import (
	"context"
	"fmt"
	"strings"
	"time"

	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/terraform-providers/terraform-provider-outscale/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceOutscaleOAPIImageExportTasks() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOAPIImageExportTasksRead,
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
			"image_export_tasks": {
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
				},
			},
			"request_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceOAPIImageExportTasksRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	filters, filtersOk := d.GetOk("filter")

	filtersReq := &oscgo.FiltersExportTask{}
	if filtersOk {
		filtersReq = buildOutscaleOSCAPIDataSourceImageExportTaskFilters(filters.(*schema.Set))
	}

	var resp oscgo.ReadImageExportTasksResponse
	var err error
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
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
		return fmt.Errorf("error reading task image %s", err)
	}

	if len(resp.GetImageExportTasks()) == 0 {
		return fmt.Errorf("your query returned no results, please change your search criteria and try again")
	}

	snapshots := make([]map[string]interface{}, len(resp.GetImageExportTasks()))
	for k, v := range resp.GetImageExportTasks() {
		snapshot := make(map[string]interface{})

		snapshot["progress"] = v.GetProgress()
		snapshot["task_id"] = v.GetTaskId()
		snapshot["state"] = v.GetState()
		snapshot["comment"] = v.GetComment()

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

		snapshot["image_id"] = v.GetImageId()
		snapshot["osu_export"] = exp

		snapshot["tags"] = tagsOSCAPIToMap(v.GetTags())

		snapshots[k] = snapshot
	}

	d.SetId(resource.UniqueId())

	return d.Set("image_export_tasks", snapshots)
}
