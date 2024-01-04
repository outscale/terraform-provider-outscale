package outscale

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/terraform-providers/terraform-provider-outscale/utils"
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
			"outscale_vm":                                resourceOutscaleOApiVM(),
			"outscale_keypair":                           resourceOutscaleOAPIKeyPair(),
			"outscale_image":                             resourceOutscaleOAPIImage(),
			"outscale_internet_service_link":             resourceOutscaleOAPIInternetServiceLink(),
			"outscale_internet_service":                  resourceOutscaleOAPIInternetService(),
			"outscale_net":                               resourceOutscaleOAPINet(),
			"outscale_security_group":                    resourceOutscaleOAPISecurityGroup(),
			"outscale_outbound_rule":                     resourceOutscaleOAPIOutboundRule(),
			"outscale_security_group_rule":               resourceOutscaleOAPIOutboundRule(),
			"outscale_tag":                               resourceOutscaleOAPITags(),
			"outscale_public_ip":                         resourceOutscaleOAPIPublicIP(),
			"outscale_public_ip_link":                    resourceOutscaleOAPIPublicIPLink(),
			"outscale_volume":                            resourceOutscaleOAPIVolume(),
			"outscale_volumes_link":                      resourceOutscaleOAPIVolumeLink(),
			"outscale_net_attributes":                    resourceOutscaleOAPILinAttributes(),
			"outscale_nat_service":                       resourceOutscaleOAPINatService(),
			"outscale_subnet":                            resourceOutscaleOAPISubNet(),
			"outscale_route":                             resourceOutscaleOAPIRoute(),
			"outscale_route_table":                       resourceOutscaleOAPIRouteTable(),
			"outscale_route_table_link":                  resourceOutscaleOAPILinkRouteTable(),
			"outscale_nic":                               resourceOutscaleOAPINic(),
			"outscale_snapshot":                          resourceOutscaleOAPISnapshot(),
			"outscale_image_launch_permission":           resourceOutscaleOAPIImageLaunchPermission(),
			"outscale_net_peering":                       resourceOutscaleOAPILinPeeringConnection(),
			"outscale_net_peering_acceptation":           resourceOutscaleOAPILinPeeringConnectionAccepter(),
			"outscale_net_access_point":                  resourceOutscaleNetAccessPoint(),
			"outscale_nic_link":                          resourceOutscaleOAPINetworkInterfaceAttachment(),
			"outscale_nic_private_ip":                    resourceOutscaleOAPINetworkInterfacePrivateIP(),
			"outscale_snapshot_attributes":               resourcedOutscaleOAPISnapshotAttributes(),
			"outscale_dhcp_option":                       resourceOutscaleDHCPOption(),
			"outscale_client_gateway":                    resourceOutscaleClientGateway(),
			"outscale_virtual_gateway":                   resourceOutscaleOAPIVirtualGateway(),
			"outscale_virtual_gateway_link":              resourceOutscaleOAPIVirtualGatewayLink(),
			"outscale_virtual_gateway_route_propagation": resourceOutscaleOAPIVirtualGatewayRoutePropagation(),
			"outscale_vpn_connection":                    resourceOutscaleVPNConnection(),
			"outscale_vpn_connection_route":              resourceOutscaleVPNConnectionRoute(),
			"outscale_access_key":                        resourceOutscaleAccessKey(),
			"outscale_load_balancer":                     resourceOutscaleOAPILoadBalancer(),
			"outscale_load_balancer_policy":              resourceOutscaleAppCookieStickinessPolicy(),
			"outscale_load_balancer_vms":                 resourceLBUAttachment(),
			"outscale_load_balancer_attributes":          resourceOutscaleOAPILoadBalancerAttributes(),
			"outscale_load_balancer_listener_rule":       resourceOutscaleLoadBalancerListenerRule(),
			"outscale_flexible_gpu":                      resourceOutscaleOAPIFlexibleGpu(),
			"outscale_flexible_gpu_link":                 resourceOutscaleOAPIFlexibleGpuLink(),
			"outscale_image_export_task":                 resourceOutscaleOAPIIMageExportTask(),
			"outscale_server_certificate":                resourceOutscaleOAPIServerCertificate(),
			"outscale_snapshot_export_task":              resourceOutscaleOAPISnapshotExportTask(),
			"outscale_ca":                                resourceOutscaleOAPICa(),
			"outscale_api_access_rule":                   resourceOutscaleOAPIApiAccessRule(),
			"outscale_api_access_policy":                 resourceOutscaleOAPIApiAccessPolicy(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"outscale_vm":                           dataSourceOutscaleOAPIVM(),
			"outscale_vms":                          datasourceOutscaleOApiVMS(),
			"outscale_security_group":               dataSourceOutscaleOAPISecurityGroup(),
			"outscale_security_groups":              dataSourceOutscaleOAPISecurityGroups(),
			"outscale_image":                        dataSourceOutscaleOAPIImage(),
			"outscale_images":                       dataSourceOutscaleOAPIImages(),
			"outscale_tag":                          dataSourceOutscaleOAPITag(),
			"outscale_tags":                         dataSourceOutscaleOAPITags(),
			"outscale_public_ip":                    dataSourceOutscaleOAPIPublicIP(),
			"outscale_public_ips":                   dataSourceOutscaleOAPIPublicIPS(),
			"outscale_volume":                       datasourceOutscaleOAPIVolume(),
			"outscale_volumes":                      datasourceOutscaleOAPIVolumes(),
			"outscale_nat_service":                  dataSourceOutscaleOAPINatService(),
			"outscale_nat_services":                 dataSourceOutscaleOAPINatServices(),
			"outscale_keypair":                      datasourceOutscaleOAPIKeyPair(),
			"outscale_keypairs":                     datasourceOutscaleOAPIKeyPairs(),
			"outscale_vm_state":                     dataSourceOutscaleOAPIVMState(),
			"outscale_vm_states":                    dataSourceOutscaleOAPIVMStates(),
			"outscale_internet_service":             datasourceOutscaleOAPIInternetService(),
			"outscale_internet_services":            datasourceOutscaleOAPIInternetServices(),
			"outscale_subnet":                       dataSourceOutscaleOAPISubnet(),
			"outscale_subnets":                      dataSourceOutscaleOAPISubnets(),
			"outscale_net":                          dataSourceOutscaleOAPIVpc(),
			"outscale_nets":                         dataSourceOutscaleOAPIVpcs(),
			"outscale_net_attributes":               dataSourceOutscaleOAPIVpcAttr(),
			"outscale_route_table":                  dataSourceOutscaleOAPIRouteTable(),
			"outscale_route_tables":                 dataSourceOutscaleOAPIRouteTables(),
			"outscale_snapshot":                     dataSourceOutscaleOAPISnapshot(),
			"outscale_snapshots":                    dataSourceOutscaleOAPISnapshots(),
			"outscale_net_peering":                  dataSourceOutscaleOAPILinPeeringConnection(),
			"outscale_net_peerings":                 dataSourceOutscaleOAPILinPeeringsConnection(),
			"outscale_nics":                         dataSourceOutscaleOAPINics(),
			"outscale_nic":                          dataSourceOutscaleOAPINic(),
			"outscale_client_gateway":               dataSourceOutscaleClientGateway(),
			"outscale_client_gateways":              dataSourceOutscaleClientGateways(),
			"outscale_virtual_gateway":              dataSourceOutscaleOAPIVirtualGateway(),
			"outscale_virtual_gateways":             dataSourceOutscaleOAPIVirtualGateways(),
			"outscale_vpn_connection":               dataSourceOutscaleVPNConnection(),
			"outscale_vpn_connections":              dataSourceOutscaleVPNConnections(),
			"outscale_access_key":                   dataSourceOutscaleAccessKey(),
			"outscale_access_keys":                  dataSourceOutscaleAccessKeys(),
			"outscale_dhcp_option":                  dataSourceOutscaleDHCPOption(),
			"outscale_dhcp_options":                 dataSourceOutscaleDHCPOptions(),
			"outscale_load_balancer":                dataSourceOutscaleOAPILoadBalancer(),
			"outscale_load_balancer_listener_rule":  dataSourceOutscaleOAPILoadBalancerLDRule(),
			"outscale_load_balancer_listener_rules": dataSourceOutscaleOAPILoadBalancerLDRules(),
			"outscale_load_balancer_tags":           dataSourceOutscaleOAPILBUTags(),
			"outscale_load_balancer_vm_health":      dataSourceOutscaleLoadBalancerVmsHeals(),
			"outscale_load_balancers":               dataSourceOutscaleOAPILoadBalancers(),
			"outscale_vm_types":                     dataSourceOutscaleOAPIVMTypes(),
			"outscale_net_access_point":             dataSourceOutscaleNetAccessPoint(),
			"outscale_net_access_points":            dataSourceOutscaleNetAccessPoints(),
			"outscale_flexible_gpu":                 dataSourceOutscaleOAPIFlexibleGpu(),
			"outscale_flexible_gpus":                dataSourceOutscaleOAPIFlexibleGpus(),
			"outscale_subregions":                   dataSourceOutscaleOAPISubregions(),
			"outscale_regions":                      dataSourceOutscaleOAPIRegions(),
			"outscale_net_access_point_services":    dataSourceOutscaleOAPINetAccessPointServices(),
			"outscale_flexible_gpu_catalog":         dataSourceOutscaleOAPIFlexibleGpuCatalog(),
			"outscale_product_type":                 dataSourceOutscaleOAPIProductType(),
			"outscale_product_types":                dataSourceOutscaleOAPIProductTypes(),
			"outscale_quotas":                       dataSourceOutscaleOAPIQuotas(),
			"outscale_image_export_task":            dataSourceOutscaleOAPIImageExportTask(),
			"outscale_image_export_tasks":           dataSourceOutscaleOAPIImageExportTasks(),
			"outscale_server_certificate":           datasourceOutscaleOAPIServerCertificate(),
			"outscale_server_certificates":          datasourceOutscaleOAPIServerCertificates(),
			"outscale_snapshot_export_task":         dataSourceOutscaleOAPISnapshotExportTask(),
			"outscale_snapshot_export_tasks":        dataSourceOutscaleOAPISnapshotExportTasks(),
			"outscale_ca":                           dataSourceOutscaleOAPICa(),
			"outscale_cas":                          dataSourceOutscaleOAPICas(),
			"outscale_api_access_rule":              dataSourceOutscaleOAPIApiAccessRule(),
			"outscale_api_access_rules":             dataSourceOutscaleOAPIApiAccessRules(),
			"outscale_api_access_policy":            dataSourceOutscaleOAPIApiAccessPolicy(),
			"outscale_public_catalog":               dataSourceOutscaleOAPIPublicCatalog(),
			"outscale_account":                      dataSourceAccount(),
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

	/*
		if data.Endpoints.IsNull() {
			if endpoints := getEnvVariableValue([]string{"OSC_ENDPOINT_API", "OUTSCALE_OAPI_URL"}); endpoints != "" {
				data.Endpoints = types.StringValue(endpoints)
			}
		}
	*/
}
