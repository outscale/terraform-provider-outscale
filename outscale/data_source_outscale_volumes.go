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

func datasourceOutscaleVolumes() *schema.Resource {
	return &schema.Resource{
		Read: datasourceVolumesRead,

		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			"volume_set": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						//Schema: map[string]*schema.Schema{
						// Arguments
						"availability_zone": {
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
						"volume_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						// Attributes
						"attachment_set": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"delete_on_termination": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"device": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"instance_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"status": {
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
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tag_set": {
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
						"tags": tagsSchema(),
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

func datasourceVolumesRead(d *schema.ResourceData, meta interface{}) error {

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
		return fmt.Errorf("Your query returned no results. Please change your search criteria and try again.")
	}

	return volumesDescriptionAttributes(d, filteredVolumes)
}

func volumesDescriptionAttributes(d *schema.ResourceData, volumes []*fcu.Volume) error {

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
