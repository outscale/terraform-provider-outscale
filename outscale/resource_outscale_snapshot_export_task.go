package outscale

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/terraform-providers/terraform-provider-outscale/utils"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceOutscaleOAPISnapshotExportTask() *schema.Resource {
	return &schema.Resource{
		Create: resourceOAPISnapshotExportTaskCreate,
		Read:   resourceOAPISnapshotExportTaskRead,
		Update: resourceOAPISnapshotExportTaskUpdate,
		Delete: resourceOAPISnapshotExportTaskDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(40 * time.Minute),
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
			"tags": tagsListOAPISchema(),
		},
	}
}

func resourceOAPISnapshotExportTaskCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	eto, etoOk := d.GetOk("osu_export")
	v, ok := d.GetOk("snapshot_id")
	request := oscgo.CreateSnapshotExportTaskRequest{}

	if !etoOk && !ok {
		return fmt.Errorf("please provide the required attributes osu_export and image_id")
	}

	request.SetSnapshotId(v.(string))

	if etoOk {
		exp := eto.([]interface{})
		e := exp[0].(map[string]interface{})

		et := oscgo.OsuExportToCreate{}

		if v, ok := e["disk_image_format"]; ok {
			et.SetDiskImageFormat(v.(string))
		}
		/*if v, ok := e["osu_key"]; ok {
			apikey := oscgo.OsuApiKey{ApiKeyId: v.(*string)}
			et.SetOsuApiKey(apikey)
		}*/
		if v, ok := e["osu_bucket"]; ok {
			et.SetOsuBucket(v.(string))
		}
		if v, ok := e["osu_prefix"]; ok {
			et.SetOsuPrefix(v.(string))
		}
		if v, ok := e["osu_api_key"]; ok {
			a := v.([]interface{})

			if len(a) > 0 {
				w := a[0].(map[string]interface{})
				et.OsuApiKey.SetApiKeyId(w["api_key_id"].(string))
				et.OsuApiKey.SetSecretKey(w["secret_key"].(string))
			}
		}
		request.SetOsuExport(et)
	}

	var resp oscgo.CreateSnapshotExportTaskResponse
	err := resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		var err error
		rp, httpResp, err := conn.SnapshotApi.CreateSnapshotExportTask(context.Background()).
			CreateSnapshotExportTaskRequest(request).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		resp = rp
		return nil
	})

	if err != nil {
		return fmt.Errorf("[DEBUG] Error image task %s", err)
	}

	id := resp.SnapshotExportTask.GetTaskId()
	d.SetId(id)
	if d.IsNewResource() {
		if err := setOSCAPITags(conn, d); err != nil {
			return err
		}
	}
	_, err = resourceOutscaleSnapshotTaskWaitForAvailable(id, conn, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return err
	}

	return resourceOAPISnapshotExportTaskRead(d, meta)
}

func resourceOAPISnapshotExportTaskRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	var resp oscgo.ReadSnapshotExportTasksResponse
	filter := &oscgo.FiltersExportTask{TaskIds: &[]string{d.Id()}}
	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		var err error
		rp, httpResp, err := conn.SnapshotApi.ReadSnapshotExportTasks(context.Background()).
			ReadSnapshotExportTasksRequest(oscgo.ReadSnapshotExportTasksRequest{
				Filters: filter,
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
	if utils.IsResponseEmpty(len(resp.GetSnapshotExportTasks()), "SnapshotExportTask", d.Id()) {
		d.SetId("")
		return nil
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

	if err = d.Set("snapshot_id", v.GetSnapshotId()); err != nil {
		return err
	}
	if err = d.Set("osu_export", exp); err != nil {
		return err
	}
	if err = d.Set("tags", tagsOSCAPIToMap(v.GetTags())); err != nil {
		return err
	}

	return nil
}

func resourceOAPISnapshotExportTaskUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	if err := setOSCAPITags(conn, d); err != nil {
		return err
	}

	return resourceOAPISnapshotExportTaskRead(d, meta)
}

func resourceOutscaleSnapshotTaskWaitForAvailable(id string, client *oscgo.APIClient, timeout time.Duration) (oscgo.SnapshotExportTask, error) {
	log.Printf("Waiting for Image Task %s to become available...", id)
	var snap oscgo.SnapshotExportTask
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"pending", "pending/queued", "queued"},
		Target:     []string{"completed", "active"},
		Refresh:    SnapshotTaskStateRefreshFunc(client, id),
		Timeout:    timeout,
		Delay:      OutscaleImageRetryDelay,
		MinTimeout: OutscaleImageRetryMinTimeout,
	}

	info, err := stateConf.WaitForState()
	if err != nil {
		return snap, fmt.Errorf("Error waiting for Snapshot export task (%s) to be ready: %s", id, err)
	}
	snap = info.(oscgo.SnapshotExportTask)
	return snap, nil
}

func resourceOAPISnapshotExportTaskDelete(d *schema.ResourceData, meta interface{}) error {
	d.SetId("")
	if err := d.Set("osu_export", nil); err != nil {
		return err
	}
	return nil
}

// SnapshotTaskStateRefreshFunc ...
func SnapshotTaskStateRefreshFunc(client *oscgo.APIClient, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		var resp oscgo.ReadSnapshotExportTasksResponse
		filter := &oscgo.FiltersExportTask{TaskIds: &[]string{id}}
		var statusCode int
		err := resource.Retry(5*time.Minute, func() *resource.RetryError {
			var err error
			rp, httpResp, err := client.SnapshotApi.ReadSnapshotExportTasks(context.Background()).
				ReadSnapshotExportTasksRequest(oscgo.ReadSnapshotExportTasksRequest{
					Filters: filter,
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
				log.Printf("[INFO] Snapshot export task %s state %s", id, "destroyed")
				return resp, "destroyed", nil
			} else if resp.GetSnapshotExportTasks() != nil && len(resp.GetSnapshotExportTasks()) == 0 {
				log.Printf("[INFO] Snapshot export task %s state %s", id, "destroyed")
				return resp, "destroyed", nil
			} else {
				return resp, "", fmt.Errorf("error on refresh: %+v", err)
			}
		}

		if resp.GetSnapshotExportTasks() == nil || len(resp.GetSnapshotExportTasks()) == 0 {
			return resp, "destroyed", nil
		}

		if resp.GetSnapshotExportTasks()[0].GetState() == "failed" {
			return resp.GetSnapshotExportTasks()[0], resp.GetSnapshotExportTasks()[0].GetState(),
				fmt.Errorf(resp.GetSnapshotExportTasks()[0].GetComment())
		}

		// Snapshot export task is valid, so return it's state
		return resp.GetSnapshotExportTasks()[0], resp.GetSnapshotExportTasks()[0].GetState(), nil
	}
}
