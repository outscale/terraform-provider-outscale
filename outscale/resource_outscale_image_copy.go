package outscale

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-outscale/osc/fcu"
)

func resourceOutscaleImageCopy() *schema.Resource {
	return &schema.Resource{
		Create: resourceImageCopyCreate,
		Read:   resourceImageRead,
		Update: resourceImageUpdate,
		Delete: resourceImageDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
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
			// Complex computed values
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

			//Argument
			"client_token": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			// "description": {
			// 	Type:     schema.TypeString,
			// 	Optional: true,
			// 	ForceNew: true,
			// 	Computed: true,
			// },
			// "name": {
			// 	Type:     schema.TypeString,
			// 	Computed: true,
			// 	Optional: true,
			// },
			"source_image_id": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"source_region": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			// "image_id": {
			// 	Type:     schema.TypeString,
			// 	Computed: true,
			// },
			// "request_id": {
			// 	Type:     schema.TypeString,
			// 	Computed: true,
			// },
		},
	}
}

func resourceImageCopyCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*OutscaleClient).FCU

	req := &fcu.CopyImageInput{
		Name:          aws.String(d.Get("name").(string)),
		Description:   aws.String(d.Get("description").(string)),
		SourceImageId: aws.String(d.Get("source_image_id").(string)),
		SourceRegion:  aws.String(d.Get("source_region").(string)),
		// Encrypted:     aws.Bool(d.Get("encrypted").(bool)),
	}

	if v, ok := d.GetOk("kms_key_id"); ok {
		req.KmsKeyId = aws.String(v.(string))
	}

	res, err := client.VM.CopyImage(req)
	if err != nil {
		return err
	}

	id := *res.ImageId
	d.SetId(id)
	d.Partial(true) // make sure we record the id even if the rest of this gets interrupted
	d.Set("image_id", id)
	d.Set("manage_ebs_snapshots", true)
	d.SetPartial("image_id")
	d.SetPartial("manage_ebs_snapshots")
	d.Partial(false)

	_, err = resourceOutscaleImageWaitForAvailable(id, client, 1)
	if err != nil {
		return err
	}

	return resourceImageUpdate(d, meta)
}
