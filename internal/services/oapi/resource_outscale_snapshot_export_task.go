package oapi

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/outscale/osc-sdk-go/v3/pkg/options"
	"github.com/outscale/osc-sdk-go/v3/pkg/osc"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/utils"
)

func ResourceOutscaleSnapshotExportTask() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceOAPISnapshotExportTaskCreate,
		ReadContext:   resourceOAPISnapshotExportTaskRead,
		UpdateContext: resourceOAPISnapshotExportTaskUpdate,
		DeleteContext: resourceOAPISnapshotExportTaskDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(40 * time.Minute),
			Read:   schema.DefaultTimeout(ReadDefaultTimeout),
			Update: schema.DefaultTimeout(UpdateDefaultTimeout),
			Delete: schema.DefaultTimeout(40 * time.Minute),
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
							Computed: true,
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
										ForceNew: true,
									},
									"secret_key": {
										Type:     schema.TypeString,
										Required: true,
										ForceNew: true,
									},
								},
							},
						},
					},
				},
			},
			"snapshot_id": {
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

func resourceOAPISnapshotExportTaskCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.OutscaleClient).OSC

	timeout := d.Timeout(schema.TimeoutCreate)

	eto, etoOk := d.GetOk("osu_export")
	v, ok := d.GetOk("snapshot_id")
	request := osc.CreateSnapshotExportTaskRequest{}

	if !etoOk && !ok {
		return diag.Errorf("please provide the required attributes osu_export and snapshot_id")
	}

	request.SnapshotId = v.(string)

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
		request.OsuExport = et
	}

	resp, err := client.CreateSnapshotExportTask(ctx, request, options.WithRetryTimeout(timeout))
	if err != nil {
		return diag.Errorf("error snapshot export task %s", err)
	}

	id := resp.SnapshotExportTask.TaskId

	wait := d.Get("wait_for_completion").(bool)
	if wait {
		err = ResourceOutscaleSnapshotTaskWaitForAvailable(client, ctx, id, timeout)
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

	return resourceOAPISnapshotExportTaskRead(ctx, d, meta)
}

func resourceOAPISnapshotExportTaskRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.OutscaleClient).OSC
	timeout := d.Timeout(schema.TimeoutRead)

	filter := &osc.FiltersExportTask{TaskIds: &[]string{d.Id()}}
	resp, err := client.ReadSnapshotExportTasks(ctx, osc.ReadSnapshotExportTasksRequest{
		Filters: filter,
	}, options.WithRetryTimeout(timeout))
	if err != nil {
		return diag.Errorf("error reading snapshot export task %s", err)
	}
	if resp.SnapshotExportTasks == nil || utils.IsResponseEmpty(len(*resp.SnapshotExportTasks), "SnapshotExportTask", d.Id()) {
		d.SetId("")
		return nil
	}
	v := (*resp.SnapshotExportTasks)[0]

	if v.State == "failed" || v.State == "cancelled" {
		taskId := d.Id()
		d.SetId("")

		errMsg := fmt.Sprintf("Snapshot export task (%s) did not succeed. Status: %s", taskId, v.State)
		errMsg += "\n" + `To remove it from state, run "terraform refresh" or recreate the resource.`

		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  "Snapshot export task " + string(v.State),
				Detail:   errMsg,
			},
		}
	}

	if err = d.Set("progress", v.Progress); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("task_id", v.TaskId); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("state", v.State); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("comment", v.Comment); err != nil {
		return diag.FromErr(err)
	}

	exp := make([]map[string]interface{}, 1)
	exportToOsu := make(map[string]interface{})
	exportToOsu["disk_image_format"] = v.OsuExport.DiskImageFormat
	exportToOsu["osu_bucket"] = v.OsuExport.OsuBucket
	exportToOsu["osu_prefix"] = v.OsuExport.OsuPrefix

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

	if err = d.Set("snapshot_id", v.SnapshotId); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("osu_export", exp); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("tags", FlattenOAPITagsSDK(v.Tags)); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceOAPISnapshotExportTaskUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.OutscaleClient).OSC
	timeout := d.Timeout(schema.TimeoutUpdate)

	if err := updateOAPITagsSDK(ctx, client, timeout, d); err != nil {
		return diag.FromErr(err)
	}
	return resourceOAPISnapshotExportTaskRead(ctx, d, meta)
}

func ResourceOutscaleSnapshotTaskWaitForAvailable(client *osc.Client, ctx context.Context, id string, timeout time.Duration) error {
	log.Printf("Waiting for Snapshot Task %s to become available...", id)
	stateConf := &retry.StateChangeConf{
		Pending: []string{"pending", "pending/queued", "queued"},
		Target:  []string{"completed", "active"},
		Timeout: timeout,
		Refresh: SnapshotTaskStateRefreshFunc(ctx, client, id, timeout),
	}

	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func resourceOAPISnapshotExportTaskDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	d.SetId("")
	if err := d.Set("osu_export", nil); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

// SnapshotTaskStateRefreshFunc ...
func SnapshotTaskStateRefreshFunc(ctx context.Context, client *osc.Client, id string, timeout time.Duration) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		filter := &osc.FiltersExportTask{TaskIds: &[]string{id}}
		resp, err := client.ReadSnapshotExportTasks(ctx, osc.ReadSnapshotExportTasksRequest{
			Filters: filter,
		}, options.WithRetryTimeout(timeout))
		if err != nil {
			return resp, "", fmt.Errorf("error on refresh: %+v", err)
		}

		if resp.SnapshotExportTasks == nil || len(*resp.SnapshotExportTasks) == 0 {
			return resp, "destroyed", nil
		}
		tasks := (*resp.SnapshotExportTasks)[0]

		if tasks.State == osc.SnapshotExportTaskStateFailed || tasks.State == osc.SnapshotExportTaskStateCancelled {
			return tasks, string(tasks.State),
				fmt.Errorf("error: %v", tasks.Comment)
		}

		// Snapshot export task is valid, so return it's state
		return tasks, string(tasks.State), nil
	}
}
