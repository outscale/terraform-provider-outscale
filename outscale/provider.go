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
			"endpoints": EndpointsSchema(),
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
			"outscale_vm":                                ResourceOutscaleOApiVM(),
			"outscale_keypair":                           ResourceOutscaleOAPIKeyPair(),
			"outscale_image":                             ResourceOutscaleOAPIImage(),
			"outscale_internet_service_link":             ResourceOutscaleOAPIInternetServiceLink(),
			"outscale_internet_service":                  ResourceOutscaleOAPIInternetService(),
			"outscale_net":                               ResourceOutscaleOAPINet(),
			"outscale_security_group":                    ResourceOutscaleOAPISecurityGroup(),
			"outscale_outbound_rule":                     ResourceOutscaleOAPIOutboundRule(),
			"outscale_security_group_rule":               ResourceOutscaleOAPIOutboundRule(),
			"outscale_tag":                               ResourceOutscaleOAPITags(),
			"outscale_public_ip":                         ResourceOutscaleOAPIPublicIP(),
			"outscale_public_ip_link":                    ResourceOutscaleOAPIPublicIPLink(),
			"outscale_volume":                            ResourceOutscaleOAPIVolume(),
			"outscale_volumes_link":                      ResourceOutscaleOAPIVolumeLink(),
			"outscale_net_attributes":                    ResourceOutscaleOAPILinAttributes(),
			"outscale_nat_service":                       ResourceOutscaleOAPINatService(),
			"outscale_subnet":                            ResourceOutscaleOAPISubNet(),
			"outscale_route":                             ResourceOutscaleOAPIRoute(),
			"outscale_route_table":                       ResourceOutscaleOAPIRouteTable(),
			"outscale_route_table_link":                  ResourceOutscaleOAPILinkRouteTable(),
			"outscale_nic":                               ResourceOutscaleOAPINic(),
			"outscale_snapshot":                          ResourceOutscaleOAPISnapshot(),
			"outscale_image_launch_permission":           ResourceOutscaleOAPIImageLaunchPermission(),
			"outscale_net_peering":                       ResourceOutscaleOAPILinPeeringConnection(),
			"outscale_net_peering_acceptation":           ResourceOutscaleOAPILinPeeringConnectionAccepter(),
			"outscale_net_access_point":                  ResourceOutscaleNetAccessPoint(),
			"outscale_nic_link":                          ResourceOutscaleOAPINetworkInterfaceAttachment(),
			"outscale_nic_private_ip":                    ResourceOutscaleOAPINetworkInterfacePrivateIP(),
			"outscale_snapshot_attributes":               ResourcedOutscaleOAPISnapshotAttributes(),
			"outscale_dhcp_option":                       ResourceOutscaleDHCPOption(),
			"outscale_client_gateway":                    ResourceOutscaleClientGateway(),
			"outscale_virtual_gateway":                   ResourceOutscaleOAPIVirtualGateway(),
			"outscale_virtual_gateway_link":              ResourceOutscaleOAPIVirtualGatewayLink(),
			"outscale_virtual_gateway_route_propagation": ResourceOutscaleOAPIVirtualGatewayRoutePropagation(),
			"outscale_vpn_connection":                    ResourceOutscaleVPNConnection(),
			"outscale_vpn_connection_route":              ResourceOutscaleVPNConnectionRoute(),
			"outscale_access_key":                        ResourceOutscaleAccessKey(),
			"outscale_load_balancer":                     ResourceOutscaleOAPILoadBalancer(),
			"outscale_load_balancer_policy":              ResourceOutscaleAppCookieStickinessPolicy(),
			"outscale_load_balancer_vms":                 ResourceOutscaleOAPILBUAttachment(),
			"outscale_load_balancer_attributes":          ResourceOutscaleOAPILoadBalancerAttributes(),
			"outscale_load_balancer_listener_rule":       ResourceOutscaleLoadBalancerListenerRule(),
			"outscale_flexible_gpu":                      ResourceOutscaleOAPIFlexibleGpu(),
			"outscale_flexible_gpu_link":                 ResourceOutscaleOAPIFlexibleGpuLink(),
			"outscale_image_export_task":                 ResourceOutscaleOAPIIMageExportTask(),
			"outscale_server_certificate":                ResourceOutscaleOAPIServerCertificate(),
			"outscale_snapshot_export_task":              ResourceOutscaleOAPISnapshotExportTask(),
			"outscale_ca":                                ResourceOutscaleOAPICa(),
			"outscale_api_access_rule":                   ResourceOutscaleOAPIApiAccessRule(),
			"outscale_api_access_policy":                 ResourceOutscaleOAPIApiAccessPolicy(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"outscale_vm":                           DataSourceOutscaleOAPIVM(),
			"outscale_vms":                          DatasourceOutscaleOApiVMS(),
			"outscale_security_group":               DataSourceOutscaleOAPISecurityGroup(),
			"outscale_security_groups":              DataSourceOutscaleOAPISecurityGroups(),
			"outscale_image":                        DataSourceOutscaleOAPIImage(),
			"outscale_images":                       DataSourceOutscaleOAPIImages(),
			"outscale_tag":                          DataSourceOutscaleOAPITag(),
			"outscale_tags":                         DataSourceOutscaleOAPITags(),
			"outscale_public_ip":                    DataSourceOutscaleOAPIPublicIP(),
			"outscale_public_ips":                   DataSourceOutscaleOAPIPublicIPS(),
			"outscale_volume":                       DatasourceOutscaleOAPIVolume(),
			"outscale_volumes":                      DatasourceOutscaleOAPIVolumes(),
			"outscale_nat_service":                  DataSourceOutscaleOAPINatService(),
			"outscale_nat_services":                 DataSourceOutscaleOAPINatServices(),
			"outscale_keypair":                      DatasourceOutscaleOAPIKeyPair(),
			"outscale_keypairs":                     DatasourceOutscaleOAPIKeyPairs(),
			"outscale_vm_state":                     DataSourceOutscaleOAPIVMState(),
			"outscale_vm_states":                    DataSourceOutscaleOAPIVMStates(),
			"outscale_internet_service":             DatasourceOutscaleOAPIInternetService(),
			"outscale_internet_services":            DatasourceOutscaleOAPIInternetServices(),
			"outscale_subnet":                       DataSourceOutscaleOAPISubnet(),
			"outscale_subnets":                      DataSourceOutscaleOAPISubnets(),
			"outscale_net":                          DataSourceOutscaleOAPIVpc(),
			"outscale_nets":                         DataSourceOutscaleOAPIVpcs(),
			"outscale_net_attributes":               DataSourceOutscaleOAPIVpcAttr(),
			"outscale_route_table":                  DataSourceOutscaleOAPIRouteTable(),
			"outscale_route_tables":                 DataSourceOutscaleOAPIRouteTables(),
			"outscale_snapshot":                     DataSourceOutscaleOAPISnapshot(),
			"outscale_snapshots":                    DataSourceOutscaleOAPISnapshots(),
			"outscale_net_peering":                  DataSourceOutscaleOAPILinPeeringConnection(),
			"outscale_net_peerings":                 DataSourceOutscaleOAPILinPeeringsConnection(),
			"outscale_nics":                         DataSourceOutscaleOAPINics(),
			"outscale_nic":                          DataSourceOutscaleOAPINic(),
			"outscale_client_gateway":               DataSourceOutscaleClientGateway(),
			"outscale_client_gateways":              DataSourceOutscaleClientGateways(),
			"outscale_virtual_gateway":              DataSourceOutscaleOAPIVirtualGateway(),
			"outscale_virtual_gateways":             DataSourceOutscaleOAPIVirtualGateways(),
			"outscale_vpn_connection":               DataSourceOutscaleVPNConnection(),
			"outscale_vpn_connections":              DataSourceOutscaleVPNConnections(),
			"outscale_access_key":                   DataSourceOutscaleAccessKey(),
			"outscale_access_keys":                  DataSourceOutscaleAccessKeys(),
			"outscale_dhcp_option":                  DataSourceOutscaleDHCPOption(),
			"outscale_dhcp_options":                 DataSourceOutscaleDHCPOptions(),
			"outscale_load_balancer":                DataSourceOutscaleOAPILoadBalancer(),
			"outscale_load_balancer_listener_rule":  DataSourceOutscaleOAPILoadBalancerLDRule(),
			"outscale_load_balancer_listener_rules": DataSourceOutscaleOAPILoadBalancerLDRules(),
			"outscale_load_balancer_tags":           DataSourceOutscaleOAPILBUTags(),
			"outscale_load_balancer_vm_health":      DataSourceOutscaleLoadBalancerVmsHeals(),
			"outscale_load_balancers":               DataSourceOutscaleOAPILoadBalancers(),
			"outscale_vm_types":                     DataSourceOutscaleOAPIVMTypes(),
			"outscale_net_access_point":             DataSourceOutscaleNetAccessPoint(),
			"outscale_net_access_points":            DataSourceOutscaleNetAccessPoints(),
			"outscale_flexible_gpu":                 DataSourceOutscaleOAPIFlexibleGpu(),
			"outscale_flexible_gpus":                DataSourceOutscaleOAPIFlexibleGpus(),
			"outscale_subregions":                   DataSourceOutscaleOAPISubregions(),
			"outscale_regions":                      DataSourceOutscaleOAPIRegions(),
			"outscale_net_access_point_services":    DataSourceOutscaleOAPINetAccessPointServices(),
			"outscale_flexible_gpu_catalog":         DataSourceOutscaleOAPIFlexibleGpuCatalog(),
			"outscale_product_type":                 DataSourceOutscaleOAPIProductType(),
			"outscale_product_types":                DataSourceOutscaleOAPIProductTypes(),
			"outscale_quota":                        DataSourceOutscaleOAPIQuota(),
			"outscale_quotas":                       DataSourceOutscaleOAPIQuotas(),
			"outscale_image_export_task":            DataSourceOutscaleOAPIImageExportTask(),
			"outscale_image_export_tasks":           DataSourceOutscaleOAPIImageExportTasks(),
			"outscale_server_certificate":           DatasourceOutscaleOAPIServerCertificate(),
			"outscale_server_certificates":          DatasourceOutscaleOAPIServerCertificates(),
			"outscale_snapshot_export_task":         DataSourceOutscaleOAPISnapshotExportTask(),
			"outscale_snapshot_export_tasks":        DataSourceOutscaleOAPISnapshotExportTasks(),
			"outscale_ca":                           DataSourceOutscaleOAPICa(),
			"outscale_cas":                          DataSourceOutscaleOAPICas(),
			"outscale_api_access_rule":              DataSourceOutscaleOAPIApiAccessRule(),
			"outscale_api_access_rules":             DataSourceOutscaleOAPIApiAccessRules(),
			"outscale_api_access_policy":            DataSourceOutscaleOAPIApiAccessPolicy(),
			"outscale_public_catalog":               DataSourceOutscaleOAPIPublicCatalog(),
			"outscale_account":                      DataSourceAccount(),
			"outscale_accounts":                     DataSourceAccounts(),
		},

		ConfigureFunc: ProviderConfigureClient,
	}
}

func ProviderConfigureClient(d *schema.ResourceData) (interface{}, error) {
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

func EndpointsSchema() *schema.Schema {
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
