package outscale

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/antihax/optional"

	oscgo "github.com/marinsalinas/osc-sdk-go"

	"github.com/openlyinc/pointy"
	"github.com/spf13/cast"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const (
	// OutscaleImageRetryTimeout ...
	OutscaleImageRetryTimeout = 40 * time.Minute
	// OutscaleImageDeleteRetryTimeout ...
	OutscaleImageDeleteRetryTimeout = 90 * time.Minute
	// OutscaleImageRetryDelay ...
	OutscaleImageRetryDelay = 5 * time.Second
	// OutscaleImageRetryMinTimeout ...
	OutscaleImageRetryMinTimeout = 3 * time.Second
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
				Optional: true,
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
				ForceNew: true,
			},
			"architecture": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"source_region_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"source_image_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"file_location": {
				Type:     schema.TypeString,
				Optional: true,
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
	conn := meta.(*OutscaleClient).OSCAPI

	imageRequest := oscgo.CreateImageRequest{
		ImageName: pointy.String(cast.ToString(d.Get("image_name"))),
	}

	if v := cast.ToString(d.Get("vm_id")); v != "" {
		imageRequest.SetVmId(v)
	}

	if v := cast.ToString(d.Get("description")); v != "" {
		imageRequest.SetDescription(v)
	}

	if v, ok := d.GetOk("no_reboot"); ok {
		imageRequest.SetNoReboot(cast.ToBool(v))
	}

	if v := cast.ToString(d.Get("architecture")); v != "" {
		imageRequest.SetArchitecture(v)
	}

	if v := cast.ToString(d.Get("file_location")); v != "" {
		imageRequest.SetFileLocation(v)
	}

	if v := cast.ToString(d.Get("source_image_id")); v != "" {
		imageRequest.SetSourceImageId(v)
	}

	if v := cast.ToString(d.Get("source_region_name")); v != "" {
		imageRequest.SetSourceRegionName(v)
	}

	if v := cast.ToString(d.Get("root_device_name")); v != "" {
		imageRequest.SetRootDeviceName(v)
	}

	resp, _, err := conn.ImageApi.CreateImage(context.Background(), &oscgo.CreateImageOpts{
		CreateImageRequest: optional.NewInterface(imageRequest),
	})
	if err != nil {
		return err
	}

	if !resp.HasImage() {
		return nil
	}

	image := resp.GetImage()

	log.Printf("[DEBUG] Waiting for OMI %s to become available...", *image.ImageId)

	filterReq := &oscgo.ReadImagesOpts{
		ReadImagesRequest: optional.NewInterface(oscgo.ReadImagesRequest{
			Filters: &oscgo.FiltersImage{ImageIds: &[]string{*image.ImageId}},
		}),
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"pending"},
		Target:     []string{"available"},
		Refresh:    ImageOAPIStateRefreshFunc(conn, filterReq, "failed"),
		Timeout:    10 * time.Minute,
		MinTimeout: 30 * time.Second,
		Delay:      1 * time.Minute,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Errorrr waiting for OMI (%s) to be ready: %v", *image.ImageId, err)
	}

	d.SetId(*image.ImageId)

	return resourceOAPIImageRead(d, meta)
}

func resourceOAPIImageRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI
	id := d.Id()

	req := &oscgo.ReadImagesOpts{
		ReadImagesRequest: optional.NewInterface(oscgo.ReadImagesRequest{
			Filters: &oscgo.FiltersImage{ImageIds: &[]string{id}},
		}),
	}

	resp, _, err := conn.ImageApi.ReadImages(context.Background(), req)
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

		set("architecture", image.Architecture)
		set("creation_date", image.CreationDate)
		set("description", image.Description)
		set("image_id", image.ImageId)
		set("file_location", image.FileLocation)
		set("account_alias", image.AccountAlias)
		set("account_id", image.AccountId)
		set("image_type", image.ImageType)
		set("image_name", image.ImageName)
		set("root_device_name", image.RootDeviceName)
		set("root_device_type", image.RootDeviceType)
		set("state", image.State)

		if err := set("block_device_mappings", omiOAPIBlockDeviceMappings(*image.BlockDeviceMappings)); err != nil {
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
		if err := set("tags", getOapiTagSet(image.Tags)); err != nil {
			return err
		}

		return d.Set("request_id", resp.ResponseContext.RequestId)
	})
}

func setResourcePermissions(por oscgo.PermissionsOnResource) []map[string]interface{} {
	return []map[string]interface{}{
		map[string]interface{}{
			"global_permission": por.GetGlobalPermission(),
			"account_ids":       por.GetAccountIds(),
		},
	}
}

func resourceOAPIImageUpdate(d *schema.ResourceData, meta interface{}) error {
	//	conn := meta.(*OutscaleClient).OSCAPI

	d.Partial(true)
	//TODO: add tags
	// if d.Get("description").(string) != "" {
	// 	_, _, err := conn.ImageApi.UpdateImage(context.Background(), &oscgo.UpdateImageOpts{
	// 		UpdateImageRequest: optional.NewInterface(oscgo.UpdateImageRequest{}),
	// 	})
	// 	if err != nil {
	// 		return err
	// 	}
	// 	d.SetPartial("description")
	// }

	d.Partial(false)

	return resourceOAPIImageRead(d, meta)
}

func resourceOAPIImageDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	_, _, err := conn.ImageApi.DeleteImage(context.Background(), &oscgo.DeleteImageOpts{
		DeleteImageRequest: optional.NewInterface(oscgo.DeleteImageRequest{
			ImageId: d.Id(),
		}),
	})

	if err != nil {
		return fmt.Errorf("Error deleting the image")
	}

	if err := resourceOutscaleOAPIImageWaitForDestroy(d.Id(), conn); err != nil {
		return err
	}

	d.SetId("")
	return nil
}

func resourceOutscaleOAPIImageWaitForDestroy(id string, conn *oscgo.APIClient) error {
	log.Printf("[INFO] Waiting for OMI %s to be deleted...", id)

	filterReq := &oscgo.ReadImagesOpts{
		ReadImagesRequest: optional.NewInterface(oscgo.ReadImagesRequest{
			Filters: &oscgo.FiltersImage{ImageIds: &[]string{id}},
		}),
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"available", "pending"},
		Target:     []string{"destroyed", "failed"},
		Refresh:    ImageOAPIStateRefreshFunc(conn, filterReq, "failed"),
		Timeout:    10 * time.Minute,
		MinTimeout: 30 * time.Second,
		Delay:      1 * time.Minute,
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("error waiting for OMI (%s) to be deleted: %v", id, err)
	}

	return nil
}

// ImageOAPIStateRefreshFunc ...
func ImageOAPIStateRefreshFunc(client *oscgo.APIClient, req *oscgo.ReadImagesOpts, failState string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, _, err := client.ImageApi.ReadImages(context.Background(), req)
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
func omiOAPIBlockDeviceMappings(m []oscgo.BlockDeviceMappingImage) []map[string]interface{} {
	blockDeviceMapping := make([]map[string]interface{}, len(m))

	for k, v := range m {
		blockDeviceMapping[k] = map[string]interface{}{
			"device_name":         v.GetDeviceName(),
			"virtual_device_name": v.GetVirtualDeviceName(),
			"bsu": map[string]interface{}{
				"delete_on_vm_deletion": cast.ToString(v.Bsu.GetDeleteOnVmDeletion()),
				"iops":                  cast.ToString(v.Bsu.GetIops()),
				"snapshot_id":           v.Bsu.GetSnapshotId(),
				"volume_size":           cast.ToString(v.Bsu.GetVolumeSize()),
				"volume_type":           v.Bsu.GetVolumeType(),
			},
		}
	}
	return blockDeviceMapping
}

// Returns the state reason.
func omiOAPIStateReason(m *oscgo.StateComment) map[string]interface{} {
	s := make(map[string]interface{})
	if m != nil {
		s["state_code"] = m.GetStateCode()
		s["state_message"] = m.GetStateMessage()
	} else {
		s["state_code"] = "UNSET"
		s["state_message"] = "UNSET"
	}
	return s
}
