package outscale

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

var endpointServiceNames []string

func init() {
	endpointServiceNames = []string{
		"api",
	}
}

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
			"endpoints": endpointsSchema(),
			"x509_cert_path": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OUTSCALE_X509CERT", nil),
				Description: "The path to your x509 cert",
			},
			"x509_key_path": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OUTSCALE_X509KEY", nil),
				Description: "The path to your x509 key",
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"outscale_vm":                                resourceVM(),
			"outscale_keypair":                           resourceKeyPair(),
			"outscale_image":                             resourceImage(),
			"outscale_internet_service_link":             resourceInternetServiceLink(),
			"outscale_internet_service":                  resourceInternetService(),
			"outscale_net":                               resourceNet(),
			"outscale_security_group":                    resourceSecurityGroup(),
			"outscale_outbound_rule":                     resourceOutboundRule(),
			"outscale_security_group_rule":               resourceOutboundRule(),
			"outscale_tag":                               resourceTags(),
			"outscale_public_ip":                         resourcePublicIP(),
			"outscale_public_ip_link":                    resourcePublicIPLink(),
			"outscale_volume":                            resourceVolume(),
			"outscale_volumes_link":                      resourceVolumeLink(),
			"outscale_net_attributes":                    resourceLinAttributes(),
			"outscale_nat_service":                       resourceNatService(),
			"outscale_subnet":                            resourceSubNet(),
			"outscale_route":                             resourceRoute(),
			"outscale_route_table":                       resourceRouteTable(),
			"outscale_route_table_link":                  resourceLinkRouteTable(),
			"outscale_nic":                               resourceNic(),
			"outscale_snapshot":                          resourceSnapshot(),
			"outscale_image_launch_permission":           resourceImageLaunchPermission(),
			"outscale_net_peering":                       resourceLinPeeringConnection(),
			"outscale_net_peering_acceptation":           resourceLinPeeringConnectionAccepter(),
			"outscale_net_access_point":                  resourceNetAccessPoint(),
			"outscale_nic_link":                          resourceNetworkInterfaceAttachment(),
			"outscale_nic_private_ip":                    resourceNetworkInterfacePrivateIP(),
			"outscale_snapshot_attributes":               resourcedSnapshotAttributes(),
			"outscale_dhcp_option":                       resourceDHCPOption(),
			"outscale_client_gateway":                    resourceClientGateway(),
			"outscale_virtual_gateway":                   resourceVirtualGateway(),
			"outscale_virtual_gateway_link":              resourceVirtualGatewayLink(),
			"outscale_virtual_gateway_route_propagation": resourceVirtualGatewayRoutePropagation(),
			"outscale_vpn_connection":                    resourceVPNConnection(),
			"outscale_vpn_connection_route":              resourceVPNConnectionRoute(),
			"outscale_access_key":                        resourceAccessKey(),
			"outscale_load_balancer":                     resourceLoadBalancer(),
			"outscale_load_balancer_policy":              resourceAppCookieStickinessPolicy(),
			"outscale_load_balancer_vms":                 resourceLBUAttachment(),
			"outscale_load_balancer_attributes":          resourceLoadBalancerAttributes(),
			"outscale_load_balancer_listener_rule":       resourceLoadBalancerListenerRule(),
			"outscale_flexible_gpu":                      resourceFlexibleGpu(),
			"outscale_flexible_gpu_link":                 resourceFlexibleGpuLink(),
			"outscale_image_export_task":                 resourceIMageExportTask(),
			"outscale_server_certificate":                resourceServerCertificate(),
			"outscale_snapshot_export_task":              resourceSnapshotExportTask(),
			"outscale_ca":                                resourceCa(),
			"outscale_api_access_rule":                   resourceApiAccessRule(),
			"outscale_api_access_policy":                 resourceApiAccessPolicy(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"outscale_vm":                           dataSourceVM(),
			"outscale_vms":                          dataSourceVMS(),
			"outscale_security_group":               dataSourceSecurityGroup(),
			"outscale_security_groups":              dataSourceSecurityGroups(),
			"outscale_image":                        dataSourceImage(),
			"outscale_images":                       dataSourceImages(),
			"outscale_tag":                          dataSourceTag(),
			"outscale_tags":                         dataSourceTags(),
			"outscale_public_ip":                    dataSourcePublicIP(),
			"outscale_public_ips":                   dataSourcePublicIPS(),
			"outscale_volume":                       dataSourceVolume(),
			"outscale_volumes":                      dataSourceVolumes(),
			"outscale_nat_service":                  dataSourceNatService(),
			"outscale_nat_services":                 dataSourceNatServices(),
			"outscale_keypair":                      dataSourceKeyPair(),
			"outscale_keypairs":                     dataSourceKeyPairs(),
			"outscale_vm_state":                     dataSourceVMState(),
			"outscale_vm_states":                    dataSourceVMStates(),
			"outscale_internet_service":             dataSourceInternetService(),
			"outscale_internet_services":            dataSourceInternetServices(),
			"outscale_subnet":                       dataSourceSubnet(),
			"outscale_subnets":                      dataSourceSubnets(),
			"outscale_net":                          dataSourceVpc(),
			"outscale_nets":                         dataSourceVpcs(),
			"outscale_net_attributes":               dataSourceVpcAttr(),
			"outscale_route_table":                  dataSourceRouteTable(),
			"outscale_route_tables":                 dataSourceRouteTables(),
			"outscale_snapshot":                     dataSourceSnapshot(),
			"outscale_snapshots":                    dataSourceSnapshots(),
			"outscale_net_peering":                  dataSourceLinPeeringConnection(),
			"outscale_net_peerings":                 dataSourceLinPeeringsConnection(),
			"outscale_nics":                         dataSourceNics(),
			"outscale_nic":                          dataSourceNic(),
			"outscale_client_gateway":               dataSourceClientGateway(),
			"outscale_client_gateways":              dataSourceClientGateways(),
			"outscale_virtual_gateway":              dataSourceVirtualGateway(),
			"outscale_virtual_gateways":             dataSourceVirtualGateways(),
			"outscale_vpn_connection":               dataSourceVPNConnection(),
			"outscale_vpn_connections":              dataSourceVPNConnections(),
			"outscale_access_key":                   dataSourceAccessKey(),
			"outscale_access_keys":                  dataSourceAccessKeys(),
			"outscale_dhcp_option":                  dataSourceDHCPOption(),
			"outscale_dhcp_options":                 dataSourceDHCPOptions(),
			"outscale_load_balancer":                dataSourceLoadBalancer(),
			"outscale_load_balancer_listener_rule":  dataSourceLoadBalancerLDRule(),
			"outscale_load_balancer_listener_rules": dataSourceLoadBalancerLDRules(),
			"outscale_load_balancer_tags":           dataSourceLBUTags(),
			"outscale_load_balancer_vm_health":      dataSourceLoadBalancerVmsHeals(),
			"outscale_load_balancers":               dataSourceLoadBalancers(),
			"outscale_vm_types":                     dataSourceVMTypes(),
			"outscale_net_access_point":             dataSourceNetAccessPoint(),
			"outscale_net_access_points":            dataSourceNetAccessPoints(),
			"outscale_flexible_gpu":                 dataSourceFlexibleGpu(),
			"outscale_flexible_gpus":                dataSourceFlexibleGpus(),
			"outscale_subregions":                   dataSourceSubregions(),
			"outscale_regions":                      dataSourceRegions(),
			"outscale_net_access_point_services":    dataSourceNetAccessPointServices(),
			"outscale_flexible_gpu_catalog":         dataSourceFlexibleGpuCatalog(),
			"outscale_product_type":                 dataSourceProductType(),
			"outscale_product_types":                dataSourceProductTypes(),
			"outscale_quota":                        dataSourceQuota(),
			"outscale_quotas":                       dataSourceQuotas(),
			"outscale_image_export_task":            dataSourceImageExportTask(),
			"outscale_image_export_tasks":           dataSourceImageExportTasks(),
			"outscale_server_certificate":           dataSourceServerCertificate(),
			"outscale_server_certificates":          dataSourceServerCertificates(),
			"outscale_snapshot_export_task":         dataSourceSnapshotExportTask(),
			"outscale_snapshot_export_tasks":        dataSourceSnapshotExportTasks(),
			"outscale_ca":                           dataSourceCa(),
			"outscale_cas":                          dataSourceCas(),
			"outscale_api_access_rule":              dataSourceApiAccessRule(),
			"outscale_api_access_rules":             dataSourceApiAccessRules(),
			"outscale_api_access_policy":            dataSourceApiAccessPolicy(),
			"outscale_public_catalog":               dataSourcePublicCatalog(),
			"outscale_accounts":                     dataSourceAccounts(),
		},

		ConfigureFunc: providerConfigureClient,
	}
}

func providerConfigureClient(d *schema.ResourceData) (interface{}, error) {
	config := Config{
		AccessKeyID: d.Get("access_key_id").(string),
		SecretKeyID: d.Get("secret_key_id").(string),
		Region:      d.Get("region").(string),
		Endpoints:   make(map[string]interface{}),
		X509cert:    d.Get("x509_cert_path").(string),
		X509key:     d.Get("x509_key_path").(string),
	}

	endpointsSet := d.Get("endpoints").(*schema.Set)

	for _, endpointsSetI := range endpointsSet.List() {
		endpoints := endpointsSetI.(map[string]interface{})
		for _, endpointServiceName := range endpointServiceNames {
			config.Endpoints[endpointServiceName] = endpoints[endpointServiceName].(string)
		}
	}

	return config.Client()
}

func endpointsSchema() *schema.Schema {
	endpointsAttributes := make(map[string]*schema.Schema)

	for _, endpointServiceName := range endpointServiceNames {
		endpointsAttributes[endpointServiceName] = &schema.Schema{
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "",
			Description: "Use this to override the default service endpoint URL",
		}
	}

	return &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		Elem: &schema.Resource{
			Schema: endpointsAttributes,
		},
	}
}
