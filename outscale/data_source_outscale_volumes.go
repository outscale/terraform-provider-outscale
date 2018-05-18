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
			"volume_id": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"volume_set": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
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
						"tag": tagsSchema(),
						"volume_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"request_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func datasourceVolumesRead(d *schema.ResourceData, meta interface{}) error {

	conn := meta.(*OutscaleClient).FCU

	filters, filtersOk := d.GetOk("filter")
	volumeIds, volumeIdsOk := d.GetOk("volume_id")

	if !filtersOk && !volumeIdsOk {
		return fmt.Errorf("One of volume_id or filters must be assigned")
	}

	params := &fcu.DescribeVolumesInput{}
	if filtersOk {
		params.Filters = buildOutscaleDataSourceFilters(filters.(*schema.Set))
	}
	if volumeIdsOk {
		params.VolumeIds = expandStringList(volumeIds.([]interface{}))
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

	d.Set("request_id", resp.RequestId)

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
				at["delete_on_termination"] = aws.BoolValue(v.DeleteOnTermination)
				at["device"] = aws.StringValue(v.Device)
				at["instance_id"] = aws.StringValue(v.InstanceId)
				at["state"] = aws.StringValue(v.State)
				at["volume_id"] = aws.StringValue(v.VolumeId)

				a[k] = at
			}
			im["attachment_set"] = a
		}
		im["availability_zone"] = aws.StringValue(v.AvailabilityZone)
		im["iops"] = aws.Int64Value(v.Iops)
		im["size"] = aws.Int64Value(v.Size)
		im["snapshot_id"] = aws.StringValue(v.SnapshotId)
		im["volume_type"] = aws.StringValue(v.VolumeType)
		im["status"] = aws.StringValue(v.State)
		im["volume_id"] = aws.StringValue(v.VolumeId)
		im["tag_set"] = tagsToMap(v.Tags)

		i[k] = im
	}

	d.SetId(resource.UniqueId())

	return d.Set("volume_set", i)
}
