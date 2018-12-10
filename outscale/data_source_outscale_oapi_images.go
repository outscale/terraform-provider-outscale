package outscale

import (
	"bytes"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-outscale/osc/oapi"
	"github.com/terraform-providers/terraform-provider-outscale/utils"
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
			"image": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
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
						"image_id": {
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
						// Complex computed values
						"block_device_mappings": {
							Type:     schema.TypeSet,
							Computed: true,
							Set:      omiOAPIBlockDeviceMappingHash,
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
							Type:     schema.TypeSet,
							Computed: true,
							Set:      omiOAPIProductCodesHash,
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
				},
			},
		},
	}
}

func dataSourceOutscaleOAPIImagesRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OAPI

	executableUsers, executableUsersOk := d.GetOk("permissions")
	filters, filtersOk := d.GetOk("filter")
	aids, ownersOk := d.GetOk("account_ids")

	if executableUsersOk == false && filtersOk == false && ownersOk == false {
		return fmt.Errorf("One of executable_users, filters, or account_ids must be assigned")
	}

	params := &oapi.ReadImagesRequest{
		Filters: oapi.FiltersImage{},
	}
	if executableUsersOk {
		params.Filters.PermissionsToLaunchAccountIds = expandStringValueList(executableUsers.([]interface{}))
	}

	if filtersOk {
		params.Filters = buildOutscaleOAPIDataSourceImagesFilters(filters.(*schema.Set))
	}
	if ownersOk {
		o := expandStringValueList(aids.([]interface{}))

		if len(o) > 0 {
			params.Filters.AccountIds = o
		}
	}

	var result *oapi.ReadImagesResponse
	var resp *oapi.POST_ReadImagesResponses
	var err error
	err = resource.Retry(20*time.Minute, func() *resource.RetryError {
		resp, err = conn.POST_ReadImages(*params)

		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded") {
				fmt.Printf("[INFO] Request limit exceeded")
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}

		return nil
	})

	var errString string

	if err != nil || resp.OK == nil {
		if err != nil {
			errString = err.Error()
		} else if resp.Code401 != nil {
			errString = fmt.Sprintf("ErrorCode: 401, %s", utils.ToJSONString(resp.Code401))
		} else if resp.Code400 != nil {
			errString = fmt.Sprintf("ErrorCode: 400, %s", utils.ToJSONString(resp.Code400))
		} else if resp.Code500 != nil {
			errString = fmt.Sprintf("ErrorCode: 500, %s", utils.ToJSONString(resp.Code500))
		}

		return fmt.Errorf("Error retrieving Outscale Images: %s", errString)
	}

	result = resp.OK

	if len(result.Images) < 1 {
		return fmt.Errorf("your query returned no results, please change your search criteria and try again")
	}

	return omisOAPIDescriptionAttributes(d, result.Images)
}

// populate the numerous fields that the image description returns.
func omisOAPIDescriptionAttributes(d *schema.ResourceData, images []oapi.Image) error {

	i := make([]interface{}, len(images))

	for k, v := range images {
		im := make(map[string]interface{})

		im["architecture"] = v.Architecture
		if v.CreationDate != "" {
			im["creation_date"] = v.CreationDate
		}
		if v.Description != "" {
			im["description"] = v.Description
		}
		im["image_id"] = v.ImageId
		im["osu_location"] = v.FileLocation
		if v.AccountAlias != "" {
			im["account_alias"] = v.AccountAlias
		}
		im["account_id"] = v.AccountId
		im["type"] = v.ImageType
		im["state"] = v.State
		im["name"] = v.ImageName
		//im["is_public"] = v.Public
		if v.RootDeviceName != "" {
			im["root_device_name"] = v.RootDeviceName
		}
		im["root_device_type"] = v.RootDeviceType

		if v.BlockDeviceMappings != nil {
			im["block_device_mappings"] = omiOAPIBlockDeviceMappings(v.BlockDeviceMappings)
		}
		if v.ProductCodes != nil {
			im["product_codes"] = v.ProductCodes
		}
		//if v.StateComment != nil {
		im["state_comment"] = omiOAPIStateReason(&v.StateComment)
		//}
		// if v.Tags != nil {
		// 	im["tag"] = dataSourceTags(v.Tags)
		// }
		i[k] = im
	}

	err := d.Set("image", i)
	d.SetId(resource.UniqueId())

	return err
}

func omiOAPIBlockDeviceMappingHash(v interface{}) int {
	var buf bytes.Buffer
	// All keys added in alphabetical order.
	m := v.(map[string]interface{})
	buf.WriteString(fmt.Sprintf("%s-", m["device_name"].(string)))
	if d, ok := m["bsu"]; ok {
		if len(d.(map[string]interface{})) > 0 {
			e := d.(map[string]interface{})
			buf.WriteString(fmt.Sprintf("%s-", e["delete_on_vm_termination"].(string)))
			buf.WriteString(fmt.Sprintf("%s-", e["iops"].(string)))
			buf.WriteString(fmt.Sprintf("%s-", e["volume_size"].(string)))
			buf.WriteString(fmt.Sprintf("%s-", e["type"].(string)))
		}
	}
	if d, ok := m["no_device"]; ok {
		buf.WriteString(fmt.Sprintf("%s-", d.(string)))
	}
	if d, ok := m["virtual_device_name"]; ok {
		buf.WriteString(fmt.Sprintf("%s-", d.(string)))
	}
	if d, ok := m["snapshot_id"]; ok {
		buf.WriteString(fmt.Sprintf("%s-", d.(string)))
	}
	return hashcode.String(buf.String())
}

func buildOutscaleOAPIDataSourceImagesFilters(set *schema.Set) oapi.FiltersImage {
	var filters oapi.FiltersImage
	for _, v := range set.List() {
		m := v.(map[string]interface{})
		var filterValues []string
		for _, e := range m["values"].([]interface{}) {
			filterValues = append(filterValues, e.(string))
		}

		switch name := m["name"].(string); name {
		case "account_aliases":
			filters.AccountAliases = filterValues
		case "account_ids":
			filters.AccountIds = filterValues
		case "architectures":
			filters.Architectures = filterValues
		case "image_ids":
			filters.ImageIds = filterValues
		case "image_names":
			filters.ImageNames = filterValues
		case "image_types":
			filters.ImageTypes = filterValues
		case "virtualization_types":
			filters.VirtualizationTypes = filterValues
		case "root_device_types":
			filters.RootDeviceTypes = filterValues
		case "block_device_mapping_volume_type":
			filters.BlockDeviceMappingVolumeType = filterValues
		//Some params are missing.
		default:
			log.Printf("[Debug] Unknown Filter Name: %s.", name)
		}
	}
	return filters
}
