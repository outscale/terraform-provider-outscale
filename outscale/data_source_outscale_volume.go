package outscale

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/davecgh/go-spew/spew"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-outscale/osc/fcu"
)

func datasourceOutscaleVolume() *schema.Resource {
	return &schema.Resource{
		Read: datasourceVolumeRead,

		Schema: map[string]*schema.Schema{
			// Arguments
			"filter": dataSourceFiltersSchema(),
			"availability_zone": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"iops": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"size": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"snapshot_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"volume_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			// Attributes
			"attachment_set": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"delete_on_termination": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"device": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"volume_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tag_set": {
				Type: schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"value": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
				Computed: true,
			},
			"tags": tagsSchema(),
			"volume_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"request_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func datasourceVolumeRead(d *schema.ResourceData, meta interface{}) error {

	conn := meta.(*OutscaleClient).FCU

	filters, filtersOk := d.GetOk("filter")
	VolumeIds, VolumeIdsOk := d.GetOk("volume_id")

	params := &fcu.DescribeVolumesInput{}
	if filtersOk {
		params.Filters = buildOutscaleDataSourceFilters(filters.(*schema.Set))
	}
	if VolumeIdsOk {
		params.VolumeIds = []*string{aws.String(VolumeIds.(string))}
	}

	var resp *fcu.DescribeVolumesOutput
	var err error

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, err = conn.VM.DescribeVolumes(params)
		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded:") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return resource.NonRetryableError(err)
	})

	if err != nil {
		return err
	}

	log.Printf("Found These Volumes %s", spew.Sdump(resp.Volumes))

	filteredVolumes := resp.Volumes[:]

	var volume *fcu.Volume
	if len(filteredVolumes) < 1 {
		return fmt.Errorf("Your query returned no results. Please change your search criteria and try again.")
	}

	if len(filteredVolumes) > 1 {
		return fmt.Errorf("Your query returned more than one result. Please try a more " +
			"specific search criteria.")
	} else {
		// Query returned single result.
		volume = filteredVolumes[0]
	}

	d.Set("request_id", resp.RequestId)

	log.Printf("[DEBUG] outscale_volume - Single Volume found: %s", *volume.VolumeId)
	return volumeDescriptionAttributes(d, volume)

}

func volumeDescriptionAttributes(d *schema.ResourceData, volume *fcu.Volume) error {
	d.SetId(*volume.VolumeId)
	d.Set("volume_id", volume.VolumeId)

	d.Set("availability_zone", *volume.AvailabilityZone)
	if volume.Size != nil {
		d.Set("size", *volume.Size)
	}
	if volume.SnapshotId != nil {
		d.Set("snapshot_id", *volume.SnapshotId)
	}
	if volume.VolumeType != nil {
		d.Set("volume_type", *volume.VolumeType)
	}

	if volume.VolumeType != nil && *volume.VolumeType == "io1" {
		if volume.Iops != nil {
			d.Set("iops", *volume.Iops)
		}
	}
	if volume.State != nil {
		d.Set("status", *volume.State)
	}
	if volume.VolumeId != nil {
		d.Set("volume_id", *volume.VolumeId)
	}
	if volume.VolumeType != nil {
		d.Set("volume_type", *volume.VolumeType)
	}
	if volume.Attachments != nil {
		res := make([]map[string]interface{}, len(volume.Attachments))
		for k, g := range volume.Attachments {
			r := make(map[string]interface{})
			if g.DeleteOnTermination != nil {
				r["delete_on_termination"] = *g.DeleteOnTermination
			}
			if g.Device != nil {
				r["device"] = *g.Device
			}
			if g.InstanceId != nil {
				r["instance_id"] = *g.InstanceId
			}
			if g.State != nil {
				r["status"] = *g.State
			}
			if g.VolumeId != nil {
				r["volume_id"] = *g.VolumeId
			}

			res[k] = r

		}

		if err := d.Set("attachment_set", res); err != nil {
			return err
		}
	} else {
		if err := d.Set("attachment_set", []map[string]interface{}{
			map[string]interface{}{
				"delete_on_termination": false,
				"device":                "none",
				"instance_id":           "none",
				"status":                "none",
				"volume_id":             "none",
			},
		}); err != nil {
			return err
		}
	}
	if volume.Tags != nil {
		if err := d.Set("tags", dataSourceTags(volume.Tags)); err != nil {
			return err
		}
	} else {
		if err := d.Set("tag_set", []map[string]string{
			map[string]string{
				"key":   "",
				"value": "",
			},
		}); err != nil {
			return err
		}
	}

	return nil
}
