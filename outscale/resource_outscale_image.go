package outscale

import (
	"context"
	"fmt"
	"log"
	"time"

	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/terraform-providers/terraform-provider-outscale/utils"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const (
	// ImageRetryTimeout ...
	ImageRetryTimeout = 40 * time.Minute
	// ImageDeleteRetryTimeout ...
	ImageDeleteRetryTimeout = 90 * time.Minute
	// ImageRetryDelay ...
	ImageRetryDelay = 20 * time.Second
	// ImageRetryMinTimeout ...
	ImageRetryMinTimeout = 3 * time.Second
)

func resourceImage() *schema.Resource {
	return &schema.Resource{
		Create: resourceImageCreate,
		Read:   resourceImageRead,
		Update: resourceImageUpdate,
		Delete: resourceImageDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(40 * time.Minute),
			Update: schema.DefaultTimeout(40 * time.Minute),
			Delete: schema.DefaultTimeout(40 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"architecture": {
				Type:     schema.TypeString,
				Optional: true,
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
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"file_location": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"image_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"no_reboot": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"root_device_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"source_image_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"source_region_name": {
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
			"root_device_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"image_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"creation_date": {
				Type:     schema.TypeString,
				Computed: true,
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
						"account_ids": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"request_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": tagsListSchema(),
			"vm_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceImageCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*Client).OSCAPI

	imageRequest := oscgo.CreateImageRequest{}
	if v, ok := d.GetOk("image_name"); ok {
		imageRequest.SetImageName(v.(string))
	}

	if v, ok := d.GetOk("vm_id"); ok {
		imageRequest.SetVmId(v.(string))
	}

	if v, ok := d.GetOk("description"); ok {
		imageRequest.SetDescription(v.(string))
	}
	if blocks, ok := d.GetOk("block_device_mappings"); ok {
		blockDevices := expandOmiBlockDeviceMappings(blocks.([]interface{}))
		imageRequest.SetBlockDeviceMappings(blockDevices)
	}
	if v, ok := d.GetOk("no_reboot"); ok {
		imageRequest.SetNoReboot(v.(bool))
	}

	if v, ok := d.GetOk("architecture"); ok {
		imageRequest.SetArchitecture(v.(string))
	}

	if v, ok := d.GetOk("file_location"); ok {
		imageRequest.SetFileLocation(v.(string))
	}

	if v, ok := d.GetOk("source_image_id"); ok {
		imageRequest.SetSourceImageId(v.(string))
	}

	if v, ok := d.GetOk("source_region_name"); ok {
		imageRequest.SetSourceRegionName(v.(string))
	}

	if v, ok := d.GetOk("root_device_name"); ok {
		imageRequest.SetRootDeviceName(v.(string))
	}
	var resp oscgo.CreateImageResponse
	var err error
	err = resource.Retry(120*time.Second, func() *resource.RetryError {
		rp, httpResp, err := conn.ImageApi.CreateImage(context.Background()).CreateImageRequest(imageRequest).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp.StatusCode, err)
		}
		resp = rp
		return nil
	})

	if err != nil {
		return err
	}

	if !resp.HasImage() {
		return nil
	}

	image := resp.GetImage()

	log.Printf("[DEBUG] Waiting for OMI %s to become available...", *image.ImageId)

	req := oscgo.ReadImagesRequest{Filters: &oscgo.FiltersImage{ImageIds: &[]string{*image.ImageId}}}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"pending"},
		Target:     []string{"available"},
		Refresh:    ImageStateRefreshFunc(conn, req, "failed"),
		Timeout:    10 * time.Minute,
		MinTimeout: 30 * time.Second,
		Delay:      1 * time.Minute,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Error waiting for OMI (%s) to be ready: %v", *image.ImageId, err)
	}

	if tags, ok := d.GetOk("tags"); ok {
		err := assignTags(tags.(*schema.Set), image.GetImageId(), conn)
		if err != nil {
			return err
		}
	}

	d.SetId(*image.ImageId)

	return resourceImageRead(d, meta)
}

func resourceImageRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*Client).OSCAPI
	id := d.Id()

	req := oscgo.ReadImagesRequest{
		Filters: &oscgo.FiltersImage{ImageIds: &[]string{id}},
	}

	var resp oscgo.ReadImagesResponse
	err := resource.Retry(120*time.Second, func() *resource.RetryError {
		var err error
		rp, httpResp, err := conn.ImageApi.ReadImages(context.Background()).ReadImagesRequest(req).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp.StatusCode, err)
		}
		resp = rp
		return nil
	})

	if err != nil {
		return fmt.Errorf("Error reading for OMI (%s): %v", id, err)
	}

	if len(resp.GetImages()) == 0 {
		d.SetId("")
		return nil
	}

	image := resp.GetImages()[0]

	return resourceDataAttrSetter(d, func(set AttributeSetter) error {
		d.SetId(*image.ImageId)

		if err := set("architecture", image.Architecture); err != nil {
			return err
		}
		if err := set("creation_date", image.CreationDate); err != nil {
			return err
		}
		if err := set("description", image.Description); err != nil {
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
		if err := set("block_device_mappings", omiBlockDeviceMappings(image.GetBlockDeviceMappings())); err != nil {
			return err
		}
		if err := set("product_codes", image.ProductCodes); err != nil {
			return err
		}
		if err := set("state_comment", omiStateReason(image.StateComment)); err != nil {
			return err
		}
		if err := set("permissions_to_launch", setResourcePermissions(*image.PermissionsToLaunch)); err != nil {
			return err
		}
		if err := d.Set("tags", tagsToMap(image.GetTags())); err != nil {
			fmt.Printf("[WARN] ERROR TAGS PROBLEME (%s)", err)
		}

		return nil
	})
}

func setResourcePermissions(por oscgo.PermissionsOnResource) []map[string]interface{} {
	return []map[string]interface{}{
		{
			"global_permission": por.GetGlobalPermission(),
			"account_ids":       por.GetAccountIds(),
		},
	}
}

func resourceImageUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*Client).OSCAPI

	d.Partial(true)
	if err := setTags(conn, d); err != nil {
		return err
	}
	d.SetPartial("tags")

	d.Partial(false)

	return resourceImageRead(d, meta)
}

func resourceImageDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*Client).OSCAPI

	var err error
	err = resource.Retry(120*time.Second, func() *resource.RetryError {
		_, httpResp, err := conn.ImageApi.DeleteImage(context.Background()).DeleteImageRequest(oscgo.DeleteImageRequest{
			ImageId: d.Id(),
		}).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp.StatusCode, err)
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("Error deleting the image %s", err)
	}

	if err := resourceImageWaitForDestroy(d.Id(), conn); err != nil {
		return err
	}

	d.SetId("")
	return nil
}

func resourceImageWaitForDestroy(id string, conn *oscgo.APIClient) error {
	log.Printf("[INFO] Waiting for OMI %s to be deleted...", id)

	filterReq := oscgo.ReadImagesRequest{
		Filters: &oscgo.FiltersImage{ImageIds: &[]string{id}},
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"available", "pending"},
		Target:     []string{"destroyed", "failed"},
		Refresh:    ImageStateRefreshFunc(conn, filterReq, "failed"),
		Timeout:    10 * time.Minute,
		MinTimeout: 30 * time.Second,
		Delay:      1 * time.Minute,
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("error waiting for OMI (%s) to be deleted: %v", id, err)
	}

	return nil
}

// ImageStateRefreshFunc ...
func ImageStateRefreshFunc(client *oscgo.APIClient, req oscgo.ReadImagesRequest, failState string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		var resp oscgo.ReadImagesResponse
		var err error
		err = resource.Retry(120*time.Second, func() *resource.RetryError {
			var err error
			rp, httpResp, err := client.ImageApi.ReadImages(context.Background()).ReadImagesRequest(req).Execute()
			if err != nil {
				return utils.CheckThrottling(httpResp.StatusCode, err)
			}
			resp = rp
			return nil
		})

		if err != nil {
			return nil, "failed", err
		}

		state := "destroyed"

		if resp.HasImages() && len(resp.GetImages()) > 0 {
			images := resp.GetImages()
			state = images[0].GetState()

			if state == failState {
				return images[0], state, fmt.Errorf("Failed to reach target state. Reason: %v", state)
			}
		}

		log.Printf("[INFO] OMI state %s", state)

		return resp, state, nil
	}
}

// Returns a set of block device mappings.
func omiBlockDeviceMappings(m []oscgo.BlockDeviceMappingImage) []map[string]interface{} {
	blockDeviceMapping := make([]map[string]interface{}, len(m))

	for k, v := range m {
		block := make(map[string]interface{})
		block["device_name"] = v.GetDeviceName()
		block["virtual_device_name"] = v.GetVirtualDeviceName()
		if val, ok := v.GetBsuOk(); ok {
			block["bsu"] = getBsuToCreate(*val)
		}
		blockDeviceMapping[k] = block
	}
	return blockDeviceMapping
}

func getBsuToCreate(bsu oscgo.BsuToCreate) []map[string]interface{} {
	return []map[string]interface{}{{
		"delete_on_vm_deletion": bsu.GetDeleteOnVmDeletion(),
		"iops":                  bsu.GetIops(),
		"snapshot_id":           bsu.GetSnapshotId(),
		"volume_size":           bsu.GetVolumeSize(),
		"volume_type":           bsu.GetVolumeType(),
	}}
}

func expandOmiBlockDeviceMappings(blocks []interface{}) []oscgo.BlockDeviceMappingImage {
	var blockDevices []oscgo.BlockDeviceMappingImage

	for _, v := range blocks {
		blockDevice := oscgo.BlockDeviceMappingImage{}

		value := v.(map[string]interface{})
		if bsu := value["bsu"].([]interface{}); bsu != nil {
			blockDevice.SetBsu(expandOmiBlockDeviceBSU(bsu))
		}

		if deviceName := value["device_name"].(string); deviceName != "" {
			blockDevice.SetDeviceName(deviceName)
		}
		if virtualDeviceName := value["virtual_device_name"].(string); virtualDeviceName != "" {
			blockDevice.SetVirtualDeviceName(virtualDeviceName)
		}

		blockDevices = append(blockDevices, blockDevice)
	}
	return blockDevices
}

func expandOmiBlockDeviceBSU(bsu []interface{}) oscgo.BsuToCreate {
	bsuToCreate := oscgo.BsuToCreate{}

	for _, v := range bsu {
		val := v.(map[string]interface{})
		if del := val["delete_on_vm_deletion"].(bool); del {
			bsuToCreate.SetDeleteOnVmDeletion(del)
		}
		if snap := val["snapshot_id"].(string); snap != "" {
			bsuToCreate.SetSnapshotId(snap)
		}
		if vSize := val["volume_size"].(int); vSize > 0 {
			bsuToCreate.SetVolumeSize(int32(vSize))
		}
		if vType := val["volume_type"].(string); vType != "" {
			bsuToCreate.SetVolumeType(vType)
			if iops := val["iops"].(int); iops > 0 && vType == "io1" {
				bsuToCreate.SetIops(int32(iops))
			}
		}
	}
	return bsuToCreate
}

// Returns the state reason.
func omiStateReason(m *oscgo.StateComment) []map[string]interface{} {
	return []map[string]interface{}{{
		"state_code":    m.GetStateCode(),
		"state_message": m.GetStateMessage(),
	}}
}
