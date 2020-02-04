package outscale

import (
	"fmt"
	"log"
	"strings"
	"time"

	oscgo "github.com/marinsalinas/osc-sdk-go"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceOutscaleOAPIImageExportTasks() *schema.Resource {
	return &schema.Resource{
		Create: resourceOAPIImageExportTasksCreate,
		Read:   resourceOAPIImageExportTasksRead,
		Delete: resourceOAPIImageExportTasksDelete,
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
						},
						"osu_bucket": {
							Type:     schema.TypeString,
							Required: true,
						},
						"osu_key": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"osu_prefix": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"osu_api_key": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
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
			"completion": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"snapshot_description": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"snapshot_id": {
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
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"comment": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceOAPIImageExportTasksCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	eto, etoOk := d.GetOk("osu_export")
	v, ok := d.GetOk("snapshot_id")
	request := &oscgo.CreateSnapshotExportTaskInput{}

	if !etoOk && !ok {
		return fmt.Errorf("Please provide the required attributes osu_export and image_id")
	}

	request.SnapshotId = aws.String(v.(string))

	if etoOk {
		exp := eto.([]interface{})
		e := exp[0].(map[string]interface{})

		et := &oscgo.ExportToOsuTaskSpecification{}

		if v, ok := e["disk_image_format"]; ok {
			et.DiskImageFormat = aws.String(v.(string))
		}
		if v, ok := e["osu_key"]; ok {
			et.OsuKey = aws.String(v.(string))
		}
		if v, ok := e["osu_bucket"]; ok {
			et.OsuBucket = aws.String(v.(string))
		}
		if v, ok := e["osu_prefix"]; ok {
			et.OsuPrefix = aws.String(v.(string))
		}
		if v, ok := e["osu_api_key"]; ok {
			a := v.([]interface{})
			if len(a) > 0 {
				w := a[0].(map[string]interface{})
				et.AkSk = &oscgo.ExportToOsuAccessKeySpecification{
					AccessKey: aws.String(w["api_key_id"].(string)),
					SecretKey: aws.String(w["secret_key"].(string)),
				}
			}
		}
		request.ExportToOsu = et
	}

	var resp *oscgo.CreateSnapshotExportTaskOutput
	var err error

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, err = conn.VM.CreateSnapshotExportTask(request)
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

	id := *resp.SnapshotExportTask.SnapshotExportTaskId
	d.SetId(id)

	_, err = resourceOutscaleSnapshotTaskWaitForAvailable(id, conn, 1)
	if err != nil {
		return err
	}

	return resourceOAPIImageExportTasksRead(d, meta)
}

func resourceOAPIImageExportTasksRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	var resp *oscgo.DescribeSnapshotExportTasksOutput
	var err error

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, err = conn.VM.DescribeSnapshotExportTasks(&oscgo.DescribeSnapshotExportTasksInput{
			SnapshotExportTaskId: []*string{aws.String(d.Id())},
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

	v := resp.SnapshotExportTask[0]

	d.Set("completion", v.Completion)
	d.Set("task_id", v.SnapshotExportTaskId)
	d.Set("state", v.State)
	d.Set("completion", v.Completion)
	if v.StatusMessage != nil {
		d.Set("comment", v.StatusMessage)
	} else {
		d.Set("comment", "")
	}

	exp := make([]map[string]interface{}, 1)
	exportToOsu := make(map[string]interface{})
	exportToOsu["disk_image_format"] = *v.ExportToOsu.DiskImageFormat
	exportToOsu["osu_bucket"] = *v.ExportToOsu.OsuBucket
	exportToOsu["osu_key"] = *v.ExportToOsu.OsuKey
	if v.ExportToOsu.OsuPrefix != nil {
		exportToOsu["osu_prefix"] = *v.ExportToOsu.OsuPrefix
	} else {
		exportToOsu["osu_prefix"] = ""
	}

	apk := make([]map[string]interface{}, 1)
	osuAkSk := make(map[string]interface{})
	if v.ExportToOsu.AkSk != nil {
		osuAkSk["api_key_id"] = *v.ExportToOsu.AkSk.AccessKey
		osuAkSk["secret_key"] = *v.ExportToOsu.AkSk.SecretKey
	} else {
		osuAkSk["api_key_id"] = ""
		osuAkSk["secret_key"] = ""
	}
	apk[0] = osuAkSk
	exportToOsu["osu_api_key"] = apk

	snapExp := make(map[string]interface{})
	snapExp["snapshot_id"] = *v.SnapshotExport.SnapshotId

	d.Set("snapshot_description", snapExp)
	exp[0] = exportToOsu
	if err := d.Set("osu_export", exp); err != nil {
		return err
	}
	d.Set("request_id", resp.RequestId)

	return nil
}

func resourceOutscaleSnapshotTaskWaitForAvailable(id string, client *oscgo.Client, i int) (*oscgo.SnapshotExportTask, error) {
	log.Printf("Waiting for Image Task %s to become available...", id)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"pending", "pending/queued", "queued"},
		Target:     []string{"active"},
		Refresh:    SnapshotTaskStateRefreshFunc(client, id),
		Timeout:    OutscaleImageRetryTimeout,
		Delay:      OutscaleImageRetryDelay,
		MinTimeout: OutscaleImageRetryMinTimeout,
	}

	info, err := stateConf.WaitForState()
	if err != nil {
		return nil, fmt.Errorf("Error waiting for OMI (%s) to be ready: %s", id, err)
	}
	return info.(*oscgo.SnapshotExportTask), nil
}

func resourceOAPIImageExportTasksDelete(d *schema.ResourceData, meta interface{}) error {

	d.SetId("")
	d.Set("snapshot_description", nil)
	d.Set("osu_export", nil)
	d.Set("request_id", nil)

	return nil
}

// SnapshotTaskStateRefreshFunc ...
func SnapshotTaskStateRefreshFunc(client *oscgo.Client, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		emptyResp := &oscgo.DescribeSnapshotExportTasksOutput{}

		var resp *oscgo.DescribeSnapshotExportTasksOutput
		var err error

		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			resp, err = client.VM.DescribeSnapshotExportTasks(&oscgo.DescribeSnapshotExportTasksInput{
				SnapshotExportTaskId: []*string{aws.String(id)},
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

			} else if resp != nil && len(resp.SnapshotExportTask) == 0 {
				log.Printf("[INFO] OMI %s state %s", id, "destroyed")
				return emptyResp, "destroyed", nil
			} else {
				return emptyResp, "", fmt.Errorf("Error on refresh: %+v", err)
			}
		}

		if resp == nil || resp.SnapshotExportTask == nil || len(resp.SnapshotExportTask) == 0 {
			return emptyResp, "destroyed", nil
		}

		if *resp.SnapshotExportTask[0].State == "failed" {
			return resp.SnapshotExportTask[0], *resp.SnapshotExportTask[0].State, fmt.Errorf(*resp.SnapshotExportTask[0].StatusMessage)
		}

		// OMI is valid, so return it's state
		return resp.SnapshotExportTask[0], *resp.SnapshotExportTask[0].State, nil
	}
}
