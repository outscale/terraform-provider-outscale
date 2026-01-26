package provider

import (
	"errors"
	"fmt"
	"os"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/outscale/terraform-provider-outscale/internal/services/oapi"
	"github.com/outscale/terraform-provider-outscale/internal/utils"
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
			"outscale_vm":                                oapi.ResourceOutscaleVM(),
			"outscale_image":                             oapi.ResourceOutscaleImage(),
			"outscale_tag":                               oapi.ResourceOutscaleTags(),
			"outscale_public_ip":                         oapi.ResourceOutscalePublicIP(),
			"outscale_public_ip_link":                    oapi.ResourceOutscalePublicIPLink(),
			"outscale_nat_service":                       oapi.ResourceOutscaleNatService(),
			"outscale_nic":                               oapi.ResourceOutscaleNic(),
			"outscale_snapshot":                          oapi.ResourceOutscaleSnapshot(),
			"outscale_image_launch_permission":           oapi.ResourceOutscaleImageLaunchPermission(),
			"outscale_nic_link":                          oapi.ResourceOutscaleNetworkInterfaceAttachment(),
			"outscale_nic_private_ip":                    oapi.ResourceOutscaleNetworkInterfacePrivateIP(),
			"outscale_snapshot_attributes":               oapi.ResourceOutscaleSnapshotAttributes(),
			"outscale_dhcp_option":                       oapi.ResourceOutscaleDHCPOption(),
			"outscale_client_gateway":                    oapi.ResourceOutscaleClientGateway(),
			"outscale_virtual_gateway":                   oapi.ResourceOutscaleVirtualGateway(),
			"outscale_virtual_gateway_link":              oapi.ResourceOutscaleVirtualGatewayLink(),
			"outscale_virtual_gateway_route_propagation": oapi.ResourceOutscaleVirtualGatewayRoutePropagation(),
			"outscale_vpn_connection":                    oapi.ResourceOutscaleVPNConnection(),
			"outscale_vpn_connection_route":              oapi.ResourceOutscaleVPNConnectionRoute(),
			"outscale_load_balancer":                     oapi.ResourceOutscaleLoadBalancer(),
			"outscale_load_balancer_policy":              oapi.ResourceOutscaleAppCookieStickinessPolicy(),
			"outscale_load_balancer_attributes":          oapi.ResourceOutscaleLoadBalancerAttributes(),
			"outscale_load_balancer_listener_rule":       oapi.ResourceOutscaleLoadBalancerListenerRule(),
			"outscale_flexible_gpu_link":                 oapi.ResourceOutscaleFlexibleGpuLink(),
			"outscale_image_export_task":                 oapi.ResourceOutscaleImageExportTask(),
			"outscale_server_certificate":                oapi.ResourceOutscaleServerCertificate(),
			"outscale_snapshot_export_task":              oapi.ResourceOutscaleSnapshotExportTask(),
			"outscale_api_access_rule":                   oapi.ResourceOutscaleApiAccessRule(),
			"outscale_api_access_policy":                 oapi.ResourceOutscaleApiAccessPolicy(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"outscale_vm":                            oapi.DataSourceOutscaleVM(),
			"outscale_vms":                           oapi.DataSourceOutscaleVMS(),
			"outscale_security_group":                oapi.DataSourceOutscaleSecurityGroup(),
			"outscale_security_groups":               oapi.DataSourceOutscaleSecurityGroups(),
			"outscale_image":                         oapi.DataSourceOutscaleImage(),
			"outscale_images":                        oapi.DataSourceOutscaleImages(),
			"outscale_tag":                           oapi.DataSourceOutscaleTag(),
			"outscale_tags":                          oapi.DataSourceOutscaleTags(),
			"outscale_public_ip":                     oapi.DataSourceOutscalePublicIP(),
			"outscale_public_ips":                    oapi.DataSourceOutscalePublicIPS(),
			"outscale_volume":                        oapi.DataSourceOutscaleVolume(),
			"outscale_volumes":                       oapi.DataSourceOutscaleVolumes(),
			"outscale_nat_service":                   oapi.DataSourceOutscaleNatService(),
			"outscale_nat_services":                  oapi.DataSourceOutscaleNatServices(),
			"outscale_keypair":                       oapi.DataSourceOutscaleKeyPair(),
			"outscale_keypairs":                      oapi.DataSourceOutscaleKeyPairs(),
			"outscale_vm_state":                      oapi.DataSourceOutscaleVMState(),
			"outscale_vm_states":                     oapi.DataSourceOutscaleVMStates(),
			"outscale_internet_service":              oapi.DataSourceOutscaleInternetService(),
			"outscale_internet_services":             oapi.DataSourceOutscaleInternetServices(),
			"outscale_subnet":                        oapi.DataSourceOutscaleSubnet(),
			"outscale_subnets":                       oapi.DataSourceOutscaleSubnets(),
			"outscale_net":                           oapi.DataSourceOutscaleVpc(),
			"outscale_nets":                          oapi.DataSourceOutscaleVpcs(),
			"outscale_net_attributes":                oapi.DataSourceOutscaleVpcAttr(),
			"outscale_route_table":                   oapi.DataSourceOutscaleRouteTable(),
			"outscale_route_tables":                  oapi.DataSourceOutscaleRouteTables(),
			"outscale_snapshot":                      oapi.DataSourceOutscaleSnapshot(),
			"outscale_snapshots":                     oapi.DataSourceOutscaleSnapshots(),
			"outscale_net_peering":                   oapi.DataSourceOutscaleLinPeeringConnection(),
			"outscale_net_peerings":                  oapi.DataSourceOutscaleLinPeeringsConnection(),
			"outscale_nics":                          oapi.DataSourceOutscaleNics(),
			"outscale_nic":                           oapi.DataSourceOutscaleNic(),
			"outscale_client_gateway":                oapi.DataSourceOutscaleClientGateway(),
			"outscale_client_gateways":               oapi.DataSourceOutscaleClientGateways(),
			"outscale_virtual_gateway":               oapi.DataSourceOutscaleVirtualGateway(),
			"outscale_virtual_gateways":              oapi.DataSourceOutscaleVirtualGateways(),
			"outscale_vpn_connection":                oapi.DataSourceOutscaleVPNConnection(),
			"outscale_vpn_connections":               oapi.DataSourceOutscaleVPNConnections(),
			"outscale_access_key":                    oapi.DataSourceOutscaleAccessKey(),
			"outscale_access_keys":                   oapi.DataSourceOutscaleAccessKeys(),
			"outscale_dhcp_option":                   oapi.DataSourceOutscaleDHCPOption(),
			"outscale_dhcp_options":                  oapi.DataSourceOutscaleDHCPOptions(),
			"outscale_load_balancer":                 oapi.DataSourceOutscaleLoadBalancer(),
			"outscale_load_balancer_listener_rule":   oapi.DataSourceOutscaleLoadBalancerLDRule(),
			"outscale_load_balancer_listener_rules":  oapi.DataSourceOutscaleLoadBalancerLDRules(),
			"outscale_load_balancer_tags":            oapi.DataSourceOutscaleLBUTags(),
			"outscale_load_balancer_vm_health":       oapi.DataSourceOutscaleLoadBalancerVmsHeals(),
			"outscale_load_balancers":                oapi.DataSourceOutscaleLoadBalancers(),
			"outscale_vm_types":                      oapi.DataSourceOutscaleVMTypes(),
			"outscale_net_access_point":              oapi.DataSourceOutscaleNetAccessPoint(),
			"outscale_net_access_points":             oapi.DataSourceOutscaleNetAccessPoints(),
			"outscale_flexible_gpu":                  oapi.DataSourceOutscaleFlexibleGpu(),
			"outscale_flexible_gpus":                 oapi.DataSourceOutscaleFlexibleGpus(),
			"outscale_subregions":                    oapi.DataSourceOutscaleSubregions(),
			"outscale_regions":                       oapi.DataSourceOutscaleRegions(),
			"outscale_net_access_point_services":     oapi.DataSourceOutscaleNetAccessPointServices(),
			"outscale_flexible_gpu_catalog":          oapi.DataSourceOutscaleFlexibleGpuCatalog(),
			"outscale_product_type":                  oapi.DataSourceOutscaleProductType(),
			"outscale_product_types":                 oapi.DataSourceOutscaleProductTypes(),
			"outscale_quotas":                        oapi.DataSourceOutscaleQuotas(),
			"outscale_image_export_task":             oapi.DataSourceOutscaleImageExportTask(),
			"outscale_image_export_tasks":            oapi.DataSourceOutscaleImageExportTasks(),
			"outscale_server_certificate":            oapi.DataSourceOutscaleServerCertificate(),
			"outscale_server_certificates":           oapi.DataSourceOutscaleServerCertificates(),
			"outscale_snapshot_export_task":          oapi.DataSourceOutscaleSnapshotExportTask(),
			"outscale_snapshot_export_tasks":         oapi.DataSourceOutscaleSnapshotExportTasks(),
			"outscale_ca":                            oapi.DataSourceOutscaleCa(),
			"outscale_cas":                           oapi.DataSourceOutscaleCas(),
			"outscale_api_access_rule":               oapi.DataSourceOutscaleApiAccessRule(),
			"outscale_api_access_rules":              oapi.DataSourceOutscaleApiAccessRules(),
			"outscale_api_access_policy":             oapi.DataSourceOutscaleApiAccessPolicy(),
			"outscale_public_catalog":                oapi.DataSourceOutscalePublicCatalog(),
			"outscale_account":                       oapi.DataSourceAccount(),
			"outscale_accounts":                      oapi.DataSourceAccounts(),
			"outscale_users":                         oapi.DataSourceUsers(),
			"outscale_user":                          oapi.DataSourceUser(),
			"outscale_user_groups":                   oapi.DataSourceUserGroups(),
			"outscale_user_groups_per_user":          oapi.DataSourceUserGroupsPerUser(),
			"outscale_user_group":                    oapi.DataSourceUserGroup(),
			"outscale_policy":                        oapi.DataSourcePolicy(),
			"outscale_policies":                      oapi.DataSourcePolicies(),
			"outscale_policies_linked_to_user":       oapi.DataSourcePoliciesLinkedToUser(),
			"outscale_entities_linked_to_policy":     oapi.DataSourceEntitiesLinkedToPolicy(),
			"outscale_policies_linked_to_user_group": oapi.DataSourcePoliciesLinkedToUserGroup(),
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
