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

func resourceOutscaleImageRegister() *schema.Resource {
	return &schema.Resource{
		Create: resourceImageRegisterCreate,
		Read:   resourceImageRegisterRead,
		Delete: resourceImageRegisterDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: getRegisterImageSchema(false),
	}
}

func resourceImageRegisterCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	request := &fcu.RegisterImageInput{}

	architecture, architectureOk := d.GetOk("architecture")
	blockDeviceMapping, blockDeviceMappingOk := d.GetOk("block_device_mapping")
	description, descriptionOk := d.GetOk("description")
	imageLocation, imageLocationOk := d.GetOk("image_location")
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
		request.BlockDeviceMappings = readBlockDeviceImage(blockDeviceMapping)
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

	_, err = resourceOutscaleImageWaitForAvailable(*registerResp.ImageId, conn, 1)
	if err != nil {
		return err
	}

	return resourceImageRegisterRead(d, meta)
}

func resourceImageRegisterRead(d *schema.ResourceData, meta interface{}) error {
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

	if err := d.Set("block_device_mapping", amiBlockDeviceMappingsReg(image.BlockDeviceMappings)); err != nil {
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

func resourceImageRegisterDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

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

func amiBlockDeviceMappingsReg(m []*fcu.BlockDeviceMapping) []map[string]interface{} {
	s := make([]map[string]interface{}, len(m))

	for k, v := range m {
		mapping := make(map[string]interface{})
		if v.Ebs != nil {
			mapping["volume_type"] = *v.Ebs.VolumeType
			mapping["volume_size"] = int(*v.Ebs.VolumeSize)
			mapping["delete_on_termination"] = aws.BoolValue(v.Ebs.DeleteOnTermination)

			if v.Ebs.Iops != nil {
				mapping["iops"] = int(aws.Int64Value(v.Ebs.Iops))
			}
			// snapshot id may not be set
			if v.Ebs.SnapshotId != nil {
				mapping["snapshot_id"] = *v.Ebs.SnapshotId
			}
		}
		if v.VirtualName != nil {
			mapping["virtual_name"] = *v.VirtualName
		}
		if v.DeviceName != nil {
			mapping["device_name"] = *v.DeviceName
		}
		if v.NoDevice != nil {
			mapping["no_device"] = *v.NoDevice
		}

		s[k] = mapping
	}
	return s
}

func readBlockDeviceImage(v interface{}) []*fcu.BlockDeviceMapping {
	maps := v.([]interface{})
	mappings := []*fcu.BlockDeviceMapping{}

	for _, m := range maps {
		f := m.(map[string]interface{})
		mapping := &fcu.BlockDeviceMapping{
			DeviceName: aws.String(f["device_name"].(string)),
		}

		if v, ok := f["no_device"]; ok && v != "" {
			mapping.NoDevice = aws.String(v.(string))
		}
		if v, ok := f["virtual_name"]; ok && v != "" {
			mapping.VirtualName = aws.String(v.(string))
		}

		ebs := &fcu.EbsBlockDevice{}

		if v, ok := f["delete_on_termination"]; ok {
			ebs.DeleteOnTermination = aws.Bool(v.(bool))
		}
		if iops, ok := f["iops"]; ok {
			if iop := iops.(int); iop != 0 {
				ebs.Iops = aws.Int64(int64(v.(int)))
			}
		}
		if v, ok := f["snapshot_id"]; ok && v != "" {
			ebs.SnapshotId = aws.String(v.(string))
		}
		if v, ok := f["volume_size"]; ok {
			if s := v.(int); s != 0 {
				ebs.VolumeSize = aws.Int64(int64(v.(int)))
			}
		}
		if v, ok := f["volume_type"]; ok && v != "" {
			ebs.VolumeType = aws.String(v.(string))
		}

		mapping.Ebs = ebs

		mappings = append(mappings, mapping)
	}

	return mappings
}

func getRegisterImageSchema(computed bool) map[string]*schema.Schema {
	// var virtualizationTypeDefault interface{}
	var deleteEbsOnTerminationDefault interface{}
	// var sriovNetSupportDefault interface{}
	var architectureDefault interface{}
	var volumeTypeDefault interface{}

	if !computed {
		// virtualizationTypeDefault = "paravirtual"
		deleteEbsOnTerminationDefault = true
		// sriovNetSupportDefault = "simple"
		architectureDefault = "i386"
		volumeTypeDefault = "standard"
	}

	return map[string]*schema.Schema{
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
			Computed: false,
			Optional: true,
			ForceNew: true,
			Default:  architectureDefault,
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
						Computed: false,
						Optional: true,
						ForceNew: true,
					},
					"no_device": {
						Type:     schema.TypeString,
						Computed: true,
						Optional: true,
						ForceNew: true,
					},
					"virtual_name": {
						Type:     schema.TypeString,
						Computed: true,
						Optional: true,
						ForceNew: true,
					},
					"delete_on_termination": &schema.Schema{
						Type:     schema.TypeBool,
						Optional: true,
						ForceNew: true,
						Computed: false,
						Default:  deleteEbsOnTerminationDefault,
					},
					"iops": &schema.Schema{
						Type:     schema.TypeInt,
						Optional: true,
						Computed: false,
						ForceNew: true,
					},
					"snapshot_id": &schema.Schema{
						Type:     schema.TypeString,
						Computed: false,
						ForceNew: true,
						Optional: true,
					},
					"volume_size": &schema.Schema{
						Type:     schema.TypeInt,
						Computed: true,
						ForceNew: true,
						Optional: true,
					},
					"volume_type": &schema.Schema{
						Type:     schema.TypeString,
						Computed: false,
						ForceNew: true,
						Optional: true,
						Default:  volumeTypeDefault,
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

		"arquitecture": {
			Type:     schema.TypeString,
			Optional: true,
			ForceNew: true,
			Computed: true,
		},
	}
}
