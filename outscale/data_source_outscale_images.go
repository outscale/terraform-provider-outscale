package outscale

import (
	"fmt"
	"log"
	"time"

	"github.com/antihax/optional"

	oscgo "github.com/marinsalinas/osc-sdk-go"
	"github.com/spf13/cast"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceOutscaleOAPIImages() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOutscaleOAPIImagesRead,

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
										Type:     schema.TypeString,
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
							Type:     schema.TypeMap,
							Computed: true,
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

func dataSourceOutscaleOAPIImagesRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	executableUsers, executableUsersOk := d.GetOk("permissions")
	filters, filtersOk := d.GetOk("filter")
	aids, ownersOk := d.GetOk("account_ids")
	if !executableUsersOk && !filtersOk && !ownersOk {
		return fmt.Errorf("One of executable_users, filters, or account_ids must be assigned")
	}

	filtersReq := &oscgo.FiltersImage{}
	if filtersOk {
		filtersReq = buildOutscaleOAPIDataSourceImagesFilters(filters.(*schema.Set))
	}
	if ownersOk {
		filtersReq.SetAccountIds([]string{aids.(string)})
	}
	if executableUsersOk {
		filtersReq.SetPermissionsToLaunchAccountIds(expandStringValueList(executableUsers.([]interface{})))
	}

	req := &oscgo.ReadImagesOpts{
		ReadImagesRequest: optional.NewInterface(oscgo.ReadImagesRequest{
			Filters: filtersReq,
		}),
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"pending"},
		Target:     []string{"available", "destroyed"},
		Refresh:    ImageOAPIStateRefreshFunc(conn, req, "deregistered"),
		Timeout:    5 * time.Minute,
		MinTimeout: 30 * time.Second,
		Delay:      1 * time.Minute,
	}

	value, err := stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Error retrieving Outscale Images: %v", err)
	}

	resp := value.(oscgo.ReadImagesResponse)
	images := resp.GetImages()

	return resourceDataAttrSetter(d, func(set AttributeSetter) error {
		d.SetId(resource.UniqueId())

		imgs := make([]map[string]interface{}, len(images))
		for i, image := range images {
			imgs[i] = map[string]interface{}{
				"architecture":          image.GetArchitecture(),
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
				"tags":                  getOapiTagSet(image.Tags),
			}
		}

		if err := d.Set("request_id", resp.ResponseContext.RequestId); err != nil {
			return err
		}

		return set("images", imgs)
	})
}

func buildOutscaleOAPIDataSourceImagesFilters(set *schema.Set) *oscgo.FiltersImage {
	filters := &oscgo.FiltersImage{}
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
		case "block_device_mapping_delete_on_vm_deletion":
			filters.SetBlockDeviceMappingDeleteOnVmDeletion(cast.ToBool(filterValues))
		case "block_device_mapping_device_names":
			filters.SetBlockDeviceMappingDeleteOnVmDeletion(cast.ToBool(filterValues))
		case "block_device_mapping_snapshot_ids":
			filters.SetBlockDeviceMappingSnapshotIds(filterValues)
		case "block_device_mapping_volume_sizes":
			filters.SetBlockDeviceMappingSnapshotIds(filterValues)
		case "block_device_mapping_volume_type":
			filters.SetBlockDeviceMappingVolumeTypes(filterValues)
		case "description":
			filters.SetDescriptions(filterValues)
		case "file_locations":
			filters.SetFileLocations(filterValues)
		case "image_ids":
			filters.SetImageIds(filterValues)
		case "permissions_to_launch_account_ids":
			filters.SetPermissionsToLaunchAccountIds(filterValues)
		case "permissions_to_launch_global_permission":
			filters.SetPermissionsToLaunchGlobalPermission(cast.ToBool(filterValues))
		case "root_device_names":
			filters.SetRootDeviceNames(filterValues)
		case "root_device_types":
			filters.SetRootDeviceTypes(filterValues)
		case "image_names":
			filters.SetImageNames(filterValues)
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
			log.Printf("[Debug] Unknown Filter Name: %s.", name)
		}
	}
	return filters
}

func expandStringValueList(configured []interface{}) []string {
	vs := make([]string, 0, len(configured))
	for _, v := range configured {
		val, ok := v.(string)
		if ok && val != "" {
			vs = append(vs, v.(string))
		}
	}
	return vs
}

func expandStringValueListPointer(configured []interface{}) *[]string {
	res := expandStringValueList(configured)
	return &res
}
