package outscale

import (
	"bytes"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/spf13/cast"

	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-outscale/osc/oapi"
	"github.com/terraform-providers/terraform-provider-outscale/utils"
)

func resourceOutscaleOAPIImage() *schema.Resource {
	return &schema.Resource{
		Create: resourceOAPIImageCreate,
		Read:   resourceOAPIImageRead,
		Update: resourceOAPIImageUpdate,
		Delete: resourceOAPIImageDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(40 * time.Minute),
			Update: schema.DefaultTimeout(40 * time.Minute),
			Delete: schema.DefaultTimeout(40 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"vm_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"image_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"no_reboot": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
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
			"image_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"file_location": {
				Type:     schema.TypeString,
				Optional: true,
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
			"is_public": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"root_device_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"root_device_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			// Complex computed values
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
							Type:     schema.TypeBool,
							Computed: true,
						},
						"account_ids": &schema.Schema{
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
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
			"request_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceOAPIImageCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OAPI

	req := &oapi.CreateImageRequest{
		ImageName:      d.Get("image_name").(string),
		VmId:           d.Get("vm_id").(string),
		Description:    d.Get("description").(string),
		NoReboot:       d.Get("no_reboot").(bool),
		Architecture:   d.Get("architecture").(string),
		FileLocation:   d.Get("file_location").(string),
		SourceImageId:  d.Get("image_id").(string),
		RootDeviceName: d.Get("root_device_name").(string),
	}

	var result *oapi.CreateImageResponse
	var resp *oapi.POST_CreateImageResponses
	var err error
	err = resource.Retry(40*time.Minute, func() *resource.RetryError {
		resp, err = conn.POST_CreateImage(*req)

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

		return fmt.Errorf("Error creating Outscale Image: %s", errString)
	}

	result = resp.OK

	id := result.Image.ImageId
	d.SetId(id)
	d.Set("image_id", id)
	d.Partial(true) // make sure we record the id even if the rest of this gets interrupted
	d.Set("id", id)
	d.SetPartial("id")
	d.Partial(false)

	_, err = resourceOutscaleOAPIImageWaitForAvailable(id, conn, 1)
	if err != nil {
		return err
	}

	d.Set("description", result.Image.Description)
	d.Set("creation_date", result.Image.CreationDate)

	return resourceOAPIImageRead(d, meta)

}

func resourceOAPIImageRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*OutscaleClient).OAPI
	id := d.Id()

	req := &oapi.ReadImagesRequest{
		Filters: oapi.FiltersImage{ImageIds: []string{id}},
	}

	var resp *oapi.POST_ReadImagesResponses
	var err error
	err = resource.Retry(40*time.Minute, func() *resource.RetryError {
		resp, err = client.POST_ReadImages(*req)

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
			if strings.Contains(err.Error(), "InvalidAMIID.NotFound") {
				fmt.Printf("[DEBUG] %s no longer exists, so we'll drop it from the state", id)
				d.SetId("")
				return nil
			}
			errString = err.Error()
		} else if resp.Code401 != nil {
			errString = fmt.Sprintf("ErrorCode: 401, %s", utils.ToJSONString(resp.Code401))
		} else if resp.Code400 != nil {
			errString = fmt.Sprintf("ErrorCode: 400, %s", utils.ToJSONString(resp.Code400))
		} else if resp.Code500 != nil {
			errString = fmt.Sprintf("ErrorCode: 500, %s", utils.ToJSONString(resp.Code500))
		}

		return fmt.Errorf("Error creating Outscale VM volume: %s", errString)
	}

	result := resp.OK

	if len(result.Images) != 1 {
		d.SetId("")
		return nil
	}

	image := result.Images[0]
	state := image.State

	if state == "pending" {
		var img *oapi.Image
		img, err = resourceOutscaleOAPIImageWaitForAvailable(id, client, 2)
		if err != nil {
			return err
		}

		image = *img

		state = image.State
	}

	if state == "deregistered" {
		d.SetId("")
		return nil
	}

	if state != "available" {
		return fmt.Errorf("OMI has become %s", state)
	}

	d.SetId(image.ImageId)
	d.Set("architecture", image.Architecture)
	if image.CreationDate != "" {
		d.Set("creation_date", image.CreationDate)
	}
	if image.Description != "" {
		d.Set("description", image.Description)
	}
	//Missing on swager spec
	//d.Set("hypervisor", image.Hypervisor)
	d.Set("image_id", image.ImageId)
	d.Set("file_location", image.FileLocation)
	if image.AccountAlias != "nil" {
		d.Set("account_alias", image.AccountAlias)
	}
	d.Set("account_id", image.AccountId)
	d.Set("image_type", image.ImageType)
	d.Set("image_name", image.ImageName)
	//Missing on swager spec
	// d.Set("is_public", image.Public)
	if image.RootDeviceName != "" {
		d.Set("root_device_name", image.RootDeviceName)
	}
	d.Set("root_device_type", image.RootDeviceType)
	d.Set("state", image.State)

	if err := d.Set("block_device_mappings", omiOAPIBlockDeviceMappings(image.BlockDeviceMappings)); err != nil {
		return err
	}
	if err := d.Set("product_codes", image.ProductCodes); err != nil {
		return err
	}
	if err := d.Set("state_comment", omiOAPIStateReason(&image.StateComment)); err != nil {
		return err
	}

	if err := d.Set("permissions_to_launch", setResourcePermissions(image.PermissionsToLaunch)); err != nil {
		return err
	}
	if err := d.Set("tags", getOapiTagSet(image.Tags)); err != nil {
		return err
	}

	d.Set("request_id", result.ResponseContext.RequestId)

	return nil
}

func setResourcePermissions(por oapi.PermissionsOnResource) []map[string]interface{} {
	lp := make([]map[string]interface{}, 1)
	l := make(map[string]interface{})
	l["global_permission"] = por.GlobalPermission
	l["account_ids"] = por.AccountIds

	lp[0] = l

	return lp
}

func resourceOAPIImageUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OAPI

	d.Partial(true)

	// if err := setOAPITags(conn, d); err != nil {
	// 	return err
	// }

	// d.SetPartial("tags")

	if d.Get("description").(string) != "" {
		_, err := conn.POST_UpdateImage(oapi.UpdateImageRequest{
			ImageId: d.Id(),
			// Description: &oapi.AttributeValue{
			// 	Value: aws.String(d.Get("description").(string)),
			// },
		})
		if err != nil {
			return err
		}
		d.SetPartial("description")
	}

	d.Partial(false)

	return resourceOAPIImageRead(d, meta)
}

func resourceOAPIImageDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*OutscaleClient).OAPI

	req := &oapi.DeleteImageRequest{
		ImageId: d.Id(),
	}

	var err error
	err = resource.Retry(40*time.Minute, func() *resource.RetryError {
		_, err := client.POST_DeleteImage(*req)

		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded") {
				fmt.Printf("[INFO] Request limit exceeded")
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("Error deleting the image")
	}

	if err := resourceOutscaleOAPIImageWaitForDestroy(d.Id(), client); err != nil {
		return err
	}

	d.SetId("")
	return nil
}

func resourceOutscaleOAPIImageWaitForAvailable(id string, client *oapi.Client, i int) (*oapi.Image, error) {
	fmt.Printf("Waiting for OMI %s to become available...", id)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"pending"},
		Target:     []string{"available"},
		Refresh:    ImageOAPIStateRefreshFunc(client, id),
		Timeout:    OutscaleImageRetryTimeout,
		Delay:      OutscaleImageRetryDelay,
		MinTimeout: OutscaleImageRetryMinTimeout,
	}

	info, err := stateConf.WaitForState()
	if err != nil {
		return nil, fmt.Errorf("Error waiting for OMI (%s) to be ready: %v", id, err)
	}

	img := info.(oapi.Image)

	return &img, nil
}

// ImageOAPIStateRefreshFunc ...
func ImageOAPIStateRefreshFunc(client *oapi.Client, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		emptyResp := &oapi.ReadImagesResponse{}
		var result *oapi.ReadImagesResponse
		var resp *oapi.POST_ReadImagesResponses
		var err error
		err = resource.Retry(15*time.Minute, func() *resource.RetryError {
			request := &oapi.ReadImagesRequest{
				Filters: oapi.FiltersImage{
					ImageIds: []string{id},
				},
			}
			resp, err = client.POST_ReadImages(*request)

			if err != nil {
				if strings.Contains(err.Error(), "RequestLimitExceeded") {
					log.Printf("[INFO] Request limit exceeded")
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)

			}

			return nil
		})

		var errString string

		if err != nil || resp.OK == nil {
			if err != nil {
				if e := fmt.Sprint(err); strings.Contains(e, "InvalidAMIID.NotFound") {
					log.Printf("[INFO] OMI %s state %s", id, "destroyed")
					return emptyResp, "destroyed", nil
				}

				errString = err.Error()
			} else if resp.Code401 != nil {
				errString = fmt.Sprintf("ErrorCode: 401, %s", utils.ToJSONString(resp.Code401))
			} else if resp.Code400 != nil {
				errString = fmt.Sprintf("ErrorCode: 400, %s", utils.ToJSONString(resp.Code400))
			} else if resp.Code500 != nil {
				errString = fmt.Sprintf("ErrorCode: 500, %s", utils.ToJSONString(resp.Code500))
			}

			return emptyResp, "", fmt.Errorf("Error refreshing image state: %s", errString)
		}

		result = resp.OK

		if result != nil && len(result.Images) == 0 {
			log.Printf("[INFO] OMI %s state %s", id, "destroyed")
			return emptyResp, "destroyed", nil
		}

		if result == nil || result.Images == nil || len(result.Images) == 0 {
			return emptyResp, "destroyed", nil
		}

		log.Printf("[INFO] OMI %s state %s", result.Images[0].ImageId, result.Images[0].State)

		// OMI is valid, so return it's state
		return result.Images[0], result.Images[0].State, nil
	}
}

func resourceOutscaleOAPIImageWaitForDestroy(id string, client *oapi.Client) error {
	fmt.Printf("Waiting for OMI %s to be deleted...", id)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"available", "pending", "failed"},
		Target:     []string{"destroyed"},
		Refresh:    ImageOAPIStateRefreshFunc(client, id),
		Timeout:    OutscaleImageDeleteRetryTimeout,
		Delay:      OutscaleImageRetryDelay,
		MinTimeout: OutscaleImageRetryTimeout,
	}

	_, err := stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Error waiting for OMI (%s) to be deleted: %v", id, err)
	}

	return nil
}

// Returns a set of block device mappings.
func omiOAPIBlockDeviceMappings(m []oapi.BlockDeviceMappingImage) []map[string]interface{} {
	blockDeviceMapping := make([]map[string]interface{}, len(m))

	for k, v := range m {
		blockDeviceMapping[k] = map[string]interface{}{
			"device_name":         v.DeviceName,
			"virtual_device_name": v.VirtualDeviceName,
			"bsu": map[string]interface{}{
				"delete_on_vm_deletion": fmt.Sprintf("%t", *v.Bsu.DeleteOnVmDeletion),
				"iops":                  cast.ToString(v.Bsu.Iops),
				"snapshot_id":           v.Bsu.SnapshotId,
				"volume_size":           cast.ToString(v.Bsu.VolumeSize),
				"volume_type":           v.Bsu.VolumeType,
			},
		}
	}
	return blockDeviceMapping
}

// Returns a set of product codes.
func omiOAPIProductCodes(m []string) *schema.Set {
	s := &schema.Set{
		F: omiOAPIProductCodesHash,
	}
	for _, v := range m {
		code := map[string]interface{}{
			"product_code": v,
			"type":         "UNSET",
		}
		s.Add(code)
	}
	return s
}

// Returns the state reason.
func omiOAPIStateReason(m *oapi.StateComment) map[string]interface{} {
	s := make(map[string]interface{})
	if m != nil {
		s["state_code"] = m.StateCode
		s["state_message"] = m.StateMessage
	} else {
		s["state_code"] = "UNSET"
		s["state_message"] = "UNSET"
	}
	return s
}

func omiOAPIProductCodesHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})
	buf.WriteString(fmt.Sprintf("%s-", m["product_code"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["type"].(string)))
	return hashcode.String(buf.String())
}
