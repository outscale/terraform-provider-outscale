package outscale

import (
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceOutscaleOAPIImageRegister() *schema.Resource {
	return &schema.Resource{
		Create: resourceOAPIImageRegisterCreate,
		Read:   resourceOAPIImageRead,
		Delete: resourceOAPIImageRegisterDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			//Image
			"instance_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"dry_run": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"no_reboot": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"architecture": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
				ForceNew: true,
			},
			"creation_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"image_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"osu_location": {
				Type:     schema.TypeString,
				Computed: true,
				ForceNew: true,
				Optional: true,
			},
			"image_owner_alias": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"image_owner_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"image_state": {
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
				Computed: true,
				Optional: true,
				ForceNew: true,
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
			"block_device_mapping": {
				Type:     schema.TypeList,
				Computed: true,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"device_name": {
							Type:     schema.TypeString,
							Computed: true,
							Optional: true,
							ForceNew: true,
						},
						"no_device": {
							Type:     schema.TypeString,
							Computed: true,
							Optional: true,
							ForceNew: true,
						},
						"virtual_device_name": {
							Type:     schema.TypeString,
							Computed: true,
							Optional: true,
							ForceNew: true,
						},
						"bsu": {
							Type:     schema.TypeMap,
							Computed: true,
							Optional: true,
							ForceNew: true,
						},
					},
				},
			},
			"product_codes": {
				Type:     schema.TypeSet,
				Computed: true,
				Set:      amiProductCodesHash,
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
			"state_reason": {
				Type:     schema.TypeMap,
				Computed: true,
			},
			"tag_set": dataSourceTagsSchema(),

			// Image Register
			"arquitecture": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
		},
	}
}

func resourceOAPIImageRegisterCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	request := &fcu.RegisterImageInput{}

	architecture, architectureOk := d.GetOk("architecture")
	blockDeviceMapping, blockDeviceMappingOk := d.GetOk("block_device_mapping")
	description, descriptionOk := d.GetOk("description")
	imageLocation, imageLocationOk := d.GetOk("osu_location")
	name, nameOk := d.GetOk("name")
	rootDeviceName, rootDeviceNameOk := d.GetOk("root_device_name")
	instanceID, instanceIDOk := d.GetOk("instance_id")

	if !nameOk && !instanceIDOk {
		return fmt.Errorf("please provide the required attributes name and instance_id")
	}

	if architectureOk {
		request.Architecture = aws.String(architecture.(string))
	}
	if blockDeviceMappingOk {
		maps := blockDeviceMapping.([]interface{})
		mappings := []*fcu.BlockDeviceMapping{}

		for _, m := range maps {
			f := m.(map[string]interface{})
			mapping := &fcu.BlockDeviceMapping{
				DeviceName: aws.String(f["device_name"].(string)),
			}

			e := f["ebs"].(map[string]interface{})
			var del bool
			if e["delete_on_termination"].(string) == "0" {
				del = false
			} else {
				del = true
			}

			ebs := &fcu.EbsBlockDevice{
				DeleteOnTermination: aws.Bool(del),
			}

			mapping.Ebs = ebs

			mappings = append(mappings, mapping)
		}

		request.BlockDeviceMappings = mappings
	}
	if descriptionOk {
		request.Description = aws.String(description.(string))
	}
	if imageLocationOk {
		request.ImageLocation = aws.String(imageLocation.(string))
	}
	if rootDeviceNameOk {
		request.RootDeviceName = aws.String(rootDeviceName.(string))
	}
	if instanceIDOk {
		request.InstanceId = aws.String(instanceID.(string))
	}

	request.Name = aws.String(name.(string))

	var registerResp *fcu.RegisterImageOutput
	var err error

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {

		registerResp, err = conn.VM.RegisterImage(request)
		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded:") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("Error register image %s", err)
	}

	d.SetId(*registerResp.ImageId)
	d.Set("image_id", *registerResp.ImageId)

	//TODO
	// _, err = resourceOutscaleOAPIImageWaitForAvailable(*registerResp.ImageId, conn, 1)
	// if err != nil {
	// 	return err
	// }

	return resourceOAPIImageRead(d, meta)
}

func resourceOAPIImageRegisterDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	var err error

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {

		_, err = conn.VM.DeregisterImage(&fcu.DeregisterImageInput{
			ImageId: aws.String(d.Id()),
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

		return fmt.Errorf("Error Deregister image %s", err)
	}
	return nil
}
