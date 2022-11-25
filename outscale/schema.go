package outscale

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func GetDataSourcesSchema(field string, resourceSchema map[string]*schema.Schema) map[string]*schema.Schema {
	result := map[string]*schema.Schema{
		field: {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: setComputed(resourceSchema),
			},
		},
	}
	result = addFilter(result)
	result = addRequestId(result)
	return result
}

func GetDataSourceSchema(resourceSchema map[string]*schema.Schema) map[string]*schema.Schema {
	result := setComputed(resourceSchema)
	result = addFilter(result)
	result = addRequestId(result)
	return result
}

func GetResourceSchema(resourceSchema map[string]*schema.Schema) map[string]*schema.Schema {
	result := addRequestId(resourceSchema)
	return result
}

func addFilter(resourceSchema map[string]*schema.Schema) map[string]*schema.Schema {
	resourceSchema["filter"] = dataSourceFiltersSchema()
	return resourceSchema
}

func addRequestId(resourceSchema map[string]*schema.Schema) map[string]*schema.Schema {
	resourceSchema["request_id"] = requestIdSchema()
	return resourceSchema
}
func setComputed(resourceSchema map[string]*schema.Schema) map[string]*schema.Schema {
	for k, v := range resourceSchema {
		v.Computed = true
		v.Required = false
		v.Optional = false
		v.ForceNew = false
		v.Default = nil
		v.DiffSuppressFunc = nil
		v.ValidateFunc = nil
		resourceSchema[k] = v
	}
	return resourceSchema
}

func requestIdSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeString,
		Computed: true,
	}
}

func dataSourceFiltersSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		ForceNew: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name": {
					Type:     schema.TypeString,
					Required: true,
				},

				"values": {
					Type:     schema.TypeList,
					Required: true,
					Elem:     &schema.Schema{Type: schema.TypeString},
				},
			},
		},
	}
}
