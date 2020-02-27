package outscale

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/antihax/optional"
	oscgo "github.com/marinsalinas/osc-sdk-go"

	"github.com/spf13/cast"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceOutscaleOAPIImageLaunchPermission() *schema.Resource {
	return &schema.Resource{
		Exists: resourceOutscaleOAPIImageLaunchPermissionExists,
		Create: resourceOutscaleOAPIImageLaunchPermissionCreate,
		Read:   resourceOutscaleOAPIImageLaunchPermissionRead,
		Delete: resourceOutscaleOAPIImageLaunchPermissionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"image_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"permission_additions": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"global_permission": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
							Default:  "false",
						},
						"account_ids": &schema.Schema{
							Type:     schema.TypeList,
							Optional: true,
							ForceNew: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"permission_removals": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"global_permission": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
							Default:  "false",
						},
						"account_ids": &schema.Schema{
							Type:     schema.TypeList,
							Optional: true,
							ForceNew: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"permissions_to_launch": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"global_permission": &schema.Schema{
							Type:     schema.TypeString,
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
			"request_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceOutscaleOAPIImageLaunchPermissionExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	conn := meta.(*OutscaleClient).OSCAPI

	imageID := d.Get("image_id").(string)
	return hasOAPILaunchPermission(conn, imageID)
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

func resourceOutscaleOAPIImageLaunchPermissionCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	imageID, ok := d.GetOk("image_id")

	if !ok {
		return fmt.Errorf("please provide the required attribute image_id")
	}
	log.Printf("Creating Outscale Image Launch Permission, image_id (%+v)", imageID.(string))

	permissionLunch := oscgo.PermissionsOnResourceCreation{}
	if permissionAdditions, ok := d.GetOk("permission_additions"); ok {
		permissionLunch.SetAdditions(expandOAPIImagePermission(permissionAdditions))
	}
	if permissionRemovals, ok := d.GetOk("permission_removals"); ok {
		permissionLunch.SetRemovals(expandOAPIImagePermission(permissionRemovals))
	}

	request := oscgo.UpdateImageRequest{
		ImageId:             imageID.(string),
		PermissionsToLaunch: permissionLunch,
	}

	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		var err error
		_, _, err = conn.ImageApi.UpdateImage(context.Background(), &oscgo.UpdateImageOpts{UpdateImageRequest: optional.NewInterface(request)})
		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded:") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	var errString string

	if err != nil {
		errString = err.Error()

		return fmt.Errorf("error creating omi launch permission: %s", errString)
	}

	d.SetId(imageID.(string))

	return resourceOutscaleOAPIImageLaunchPermissionRead(d, meta)
}

func resourceOutscaleOAPIImageLaunchPermissionRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	var resp oscgo.ReadImagesResponse
	var err error
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, _, err = conn.ImageApi.ReadImages(context.Background(), &oscgo.ReadImagesOpts{ReadImagesRequest: optional.NewInterface(oscgo.ReadImagesRequest{
			Filters: &oscgo.FiltersImage{
				ImageIds: &[]string{d.Id()},
			},
		})})
		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded:") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
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

	result := resp.GetImages()[0]

	if err := d.Set("request_id", resp.ResponseContext.GetRequestId()); err != nil {
		return err
	}
	if err := d.Set("description", result.Description); err != nil {
		return err
	}

	lp := make(map[string]interface{})
	lp["global_permission"] = strconv.FormatBool(result.PermissionsToLaunch.GetGlobalPermission())
	lp["account_ids"] = result.PermissionsToLaunch.GetAccountIds()

	return d.Set("permissions_to_launch", []map[string]interface{}{lp})
}

func resourceOutscaleOAPIImageLaunchPermissionDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

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

		err := resource.Retry(5*time.Minute, func() *resource.RetryError {
			var err error
			_, _, err = conn.ImageApi.UpdateImage(context.Background(), &oscgo.UpdateImageOpts{UpdateImageRequest: optional.NewInterface(request)})
			if err != nil {
				if strings.Contains(err.Error(), "RequestLimitExceeded:") {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
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

func hasOAPILaunchPermission(conn *oscgo.APIClient, imageID string) (bool, error) {
	var resp oscgo.ReadImagesResponse
	var err error
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, _, err = conn.ImageApi.ReadImages(context.Background(), &oscgo.ReadImagesOpts{ReadImagesRequest: optional.NewInterface(oscgo.ReadImagesRequest{
			Filters: &oscgo.FiltersImage{
				ImageIds: &[]string{imageID},
			},
		})})
		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded:") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	var errString string

	if err != nil {
		// When an AMI disappears out from under a launch permission resource, we will
		// see either InvalidAMIID.NotFound or InvalidAMIID.Unavailable.
		if strings.Contains(fmt.Sprint(err), "InvalidAMIID") {
			log.Printf("[DEBUG] %s no longer exists, so we'll drop launch permission for the state", imageID)
			return false, nil
		}
		errString = err.Error()

		return false, fmt.Errorf("Error creating Outscale VM volume: %s", errString)
	}

	if len(resp.GetImages()) == 0 {
		log.Printf("[DEBUG] %s no longer exists, so we'll drop launch permission for the state", imageID)
		return false, nil
	}

	result := resp.GetImages()[0]

	if len(result.PermissionsToLaunch.GetAccountIds()) > 0 {
		return true, nil
	}
	return false, nil
}
