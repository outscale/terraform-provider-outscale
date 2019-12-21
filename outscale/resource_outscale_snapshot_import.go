package outscale

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-outscale/osc/fcu"
)

func resourcedOutscaleOAPISnapshotImport() *schema.Resource {
	return &schema.Resource{
		Create: resourcedOutscaleOAPISnapshotImportCreate,
		Read:   resourcedOutscaleOAPISnapshotImportRead,
		Delete: resourcedOutscaleOAPISnapshotImportDelete,

		Schema: map[string]*schema.Schema{
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"osu_location": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"snapshot_size": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"is_encrypted": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"vm_profile_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"account_alias": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"completion": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"state": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"volume_size": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"request_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourcedOutscaleOAPISnapshotImportCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	req := &fcu.ImportSnapshotInput{
		SnapshotLocation: aws.String(d.Get("osu_location").(string)),
		SnapshotSize:     aws.String(d.Get("snapshot_size").(string)),
	}

	if v, ok := d.GetOk("description"); ok {
		req.Description = aws.String(v.(string))
	}

	var resp *fcu.ImportSnapshotOutput
	var err error
	err = resource.Retry(2*time.Minute, func() *resource.RetryError {
		resp, err = conn.VM.ImportSnapshot(req)
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
		return fmt.Errorf("Error adding snapshot createVolumePermission: %s", err)
	}

	d.Set("vm_profile_id", resp.ImportTaskId)
	d.SetId(*resp.Id)

	// Wait for the account to appear in the permission list
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"pending"},
		Target:     []string{"completed"},
		Refresh:    resourcedOutscaleOAPISnapshotImportStateRefreshFunc(d, conn, *resp.ImportTaskId),
		Timeout:    5 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}
	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf(
			"Error waiting for snapshot createVolumePermission (%s) to be added: %s",
			d.Id(), err)
	}

	return resourcedOutscaleOAPISnapshotImportRead(d, meta)
}

func resourcedOutscaleOAPISnapshotImportRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	var attrs *fcu.DescribeSnapshotsOutput
	var err error
	err = resource.Retry(2*time.Minute, func() *resource.RetryError {
		attrs, err = conn.VM.DescribeSnapshots(&fcu.DescribeSnapshotsInput{
			SnapshotIds: []*string{aws.String(d.Id())},
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
		return fmt.Errorf("Error refreshing snapshot state: %s", err)
	}

	s := attrs.Snapshots[0]

	d.Set("description", s.Description)
	d.Set("is_encrypted", s.Encrypted)
	d.Set("account_alias", s.OwnerAlias)
	d.Set("completion", s.Progress)
	d.Set("state", s.State)
	d.Set("volume_size", s.VolumeSize)
	d.Set("request_id", attrs.RequestId)

	return nil
}

func resourcedOutscaleOAPISnapshotImportDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		request := &fcu.DeleteSnapshotInput{
			SnapshotId: aws.String(d.Id()),
		}
		var err error
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			_, err := conn.VM.DeleteSnapshot(request)

			if err != nil {
				if strings.Contains(err.Error(), "RequestLimitExceeded") {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}

			return nil
		})
		if err == nil {
			return nil
		}

		ebsErr, ok := err.(awserr.Error)
		if ebsErr.Code() == "SnapshotInUse" {
			return resource.RetryableError(fmt.Errorf("EBS SnapshotInUse - trying again while it detaches"))
		}

		if !ok {
			return resource.NonRetryableError(err)
		}

		return nil
	})
}

func resourcedOutscaleOAPISnapshotImportStateRefreshFunc(d *schema.ResourceData, conn *fcu.Client, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {

		var attrs *fcu.DescribeSnapshotsOutput
		var err error
		err = resource.Retry(2*time.Minute, func() *resource.RetryError {
			attrs, err = conn.VM.DescribeSnapshots(&fcu.DescribeSnapshotsInput{
				SnapshotIds: []*string{aws.String(id)},
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
			return nil, "", fmt.Errorf("Error refreshing snapshot state: %s", err)
		}

		s := attrs.Snapshots[0]

		d.Set("completion", s.Progress)

		return attrs, "error", nil
	}
}
