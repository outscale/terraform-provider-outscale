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

func dataSourceOutscaleOAPIImage() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOutscaleOAPIImageRead,

		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			"permission": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"image_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			// Computed values.
			"architecture": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"creation_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"osu_location": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"account_alias": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"account_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_public": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"root_device_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"root_device_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"request_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			// Complex computed values
			"block_device_mapping": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"device_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"no_device": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"virtual_device_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"bsu": {
							Type:     schema.TypeMap,
							Computed: true,
						},
					},
				},
			},
			"product_code": {
				Type:     schema.TypeSet,
				Computed: true,
				Set:      amiProductCodesHash,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"product_code": {
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
			"state_comment": {
				Type:     schema.TypeMap,
				Computed: true,
			},
			"tag": dataSourceTagsSchema(),
		},
	}
}

func dataSourceOutscaleOAPIImageRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	executableUsers, executableUsersOk := d.GetOk("permission")
	filters, filtersOk := d.GetOk("filter")
	ai, aisOk := d.GetOk("account_id")
	imageID, imageIDOk := d.GetOk("image_id")

	if executableUsersOk == false && filtersOk == false && aisOk == false && imageIDOk == false {
		return fmt.Errorf("One of executable_users, filters, or account_id must be assigned, or image_id must be provided")
	}

	params := &fcu.DescribeImagesInput{}
	if executableUsersOk {
		params.ExecutableUsers = expandStringList(executableUsers.([]interface{}))
	}
	if filtersOk {
		params.Filters = buildOutscaleDataSourceFilters(filters.(*schema.Set))
	}
	if imageIDOk {
		params.ImageIds = []*string{aws.String(imageID.(string))}
	}
	if aisOk {
		params.Owners = []*string{aws.String(ai.(string))}
	}

	var res *fcu.DescribeImagesOutput
	var err error
	err = resource.Retry(40*time.Minute, func() *resource.RetryError {
		res, err = conn.VM.DescribeImages(params)

		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded") {
				fmt.Printf("[INFO] Request limit exceeded")
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}

		return resource.RetryableError(err)
	})

	if err != nil {
		return err
	}

	if len(res.Images) < 1 {
		return fmt.Errorf("your query returned no results, please change your search criteria and try again")
	}

	if len(res.Images) > 1 {
		return fmt.Errorf("your query returned more than one result, please try a more " +
			"specific search criteria")
	}

	d.Set("request_id", res.RequestId)

	return omiOAPIDescriptionAttributes(d, res.Images[0])
}

// populate the numerous fields that the image description returns.
func omiOAPIDescriptionAttributes(d *schema.ResourceData, image *fcu.Image) error {

	d.SetId(*image.ImageId)
	d.Set("architecture", image.Architecture)
	if image.CreationDate != nil {
		d.Set("creation_date", image.CreationDate)
	} else {
		d.Set("creation_date", "")
	}
	if image.Description != nil {
		d.Set("description", image.Description)
	} else {
		d.Set("description", "")
	}
	d.Set("hypervisor", image.Hypervisor)
	d.Set("image_id", image.ImageId)
	d.Set("osu_location", image.ImageLocation)
	if image.ImageOwnerAlias != nil {
		d.Set("account_alias", image.ImageOwnerAlias)
	} else {
		d.Set("account_alias", "")
	}
	d.Set("account_id", image.OwnerId)
	d.Set("type", image.ImageType)
	d.Set("name", image.Name)
	d.Set("is_public", image.Public)
	if image.RootDeviceName != nil {
		d.Set("root_device_name", image.RootDeviceName)
	} else {
		d.Set("root_device_name", "")
	}
	d.Set("root_device_type", image.RootDeviceType)
	d.Set("state", image.State)
	d.Set("virtualization_type", image.VirtualizationType)
	// Complex types get their own functions
	if err := d.Set("block_device_mapping", amiBlockDeviceMappings(image.BlockDeviceMappings)); err != nil {
		return err
	}
	if err := d.Set("product_code", amiProductCodes(image.ProductCodes)); err != nil {
		return err
	}
	if err := d.Set("state_comment", amiStateReason(image.StateReason)); err != nil {
		return err
	}

	return d.Set("tag", tagsToMap(image.Tags))
}
