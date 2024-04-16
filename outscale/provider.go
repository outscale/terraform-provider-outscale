package outscale

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/outscale/terraform-provider-outscale/utils"
)

var endpointServiceNames []string

func init() {
	endpointServiceNames = []string{
		"api",
	}
}

// Provider ...
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"access_key_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The Access Key ID for API operations.",
			},
			"secret_key_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The Secret Key ID for API operations.",
			},
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The Region for API operations.",
			},
			"endpoints": endpointsSchema(),
			"x509_cert_path": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The path to your x509 cert",
			},
			"x509_key_path": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The path to your x509 key",
			},
			"insecure": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "tls insecure connection",
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"outscale_vm":                                ResourceOutscaleVM(),
			"outscale_keypair":                           ResourceOutscaleKeyPair(),
			"outscale_image":                             ResourceOutscaleImage(),
			"outscale_internet_service_link":             ResourceOutscaleInternetServiceLink(),
			"outscale_internet_service":                  ResourceOutscaleInternetService(),
			"outscale_net":                               ResourceOutscaleNet(),
			"outscale_security_group":                    ResourceOutscaleSecurityGroup(),
			"outscale_outbound_rule":                     ResourceOutscaleOutboundRule(),
			"outscale_security_group_rule":               ResourceOutscaleOutboundRule(),
			"outscale_tag":                               ResourceOutscaleTags(),
			"outscale_public_ip":                         ResourceOutscalePublicIP(),
			"outscale_public_ip_link":                    ResourceOutscalePublicIPLink(),
			"outscale_volume":                            ResourceOutscaleVolume(),
			"outscale_volumes_link":                      ResourceOutscaleVolumeLink(),
			"outscale_net_attributes":                    ResourceOutscaleLinAttributes(),
			"outscale_nat_service":                       ResourceOutscaleNatService(),
			"outscale_subnet":                            ResourceOutscaleSubNet(),
			"outscale_route":                             ResourceOutscaleRoute(),
			"outscale_route_table":                       ResourceOutscaleRouteTable(),
			"outscale_route_table_link":                  ResourceOutscaleLinkRouteTable(),
			"outscale_nic":                               ResourceOutscaleNic(),
			"outscale_snapshot":                          ResourceOutscaleSnapshot(),
			"outscale_image_launch_permission":           ResourceOutscaleImageLaunchPermission(),
			"outscale_net_peering":                       ResourceOutscaleLinPeeringConnection(),
			"outscale_net_peering_acceptation":           ResourceOutscaleLinPeeringConnectionAccepter(),
			"outscale_net_access_point":                  ResourceOutscaleNetAccessPoint(),
			"outscale_nic_link":                          ResourceOutscaleNetworkInterfaceAttachment(),
			"outscale_nic_private_ip":                    ResourceOutscaleNetworkInterfacePrivateIP(),
			"outscale_snapshot_attributes":               ResourcedOutscaleSnapshotAttributes(),
			"outscale_dhcp_option":                       ResourceOutscaleDHCPOption(),
			"outscale_client_gateway":                    ResourceOutscaleClientGateway(),
			"outscale_virtual_gateway":                   ResourceOutscaleVirtualGateway(),
			"outscale_virtual_gateway_link":              ResourceOutscaleVirtualGatewayLink(),
			"outscale_virtual_gateway_route_propagation": ResourceOutscaleVirtualGatewayRoutePropagation(),
			"outscale_vpn_connection":                    ResourceOutscaleVPNConnection(),
			"outscale_vpn_connection_route":              ResourceOutscaleVPNConnectionRoute(),
			"outscale_access_key":                        ResourceOutscaleAccessKey(),
			"outscale_load_balancer":                     ResourceOutscaleLoadBalancer(),
			"outscale_load_balancer_policy":              ResourceOutscaleAppCookieStickinessPolicy(),
			"outscale_load_balancer_vms":                 ResourceLBUAttachment(),
			"outscale_load_balancer_attributes":          ResourceOutscaleLoadBalancerAttributes(),
			"outscale_load_balancer_listener_rule":       ResourceOutscaleLoadBalancerListenerRule(),
			"outscale_flexible_gpu":                      ResourceOutscaleFlexibleGpu(),
			"outscale_flexible_gpu_link":                 ResourceOutscaleFlexibleGpuLink(),
			"outscale_image_export_task":                 ResourceOutscaleIMageExportTask(),
			"outscale_server_certificate":                ResourceOutscaleServerCertificate(),
			"outscale_snapshot_export_task":              ResourceOutscaleSnapshotExportTask(),
			"outscale_ca":                                ResourceOutscaleCa(),
			"outscale_api_access_rule":                   ResourceOutscaleApiAccessRule(),
			"outscale_api_access_policy":                 ResourceOutscaleApiAccessPolicy(),
			"outscale_main_route_table_link":             resourceLinkMainRouteTable(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"outscale_vm":                           DataSourceOutscaleVM(),
			"outscale_vms":                          DataSourceOutscaleVMS(),
			"outscale_security_group":               DataSourceOutscaleSecurityGroup(),
			"outscale_security_groups":              DataSourceOutscaleSecurityGroups(),
			"outscale_image":                        DataSourceOutscaleImage(),
			"outscale_images":                       DataSourceOutscaleImages(),
			"outscale_tag":                          DataSourceOutscaleTag(),
			"outscale_tags":                         DataSourceOutscaleTags(),
			"outscale_public_ip":                    DataSourceOutscalePublicIP(),
			"outscale_public_ips":                   DataSourceOutscalePublicIPS(),
			"outscale_volume":                       DataSourceOutscaleVolume(),
			"outscale_volumes":                      DataSourceOutscaleVolumes(),
			"outscale_nat_service":                  DataSourceOutscaleNatService(),
			"outscale_nat_services":                 DataSourceOutscaleNatServices(),
			"outscale_keypair":                      DataSourceOutscaleKeyPair(),
			"outscale_keypairs":                     DataSourceOutscaleKeyPairs(),
			"outscale_vm_state":                     DataSourceOutscaleVMState(),
			"outscale_vm_states":                    DataSourceOutscaleVMStates(),
			"outscale_internet_service":             DataSourceOutscaleInternetService(),
			"outscale_internet_services":            DataSourceOutscaleInternetServices(),
			"outscale_subnet":                       DataSourceOutscaleSubnet(),
			"outscale_subnets":                      DataSourceOutscaleSubnets(),
			"outscale_net":                          DataSourceOutscaleVpc(),
			"outscale_nets":                         DataSourceOutscaleVpcs(),
			"outscale_net_attributes":               DataSourceOutscaleVpcAttr(),
			"outscale_route_table":                  DataSourceOutscaleRouteTable(),
			"outscale_route_tables":                 DataSourceOutscaleRouteTables(),
			"outscale_snapshot":                     DataSourceOutscaleSnapshot(),
			"outscale_snapshots":                    DataSourceOutscaleSnapshots(),
			"outscale_net_peering":                  DataSourceOutscaleLinPeeringConnection(),
			"outscale_net_peerings":                 DataSourceOutscaleLinPeeringsConnection(),
			"outscale_nics":                         DataSourceOutscaleNics(),
			"outscale_nic":                          DataSourceOutscaleNic(),
			"outscale_client_gateway":               DataSourceOutscaleClientGateway(),
			"outscale_client_gateways":              DataSourceOutscaleClientGateways(),
			"outscale_virtual_gateway":              DataSourceOutscaleVirtualGateway(),
			"outscale_virtual_gateways":             DataSourceOutscaleVirtualGateways(),
			"outscale_vpn_connection":               DataSourceOutscaleVPNConnection(),
			"outscale_vpn_connections":              DataSourceOutscaleVPNConnections(),
			"outscale_access_key":                   DataSourceOutscaleAccessKey(),
			"outscale_access_keys":                  DataSourceOutscaleAccessKeys(),
			"outscale_dhcp_option":                  DataSourceOutscaleDHCPOption(),
			"outscale_dhcp_options":                 DataSourceOutscaleDHCPOptions(),
			"outscale_load_balancer":                DataSourceOutscaleLoadBalancer(),
			"outscale_load_balancer_listener_rule":  DataSourceOutscaleLoadBalancerLDRule(),
			"outscale_load_balancer_listener_rules": DataSourceOutscaleLoadBalancerLDRules(),
			"outscale_load_balancer_tags":           DataSourceOutscaleLBUTags(),
			"outscale_load_balancer_vm_health":      DataSourceOutscaleLoadBalancerVmsHeals(),
			"outscale_load_balancers":               DataSourceOutscaleLoadBalancers(),
			"outscale_vm_types":                     DataSourceOutscaleVMTypes(),
			"outscale_net_access_point":             DataSourceOutscaleNetAccessPoint(),
			"outscale_net_access_points":            DataSourceOutscaleNetAccessPoints(),
			"outscale_flexible_gpu":                 DataSourceOutscaleFlexibleGpu(),
			"outscale_flexible_gpus":                DataSourceOutscaleFlexibleGpus(),
			"outscale_subregions":                   DataSourceOutscaleSubregions(),
			"outscale_regions":                      DataSourceOutscaleRegions(),
			"outscale_net_access_point_services":    DataSourceOutscaleNetAccessPointServices(),
			"outscale_flexible_gpu_catalog":         DataSourceOutscaleFlexibleGpuCatalog(),
			"outscale_product_type":                 DataSourceOutscaleProductType(),
			"outscale_product_types":                DataSourceOutscaleProductTypes(),
			"outscale_quotas":                       DataSourceOutscaleQuotas(),
			"outscale_image_export_task":            DataSourceOutscaleImageExportTask(),
			"outscale_image_export_tasks":           DataSourceOutscaleImageExportTasks(),
			"outscale_server_certificate":           DataSourceOutscaleServerCertificate(),
			"outscale_server_certificates":          DataSourceOutscaleServerCertificates(),
			"outscale_snapshot_export_task":         DataSourceOutscaleSnapshotExportTask(),
			"outscale_snapshot_export_tasks":        DataSourceOutscaleSnapshotExportTasks(),
			"outscale_ca":                           DataSourceOutscaleCa(),
			"outscale_cas":                          DataSourceOutscaleCas(),
			"outscale_api_access_rule":              DataSourceOutscaleApiAccessRule(),
			"outscale_api_access_rules":             DataSourceOutscaleApiAccessRules(),
			"outscale_api_access_policy":            DataSourceOutscaleApiAccessPolicy(),
			"outscale_public_catalog":               DataSourceOutscalePublicCatalog(),
			"outscale_account":                      DataSourceAccount(),
			"outscale_accounts":                     DataSourceAccounts(),
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
		Insecure:    d.Get("insecure").(bool),
	}

	setProviderDefaultEnv(&config)
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

func setProviderDefaultEnv(conf *Config) {
	if conf.AccessKeyID == "" {
		if accessKeyId := utils.GetEnvVariableValue([]string{"OSC_ACCESS_KEY", "OUTSCALE_ACCESSKEYID"}); accessKeyId != "" {
			conf.AccessKeyID = accessKeyId
		}
	}
	if conf.SecretKeyID == "" {
		if secretKeyId := utils.GetEnvVariableValue([]string{"OSC_SECRET_KEY", "OUTSCALE_SECRETKEYID"}); secretKeyId != "" {
			conf.SecretKeyID = secretKeyId
		}
	}

	if conf.Region == "" {
		if region := utils.GetEnvVariableValue([]string{"OSC_REGION", "OUTSCALE_REGION"}); region != "" {
			conf.Region = region
		}
	}

	if conf.X509cert == "" {
		if x509Cert := utils.GetEnvVariableValue([]string{"OSC_X509_CLIENT_CERT", "OUTSCALE_X509CERT"}); x509Cert != "" {
			conf.X509cert = x509Cert
		}
	}

	if conf.X509key == "" {
		if x509Key := utils.GetEnvVariableValue([]string{"OSC_X509_CLIENT_KEY", "OUTSCALE_X509KEY"}); x509Key != "" {
			conf.X509key = x509Key
		}
	}

	if len(conf.Endpoints) == 0 {
		if endpoints := utils.GetEnvVariableValue([]string{"OSC_ENDPOINT_API", "OUTSCALE_OAPI_URL"}); endpoints != "" {
			endpointsAttributes := make(map[string]interface{})
			endpointsAttributes["api"] = endpoints
			conf.Endpoints = endpointsAttributes
		}
	}
}
