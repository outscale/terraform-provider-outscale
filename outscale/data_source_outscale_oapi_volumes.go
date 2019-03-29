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
						//Schema: map[string]*schema.Schema{
						// Arguments
						"subregion_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"iops": {
							Type: schema.TypeInt,

							Computed: true,
						},
						"size": {
							Type: schema.TypeInt,

							Computed: true,
						},
						"snapshot_id": {
							Type: schema.TypeString,

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
						"tags": tagsOAPIListSchemaComputed(),
						"volume_id": {
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

	if filtersOk {
		params.Filters = buildOutscaleOAPIDataSourceVolumesFilters(filters.(*schema.Set))
	}

	if volumeIdsOk {
		volIDs := expandStringValueList(volumeIds.([]interface{}))
		params.Filters.VolumeIds = volIDs
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

	if err != nil {
		return err
	}
	resp = rs.OK

	log.Printf("Found These Volumes %s", spew.Sdump(resp.Volumes))

	filteredVolumes := resp.Volumes[:]

	if len(filteredVolumes) < 1 {
		return fmt.Errorf("your query returned no results, please change your search criteria and try again")
	}

	d.Set("request_id", resp.ResponseContext.RequestId)

	return volumesOAPIDescriptionAttributes(d, filteredVolumes)
}

func volumesOAPIDescriptionAttributes(d *schema.ResourceData, volumes []oapi.Volume) error {

	i := make([]interface{}, len(volumes))

	for k, v := range volumes {
		im := make(map[string]interface{})

		if v.LinkedVolumes != nil {
			a := make([]map[string]interface{}, len(v.LinkedVolumes))
			for k, v := range v.LinkedVolumes {
				at := make(map[string]interface{})
				//if v.DeleteOnVmDeletion != nil {
				at["delete_on_vm_termination"] = v.DeleteOnVmDeletion
				//}
				if v.DeviceName != "" {
					at["device_name"] = v.DeviceName
				}
				if v.VmId != "" {
					at["vm_id"] = v.VmId
				}
				if v.State != "" {
					at["state"] = v.State
				}
				if v.VolumeId != "" {
					at["volume_id"] = v.VolumeId
				}
				a[k] = at
			}
			im["linked_volumes"] = a
		}
		if v.SubregionName != "" {
			im["subregion_name"] = v.SubregionName
		}
		//if v.Iops != nil {
		im["iops"] = v.Iops
		//}
		//if v.Size != nil {
		im["size"] = v.Size
		//}
		if v.SnapshotId != "" {
			im["snapshot_id"] = v.SnapshotId
		}
		if v.Tags != nil {
			im["tags"] = tagsOAPIToMap(v.Tags)
		}
		if v.VolumeType != "" {
			im["type"] = v.VolumeType
		}
		if v.State != "" {
			im["state"] = v.State
		}
		if v.VolumeId != "" {
			im["volume_id"] = v.VolumeId
		}
		i[k] = im
	}

	err := d.Set("volumes", i)
	d.SetId(resource.UniqueId())

	return err
}
