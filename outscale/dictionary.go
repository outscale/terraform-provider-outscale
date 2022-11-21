package outscale

import "github.com/hashicorp/terraform-plugin-sdk/helper/schema"

// Dictionary for the Outscale APIs maps the apis to their respective functions
type Dictionary map[string]ResourceMap

// ResourceMap maps a schema to their resource or datasource implementation
type ResourceMap map[string]SchemaFunc

// SchemaFunc maps a function that returns a schema
type SchemaFunc func() *schema.Resource

var resources Dictionary
var datasources Dictionary

func init() {
	resources = Dictionary{
		"oapi": ResourceMap{
			"outscale_vm":                      resourceVM,
			"outscale_firewall_rules_set":      resourceSecurityGroup,
			"outscale_security_group":          resourceSecurityGroup,
			"outscale_image":                   resourceImage,
			"outscale_keypair":                 resourceKeyPair,
			"outscale_public_ip":               resourcePublicIP,
			"outscale_public_ip_link":          resourcePublicIPLink,
			"outscale_volume":                  resourceVolume,
			"outscale_volumes_link":            resourceVolumeLink,
			"outscale_outbound_rule":           resourceOutboundRule,
			"outscale_security_group_rule":     resourceOutboundRule,
			"outscale_tag":                     resourceTags,
			"outscale_net_attributes":          resourceLinAttributes,
			"outscale_net":                     resourceNet,
			"outscale_internet_service":        resourceInternetService,
			"outscale_internet_service_link":   resourceInternetServiceLink,
			"outscale_nat_service":             resourceNatService,
			"outscale_subnet":                  resourceSubNet,
			"outscale_route":                   resourceRoute,
			"outscale_route_table":             resourceRouteTable,
			"outscale_route_table_link":        resourceLinkRouteTable,
			"outscale_snapshot":                resourceSnapshot,
			"outscale_image_launch_permission": resourceImageLaunchPermission,
			"outscale_net_peering":             resourceLinPeeringConnection,
			"outscale_nic_private_ip":          resourceNetworkInterfacePrivateIP,
			"outscale_nic_link":                resourceNetworkInterfaceAttachment,
			"outscale_nic":                     resourceNic,
			"outscale_snapshot_attributes":     resourcedSnapshotAttributes,
			"outscale_image_export_task":       resourceIMageExportTask,
			"outscale_net_peering_acceptation": resourceLinPeeringConnectionAccepter,
			"outscale_server_certificate":      resourceServerCertificate,
			"outscale_snapshot_export_task":    resourceSnapshotExportTask,
			"outscale_ca":                      resourceCa,
			"outscale_api_access_rule":         resourceApiAccessRule,
			"outscale_api_access_policy":       resourceApiAccessPolicy,
		},
	}

	datasources = Dictionary{
		"oapi": ResourceMap{
			"outscale_vm":                    dataSourceVM,
			"outscale_vms":                   dataSourceVMS,
			"outscale_firewall_rules_sets":   dataSourceSecurityGroups,
			"outscale_security_groups":       dataSourceSecurityGroups,
			"outscale_images":                dataSourceImages,
			"outscale_firewall_rules_set":    dataSourceSecurityGroup,
			"outscale_security_group":        dataSourceSecurityGroup,
			"outscale_tag":                   dataSourceTag,
			"outscale_tags":                  dataSourceTags,
			"outscale_volume":                dataSourceVolume,
			"outscale_volumes":               dataSourceVolumes,
			"outscale_keypair":               dataSourceKeyPair,
			"outscale_keypairs":              dataSourceKeyPairs,
			"outscale_internet_service":      dataSourceInternetService,
			"outscale_internet_services":     dataSourceInternetServices,
			"outscale_subnet":                dataSourceSubnet,
			"outscale_subnets":               dataSourceSubnets,
			"outscale_vm_state":              dataSourceVMState,
			"outscale_vm_states":             dataSourceVMStates,
			"outscale_net":                   dataSourceVpc,
			"outscale_nets":                  dataSourceVpcs,
			"outscale_net_attributes":        dataSourceVpcAttr,
			"outscale_route_table":           dataSourceRouteTable,
			"outscale_route_tables":          dataSourceRouteTables,
			"outscale_snapshot":              dataSourceSnapshot,
			"outscale_snapshots":             dataSourceSnapshots,
			"outscale_net_peering":           dataSourceLinPeeringConnection,
			"outscale_net_peerings":          dataSourceLinPeeringsConnection,
			"outscale_nic":                   dataSourceNic,
			"outscale_nics":                  dataSourceNics,
			"outscale_image":                 dataSourceImage,
			"outscale_public_ip":             dataSourcePublicIP,
			"outscale_public_ips":            dataSourcePublicIPS,
			"outscale_nat_service":           dataSourceNatService,
			"outscale_nat_services":          dataSourceNatServices,
			"outscale_subregions":            dataSourceSubregions,
			"outscale_regions":               dataSourceRegions,
			"outscale_image_export_task":     dataSourceImageExportTask,
			"outscale_image_export_tasks":    dataSourceImageExportTasks,
			"outscale_server_certificate":    dataSourceServerCertificate,
			"outscale_server_certificates":   dataSourceServerCertificates,
			"outscale_snapshot_export_task":  dataSourceSnapshotExportTask,
			"outscale_snapshot_export_tasks": dataSourceSnapshotExportTasks,
			"outscale_ca":                    dataSourceCa,
			"outscale_cas":                   dataSourceCas,
			"outscale_api_access_rule":       dataSourceApiAccessRule,
			"outscale_api_access_rules":      dataSourceApiAccessRules,
			"outscale_api_access_policy":     dataSourceApiAccessPolicy,
		},
	}
}

// GetResource ...
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

// GetDatasource receives the apu and the name of the datasource
// and returns the corrresponding
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
