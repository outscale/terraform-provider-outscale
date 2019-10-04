package outscale

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-outscale/osc/oapi"
	"github.com/terraform-providers/terraform-provider-outscale/utils"
)

func dataSourceOutscaleOAPIImage() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOutscaleOAPIImageRead,

		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			"permission": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"image_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			// Computed values.
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
			"file_location": {
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
			"image_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"image_name": {
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
			"request_id": {
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
			"permissions_to_launch": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"global_permission": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"account_id": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
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
		},
	}
}

func dataSourceOutscaleOAPIImageRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OAPI

	executableUsers, executableUsersOk := d.GetOk("permission")
	filters, filtersOk := d.GetOk("filter")
	ai, aisOk := d.GetOk("account_id")
	imageID, imageIDOk := d.GetOk("image_id")

	if executableUsersOk == false && filtersOk == false && aisOk == false && imageIDOk == false {
		return fmt.Errorf("One of executable_users, filters, or account_id must be assigned, or image_id must be provided")
	}

	params := &oapi.ReadImagesRequest{
		Filters: oapi.FiltersImage{},
	}

	if executableUsersOk {
		params.Filters.PermissionsToLaunchAccountIds = expandStringValueList(executableUsers.([]interface{}))
	}
	if filtersOk {
		params.Filters = buildOutscaleOAPIDataSourceImagesFilters(filters.(*schema.Set))
	}
	if imageIDOk {
		params.Filters.ImageIds = []string{imageID.(string)}
	}
	if aisOk {
		params.Filters.AccountIds = []string{ai.(string)}
	}

	var result *oapi.ReadImagesResponse
	var resp *oapi.POST_ReadImagesResponses
	var err error
	err = resource.Retry(20*time.Minute, func() *resource.RetryError {
		resp, err = conn.POST_ReadImages(*params)

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

		return fmt.Errorf("Error retrieving Outscale Images: %s", errString)
	}

	result = resp.OK

	if len(result.Images) < 1 {
		return fmt.Errorf("your query returned no results, please change your search criteria and try again")
	}

	if len(result.Images) > 1 {
		return fmt.Errorf("your query returned more than one result, please try a more " +
			"specific search criteria")
	}

	d.Set("request_id", result.ResponseContext.RequestId)

	return omiOAPIDescriptionAttributes(d, &result.Images[0])
}

// populate the numerous fields that the image description returns.
func omiOAPIDescriptionAttributes(d *schema.ResourceData, image *oapi.Image) error {

	d.SetId(image.ImageId)
	d.Set("architecture", image.Architecture)
	d.Set("creation_date", image.CreationDate)
	d.Set("description", image.Description)
	//Missing on swager spec
	//d.Set("hypervisor", image.Hypervisor)
	d.Set("image_id", image.ImageId)
	d.Set("file_location", image.FileLocation)
	if image.AccountAlias != "" {
		d.Set("account_alias", image.AccountAlias)
	} else {
		d.Set("account_alias", "")
	}
	d.Set("account_id", image.AccountId)
	d.Set("image_type", image.ImageType)
	d.Set("image_name", image.ImageName)
	//Missing on swager spec
	//d.Set("is_public", image.Public)
	if image.RootDeviceName != "" {
		d.Set("root_device_name", image.RootDeviceName)
	} else {
		d.Set("root_device_name", "")
	}
	d.Set("root_device_type", image.RootDeviceType)
	d.Set("state", image.State)
	//Missing on swager spec
	//d.Set("virtualization_type", image.VirtualizationType)
	// Complex types get their own functions
	if err := d.Set("block_device_mappings", omiOAPIBlockDeviceMappings(image.BlockDeviceMappings)); err != nil {
		return err
	}
	if err := d.Set("product_codes", omiOAPIProductCodes(image.ProductCodes)); err != nil {
		return err
	}
	if err := d.Set("state_comment", omiOAPIStateReason(&image.StateComment)); err != nil {
		return err
	}
	if err := d.Set("tags", getOapiTagSet(image.Tags)); err != nil {
		return err
	}

	accountIds := image.PermissionsToLaunch.AccountIds
	lp := make([]map[string]interface{}, len(accountIds))
	for k, v := range accountIds {
		l := make(map[string]interface{})
		//if image.PermissionsToLaunch.GlobalPermission != nil {
		l["global_permission"] = image.PermissionsToLaunch.GlobalPermission
		//}
		//if v.UserId != nil {
		l["account_id"] = v
		//}
		lp[k] = l
	}

	d.Set("permissions_to_launch", lp)

	return nil
}
