package outscale

import (
	"errors"
	"fmt"
	"os"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/outscale/terraform-provider-outscale/utils"
	"github.com/tidwall/gjson"
)

var endpointServiceNames = []string{
	"api",
	"oks",
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
			"endpoints": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"api": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The Endpoint for Outscale API operations.",
						},
						"oks": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The Endpoint for OKS API operations.",
						},
					},
				},
			},
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
			"config_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The path to your configuration file in which you have defined your credentials.",
			},
			"profile": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of your profile in which you define your credencial",
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
			"outscale_image":                             ResourceOutscaleImage(),
			"outscale_tag":                               ResourceOutscaleTags(),
			"outscale_public_ip":                         ResourceOutscalePublicIP(),
			"outscale_public_ip_link":                    ResourceOutscalePublicIPLink(),
			"outscale_nat_service":                       ResourceOutscaleNatService(),
			"outscale_nic":                               ResourceOutscaleNic(),
			"outscale_snapshot":                          ResourceOutscaleSnapshot(),
			"outscale_image_launch_permission":           ResourceOutscaleImageLaunchPermission(),
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
			"outscale_load_balancer":                     ResourceOutscaleLoadBalancer(),
			"outscale_load_balancer_policy":              ResourceOutscaleAppCookieStickinessPolicy(),
			"outscale_load_balancer_attributes":          ResourceOutscaleLoadBalancerAttributes(),
			"outscale_load_balancer_listener_rule":       ResourceOutscaleLoadBalancerListenerRule(),
			"outscale_flexible_gpu_link":                 ResourceOutscaleFlexibleGpuLink(),
			"outscale_image_export_task":                 ResourceOutscaleIMageExportTask(),
			"outscale_server_certificate":                ResourceOutscaleServerCertificate(),
			"outscale_snapshot_export_task":              ResourceOutscaleSnapshotExportTask(),
			"outscale_ca":                                ResourceOutscaleCa(),
			"outscale_api_access_rule":                   ResourceOutscaleApiAccessRule(),
			"outscale_api_access_policy":                 ResourceOutscaleApiAccessPolicy(),
			"outscale_user":                              ResourceOutscaleUser(),
			"outscale_user_group":                        ResourceUserGroup(),
			"outscale_policy":                            ResourceOutscalePolicy(),
			"outscale_policy_version":                    ResourcePolicyVersion(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"outscale_vm":                            DataSourceOutscaleVM(),
			"outscale_vms":                           DataSourceOutscaleVMS(),
			"outscale_security_group":                DataSourceOutscaleSecurityGroup(),
			"outscale_security_groups":               DataSourceOutscaleSecurityGroups(),
			"outscale_image":                         DataSourceOutscaleImage(),
			"outscale_images":                        DataSourceOutscaleImages(),
			"outscale_tag":                           DataSourceOutscaleTag(),
			"outscale_tags":                          DataSourceOutscaleTags(),
			"outscale_public_ip":                     DataSourceOutscalePublicIP(),
			"outscale_public_ips":                    DataSourceOutscalePublicIPS(),
			"outscale_volume":                        DataSourceOutscaleVolume(),
			"outscale_volumes":                       DataSourceOutscaleVolumes(),
			"outscale_nat_service":                   DataSourceOutscaleNatService(),
			"outscale_nat_services":                  DataSourceOutscaleNatServices(),
			"outscale_keypair":                       DataSourceOutscaleKeyPair(),
			"outscale_keypairs":                      DataSourceOutscaleKeyPairs(),
			"outscale_vm_state":                      DataSourceOutscaleVMState(),
			"outscale_vm_states":                     DataSourceOutscaleVMStates(),
			"outscale_internet_service":              DataSourceOutscaleInternetService(),
			"outscale_internet_services":             DataSourceOutscaleInternetServices(),
			"outscale_subnet":                        DataSourceOutscaleSubnet(),
			"outscale_subnets":                       DataSourceOutscaleSubnets(),
			"outscale_net":                           DataSourceOutscaleVpc(),
			"outscale_nets":                          DataSourceOutscaleVpcs(),
			"outscale_net_attributes":                DataSourceOutscaleVpcAttr(),
			"outscale_route_table":                   DataSourceOutscaleRouteTable(),
			"outscale_route_tables":                  DataSourceOutscaleRouteTables(),
			"outscale_snapshot":                      DataSourceOutscaleSnapshot(),
			"outscale_snapshots":                     DataSourceOutscaleSnapshots(),
			"outscale_net_peering":                   DataSourceOutscaleLinPeeringConnection(),
			"outscale_net_peerings":                  DataSourceOutscaleLinPeeringsConnection(),
			"outscale_nics":                          DataSourceOutscaleNics(),
			"outscale_nic":                           DataSourceOutscaleNic(),
			"outscale_client_gateway":                DataSourceOutscaleClientGateway(),
			"outscale_client_gateways":               DataSourceOutscaleClientGateways(),
			"outscale_virtual_gateway":               DataSourceOutscaleVirtualGateway(),
			"outscale_virtual_gateways":              DataSourceOutscaleVirtualGateways(),
			"outscale_vpn_connection":                DataSourceOutscaleVPNConnection(),
			"outscale_vpn_connections":               DataSourceOutscaleVPNConnections(),
			"outscale_access_key":                    DataSourceOutscaleAccessKey(),
			"outscale_access_keys":                   DataSourceOutscaleAccessKeys(),
			"outscale_dhcp_option":                   DataSourceOutscaleDHCPOption(),
			"outscale_dhcp_options":                  DataSourceOutscaleDHCPOptions(),
			"outscale_load_balancer":                 DataSourceOutscaleLoadBalancer(),
			"outscale_load_balancer_listener_rule":   DataSourceOutscaleLoadBalancerLDRule(),
			"outscale_load_balancer_listener_rules":  DataSourceOutscaleLoadBalancerLDRules(),
			"outscale_load_balancer_tags":            DataSourceOutscaleLBUTags(),
			"outscale_load_balancer_vm_health":       DataSourceOutscaleLoadBalancerVmsHeals(),
			"outscale_load_balancers":                DataSourceOutscaleLoadBalancers(),
			"outscale_vm_types":                      DataSourceOutscaleVMTypes(),
			"outscale_net_access_point":              DataSourceOutscaleNetAccessPoint(),
			"outscale_net_access_points":             DataSourceOutscaleNetAccessPoints(),
			"outscale_flexible_gpu":                  DataSourceOutscaleFlexibleGpu(),
			"outscale_flexible_gpus":                 DataSourceOutscaleFlexibleGpus(),
			"outscale_subregions":                    DataSourceOutscaleSubregions(),
			"outscale_regions":                       DataSourceOutscaleRegions(),
			"outscale_net_access_point_services":     DataSourceOutscaleNetAccessPointServices(),
			"outscale_flexible_gpu_catalog":          DataSourceOutscaleFlexibleGpuCatalog(),
			"outscale_product_type":                  DataSourceOutscaleProductType(),
			"outscale_product_types":                 DataSourceOutscaleProductTypes(),
			"outscale_quotas":                        DataSourceOutscaleQuotas(),
			"outscale_image_export_task":             DataSourceOutscaleImageExportTask(),
			"outscale_image_export_tasks":            DataSourceOutscaleImageExportTasks(),
			"outscale_server_certificate":            DataSourceOutscaleServerCertificate(),
			"outscale_server_certificates":           DataSourceOutscaleServerCertificates(),
			"outscale_snapshot_export_task":          DataSourceOutscaleSnapshotExportTask(),
			"outscale_snapshot_export_tasks":         DataSourceOutscaleSnapshotExportTasks(),
			"outscale_ca":                            DataSourceOutscaleCa(),
			"outscale_cas":                           DataSourceOutscaleCas(),
			"outscale_api_access_rule":               DataSourceOutscaleApiAccessRule(),
			"outscale_api_access_rules":              DataSourceOutscaleApiAccessRules(),
			"outscale_api_access_policy":             DataSourceOutscaleApiAccessPolicy(),
			"outscale_public_catalog":                DataSourceOutscalePublicCatalog(),
			"outscale_account":                       DataSourceAccount(),
			"outscale_accounts":                      DataSourceAccounts(),
			"outscale_users":                         DataSourceUsers(),
			"outscale_user":                          DataSourceUser(),
			"outscale_user_groups":                   DataSourceUserGroups(),
			"outscale_user_groups_per_user":          DataSourceUserGroupsPerUser(),
			"outscale_user_group":                    DataSourceUserGroup(),
			"outscale_policy":                        DataSourcePolicy(),
			"outscale_policies":                      DataSourcePolicies(),
			"outscale_policies_linked_to_user":       DataSourcePoliciesLinkedToUser(),
			"outscale_entities_linked_to_policy":     DataSourceEntitiesLinkedToPolicy(),
			"outscale_policies_linked_to_user_group": DataSourcePoliciesLinkedToUserGroup(),
		},

		ConfigureFunc: providerConfigureClient,
	}
}

func providerConfigureClient(d *schema.ResourceData) (interface{}, error) {
	config := Config{
		AccessKeyID:  d.Get("access_key_id").(string),
		SecretKeyID:  d.Get("secret_key_id").(string),
		Region:       d.Get("region").(string),
		Endpoints:    make(map[string]string),
		X509CertPath: d.Get("x509_cert_path").(string),
		X509KeyPath:  d.Get("x509_key_path").(string),
		ConfigFile:   d.Get("config_file").(string),
		Profile:      d.Get("profile").(string),
		Insecure:     d.Get("insecure").(bool),
	}
	endpointsSet := d.Get("endpoints").(*schema.Set)
	for _, endpointsSetI := range endpointsSet.List() {
		endpoints := make(map[string]string)
		for key, value := range endpointsSetI.(map[string]interface{}) {
			endpoints[key] = value.(string)
		}
		for _, endpointServiceName := range endpointServiceNames {
			config.Endpoints[endpointServiceName] = endpoints[endpointServiceName]
		}
	}

	ok, err := IsOldProfileSet(&config)
	if err != nil {
		return nil, err
	}
	if !ok {
		setProviderDefaultEnv(&config)
	}
	return config.Client()
}
func IsOldProfileSet(conf *Config) (bool, error) {
	isProfSet := false
	if profileName, ok := os.LookupEnv("OSC_PROFILE"); ok || conf.Profile != "" {
		if conf.Profile != "" {
			profileName = conf.Profile
		}

		var configFilePath string
		if envPath, ok := os.LookupEnv("OSC_CONFIG_FILE"); ok || conf.ConfigFile != "" {
			if conf.ConfigFile != "" {
				configFilePath = conf.ConfigFile
			} else {
				configFilePath = envPath
			}
		} else {
			homePath, err := os.UserHomeDir()
			if err != nil {
				return isProfSet, err
			}
			configFilePath = homePath + utils.SuffixConfigFilePath
		}
		jsonFile, err := os.ReadFile(configFilePath)
		if err != nil {
			return isProfSet, fmt.Errorf("unable to read config file '%v', Error: %w", configFilePath, err)
		}
		profile := gjson.GetBytes(jsonFile, profileName)
		if !gjson.Valid(profile.String()) {
			return isProfSet, fmt.Errorf("invalid json profile file")
		}
		if !profile.Get("access_key").Exists() ||
			!profile.Get("secret_key").Exists() {
			return isProfSet, errors.New("profile 'access_key' or 'secret_key' are not defined! ")
		}
		setOldProfile(conf, profile)
		isProfSet = true
	}
	return isProfSet, nil
}

func setOldProfile(conf *Config, profile gjson.Result) {

	if conf.AccessKeyID == "" {
		if accessKeyId := profile.Get("access_key").String(); accessKeyId != "" {
			conf.AccessKeyID = accessKeyId
		}
	}
	if conf.SecretKeyID == "" {
		if secretKeyId := profile.Get("secret_key").String(); secretKeyId != "" {
			conf.SecretKeyID = secretKeyId
		}
	}
	if conf.Region == "" {
		if profile.Get("region").Exists() {
			if region := profile.Get("region").String(); region != "" {
				conf.Region = region
			}
		}
	}
	if conf.X509CertPath == "" {
		if profile.Get("x509_cert_path").Exists() {
			if x509Cert := profile.Get("x509_cert_path").String(); x509Cert != "" {
				conf.X509CertPath = x509Cert
			}
		}
	}
	if conf.X509KeyPath == "" {
		if profile.Get("x509_key_path").Exists() {
			if x509Key := profile.Get("x509_key_path").String(); x509Key != "" {
				conf.X509KeyPath = x509Key
			}
		}
	}
	if len(conf.Endpoints) == 0 {
		if profile.Get("endpoints").Exists() {
			endpoints := profile.Get("endpoints").Value().(map[string]interface{})
			if endpoint := endpoints["api"].(string); endpoint != "" {
				conf.Endpoints["api"] = endpoint
			}
		}
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

	if conf.X509CertPath == "" {
		if x509Cert := utils.GetEnvVariableValue([]string{"OSC_X509_CLIENT_CERT", "OUTSCALE_X509CERT"}); x509Cert != "" {
			conf.X509CertPath = x509Cert
		}
	}

	if conf.X509KeyPath == "" {
		if x509Key := utils.GetEnvVariableValue([]string{"OSC_X509_CLIENT_KEY", "OUTSCALE_X509KEY"}); x509Key != "" {
			conf.X509KeyPath = x509Key
		}
	}
	if len(conf.Endpoints) == 0 {
		if endpoints := utils.GetEnvVariableValue([]string{"OSC_ENDPOINT_API", "OUTSCALE_OAPI_URL"}); endpoints != "" {
			endpointsAttributes := make(map[string]string)
			endpointsAttributes["api"] = endpoints
			conf.Endpoints = endpointsAttributes
		}
	}
}
