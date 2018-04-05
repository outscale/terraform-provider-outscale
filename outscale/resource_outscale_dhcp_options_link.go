package outscale

import (
	"time"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceOutscaleDHCPOptionLink() *schema.Resource {
	return &schema.Resource{
		Create: resourceOutscaleDHCPOptionLinkCreate,
		Read:   resourceOutscaleDHCPOptionLinkRead,
		Delete: resourceOutscaleDHCPOptionLinkDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: getDHCPOptionLinkSchema(),
	}
}

func getDHCPOptionLinkSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"dhcp_options_id": {
			Type:     schema.TypeString,
			Required: true,
		},
		"vpc_id": {
			Type:     schema.TypeString,
			Required: true,
		},
	}
}

func resourceOutscaleDHCPOptionLinkCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	return nil
}

func resourceOutscaleDHCPOptionLinkRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	return nil
}

func resourceOutscaleDHCPOptionLinkDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU
	return nil

}
