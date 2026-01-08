package oapi

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/utils"

	oscgo "github.com/outscale/osc-sdk-go/v2"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceOutscaleImageExportTask() *schema.Resource {
	return &schema.Resource{
		Create: resourceOAPIImageExportTaskCreate,
		Read:   resourceOAPIImageExportTaskRead,
		Update: resourceOAPIImageExportTaskUpdate,
		Delete: resourceOAPIImageExportTaskDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
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
			"tags": TagsSchemaSDK(),
		},
	}
}

func resourceOAPIImageExportTaskCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*client.OutscaleClient).OSCAPI

	eto, etoOk := d.GetOk("osu_export")
	v, ok := d.GetOk("image_id")
	request := oscgo.CreateImageExportTaskRequest{}

	if !etoOk && !ok {
		return fmt.Errorf("please provide the required attributes osu_export and image_id")
	}

	request.SetImageId(v.(string))

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
		if v, ok := e["osu_manifest_url"]; ok {
			et.SetOsuManifestUrl(v.(string))
		}
		request.SetOsuExport(et)
	}

	var resp oscgo.CreateImageExportTaskResponse
	var err error

	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		rp, httpResp, err := conn.ImageApi.CreateImageExportTask(context.Background()).
			CreateImageExportTaskRequest(request).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		resp = rp
		return nil
	})
	if err != nil {
		return fmt.Errorf("error image task %s", err)
	}

	id := resp.ImageExportTask.GetTaskId()
	d.SetId(id)
	if d.IsNewResource() {
		if err := updateOAPITagsSDK(conn, d); err != nil {
			return err
		}
	}
	_, err = ResourceOutscaleImageTaskWaitForAvailable(id, conn, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return err
	}

	return resourceOAPIImageExportTaskRead(d, meta)
}

func resourceOAPIImageExportTaskRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*client.OutscaleClient).OSCAPI

	var resp oscgo.ReadImageExportTasksResponse
	var err error
	filter := &oscgo.FiltersExportTask{TaskIds: &[]string{d.Id()}}
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		rp, httpResp, err := conn.ImageApi.ReadImageExportTasks(context.Background()).
			ReadImageExportTasksRequest(oscgo.ReadImageExportTasksRequest{
				Filters: filter,
			}).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		resp = rp
		return nil
	})
	if err != nil {
		return fmt.Errorf("error reading task image %w", err)
	}
	if utils.IsResponseEmpty(len(resp.GetImageExportTasks()), "ImageExportTask", d.Id()) {
		d.SetId("")
		return nil
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

	if err = d.Set("image_id", v.GetImageId()); err != nil {
		return err
	}
	if err = d.Set("osu_export", exp); err != nil {
		return err
	}
	if err = d.Set("tags", FlattenOAPITagsSDK(v.GetTags())); err != nil {
		return err
	}

	return nil
}

func resourceOAPIImageExportTaskUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*client.OutscaleClient).OSCAPI

	if err := updateOAPITagsSDK(conn, d); err != nil {
		return err
	}
	return resourceOAPIImageExportTaskRead(d, meta)
}

func ResourceOutscaleImageTaskWaitForAvailable(id string, client *oscgo.APIClient, timeout time.Duration) (oscgo.ImageExportTask, error) {
	log.Printf("Waiting for Image Task %s to become available...", id)
	var image oscgo.ImageExportTask
	stateConf := &retry.StateChangeConf{
		Pending:    []string{"pending", "pending/queued", "queued"},
		Target:     []string{"completed"},
		Refresh:    ImageTaskStateRefreshFunc(client, id),
		Timeout:    timeout,
		Delay:      OutscaleImageRetryDelay,
		MinTimeout: OutscaleImageRetryMinTimeout,
	}

	info, err := stateConf.WaitForState()
	if err != nil {
		return image, fmt.Errorf("error waiting for image export task (%s) to be ready: %s", id, err)
	}
	image = info.(oscgo.ImageExportTask)
	return image, nil
}

func resourceOAPIImageExportTaskDelete(d *schema.ResourceData, meta interface{}) error {
	d.SetId("")
	if err := d.Set("osu_export", nil); err != nil {
		return err
	}

	return nil
}

// ImageTaskStateRefreshFunc ...
func ImageTaskStateRefreshFunc(client *oscgo.APIClient, id string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		var resp oscgo.ReadImageExportTasksResponse
		var err error
		var statusCode int

		err = retry.Retry(5*time.Minute, func() *retry.RetryError {
			filter := &oscgo.FiltersExportTask{TaskIds: &[]string{id}}
			rp, httpResp, err := client.ImageApi.ReadImageExportTasks(context.Background()).
				ReadImageExportTasksRequest(oscgo.ReadImageExportTasksRequest{
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
				log.Printf("[INFO] Image export task %s state %s", id, "destroyed")
				return resp, "destroyed", nil
			} else if resp.GetImageExportTasks() != nil && len(resp.GetImageExportTasks()) == 0 {
				log.Printf("[INFO] Image export task %s state %s", id, "destroyed")
				return resp, "destroyed", nil
			} else {
				return resp, "", fmt.Errorf("error on refresh: %w", err)
			}
		}

		if resp.GetImageExportTasks() == nil || len(resp.GetImageExportTasks()) == 0 {
			return resp, "destroyed", nil
		}

		if resp.GetImageExportTasks()[0].GetState() == "failed" || resp.GetImageExportTasks()[0].GetState() == "cancelled" {
			return resp.GetImageExportTasks()[0], resp.GetImageExportTasks()[0].GetState(), fmt.Errorf("error: %v", resp.GetImageExportTasks()[0].GetComment())
		}

		// Image export task is valid, so return it's state
		return resp.GetImageExportTasks()[0], resp.GetImageExportTasks()[0].GetState(), nil
	}
}
