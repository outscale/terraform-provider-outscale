package outscale

import (
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-outscale/osc/fcu"
)

func datasourceOutscaleOAPIVolumes() *schema.Resource {
	return &schema.Resource{
		Read: datasourceVolumesRead,

		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			"volume": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						//Schema: map[string]*schema.Schema{
						"linked_volume": {
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
						"state": {
							Type:     schema.TypeString,
							Computed: true,
						},
						// Attributes
						"tag": {
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
						//						"tags": tagsSchema(),
						"volume_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
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

	if filtersOk == false {
		return fmt.Errorf("One of filters must be assigned")
	}

	// Build up search parameters

	request := &fcu.DescribeVolumesInput{
		VolumeIds: []*string{aws.String(d.Id())},
	}
	if filtersOk {
		request.Filters = buildOutscaleDataSourceFilters(filters.(*schema.Set))
	}

	var response *fcu.DescribeVolumesOutput
	var err error

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.VM.DescribeVolumes(request)
		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded:") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return resource.NonRetryableError(err)
	})

	if err != nil {
		if strings.Contains(fmt.Sprint(err), "InvalidVolume.NotFound") {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error reading Outscale volume %s: %s", d.Id(), err)
	}

	return volumesDescriptionOAPIAttributes(d, response.Volumes)
}

// populate the numerous fields that the volume description returns.
func volumesDescriptionOAPIAttributes(d *schema.ResourceData, volumes []*fcu.Volume) error {

	i := make([]interface{}, len(volumes))

	for k, v := range volumes {
		im := make(map[string]interface{})

		if v.Attachments != nil {
			a := make([]map[string]interface{}, len(v.Attachments))
			for k, v := range v.Attachments {
				at := make(map[string]interface{})
				if v.DeleteOnTermination != nil {
					at["delete_on_termination"] = *v.DeleteOnTermination
				}
				if v.Device != nil {
					at["device"] = *v.Device
				}
				if v.InstanceId != nil {
					at["instance_id"] = *v.InstanceId
				}
				if v.State != nil {
					at["state"] = *v.State
				}
				if v.VolumeId != nil {
					at["volume_id"] = *v.VolumeId
				}
				a[k] = at
			}
			im["attachment_set"] = a
		}
		if v.AvailabilityZone != nil {
			im["availability_zone"] = *v.AvailabilityZone
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
			im["tag_set"] = dataSourceTags(v.Tags)
		}
		if v.VolumeType != nil {
			im["volume_type"] = *v.VolumeType
		}
		if v.State != nil {
			im["status"] = *v.State
		}
		if v.VolumeId != nil {
			im["volume_id"] = *v.VolumeId
		}
		i[k] = im
	}

	err := d.Set("volume_set", i)
	d.SetId(resource.UniqueId())

	return err
}
