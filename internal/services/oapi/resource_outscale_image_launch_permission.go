package oapi

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/outscale/goutils/sdk/ptr"
	"github.com/outscale/osc-sdk-go/v3/pkg/options"
	"github.com/outscale/osc-sdk-go/v3/pkg/osc"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/services/oapi/oapihelpers"
	"github.com/outscale/terraform-provider-outscale/internal/utils"

	"github.com/spf13/cast"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceOutscaleImageLaunchPermission() *schema.Resource {
	return &schema.Resource{
		Exists:        ResourceOutscaleImageLaunchPermissionExists,
		CreateContext: ResourceOutscaleImageLaunchPermissionCreate,
		ReadContext:   ResourceOutscaleImageLaunchPermissionRead,
		DeleteContext: ResourceOutscaleImageLaunchPermissionDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(CreateDefaultTimeout),
			Read:   schema.DefaultTimeout(ReadDefaultTimeout),
			Delete: schema.DefaultTimeout(DeleteDefaultTimeout),
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
	client := meta.(*client.OutscaleClient).OSC

	imageID := d.Get("image_id").(string)
	return oapihelpers.ImageHasLaunchPermission(context.Background(), client, ReadDefaultTimeout, imageID)
}

func expandOAPIImagePermission(permissionType interface{}) (res osc.PermissionsOnResource) {
	if len(permissionType.([]interface{})) > 0 {
		permission := permissionType.([]interface{})[0].(map[string]interface{})

		if globalPermission, ok := permission["global_permission"]; ok {
			res.GlobalPermission = new(cast.ToBool(globalPermission))
		}
		if accountIDs, ok := permission["account_ids"]; ok {
			for _, accountID := range accountIDs.([]interface{}) {
				acc := append(ptr.From(res.AccountIds), accountID.(string))
				res.AccountIds = &acc
			}
		}
	}
	return
}

func ResourceOutscaleImageLaunchPermissionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.OutscaleClient).OSC
	timeout := d.Timeout(schema.TimeoutCreate)

	imageID, ok := d.GetOk("image_id")

	if !ok {
		return diag.Errorf("please provide the required attribute image_id")
	}
	log.Printf("Creating Outscale Image Launch Permission, image_id (%+v)", imageID.(string))

	permissionLaunch := osc.PermissionsOnResourceCreation{}
	if permissionAdditions, ok := d.GetOk("permission_additions"); ok {
		permissionLaunch.Additions = new(expandOAPIImagePermission(permissionAdditions))
	}
	if permissionRemovals, ok := d.GetOk("permission_removals"); ok {
		permissionLaunch.Removals = new(expandOAPIImagePermission(permissionRemovals))
	}

	request := osc.UpdateImageRequest{
		ImageId:             imageID.(string),
		PermissionsToLaunch: &permissionLaunch,
	}

	_, err := client.UpdateImage(ctx, request, options.WithRetryTimeout(timeout))

	var errString string
	if err != nil {
		errString = err.Error()

		return diag.Errorf("error creating omi launch permission: %s", errString)
	}

	d.SetId(imageID.(string))

	return ResourceOutscaleImageLaunchPermissionRead(ctx, d, meta)
}

func ResourceOutscaleImageLaunchPermissionRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.OutscaleClient).OSC

	timeout := d.Timeout(schema.TimeoutRead)

	resp, err := client.ReadImages(ctx, osc.ReadImagesRequest{
		Filters: &osc.FiltersImage{
			ImageIds: &[]string{d.Id()},
		},
	}, options.WithRetryTimeout(timeout))

	var errString string
	if err != nil {
		// When an AMI disappears out from under a launch permission resource, we will
		// see either InvalidAMIID.NotFound or InvalidAMIID.Unavailable.
		if strings.Contains(fmt.Sprint(err), "InvalidAMIID") {
			log.Printf("[DEBUG] %s no longer exists, so we'll drop launch permission for the state", d.Id())
			return nil
		}
		errString = err.Error()

		return diag.Errorf("error reading outscale image permission: %s", errString)
	}
	if resp.Images == nil || utils.IsResponseEmpty(len(*resp.Images), "ImageLaunchPermission", d.Id()) {
		d.SetId("")
		return nil
	}
	result := (*resp.Images)[0]

	if err := d.Set("description", ptr.From(result.Description)); err != nil {
		return diag.FromErr(err)
	}

	lp := make(map[string]interface{})
	perm := ptr.From(result.PermissionsToLaunch)
	lp["global_permission"] = strconv.FormatBool(ptr.From(perm.GlobalPermission))
	lp["account_ids"] = perm.AccountIds

	return diag.FromErr(d.Set("permissions_to_launch", []map[string]interface{}{lp}))
}

func ResourceOutscaleImageLaunchPermissionDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.OutscaleClient).OSC

	timeout := d.Timeout(schema.TimeoutDelete)

	imageID, ok := d.GetOk("image_id")
	if !ok {
		return diag.Errorf("please provide the required attribute image_id")
	}

	if permissionAdditions, ok := d.GetOk("permission_additions"); ok {
		permission := osc.PermissionsOnResourceCreation{}
		request := osc.UpdateImageRequest{
			ImageId: imageID.(string),
		}
		permission.Removals = new(expandOAPIImagePermission(permissionAdditions))
		request.PermissionsToLaunch = &permission

		_, err := client.UpdateImage(ctx, request, options.WithRetryTimeout(timeout))

		var errString string
		if err != nil {
			errString = err.Error()

			return diag.Errorf("error removing omi launch permission: %s", errString)
		}
	}

	d.SetId("")
	return nil
}
