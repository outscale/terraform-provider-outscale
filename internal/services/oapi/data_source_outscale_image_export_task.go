package oapi

import (
	"context"
	"strings"
	"time"

	"github.com/outscale/goutils/sdk/ptr"
	"github.com/outscale/osc-sdk-go/v3/pkg/options"
	"github.com/outscale/osc-sdk-go/v3/pkg/osc"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceOutscaleImageExportTask() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceOAPISnapshotImageTaskRead,
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

func dataSourceOAPISnapshotImageTaskRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.OutscaleClient).OSC

	filters, filtersOk := d.GetOk("filter")

	var err error
	filtersReq := &osc.FiltersExportTask{}
	if filtersOk {
		filtersReq, err = buildOutscaleOSCAPIDataSourceImageExportTaskFilters(filters.(*schema.Set))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	resp, err := client.ReadImageExportTasks(ctx, osc.ReadImageExportTasksRequest{
		Filters: filtersReq,
	}, options.WithRetryTimeout(5*time.Minute))
	if err != nil {
		return diag.Errorf("error reading task image %s", err)
	}

	if resp.ImageExportTasks == nil || len(*resp.ImageExportTasks) == 0 {
		return diag.FromErr(ErrNoResults)
	}
	v := (*resp.ImageExportTasks)[0]

	if err = d.Set("progress", ptr.From(v.Progress)); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("task_id", ptr.From(v.TaskId)); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("state", ptr.From(v.State)); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("comment", ptr.From(v.Comment)); err != nil {
		return diag.FromErr(err)
	}

	exp := make([]map[string]interface{}, 1)
	exportToOsu := make(map[string]interface{})
	export := ptr.From(v.OsuExport)
	exportToOsu["disk_image_format"] = export.DiskImageFormat
	exportToOsu["osu_bucket"] = export.OsuBucket
	osuPrefix := ptr.From(export.OsuPrefix)
	if strings.Contains(osuPrefix, "/") {
		osuList := strings.Split(osuPrefix, "/")
		osuPrefix = osuList[0]
	}
	exportToOsu["osu_prefix"] = osuPrefix
	exportToOsu["osu_manifest_url"] = ptr.From(export.OsuManifestUrl)

	exp[0] = exportToOsu

	if err = d.Set("image_id", ptr.From(v.ImageId)); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("osu_export", exp); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("tags", FlattenOAPITagsSDK(ptr.From(v.Tags))); err != nil {
		return diag.FromErr(err)
	}
	d.SetId(ptr.From(v.TaskId))

	return nil
}

func buildOutscaleOSCAPIDataSourceImageExportTaskFilters(set *schema.Set) (*osc.FiltersExportTask, error) {
	var filters osc.FiltersExportTask
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
			return nil, utils.UnknownDataSourceFilterError(name)
		}
	}
	return &filters, nil
}
