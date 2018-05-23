package outscale

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-outscale/osc/fcu"
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
			"permission": &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"create": &schema.Schema{
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"global_permission": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
									},
									"account_id": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"delete": &schema.Schema{
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"global_permission": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
									},
									"account_id": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
					},
				},
			},
			"description": &schema.Schema{
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"value": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"permissions": &schema.Schema{
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
			"request_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceOutscaleOAPIImageLaunchPermissionExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	conn := meta.(*OutscaleClient).FCU

	imageID := d.Get("image_id").(string)
	return hasOAPILaunchPermission(conn, imageID)
}

func resourceOutscaleOAPIImageLaunchPermissionCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	imageID, iok := d.GetOk("image_id")
	permission, lok := d.GetOk("permission")

	if iok {
		return fmt.Errorf("please provide the required attribute imageID")
	}

	request := &fcu.ModifyImageAttributeInput{
		ImageId: aws.String(imageID.(string)),
	}

	if lok {
		request.Attribute = aws.String("launchPermission")
		launchPermission := &fcu.LaunchPermissionModifications{}

		l := permission.([]interface{})

		lp := l[0].(map[string]interface{})

		if a, ok := lp["create"]; ok {
			ad := a.([]interface{})
			if len(ad) > 0 {
				add := make([]*fcu.LaunchPermission, len(ad))
				for k, v := range ad {
					att := v.(map[string]interface{})
					at := &fcu.LaunchPermission{}
					if g, ok := att["global_permission"]; ok {
						at.Group = aws.String(g.(string))
					}
					if g, ok := att["account_id"]; ok {
						at.UserId = aws.String(g.(string))
					}
					add[k] = at
				}
				launchPermission.Add = add
			}
		}
		request.LaunchPermission = launchPermission
	}

	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		var err error
		_, err = conn.VM.ModifyImageAttribute(request)
		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded:") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("error creating ami launch permission: %s", err)
	}

	d.SetId(imageID.(string))
	d.Set("description", map[string]string{"value": ""})
	d.Set("permissions", make([]map[string]interface{}, 0))
	return nil
}

func resourceOutscaleOAPIImageLaunchPermissionRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	var attrs *fcu.DescribeImageAttributeOutput
	var err error
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		attrs, err = conn.VM.DescribeImageAttribute(&fcu.DescribeImageAttributeInput{
			ImageId:   aws.String(d.Id()),
			Attribute: aws.String("launchPermission"),
		})
		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded:") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		// When an AMI disappears out from under a launch permission resource, we will
		// see either InvalidAMIID.NotFound or InvalidAMIID.Unavailable.
		if strings.Contains(fmt.Sprint(err), "InvalidAMIID") {
			log.Printf("[DEBUG] %s no longer exists, so we'll drop launch permission for the state", d.Id())
			return nil
		}
		return err
	}

	d.Set("request_id", attrs.RequestId)
	d.Set("description", map[string]string{"value": aws.StringValue(attrs.Description.Value)})

	lp := make([]map[string]interface{}, len(attrs.LaunchPermissions))
	for k, v := range attrs.LaunchPermissions {
		l := make(map[string]interface{})
		if v.Group != nil {
			l["global_permission"] = *v.Group
		}
		if v.UserId != nil {
			l["account_id"] = *v.UserId
		}
		lp[k] = l
	}

	d.Set("permissions", lp)

	return nil
}

func resourceOutscaleOAPIImageLaunchPermissionDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	imageID, iok := d.GetOk("image_id")
	permission, lok := d.GetOk("permission")

	if iok {
		return fmt.Errorf("please provide the required attribute imageID")
	}

	request := &fcu.ModifyImageAttributeInput{
		ImageId: aws.String(imageID.(string)),
	}

	if lok {
		request.Attribute = aws.String("launchPermission")
		launchPermission := &fcu.LaunchPermissionModifications{}

		lps := permission.([]interface{})
		lp := lps[0].(map[string]interface{})

		if a, ok := lp["delete"]; ok {
			ad := a.([]interface{})
			if len(ad) > 0 {
				remove := make([]*fcu.LaunchPermission, len(ad))
				for k, v := range ad {
					att := v.(map[string]interface{})
					at := &fcu.LaunchPermission{}
					if g, ok := att["global_permission"]; ok {
						at.Group = aws.String(g.(string))
					}
					if g, ok := att["account_id"]; ok {
						at.UserId = aws.String(g.(string))
					}
					remove[k] = at
				}
				launchPermission.Remove = remove
			}
		}
		request.LaunchPermission = launchPermission
	}

	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		var err error
		_, err = conn.VM.ModifyImageAttribute(request)
		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded:") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		if strings.Contains(fmt.Sprint(err), "InvalidAMIID") {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("error removing ami launch permission: %s", err)
	}

	d.SetId("")

	return nil
}

func hasOAPILaunchPermission(conn *fcu.Client, imageID string) (bool, error) {
	var err error
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, err = conn.VM.DescribeImageAttribute(&fcu.DescribeImageAttributeInput{
			ImageId:   aws.String(imageID),
			Attribute: aws.String("launchPermission"),
		})
		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded:") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		if strings.Contains(fmt.Sprint(err), "InvalidAMIID") {
			log.Printf("[DEBUG] %s no longer exists, so we'll drop launch permission from the state", imageID)
			return false, nil
		}
		return false, err
	}

	// for _, lp := range attrs.LaunchPermissions {
	// 	if *lp.UserId == account_id {
	// 		return true, nil
	// 	}
	// }
	return true, nil
}
