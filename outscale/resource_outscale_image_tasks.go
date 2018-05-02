package outscale

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-outscale/osc/fcu"
)

func resourceOutscaleImageTasks() *schema.Resource {
	return &schema.Resource{
		Create: resourceImageTasksCreate,
		Read:   resourceImageTasksRead,
		Delete: resourceImageTasksDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(40 * time.Minute),
			Delete: schema.DefaultTimeout(40 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"export_to_osu": {
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
						"osu_ak_sk": {
							Type:     schema.TypeMap,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"access_key": {
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
						"export_to_osu": {
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
									"osu_ak_sk": {
										Type:     schema.TypeMap,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"access_key": {
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
						"image_export_task_id": {
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
						"status_message": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func resourceImageTasksCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	eto, etoOk := d.GetOk("export_to_osu")
	v, ok := d.GetOk("image_id")
	request := &fcu.CreateImageExportTaskInput{}

	if !etoOk && !ok {
		return fmt.Errorf("Please provide the required attributes export_to_osu and image_id")
	}

	request.ImageId = aws.String(v.(string))

	if etoOk {
		e := eto.(map[string]interface{})
		et := &fcu.ImageExportToOsuTaskSpecification{}
		if v, ok := e["disk_image_format"]; ok {
			et.DiskImageFormat = aws.String(v.(string))
		}
		if v, ok := e["manifest_url"]; ok {
			et.OsuManifestUrl = aws.String(v.(string))
		}
		if v, ok := e["osu_bucket"]; ok {
			et.OsuBucket = aws.String(v.(string))
		}
		if v, ok := e["osu_ak_sk"]; ok {
			w := v.(map[string]interface{})
			et.OsuAkSk = &fcu.ExportToOsuAccessKeySpecification{
				AccessKey: aws.String(w["access_key"].(string)),
				SecretKey: aws.String(w["secret_key"].(string)),
			}
		}
		request.ExportToOsu = et
	}

	var resp *fcu.CreateImageExportTaskOutput
	var err error

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, err = conn.VM.CreateImageExportTask(request)
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

	id := *resp.ImageExportTask.ImageExportTaskId
	d.SetId(id)

	_, err = resourceOutscaleImageTaskWaitForAvailable(id, conn, 1)
	if err != nil {
		return err
	}

	return resourceImageTasksRead(d, meta)
}

func resourceImageTasksRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	var resp *fcu.DescribeImageExportTasksOutput
	var err error

	fmt.Printf("[DEBUG] DESCRIBE IMAGE TASK")

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, err = conn.VM.DescribeImageExportTasks(&fcu.DescribeImageExportTasksInput{
			ImageExportTaskId: []*string{aws.String(d.Id())},
		})
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

	imageExportTask := make([]map[string]interface{}, len(resp.ImageExportTask))
	for k, v := range resp.ImageExportTask {
		i := make(map[string]interface{})
		i["completion"] = *v.Completion
		i["image_export_task_id"] = *v.ImageExportTaskId
		i["image_id"] = *v.ImageId
		i["state"] = *v.State
		i["status_message"] = *v.StatusMessage

		exportToOsu := make(map[string]interface{})
		exportToOsu["disk_image_format"] = *v.ExportToOsu.DiskImageFormat
		exportToOsu["osu_bucket"] = *v.ExportToOsu.OsuBucket
		exportToOsu["manifest_url"] = *v.ExportToOsu.OsuManifestUrl
		exportToOsu["osu_prefix"] = *v.ExportToOsu.OsuPrefix

		osuAkSk := make(map[string]interface{})
		osuAkSk["access_key"] = *v.ExportToOsu.OsuAkSk.AccessKey
		osuAkSk["secret_key"] = *v.ExportToOsu.OsuAkSk.SecretKey

		exportToOsu["osu_ak_sk"] = osuAkSk

		i["export_to_osu"] = exportToOsu

		imageExportTask[k] = i
	}

	if err := d.Set("image_export_task", imageExportTask); err != nil {
		return err
	}

	d.Set("request_id", resp.RequestId)

	return nil
}

func resourceImageTasksDelete(d *schema.ResourceData, meta interface{}) error {

	d.SetId("")

	return nil
}

func resourceOutscaleImageTaskWaitForAvailable(id string, client *fcu.Client, i int) (*fcu.Image, error) {
	fmt.Printf("Waiting for Image Task %s to become available...", id)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"pending", "pending/queued", "queued"},
		Target:     []string{"available"},
		Refresh:    ImageTaskStateRefreshFunc(client, id),
		Timeout:    OutscaleImageRetryTimeout,
		Delay:      OutscaleImageRetryDelay,
		MinTimeout: OutscaleImageRetryMinTimeout,
	}

	info, err := stateConf.WaitForState()
	if err != nil {
		return nil, fmt.Errorf("Error waiting for OMI (%s) to be ready: %v", id, err)
	}
	return info.(*fcu.Image), nil
}

func ImageTaskStateRefreshFunc(client *fcu.Client, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		emptyResp := &fcu.DescribeImageExportTasksOutput{}

		var resp *fcu.DescribeImageExportTasksOutput
		var err error

		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			resp, err = client.VM.DescribeImageExportTasks(&fcu.DescribeImageExportTasksInput{
				ImageExportTaskId: []*string{aws.String(id)},
			})
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
				log.Printf("[INFO] OMI %s state %s", id, "destroyed")
				return emptyResp, "destroyed", nil

			} else if resp != nil && len(resp.ImageExportTask) == 0 {
				log.Printf("[INFO] OMI %s state %s", id, "destroyed")
				return emptyResp, "destroyed", nil
			} else {
				return emptyResp, "", fmt.Errorf("Error on refresh: %+v", err)
			}
		}

		if resp == nil || resp.ImageExportTask == nil || len(resp.ImageExportTask) == 0 {
			return emptyResp, "destroyed", nil
		}

		log.Printf("[INFO] OMI %s state %s", *resp.ImageExportTask[0].ImageId, *resp.ImageExportTask[0].State)

		// OMI is valid, so return it's state
		return resp.ImageExportTask[0], *resp.ImageExportTask[0].State, nil
	}
}
