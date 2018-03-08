package outscale

import (
	"os"
	"strconv"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

// Provider ...
func Provider() terraform.ResourceProvider {

	fcu := "fcu"

	o := os.Getenv("OUTSCALE_OAPI")

	isoapi, err := strconv.ParseBool(o)
	if err != nil {
		isoapi = false
	}

	if isoapi {
		fcu = "oapi"
	}

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
			"region": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("OUTSCALE_REGION", nil),
				Description: "The Region for API operations.",
			},
			"oapi": {
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OUTSCALE_OAPI", false),
				Description: "Enable oAPI Usage",
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"outscale_vm":                 GetResource(fcu, "outscale_vm")(),
			"outscale_key_pair":           GetResource(fcu, "outscale_key_pair")(),
			"outscale_image":              GetResource(fcu, "outscale_image")(),
			"outscale_firewall_rules_set": GetResource(fcu, "outscale_firewall_rules_set")(),
			"outscale_outbound_rule":      GetResource(fcu, "outscale_outbound_rule")(),
			"outscale_inbound_rule":       GetResource(fcu, "outscale_inbound_rule")(),
			"outscale_tag":                GetResource(fcu, "outscale_tag")(),
			"outscale_public_ip":          GetResource(fcu, "outscale_public_ip")(),
			"outscale_public_ip_link":     GetResource(fcu, "outscale_public_ip_link")(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"outscale_vm":                 GetDatasource(fcu, "outscale_vm")(),
			"outscale_vms":                GetDatasource(fcu, "outscale_vms")(),
			"outscale_firewall_rule_set":  GetDatasource(fcu, "outscale_firewall_rule_set")(),
			"outscale_firewall_rules_set": GetDatasource(fcu, "outscale_firewall_rules_set")(),
			"outscale_image":      GetDatasource(fcu, "outscale_image")(),
			"outscale_images":     GetDatasource(fcu, "outscale_images")(),
			"outscale_tag":        GetDatasource(fcu, "outscale_tag")(),
			"outscale_tags":       GetDatasource(fcu, "outscale_tags")(),
			"outscale_public_ip":  GetDatasource(fcu, "outscale_public_ip")(),
			"outscale_public_ips": GetDatasource(fcu, "outscale_public_ips")(),
		},

		ConfigureFunc: providerConfigureClient,
	}
}

func providerConfigureClient(d *schema.ResourceData) (interface{}, error) {
	config := Config{
		AccessKeyID: d.Get("access_key_id").(string),
		SecretKeyID: d.Get("secret_key_id").(string),
		Region:      d.Get("region").(string),
		OApi:        d.Get("oapi").(bool),
	}
	return config.Client()
}
