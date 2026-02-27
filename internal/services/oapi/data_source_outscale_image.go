package oapi

import (
	"context"
	"time"

	"github.com/outscale/goutils/sdk/ptr"
	"github.com/outscale/osc-sdk-go/v3/pkg/options"
	"github.com/outscale/osc-sdk-go/v3/pkg/osc"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/framework/fwhelpers/from"
	"github.com/outscale/terraform-provider-outscale/internal/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/samber/lo"
)

func DataSourceOutscaleImage() *schema.Resource {
	return &schema.Resource{
		ReadContext: DataSourceOutscaleImageRead,

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

func DataSourceOutscaleImageRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.OutscaleClient).OSC

	filters, filtersOk := d.GetOk("filter")
	executableUsers, executableUsersOk := d.GetOk("permission")
	ai, aisOk := d.GetOk("account_id")
	imageID, imageIDOk := d.GetOk("image_id")
	if !executableUsersOk && !filtersOk && !aisOk && !imageIDOk {
		return diag.Errorf("one of executable_users, filters, or account_id must be assigned, or image_id must be provided")
	}

	var err error
	filtersReq := &osc.FiltersImage{}
	if filtersOk {
		filtersReq, err = buildOutscaleDataSourceImagesFilters(filters.(*schema.Set))
		if err != nil {
			return diag.FromErr(err)
		}
	}
	if imageIDOk {
		filtersReq.ImageIds = &[]string{imageID.(string)}
	}
	if aisOk {
		filtersReq.AccountIds = &[]string{ai.(string)}
	}
	if executableUsersOk {
		filtersReq.PermissionsToLaunchAccountIds = utils.InterfaceSliceToStringSlicePtr(executableUsers.([]interface{}))
	}

	req := osc.ReadImagesRequest{Filters: filtersReq}

	resp, err := client.ReadImages(ctx, req, options.WithRetryTimeout(5*time.Minute))
	if err != nil {
		return diag.FromErr(err)
	}

	images := resp.Images

	if images == nil || len(*images) < 1 {
		return diag.FromErr(ErrNoResults)
	}
	if len(*images) > 1 {
		return diag.FromErr(ErrMultipleResults)
	}

	return diag.FromErr(resourceDataAttrSetter(d, func(set AttributeSetter) error {
		image := (*images)[0]
		d.SetId(image.ImageId)

		if err := set("architecture", image.Architecture); err != nil {
			return err
		}

		if err := set("boot_modes", lo.Map(image.BootModes, func(b osc.BootMode, _ int) string { return string(b) })); err != nil {
			return err
		}
		if err := set("secure_boot", image.SecureBoot); err != nil {
			return err
		}
		if err := set("tpm_mandatory", image.TpmMandatory); err != nil {
			return err
		}
		if err := set("creation_date", from.ISO8601(image.CreationDate)); err != nil {
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
		if err := set("block_device_mappings", omiOAPIBlockDeviceMappings(ptr.From(image.BlockDeviceMappings))); err != nil {
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
		if err := set("tags", FlattenOAPITagsSDK(image.Tags)); err != nil {
			return err
		}

		return nil
	}))
}

func omiOAPIPermissionToLuch(p *osc.PermissionsOnResource) (res []map[string]interface{}) {
	for _, v := range *p.AccountIds {
		res = append(res, map[string]interface{}{
			"account_id":        v,
			"global_permission": p.GlobalPermission,
		})
	}
	return
}
