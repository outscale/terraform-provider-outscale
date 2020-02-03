package outscale

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

// Provider ...
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
			"outscale_vm":                      resourceOutscaleOApiVM(),
			"outscale_keypair":                 resourceOutscaleOAPIKeyPair(),
			"outscale_image":                   resourceOutscaleOAPIImage(),
			"outscale_internet_service_link":   resourceOutscaleOAPIInternetServiceLink(),
			"outscale_internet_service":        resourceOutscaleOAPIInternetService(),
			"outscale_net":                     resourceOutscaleOAPINet(),
			"outscale_security_group":          resourceOutscaleOAPISecurityGroup(),
			"outscale_outbound_rule":           resourceOutscaleOAPIOutboundRule(),
			"outscale_security_group_rule":     resourceOutscaleOAPIOutboundRule(),
			"outscale_tag":                     resourceOutscaleOAPITags(),
			"outscale_public_ip":               resourceOutscaleOAPIPublicIP(),
			"outscale_public_ip_link":          resourceOutscaleOAPIPublicIPLink(),
			"outscale_volume":                  resourceOutscaleOAPIVolume(),
			"outscale_volumes_link":            resourceOutscaleOAPIVolumeLink(),
			"outscale_net_attributes":          resourceOutscaleOAPILinAttributes(),
			"outscale_nat_service":             resourceOutscaleOAPINatService(),
			"outscale_subnet":                  resourceOutscaleOAPISubNet(),
			"outscale_route":                   resourceOutscaleOAPIRoute(),
			"outscale_route_table":             resourceOutscaleOAPIRouteTable(),
			"outscale_route_table_link":        resourceOutscaleOAPILinkRouteTable(),
			"outscale_nic":                     resourceOutscaleOAPINic(),
			"outscale_snapshot":                resourceOutscaleOAPISnapshot(),
			"outscale_image_launch_permission": resourceOutscaleOAPIImageLaunchPermission(),
			"outscale_net_peering":             resourceOutscaleOAPILinPeeringConnection(),
			"outscale_net_peering_acceptation": resourceOutscaleOAPILinPeeringConnectionAccepter(),
			"outscale_nic_link":                resourceOutscaleOAPINetworkInterfaceAttachment(),
			"outscale_nic_private_ip":          resourceOutscaleOAPINetworkInterfacePrivateIP(),
			"outscale_snapshot_attributes":     resourcedOutscaleOAPISnapshotAttributes(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"outscale_vm":                dataSourceOutscaleOAPIVM(),
			"outscale_vms":               datasourceOutscaleOApiVMS(),
			"outscale_security_group":    dataSourceOutscaleOAPISecurityGroup(),
			"outscale_security_groups":   dataSourceOutscaleOAPISecurityGroups(),
			"outscale_image":             dataSourceOutscaleOAPIImage(),
			"outscale_images":            dataSourceOutscaleOAPIImages(),
			"outscale_tag":               dataSourceOutscaleOAPITag(),
			"outscale_tags":              dataSourceOutscaleOAPITags(),
			"outscale_public_ip":         dataSourceOutscaleOAPIPublicIP(),
			"outscale_public_ips":        dataSourceOutscaleOAPIPublicIPS(),
			"outscale_volume":            datasourceOutscaleOAPIVolume(),
			"outscale_volumes":           datasourceOutscaleOAPIVolumes(),
			"outscale_nat_service":       dataSourceOutscaleOAPINatService(),
			"outscale_nat_services":      dataSourceOutscaleOAPINatServices(),
			"outscale_keypair":           datasourceOutscaleOAPIKeyPair(),
			"outscale_keypairs":          datasourceOutscaleOAPIKeyPairs(),
			"outscale_vm_state":          dataSourceOutscaleOAPIVMState(),
			"outscale_vms_state":         dataSourceOutscaleOAPIVMSState(),
			"outscale_internet_service":  datasourceOutscaleOAPIInternetService(),
			"outscale_internet_services": datasourceOutscaleOAPIInternetServices(),
			"outscale_subnet":            dataSourceOutscaleOAPISubnet(),
			"outscale_subnets":           dataSourceOutscaleOAPISubnets(),
			"outscale_net":               dataSourceOutscaleOAPIVpc(),
			"outscale_nets":              dataSourceOutscaleOAPIVpcs(),
			"outscale_net_attributes":    dataSourceOutscaleOAPIVpcAttr(),
			"outscale_route_table":       dataSourceOutscaleOAPIRouteTable(),
			"outscale_route_tables":      dataSourceOutscaleOAPIRouteTables(),
			"outscale_net_peering":       dataSourceOutscaleOAPILinPeeringConnection(),
			"outscale_net_peerings":      dataSourceOutscaleOAPILinPeeringsConnection(),
			"outscale_nics":              dataSourceOutscaleOAPINics(),
			"outscale_nic":               dataSourceOutscaleOAPINic(),
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
