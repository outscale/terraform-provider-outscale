package outscale

import "github.com/hashicorp/terraform/helper/schema"

func resourceOutscaleLinInternetGateway() *schema.Resource {
	return &schema.Resource{
		Create: resourceOutscaleLinInternetGatewayCreate,
		Read:   resourceOutscaleLinInternetGatewayRead,
		Delete: resourceOutscaleLinInternetGatewayDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: getLinInternetGatewaySchema(),
	}
}

func resourceOutscaleLinInternetGatewayCreate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceOutscaleLinInternetGatewayRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceOutscaleLinInternetGatewayDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func getLinInternetGatewaySchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		// Attributes
		"attachement_set": {
			Type:     schema.TypeSet,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"state": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"vpc_id": {
						Type:     schema.TypeString,
						Computed: true,
					},
				},
			},
		},
		"internet_gateway_id": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"tag_set": dataSourceTagsSchema(),
	}
}
