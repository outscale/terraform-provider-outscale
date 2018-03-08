package outscale

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-outscale/osc/fcu"
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
	conn := meta.(*OutscaleClient).FCU

	executableUsers, executableUsersOk := d.GetOk("permissions")
	filters, filtersOk := d.GetOk("filter")
	account_ids, ownersOk := d.GetOk("account_ids")

	if executableUsersOk == false && filtersOk == false && ownersOk == false {
		return fmt.Errorf("One of executable_users, filters, or account_ids must be assigned")
	}

	params := &fcu.DescribeImagesInput{}
	if executableUsersOk {
		params.ExecutableUsers = expandStringList(executableUsers.([]interface{}))
	}
	if filtersOk {
		params.Filters = buildOutscaleDataSourceFilters(filters.(*schema.Set))
	}
	if ownersOk {
		o := expandStringList(account_ids.([]interface{}))

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

	return omisOAPIDescriptionAttributes(d, res.Images)
}

// populate the numerous fields that the image description returns.
func omisOAPIDescriptionAttributes(d *schema.ResourceData, images []*fcu.Image) error {

	i := make([]interface{}, len(images))

	for k, v := range images {
		im := make(map[string]interface{})

		im["architecture"] = *v.Architecture
		if v.CreationDate != nil {
			im["creation_date"] = *v.CreationDate
		}
		if v.Description != nil {
			im["description"] = *v.Description
		}
		im["image_id"] = *v.ImageId
		im["osu_location"] = *v.ImageLocation
		if v.ImageOwnerAlias != nil {
			im["account_alias"] = *v.ImageOwnerAlias
		}
		im["account_id"] = *v.OwnerId
		im["type"] = *v.ImageType
		im["state"] = *v.State
		im["name"] = *v.Name
		im["is_public"] = *v.Public
		if v.RootDeviceName != nil {
			im["root_device_name"] = *v.RootDeviceName
		}
		im["root_device_type"] = *v.RootDeviceType

		if v.BlockDeviceMappings != nil {
			im["block_device_mappings"] = omiOAPIBlockDeviceMappings(v.BlockDeviceMappings)
		}
		if v.ProductCodes != nil {
			im["product_codes"] = omiOAPIProductCodes(v.ProductCodes)
		}
		if v.StateReason != nil {
			im["state_comment"] = omiOAPIStateReason(v.StateReason)
		}
		if v.Tags != nil {
			im["tag"] = dataSourceTags(v.Tags)
		}
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
