package outscale

import (
	"context"
	"fmt"
	"log"
	"time"

	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/terraform-providers/terraform-provider-outscale/utils"

	"github.com/davecgh/go-spew/spew"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func datasourceOutscaleOAPIVolumes() *schema.Resource {
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
						"tags": dataSourceTagsSchema(),
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
	conn := meta.(*OutscaleClient).OSCAPI

	filters, filtersOk := d.GetOk("filter")
	volumeIds, volumeIdsOk := d.GetOk("volume_id")
	params := oscgo.ReadVolumesRequest{
		Filters: &oscgo.FiltersVolume{},
	}

	if volumeIdsOk {
		volIDs := utils.InterfaceSliceToStringSlice(volumeIds.([]interface{}))
		filter := oscgo.FiltersVolume{}
		filter.SetVolumeIds(volIDs)
		params.SetFilters(filter)
	}

	if filtersOk {
		params.SetFilters(buildOutscaleOSCAPIDataSourceVolumesFilters(filters.(*schema.Set)))
	}

	log.Printf("LOG____ params: %#+v\n", params.GetFilters())

	var resp oscgo.ReadVolumesResponse
	var err error

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		rp, httpResp, err := conn.VolumeApi.ReadVolumes(context.Background()).ReadVolumesRequest(params).Execute()
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
		return fmt.Errorf("your query returned no results, please change your search criteria and try again")
	}

	if err := d.Set("volumes", getOAPIVolumes(volumes)); err != nil {
		return err
	}

	d.SetId(resource.UniqueId())

	return nil
}

func getOAPIVolumes(volumes []oscgo.Volume) (res []map[string]interface{}) {
	for _, v := range volumes {
		res = append(res, map[string]interface{}{
			"creation_date":  v.CreationDate,
			"iops":           getIops(v.GetVolumeType(), v.GetIops()),
			"linked_volumes": getLinkedVolumes(v.GetLinkedVolumes()),
			"size":           v.Size,
			"snapshot_id":    v.SnapshotId,
			"state":          v.State,
			"subregion_name": v.SubregionName,
			"tags":           tagsOSCAPIToMap(v.GetTags()),
			"volume_id":      v.VolumeId,
			"volume_type":    v.VolumeType,
		})
	}
	return
}

func getLinkedVolumes(linkedVolumes []oscgo.LinkedVolume) (res []map[string]interface{}) {
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
	return utils.DefaultIops
}
