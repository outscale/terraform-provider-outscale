package outscale

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/outscale/osc-go/oapi"
	"github.com/terraform-providers/terraform-provider-outscale/utils"
)

func datasourceOutscaleOAPIVolume() *schema.Resource {
	return &schema.Resource{
		Read: datasourceOAPIVolumeRead,

		Schema: map[string]*schema.Schema{
			// Arguments
			"filter": dataSourceFiltersSchema(),
			"subregion_name": {
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
			"tags": tagsOAPIListSchemaComputed(),
			"volume_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"request_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func datasourceOAPIVolumeRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OAPI

	filters, filtersOk := d.GetOk("filter")
	volumeIds, VolumeIdsOk := d.GetOk("volume_id")

	params := oapi.ReadVolumesRequest{
		Filters: oapi.FiltersVolume{},
	}
	if VolumeIdsOk {
		params.Filters.VolumeIds = []string{volumeIds.(string)}
	}

	if filtersOk {
		params.Filters = buildOutscaleOAPIDataSourceVolumesFilters(filters.(*schema.Set))
	}

	var resp *oapi.ReadVolumesResponse
	var rs *oapi.POST_ReadVolumesResponses
	var err error

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		rs, err = conn.POST_ReadVolumes(params)
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

	var volume oapi.Volume
	if len(filteredVolumes) < 1 {
		return fmt.Errorf("your query returned no results, please change your search criteria and try again")
	}

	if len(filteredVolumes) > 1 {
		return fmt.Errorf("your query returned more than one result, please try a more " +
			"specific search criteria")
	}

	// Query returned single result.
	volume = filteredVolumes[0]
	d.Set("request_id", resp.ResponseContext.RequestId)
	log.Printf("[DEBUG] outscale_volume - Single Volume found: %s", volume.VolumeId)
	return volumeOAPIDescriptionAttributes(d, &volume)

}

func volumeOAPIDescriptionAttributes(d *schema.ResourceData, volume *oapi.Volume) error {
	d.SetId(volume.VolumeId)
	d.Set("volume_id", volume.VolumeId)
	d.Set("subregion_name", volume.SubregionName)
	d.Set("size", volume.Size)
	d.Set("snapshot_id", volume.SnapshotId)
	d.Set("volume_type", volume.VolumeType)

	d.Set("state", volume.State)
	d.Set("volume_id", volume.VolumeId)
	d.Set("iops", volume.Iops)

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

	if volume.Tags != nil {
		if err := d.Set("tags", tagsOAPIToMap(volume.Tags)); err != nil {
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

func buildOutscaleOAPIDataSourceVolumesFilters(set *schema.Set) oapi.FiltersVolume {
	var filters oapi.FiltersVolume
	for _, v := range set.List() {
		m := v.(map[string]interface{})
		var filterValues []string
		for _, e := range m["values"].([]interface{}) {
			filterValues = append(filterValues, e.(string))
		}

		switch name := m["name"].(string); name {
		case "creation_dates":
			filters.CreationDates = filterValues
		case "snapshot_ids":
			filters.SnapshotIds = filterValues
		case "subregion_names":
			filters.SubregionNames = filterValues
		case "tag_keys":
			filters.TagKeys = filterValues
		//TODO: case "tags":
		// 	filters.Tags = filterValues
		case "tag_values":
			filters.TagValues = filterValues
		case "volume_ids":
			filters.VolumeIds = filterValues
		case "volume_sizes":
			filters.VolumeSizes = utils.StringSliceToInt64Slice(filterValues)
		case "volume_types":
			filters.VolumeTypes = filterValues
		default:
			log.Printf("[Debug] Unknown Filter Name: %s.", name)
		}
	}
	return filters
}
