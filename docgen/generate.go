package main

import (
	"bytes"
	"context"
	_ "embed"
	"errors"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"text/template"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/outscale/terraform-provider-outscale/provider"
	"github.com/outscale/terraform-provider-outscale/version"
	"github.com/samber/lo"
)

//go:embed templates/resources.md.tmpl
var resourceTemplate string

// Field represents a Terraform block/attribute
type Field struct {
	// NestedFields contains the childrens (ordered by name)
	NestedFields []Field

	// Information about the field extracted from the schema
	Name               string
	Description        string
	Depth              int
	Required           bool
	Optional           bool
	Computed           bool
	Sensitive          bool
	Deprecated         bool
	DeprecationMessage string

	// Classifies a field as an "Argument"; which is a configurable
	// (required/optional) field or a field that contains configurable childrens
	IsArgument bool
	// Classifies a field as an "Attribute"; which is a read-only (computed) field
	// or a field that contains read-only childrens
	IsAttribute bool
}

type Timeouts struct {
	Create bool
	Read   bool
	Update bool
	Delete bool
}

// Resource contains the template input of a resource
type Resource struct {
	Name string
	// API identifies the service: OUTSCALE, OKS, OOS...
	API string
	// All the fields of the resource (ordered by name)
	Fields []Field
	// The timeout operation the resource supports (create, update...)
	Timeouts Timeouts
}

func main() {
	if err := generate(); err != nil {
		log.Fatal(err)
	}
}

func generate() error {
	ctx := context.Background()

	targets := os.Args[1:]
	if len(targets) == 0 {
		return errors.New("specify at least one resource to generate (e.g. outscale_net outscale_oos_bucket)")
	}

	pwd, err := os.Getwd()
	if err != nil {
		return err
	}
	outDir := filepath.Join(pwd, "docs", "resources")

	root, err := os.OpenRoot(outDir)
	if err != nil {
		return err
	}
	defer root.Close() //nolint:errcheck

	server, err := provider.NewProtocol6Server(ctx, version.Version)
	if err != nil {
		return err
	}

	schema, err := server.GetProviderSchema(ctx, &tfprotov6.GetProviderSchemaRequest{})
	if err != nil {
		return err
	}

	for _, target := range targets {
		keySchema := schema.ResourceSchemas[target]
		if keySchema == nil {
			return errors.New("target not found")
		}

		api := "OUTSCALE"
		if strings.HasPrefix(target, "outscale_oks_") {
			api = "OKS"
		} else if strings.HasPrefix(target, "outscale_oos_") {
			api = "OOS"
		}

		resource := Resource{
			Name: target,
			API:  api,
		}

		// Flatten all blocks/attributes into a slice of Fields
		fields := walkBlock(keySchema.Block, 0)

		toTimeouts := func(fields []Field) (timeouts Timeouts) {
			for _, timeout := range fields {
				switch timeout.Name {
				case "create":
					timeouts.Create = true
				case "read":
					timeouts.Read = true
				case "update":
					timeouts.Update = true
				case "delete":
					timeouts.Delete = true
				}
			}

			return
		}

		// Filter out the timeouts from fields and assign them into the Timeouts element
		for _, field := range fields {
			if field.Name == "timeouts" {
				resource.Timeouts = toTimeouts(field.NestedFields)
				continue
			}

			resource.Fields = append(resource.Fields, field)
		}

		tmpl, err := template.New("resources.md.tmpl").Funcs(template.FuncMap{
			"indent": func(depth int) string {
				return strings.Repeat("\x20", 4*depth)
			},
			"replace":   strings.ReplaceAll,
			"trimspace": strings.TrimSpace,
		}).Parse(resourceTemplate)
		if err != nil {
			return err
		}

		var page bytes.Buffer
		if err := tmpl.Execute(&page, resource); err != nil {
			return err
		}

		filename := strings.TrimPrefix(target, "outscale_") + ".md"
		err = root.WriteFile(filename, page.Bytes(), 0o600)
		if err != nil {
			return err
		}
	}
	return nil
}

func walkAttributes(attributes []*tfprotov6.SchemaAttribute, depth int) []Field {
	return lo.FilterMap(attributes, func(attr *tfprotov6.SchemaAttribute, _ int) (Field, bool) {
		if attr == nil {
			return Field{}, false
		}

		field := Field{
			Name:               attr.Name,
			Description:        attr.Description,
			Depth:              depth,
			Required:           attr.Required,
			Optional:           attr.Optional,
			Computed:           attr.Computed,
			Sensitive:          attr.Sensitive,
			Deprecated:         attr.Deprecated,
			DeprecationMessage: attr.DeprecationMessage,
		}

		switch {
		case attr.NestedType != nil:
			field.NestedFields = walkAttributes(attr.NestedType.Attributes, depth+1)
			sort(field.NestedFields)
		default:
			// A Type can be complex (List, Set, Object...)
			// In this case, we need verify if the associated type is an object or a simple type
			field.NestedFields = walkType(attr.Type, field)
		}

		return field.classify(), true
	})
}

func walkType(attrType tftypes.Type, parent Field) []Field {
	switch t := attrType.(type) {
	case tftypes.List:
		return walkType(t.ElementType, parent)
	case tftypes.Set:
		return walkType(t.ElementType, parent)
	case tftypes.Map:
		return walkType(t.ElementType, parent)
	case tftypes.Object:
		// Objects contains nested childrens (a map of name->type) without flags (computed, optional or required)
		// We only retrieve their name and set the parent flags
		fields := lo.MapToSlice(t.AttributeTypes, func(name string, childType tftypes.Type) Field {
			field := Field{
				Name:      name,
				Depth:     parent.Depth + 1,
				Required:  parent.Required,
				Optional:  parent.Optional,
				Computed:  parent.Computed,
				Sensitive: parent.Sensitive,
			}
			field.NestedFields = walkType(childType, field)

			return field.classify()
		})
		sort(fields)

		return fields
	// A non-complex type is the leaf, we stop
	default:
		return nil
	}
}

func walkBlock(block *tfprotov6.SchemaBlock, depth int) []Field {
	if block == nil {
		return nil
	}

	fields := walkAttributes(block.Attributes, depth)

	blocks := lo.FilterMap(block.BlockTypes, func(block *tfprotov6.SchemaNestedBlock, _ int) (Field, bool) {
		if block == nil || block.Block == nil {
			return Field{}, false
		}

		children := walkBlock(block.Block, depth+1)
		required := block.MinItems == 1 && block.MaxItems == 1

		// Blocks do not have any optional/computed fields so we need to infer it from the nested block and it's fields
		optional := !required && (len(children) == 0 || lo.SomeBy(children, func(child Field) bool {
			return child.Required || child.Optional
		}))
		computed := !required && !optional && lo.EveryBy(children, func(child Field) bool {
			return child.Computed && !child.Required && !child.Optional
		})

		field := Field{
			NestedFields:       children,
			Name:               block.TypeName,
			Description:        block.Block.Description,
			Depth:              depth,
			Required:           required,
			Optional:           optional,
			Computed:           computed,
			Deprecated:         block.Block.Deprecated,
			DeprecationMessage: block.Block.DeprecationMessage,
		}
		return field.classify(), true
	})

	fields = append(fields, blocks...)
	sort(fields)

	return fields
}

func (field Field) classify() Field {
	field.IsArgument = field.Required || field.Optional
	field.IsAttribute = field.Computed && !field.IsArgument

	// It is necessary to propagate the child classification to the parents
	// because rendering happens from parent to child
	for _, child := range field.NestedFields {
		field.IsArgument = field.IsArgument || child.IsArgument
		field.IsAttribute = field.IsAttribute || child.IsAttribute
	}

	return field
}

func sort(fields []Field) {
	slices.SortFunc(fields, func(f1, f2 Field) int {
		return strings.Compare(f1.Name, f2.Name)
	})
}
