package outscale

import (
	"bytes"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-outscale/osc/fcu"
)

func resourceOutscaleOAPIImage() *schema.Resource {
	return &schema.Resource{
		Create: resourceOAPIImageCreate,
		Read:   resourceOAPIImageRead,
		Update: resourceOAPIImageUpdate,
		Delete: resourceOAPIImageDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(40 * time.Minute),
			Update: schema.DefaultTimeout(40 * time.Minute),
			Delete: schema.DefaultTimeout(40 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"vm_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"no_reboot": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"description": {
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
			"osu_location": {
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
			"type": {
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
			"tag": dataSourceTagsSchema(),
		},
	}
}

func resourceOAPIImageCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	req := &fcu.RegisterImageInput{
		Name:       aws.String(d.Get("name").(string)),
		InstanceId: aws.String(d.Get("vm_id").(string)),
	}

	if a, aok := d.GetOk("description"); aok {
		req.Description = aws.String(a.(string))
	}
	if a, aok := d.GetOk("no_reboot"); aok {
		req.NoReboot = aws.Bool(a.(bool))
	}

	var res *fcu.RegisterImageOutput
	var err error
	err = resource.Retry(40*time.Minute, func() *resource.RetryError {
		res, err = conn.VM.RegisterImage(req)

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

	id := *res.ImageId
	d.SetId(id)
	d.Set("image_id", id)
	d.Partial(true) // make sure we record the id even if the rest of this gets interrupted
	d.Set("id", id)
	d.SetPartial("id")
	d.Partial(false)

	_, err = resourceOutscaleOAPIImageWaitForAvailable(id, conn, 1)
	if err != nil {
		return err
	}

	return resourceOAPIImageUpdate(d, meta)

}

func resourceOAPIImageRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*OutscaleClient).FCU
	id := d.Id()

	req := &fcu.DescribeImagesInput{
		ImageIds: []*string{aws.String(id)},
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

		return resource.RetryableError(err)
	})

	if err != nil {
		if strings.Contains(err.Error(), "InvalidAMIID.NotFound") {
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
		image, err = resourceOutscaleOAPIImageWaitForAvailable(id, client, 2)
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
	d.Set("architecture", *image.Architecture)
	if image.CreationDate != nil {
		d.Set("creation_date", *image.CreationDate)
	}
	if image.Description != nil {
		d.Set("description", *image.Description)
	}
	d.Set("hypervisor", *image.Hypervisor)
	d.Set("image_id", *image.ImageId)
	d.Set("osu_location", *image.ImageLocation)
	if image.ImageOwnerAlias != nil {
		d.Set("account_alias", *image.ImageOwnerAlias)
	}
	d.Set("account_id", *image.OwnerId)
	d.Set("type", *image.ImageType)
	d.Set("name", *image.Name)
	d.Set("is_public", *image.Public)
	if image.RootDeviceName != nil {
		d.Set("root_device_name", *image.RootDeviceName)
	}
	d.Set("root_device_type", *image.RootDeviceType)
	d.Set("state", *image.State)

	if err := d.Set("block_device_mappings", omiOAPIBlockDeviceMappings(image.BlockDeviceMappings)); err != nil {
		return err
	}
	if err := d.Set("product_codes", omiOAPIProductCodes(image.ProductCodes)); err != nil {
		return err
	}
	if err := d.Set("state_comment", omiOAPIStateReason(image.StateReason)); err != nil {
		return err
	}
	if err := d.Set("tag", dataSourceTags(image.Tags)); err != nil {
		return err
	}

	return nil
}

func resourceOAPIImageUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	d.Partial(true)

	if err := setTags(conn, d); err != nil {
		return err
	}

	d.SetPartial("tag")

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

	return resourceOAPIImageRead(d, meta)
}

func resourceOAPIImageDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*OutscaleClient).FCU

	req := &fcu.DeregisterImageInput{
		ImageId: aws.String(d.Id()),
	}

	var err error
	err = resource.Retry(40*time.Minute, func() *resource.RetryError {
		_, err := client.VM.DeregisterImage(req)

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
		return fmt.Errorf("Error deleting the image")
	}

	if err := resourceOutscaleOAPIImageWaitForDestroy(d.Id(), client); err != nil {
		return err
	}

	d.SetId("")
	return nil
}

func resourceOutscaleOAPIImageWaitForAvailable(id string, client *fcu.Client, i int) (*fcu.Image, error) {
	fmt.Printf("MSG %s, Waiting for OMI %s to become available...", i, id)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"pending"},
		Target:     []string{"available"},
		Refresh:    ImageOAPIStateRefreshFunc(client, id),
		Timeout:    OutscaleImageRetryTimeout,
		Delay:      OutscaleImageRetryDelay,
		MinTimeout: OutscaleImageRetryMinTimeout,
	}

	info, err := stateConf.WaitForState()
	if err != nil {
		return nil, fmt.Errorf("Error waiting for OMI (%s) to be ready: %v", id, err)
	}
	return info.(*fcu.Image), nil
}

func ImageOAPIStateRefreshFunc(client *fcu.Client, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		emptyResp := &fcu.DescribeImagesOutput{}

		var resp *fcu.DescribeImagesOutput
		var err error
		err = resource.Retry(15*time.Minute, func() *resource.RetryError {
			resp, err = client.VM.DescribeImages(&fcu.DescribeImagesInput{ImageIds: []*string{aws.String(id)}})

			if err != nil {
				if strings.Contains(err.Error(), "RequestLimitExceeded") {
					log.Printf("[INFO] Request limit exceeded")
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)

			}

			return resource.NonRetryableError(err)
		})

		if err != nil {
			if e := fmt.Sprint(err); strings.Contains(e, "InvalidAMIID.NotFound") {
				log.Printf("[INFO] OMI %s state %s", id, "destroyed")
				return emptyResp, "destroyed", nil

			} else if resp != nil && len(resp.Images) == 0 {
				log.Printf("[INFO] OMI %s state %s", id, "destroyed")
				return emptyResp, "destroyed", nil
			} else {
				return emptyResp, "", fmt.Errorf("Error on refresh: %+v", err)
			}
		}

		if resp == nil || resp.Images == nil || len(resp.Images) == 0 {
			return emptyResp, "destroyed", nil
		}

		log.Printf("[INFO] OMI %s state %s", *resp.Images[0].ImageId, *resp.Images[0].State)

		// OMI is valid, so return it's state
		return resp.Images[0], *resp.Images[0].State, nil
	}
}

func resourceOutscaleOAPIImageWaitForDestroy(id string, client *fcu.Client) error {
	fmt.Printf("Waiting for OMI %s to be deleted...", id)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"available", "pending", "failed"},
		Target:     []string{"destroyed"},
		Refresh:    ImageOAPIStateRefreshFunc(client, id),
		Timeout:    OutscaleImageDeleteRetryTimeout,
		Delay:      OutscaleImageRetryDelay,
		MinTimeout: OutscaleImageRetryTimeout,
	}

	_, err := stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Error waiting for OMI (%s) to be deleted: %v", id, err)
	}

	return nil
}

// Returns a set of block device mappings.
func omiOAPIBlockDeviceMappings(m []*fcu.BlockDeviceMapping) []map[string]interface{} {
	bdm := make([]map[string]interface{}, len(m))
	for k, v := range m {
		mapping := map[string]interface{}{
			"device_name": *v.DeviceName,
		}
		if v.Ebs != nil {
			bsu := map[string]interface{}{
				"delete_on_vm_termination": fmt.Sprintf("%t", *v.Ebs.DeleteOnTermination),
				"volume_size":              fmt.Sprintf("%d", *v.Ebs.VolumeSize),
				"type":                     *v.Ebs.VolumeType,
			}

			if v.Ebs.Iops != nil {
				bsu["iops"] = fmt.Sprintf("%d", *v.Ebs.Iops)
			} else {
				bsu["iops"] = "0"
			}
			if v.Ebs.SnapshotId != nil {
				bsu["snapshot_id"] = *v.Ebs.SnapshotId
			}

			mapping["bsu"] = bsu
		}
		if v.VirtualName != nil {
			mapping["virtual_device_name"] = *v.VirtualName
		}
		log.Printf("[DEBUG] outscale_image - adding block device mapping: %v", mapping)
		bdm[k] = mapping
	}
	return bdm
}

// Returns a set of product codes.
func omiOAPIProductCodes(m []*fcu.ProductCode) *schema.Set {
	s := &schema.Set{
		F: omiOAPIProductCodesHash,
	}
	for _, v := range m {
		code := map[string]interface{}{
			"product_code": *v.ProductCode,
			"type":         *v.Type,
		}
		s.Add(code)
	}
	return s
}

// Returns the state reason.
func omiOAPIStateReason(m *fcu.StateReason) map[string]interface{} {
	s := make(map[string]interface{})
	if m != nil {
		s["state_code"] = *m.Code
		s["message"] = *m.Message
	} else {
		s["state_code"] = "UNSET"
		s["message"] = "UNSET"
	}
	return s
}

func omiOAPIProductCodesHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})
	buf.WriteString(fmt.Sprintf("%s-", m["product_code_id"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["product_code_type"].(string)))
	return hashcode.String(buf.String())
}
