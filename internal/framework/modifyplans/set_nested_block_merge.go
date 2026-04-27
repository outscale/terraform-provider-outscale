package modifyplans

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/samber/lo"
)

// SetNestedBlockMergeOptions configures SetNestedBlockMerge.
type SetNestedBlockMergeOptions struct {
	// ObjectType is the nested object type of the set block elements.
	ObjectType types.ObjectType

	// Fields are the nested attribute names used to correlate config and
	// state elements. Only these fields participate in the element signature.
	Fields []string
}

// SetNestedBlockMerge returns a planmodifier.Set that works around
// terraform-plugin-framework issue #884 for optional+computed SetNestedBlock
// attributes with computed nested fields.
//
// The framework hashes set elements by ALL nested attribute values (including
// Computed ones). When computed fields differ between config and state the element hash changes and
// Terraform produces a perpetual add/remove diff. This plan modifier avoids
// that by comparing elements using only user-configurable fields (the
// "signature").
//
// When signatures match exactly, the state value is preserved unchanged so
// computed nested fields stay stable. When signatures differ, the config
// elements are used with configured unknown field rewrites so Terraform shows
// "known after apply" for API-populated nested fields.
//
// On create with an omitted block one unknown element is injected by default.
// This is useful when the API always populates at least one element and you want the plan to show
// "(known after apply)" rather than omitting the block entirely.
//
// On update with an omitted block the state value is preserved.
func SetNestedBlockMerge(options SetNestedBlockMergeOptions) planmodifier.Set {
	return setNestedBlockMergePlanModifier{
		options:                      options,
		injectUnknownElementOnCreate: true,
	}
}

// SetNestedBlockMergeEmptyOnCreate returns the same workaround behavior
// without injecting an unknown element on create when the block is omitted.
// Use this when omission is expected to mean an actually empty set.
func SetNestedBlockMergeEmptyOnCreate(options SetNestedBlockMergeOptions) planmodifier.Set {
	return setNestedBlockMergePlanModifier{
		options:                      options,
		injectUnknownElementOnCreate: false,
	}
}

type setNestedBlockMergePlanModifier struct {
	options                      SetNestedBlockMergeOptions
	injectUnknownElementOnCreate bool
}

func (m setNestedBlockMergePlanModifier) Description(_ context.Context) string {
	return "Merges optional+computed set nested block to avoid perpetual diffs from computed nested fields."
}

func (m setNestedBlockMergePlanModifier) MarkdownDescription(_ context.Context) string {
	return m.Description(context.Background())
}

func (m setNestedBlockMergePlanModifier) PlanModifySet(ctx context.Context, req planmodifier.SetRequest, resp *planmodifier.SetResponse) {
	if len(m.options.Fields) == 0 {
		resp.Diagnostics.AddError("Invalid SetNestedBlockMerge configuration", "Fields must not be empty.")
		return
	}

	// Create with omitted block: optionally inject a known set containing one
	// fully-unknown element so the plan signals "known after apply". A
	// SetNestedBlock itself must never be unknown; only its nested attributes
	// may be unknown.
	if req.StateValue.IsNull() && (req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown()) {
		if m.injectUnknownElementOnCreate {
			elem, diags := m.unknownObjectFromType(ctx)
			resp.Diagnostics.Append(diags...)
			if resp.Diagnostics.HasError() {
				return
			}

			unknownSet, diags := types.SetValue(m.options.ObjectType, []attr.Value{elem})
			resp.Diagnostics.Append(diags...)
			if resp.Diagnostics.HasError() {
				return
			}

			resp.PlanValue = unknownSet
		}
		return
	}

	// No config but state exists: update with omitted block (preserve state)
	// or destroy (preserving state is harmless because the resource is being
	// removed anyway).
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() || len(req.ConfigValue.Elements()) == 0 {
		resp.PlanValue = req.StateValue
		return
	}

	// Both config and state have values: compare signatures.
	var configElements []types.Object
	resp.Diagnostics.Append(req.ConfigValue.ElementsAs(ctx, &configElements, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	configKeySet, diags := m.projectToKeySet(configElements)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	stateKeySet := types.SetNull(m.keyObjectType())
	if !req.StateValue.IsNull() && !req.StateValue.IsUnknown() {
		var stateElements []types.Object
		resp.Diagnostics.Append(req.StateValue.ElementsAs(ctx, &stateElements, false)...)
		if resp.Diagnostics.HasError() {
			return
		}

		stateKeySet, diags = m.projectToKeySet(stateElements)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	unchanged := configKeySet.Equal(stateKeySet)

	if unchanged {
		resp.PlanValue = req.StateValue
		return
	}

	// Signatures differ: build merged set from config with computed fields unknown.
	merged := make([]attr.Value, 0, len(configElements))
	for i := range configElements {
		elem, diags := m.prepareChangedObject(ctx, configElements[i])
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		merged = append(merged, elem)
	}

	mergedSet, diags := types.SetValue(m.options.ObjectType, merged)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.PlanValue = mergedSet
}

func (m setNestedBlockMergePlanModifier) projectToKeySet(elements []types.Object) (types.Set, diag.Diagnostics) {
	var diags diag.Diagnostics
	keyObjType := m.keyObjectType()

	keyObjects := make([]attr.Value, 0, len(elements))
	for _, elem := range elements {
		keyObj, d := m.projectToKeyObject(elem)
		diags.Append(d...)
		if diags.HasError() {
			return types.SetNull(keyObjType), diags
		}

		keyObjects = append(keyObjects, keyObj)
	}

	setValue, setDiags := types.SetValue(keyObjType, keyObjects)
	diags.Append(setDiags...)
	if diags.HasError() {
		return types.SetNull(keyObjType), diags
	}

	return setValue, diags
}

func (m setNestedBlockMergePlanModifier) projectToKeyObject(elem types.Object) (attr.Value, diag.Diagnostics) {
	var diags diag.Diagnostics
	keyObjType := m.keyObjectType()
	attrs := elem.Attributes()
	keyAttrs := make(map[string]attr.Value, len(m.options.Fields))
	for _, key := range m.options.Fields {
		fieldType, ok := m.options.ObjectType.AttrTypes[key]
		if !ok {
			diags.AddError(
				"Invalid SetNestedBlockMerge configuration",
				fmt.Sprintf("key field %q is not present in object type", key),
			)
			return types.ObjectNull(keyObjType.AttrTypes), diags
		}

		v, ok := attrs[key]
		if !ok {
			diags.AddError(
				"Invalid SetNestedBlockMerge configuration",
				fmt.Sprintf("key field %q is not present in object type", key),
			)
			return types.ObjectNull(keyObjType.AttrTypes), diags
		}

		keyAttrs[key] = normalizeKeyValue(v, fieldType)
	}

	obj, objectDiags := types.ObjectValue(keyObjType.AttrTypes, keyAttrs)
	diags.Append(objectDiags...)
	if diags.HasError() {
		return types.ObjectNull(keyObjType.AttrTypes), diags
	}

	return obj, diags
}

func (m setNestedBlockMergePlanModifier) keyObjectType() types.ObjectType {
	keyAttrTypes := make(map[string]attr.Type, len(m.options.Fields))
	for _, key := range m.options.Fields {
		fieldType, ok := m.options.ObjectType.AttrTypes[key]
		if ok {
			keyAttrTypes[key] = fieldType
		}
	}

	return types.ObjectType{AttrTypes: keyAttrTypes}
}

func (m setNestedBlockMergePlanModifier) prepareChangedObject(ctx context.Context, elem types.Object) (attr.Value, diag.Diagnostics) {
	var diags diag.Diagnostics

	attrs := elem.Attributes()
	newAttrs := make(map[string]attr.Value, len(attrs))
	for key, value := range attrs {
		newAttrs[key] = value
	}

	for key, fieldType := range m.options.ObjectType.AttrTypes {
		if lo.Contains(m.options.Fields, key) {
			current, ok := newAttrs[key]
			if !ok {
				diags.AddError(
					"Invalid SetNestedBlockMerge configuration",
					fmt.Sprintf("key field %q is not present in object value", key),
				)
				continue
			}
			if current.IsNull() {
				newAttrs[key] = normalizeKeyValue(current, fieldType)
			}
			continue
		}

		current, ok := newAttrs[key]
		if !ok {
			diags.AddError(
				"Invalid SetNestedBlockMerge configuration",
				fmt.Sprintf("field %q is not present in object value", key),
			)
			continue
		}

		// If the field is not configured, mark it unknown so API-populated values become known-after-apply.
		if current.IsNull() || current.IsUnknown() {
			unknownValue, typeDiags := unknownValueFromType(ctx, fieldType)
			diags.Append(typeDiags...)
			if diags.HasError() {
				return types.ObjectNull(m.options.ObjectType.AttrTypes), diags
			}
			newAttrs[key] = unknownValue
		}
	}

	if diags.HasError() {
		return types.ObjectNull(m.options.ObjectType.AttrTypes), diags
	}

	obj, objectDiags := types.ObjectValue(m.options.ObjectType.AttrTypes, newAttrs)
	diags.Append(objectDiags...)
	if diags.HasError() {
		return types.ObjectNull(m.options.ObjectType.AttrTypes), diags
	}

	return obj, diags
}

func (m setNestedBlockMergePlanModifier) unknownObjectFromType(ctx context.Context) (attr.Value, diag.Diagnostics) {
	var diags diag.Diagnostics

	attrs := make(map[string]attr.Value, len(m.options.ObjectType.AttrTypes))
	for key, fieldType := range m.options.ObjectType.AttrTypes {
		value, fieldDiags := unknownValueFromType(ctx, fieldType)
		diags.Append(fieldDiags...)
		if diags.HasError() {
			return types.ObjectNull(m.options.ObjectType.AttrTypes), diags
		}
		attrs[key] = value
	}

	if diags.HasError() {
		return types.ObjectNull(m.options.ObjectType.AttrTypes), diags
	}

	obj, objectDiags := types.ObjectValue(m.options.ObjectType.AttrTypes, attrs)
	diags.Append(objectDiags...)
	if diags.HasError() {
		return types.ObjectNull(m.options.ObjectType.AttrTypes), diags
	}

	return obj, diags
}

func unknownValueFromType(ctx context.Context, t attr.Type) (attr.Value, diag.Diagnostics) {
	var diags diag.Diagnostics

	tfType := t.TerraformType(ctx)
	unknownTfValue := tftypes.NewValue(tfType, tftypes.UnknownValue)
	value, err := t.ValueFromTerraform(ctx, unknownTfValue)
	if err != nil {
		diags.AddError(
			"Unsupported attribute type for unknown generation",
			fmt.Sprintf("cannot auto-generate unknown value for attribute type %T: %s", t, err),
		)
		return nil, diags
	}

	return value, diags
}

func normalizeKeyValue(v attr.Value, t attr.Type) attr.Value {
	if v.IsNull() {
		if _, ok := t.(basetypes.BoolType); ok {
			return types.BoolValue(false)
		}
	}

	return v
}
