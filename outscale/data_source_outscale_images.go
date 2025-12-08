package outscale

import (
	"context"
	"fmt"
	"time"

	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/spf13/cast"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/outscale/terraform-provider-outscale/utils"
)

func DataSourceOutscaleImages() *schema.Resource {
	return &schema.Resource{
		Read: DataSourceOutscaleImagesRead,

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

func DataSourceOutscaleImagesRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	executableUsers, executableUsersOk := d.GetOk("permissions")
	filters, filtersOk := d.GetOk("filter")
	aids, ownersOk := d.GetOk("account_ids")
	if !executableUsersOk && !filtersOk && !ownersOk {
		return fmt.Errorf("One of executable_users, filters, or account_ids must be assigned")
	}

	var err error
	filtersReq := &oscgo.FiltersImage{}
	if filtersOk {
		filtersReq, err = buildOutscaleDataSourceImagesFilters(filters.(*schema.Set))
		if err != nil {
			return err
		}
	}
	if ownersOk {
		filtersReq.SetAccountIds([]string{aids.(string)})
	}
	if executableUsersOk {
		filtersReq.SetPermissionsToLaunchAccountIds(utils.InterfaceSliceToStringSlice(executableUsers.([]interface{})))
	}

	req := oscgo.ReadImagesRequest{Filters: filtersReq}

	var resp oscgo.ReadImagesResponse
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
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
	if len(images) == 0 {
		return fmt.Errorf("your query returned no results, please change your search criteria and try again")
	}

	return resourceDataAttrSetter(d, func(set AttributeSetter) error {
		d.SetId(id.UniqueId())

		imgs := make([]map[string]interface{}, len(images))
		for i, image := range images {
			imgs[i] = map[string]interface{}{
				"architecture":          image.GetArchitecture(),
				"boot_modes":            utils.Map(image.GetBootModes(), func(b oscgo.BootMode) string { return string(b) }),
				"secure_boot":           image.GetSecureBoot(),
				"creation_date":         image.GetCreationDate(),
				"description":           image.GetDescription(),
				"image_id":              image.GetImageId(),
				"file_location":         image.GetFileLocation(),
				"account_alias":         image.GetAccountAlias(),
				"account_id":            image.GetAccountId(),
				"image_type":            image.GetImageType(),
				"image_name":            image.GetImageName(),
				"root_device_name":      image.GetRootDeviceName(),
				"root_device_type":      image.GetRootDeviceType(),
				"state":                 image.GetState(),
				"block_device_mappings": omiOAPIBlockDeviceMappings(*image.BlockDeviceMappings),
				"product_codes":         image.GetProductCodes(),
				"state_comment":         omiOAPIStateReason(image.StateComment),
				"permissions_to_launch": omiOAPIPermissionToLuch(image.PermissionsToLaunch),
				"tags":                  flattenOAPITagsSDK(image.GetTags()),
			}
		}

		return set("images", imgs)
	})
}

func buildOutscaleDataSourceImagesFilters(set *schema.Set) (*oscgo.FiltersImage, error) {
	filters := oscgo.FiltersImage{}
	for _, v := range set.List() {
		m := v.(map[string]interface{})
		var filterValues []string

		for _, e := range m["values"].([]interface{}) {
			filterValues = append(filterValues, cast.ToString(e))
		}

		switch name := m["name"].(string); name {
		case "account_aliases":
			filters.SetAccountAliases(filterValues)
		case "account_ids":
			filters.SetAccountIds(filterValues)
		case "architectures":
			filters.SetArchitectures(filterValues)
		case "boot_modes":
			filters.SetBootModes(utils.Map(filterValues, func(s string) oscgo.BootMode { return (oscgo.BootMode)(s) }))
		case "secure_boot":
			filters.SetSecureBoot(cast.ToBool(filterValues[0]))
		case "block_device_mapping_delete_on_vm_deletion":
			filters.SetBlockDeviceMappingDeleteOnVmDeletion(cast.ToBool(filterValues[0]))
		case "block_device_mapping_device_names":
			filters.SetBlockDeviceMappingDeviceNames(filterValues)
		case "block_device_mapping_snapshot_ids":
			filters.SetBlockDeviceMappingSnapshotIds(filterValues)
		case "block_device_mapping_volume_sizes":
			filters.SetBlockDeviceMappingVolumeSizes(utils.StringSliceToInt32Slice(filterValues))
		case "block_device_mapping_volume_types":
			filters.SetBlockDeviceMappingVolumeTypes(filterValues)
		case "descriptions":
			filters.SetDescriptions(filterValues)
		case "file_locations":
			filters.SetFileLocations(filterValues)
		case "hypervisors":
			filters.SetHypervisors(filterValues)
		case "image_ids":
			filters.SetImageIds(filterValues)
		case "image_names":
			filters.SetImageNames(filterValues)
		case "permissions_to_launch_account_ids":
			filters.SetPermissionsToLaunchAccountIds(filterValues)
		case "permissions_to_launch_global_permission":
			filters.SetPermissionsToLaunchGlobalPermission(cast.ToBool(filterValues[0]))
		case "product_codes":
			filters.SetProductCodes(filterValues)
		case "product_code_names":
			filters.SetProductCodeNames(filterValues)
		case "root_device_names":
			filters.SetRootDeviceNames(filterValues)
		case "root_device_types":
			filters.SetRootDeviceTypes(filterValues)
		case "states":
			filters.SetStates(filterValues)
		case "tag_keys":
			filters.SetTagKeys(filterValues)
		case "tag_values":
			filters.SetTagValues(filterValues)
		case "tags":
			filters.SetTags(filterValues)
		case "virtualization_types":
			filters.SetVirtualizationTypes(filterValues)
		default:
			return nil, utils.UnknownDataSourceFilterError(context.Background(), name)
		}
	}
	return &filters, nil
}
