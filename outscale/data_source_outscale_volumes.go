package outscale

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/spf13/cast"

	"github.com/davecgh/go-spew/spew"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/outscale/osc-go/oapi"
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
			"volumes": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
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
						"tags": tagsOAPIListSchemaComputed(),
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
	conn := meta.(*OutscaleClient).OAPI

	filters, filtersOk := d.GetOk("filter")
	volumeIds, volumeIdsOk := d.GetOk("volume_id")
	params := &oapi.ReadVolumesRequest{
		Filters: oapi.FiltersVolume{},
	}

	if volumeIdsOk {
		volIDs := expandStringValueList(volumeIds.([]interface{}))
		params.Filters.VolumeIds = volIDs
	}
	log.Printf("LOOOGGG___ filtersOk \n %+v \n", filtersOk)

	if filtersOk {
		params.Filters = buildOutscaleOAPIDataSourceVolumesFilters(filters.(*schema.Set))
	}

	log.Printf("LOOOGGG___ Filters \n %+v \n", params.Filters)

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
	if err != nil {
		return err
	}

	log.Printf("Found These Volumes %s", spew.Sdump(rs.OK.Volumes))

	volumes := rs.OK.Volumes

	if len(volumes) < 1 {
		return fmt.Errorf("your query returned no results, please change your search criteria and try again")
	}

	log.Printf("LOOOGGG___ volumes \n %+v \n", volumes)

	if err := d.Set("volumes", getOAPIVolumes(volumes)); err != nil {
		return err
	}

	d.Set("request_id", rs.OK.ResponseContext.RequestId)

	d.SetId(resource.UniqueId())

	return nil
}

func getOAPIVolumes(volumes []oapi.Volume) (res []map[string]interface{}) {
	for _, v := range volumes {
		res = append(res, map[string]interface{}{
			"iops":           v.Iops,
			"linked_volumes": v.LinkedVolumes,
			"size":           v.Size,
			"snapshot_id":    v.SnapshotId,
			"state":          v.State,
			"subregion_name": v.SubregionName,
			"tags":           tagsOAPIToMap(v.Tags),
			"volume_id":      v.VolumeId,
			"volume_type":    v.VolumeType,
		})
	}
	return
}

func getOAPILinkedVolumes(linkedVolumes []oapi.LinkedVolume) (res []map[string]interface{}) {
	for _, l := range linkedVolumes {
		res = append(res, map[string]interface{}{
			"delete_on_vm_deletion": cast.ToString(l.DeleteOnVmDeletion),
			"device_name":           l.DeviceName,
			"vm_id":                 l.VmId,
			"state":                 l.State,
			"volume_id":             l.VolumeId,
		})
	}
	return
}
