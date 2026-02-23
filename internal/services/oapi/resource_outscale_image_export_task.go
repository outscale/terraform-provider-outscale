package oapi

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/outscale/goutils/sdk/ptr"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/utils"

	"github.com/outscale/osc-sdk-go/v3/pkg/options"
	"github.com/outscale/osc-sdk-go/v3/pkg/osc"
)

func ResourceOutscaleImageExportTask() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceOAPIImageExportTaskCreate,
		ReadContext:   resourceOAPIImageExportTaskRead,
		UpdateContext: resourceOAPIImageExportTaskUpdate,
		DeleteContext: resourceOAPIImageExportTaskDelete,
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
			"osu_export": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"disk_image_format": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"osu_bucket": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"osu_prefix": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"osu_api_key": {
							Type:     schema.TypeList,
							Optional: true,
							ForceNew: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"api_key_id": {
										Type:     schema.TypeString,
										Required: true,
									},
									"secret_key": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},
						"osu_manifest_url": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},
					},
				},
			},
			"image_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
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
			"wait_for_completion": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Default:  true,
			},
			"tags": TagsSchemaSDK(),
		},
	}
}

func resourceOAPIImageExportTaskCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.OutscaleClient).OSC

	timeout := d.Timeout(schema.TimeoutCreate)

	eto, etoOk := d.GetOk("osu_export")
	v, ok := d.GetOk("image_id")
	request := osc.CreateImageExportTaskRequest{}

	if !etoOk && !ok {
		return diag.Errorf("please provide the required attributes osu_export and image_id")
	}

	request.ImageId = v.(string)

	if etoOk {
		exp := eto.([]interface{})
		e := exp[0].(map[string]interface{})

		et := osc.OsuExportToCreate{}

		if v, ok := e["disk_image_format"]; ok {
			et.DiskImageFormat = v.(string)
		}
		/*if v, ok := e["osu_key"]; ok {
			apikey := osc.OsuApiKey{ApiKeyId: v.(*string)}
			et.SetOsuApiKey(apikey)
		}*/
		if v, ok := e["osu_bucket"]; ok {
			et.OsuBucket = v.(string)
		}
		if v, ok := e["osu_prefix"]; ok {
			et.OsuPrefix = new(v.(string))
		}
		if v, ok := e["osu_api_key"]; ok {
			a := v.([]interface{})
			if len(a) > 0 {
				w := a[0].(map[string]interface{})
				et.OsuApiKey.ApiKeyId = new(w["api_key_id"].(string))
				et.OsuApiKey.SecretKey = new(w["secret_key"].(string))
			}
		}
		if v, ok := e["osu_manifest_url"]; ok {
			et.OsuManifestUrl = new(v.(string))
		}
		request.OsuExport = et
	}

	resp, err := client.CreateImageExportTask(ctx, request, options.WithRetryTimeout(timeout))
	if err != nil {
		return diag.Errorf("error image task %s", err)
	}

	id := *resp.ImageExportTask.TaskId

	wait := d.Get("wait_for_completion").(bool)
	if wait {
		err = ResourceOutscaleImageTaskWaitForAvailable(ctx, id, client, timeout)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId(id)
	if d.IsNewResource() {
		if err := updateOAPITagsSDK(ctx, client, timeout, d); err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceOAPIImageExportTaskRead(ctx, d, meta)
}

func resourceOAPIImageExportTaskRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.OutscaleClient).OSC
	timeout := d.Timeout(schema.TimeoutRead)

	filter := &osc.FiltersExportTask{TaskIds: &[]string{d.Id()}}
	resp, err := client.ReadImageExportTasks(ctx, osc.ReadImageExportTasksRequest{
		Filters: filter,
	}, options.WithRetryTimeout(timeout))
	if err != nil {
		return diag.Errorf("error reading task image %s", err)
	}
	if resp.ImageExportTasks == nil || utils.IsResponseEmpty(len(*resp.ImageExportTasks), "ImageExportTask", d.Id()) {
		d.SetId("")
		return nil
	}
	v := (*resp.ImageExportTasks)[0]

	if ptr.From(v.State) == "failed" || *v.State == "cancelled" {
		taskId := d.Id()
		d.SetId("")

		errMsg := fmt.Sprintf("Image export task (%s) did not succeed. Status: %s", taskId, ptr.From(v.State))
		errMsg += "\n" + `To remove it from state, run "terraform refresh" or recreate the resource.`

		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  "Image export task " + *v.State,
				Detail:   errMsg,
			},
		}
	}

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
	export := ptr.From(v.OsuExport)
	exportToOsu := make(map[string]interface{})
	exportToOsu["disk_image_format"] = export.DiskImageFormat
	exportToOsu["osu_bucket"] = export.OsuBucket
	osuPrefix := ptr.From(export.OsuPrefix)
	if strings.Contains(osuPrefix, "/") {
		osuList := strings.Split(osuPrefix, "/")
		osuPrefix = osuList[0]
	}
	exportToOsu["osu_prefix"] = osuPrefix
	exportToOsu["osu_manifest_url"] = export.OsuManifestUrl

	eto, etoOk := d.GetOk("osu_export")
	if etoOk {
		exp2 := eto.([]interface{})
		e := exp2[0].(map[string]interface{})
		if v, ok := e["osu_api_key"]; ok {
			a := v.([]interface{})
			if len(a) > 0 {
				w := a[0].(map[string]interface{})
				apk := make([]map[string]interface{}, 1)
				osuAkSk := make(map[string]interface{})
				osuAkSk["api_key_id"] = w["api_key_id"].(string)
				osuAkSk["secret_key"] = w["secret_key"].(string)
				apk[0] = osuAkSk
				exportToOsu["osu_api_key"] = apk
			}
		}
	}

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

	return nil
}

func resourceOAPIImageExportTaskUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.OutscaleClient).OSC
	timeout := d.Timeout(schema.TimeoutUpdate)

	if err := updateOAPITagsSDK(ctx, client, timeout, d); err != nil {
		return diag.FromErr(err)
	}
	return resourceOAPIImageExportTaskRead(ctx, d, meta)
}

func ResourceOutscaleImageTaskWaitForAvailable(ctx context.Context, id string, client *osc.Client, timeout time.Duration) error {
	log.Printf("Waiting for Image Task %s to become available...", id)
	stateConf := &retry.StateChangeConf{
		Pending: []string{string(osc.ImageStatePending)},
		Target:  []string{string(osc.ImageStateAvailable)},
		Timeout: timeout,
		Refresh: ImageTaskStateRefreshFunc(client, ctx, id, timeout),
	}

	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func resourceOAPIImageExportTaskDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	d.SetId("")
	if err := d.Set("osu_export", nil); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

// ImageTaskStateRefreshFunc ...
func ImageTaskStateRefreshFunc(client *osc.Client, ctx context.Context, id string, timeout time.Duration) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := client.ReadImageExportTasks(ctx, osc.ReadImageExportTasksRequest{
			Filters: &osc.FiltersExportTask{TaskIds: &[]string{id}},
		}, options.WithRetryTimeout(timeout))
		if err != nil {
			return resp, "", fmt.Errorf("error on refresh: %w", err)
		}

		if resp.ImageExportTasks == nil || len(*resp.ImageExportTasks) == 0 {
			return resp, "destroyed", nil
		}
		task := (*resp.ImageExportTasks)[0]

		if ptr.From(task.State) == "failed" || *task.State == "cancelled" {
			return task, *task.State, fmt.Errorf("error: %v", *task.Comment)
		}

		// Image export task is valid, so return it's state
		return task, *task.State, nil
	}
}
