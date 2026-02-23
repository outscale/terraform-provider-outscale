package oapi

import (
	"context"
	"fmt"
	"log"
	"slices"
	"time"

	"github.com/outscale/goutils/sdk/ptr"
	"github.com/outscale/osc-sdk-go/v3/pkg/options"
	"github.com/outscale/osc-sdk-go/v3/pkg/osc"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/framework/fwhelpers/from"
	"github.com/outscale/terraform-provider-outscale/internal/utils"
	"github.com/samber/lo"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
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
		CreateContext: resourceOAPIImageCreate,
		ReadContext:   resourceOAPIImageRead,
		UpdateContext: resourceOAPIImageUpdate,
		DeleteContext: resourceOAPIImageDelete,
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

func resourceOAPIImageCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.OutscaleClient).OSC

	timeout := d.Timeout(schema.TimeoutCreate)

	imageRequest := osc.CreateImageRequest{}
	if v, ok := d.GetOk("image_name"); ok {
		imageRequest.ImageName = new(v.(string))
	}
	if v, ok := d.GetOk("vm_id"); ok {
		imageRequest.VmId = new(v.(string))
	}
	if v, ok := d.GetOk("description"); ok {
		imageRequest.Description = new(v.(string))
	}
	if blocks, ok := d.GetOk("block_device_mappings"); ok {
		blockDevices := expandOmiBlockDeviceOApiMappings(blocks.([]interface{}))
		imageRequest.BlockDeviceMappings = new(blockDevices)
	}
	if v, ok := d.GetOk("no_reboot"); ok {
		imageRequest.NoReboot = new(v.(bool))
	}

	if v, ok := d.GetOk("architecture"); ok {
		imageRequest.Architecture = new(v.(string))
	}

	if v, ok := d.GetOk("file_location"); ok {
		imageRequest.FileLocation = new(v.(string))
	}

	if v, ok := d.GetOk("source_image_id"); ok {
		imageRequest.SourceImageId = new(v.(string))
	}

	if v, ok := d.GetOk("source_region_name"); ok {
		imageRequest.SourceRegionName = new(v.(string))
	}
	tpm := d.GetRawConfig().GetAttr("tpm_mandatory")
	if !tpm.IsNull() {
		imageRequest.TpmMandatory = new(tpm.True())
	}

	if v, ok := d.GetOk("root_device_name"); ok {
		imageRequest.RootDeviceName = new(v.(string))
	}

	if v, ok := d.GetOk("boot_modes"); ok {
		modes := utils.SetToStringSlice(v.(*schema.Set))
		if lo.EveryBy(modes, func(s string) bool { return slices.Contains([]string{"uefi", "legacy"}, s) }) {
			imageRequest.BootModes = new(lo.Map(modes, func(s string, _ int) osc.BootMode { return (osc.BootMode)(s) }))
		} else {
			return diag.Errorf("the boot modes compatible with the omi are: uefi, legacy - provided: %v", modes)
		}
	}
	resp, err := client.CreateImage(ctx, imageRequest, options.WithRetryTimeout(timeout))
	if err != nil {
		return diag.FromErr(err)
	}

	if resp.Image == nil {
		return nil
	}

	image := ptr.From(resp.Image)

	log.Printf("[DEBUG] Waiting for OMI %s to become available...", image.ImageId)

	req := osc.ReadImagesRequest{Filters: &osc.FiltersImage{ImageIds: &[]string{image.ImageId}}}

	stateConf := &retry.StateChangeConf{
		Pending: []string{string(osc.ImageStatePending)},
		Target:  []string{string(osc.ImageStateAvailable)},
		Timeout: timeout,
		Refresh: ImageOAPIStateRefreshFunc(ctx, client, req, string(osc.ImageStateFailed), timeout),
	}

	if _, err = stateConf.WaitForStateContext(ctx); err != nil {
		return diag.Errorf("error waiting for omi (%s) to be ready: %v", image.ImageId, err)
	}
	d.SetId(image.ImageId)

	err = createOAPITagsSDK(ctx, client, timeout, d)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceOAPIImageRead(ctx, d, meta)
}

func resourceOAPIImageRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.OutscaleClient).OSC
	timeout := d.Timeout(schema.TimeoutRead)

	id := d.Id()

	req := osc.ReadImagesRequest{
		Filters: &osc.FiltersImage{ImageIds: &[]string{id}},
	}

	resp, err := client.ReadImages(ctx, req, options.WithRetryTimeout(timeout))
	if err != nil {
		return diag.Errorf("error reading for omi (%s): %v", id, err)
	}
	if utils.IsResponseEmpty(len(ptr.From(resp.Images)), "Image", d.Id()) {
		d.SetId("")
		return nil
	}
	image := ptr.From(resp.Images)[0]

	return diag.FromErr(resourceDataAttrSetter(d, func(set AttributeSetter) error {
		d.SetId(image.ImageId)

		if err := set("architecture", image.Architecture); err != nil {
			return err
		}
		if err := set("boot_modes", lo.Map(image.BootModes, func(b osc.BootMode, _ int) string { return string(b) })); err != nil {
			return err
		}
		if err := set("creation_date", from.ISO8601(image.CreationDate)); err != nil {
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
		if err := set("block_device_mappings", omiOAPIBlockDeviceMappings(ptr.From(image.BlockDeviceMappings))); err != nil {
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
		if err := d.Set("tags", FlattenOAPITagsSDK(image.Tags)); err != nil {
			return fmt.Errorf("unable to set image tags: %w", err)
		}

		return nil
	}))
}

func setResourcePermissions(por osc.PermissionsOnResource) []map[string]interface{} {
	return []map[string]interface{}{
		{
			"global_permission": por.GlobalPermission,
			"account_ids":       por.AccountIds,
		},
	}
}

func resourceOAPIImageUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.OutscaleClient).OSC
	timeout := d.Timeout(schema.TimeoutUpdate)

	if err := updateOAPITagsSDK(ctx, client, timeout, d); err != nil {
		return diag.FromErr(err)
	}
	return resourceOAPIImageRead(ctx, d, meta)
}

func resourceOAPIImageDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.OutscaleClient).OSC
	timeout := d.Timeout(schema.TimeoutDelete)

	_, err := client.DeleteImage(ctx, osc.DeleteImageRequest{
		ImageId: d.Id(),
	}, options.WithRetryTimeout(timeout))
	if err != nil {
		return diag.Errorf("error deleting the image %v", err)
	}

	if err := ResourceOutscaleImageWaitForDestroy(ctx, d.Id(), client, timeout); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}

func ResourceOutscaleImageWaitForDestroy(ctx context.Context, id string, client *osc.Client, timeOut time.Duration) error {
	log.Printf("[INFO] Waiting for OMI %s to be deleted...", id)

	filterReq := osc.ReadImagesRequest{
		Filters: &osc.FiltersImage{ImageIds: &[]string{id}},
	}

	stateConf := &retry.StateChangeConf{
		Pending: []string{string(osc.ImageStateAvailable), string(osc.ImageStatePending)},
		Target:  []string{},
		Timeout: timeOut,
		Refresh: func() (any, string, error) {
			resp, err := client.ReadImages(ctx, filterReq, options.WithRetryTimeout(timeOut))
			if err != nil {
				return nil, "", err
			}
			if resp.Images == nil || len(*resp.Images) == 0 {
				return nil, "", nil
			}

			images := ptr.From(resp.Images)
			state := string(images[0].State)
			if state == string(osc.ImageStateFailed) {
				return images[0], state, fmt.Errorf("failed to reach target state: %v", state)
			}

			return resp, state, nil
		},
	}

	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("error waiting for omi (%s) to be deleted: %w", id, err)
	}

	return nil
}

// ImageOAPIStateRefreshFunc ...
func ImageOAPIStateRefreshFunc(ctx context.Context, client *osc.Client, req osc.ReadImagesRequest, failState string, timeOut time.Duration) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := client.ReadImages(ctx, req, options.WithRetryTimeout(timeOut))
		if err != nil {
			return nil, "", err
		}
		if resp.Images == nil || len(*resp.Images) == 0 {
			return nil, "", fmt.Errorf("failed to get image")
		}

		images := ptr.From(resp.Images)
		state := string(images[0].State)

		if state == failState {
			return images[0], state, fmt.Errorf("failed to reach target state: %v", state)
		}

		return resp, state, nil
	}
}

// Returns a set of block device mappings.
func omiOAPIBlockDeviceMappings(m []osc.BlockDeviceMappingImage) []map[string]interface{} {
	blockDeviceMapping := make([]map[string]interface{}, len(m))

	for k, v := range m {
		block := make(map[string]interface{})
		block["device_name"] = v.DeviceName
		block["virtual_device_name"] = v.VirtualDeviceName
		if v.Bsu != nil {
			block["bsu"] = getOAPIBsuToCreate(*v.Bsu)
		}
		blockDeviceMapping[k] = block
	}
	return blockDeviceMapping
}

func getOAPIBsuToCreate(bsu osc.BsuToCreate) []map[string]interface{} {
	return []map[string]interface{}{{
		"delete_on_vm_deletion": bsu.DeleteOnVmDeletion,
		"iops":                  bsu.Iops,
		"snapshot_id":           bsu.SnapshotId,
		"volume_size":           bsu.VolumeSize,
		"volume_type":           ptr.From(bsu.VolumeType),
	}}
}

func expandOmiBlockDeviceOApiMappings(blocks []interface{}) []osc.BlockDeviceMappingImage {
	var blockDevices []osc.BlockDeviceMappingImage

	for _, v := range blocks {
		blockDevice := osc.BlockDeviceMappingImage{}

		value := v.(map[string]interface{})
		if bsu := value["bsu"].([]interface{}); bsu != nil {
			blockDevice.Bsu = new(expandOmiBlockDeviceBSU(bsu))
		}

		if deviceName := value["device_name"].(string); deviceName != "" {
			blockDevice.DeviceName = &deviceName
		}
		if virtualDeviceName := value["virtual_device_name"].(string); virtualDeviceName != "" {
			blockDevice.VirtualDeviceName = &virtualDeviceName
		}

		blockDevices = append(blockDevices, blockDevice)
	}
	return blockDevices
}

func expandOmiBlockDeviceBSU(bsu []interface{}) osc.BsuToCreate {
	bsuToCreate := osc.BsuToCreate{}

	for _, v := range bsu {
		val := v.(map[string]interface{})
		if del := val["delete_on_vm_deletion"].(bool); del {
			bsuToCreate.DeleteOnVmDeletion = &del
		}
		if snap := val["snapshot_id"].(string); snap != "" {
			bsuToCreate.SnapshotId = &snap
		}
		if vSize := val["volume_size"].(int); vSize > 0 {
			bsuToCreate.VolumeSize = new(vSize)
		}
		if vType := val["volume_type"].(string); vType != "" {
			bsuToCreate.VolumeType = new(osc.VolumeType(vType))
			if iops := val["iops"].(int); iops > 0 && vType == "io1" {
				bsuToCreate.Iops = new(iops)
			}
		}
	}
	return bsuToCreate
}

// Returns the state reason.
func omiOAPIStateReason(m *osc.StateComment) []map[string]interface{} {
	return []map[string]interface{}{{
		"state_code":    m.StateCode,
		"state_message": m.StateMessage,
	}}
}
