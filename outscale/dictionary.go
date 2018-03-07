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
			"outscale_vm":                 resourceOutscaleVM,
			"outscale_image":              resourceOutscaleImage,
			"outscale_firewall_rules_set": resourceOutscaleFirewallRulesSet,
			"outscale_outbound_rule":      resourceOutscaleOutboundRule,
			"outscale_inbound_rule":       resourceOutscaleInboundRule,
			"outscale_tag":                resourceOutscaleTags,
			"outscale_key_pair":           resourceOutscaleKeyPair,
			"outscale_public_ip":          resourceOutscalePublicIP,
			"outscale_public_ip_link":     resourceOutscalePublicIPLink,
		},
		"oapi": ResourceMap{
			"outscale_vm":    resourceOutscaleOApiVM,
			"outscale_image": resourceOutscaleOAPIImage,
		},
	}
	datasources = Dictionary{
		"fcu": ResourceMap{
			"outscale_vm":         dataSourceOutscaleVM,
			"outscale_vms":        dataSourceOutscaleVMS,
			"outscale_image":      dataSourceOutscaleImage,
			"outscale_images":     dataSourceOutscaleImages,
			"outscale_tag":        dataSourceOutscaleTag,
			"outscale_tags":       dataSourceOutscaleTags,
			"outscale_public_ip":  dataSourceOutscalePublicIP,
			"outscale_public_ips": dataSourceOutscalePublicIPS,
		},
		"oapi": ResourceMap{
			"outscale_vm":  dataSourceOutscaleOAPIVM,
			"outscale_vms": datasourceOutscaleOApiVMS,
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
