package outscale

import (
	"context"
	"log"
	"time"

	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/spf13/cast"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
			"creation_date": {
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
			"tags": dataSourceTagsSchema(),
			"volume_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"request_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func datasourceOAPIVolumeRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	req := oscgo.ReadVolumesRequest{}

	if filters, filtersOk := d.GetOk("filter"); filtersOk {
		req.SetFilters(buildOutscaleOSCAPIDataSourceVolumesFilters(filters.(*schema.Set)))
	}

	var resp oscgo.ReadVolumesResponse
	var err error

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		rp, httpResp, err := conn.VolumeApi.ReadVolumes(context.Background()).ReadVolumesRequest(req).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		resp = rp
		return nil
	})

	if err != nil {
		return err
	}
	volumes := resp.GetVolumes()

	if err = utils.IsResponseEmptyOrMutiple(len(volumes), "Access Key"); err != nil {
		return err
	}
	return volumeOAPIDescriptionAttributes(d, &volumes[0])
}

func volumeOAPIDescriptionAttributes(d *schema.ResourceData, volume *oscgo.Volume) error {
	if err := d.Set("volume_id", volume.GetVolumeId()); err != nil {
		return err
	}
	if err := d.Set("creation_date", volume.GetCreationDate()); err != nil {
		return err
	}
	if err := d.Set("subregion_name", volume.GetSubregionName()); err != nil {
		return err
	}
	if err := d.Set("size", volume.GetSize()); err != nil {
		return err
	}
	if err := d.Set("snapshot_id", volume.GetSnapshotId()); err != nil {
		return err
	}
	if err := d.Set("volume_type", volume.GetVolumeType()); err != nil {
		return err
	}
	if err := d.Set("state", volume.GetState()); err != nil {
		return err
	}
	if err := d.Set("volume_id", volume.GetVolumeId()); err != nil {
		return err
	}
	if err := d.Set("iops", getIops(volume.GetVolumeType(), volume.GetIops())); err != nil {
		return err
	}

	if volume.LinkedVolumes != nil {
		if err := d.Set("linked_volumes", getLinkedVolumes(volume.GetLinkedVolumes())); err != nil {
			return err
		}
	} else {
		if err := d.Set("linked_volumes", []map[string]interface{}{
			{
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
			{
				"key":   "",
				"value": "",
			},
		}); err != nil {
			return err
		}
	}

	d.SetId(volume.GetVolumeId())
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
		case "tags":
			filters.SetTags(filterValues)
		case "tag_keys":
			filters.SetTagKeys(filterValues)
		case "tag_values":
			filters.SetTagValues(filterValues)
		case "volume_ids":
			filters.SetVolumeIds(filterValues)
		case "volume_sizes":
			filters.SetVolumeSizes(utils.StringSliceToInt32Slice(filterValues))
		case "volume_types":
			filters.SetVolumeTypes(filterValues)
		case "link_volume_vm_ids":
			filters.SetLinkVolumeVmIds(filterValues)
		case "volume_states":
			filters.SetVolumeStates(filterValues)
		case "link_volume_link_states":
			filters.SetLinkVolumeLinkStates(filterValues)
		case "link_volume_delete_on_vm_deletion":
			filters.SetLinkVolumeDeleteOnVmDeletion(cast.ToBool(filterValues))
		case "link_volume_link_dates":
			filters.SetLinkVolumeLinkDates(filterValues)
		case "link_volume_device_names":
			filters.SetLinkVolumeDeviceNames(filterValues)
		default:
			log.Printf("[Debug] Unknown Filter Name: %s.", name)
		}
	}
	return filters
}
