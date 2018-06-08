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

func resourcedOutscaleSnapshotImport() *schema.Resource {
	return &schema.Resource{
		Create: resourcedOutscaleSnapshotImportCreate,
		Read:   resourcedOutscaleSnapshotImportRead,
		Delete: resourcedOutscaleSnapshotImportDelete,

		Schema: map[string]*schema.Schema{
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"snapshot_location": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"snapshot_size": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"encrypted": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"owner_alias": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"progress": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": &schema.Schema{
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

func resourcedOutscaleSnapshotImportCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	req := &fcu.ImportSnapshotInput{
		SnapshotLocation: aws.String(d.Get("snapshot_location").(string)),
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

	d.Set("id", resp.ImportTaskId)
	d.SetId(*resp.Id)

	// Wait for the account to appear in the permission list
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"pending"},
		Target:     []string{"completed"},
		Refresh:    resourcedOutscaleSnapshotImportStateRefreshFunc(d, conn, *resp.ImportTaskId),
		Timeout:    5 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}
	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf(
			"Error waiting for snapshot createVolumePermission (%s) to be added: %s",
			d.Id(), err)
	}

	return resourcedOutscaleSnapshotImportRead(d, meta)
}

func resourcedOutscaleSnapshotImportRead(d *schema.ResourceData, meta interface{}) error {
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
	d.Set("encrypted", s.Encrypted)
	d.Set("owner_alias", s.OwnerAlias)
	d.Set("progress", s.Progress)
	d.Set("status", s.State)
	d.Set("volume_size", s.VolumeSize)
	d.Set("request_id", attrs.RequestId)

	return nil
}

func resourcedOutscaleSnapshotImportDelete(d *schema.ResourceData, meta interface{}) error {
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

func resourcedOutscaleSnapshotImportStateRefreshFunc(d *schema.ResourceData, conn *fcu.Client, id string) resource.StateRefreshFunc {
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

		d.Set("progress", s.Progress)

		return attrs, "error", nil
	}
}
