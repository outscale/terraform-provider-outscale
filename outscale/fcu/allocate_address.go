package fcu

import (
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceAllocateAddress() *schema.Resource {
	return &schema.Resource{
		Create: resourceAllocateAddressCreate,
		Read:   resourceAllocateAddressRead,
		Delete: resourceAllocateAddressDelete,

		Schema: getAllocateAddressSchema(),
	}
}

func getAllocateAddressSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{

		// Optional attributes

		"domain": &schema.Schema{
			Type:     schema.TypeString,
			Required: false,
			ForceNew: true,
		},

		// Computed attributes

		"allocation_id": &schema.Schema{
			Type:     schema.TypeString,
			Required: false,
			Computed: true,
		},
		"association_id": &schema.Schema{
			Type:     schema.TypeString,
			Required: false,
			Computed: true,
		},
		"instance_id": &schema.Schema{
			Type:     schema.TypeString,
			Required: false,
			Computed: true,
		},
		"network_interface_id": &schema.Schema{
			Type:     schema.TypeString,
			Required: false,
			Computed: true,
		},
		"network_interface_owner_id": &schema.Schema{
			Type:     schema.TypeString,
			Required: false,
			Computed: true,
		},
		"private_ip_address": &schema.Schema{
			Type:     schema.TypeString,
			Required: false,
			Computed: true,
		},
		"public_ip": &schema.Schema{
			Type:     schema.TypeString,
			Required: false,
			Computed: true,
		},
	}
}

func resourceAllocateAddressCreate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceAllocateAddressRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}
func resourceAllocateAddressDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}
