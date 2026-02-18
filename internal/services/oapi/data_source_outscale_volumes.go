package oapi

import (
	"log"
	"time"

	"github.com/outscale/osc-sdk-go/v3/pkg/osc"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/utils"

	"github.com/davecgh/go-spew/spew"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceOutscaleVolumes() *schema.Resource {
	return &schema.Resource{
		Read: datasourceOAPIVolumesRead,

		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			"volume_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"volumes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"creation_date": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"iops": {
							Type:     schema.TypeInt,
							Computed: true,
						},
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
						"size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"snapshot_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"state": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"subregion_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tags": TagsSchemaComputedSDK(),
						"volume_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"volume_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"request_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func datasourceOAPIVolumesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*client.OutscaleClient).OSC

	filters, filtersOk := d.GetOk("filter")
	volumeIds, volumeIdsOk := d.GetOk("volume_id")
	params := osc.ReadVolumesRequest{
		Filters: &osc.FiltersVolume{},
	}

	if volumeIdsOk {
		volIDs := utils.InterfaceSliceToStringSlice(volumeIds.([]interface{}))
		filter := osc.FiltersVolume{}
		filter.SetVolumeIds(volIDs)
		params.SetFilters(filter)
	}

	var err error
	if filtersOk {
		params.Filters, err = buildOutscaleOSCAPIDataSourceVolumesFilters(filters.(*schema.Set))
		if err != nil {
			return err
		}
	}

	var resp osc.ReadVolumesResponse
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		rp, httpResp, err := client.VolumeApi.ReadVolumes(ctx).ReadVolumesRequest(params).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		resp = rp
		return nil
	})
	if err != nil {
		return err
	}

	log.Printf("Found These Volumes %s", spew.Sdump(resp.GetVolumes()))

	volumes := resp.GetVolumes()

	if len(volumes) < 1 {
		return ErrNoResults
	}

	if err := d.Set("volumes", getOAPIVolumes(volumes)); err != nil {
		return err
	}

	d.SetId(id.UniqueId())

	return nil
}

func getOAPIVolumes(volumes []osc.Volume) (res []map[string]interface{}) {
	for _, v := range volumes {
		res = append(res, map[string]interface{}{
			"creation_date":  v.CreationDate,
			"iops":           getIops(v.GetVolumeType(), v.GetIops()),
			"linked_volumes": getLinkedVolumes(v.GetLinkedVolumes()),
			"size":           v.Size,
			"snapshot_id":    v.SnapshotId,
			"state":          v.State,
			"subregion_name": v.SubregionName,
			"tags":           FlattenOAPITagsSDK(v.Tags),
			"volume_id":      v.VolumeId,
			"volume_type":    v.VolumeType,
		})
	}
	return
}

func getLinkedVolumes(linkedVolumes []osc.LinkedVolume) (res []map[string]interface{}) {
	for _, l := range linkedVolumes {
		res = append(res, map[string]interface{}{
			"delete_on_vm_deletion": l.DeleteOnVmDeletion,
			"device_name":           l.DeviceName,
			"vm_id":                 l.VmId,
			"state":                 l.State,
			"volume_id":             l.VolumeId,
		})
	}
	return
}

func getIops(volumeType string, iops int32) int32 {
	if volumeType != "standard" {
		return iops
	}
	return DefaultIops
}
