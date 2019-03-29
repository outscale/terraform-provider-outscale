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
			"launch_permission_add": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
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

	fmt.Println("Creating Outscale Image Launch Permission, image_id", id.(string))

	if !iok {
		return fmt.Errorf("please provide the required attribute image_id")
	}

	request := &fcu.ModifyImageAttributeInput{
		ImageId: aws.String(id.(string)),
	}

	if v, ok := d.GetOk("launch_permission_add"); ok {
		request.Attribute = aws.String("launchPermission")

		add := v.([]interface{})

		if len(add) > 0 {
			a := make([]*fcu.LaunchPermission, len(add))
			for k, v1 := range add {
				data := v1.(map[string]interface{})
				a[k] = &fcu.LaunchPermission{
					UserId: aws.String(data["user_id"].(string)),
					Group:  aws.String(data["group"].(string)),
				}
			}
			request.LaunchPermission = &fcu.LaunchPermissionModifications{Add: a}
		}
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

	d.SetId(fmt.Sprintf("%s-%s", id, "lp"))
	return resourceOutscaleImageLaunchPermissionRead(d, meta)
}

func resourceOutscaleImageLaunchPermissionRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU
	id := d.Get("image_id").(string)
	var attrs *fcu.DescribeImageAttributeOutput
	var err error
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		attrs, err = conn.VM.DescribeImageAttribute(&fcu.DescribeImageAttributeInput{
			ImageId:   aws.String(id),
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
	desc := make(map[string]interface{})

	if attrs.Description != nil {
		desc["value"] = aws.StringValue(attrs.Description.Value)
	}
	d.Set("description", desc)

	lp := make([]map[string]interface{}, len(attrs.LaunchPermissions))
	for k, v := range attrs.LaunchPermissions {
		l := make(map[string]interface{})
		l["group"] = aws.StringValue(v.Group)
		l["user_id"] = aws.StringValue(v.UserId)
		lp[k] = l
	}

	return d.Set("launch_permissions", lp)
}

func resourceOutscaleImageLaunchPermissionDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	ID := d.Get("image_id")
	lp, lok := d.GetOk("launch_permission_add")

	request := &fcu.ModifyImageAttributeInput{
		ImageId: aws.String(ID.(string)),
	}

	if lok {
		remove := lp.([]interface{})
		a := make([]*fcu.LaunchPermission, 0)
		request.Attribute = aws.String("launchPermission")
		for _, v1 := range remove {
			data := v1.(map[string]interface{})
			item := &fcu.LaunchPermission{
				UserId: aws.String(data["user_id"].(string)),
				Group:  aws.String(data["group"].(string)),
			}
			a = append(a, item)
		}
		request.LaunchPermission = &fcu.LaunchPermissionModifications{Remove: a}
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
