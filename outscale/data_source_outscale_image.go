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

func dataSourceOutscaleImage() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOutscaleImageRead,

		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			"executable_by": {
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
			"owners": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
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
			"image_location": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"image_owner_alias": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"image_owner_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"image_state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"image_type": {
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
			// Complex computed values
			"block_device_mappings": {
				Type:     schema.TypeSet,
				Computed: true,
				Set:      amiBlockDeviceMappingHash,
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
						"virtual_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ebs": {
							Type:     schema.TypeMap,
							Computed: true,
						},
					},
				},
			},
			"product_codes": {
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
			"state_reason": {
				Type:     schema.TypeMap,
				Computed: true,
			},
			"tag_set": dataSourceTagsSchema(),
		},
	}
}

func dataSourceOutscaleImageRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	executableUsers, executableUsersOk := d.GetOk("executable_by")
	filters, filtersOk := d.GetOk("filter")
	owners, ownersOk := d.GetOk("owners")
	imageID, imageIDOk := d.GetOk("image_id")

	if executableUsersOk == false && filtersOk == false && ownersOk == false && imageIDOk == false {
		return fmt.Errorf("One of executable_users, filters, or owners must be assigned, or image_id must be provided")
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
	if ownersOk {
		o := expandStringList(owners.([]interface{}))

		if len(o) > 0 {
			params.Owners = o
		}
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
		return fmt.Errorf("Your query returned no results. Please change your search criteria and try again.")
	}

	if len(res.Images) > 1 {
		return fmt.Errorf("Your query returned more than one result. Please try a more " +
			"specific search criteria.")
	}

	return omiDescriptionAttributes(d, res.Images[0])
}

// populate the numerous fields that the image description returns.
func omiDescriptionAttributes(d *schema.ResourceData, image *fcu.Image) error {

	d.SetId(*image.ImageId)
	d.Set("architecture", image.Architecture)
	if image.CreationDate != nil {
		d.Set("creation_date", image.CreationDate)
	}
	if image.Description != nil {
		d.Set("description", image.Description)
	}
	d.Set("hypervisor", image.Hypervisor)
	d.Set("image_id", image.ImageId)
	d.Set("image_location", image.ImageLocation)
	if image.ImageOwnerAlias != nil {
		d.Set("image_owner_alias", image.ImageOwnerAlias)
	}
	d.Set("image_owner_id", image.OwnerId)
	d.Set("image_type", image.ImageType)
	d.Set("name", image.Name)
	d.Set("is_public", image.Public)
	if image.RootDeviceName != nil {
		d.Set("root_device_name", image.RootDeviceName)
	}
	d.Set("root_device_type", image.RootDeviceType)
	d.Set("image_state", image.State)
	d.Set("virtualization_type", image.VirtualizationType)
	// Complex types get their own functions
	if err := d.Set("block_device_mappings", amiBlockDeviceMappings(image.BlockDeviceMappings)); err != nil {
		return err
	}
	if err := d.Set("product_codes", amiProductCodes(image.ProductCodes)); err != nil {
		return err
	}
	if err := d.Set("state_reason", amiStateReason(image.StateReason)); err != nil {
		return err
	}
	if err := d.Set("tag_set", dataSourceTags(image.Tags)); err != nil {
		return err
	}
	return nil
}
