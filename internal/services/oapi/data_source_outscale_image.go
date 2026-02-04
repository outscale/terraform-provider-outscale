package oapi

import (
	"context"
	"fmt"
	"time"

	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/samber/lo"
)

func DataSourceOutscaleImage() *schema.Resource {
	return &schema.Resource{
		Read: DataSourceOutscaleImageRead,

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
			"boot_modes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"tpm_mandatory": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"secure_boot": {
				Type:     schema.TypeBool,
				Computed: true,
			},
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
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"bsu": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"delete_on_vm_deletion": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"iops": {
										Type:     schema.TypeInt,
										Optional: true,
										Computed: true,
									},
									"snapshot_id": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"volume_size": {
										Type:     schema.TypeInt,
										Optional: true,
										Computed: true,
									},
									"volume_type": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
								},
							},
						},
						"device_name": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"virtual_device_name": {
							Type:     schema.TypeString,
							Optional: true,
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
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"state_code": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"state_message": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"permissions_to_launch": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"global_permission": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"account_id": {
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

func DataSourceOutscaleImageRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*client.OutscaleClient).OSCAPI

	filters, filtersOk := d.GetOk("filter")
	executableUsers, executableUsersOk := d.GetOk("permission")
	ai, aisOk := d.GetOk("account_id")
	imageID, imageIDOk := d.GetOk("image_id")
	if !executableUsersOk && !filtersOk && !aisOk && !imageIDOk {
		return fmt.Errorf("one of executable_users, filters, or account_id must be assigned, or image_id must be provided")
	}

	var err error
	filtersReq := &oscgo.FiltersImage{}
	if filtersOk {
		filtersReq, err = buildOutscaleDataSourceImagesFilters(filters.(*schema.Set))
		if err != nil {
			return err
		}
	}
	if imageIDOk {
		filtersReq.SetImageIds([]string{imageID.(string)})
	}
	if aisOk {
		filtersReq.SetAccountIds([]string{ai.(string)})
	}
	if executableUsersOk {
		filtersReq.SetPermissionsToLaunchAccountIds(utils.InterfaceSliceToStringSlice(executableUsers.([]interface{})))
	}

	req := oscgo.ReadImagesRequest{Filters: filtersReq}

	var resp oscgo.ReadImagesResponse
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		rp, httpResp, err := conn.ImageApi.ReadImages(context.Background()).ReadImagesRequest(req).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		resp = rp
		return nil
	})
	if err != nil {
		return err
	}

	images := resp.GetImages()

	if len(images) < 1 {
		return ErrNoResults
	}
	if len(images) > 1 {
		return ErrMultipleResults
	}

	return resourceDataAttrSetter(d, func(set AttributeSetter) error {
		image := images[0]
		d.SetId(*image.ImageId)

		if err := set("architecture", image.Architecture); err != nil {
			return err
		}

		if err := set("boot_modes", lo.Map(image.GetBootModes(), func(b oscgo.BootMode, _ int) string { return string(b) })); err != nil {
			return err
		}
		if err := set("secure_boot", image.GetSecureBoot()); err != nil {
			return err
		}
		if err := set("tpm_mandatory", image.GetTpmMandatory()); err != nil {
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
		if err := set("tags", FlattenOAPITagsSDK(image.GetTags())); err != nil {
			return err
		}

		return nil
	})
}

func omiOAPIPermissionToLuch(p *oscgo.PermissionsOnResource) (res []map[string]interface{}) {
	for _, v := range *p.AccountIds {
		res = append(res, map[string]interface{}{
			"account_id":        v,
			"global_permission": p.GetGlobalPermission(),
		})
	}
	return
}
