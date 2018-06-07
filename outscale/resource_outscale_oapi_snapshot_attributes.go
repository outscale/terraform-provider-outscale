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

func resourcedOutscaleOAPISnapshotAttributes() *schema.Resource {
	return &schema.Resource{
		Exists: resourcedOutscaleOAPISnapshotAttributesExists,
		Create: resourcedOutscaleOAPISnapshotAttributesCreate,
		Read:   resourcedOutscaleOAPISnapshotAttributesRead,
		Delete: resourcedOutscaleOAPISnapshotAttributesDelete,

		Schema: map[string]*schema.Schema{
			"permission_to_create_volume": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"create": &schema.Schema{
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"global_permission": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
									},
									"account_id": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"delete": &schema.Schema{
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"global_permission": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
									},
									"account_id": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
					},
				},
			},
			"permission_to_create_volumes": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"global_permission": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"account_id": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"snapshot_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"account_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourcedOutscaleOAPISnapshotAttributesExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	conn := meta.(*OutscaleClient).FCU

	sid := d.Get("snapshot_id").(string)
	aid := d.Get("account_id").(string)
	return hasOAPICreateVolumePermission(conn, sid, aid)
}

func resourcedOutscaleOAPISnapshotAttributesCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	sid := d.Get("snapshot_id").(string)
	aid := ""

	req := &fcu.ModifySnapshotAttributeInput{
		SnapshotId: aws.String(sid),
		Attribute:  aws.String("createVolumePermission"),
	}

	if v, ok := d.GetOk("permission_to_create_volume"); ok {
		create := v.([]interface{})[0].(map[string]interface{})["create"].([]interface{})

		a := make([]*fcu.CreateVolumePermission, len(create))

		for k, v1 := range create {
			data := v1.(map[string]interface{})
			a[k] = &fcu.CreateVolumePermission{UserId: aws.String(data["account_id"].(string)), Group: aws.String(data["global_permission"].(string))}
			aid = data["account_id"].(string)
		}
		req.CreateVolumePermission = &fcu.CreateVolumePermissionModifications{Add: a}
	}

	var err error
	err = resource.Retry(2*time.Minute, func() *resource.RetryError {
		_, err = conn.VM.ModifySnapshotAttribute(req)
		if err != nil {
			if strings.Contains(fmt.Sprint(err), "RequestLimitExceeded") {
				log.Printf("[DEBUG] Error: %q", err)
				return resource.RetryableError(err)
			}

			return resource.NonRetryableError(err)
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("Error createing snapshot createVolumePermission: %s", err)
	}

	d.SetId(fmt.Sprintf("%s-%s", sid, aid))
	d.Set("account_id", aid)
	d.Set("permission_to_create_volumes", make([]map[string]interface{}, 0))

	// Wait for the account to appear in the permission list
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"denied"},
		Target:     []string{"granted"},
		Refresh:    resourcedOutscaleOAPISnapshotAttributesStateRefreshFunc(conn, sid, aid),
		Timeout:    5 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}
	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf(
			"Error waiting for snapshot createVolumePermission (%s) to be createed: %s",
			d.Id(), err)
	}

	return nil
}

func resourcedOutscaleOAPISnapshotAttributesRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU
	sid := d.Get("snapshot_id").(string)

	var attrs *fcu.DescribeSnapshotAttributeOutput
	var err error
	err = resource.Retry(2*time.Minute, func() *resource.RetryError {
		attrs, err = conn.VM.DescribeSnapshotAttribute(&fcu.DescribeSnapshotAttributeInput{
			SnapshotId: aws.String(sid),
			Attribute:  aws.String("createVolumePermission"),
		})
		if err != nil {
			if strings.Contains(fmt.Sprint(err), "RequestLimitExceeded") {
				log.Printf("[DEBUG] Error: %q", err)
				return resource.RetryableError(err)
			}

			return resource.NonRetryableError(err)
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("Error refreshing snapshot createVolumePermission state: %s", err)
	}

	cvp := make([]map[string]interface{}, len(attrs.CreateVolumePermissions))
	for k, v := range attrs.CreateVolumePermissions {
		c := make(map[string]interface{})
		c["global_permission"] = aws.StringValue(v.Group)
		c["account_id"] = aws.StringValue(v.UserId)
		cvp[k] = c
	}

	d.Set("request_id", attrs.RequestId)

	return d.Set("permission_to_create_volumes", cvp)
}

func resourcedOutscaleOAPISnapshotAttributesDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	sid := d.Get("snapshot_id").(string)
	v := d.Get("permission_to_create_volume")
	aid := ""

	req := &fcu.ModifySnapshotAttributeInput{
		SnapshotId: aws.String(sid),
		Attribute:  aws.String("createVolumePermission"),
	}

	delete := v.([]interface{})[0].(map[string]interface{})["create"].([]interface{})

	a := make([]*fcu.CreateVolumePermission, len(delete))

	for k, v1 := range delete {
		data := v1.(map[string]interface{})
		a[k] = &fcu.CreateVolumePermission{UserId: aws.String(data["account_id"].(string)), Group: aws.String(data["global_permission"].(string))}
		aid = data["account_id"].(string)
	}
	req.CreateVolumePermission = &fcu.CreateVolumePermissionModifications{Remove: a}

	var err error
	err = resource.Retry(2*time.Minute, func() *resource.RetryError {
		_, err := conn.VM.ModifySnapshotAttribute(req)
		if err != nil {
			if strings.Contains(fmt.Sprint(err), "RequestLimitExceeded") {
				log.Printf("[DEBUG] Error: %q", err)
				return resource.RetryableError(err)
			}

			return resource.NonRetryableError(err)
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("Error removing snapshot createVolumePermission: %s", err)
	}

	// Wait for the account to disappear from the permission list
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"granted"},
		Target:     []string{"denied"},
		Refresh:    resourcedOutscaleOAPISnapshotAttributesStateRefreshFunc(conn, sid, aid),
		Timeout:    5 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}
	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf(
			"Error waiting for snapshot createVolumePermission (%s) to be deleted: %s",
			d.Id(), err)
	}

	return nil
}

func hasOAPICreateVolumePermission(conn *fcu.Client, sid string, aid string) (bool, error) {
	_, state, err := resourcedOutscaleOAPISnapshotAttributesStateRefreshFunc(conn, sid, aid)()
	if err != nil {
		return false, err
	}
	if state == "granted" {
		return true, nil
	}
	return false, nil
}

func resourcedOutscaleOAPISnapshotAttributesStateRefreshFunc(conn *fcu.Client, sid string, aid string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {

		var attrs *fcu.DescribeSnapshotAttributeOutput
		var err error
		err = resource.Retry(2*time.Minute, func() *resource.RetryError {
			attrs, err = conn.VM.DescribeSnapshotAttribute(&fcu.DescribeSnapshotAttributeInput{
				SnapshotId: aws.String(sid),
				Attribute:  aws.String("createVolumePermission"),
			})
			if err != nil {
				if strings.Contains(fmt.Sprint(err), "RequestLimitExceeded") {
					log.Printf("[DEBUG] Error: %q", err)
					return resource.RetryableError(err)
				}

				return resource.NonRetryableError(err)
			}

			return nil
		})

		if err != nil {
			return nil, "", fmt.Errorf("Error refreshing snapshot createVolumePermission state: %s", err)
		}

		for _, vp := range attrs.CreateVolumePermissions {
			if *vp.UserId == aid {
				return attrs, "granted", nil
			}
		}
		return attrs, "denied", nil
	}
}
