package outscale

import (
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceOutscaleImageRegister() *schema.Resource {
	return &schema.Resource{
		Create: resourceImageRegisterCreate,
		// Read:   resourceImageRegisterRead,
		// Update: resourceImageRegisterUpdate,
		// Delete: resourceImageRegisterDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"arquitecture": {
				Type:     schema.TypeString,
				Required: true,
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
						"ebs": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"delete_on_termination": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"iops": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"snapshot_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"volume_size": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"volume_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"no_device": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"virtual_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"image_location": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"root_device_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceImageRegisterCreate(d *schema.ResourceData, meta interface{}) error {
	//client := meta.(*OutscaleClient).FCU
	return nil
}
