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

func resourcedOutscaleSnapshotCopy() *schema.Resource {
	return &schema.Resource{
		Create: resourcedOutscaleSnapshotCopyCreate,
		Read:   resourcedOutscaleSnapshotCopyRead,
		Delete: resourcedOutscaleSnapshotCopyDelete,

		Schema: map[string]*schema.Schema{
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"destination_region": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"source_region": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"snapshot_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"request_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"source_snapshot_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourcedOutscaleSnapshotCopyCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	req := &fcu.CopySnapshotInput{
		SourceRegion:     aws.String(d.Get("source_region").(string)),
		SourceSnapshotId: aws.String(d.Get("source_snapshot_id").(string)),
	}

	if v, ok := d.GetOk("description"); ok {
		req.Description = aws.String(v.(string))
	}
	if v, ok := d.GetOk("destination_region"); ok {
		req.DestinationRegion = aws.String(v.(string))
	}

	var o *fcu.CopySnapshotOutput
	var err error
	err = resource.Retry(2*time.Minute, func() *resource.RetryError {
		o, err = conn.VM.CopySnapshot(req)
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
		return fmt.Errorf("Error copying snapshot: %s", err)
	}

	d.SetId(resource.UniqueId())
	d.Set("snapshot_id", aws.StringValue(o.SnapshotId))
	d.Set("request_id", aws.StringValue(o.RequestId))

	return nil
}

func resourcedOutscaleSnapshotCopyRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourcedOutscaleSnapshotCopyDelete(d *schema.ResourceData, meta interface{}) error {
	d.SetId("")

	return nil
}
