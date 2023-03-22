// Package configschema is an adoption of a subset of the github.com/hashicorp/terraform/internal/configs/configschema@15ecdb66c84cd8202b0ae3d34c44cb4bbece5444.
// It only focus on the implied type (and its dependencies) of the schema `Block` type. But instead of the `Block` defined internally by terraform core, it target
// to the github.com/hashicorp/terraform-json.SchemaBlock.
package configschema
