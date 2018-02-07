package outscale

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"access_key_id": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("OUTSCALE_ACCESSKEYID", nil),
				Description: "The Access Key ID for API operations.",
			},
			"secret_key_id": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("OUTSCALE_SECRETKEYID", nil),
				Description: "The Secret Key ID for API operations.",
			},
			"oapi": {
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OUTSCALE_OAPI", false),
				Description: "Enable oAPI Usage",
			},
		},

		ResourcesMap:  map[string]*schema.Resource{
		    "outscale_vm": resourceOutscaleVM(),
		},
		ConfigureFunc: providerConfigureClient,
	}
}

func providerConfigureClient(d *schema.ResourceData) (interface{}, error) {
	config := Config{
		AccessKeyId: d.Get("access_key_id").(string),
		SecretKeyId: d.Get("secret_key_id").(string),
		OApi:        d.Get("oapi").(bool),
	}
	return config.Client()
	//	return nil, nil
}
