package oapi

import (
	"context"
	"log"
	"time"

	"github.com/outscale/goutils/sdk/ptr"
	"github.com/outscale/osc-sdk-go/v3/pkg/options"
	"github.com/outscale/osc-sdk-go/v3/pkg/osc"
	"github.com/samber/lo"
	"github.com/spf13/cast"

	"github.com/davecgh/go-spew/spew"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/framework/fwhelpers/from"
	"github.com/outscale/terraform-provider-outscale/internal/utils"
)

func DataSourceOutscaleVolume() *schema.Resource {
	return &schema.Resource{
		ReadContext: datasourceOAPIVolumeRead,

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
			"tags": TagsSchemaComputedSDK(),
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

func datasourceOAPIVolumeRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.OutscaleClient).OSC

	filters, filtersOk := d.GetOk("filter")
	volumeIds, VolumeIdsOk := d.GetOk("volume_id")

	params := osc.ReadVolumesRequest{
		Filters: &osc.FiltersVolume{},
	}
	if VolumeIdsOk {
		params.Filters.VolumeIds = &[]string{volumeIds.(string)}
	}

	var err error
	if filtersOk {
		params.Filters, err = buildOutscaleOSCAPIDataSourceVolumesFilters(filters.(*schema.Set))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	resp, err := client.ReadVolumes(ctx, params, options.WithRetryTimeout(5*time.Minute))
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("Found These Volumes %s", spew.Sdump(resp.Volumes))

	filteredVolumes := ptr.From(resp.Volumes)[:]

	var volume osc.Volume
	if len(filteredVolumes) < 1 {
		return diag.FromErr(ErrNoResults)
	}

	if len(filteredVolumes) > 1 {
		return diag.FromErr(ErrMultipleResults)
	}

	// Query returned single result.
	volume = filteredVolumes[0]
	log.Printf("[DEBUG] outscale_volume - Single Volume found: %s", volume.VolumeId)
	return diag.FromErr(volumeOAPIDescriptionAttributes(d, &volume))
}

func volumeOAPIDescriptionAttributes(d *schema.ResourceData, volume *osc.Volume) error {
	if err := d.Set("volume_id", volume.VolumeId); err != nil {
		return err
	}
	if err := d.Set("creation_date", from.ISO8601(volume.CreationDate)); err != nil {
		return err
	}
	if err := d.Set("subregion_name", volume.SubregionName); err != nil {
		return err
	}
	if err := d.Set("size", volume.Size); err != nil {
		return err
	}
	if err := d.Set("snapshot_id", ptr.From(volume.SnapshotId)); err != nil {
		return err
	}
	if err := d.Set("volume_type", volume.VolumeType); err != nil {
		return err
	}
	if err := d.Set("state", volume.State); err != nil {
		return err
	}
	if err := d.Set("volume_id", volume.VolumeId); err != nil {
		return err
	}
	if err := d.Set("iops", getIops(volume.VolumeType, volume.Iops)); err != nil {
		return err
	}

	if volume.LinkedVolumes != nil {
		if err := d.Set("linked_volumes", getLinkedVolumes(volume.LinkedVolumes)); err != nil {
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

	if volume.Tags != nil {
		if err := d.Set("tags", FlattenOAPITagsSDK(volume.Tags)); err != nil {
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

	d.SetId(volume.VolumeId)
	return nil
}

func buildOutscaleOSCAPIDataSourceVolumesFilters(set *schema.Set) (*osc.FiltersVolume, error) {
	var filters osc.FiltersVolume
	for _, v := range set.List() {
		m := v.(map[string]interface{})
		var filterValues []string
		for _, e := range m["values"].([]interface{}) {
			filterValues = append(filterValues, e.(string))
		}

		switch name := m["name"].(string); name {
		case "creation_dates":
			dates, err := utils.StringSliceToTimeSlice(filterValues, "creation_dates")
			if err != nil {
				return nil, err
			}
			filters.CreationDates = &dates
		case "snapshot_ids":
			filters.SnapshotIds = &filterValues
		case "subregion_names":
			filters.SubregionNames = &filterValues
		case "tags":
			filters.Tags = &filterValues
		case "tag_keys":
			filters.TagKeys = &filterValues
		case "tag_values":
			filters.TagValues = &filterValues
		case "volume_ids":
			filters.VolumeIds = &filterValues
		case "volume_sizes":
			filters.VolumeSizes = new(utils.StringSliceToIntSlice(filterValues))
		case "volume_types":
			filters.VolumeTypes = new(lo.Map(filterValues, func(s string, _ int) osc.VolumeType { return osc.VolumeType(s) }))
		case "link_volume_vm_ids":
			filters.LinkVolumeVmIds = &filterValues
		case "volume_states":
			filters.VolumeStates = new(lo.Map(filterValues, func(s string, _ int) osc.VolumeState { return osc.VolumeState(s) }))
		case "link_volume_link_states":
			filters.LinkVolumeLinkStates = new(lo.Map(filterValues, func(s string, _ int) osc.LinkedVolumeState { return osc.LinkedVolumeState(s) }))
		case "link_volume_delete_on_vm_deletion":
			filters.LinkVolumeDeleteOnVmDeletion = new(cast.ToBool(filterValues))
		case "link_volume_link_dates":
			dates, err := utils.StringSliceToTimeSlice(filterValues, "link_volume_link_dates")
			if err != nil {
				return nil, err
			}
			filters.LinkVolumeLinkDates = &dates
		case "link_volume_device_names":
			filters.LinkVolumeDeviceNames = &filterValues
		default:
			return nil, utils.UnknownDataSourceFilterError(name)
		}
	}
	return &filters, nil
}
