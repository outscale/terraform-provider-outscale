package outscale

import (
	"time"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceOutscaleKeyPairImportation() *schema.Resource {
	return &schema.Resource{
		Create: resourceKeyPairImportationCreate,
		Read:   resourceKeyPairImportationRead,
		Delete: resourceKeyPairImportationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Read:   schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: getKeyPairImportationSchema(),
	}
}

func resourceKeyPairImportationCreate(d *schema.ResourceData, meta interface{}) error {
	//	conn := meta.(*OutscaleClient).FCU
	return nil
}

func resourceKeyPairImportationRead(d *schema.ResourceData, meta interface{}) error {
	//	conn := meta.(*OutscaleClient).FCU
	return nil
}

func resourceKeyPairImportationDelete(d *schema.ResourceData, meta interface{}) error {
	//	conn := meta.(*OutscaleClient).FCU
	return nil
}

func getKeyPairImportationSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		// Attributes
		"key_name": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"public_key_material": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
	}
}
