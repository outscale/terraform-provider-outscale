package fwtypes

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
)

type AttributeTypes interface {
	AttributeTypes() map[string]attr.Type
}

func AttributeTypesAtPath(root attr.Type, path ...string) map[string]attr.Type {
	current := root

	for _, name := range path {
		if collection, ok := current.(attr.TypeWithElementType); ok {
			current = collection.ElementType()
		}

		object, ok := current.(attr.TypeWithAttributeTypes)
		if !ok {
			panic(fmt.Sprintf("type at path %v does not provide attribute types", path))
		}

		current, ok = object.AttributeTypes()[name]
		if !ok {
			panic(fmt.Sprintf("attribute %q does not exist at path %v", name, path))
		}
	}

	if collection, ok := current.(attr.TypeWithElementType); ok {
		current = collection.ElementType()
	}

	object, ok := current.(attr.TypeWithAttributeTypes)
	if !ok {
		panic(fmt.Sprintf("type at path %v does not provide attribute types", path))
	}

	return object.AttributeTypes()
}
