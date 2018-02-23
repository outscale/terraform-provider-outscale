package outscale

import (
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-outscale/osc/fcu"
)

const (
	OutscaleImageRetryTimeout       = 40 * time.Minute
	OutscaleImageDeleteRetryTimeout = 90 * time.Minute
	OutscaleImageRetryDelay         = 5 * time.Second
	OutscaleImageRetryMinTimeout    = 3 * time.Second
)

func resourceOutscaleImage() *schema.Resource {
	return &schema.Resource{
		Create: resourceImageCreate,
		Read:   resourceImageRead,
		Update: resourceImageUpdate,
		Delete: resourceImageDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(40 * time.Minute),
			Update: schema.DefaultTimeout(40 * time.Minute),
			Delete: schema.DefaultTimeout(40 * time.Minute),
		},

		Schema: getImageSchema(),
	}
}

func resourceImageCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	req := &fcu.RegisterImageInput{
		Name:       aws.String(d.Get("name").(string)),
		InstanceId: aws.String(d.Get("instance_id").(string)),
	}

	if a, aok := d.GetOk("description"); aok {
		req.Description = aws.String(a.(string))
	}
	if a, aok := d.GetOk("dry_run"); aok {
		req.DryRun = aws.Bool(a.(bool))
	}
	if a, aok := d.GetOk("no_reboot"); aok {
		req.NoReboot = aws.Bool(a.(bool))
	}

	var res *fcu.RegisterImageOutput
	var err error
	err = resource.Retry(10*time.Minute, func() *resource.RetryError {
		res, err = conn.VM.RegisterImage(req)

		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded") {
				fmt.Printf("[INFO] Request limit exceeded")
				return resource.RetryableError(err)
			}
		}

		return resource.RetryableError(err)
	})

	if err != nil {
		return err
	}

	id := *res.ImageId
	d.SetId(id)
	d.Partial(true) // make sure we record the id even if the rest of this gets interrupted
	d.Set("id", id)
	d.Set("manage_ebs_block_devices", false)
	d.SetPartial("id")
	d.SetPartial("manage_ebs_block_devices")
	d.Partial(false)

	_, err = resourceOutscaleImageWaitForAvailable(id, conn, 1)
	if err != nil {
		return err
	}

	return resourceImageUpdate(d, meta)

}

func resourceImageRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*OutscaleClient).FCU
	id := d.Id()

	req := &fcu.DescribeImagesInput{
		ImageIds: []*string{aws.String(id)},
	}

	var res *fcu.DescribeImagesOutput
	var err error
	err = resource.Retry(10*time.Minute, func() *resource.RetryError {
		res, err = client.VM.DescribeImages(req)

		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded") {
				fmt.Printf("[INFO] Request limit exceeded")
				return resource.RetryableError(err)
			}
		}

		return resource.RetryableError(err)
	})

	if err != nil {
		if ec2err, ok := err.(awserr.Error); ok && ec2err.Code() == "InvalidAMIID.NotFound" {
			fmt.Printf("[DEBUG] %s no longer exists, so we'll drop it from the state", id)
			d.SetId("")
			return nil
		}

		return err
	}

	if len(res.Images) != 1 {
		d.SetId("")
		return nil
	}

	image := res.Images[0]
	state := *image.State

	if state == "pending" {
		image, err = resourceOutscaleImageWaitForAvailable(id, client, 2)
		if err != nil {
			return err
		}
		state = *image.State
	}

	if state == "deregistered" {
		d.SetId("")
		return nil
	}

	if state != "available" {
		return fmt.Errorf("AMI has become %s", state)
	}

	d.Set("name", image.Name)
	d.Set("description", image.Description)
	d.Set("image_location", image.ImageLocation)
	d.Set("image_owner_alias", image.ImageOwnerAlias)
	d.Set("image_owner_id", image.OwnerId)
	d.Set("image_state", image.State)
	d.Set("image_type", image.ImageType)
	d.Set("architecture", image.Architecture)
	d.Set("is_public", image.Public)
	d.Set("creation_date", image.CreationDate)
	d.Set("root_device_name", image.RootDeviceName)
	d.Set("root_device_type", image.RootDeviceType)

	var ebsBlockDevs []map[string]interface{}
	var ephemeralBlockDevs []map[string]interface{}

	for _, blockDev := range image.BlockDeviceMappings {
		ephemeralBlockDev := make(map[string]interface{})

		if blockDev.Ebs != nil {
			ebsBlockDev := make(map[string]interface{})

			if blockDev.DeviceName != nil {
				ebsBlockDev["device_name"] = *blockDev.DeviceName
			}
			if blockDev.Ebs.DeleteOnTermination != nil {
				ebsBlockDev["delete_on_termination"] = *blockDev.Ebs.DeleteOnTermination
			}
			if blockDev.Ebs.Encrypted != nil {
				ebsBlockDev["encrypted"] = *blockDev.Ebs.Encrypted
			}
			if blockDev.Ebs.VolumeSize != nil {
				ebsBlockDev["volume_size"] = int(*blockDev.Ebs.VolumeSize)
			}
			if blockDev.Ebs.VolumeType != nil {
				ebsBlockDev["volume_type"] = *blockDev.Ebs.VolumeType
			}

			if blockDev.Ebs.Iops != nil {
				ebsBlockDev["iops"] = int(*blockDev.Ebs.Iops)
			}
			// The snapshot ID might not be set.
			if blockDev.Ebs.SnapshotId != nil {
				ebsBlockDev["snapshot_id"] = *blockDev.Ebs.SnapshotId
			}
			ebsBlockDevs = append(ebsBlockDevs, ebsBlockDev)

			if blockDev.DeviceName != nil {
				ephemeralBlockDev["device_name"] = *blockDev.DeviceName
			}
			if blockDev.VirtualName != nil {
				ephemeralBlockDev["virtual_name"] = *blockDev.VirtualName
			}

			ephemeralBlockDev["ebs"] = ebsBlockDevs

			ephemeralBlockDevs = append(ephemeralBlockDevs, ephemeralBlockDev)
		} else {

			if blockDev.DeviceName != nil {
				ephemeralBlockDev["device_name"] = *blockDev.DeviceName
			}
			if blockDev.VirtualName != nil {
				ephemeralBlockDev["virtual_name"] = *blockDev.VirtualName
			}

			ephemeralBlockDevs = append(ephemeralBlockDevs, ephemeralBlockDev)

		}
	}

	d.Set("block_device_mapping", ephemeralBlockDevs)
	d.Set("product_codes", getProductCodes(image.ProductCodes))
	d.Set("product_codes", getStateReason(image.StateReason))

	d.Set("tags", tagsToMap(image.Tags))

	return nil
}

func resourceImageUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	d.Partial(true)

	if err := setTags(conn, d); err != nil {
		return err
	}

	d.SetPartial("tags")

	if d.Get("description").(string) != "" {
		_, err := conn.VM.ModifyImageAttribute(&fcu.ModifyImageAttributeInput{
			ImageId: aws.String(d.Id()),
			Description: &fcu.AttributeValue{
				Value: aws.String(d.Get("description").(string)),
			},
		})
		if err != nil {
			return err
		}
		d.SetPartial("description")
	}

	d.Partial(false)

	return resourceImageRead(d, meta)
}

func resourceImageDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*OutscaleClient).FCU

	req := &fcu.DeregisterImageInput{
		ImageId: aws.String(d.Id()),
	}

	var err error
	err = resource.Retry(10*time.Minute, func() *resource.RetryError {
		_, err := client.VM.DeregisterImage(req)

		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded") {
				fmt.Printf("[INFO] Request limit exceeded")
				return resource.RetryableError(err)
			}
		}

		return resource.RetryableError(err)
	})

	if err != nil {
		return fmt.Errorf("Error deleting the image")
	}

	// Verify that the image is actually removed, if not we need to wait for it to be removed
	if err := resourceOutscaleImageWaitForDestroy(d.Id(), client); err != nil {
		return err
	}

	// No error, ami was deleted successfully
	d.SetId("")
	return nil
}

func getImageSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"architecture": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"block_device_mapping": {
			Type: schema.TypeSet,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"device_name": {
						Type:     schema.TypeString,
						Optional: true,
					},
					"ebs": {
						Type:     schema.TypeMap,
						Optional: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"delete_on_termination": {
									Type:     schema.TypeBool,
									Optional: true,
								},
								"iops": {
									Type:     schema.TypeInt,
									Optional: true,
								},
								"snapshot_id": {
									Type:     schema.TypeString,
									Optional: true,
								},
								"volume_size": {
									Type:     schema.TypeInt,
									Optional: true,
								},
								"volume_type": {
									Type:     schema.TypeString,
									Optional: true,
								},
							},
						},
					},
					"no_device": {
						Type:     schema.TypeBool,
						Optional: true,
					},
					"virtual_name": {
						Type:     schema.TypeString,
						Optional: true,
					},
				},
			},
			Optional: true,
			Computed: true,
		},
		"creation_date": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"description": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"dry_run": {
			Type:     schema.TypeBool,
			Optional: true,
		},
		"instance_id": {
			Type:     schema.TypeString,
			Required: true,
		},
		"image_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"image_location": {
			Type:     schema.TypeString,
			Computed: true,
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
			Type:     schema.TypeString,
			Computed: true,
		},
		"name": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"no_reboot": {
			Type:     schema.TypeBool,
			Optional: true,
		},
		"request_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"product_codes": {
			Type: schema.TypeList,
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
		"state_reason": {
			Type: schema.TypeMap,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"code": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"message": {
						Type:     schema.TypeString,
						Computed: true,
					},
				},
			},
			Computed: true,
		},
		"tag_set": {
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
	}
}

func resourceOutscaleImageWaitForAvailable(id string, client *fcu.Client, i int) (*fcu.Image, error) {
	fmt.Printf("MSG %s, Waiting for AMI %s to become available...", i, id)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"pending"},
		Target:     []string{"available"},
		Refresh:    ImageStateRefreshFunc(client, id),
		Timeout:    OutscaleImageRetryTimeout,
		Delay:      OutscaleImageRetryDelay,
		MinTimeout: OutscaleImageRetryMinTimeout,
	}

	info, err := stateConf.WaitForState()
	if err != nil {
		return nil, fmt.Errorf("Error waiting for AMI (%s) to be ready: %v", id, err)
	}
	return info.(*fcu.Image), nil
}

func ImageStateRefreshFunc(client *fcu.Client, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		emptyResp := &fcu.DescribeImagesOutput{}

		var resp *fcu.DescribeImagesOutput
		var err error
		err = resource.Retry(10*time.Minute, func() *resource.RetryError {
			resp, err = client.VM.DescribeImages(&fcu.DescribeImagesInput{ImageIds: []*string{aws.String(id)}})

			if err != nil {
				if strings.Contains(err.Error(), "RequestLimitExceeded") {
					fmt.Printf("[INFO] Request limit exceeded")
					return resource.RetryableError(err)
				}
			}

			return resource.RetryableError(err)
		})

		if err != nil {
			if e := fmt.Sprint(err); strings.Contains(e, "InvalidAMIID.NotFound") {
				return emptyResp, "destroyed", nil
			} else if resp != nil && len(resp.Images) == 0 {
				return emptyResp, "destroyed", nil
			} else {
				// if e := fmt.Sprint(err); strings.Contains(e, "InvalidAMIID.NotFound") {
				// 	return emptyResp, "destroyed", nil
				// }
				return emptyResp, "", fmt.Errorf("Error on refresh: %+v", err)
			}
		}

		if resp == nil || resp.Images == nil || len(resp.Images) == 0 {
			return emptyResp, "destroyed", nil
		}

		// AMI is valid, so return it's state
		return resp.Images[0], *resp.Images[0].State, nil
	}
}

func resourceOutscaleImageWaitForDestroy(id string, client *fcu.Client) error {
	fmt.Printf("Waiting for AMI %s to be deleted...", id)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"available", "pending", "failed"},
		Target:     []string{"destroyed"},
		Refresh:    ImageStateRefreshFunc(client, id),
		Timeout:    OutscaleImageDeleteRetryTimeout,
		Delay:      OutscaleImageRetryDelay,
		MinTimeout: OutscaleImageRetryTimeout,
	}

	_, err := stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Error waiting for AMI (%s) to be deleted: %v", id, err)
	}

	return nil
}
