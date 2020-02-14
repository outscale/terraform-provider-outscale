package outscale

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/antihax/optional"
	oscgo "github.com/marinsalinas/osc-sdk-go"
	"github.com/terraform-providers/terraform-provider-outscale/utils"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourcedOutscaleOAPISnapshotCopy() *schema.Resource {
	return &schema.Resource{
		Create: resourcedOutscaleOAPISnapshotCopyCreate,
		Read:   resourcedOutscaleOAPISnapshotCopyRead,
		Delete: resourcedOutscaleOAPISnapshotCopyDelete,

		Schema: map[string]*schema.Schema{
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"source_region_name": &schema.Schema{
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

func resourcedOutscaleOAPISnapshotCopyCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	req := oscgo.CreateSnapshotRequest{
		SourceRegionName: oscgo.PtrString(d.Get("source_region_name").(string)),
		SourceSnapshotId: oscgo.PtrString(d.Get("source_snapshot_id").(string)),
	}

	if v, ok := d.GetOk("description"); ok {
		req.Description = oscgo.PtrString(v.(string))
	}

	var o oscgo.CreateSnapshotResponse
	var err error
	err = resource.Retry(2*time.Minute, func() *resource.RetryError {
		o, _, err = conn.SnapshotApi.CreateSnapshot(context.Background(), &oscgo.CreateSnapshotOpts{
			CreateSnapshotRequest: optional.NewInterface(req),
		})
		if err != nil {
			if strings.Contains(fmt.Sprint(err), "Throttling") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("Error copying snapshot: %s", utils.GetErrorResponse(err))
	}

	d.SetId(resource.UniqueId())
	d.Set("snapshot_id", o.Snapshot.SnapshotId)
	d.Set("request_id", o.ResponseContext.RequestId)

	return nil
}

func resourcedOutscaleOAPISnapshotCopyRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourcedOutscaleOAPISnapshotCopyDelete(d *schema.ResourceData, meta interface{}) error {
	d.SetId("")

	return nil
}
