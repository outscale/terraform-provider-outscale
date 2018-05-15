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

func datasourceOutscaleOAPIVolume() *schema.Resource {
	return &schema.Resource{
		Read: datasourceOAPIVolumeRead,

		Schema: map[string]*schema.Schema{
			// Arguments
			"filter": dataSourceFiltersSchema(),
			"sub_region_name": {
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
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			// Attributes
			"linked_volumes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"delete_on_vm_deletion": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"device_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vm_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"state": {
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
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": {
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
			"tag": tagsSchema(),
			"volume_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func datasourceOAPIVolumeRead(d *schema.ResourceData, meta interface{}) error {

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
		return fmt.Errorf("your query returned no results, please change your search criteria and try again")
	}

	if len(filteredVolumes) > 1 {
		return fmt.Errorf("your query returned more than one result, please try a more " +
			"specific search criteria")
	}

	// Query returned single result.
	volume = filteredVolumes[0]

	log.Printf("[DEBUG] outscale_volume - Single Volume found: %s", *volume.VolumeId)
	return volumeOAPIDescriptionAttributes(d, volume)

}

func volumeOAPIDescriptionAttributes(d *schema.ResourceData, volume *fcu.Volume) error {
	d.SetId(*volume.VolumeId)
	d.Set("volume_id", volume.VolumeId)

	d.Set("sub_region_name", *volume.AvailabilityZone)
	if volume.Size != nil {
		d.Set("size", *volume.Size)
	}
	if volume.SnapshotId != nil {
		d.Set("snapshot_id", *volume.SnapshotId)
	}
	if volume.VolumeType != nil {
		d.Set("type", *volume.VolumeType)
	}

	if volume.VolumeType != nil && *volume.VolumeType == "io1" {
		if volume.Iops != nil {
			d.Set("iops", *volume.Iops)
		}
	}
	if volume.State != nil {
		d.Set("state", *volume.State)
	}
	if volume.VolumeId != nil {
		d.Set("volume_id", *volume.VolumeId)
	}
	if volume.VolumeType != nil {
		d.Set("type", *volume.VolumeType)
	}
	if volume.Attachments != nil {
		res := make([]map[string]interface{}, len(volume.Attachments))
		for k, g := range volume.Attachments {
			r := make(map[string]interface{})
			if g.DeleteOnTermination != nil {
				r["delete_on_vm_deletion"] = *g.DeleteOnTermination
			}
			if g.Device != nil {
				r["device_name"] = *g.Device
			}
			if g.InstanceId != nil {
				r["instance_id"] = *g.InstanceId
			}
			if g.State != nil {
				r["state"] = *g.State
			}
			if g.VolumeId != nil {
				r["volume_id"] = *g.VolumeId
			}

			res[k] = r

		}

		if err := d.Set("linked_volumes", res); err != nil {
			return err
		}
	} else {
		if err := d.Set("linked_volumes", []map[string]interface{}{
			map[string]interface{}{
				"delete_on_vm_deletion": false,
				"device_name":           "none",
				"instance_id":           "none",
				"state":                 "none",
				"volume_id":             "none",
			},
		}); err != nil {
			return err
		}
	}
	if volume.Tags != nil {
		if err := d.Set("tag", dataSourceTags(volume.Tags)); err != nil {
			return err
		}
	} else {
		if err := d.Set("tags", []map[string]string{
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
