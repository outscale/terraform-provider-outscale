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

func resourceOutscaleImageLaunchPermission() *schema.Resource {
	return &schema.Resource{
		Create: resourceOutscaleImageLaunchPermissionCreate,
		Read:   resourceOutscaleImageLaunchPermissionRead,
		Delete: resourceOutscaleImageLaunchPermissionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"image_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"launch_permission": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"add": &schema.Schema{
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"group": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
									},
									"user_id": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"remove": &schema.Schema{
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"group": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
									},
									"user_id": &schema.Schema{
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
			"launch_permissions": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"group": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"user_id": &schema.Schema{
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

func resourceOutscaleImageLaunchPermissionCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	id, iok := d.GetOk("image_id")
	lp, lok := d.GetOk("launch_permission")

	fmt.Println("Creating Outscale Image Launch Permission, image_id", id.(string))

	if !iok {
		return fmt.Errorf("please provide the required attribute image_id")
	}

	request := &fcu.ModifyImageAttributeInput{
		ImageId: aws.String(id.(string)),
	}

	if lok {
		request.Attribute = aws.String("launchPermission")
		launchPermission := &fcu.LaunchPermissionModifications{}

		l := lp.([]interface{})

		lp := l[0].(map[string]interface{})

		if a, ok := lp["add"]; ok {
			ad := a.([]interface{})
			if len(ad) > 0 {
				add := make([]*fcu.LaunchPermission, len(ad))
				for k, v := range ad {
					att := v.(map[string]interface{})
					at := &fcu.LaunchPermission{}
					if g, ok := att["group"]; ok {
						at.Group = aws.String(g.(string))
					}
					if g, ok := att["user_id"]; ok {
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

	d.SetId(id.(string))
	d.Set("description", map[string]string{"value": ""})
	d.Set("launch_permissions", make([]map[string]interface{}, 0))

	return nil
}

func resourceOutscaleImageLaunchPermissionRead(d *schema.ResourceData, meta interface{}) error {
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
		if strings.Contains(fmt.Sprint(err), "InvalidAMIID") {
			log.Printf("[DEBUG] %s no longer exists, so we'll drop launch permission for the state", d.Id())
			return nil
		}
		return err
	}

	d.Set("request_id", attrs.RequestId)

	d.Set("description", map[string]string{"value": ""})
	if attrs.Description != nil {
		d.Set("description", map[string]string{"value": aws.StringValue(attrs.Description.Value)})
	}

	lp := make([]map[string]interface{}, len(attrs.LaunchPermissions))
	for k, v := range attrs.LaunchPermissions {
		l := make(map[string]interface{})
		if v.Group != nil {
			l["group"] = *v.Group
		} else {
			l["group"] = ""
		}
		if v.UserId != nil {
			l["user_id"] = *v.UserId
		} else {
			l["user_id"] = ""
		}
		lp[k] = l
	}

	return d.Set("launch_permissions", lp)
}

func resourceOutscaleImageLaunchPermissionDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	ID := d.Get("image_id")
	lp, lok := d.GetOk("launch_permission")

	request := &fcu.ModifyImageAttributeInput{
		ImageId: aws.String(ID.(string)),
	}

	if lok {
		request.Attribute = aws.String("launchPermission")
		launchPermission := &fcu.LaunchPermissionModifications{}

		lps := lp.([]interface{})
		lp := lps[0].(map[string]interface{})

		if a, ok := lp["remove"]; ok {
			ad := a.([]interface{})
			if len(ad) > 0 {
				remove := make([]*fcu.LaunchPermission, len(ad))
				for k, v := range ad {
					att := v.(map[string]interface{})
					at := &fcu.LaunchPermission{}
					if g, ok := att["group"]; ok {
						at.Group = aws.String(g.(string))
					}
					if g, ok := att["user_id"]; ok {
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

func hasLaunchPermission(conn *fcu.Client, ID string) (bool, error) {

	var err error
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, err = conn.VM.DescribeImageAttribute(&fcu.DescribeImageAttributeInput{
			ImageId:   aws.String(ID),
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
			log.Printf("[DEBUG] %s no longer exists, so we'll drop launch permission from the state", ID)
			return false, nil
		}
		return false, err
	}

	return false, nil
}
