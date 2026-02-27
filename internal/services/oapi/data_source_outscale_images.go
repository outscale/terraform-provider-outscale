package oapi

import (
	"context"
	"time"

	"github.com/outscale/osc-sdk-go/v3/pkg/options"
	"github.com/outscale/osc-sdk-go/v3/pkg/osc"
	"github.com/spf13/cast"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/framework/fwhelpers/from"
	"github.com/outscale/terraform-provider-outscale/internal/utils"
	"github.com/samber/lo"
)

func DataSourceOutscaleImages() *schema.Resource {
	return &schema.Resource{
		ReadContext: DataSourceOutscaleImagesRead,

		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			"permissions": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"image_ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"account_ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			// Computed values.
			"request_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"images": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"account_alias": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"account_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"architecture": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"boot_modes": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"secure_boot": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"tpm_mandatory": {
							Type:     schema.TypeBool,
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
						"image_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"image_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"image_type": {
							Type:     schema.TypeString,
							Computed: true,
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
						"product_codes": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"root_device_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"root_device_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"state": {
							Type:     schema.TypeString,
							Computed: true,
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
				},
			},
		},
	}
}

func DataSourceOutscaleImagesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.OutscaleClient).OSC

	executableUsers, executableUsersOk := d.GetOk("permissions")
	filters, filtersOk := d.GetOk("filter")
	aids, ownersOk := d.GetOk("account_ids")
	if !executableUsersOk && !filtersOk && !ownersOk {
		return diag.Errorf("one of executable_users, filters, or account_ids must be assigned")
	}

	var err error
	filtersReq := &osc.FiltersImage{}
	if filtersOk {
		filtersReq, err = buildOutscaleDataSourceImagesFilters(filters.(*schema.Set))
		if err != nil {
			return diag.FromErr(err)
		}
	}
	if ownersOk {
		filtersReq.AccountIds = &[]string{aids.(string)}
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
	if images == nil || len(*images) == 0 {
		return diag.FromErr(ErrNoResults)
	}

	return diag.FromErr(resourceDataAttrSetter(d, func(set AttributeSetter) error {
		d.SetId(id.UniqueId())

		imgs := make([]map[string]interface{}, len(*images))
		for i, image := range *images {
			imgs[i] = map[string]interface{}{
				"architecture":          image.Architecture,
				"boot_modes":            lo.Map(image.BootModes, func(b osc.BootMode, _ int) string { return string(b) }),
				"secure_boot":           image.SecureBoot,
				"tpm_mandatory":         image.TpmMandatory,
				"creation_date":         from.ISO8601(image.CreationDate),
				"description":           image.Description,
				"image_id":              image.ImageId,
				"file_location":         image.FileLocation,
				"account_alias":         image.AccountAlias,
				"account_id":            image.AccountId,
				"image_type":            image.ImageType,
				"image_name":            image.ImageName,
				"root_device_name":      image.RootDeviceName,
				"root_device_type":      image.RootDeviceType,
				"state":                 image.State,
				"block_device_mappings": omiOAPIBlockDeviceMappings(*image.BlockDeviceMappings),
				"product_codes":         image.ProductCodes,
				"state_comment":         omiOAPIStateReason(image.StateComment),
				"permissions_to_launch": omiOAPIPermissionToLuch(image.PermissionsToLaunch),
				"tags":                  FlattenOAPITagsSDK(image.Tags),
			}
		}

		return set("images", imgs)
	}))
}

func buildOutscaleDataSourceImagesFilters(set *schema.Set) (*osc.FiltersImage, error) {
	filters := osc.FiltersImage{}
	for _, v := range set.List() {
		m := v.(map[string]interface{})
		var filterValues []string

		for _, e := range m["values"].([]interface{}) {
			filterValues = append(filterValues, cast.ToString(e))
		}

		switch name := m["name"].(string); name {
		case "account_aliases":
			filters.AccountAliases = &filterValues
		case "account_ids":
			filters.AccountIds = &filterValues
		case "architectures":
			filters.Architectures = &filterValues
		case "boot_modes":
			filters.BootModes = new(lo.Map(filterValues, func(s string, _ int) osc.BootMode { return (osc.BootMode)(s) }))
		case "secure_boot":
			filters.SecureBoot = new(cast.ToBool(filterValues[0]))
		case "block_device_mapping_delete_on_vm_deletion":
			filters.BlockDeviceMappingDeleteOnVmDeletion = new(cast.ToBool(filterValues[0]))
		case "block_device_mapping_device_names":
			filters.BlockDeviceMappingDeviceNames = &filterValues
		case "block_device_mapping_snapshot_ids":
			filters.BlockDeviceMappingSnapshotIds = &filterValues
		case "block_device_mapping_volume_sizes":
			filters.BlockDeviceMappingVolumeSizes = new(utils.StringSliceToIntSlice(filterValues))
		case "block_device_mapping_volume_types":
			filters.BlockDeviceMappingVolumeTypes = new(lo.Map(filterValues, func(s string, _ int) osc.VolumeType { return (osc.VolumeType)(s) }))
		case "descriptions":
			filters.Descriptions = &filterValues
		case "file_locations":
			filters.FileLocations = &filterValues
		case "hypervisors":
			filters.Hypervisors = &filterValues
		case "image_ids":
			filters.ImageIds = &filterValues
		case "image_names":
			filters.ImageNames = &filterValues
		case "permissions_to_launch_account_ids":
			filters.PermissionsToLaunchAccountIds = &filterValues
		case "permissions_to_launch_global_permission":
			filters.PermissionsToLaunchGlobalPermission = new(cast.ToBool(filterValues[0]))
		case "product_codes":
			filters.ProductCodes = &filterValues
		case "product_code_names":
			filters.ProductCodeNames = &filterValues
		case "root_device_names":
			filters.RootDeviceNames = &filterValues
		case "root_device_types":
			filters.RootDeviceTypes = &filterValues
		case "states":
			filters.States = new(lo.Map(filterValues, func(s string, _ int) osc.ImageState { return (osc.ImageState)(s) }))
		case "tag_keys":
			filters.TagKeys = &filterValues
		case "tag_values":
			filters.TagValues = &filterValues
		case "tags":
			filters.Tags = &filterValues
		case "virtualization_types":
			filters.VirtualizationTypes = &filterValues
		default:
			return nil, utils.UnknownDataSourceFilterError(name)
		}
	}
	return &filters, nil
}
