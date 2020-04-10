package outscale

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/antihax/optional"
	oscgo "github.com/marinsalinas/osc-sdk-go"

	"github.com/spf13/cast"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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
			"file_location": {
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
			"image_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"image_name": {
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
			"block_device_mappings": {
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
			"product_codes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"state_comment": {
				Type:     schema.TypeMap,
				Computed: true,
			},
			"permissions_to_launch": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"global_permission": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"account_id": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
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
		},
	}
}

func dataSourceOutscaleOAPIImageRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	filters, filtersOk := d.GetOk("filter")
	executableUsers, executableUsersOk := d.GetOk("permission")
	ai, aisOk := d.GetOk("account_id")
	imageID, imageIDOk := d.GetOk("image_id")
	if !executableUsersOk && !filtersOk && !aisOk && !imageIDOk {
		return fmt.Errorf("One of executable_users, filters, or account_id must be assigned, or image_id must be provided")
	}

	filtersReq := &oscgo.FiltersImage{}
	if filtersOk {
		filtersReq = buildOutscaleOAPIDataSourceImagesFilters(filters.(*schema.Set))
	}
	if imageIDOk {
		filtersReq.SetImageIds([]string{imageID.(string)})
	}
	if aisOk {
		filtersReq.SetAccountIds([]string{ai.(string)})
	}
	if executableUsersOk {
		filtersReq.SetPermissionsToLaunchAccountIds(expandStringValueList(executableUsers.([]interface{})))
	}

	req := oscgo.ReadImagesRequest{Filters: filtersReq}

	var resp oscgo.ReadImagesResponse
	var err error
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, _, err = conn.ImageApi.ReadImages(context.Background(), &oscgo.ReadImagesOpts{ReadImagesRequest: optional.NewInterface(req)})
		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded:") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return err
	}

	images := resp.GetImages()

	if len(images) < 1 {
		return fmt.Errorf("your query returned no results, please change your search criteria and try again")
	}
	if len(images) > 1 {
		return fmt.Errorf("your query returned more than one result, please try a more specific search criteria")
	}

	return resourceDataAttrSetter(d, func(set AttributeSetter) error {
		image := images[0]
		d.SetId(*image.ImageId)

		if err := set("architecture", image.Architecture); err != nil {
			return err
		}

		if err := set("creation_date", image.CreationDate); err != nil {
			return err
		}

		if err := set("image_id", image.ImageId); err != nil {
			return err
		}
		if err := set("file_location", image.FileLocation); err != nil {
			return err
		}
		if err := set("account_alias", image.AccountAlias); err != nil {
			return err
		}
		if err := set("account_id", image.AccountId); err != nil {
			return err
		}
		if err := set("image_type", image.ImageType); err != nil {
			return err
		}
		if err := set("image_name", image.ImageName); err != nil {
			return err
		}
		if err := set("root_device_name", image.RootDeviceName); err != nil {
			return err
		}
		if err := set("root_device_type", image.RootDeviceType); err != nil {
			return err
		}
		if err := set("state", image.State); err != nil {
			return err
		}
		if err := set("block_device_mappings", omiOAPIBlockDeviceMappings(*image.BlockDeviceMappings)); err != nil {
			return err
		}
		if err := set("product_codes", image.ProductCodes); err != nil {
			return err
		}
		if err := set("state_comment", omiOAPIStateReason(image.StateComment)); err != nil {
			return err
		}
		if err := set("permissions_to_launch", omiOAPIPermissionToLuch(image.PermissionsToLaunch)); err != nil {
			return err
		}
		if err := set("tags", getOapiTagSet(image.Tags)); err != nil {
			return err
		}

		return d.Set("request_id", resp.ResponseContext.RequestId)
	})
}

func omiOAPIPermissionToLuch(p *oscgo.PermissionsOnResource) (res []map[string]interface{}) {
	for _, v := range *p.AccountIds {
		res = append(res, map[string]interface{}{
			"account_id":        v,
			"global_permission": cast.ToString(p.GlobalPermission),
		})
	}
	return
}
