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
					},
				},
				// "osu_bucket": {
				// 	Type:     schema.TypeString,
				// 	Computed: true,
				// },
			},
			"image_id": {
				Type:     schema.TypeString,
				Computed: true,
				Required: true,
			},
			"image_export_task": {
				Type:     schema.TypeList,
				Computed: true,
				"completion": {
					Type:     schema.TypeInt,
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
											Required: true,
										},
										"secret_key": {
											Type:     schema.TypeString,
											Computed: true,
											Required: true,
										},
									},
								},
							},
						},
					},
					// "osu_bucket": {
					// 	Type:     schema.TypeString,
					// 	Computed: true,
					// },
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
			"request_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceImageTasksCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	export_to_osu, export_to_osuOk := d.GetOk("export_to_osu")
	image_id, image_idOk := d.GetOk("image_id")
	image_export_task, image_export_taskOk := d.GetOk("image_export_task")
	request_id, request_idOk := d.GetOk("request_id")

	request := &fcu.CreateImageExportTaskInput{}

	if export_to_osuOk {
		request.export_to_osu = aws.String(export_to_osu.(string))
	}
	if image_idOk {
		request.image_id = aws.String(image_id.(string))
	}
	if image_export_taskOk {
		request.image_export_task = aws.String(image_export_task.(string))
	}
	if request_idOk {
		request.request_id = aws.String(request_id.(string))
	}

	var tasksResp *fcu.CreateImageExportTaskOutput
	var err error

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {

		tasksResp, err = conn.VM.CreateImageExportTask(request)
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

	return nil
}

func resourceImageTasksRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	request := &fcu.DescribeInstanceExportTaskInput{}

	var err error

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {

		_, err = conn.VM.DescribeInstanceExportTask(&fcu.DescribeInstanceExportTaskInput{
			ImageId: aws.String(d.Id()),
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

		return fmt.Errorf("[DEBUG] Error Deregister image %s", err)
	}

	return nil
}

func resourceImageTasksDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}
