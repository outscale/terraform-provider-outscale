package outscale

import (
	"fmt"
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
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"disk_image_format": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"manifest_url": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"osu_ak_sk": {
							Type:     schema.TypeList,
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
							Type:     schema.TypeList,
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
										Type:     schema.TypeList,
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
							Type:     schema.TypeList,
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

	request := &fcu.CreateImageExportTaskInput{}

	if v, ok := d.GetOk("image_id"); ok {
		request.ImageId = aws.String(v.(string))
	}

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

	d.SetId(*resp.ImageExportTask.ImageExportTaskId)

	return resourceImageTasksRead(d, meta)
}

func resourceImageTasksRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	var resp *fcu.DescribeImageExportTasksOutput
	var err error

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

		i["exportToOsu"] = exportToOsu

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
