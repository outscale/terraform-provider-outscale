package oapi

import (
	"context"
	"fmt"
	"log"
	"slices"
	"time"

	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/utils"
	"github.com/samber/lo"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	// OutscaleImageRetryTimeout ...
	OutscaleImageRetryTimeout = 40 * time.Minute
	// OutscaleImageDeleteRetryTimeout ...
	OutscaleImageDeleteRetryTimeout = 90 * time.Minute
	// OutscaleImageRetryDelay ...
	OutscaleImageRetryDelay = 20 * time.Second
	// OutscaleImageRetryMinTimeout ...
	OutscaleImageRetryMinTimeout = 3 * time.Second
)

func ResourceOutscaleImage() *schema.Resource {
	return &schema.Resource{
		Create: resourceOAPIImageCreate,
		Read:   resourceOAPIImageRead,
		Update: resourceOAPIImageUpdate,
		Delete: resourceOAPIImageDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(40 * time.Minute),
			Read:   schema.DefaultTimeout(ReadDefaultTimeout),
			Update: schema.DefaultTimeout(UpdateDefaultTimeout),
			Delete: schema.DefaultTimeout(10 * time.Minute),
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
			"boot_modes": {
				Type:     schema.TypeSet,
				Computed: true,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
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
			"tpm_mandatory": {
				Type:     schema.TypeBool,
				Computed: true,
				Optional: true,
				ForceNew: true,
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
			"tags": TagsSchemaSDK(),
			"vm_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceOAPIImageCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*client.OutscaleClient).OSCAPI
	timeout := d.Timeout(schema.TimeoutCreate)

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
		blockDevices := expandOmiBlockDeviceOApiMappings(blocks.([]interface{}))
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
	tpm := d.GetRawConfig().GetAttr("tpm_mandatory")
	if !tpm.IsNull() {
		imageRequest.SetTpmMandatory(tpm.True())
	}

	if v, ok := d.GetOk("root_device_name"); ok {
		imageRequest.SetRootDeviceName(v.(string))
	}

	if v, ok := d.GetOk("boot_modes"); ok {
		modes := utils.SetToStringSlice(v.(*schema.Set))
		if lo.EveryBy(modes, func(s string) bool { return slices.Contains([]string{"uefi", "legacy"}, s) }) {
			imageRequest.SetBootModes(lo.Map(modes, func(s string, _ int) oscgo.BootMode { return (oscgo.BootMode)(s) }))
		} else {
			return fmt.Errorf("the boot modes compatible with the omi are: uefi, legacy - provided: %v", modes)
		}
	}
	var resp oscgo.CreateImageResponse
	var err error
	err = retry.RetryContext(context.Background(), timeout, func() *retry.RetryError {
		rp, httpResp, err := conn.ImageApi.CreateImage(context.Background()).CreateImageRequest(imageRequest).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
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

	stateConf := &retry.StateChangeConf{
		Pending:    []string{"pending"},
		Target:     []string{"available"},
		Refresh:    ImageOAPIStateRefreshFunc(conn, req, "failed", timeout),
		Timeout:    timeout,
		MinTimeout: 30 * time.Second,
		Delay:      5 * time.Second,
	}

	if _, err = stateConf.WaitForStateContext(context.Background()); err != nil {
		return fmt.Errorf("error waiting for omi (%s) to be ready: %w", *image.ImageId, err)
	}
	d.SetId(image.GetImageId())

	err = createOAPITagsSDK(conn, d)
	if err != nil {
		return err
	}

	return resourceOAPIImageRead(d, meta)
}

func resourceOAPIImageRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*client.OutscaleClient).OSCAPI
	timeout := d.Timeout(schema.TimeoutRead)
	id := d.Id()

	req := oscgo.ReadImagesRequest{
		Filters: &oscgo.FiltersImage{ImageIds: &[]string{id}},
	}

	var resp oscgo.ReadImagesResponse
	err := retry.RetryContext(context.Background(), timeout, func() *retry.RetryError {
		var err error
		rp, httpResp, err := conn.ImageApi.ReadImages(context.Background()).ReadImagesRequest(req).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		resp = rp
		return nil
	})
	if err != nil {
		return fmt.Errorf("error reading for omi (%s): %w", id, err)
	}
	if utils.IsResponseEmpty(len(resp.GetImages()), "Image", d.Id()) {
		d.SetId("")
		return nil
	}
	image := resp.GetImages()[0]

	return resourceDataAttrSetter(d, func(set AttributeSetter) error {
		d.SetId(*image.ImageId)

		if err := set("architecture", image.Architecture); err != nil {
			return err
		}
		if err := set("boot_modes", lo.Map(image.GetBootModes(), func(b oscgo.BootMode, _ int) string { return string(b) })); err != nil {
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
		if err := set("block_device_mappings", omiOAPIBlockDeviceMappings(image.GetBlockDeviceMappings())); err != nil {
			return err
		}
		if err := set("product_codes", image.ProductCodes); err != nil {
			return err
		}
		if err := set("state_comment", omiOAPIStateReason(image.StateComment)); err != nil {
			return err
		}
		if err := set("permissions_to_launch", setResourcePermissions(*image.PermissionsToLaunch)); err != nil {
			return err
		}
		if err := set("tpm_mandatory", image.TpmMandatory); err != nil {
			return err
		}
		if err := d.Set("tags", FlattenOAPITagsSDK(image.GetTags())); err != nil {
			return fmt.Errorf("unable to set image tags: %w", err)
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

func resourceOAPIImageUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*client.OutscaleClient).OSCAPI

	if err := updateOAPITagsSDK(conn, d); err != nil {
		return err
	}
	return resourceOAPIImageRead(d, meta)
}

func resourceOAPIImageDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*client.OutscaleClient).OSCAPI
	timeout := d.Timeout(schema.TimeoutDelete)

	err := retry.RetryContext(context.Background(), timeout, func() *retry.RetryError {
		_, httpResp, err := conn.ImageApi.DeleteImage(context.Background()).DeleteImageRequest(oscgo.DeleteImageRequest{
			ImageId: d.Id(),
		}).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("error deleting the image %w", err)
	}

	if err := ResourceOutscaleImageWaitForDestroy(d.Id(), conn, timeout); err != nil {
		return err
	}

	d.SetId("")
	return nil
}

func ResourceOutscaleImageWaitForDestroy(id string, conn *oscgo.APIClient, timeOut time.Duration) error {
	log.Printf("[INFO] Waiting for OMI %s to be deleted...", id)

	filterReq := oscgo.ReadImagesRequest{
		Filters: &oscgo.FiltersImage{ImageIds: &[]string{id}},
	}

	stateConf := &retry.StateChangeConf{
		Pending:    []string{"available", "pending"},
		Target:     []string{"destroyed", "failed"},
		Refresh:    ImageOAPIStateRefreshFunc(conn, filterReq, "failed", timeOut),
		Timeout:    timeOut,
		MinTimeout: 30 * time.Second,
		Delay:      5 * time.Second,
	}

	if _, err := stateConf.WaitForStateContext(context.Background()); err != nil {
		return fmt.Errorf("error waiting for omi (%s) to be deleted: %w", id, err)
	}

	return nil
}

// ImageOAPIStateRefreshFunc ...
func ImageOAPIStateRefreshFunc(client *oscgo.APIClient, req oscgo.ReadImagesRequest, failState string, timeOut time.Duration) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		var resp oscgo.ReadImagesResponse
		err := retry.RetryContext(context.Background(), timeOut, func() *retry.RetryError {
			var err error
			rp, httpResp, err := client.ImageApi.ReadImages(context.Background()).ReadImagesRequest(req).Execute()
			if err != nil {
				return utils.CheckThrottling(httpResp, err)
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
				return images[0], state, fmt.Errorf("failed to reach target state:: %v", state)
			}
		}

		return resp, state, nil
	}
}

// Returns a set of block device mappings.
func omiOAPIBlockDeviceMappings(m []oscgo.BlockDeviceMappingImage) []map[string]interface{} {
	blockDeviceMapping := make([]map[string]interface{}, len(m))

	for k, v := range m {
		block := make(map[string]interface{})
		block["device_name"] = v.GetDeviceName()
		block["virtual_device_name"] = v.GetVirtualDeviceName()
		if val, ok := v.GetBsuOk(); ok {
			block["bsu"] = getOAPIBsuToCreate(*val)
		}
		blockDeviceMapping[k] = block
	}
	return blockDeviceMapping
}

func getOAPIBsuToCreate(bsu oscgo.BsuToCreate) []map[string]interface{} {
	return []map[string]interface{}{{
		"delete_on_vm_deletion": bsu.GetDeleteOnVmDeletion(),
		"iops":                  bsu.GetIops(),
		"snapshot_id":           bsu.GetSnapshotId(),
		"volume_size":           bsu.GetVolumeSize(),
		"volume_type":           bsu.GetVolumeType(),
	}}
}

func expandOmiBlockDeviceOApiMappings(blocks []interface{}) []oscgo.BlockDeviceMappingImage {
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
func omiOAPIStateReason(m *oscgo.StateComment) []map[string]interface{} {
	return []map[string]interface{}{{
		"state_code":    m.GetStateCode(),
		"state_message": m.GetStateMessage(),
	}}
}
