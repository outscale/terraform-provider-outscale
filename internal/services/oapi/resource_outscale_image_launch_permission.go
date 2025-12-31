package oapi

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/services/oapi/oapihelpers"
	"github.com/outscale/terraform-provider-outscale/internal/utils"

	"github.com/spf13/cast"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceOutscaleImageLaunchPermission() *schema.Resource {
	return &schema.Resource{
		Exists: ResourceOutscaleImageLaunchPermissionExists,
		Create: ResourceOutscaleImageLaunchPermissionCreate,
		Read:   ResourceOutscaleImageLaunchPermissionRead,
		Delete: ResourceOutscaleImageLaunchPermissionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"image_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"permission_additions": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"global_permission": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
							Default:  "false",
						},
						"account_ids": {
							Type:     schema.TypeList,
							Optional: true,
							ForceNew: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"permission_removals": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"global_permission": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
							Default:  "false",
						},
						"account_ids": {
							Type:     schema.TypeList,
							Optional: true,
							ForceNew: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"permissions_to_launch": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"global_permission": {
							Type:     schema.TypeString,
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
		},
	}
}

func ResourceOutscaleImageLaunchPermissionExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	conn := meta.(*client.OutscaleClient).OSCAPI

	imageID := d.Get("image_id").(string)
	return oapihelpers.ImageHasLaunchPermission(conn, imageID)
}

func expandOAPIImagePermission(permissionType interface{}) (res oscgo.PermissionsOnResource) {
	if len(permissionType.([]interface{})) > 0 {
		permission := permissionType.([]interface{})[0].(map[string]interface{})

		if globalPermission, ok := permission["global_permission"]; ok {
			res.SetGlobalPermission(cast.ToBool(globalPermission))
		}
		if accountIDs, ok := permission["account_ids"]; ok {
			for _, accountID := range accountIDs.([]interface{}) {
				res.SetAccountIds(append(res.GetAccountIds(), accountID.(string)))
			}
		}
	}
	return
}

func ResourceOutscaleImageLaunchPermissionCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*client.OutscaleClient).OSCAPI

	imageID, ok := d.GetOk("image_id")

	if !ok {
		return fmt.Errorf("please provide the required attribute image_id")
	}
	log.Printf("Creating Outscale Image Launch Permission, image_id (%+v)", imageID.(string))

	permissionLaunch := oscgo.PermissionsOnResourceCreation{}
	if permissionAdditions, ok := d.GetOk("permission_additions"); ok {
		permissionLaunch.SetAdditions(expandOAPIImagePermission(permissionAdditions))
	}
	if permissionRemovals, ok := d.GetOk("permission_removals"); ok {
		permissionLaunch.SetRemovals(expandOAPIImagePermission(permissionRemovals))
	}

	request := oscgo.UpdateImageRequest{
		ImageId:             imageID.(string),
		PermissionsToLaunch: &permissionLaunch,
	}

	err := retry.Retry(5*time.Minute, func() *retry.RetryError {
		var err error
		_, httpResp, err := conn.ImageApi.UpdateImage(context.Background()).UpdateImageRequest(request).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		return nil
	})

	var errString string
	if err != nil {
		errString = err.Error()

		return fmt.Errorf("error creating omi launch permission: %s", errString)
	}

	d.SetId(imageID.(string))

	return ResourceOutscaleImageLaunchPermissionRead(d, meta)
}

func ResourceOutscaleImageLaunchPermissionRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*client.OutscaleClient).OSCAPI

	var resp oscgo.ReadImagesResponse
	err := retry.Retry(5*time.Minute, func() *retry.RetryError {
		var err error
		rp, httpResp, err := conn.ImageApi.ReadImages(context.Background()).ReadImagesRequest(oscgo.ReadImagesRequest{
			Filters: &oscgo.FiltersImage{
				ImageIds: &[]string{d.Id()},
			},
		}).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		resp = rp
		return nil
	})

	var errString string
	if err != nil {
		// When an AMI disappears out from under a launch permission resource, we will
		// see either InvalidAMIID.NotFound or InvalidAMIID.Unavailable.
		if strings.Contains(fmt.Sprint(err), "InvalidAMIID") {
			log.Printf("[DEBUG] %s no longer exists, so we'll drop launch permission for the state", d.Id())
			return nil
		}
		errString = err.Error()

		return fmt.Errorf("Error reading Outscale image permission: %s", errString)
	}
	if utils.IsResponseEmpty(len(resp.GetImages()), "ImageLaunchPermission", d.Id()) {
		d.SetId("")
		return nil
	}
	result := resp.GetImages()[0]

	if err := d.Set("description", result.Description); err != nil {
		return err
	}

	lp := make(map[string]interface{})
	lp["global_permission"] = strconv.FormatBool(result.PermissionsToLaunch.GetGlobalPermission())
	lp["account_ids"] = result.PermissionsToLaunch.GetAccountIds()

	return d.Set("permissions_to_launch", []map[string]interface{}{lp})
}

func ResourceOutscaleImageLaunchPermissionDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*client.OutscaleClient).OSCAPI

	imageID, ok := d.GetOk("image_id")
	if !ok {
		return fmt.Errorf("please provide the required attribute image_id")
	}

	if permissionAdditions, ok := d.GetOk("permission_additions"); ok {
		permission := oscgo.PermissionsOnResourceCreation{}
		request := oscgo.UpdateImageRequest{
			ImageId: imageID.(string),
		}
		permission.SetRemovals(expandOAPIImagePermission(permissionAdditions))
		request.SetPermissionsToLaunch(permission)

		err := retry.Retry(5*time.Minute, func() *retry.RetryError {
			_, httpResp, err := conn.ImageApi.UpdateImage(context.Background()).UpdateImageRequest(request).Execute()
			if err != nil {
				return utils.CheckThrottling(httpResp, err)
			}
			return nil
		})

		var errString string
		if err != nil {
			errString = err.Error()

			return fmt.Errorf("error removing omi launch permission: %s", errString)
		}
	}

	d.SetId("")
	return nil
}
