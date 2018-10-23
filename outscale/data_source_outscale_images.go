package outscale

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-outscale/osc/fcu"
)

func dataSourceOutscaleImages() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOutscaleImagesRead,

		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			"executable_by": {
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
			"owners": {
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
			"images_set": {
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
				},
			},
		},
	}
}

func dataSourceOutscaleImagesRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	executableUsers, executableUsersOk := d.GetOk("executable_by")
	filters, filtersOk := d.GetOk("filter")
	owners, ownersOk := d.GetOk("owners")

	if executableUsersOk == false && filtersOk == false && ownersOk == false {
		return fmt.Errorf("One of executable_users, filters, or owners must be assigned")
	}

	params := &fcu.DescribeImagesInput{}
	if executableUsersOk {
		params.ExecutableUsers = expandStringList(executableUsers.([]interface{}))
	}
	if filtersOk {
		params.Filters = buildOutscaleDataSourceFilters(filters.(*schema.Set))
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
		return fmt.Errorf("your query returned no results, please change your search criteria and try again")
	}

	d.Set("request_id", res.RequestId)

	return omisDescriptionAttributes(d, res.Images)
}

// populate the numerous fields that the image description returns.
func omisDescriptionAttributes(d *schema.ResourceData, images []*fcu.Image) error {

	i := make([]interface{}, len(images))

	for k, v := range images {
		im := make(map[string]interface{})

		im["architecture"] = *v.Architecture
		if v.CreationDate != nil {
			im["creation_date"] = *v.CreationDate
		} else {
			im["creation_date"] = ""
		}
		if v.Description != nil {
			im["description"] = *v.Description
		} else {
			im["description"] = ""
		}
		im["image_id"] = *v.ImageId
		im["image_location"] = *v.ImageLocation
		if v.ImageOwnerAlias != nil {
			im["image_owner_alias"] = *v.ImageOwnerAlias
		} else {
			im["image_owner_alias"] = ""
		}
		im["image_owner_id"] = *v.OwnerId
		im["image_type"] = *v.ImageType
		im["image_state"] = *v.State
		im["name"] = *v.Name
		im["is_public"] = *v.Public
		if v.RootDeviceName != nil {
			im["root_device_name"] = *v.RootDeviceName
		} else {
			im["root_device_name"] = ""
		}
		im["root_device_type"] = *v.RootDeviceType

		if v.BlockDeviceMappings != nil {
			im["block_device_mapping"] = amiBlockDeviceMappings(v.BlockDeviceMappings)
		}
		if v.ProductCodes != nil {
			im["product_codes"] = amiProductCodes(v.ProductCodes)
		}
		if v.StateReason != nil {
			im["state_reason"] = amiStateReason(v.StateReason)
		}
		if v.Tags != nil {
			im["tag_set"] = dataSourceTags(v.Tags)
		}
		i[k] = im
	}

	err := d.Set("images_set", i)
	d.SetId(resource.UniqueId())

	return err
}

// Returns a set of block device mappings.
func amiBlockDeviceMappings(m []*fcu.BlockDeviceMapping) []map[string]interface{} {
	s := make([]map[string]interface{}, len(m))

	for k, v := range m {
		mapping := make(map[string]interface{})
		if v.Ebs != nil {
			ebs := map[string]interface{}{
				"delete_on_termination": fmt.Sprintf("%t", *v.Ebs.DeleteOnTermination),
				"volume_size":           fmt.Sprintf("%d", *v.Ebs.VolumeSize),
				"volume_type":           *v.Ebs.VolumeType,
			}

			if v.Ebs.Iops != nil {
				ebs["iops"] = fmt.Sprintf("%d", *v.Ebs.Iops)
			} else {
				ebs["iops"] = "0"
			}
			// snapshot id may not be set
			if v.Ebs.SnapshotId != nil {
				ebs["snapshot_id"] = *v.Ebs.SnapshotId
			}

			mapping["ebs"] = ebs
		}
		if v.VirtualName != nil {
			mapping["virtual_name"] = *v.VirtualName
		}
		if v.DeviceName != nil {
			mapping["device_name"] = *v.DeviceName
		}
		if v.NoDevice != nil {
			mapping["no_device"] = *v.NoDevice
		}

		s[k] = mapping
	}
	return s
}

// Returns a set of product codes.
func amiProductCodes(m []*fcu.ProductCode) *schema.Set {
	s := &schema.Set{
		F: amiProductCodesHash,
	}
	for _, v := range m {
		code := map[string]interface{}{
			"product_code": *v.ProductCode,
			"type":         *v.Type,
		}
		s.Add(code)
	}
	return s
}

// Returns the state reason.
func amiStateReason(m *fcu.StateReason) map[string]interface{} {
	s := make(map[string]interface{})
	if m != nil {
		s["code"] = *m.Code
		s["message"] = *m.Message
	} else {
		s["code"] = "UNSET"
		s["message"] = "UNSET"
	}
	return s
}

// Generates a hash for the set hash function used by the block_device_mapping
// attribute.
func amiBlockDeviceMappingHash(v interface{}) int {
	var buf bytes.Buffer
	// All keys added in alphabetical order.
	m := v.(map[string]interface{})
	buf.WriteString(fmt.Sprintf("%s-", m["device_name"].(string)))
	if d, ok := m["ebs"]; ok {
		if len(d.(map[string]interface{})) > 0 {
			e := d.(map[string]interface{})
			buf.WriteString(fmt.Sprintf("%s-", e["delete_on_termination"].(string)))
			buf.WriteString(fmt.Sprintf("%s-", e["encrypted"].(string)))
			buf.WriteString(fmt.Sprintf("%s-", e["iops"].(string)))
			buf.WriteString(fmt.Sprintf("%s-", e["volume_size"].(string)))
			buf.WriteString(fmt.Sprintf("%s-", e["volume_type"].(string)))
		}
	}
	if d, ok := m["no_device"]; ok {
		buf.WriteString(fmt.Sprintf("%s-", d.(string)))
	}
	if d, ok := m["virtual_name"]; ok {
		buf.WriteString(fmt.Sprintf("%s-", d.(string)))
	}
	if d, ok := m["snapshot_id"]; ok {
		buf.WriteString(fmt.Sprintf("%s-", d.(string)))
	}
	return hashcode.String(buf.String())
}

// Generates a hash for the set hash function used by the product_codes
// attribute.
func amiProductCodesHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})
	// All keys added in alphabetical order.
	buf.WriteString(fmt.Sprintf("%s-", m["product_code_id"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["product_code_type"].(string)))
	return hashcode.String(buf.String())
}

func dataSourceTagsSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeSet,
		Computed: true,
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
	}
}

func expandStringList(configured []interface{}) []*string {
	vs := make([]*string, 0, len(configured))
	for _, v := range configured {
		val, ok := v.(string)
		if ok && val != "" {
			vs = append(vs, aws.String(v.(string)))
		}
	}
	return vs
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

func dataSourceTagsHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})
	buf.WriteString(fmt.Sprintf("%s-", m["key"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["value"].(string)))
	return hashcode.String(buf.String())
}

func dataSourceTags(m []*fcu.Tag) *schema.Set {
	s := &schema.Set{
		F: dataSourceTagsHash,
	}
	for _, v := range m {
		tag := map[string]interface{}{
			"key":   *v.Key,
			"value": *v.Value,
		}
		s.Add(tag)
	}
	return s
}
