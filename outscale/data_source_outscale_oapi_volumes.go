package outscale

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/davecgh/go-spew/spew"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-outscale/osc/fcu"
)

func datasourceOutscaleOAPIVolumes() *schema.Resource {
	return &schema.Resource{
		Read: datasourceOAPIVolumesRead,

		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			"volumes": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						//Schema: map[string]*schema.Schema{
						// Arguments
						"sub_region_name": {
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
									"delete_on_vm_termination": {
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
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func datasourceOAPIVolumesRead(d *schema.ResourceData, meta interface{}) error {

	conn := meta.(*OutscaleClient).FCU

	filters, filtersOk := d.GetOk("filter")
	VolumeIds, VolumeIdsOk := d.GetOk("volume_id")

	params := &fcu.DescribeVolumesInput{}
	if filtersOk {
		params.Filters = buildOutscaleDataSourceFilters(filters.(*schema.Set))
	}
	if VolumeIdsOk {
		params.VolumeIds = []*string{aws.String(VolumeIds.(string))}
	}

	var resp *fcu.DescribeVolumesOutput
	var err error

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, err = conn.VM.DescribeVolumes(params)
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

	filteredVolumes := resp.Volumes[:]

	if len(filteredVolumes) < 1 {
		return fmt.Errorf("your query returned no results, please change your search criteria and try again")
	}

	return volumesOAPIDescriptionAttributes(d, filteredVolumes)
}

func volumesOAPIDescriptionAttributes(d *schema.ResourceData, volumes []*fcu.Volume) error {

	i := make([]interface{}, len(volumes))

	for k, v := range volumes {
		im := make(map[string]interface{})

		if v.Attachments != nil {
			a := make([]map[string]interface{}, len(v.Attachments))
			for k, v := range v.Attachments {
				at := make(map[string]interface{})
				if v.DeleteOnTermination != nil {
					at["delete_on_vm_termination"] = *v.DeleteOnTermination
				}
				if v.Device != nil {
					at["device_name"] = *v.Device
				}
				if v.InstanceId != nil {
					at["vm_id"] = *v.InstanceId
				}
				if v.State != nil {
					at["state"] = *v.State
				}
				if v.VolumeId != nil {
					at["volume_id"] = *v.VolumeId
				}
				a[k] = at
			}
			im["linked_volumes"] = a
		}
		if v.AvailabilityZone != nil {
			im["sub_region_name"] = *v.AvailabilityZone
		}
		if v.Iops != nil {
			im["iops"] = *v.Iops
		}
		if v.Size != nil {
			im["size"] = *v.Size
		}
		if v.SnapshotId != nil {
			im["snapshot_id"] = *v.SnapshotId
		}
		if v.Tags != nil {
			im["tags"] = dataSourceTags(v.Tags)
		}
		if v.VolumeType != nil {
			im["type"] = *v.VolumeType
		}
		if v.State != nil {
			im["state"] = *v.State
		}
		if v.VolumeId != nil {
			im["volume_id"] = *v.VolumeId
		}
		i[k] = im
	}

	err := d.Set("volumes", i)
	d.SetId(resource.UniqueId())

	return err
}
