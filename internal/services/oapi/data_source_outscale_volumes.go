package oapi

import (
	"context"
	"log"
	"time"

	"github.com/outscale/goutils/sdk/ptr"
	"github.com/outscale/osc-sdk-go/v3/pkg/options"
	"github.com/outscale/osc-sdk-go/v3/pkg/osc"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/framework/fwhelpers/from"
	"github.com/outscale/terraform-provider-outscale/internal/utils"

	"github.com/davecgh/go-spew/spew"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceOutscaleVolumes() *schema.Resource {
	return &schema.Resource{
		ReadContext: datasourceOAPIVolumesRead,

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

func datasourceOAPIVolumesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.OutscaleClient).OSC

	filters, filtersOk := d.GetOk("filter")
	volumeIds, volumeIdsOk := d.GetOk("volume_id")
	params := osc.ReadVolumesRequest{
		Filters: &osc.FiltersVolume{},
	}

	if volumeIdsOk {
		volIDs := utils.InterfaceSliceToStringSlice(volumeIds.([]interface{}))
		filter := osc.FiltersVolume{}
		filter.VolumeIds = &volIDs
		params.Filters = &filter
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

	volumes := ptr.From(resp.Volumes)

	if len(volumes) < 1 {
		return diag.FromErr(ErrNoResults)
	}

	if err := d.Set("volumes", getOAPIVolumes(volumes)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(id.UniqueId())

	return nil
}

func getOAPIVolumes(volumes []osc.Volume) (res []map[string]interface{}) {
	for _, v := range volumes {
		res = append(res, map[string]interface{}{
			"creation_date":  from.ISO8601(v.CreationDate),
			"iops":           getIops(v.VolumeType, v.Iops),
			"linked_volumes": getLinkedVolumes(v.LinkedVolumes),
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

func getIops(volumeType osc.VolumeType, iops int) int {
	if volumeType != osc.VolumeTypeStandard {
		return iops
	}
	return DefaultIops
}
