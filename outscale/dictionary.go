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
			"outscale_vm":        resourceOutscaleVM,
			"outscale_public_ip": resourceOutscalePublicIP,
		},
		"oapi": ResourceMap{
			"outscale_vm": resourceOutscaleOApiVM,
		},
	}
	datasources = Dictionary{
		"fcu": ResourceMap{
			"outscale_vm":  dataSourceOutscaleVM,
			"outscale_vms": dataSourceOutscaleVMS,
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
