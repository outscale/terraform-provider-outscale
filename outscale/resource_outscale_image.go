package outscale

import (
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-outscale/osc/fcu"
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

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"dry_run": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"no_reboot": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"client_token": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"architecture": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"creation_date": {
				Type:     schema.TypeString,
				Computed: true,
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
			"block_device_mapping": {
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
						"virtual_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ebs": {
							Type:     schema.TypeMap,
							Computed: true,
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
		},
	}
}

func resourceImageCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	req := &fcu.CreateImageInput{
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

	var res *fcu.CreateImageOutput
	var err error
	err = resource.Retry(40*time.Minute, func() *resource.RetryError {
		res, err = conn.VM.CreateImage(req)

		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded") {
				fmt.Printf("[INFO] Request limit exceeded")
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}

		return resource.RetryableError(err)
	})

	if err != nil {
		return err
	}

	d.SetId(*res.ImageId)
	d.Set("image_id", *res.ImageId)

	_, err = resourceOutscaleImageWaitForAvailable(*res.ImageId, conn, 1)
	if err != nil {
		return err
	}

	return resourceImageUpdate(d, meta)

}

func resourceImageRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*OutscaleClient).FCU
	ID := d.Id()

	req := &fcu.DescribeImagesInput{
		ImageIds: []*string{aws.String(ID)},
	}

	var res *fcu.DescribeImagesOutput
	var err error
	err = resource.Retry(40*time.Minute, func() *resource.RetryError {
		res, err = client.VM.DescribeImages(req)

		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded") {
				fmt.Printf("[INFO] Request limit exceeded")
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}

		return nil
	})

	if err != nil {
		if strings.Contains(err.Error(), "InvalidAMIID.NotFound") {
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
		image, err = resourceOutscaleImageWaitForAvailable(ID, client, 2)
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
		return fmt.Errorf("OMI has become %s", state)
	}

	d.SetId(*image.ImageId)
	d.Set("request_id", res.RequestId)
	d.Set("architecture", aws.StringValue(image.Architecture))
	d.Set("client_token", aws.StringValue(image.ClientToken))
	d.Set("creation_date", aws.StringValue(image.CreationDate))
	d.Set("description", aws.StringValue(image.Description))
	d.Set("hypervisor", aws.StringValue(image.Hypervisor))
	d.Set("image_id", aws.StringValue(image.ImageId))
	d.Set("image_location", aws.StringValue(image.ImageLocation))
	d.Set("image_owner_alias", aws.StringValue(image.ImageOwnerAlias))
	d.Set("image_owner_id", aws.StringValue(image.OwnerId))
	d.Set("image_type", aws.StringValue(image.ImageType))
	d.Set("name", aws.StringValue(image.Name))
	d.Set("is_public", aws.BoolValue(image.Public))
	d.Set("root_device_name", aws.StringValue(image.RootDeviceName))
	d.Set("root_device_type", aws.StringValue(image.RootDeviceType))
	d.Set("image_state", aws.StringValue(image.State))

	if err := d.Set("block_device_mapping", amiBlockDeviceMappings(image.BlockDeviceMappings)); err != nil {
		return err
	}
	if err := d.Set("product_codes", amiProductCodes(image.ProductCodes)); err != nil {
		return err
	}
	if err := d.Set("state_reason", amiStateReason(image.StateReason)); err != nil {
		return err
	}

	return d.Set("tag_set", tagsToMap(image.Tags))
}

func resourceImageUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	d.Partial(true)

	if err := setTags(conn, d); err != nil {
		return err
	}

	d.SetPartial("tag_set")

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
	err = resource.Retry(40*time.Minute, func() *resource.RetryError {
		_, err := client.VM.DeregisterImage(req)

		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("Error deleting the image")
	}

	if err := resourceOutscaleImageWaitForDestroy(d.Id(), client); err != nil {
		return err
	}

	d.SetId("")
	return nil
}

func resourceOutscaleImageWaitForAvailable(ID string, client *fcu.Client, i int) (*fcu.Image, error) {
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"pending"},
		Target:     []string{"available"},
		Refresh:    ImageStateRefreshFunc(client, ID),
		Timeout:    OutscaleImageRetryTimeout,
		Delay:      OutscaleImageRetryDelay,
		MinTimeout: OutscaleImageRetryMinTimeout,
	}

	info, err := stateConf.WaitForState()
	if err != nil {
		return nil, fmt.Errorf("Error waiting for OMI (%s) to be ready: %v", ID, err)
	}
	return info.(*fcu.Image), nil
}

// ImageStateRefreshFunc ...
func ImageStateRefreshFunc(client *fcu.Client, ID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		emptyResp := &fcu.DescribeImagesOutput{}

		var resp *fcu.DescribeImagesOutput
		var err error
		err = resource.Retry(15*time.Minute, func() *resource.RetryError {
			resp, err = client.VM.DescribeImages(&fcu.DescribeImagesInput{ImageIds: []*string{aws.String(ID)}})

			if err != nil {
				if strings.Contains(err.Error(), "RequestLimitExceeded") {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)

			}

			return nil
		})

		if err != nil {
			if e := fmt.Sprint(err); strings.Contains(e, "InvalidAMIID.NotFound") {
				return emptyResp, "destroyed", nil

			} else if resp != nil && len(resp.Images) == 0 {
				return emptyResp, "destroyed", nil
			} else {
				return emptyResp, "", fmt.Errorf("Error on refresh: %+v", err)
			}
		}

		if resp == nil || resp.Images == nil || len(resp.Images) == 0 {
			return emptyResp, "destroyed", nil
		}

		return resp.Images[0], *resp.Images[0].State, nil
	}
}

func resourceOutscaleImageWaitForDestroy(ID string, client *fcu.Client) error {
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"available", "pending", "failed"},
		Target:     []string{"destroyed"},
		Refresh:    ImageStateRefreshFunc(client, ID),
		Timeout:    OutscaleImageDeleteRetryTimeout,
		Delay:      OutscaleImageRetryDelay,
		MinTimeout: OutscaleImageRetryTimeout,
	}

	_, err := stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Error waiting for OMI (%s) to be deleted: %v", ID, err)
	}

	return nil
}
