// This is derived from github.com/hashicorp/terraform/internal/providers/provider.go (15ecdb66c84cd8202b0ae3d34c44cb4bbece5444)

package tfclient

import (
	"context"
	"github.com/magodo/terraform-client-go/tfclient/typ"
)

// Client represents the set of methods required for a complete resource
// provider plugin.
type Client interface {
	// GetProviderSchema returns the complete schema for the provider.
	GetProviderSchema() (*typ.GetProviderSchemaResponse, typ.Diagnostics)

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

	// UpgradeResourceState is called when the state loader encounters an
	// instance state whose schema version is less than the one reported by the
	// currently-used version of the corresponding provider, and the upgraded
	// result is used for any further processing.
	UpgradeResourceState(context.Context, typ.UpgradeResourceStateRequest) (*typ.UpgradeResourceStateResponse, typ.Diagnostics)

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

	// ReadDataSource returns the data source's current state.
	ReadDataSource(context.Context, typ.ReadDataSourceRequest) (*typ.ReadDataSourceResponse, typ.Diagnostics)

	// Close shuts down the plugin process if applicable.
	Close()
}
