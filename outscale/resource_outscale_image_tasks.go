package outscale

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/antihax/optional"
	oscgo "github.com/marinsalinas/osc-sdk-go"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceOutscaleOAPIImageTasks() *schema.Resource {
	return &schema.Resource{
		Create: resourceOAPIImageTasksCreate,
		Read:   resourceOAPIImageTasksRead,
		Delete: resourceOAPIImageTasksDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(40 * time.Minute),
			Delete: schema.DefaultTimeout(40 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"osu_export": {
				Type:     schema.TypeMap,
				Required: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"disk_image_format": {
							Type:     schema.TypeString,
							Required: true,
						},
						"manifest_url": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"osu_api_key": {
							Type:     schema.TypeMap,
							Optional: true,
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
						"osu_bucket": {
							Type:     schema.TypeString,
							Optional: true,
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
			"image_export_task": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"completion": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"osu_export": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"disk_image_format": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"manifest_url": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"osu_api_key": {
										Type:     schema.TypeMap,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"api_key_id": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"secret_key": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
									"osu_bucket": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"image_export": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"image_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"task_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"image_id": {
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
					},
				},
			},
		},
	}
}

func resourceOAPIImageTasksCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	eto, etoOk := d.GetOk("osu_export")
	v, ok := d.GetOk("image_id")
	request := &oscgo.CreateImageExportTaskRequest{}

	if !etoOk && !ok {
		return fmt.Errorf("Please provide the required attributes osu_export and image_id")
	}

	request.ImageId = v.(string)

	if etoOk {
		e := eto.(map[string]interface{})
		et := oscgo.OsuExport{}
		if v, ok := e["disk_image_format"]; ok {
			et.DiskImageFormat = v.(string)
		}
		if v, ok := e["manifest_url"]; ok {
			et.SetOsuManifestUrl(v.(string))
		}
		if v, ok := e["osu_bucket"]; ok {
			et.OsuBucket = v.(string)
		}
		if v, ok := e["osu_api_key"]; ok {
			w := v.(map[string]interface{})
			oak := oscgo.OsuApiKey{}
			oak.SetApiKeyId(w["api_key_id"].(string))
			oak.SetSecretKey(w["secret_key"].(string))
			et.OsuApiKey = &oak
		}
		request.OsuExport = et
	}

	var resp oscgo.CreateImageExportTaskResponse
	var err error

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, _, err = conn.ImageApi.CreateImageExportTask(
			context.Background(),
			&oscgo.CreateImageExportTaskOpts{
				CreateImageExportTaskRequest: optional.NewInterface(request)})
		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded:") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("[DEBUG] Error image task %s", err)
	}

	ID := *resp.ImageExportTask.ImageId
	d.SetId(ID)

	_, err = resourceOutscaleImageTaskWaitForAvailable(ID, conn, 1)
	if err != nil {
		return err
	}

	return resourceOAPIImageTasksRead(d, meta)
}

func resourceOutscaleImageTaskWaitForAvailable(ID string, client *oscgo.APIClient, i int) (*oscgo.Image, error) {
	fmt.Printf("Waiting for Image Task %s to become available...", ID)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"pending", "pending/queued", "queued"},
		Target:     []string{"available"},
		Refresh:    OAPIImageTaskStateRefreshFunc(client, ID),
		Timeout:    OutscaleImageRetryTimeout,
		Delay:      OutscaleImageRetryDelay,
		MinTimeout: OutscaleImageRetryMinTimeout,
	}

	info, err := stateConf.WaitForState()
	if err != nil {
		return nil, fmt.Errorf("Error waiting for OMI (%s) to be ready: %v", ID, err)
	}
	return info.(*oscgo.Image), nil
}

func resourceOAPIImageTasksRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	var resp oscgo.ReadImageExportTasksResponse
	var err error

	log.Printf("[DEBUG] DESCRIBE IMAGE TASK")

	tids := []string{d.Id()}
	filter := oscgo.FiltersExportTask{
		TaskIds: &tids,
	}
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, _, err = conn.ImageApi.ReadImageExportTasks(
			context.Background(),
			&oscgo.ReadImageExportTasksOpts{
				ReadImageExportTasksRequest: optional.NewInterface(
					&oscgo.ReadImageExportTasksRequest{
						Filters: &filter,
					},
				),
			},
		)
		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded:") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("Error reading task image %s", err)
	}

	imageExportTask := make([]map[string]interface{}, len(*resp.ImageExportTasks))
	for k, v := range *resp.ImageExportTasks {
		i := make(map[string]interface{})
		i["completion"] = *v.Progress
		i["task_id"] = *v.TaskId
		i["image_id"] = *v.ImageId
		i["state"] = *v.State
		i["comment"] = *v.Comment

		exportToOsu := make(map[string]interface{})
		exportToOsu["disk_image_format"] = v.OsuExport.DiskImageFormat
		exportToOsu["osu_bucket"] = v.OsuExport.OsuBucket
		exportToOsu["manifest_url"] = *v.OsuExport.OsuManifestUrl
		exportToOsu["osu_prefix"] = *v.OsuExport.OsuPrefix

		osuAkSk := make(map[string]interface{})
		osuAkSk["api_key_id"] = *v.OsuExport.OsuApiKey.ApiKeyId
		osuAkSk["secret_key"] = *v.OsuExport.OsuApiKey.SecretKey

		exportToOsu["osu_api_key"] = osuAkSk

		i["osu_export"] = exportToOsu

		imageExportTask[k] = i
	}

	if err := d.Set("image_export_task", imageExportTask); err != nil {
		return err
	}

	d.Set("request_id", resp.ResponseContext.RequestId)

	return nil
}

func resourceOAPIImageTasksDelete(d *schema.ResourceData, meta interface{}) error {

	d.SetId("")

	return nil
}

// OAPIImageTaskStateRefreshFunc ...
func OAPIImageTaskStateRefreshFunc(client *oscgo.APIClient, ID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		emptyResp := &oscgo.ReadImageExportTasksResponse{}

		var resp oscgo.ReadImageExportTasksResponse
		var err error
		tids := []string{ID}
		filter := oscgo.FiltersExportTask{
			TaskIds: &tids,
		}

		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			resp, _, err = client.ImageApi.ReadImageExportTasks(
				context.Background(),
				&oscgo.ReadImageExportTasksOpts{
					ReadImageExportTasksRequest: optional.NewInterface(
						&oscgo.ReadImageExportTasksRequest{
							Filters: &filter,
						},
					),
				},
			)
			if err != nil {
				if strings.Contains(err.Error(), "RequestLimitExceeded:") {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})

		if err != nil {
			if e := fmt.Sprint(err); strings.Contains(e, "InvalidAMIID.NotFound") {
				log.Printf("[INFO] OMI %s state %s", ID, "destroyed")
				return emptyResp, "destroyed", nil

			} else if resp.ImageExportTasks == nil || len(*resp.ImageExportTasks) == 0 {
				log.Printf("[INFO] OMI %s state %s", ID, "destroyed")
				return emptyResp, "destroyed", nil
			} else {
				return emptyResp, "", fmt.Errorf("Error on refresh: %+v", err)
			}
		}

		if resp.ImageExportTasks == nil || len(*resp.ImageExportTasks) == 0 {
			return emptyResp, "destroyed", nil
		}

		log.Printf("[INFO] OMI %s state %s", (*resp.ImageExportTasks)[0].ImageId,
			(*resp.ImageExportTasks)[0].State)

		// OMI is valid, so return it's state
		return (*resp.ImageExportTasks)[0], *(*resp.ImageExportTasks)[0].State, nil
	}
}
