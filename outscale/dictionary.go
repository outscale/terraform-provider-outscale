package outscale

import "github.com/hashicorp/terraform/helper/schema"

//Dictionary for the Outscale APIs maps the apis to their respective functions
type Dictionary map[string]ResourceMap

//ResourceMap maps a schema to their resource or datasource implementation
type ResourceMap map[string]SchemaFunc

//SchemaFunc maps a function that returns a schema
type SchemaFunc func() *schema.Resource

var resources Dictionary
var datasources Dictionary

func init() {
	resources = Dictionary{
		"fcu": ResourceMap{
			"outscale_vm":                            resourceOutscaleVM,
			"outscale_image":                         resourceOutscaleImage,
			"outscale_firewall_rules_set":            resourceOutscaleFirewallRulesSet,
			"outscale_outbound_rule":                 resourceOutscaleOutboundRule,
			"outscale_inbound_rule":                  resourceOutscaleInboundRule,
			"outscale_tag":                           resourceOutscaleTags,
			"outscale_keypair":                       resourceOutscaleKeyPair,
			"outscale_public_ip":                     resourceOutscalePublicIP,
			"outscale_public_ip_link":                resourceOutscalePublicIPLink,
			"outscale_volume":                        resourceOutscaleVolume,
			"outscale_volumes_link":                  resourceOutscaleVolumeLink,
			"outscale_lin":                           resourceOutscaleLin,
			"outscale_lin_attributes":                resourceOutscaleLinAttributes,
			"outscale_lin_internet_gateway":          resourceOutscaleLinInternetGateway,
			"outscale_lin_internet_gateway_link":     resourceOutscaleLinInternetGatewayLink,
			"outscale_vm_attributes":                 resourceOutscaleVMAttributes,
			"outscale_nat_service":                   resourceOutscaleNatService,
			"outscale_subnet":                        resourceOutscaleSubNet,
			"outscale_keypair_importation":           resourceOutscaleKeyPairImportation,
			"outscale_client_endpoint":               resourceOutscaleCustomerGateway,
			"outscale_route":                         resourceOutscaleRoute,
			"outscale_route_table":                   resourceOutscaleRouteTable,
			"outscale_route_table_link":              resourceOutscaleRouteTableAssociation,
			"outscale_dhcp_option":                   resourceOutscaleDHCPOption,
			"outscale_dhcp_option_link":              resourceOutscaleDHCPOptionLink,
			"outscale_image_copy":                    resourceOutscaleImageCopy,
			"outscale_vpn_connection":                resourceOutscaleVpnConnection,
			"outscale_vpn_gateway":                   resourceOutscaleVpnGateway,
			"outscale_image_tasks":                   resourceOutscaleImageTasks,
			"outscale_vpn_connection_route":          resourceOutscaleVpnConnectionRoute,
			"outscale_vpn_gateway_route_propagation": resourceOutscaleVpnGatewayRoutePropagation,
			"outscale_vpn_gateway_link":              resourceOutscaleVpnGatewayLink,
			"outscale_nic":                           resourceOutscaleNic,
			"outscale_snapshot_export_task":          resourceOutscaleImageExportTasks,
			"outscale_snapshot":                      resourceOutscaleSnapshot,
			"outscale_image_register":                resourceOutscaleImageRegister,
			"outscale_image_launch_permission":       resourceOutscaleImageLaunchPermission,
			"outscale_lin_peering":                   resourceOutscaleLinPeeringConnection,
			"outscale_lin_peering_acceptation":       resourceOutscaleLinPeeringConnectionAccepter,
			"outscale_nic_link":                      resourceOutscaleNetworkInterfaceAttachment,
			"outscale_nic_private_ip":                resourceOutscaleNetworkInterfacePrivateIP,
			"outscale_reserved_vms_offer_purchase":   resourceOutscaleReservedVmsOfferPurchase,
			"outscale_lin_api_access":                resourceOutscaleVpcEndpoint,
		},
		"oapi": ResourceMap{
			"outscale_vm":                            resourceOutscaleOApiVM,
			"outscale_firewall_rules_set":            resourceOutscaleOAPIFirewallRulesSet,
			"outscale_image":                         resourceOutscaleOAPIImage,
			"outscale_keypair":                       resourceOutscaleOAPIKeyPair,
			"outscale_public_ip":                     resourceOutscaleOAPIPublicIP,
			"outscale_inbound_rule":                  resourceOutscaleOAPIInboundRule,
			"outscale_outbound_rule":                 resourceOutscaleOAPIOutboundRule,
			"outscale_tag":                           resourceOutscaleOAPITags,
			"outscale_lin_attributes":                resourceOutscaleOAPILinAttributes,
			"outscale_lin":                           resourceOutscaleOAPILin,
			"outscale_lin_internet_gateway":          resourceOutscaleOAPILinInternetGateway,
			"outscale_nat_service":                   resourceOutscaleOAPINatService,
			"outscale_subnet":                        resourceOutscaleOAPISubNet,
			"outscale_client_endpoint":               resourceOutscaleCustomerGateway,
			"outscale_route":                         resourceOutscaleRoute,
			"outscale_route_table":                   resourceOutscaleRouteTable,
			"outscale_route_table_link":              resourceOutscaleRouteTableAssociation,
			"outscale_snapshot":                      resourceOutscaleOAPISnapshot,
			"outscale_api_key":                       resourceOutscaleOAPIIamAccessKey,
			"outscale_keypair_importation":           resourceOutscaleOAPIKeyPairImportation,
			"outscale_image_launch_permission":       resourceOutscaleOAPIImageLaunchPermission,
			"outscale_lin_peering":                   resourceOutscaleOAPILinPeeringConnection,
			"outscale_load_balancer":                 resourceOutscaleOAPILoadBalancer,
			"outscale_nic_private_ip":                resourceOutscaleOAPINetworkInterfacePrivateIP,
			"outscale_nic_link":                      resourceOutscaleOAPINetworkInterfaceAttachment,
			"outscale_nic":                           resourceOutscaleOAPINic,
			"outscale_load_balancer_cookiepolicy":    resourceOutscaleAppCookieStickinessPolicy,
			"outscale_load_balancer_listeners":       resourceOutscaleOAPILoadBalancerListeners,
			"outscale_load_balancer_attributes":      resourceOutscaleOAPILoadBalancerAttributes,
			"outscale_load_balancer_vms":             resourceOutscaleOAPILBUAttachment,
			"outscale_load_balancer_tags":            resourceOutscaleOAPILBUTags,
			"outscale_image_tasks":                   resourceOutscaleOAPIImageTasks,
			"outscale_reserved_vms_offer_purchase":   resourceOutscaleOAPIReservedVmsOfferPurchase,
			"outscale_vpn_gateway":                   resourceOutscaleOAPIVpnGateway,
			"outscale_vpn_gateway_route_propagation": resourceOutscaleOAPIVpnGatewayRoutePropagation,
			"outscale_snapshot_export_task":          resourceOutscaleOAPIImageExportTasks,
			"outscale_vpn_connection_route":          resourceOutscaleOAPIVpnConnectionRoute,
			"outscale_vpn_connection":                resourceOutscaleOAPIVpnConnection,
			"outscale_image_copy":                    resourceOutscaleOAPIImageCopy,
			"outscale_vpn_gateway_link":              resourceOutscaleOAPIVpnGatewayLink,
		},
		"icu": ResourceMap{
			"outscale_api_key": resourceOutscaleIamAccessKey,
		},
		"lbu": ResourceMap{
			"outscale_load_balancer":              resourceOutscaleLoadBalancer,
			"outscale_load_balancer_cookiepolicy": resourceOutscaleAppCookieStickinessPolicy,
			"outscale_load_balancer_vms":          resourceOutscaleLBUAttachment,
			"outscale_load_balancer_listeners":    resourceOutscaleLoadBalancerListeners,
			"outscale_load_balancer_attributes":   resourceOutscaleLoadBalancerAttributes,
			"outscale_load_balancer_tags":         resourceOutscaleLBUTags,
		},
	}
	datasources = Dictionary{
		"fcu": ResourceMap{
			"outscale_vm":                    dataSourceOutscaleVM,
			"outscale_vms":                   dataSourceOutscaleVMS,
			"outscale_firewall_rules_set":    dataSourceOutscaleFirewallRuleSet,
			"outscale_firewall_rules_sets":   dataSourceOutscaleFirewallRulesSets,
			"outscale_image":                 dataSourceOutscaleImage,
			"outscale_images":                dataSourceOutscaleImages,
			"outscale_tag":                   dataSourceOutscaleTag,
			"outscale_tags":                  dataSourceOutscaleTags,
			"outscale_public_ip":             dataSourceOutscalePublicIP,
			"outscale_public_ips":            dataSourceOutscalePublicIPS,
			"outscale_volume":                datasourceOutscaleVolume,
			"outscale_volumes":               datasourceOutscaleVolumes,
			"outscale_nat_service":           dataSourceOutscaleNatService,
			"outscale_nat_services":          dataSourceOutscaleNatServices,
			"outscale_keypair":               datasourceOutscaleKeyPair,
			"outscale_keypairs":              datasourceOutscaleKeyPairs,
			"outscale_vm_state":              dataSourceOutscaleVMState,
			"outscale_vms_state":             dataSourceOutscaleVMSState,
			"outscale_lin_internet_gateway":  datasourceOutscaleLinInternetGateway,
			"outscale_lin_internet_gateways": datasourceOutscaleLinInternetGateways,
			"outscale_subnet":                dataSourceOutscaleSubnet,
			"outscale_subnets":               dataSourceOutscaleSubnets,
			"outscale_lin":                   dataSourceOutscaleVpc,
			"outscale_lins":                  dataSourceOutscaleVpcs,
			"outscale_lin_attributes":        dataSourceOutscaleVpcAttr,
			"outscale_client_endpoint":       dataSourceOutscaleCustomerGateway,
			"outscale_client_endpoints":      dataSourceOutscaleCustomerGateways,
			"outscale_route_table":           dataSourceOutscaleRouteTable,
			"outscale_route_tables":          dataSourceOutscaleRouteTables,
			"outscale_vpn_gateway":           dataSourceOutscaleVpnGateway,
			"outscale_api_key":               dataSourceOutscaleIamAccessKey,
			"outscale_vpn_gateways":          dataSourceOutscaleVpnGateways,
			"outscale_vpn_connection":        dataSourceOutscaleVpnConnection,
			"outscale_sub_region":            dataSourceOutscaleAvailabilityZone,
			"outscale_prefix_list":           dataSourceOutscalePrefixList,
			"outscale_quota":                 dataSourceOutscaleQuota,
			"outscale_quotas":                dataSourceOutscaleQuotas,
			"outscale_prefix_lists":          dataSourceOutscalePrefixLists,
			"outscale_region":                dataSourceOutscaleRegion,
			"outscale_sub_regions":           dataSourceOutscaleAvailabilityZones,
			"outscale_regions":               dataSourceOutscaleRegions,
			"outscale_vpn_connections":       dataSourceOutscaleVpnConnections,
			"outscale_product_types":         dataSourceOutscaleProductTypes,
			"outscale_reserved_vms":          dataSourceOutscaleReservedVMS,
			"outscale_vm_type":               dataSourceOutscaleVMType,
			"outscale_vm_types":              dataSourceOutscaleVMTypes,
			"outscale_reserved_vms_offers":   dataSourceOutscaleReservedVMOffers,
			"outscale_reserved_vms_offer":    dataSourceOutscaleReservedVMOffer,
			"outscale_snapshot":              dataSourceOutscaleSnapshot,
			"outscale_snapshots":             dataSourceOutscaleSnapshots,
			"outscale_lin_peering":           dataSourceOutscaleLinPeeringConnection,
			"outscale_lin_peerings":          dataSourceOutscaleLinPeeringsConnection,
			"outscale_nics":                  dataSourceOutscaleNics,
			"outscale_nic":                   dataSourceOutscaleNic,
		},
		"oapi": ResourceMap{
			"outscale_vm":                                  dataSourceOutscaleOAPIVM,
			"outscale_vms":                                 datasourceOutscaleOApiVMS,
			"outscale_firewall_rules_sets":                 dataSourceOutscaleOAPIFirewallRulesSets,
			"outscale_images":                              dataSourceOutscaleOAPIImages,
			"outscale_firewall_rules_set":                  dataSourceOutscaleOAPIFirewallRuleSet,
			"outscale_tag":                                 dataSourceOutscaleOAPITag,
			"outscale_tags":                                dataSourceOutscaleOAPITags,
			"outscale_volume":                              datasourceOutscaleOAPIVolume,
			"outscale_volumes":                             datasourceOutscaleOAPIVolumes,
			"outscale_keypair":                             datasourceOutscaleOAPIKeyPair,
			"outscale_keypairs":                            datasourceOutscaleOAPIKeyPairs,
			"outscale_lin_internet_gateway":                datasourceOutscaleOAPILinInternetGateway,
			"outscale_lin_internet_gateways":               datasourceOutscaleOAPILinInternetGateways,
			"outscale_subnet":                              dataSourceOutscaleOAPISubnet,
			"outscale_subnets":                             dataSourceOutscaleOAPISubnets,
			"outscale_vm_state":                            dataSourceOutscaleOAPIVMState,
			"outscale_vms_state":                           dataSourceOutscaleOAPIVMSState,
			"outscale_lin":                                 dataSourceOutscaleOAPIVpc,
			"outscale_lins":                                dataSourceOutscaleOAPIVpcs,
			"outscale_lin_attributes":                      dataSourceOutscaleOAPIVpcAttr,
			"outscale_client_endpoint":                     dataSourceOutscaleOAPICustomerGateway,
			"outscale_client_endpoints":                    dataSourceOutscaleOAPICustomerGateways,
			"outscale_route_table":                         dataSourceOutscaleOAPIRouteTable,
			"outscale_route_tables":                        dataSourceOutscaleOAPIRouteTables,
			"outscale_snapshot":                            dataSourceOutscaleOAPISnapshot,
			"outscale_snapshots":                           dataSourceOutscaleOAPISnapshots,
			"outscale_lin_peering":                         dataSourceOutscaleOAPILinPeeringConnection,
			"outscale_lin_peerings":                        dataSourceOutscaleOAPILinPeeringsConnection,
			"outscale_load_balancer":                       dataSourceOutscaleOAPILoadBalancer,
			"outscale_load_balancers":                      dataSourceOutscaleOAPILoadBalancers,
			"outscale_nic":                                 dataSourceOutscaleOAPINic,
			"outscale_nics":                                dataSourceOutscaleOAPINics,
			"outscale_load_balancer_tags":                  dataSourceOutscaleOAPILBUTags,
			"outscale_vm_health":                           dataSourceOutscaleOAPIVMHealth,
			"outscale_load_balancer_health_check":          dataSourceOutscaleOAPILoadBalancerHealthCheck,
			"outscale_load_balancer_listener_description":  dataSourceOutscaleOAPILoadBalancerLD,
			"outscale_load_balancer_listener_descriptions": dataSourceOutscaleOAPILoadBalancerLDs,
			"outscale_load_balancer_attributes":            dataSourceOutscaleOAPILoadBalancerAttr,
			"outscale_regions":                             dataSourceOutscaleOAPIRegions,
			"outscale_region":                              dataSourceOutscaleOAPIRegion,
			"outscale_reserved_vms_offer":                  dataSourceOutscaleOAPIReservedVMOffer,
			"outscale_reserved_vms_offers":                 dataSourceOutscaleOAPIReservedVMOffers,
			"outscale_reserved_vms":                        dataSourceOutscaleOAPIReservedVMS,
			"outscale_vpn_gateways":                        dataSourceOutscaleOAPIVpnGateways,
			"outscale_vpn_gateway":                         dataSourceOutscaleOAPIVpnGateway,
			"outscale_vm_type":                             dataSourceOutscaleOAPIVMType,
			"outscale_vm_types":                            dataSourceOutscaleOAPIVMTypes,
			"outscale_quotas":                              dataSourceOutscaleOAPIQuotas,
			"outscale_quota":                               dataSourceOutscaleOAPIQuota,
			"outscale_prefix_lists":                        dataSourceOutscaleOAPIPrefixLists,
			"outscale_prefix_list":                         dataSourceOutscaleOAPIPrefixList,
			"outscale_vpn_connections":                     dataSourceOutscaleOAPIVpnConnections,
			"outscale_sub_region":                          dataSourceOutscaleOAPIAvailabilityZone,
			"outscale_product_types":                       dataSourceOutscaleOAPIProductTypes,
			"outscale_image":                               dataSourceOutscaleOAPIImage,
		},
		"lbu": ResourceMap{
			"outscale_load_balancer":                       dataSourceOutscaleLoadBalancer,
			"outscale_load_balancers":                      dataSourceOutscaleLoadBalancers,
			"outscale_load_balancer_access_logs":           dataSourceOutscaleLoadBalancerAccessLogs,
			"outscale_load_balancer_tags":                  dataSourceOutscaleLBUTags,
			"outscale_vm_health":                           dataSourceOutscaleVMHealth,
			"outscale_load_balancer_health_check":          dataSourceOutscaleLoadBalancerHealthCheck,
			"outscale_load_balancer_listener_description":  dataSourceOutscaleLoadBalancerLD,
			"outscale_load_balancer_listener_descriptions": dataSourceOutscaleLoadBalancerLDs,
			"outscale_load_balancer_attributes":            dataSourceOutscaleLoadBalancerAttr,
		},
	}
}

//GetResource receives the apu and the name of the resource
//and returns the corrresponding
func GetResource(api, resource string) SchemaFunc {
	var a ResourceMap

	if _, ok := resources[api]; !ok {
		return nil
	}

	a = resources[api]

	if _, ok := a[resource]; !ok {
		return nil
	}
	return a[resource]
}

//GetDatasource receives the apu and the name of the datasource
//and returns the corrresponding
func GetDatasource(api, datasource string) SchemaFunc {
	var a ResourceMap
	if _, ok := datasources[api]; !ok {
		return nil
	}

	a = datasources[api]

	if _, ok := a[datasource]; !ok {
		return nil
	}
	return a[datasource]
}
