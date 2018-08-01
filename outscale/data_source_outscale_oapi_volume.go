package outscale

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-outscale/osc/oapi"
	"github.com/terraform-providers/terraform-provider-outscale/utils"
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

	conn := meta.(*OutscaleClient).OAPI

	filters, filtersOk := d.GetOk("filter")
	VolumeIds, VolumeIdsOk := d.GetOk("volume_id")

	params := &oapi.ReadVolumesRequest{
		Filters: oapi.ReadVolumesFilters{},
	}
	if filtersOk {
		params.Filters = buildOutscaleOAPIDataSourceVolumesFilters(filters.(*schema.Set))
	}
	if VolumeIdsOk {
		params.Filters.VolumeIds = []string{VolumeIds.(string)}
	}

	var resp *oapi.ReadVolumesResponse
	var rs *oapi.POST_ReadVolumesResponses
	var err error

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		rs, err = conn.POST_ReadVolumes(*params)
		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded:") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return resource.NonRetryableError(err)
	})

	resp = rs.OK

	if err != nil {
		return err
	}

	log.Printf("Found These Volumes %s", spew.Sdump(resp.Volumes))

	filteredVolumes := resp.Volumes[:]

	var volume *oapi.Volumes
	if len(filteredVolumes) < 1 {
		return fmt.Errorf("your query returned no results, please change your search criteria and try again")
	}

	if len(filteredVolumes) > 1 {
		return fmt.Errorf("your query returned more than one result, please try a more " +
			"specific search criteria")
	}

	// Query returned single result.
	volume = &filteredVolumes[0]

	log.Printf("[DEBUG] outscale_volume - Single Volume found: %s", volume.VolumeId)
	return volumeOAPIDescriptionAttributes(d, volume)

}

func volumeOAPIDescriptionAttributes(d *schema.ResourceData, volume *oapi.Volumes) error {
	d.SetId(volume.VolumeId)
	d.Set("volume_id", volume.VolumeId)

	d.Set("sub_region_name", volume.SubRegionName)
	//if volume.Size != "" {
	d.Set("size", volume.Size)
	//}
	if volume.SnapshotId != "" {
		d.Set("snapshot_id", volume.SnapshotId)
	}
	if volume.Type != "" {
		d.Set("type", volume.Type)
	}

	if volume.Type != "" && volume.Type == "io1" {
		//if volume.Iops != "" {
		d.Set("iops", volume.Iops)
		//}
	}
	if volume.State != "" {
		d.Set("state", volume.State)
	}
	if volume.VolumeId != "" {
		d.Set("volume_id", volume.VolumeId)
	}

	if volume.LinkedVolumes != nil {
		res := make([]map[string]interface{}, len(volume.LinkedVolumes))
		for k, g := range volume.LinkedVolumes {
			r := make(map[string]interface{})
			//if g.DeleteOnVmDeletion != "" {
			r["delete_on_vm_termination"] = g.DeleteOnVmDeletion
			//}
			if g.DeviceName != "" {
				r["device"] = g.DeviceName
			}
			if g.VmId != "" {
				r["vm_id"] = g.VmId
			}
			if g.State != "" {
				r["state"] = g.State
			}
			if g.VolumeId != "" {
				r["volume_id"] = g.VolumeId
			}

			res[k] = r

		}

		if err := d.Set("linked_volumes", res); err != nil {
			return err
		}
	} else {
		if err := d.Set("linked_volumes", []map[string]interface{}{
			map[string]interface{}{
				"delete_on_vm_termination": false,
				"device":                   "none",
				"vm_id":                    "none",
				"state":                    "none",
				"volume_id":                "none",
			},
		}); err != nil {
			return err
		}
	}
	// if volume.Tags != nil {
	// 	if err := d.Set("tags", tagsToMap(volume.Tags)); err != nil {
	// 		return err
	// 	}
	// } else {
	// 	if err := d.Set("tags", []map[string]string{
	// 		map[string]string{
	// 			"key":   "",
	// 			"value": "",
	// 		},
	// 	}); err != nil {
	// 		return err
	// 	}
	// }

	return nil
}

func buildOutscaleOAPIDataSourceVolumesFilters(set *schema.Set) oapi.ReadVolumesFilters {
	var filters oapi.ReadVolumesFilters
	for _, v := range set.List() {
		m := v.(map[string]interface{})
		var filterValues []string
		for _, e := range m["values"].([]interface{}) {
			filterValues = append(filterValues, e.(string))
		}

		switch name := m["name"].(string); name {
		case "creation-date":
			filters.CreationDates = filterValues
		case "snapshot-id":
			filters.SnapshotIds = filterValues
		case "sub-region-name":
			filters.SubRegionNames = filterValues
		case "tag-key":
			filters.TagKeys = filterValues
		// case "tags":
		// 	filters.Tags = filterValues
		case "tag-value":
			filters.TagValues = filterValues
		case "volume-id":
			filters.VolumeIds = filterValues
		case "volume-size":
			filters.VolumeSizes = utils.StringSliceToInt64Slice(filterValues)
		case "volume-type":
			filters.VolumeTypes = filterValues
		default:
			log.Printf("[Debug] Unknown Filter Name: %s.", name)
		}
	}
	return filters
}
