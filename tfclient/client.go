// This is derived from github.com/hashicorp/terraform/internal/providers/provider.go (15ecdb66c84cd8202b0ae3d34c44cb4bbece5444)

package tfclient

import (
	"context"

	"github.com/magodo/terraform-client-go/tfclient/tfprotov5/tf5client"
	"github.com/magodo/terraform-client-go/tfclient/tfprotov6/tf6client"
	"github.com/magodo/terraform-client-go/tfclient/typ"
)

var _ Client = &tf5client.Client{}
var _ Client = &tf6client.Client{}

// Client represents the set of methods required for a complete resource
// provider plugin.
type Client interface {
	// GetProviderSchema returns the complete schema for the provider.
	GetProviderSchema() (*typ.GetProviderSchemaResponse, typ.Diagnostics)

	// GetResourceIdentitySchemas returns the identity schemas for all managed resources
	// for the provider. Usually you don't need to call this method directly as GetProviderSchema
	// will merge the identity schemas into the provider schema.
	GetResourceIdentitySchemas(context.Context) (*typ.GetResourceIdentitySchemasResponse, typ.Diagnostics)

	// ValidateProviderConfig allows the provider to validate the configuration.
	// The ValidateProviderConfigResponse.PreparedConfig field is unused. The
	// final configuration is not stored in the state, and any modifications
	// that need to be made must be made during the Configure method call.
	ValidateProviderConfig(context.Context, typ.ValidateProviderConfigRequest) (*typ.ValidateProviderConfigResponse, typ.Diagnostics)

	// ValidateResourceConfig allows the provider to validate the resource
	// configuration values.
	ValidateResourceConfig(context.Context, typ.ValidateResourceConfigRequest) (*typ.ValidateResourceConfigResponse, typ.Diagnostics)

	// ValidateDataResourceConfig allows the provider to validate the data source
	// configuration values.
	ValidateDataResourceConfig(context.Context, typ.ValidateDataResourceConfigRequest) (*typ.ValidateDataResourceConfigResponse, typ.Diagnostics)

	// ValidateEphemeralResourceConfig allows the provider to validate the
	// ephemeral resource configuration values.
	ValidateEphemeralResourceConfig(context.Context, typ.ValidateEphemeralResourceConfigRequest) typ.Diagnostics

	// ValidateListResourceConfig allows the provider to validate the list
	// resource configuration values.
	ValidateListResourceConfig(context.Context, typ.ValidateListResourceConfigRequest) typ.Diagnostics

	// ValidateActionConfig allows the provider to validate the action configuration values.
	ValidateActionConfig(context.Context, typ.ValidateActionConfigRequest) typ.Diagnostics

	// UpgradeResourceState is called when the state loader encounters an
	// instance state whose schema version is less than the one reported by the
	// currently-used version of the corresponding provider, and the upgraded
	// result is used for any further processing.
	UpgradeResourceState(context.Context, typ.UpgradeResourceStateRequest) (*typ.UpgradeResourceStateResponse, typ.Diagnostics)

	// UpgradeResourceIdentity is called when the state loader encounters an
	// instance identity whose schema version is less than the one reported by
	// the currently-used version of the corresponding provider, and the upgraded
	// result is used for any further processing.
	UpgradeResourceIdentity(context.Context, typ.UpgradeResourceIdentityRequest) (*typ.UpgradeResourceIdentityResponse, typ.Diagnostics)

	// ConfigureProvider configures and initialized the provider.
	ConfigureProvider(context.Context, typ.ConfigureProviderRequest) (*typ.ConfigureProviderResponse, typ.Diagnostics)

	// Stop is called when the provider should halt any in-flight actions.
	//
	// Stop should not block waiting for in-flight actions to complete. It
	// should take any action it wants and return immediately acknowledging it
	// has received the stop request. Terraform will not make any further API
	// calls to the provider after Stop is called.
	//
	// The error returned, if non-nil, is assumed to mean that signaling the
	// stop somehow failed and that the user should expect potentially waiting
	// a longer period of time.
	Stop(context.Context) error

	// ReadResource refreshes a resource and returns its current state.
	ReadResource(context.Context, typ.ReadResourceRequest) (*typ.ReadResourceResponse, typ.Diagnostics)

	// PlanResourceChange takes the current state and proposed state of a
	// resource, and returns the planned final state.
	PlanResourceChange(context.Context, typ.PlanResourceChangeRequest) (*typ.PlanResourceChangeResponse, typ.Diagnostics)

	// ApplyResourceChange takes the planned state for a resource, which may
	// yet contain unknown computed values, and applies the changes returning
	// the final state.
	ApplyResourceChange(context.Context, typ.ApplyResourceChangeRequest) (*typ.ApplyResourceChangeResponse, typ.Diagnostics)

	// ImportResourceState requests that the given resource be imported.
	ImportResourceState(context.Context, typ.ImportResourceStateRequest) (*typ.ImportResourceStateResponse, typ.Diagnostics)

	// MoveResourceState retrieves the updated value for a resource after it
	// has moved resource types.
	MoveResourceState(context.Context, typ.MoveResourceStateRequest) (*typ.MoveResourceStateResponse, typ.Diagnostics)

	// ReadDataSource returns the data source's current state.
	ReadDataSource(context.Context, typ.ReadDataSourceRequest) (*typ.ReadDataSourceResponse, typ.Diagnostics)

	// OpenEphemeralResource opens an ephemeral resource instance.
	OpenEphemeralResource(context.Context, typ.OpenEphemeralResourceRequest) (*typ.OpenEphemeralResourceResponse, typ.Diagnostics)
	// RenewEphemeralResource extends the validity of a previously-opened ephemeral
	// resource instance.
	RenewEphemeralResource(context.Context, typ.RenewEphemeralResourceRequest) (*typ.RenewEphemeralResourceResponse, typ.Diagnostics)
	// CloseEphemeralResource closes an ephemeral resource instance, with the intent
	// of rendering it invalid as soon as possible.
	CloseEphemeralResource(context.Context, typ.CloseEphemeralResourceRequest) typ.Diagnostics

	// CallFunction calls a provider-contributed function.
	CallFunction(context.Context, typ.CallFunctionRequest) (*typ.CallFunctionResponse, typ.Diagnostics)

	// ListResource lists resources
	ListResource(context.Context, typ.ListResourceRequest) (typ.ListResourceResponse, typ.Diagnostics)

	// PlanAction takes the proposed action config and returns the plan
	PlanAction(context.Context, typ.PlanActionRequest) (typ.PlanActionResponse, typ.Diagnostics)

	// InvokeAction invokes an action
	InvokeAction(context.Context, typ.InvokeActionRequest) (typ.InvokeActionResponse, typ.Diagnostics)

	// Close shuts down the plugin process if applicable.
	Close()
}
