package outscale

import (
	"context"
	"fmt"
	"github.com/antihax/optional"
	oscgo "github.com/marinsalinas/osc-sdk-go"
	"log"
	"strings"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
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
	conn := meta.(*OutscaleClient).OSCAPI

	filters, filtersOk := d.GetOk("filter")
	volumeIds, VolumeIdsOk := d.GetOk("volume_id")

	params := oscgo.ReadVolumesRequest{
		Filters: &oscgo.FiltersVolume{},
	}
	if VolumeIdsOk {
		params.Filters.SetVolumeIds([]string{volumeIds.(string)})
	}

	if filtersOk {
		params.SetFilters(buildOutscaleOSCAPIDataSourceVolumesFilters(filters.(*schema.Set)))
	}

	var resp oscgo.ReadVolumesResponse
	var err error

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, _, err = conn.VolumeApi.ReadVolumes(context.Background(), &oscgo.ReadVolumesOpts{ReadVolumesRequest: optional.NewInterface(params)})
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

	filteredVolumes := resp.GetVolumes()[:]

	var volume oscgo.Volume
	if len(filteredVolumes) < 1 {
		return fmt.Errorf("your query returned no results, please change your search criteria and try again")
	}

	if len(filteredVolumes) > 1 {
		return fmt.Errorf("your query returned more than one result, please try a more " +
			"specific search criteria")
	}

	// Query returned single result.
	volume = filteredVolumes[0]
	d.Set("request_id", resp.ResponseContext.GetRequestId())
	log.Printf("[DEBUG] outscale_volume - Single Volume found: %s", volume.GetVolumeId())
	return volumeOAPIDescriptionAttributes(d, &volume)

}

func volumeOAPIDescriptionAttributes(d *schema.ResourceData, volume *oscgo.Volume) error {
	d.SetId(volume.GetVolumeId())
	d.Set("volume_id", volume.GetVolumeId())
	d.Set("subregion_name", volume.GetSubregionName())
	d.Set("size", volume.GetSize())
	d.Set("snapshot_id", volume.GetSnapshotId())
	d.Set("volume_type", volume.GetVolumeType())

	d.Set("state", volume.GetState())
	d.Set("volume_id", volume.GetVolumeId())
	d.Set("iops", volume.GetIops())

	if volume.LinkedVolumes != nil {
		res := make([]map[string]interface{}, len(volume.GetLinkedVolumes()))
		for k, g := range volume.GetLinkedVolumes() {
			r := make(map[string]interface{})
			if g.DeleteOnVmDeletion != nil {
				r["delete_on_vm_deletion"] = g.GetDeleteOnVmDeletion()
			}
			if g.GetDeviceName() != "" {
				r["device_name"] = g.GetDeviceName()
			}
			if g.GetVmId() != "" {
				r["vm_id"] = g.GetVmId()
			}
			if g.GetState() != "" {
				r["state"] = g.GetState()
			}
			if g.GetVolumeId() != "" {
				r["volume_id"] = g.GetVolumeId()
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
				"vm_id":                 "none",
				"state":                 "none",
				"volume_id":             "none",
			},
		}); err != nil {
			return err
		}
	}

	if volume.GetTags() != nil {
		if err := d.Set("tags", tagsOSCAPIToMap(volume.GetTags())); err != nil {
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

func buildOutscaleOSCAPIDataSourceVolumesFilters(set *schema.Set) oscgo.FiltersVolume {
	var filters oscgo.FiltersVolume
	for _, v := range set.List() {
		m := v.(map[string]interface{})
		var filterValues []string
		for _, e := range m["values"].([]interface{}) {
			filterValues = append(filterValues, e.(string))
		}

		switch name := m["name"].(string); name {
		case "creation_dates":
			filters.SetCreationDates(filterValues)
		case "snapshot_ids":
			filters.SetSnapshotIds(filterValues)
		case "subregion_names":
			filters.SetSubregionNames(filterValues)
		case "tag_keys":
			filters.SetTagKeys(filterValues)
		//TODO: case "tags":
		// 	filters.Tags = filterValues
		case "tag_values":
			filters.SetTagValues(filterValues)
		case "volume_ids":
			filters.SetVolumeIds(filterValues)
		case "volume_sizes":
			filters.SetVolumeSizes(utils.StringSliceToInt64Slice(filterValues))
		case "volume_types":
			filters.SetVolumeTypes(filterValues)
		default:
			log.Printf("[Debug] Unknown Filter Name: %s.", name)
		}
	}
	return filters
}
