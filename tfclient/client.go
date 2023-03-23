// This is derived from github.com/hashicorp/terraform/internal/providers/provider.go (15ecdb66c84cd8202b0ae3d34c44cb4bbece5444)

package tfclient

import (
	"context"
	"github.com/magodo/terraform-client-go/tfclient/client"
)

// Client represents the set of methods required for a complete resource
// provider plugin.
type Client interface {
	// GetProviderSchema returns the complete schema for the provider.
	GetProviderSchema() (*client.GetProviderSchemaResponse, client.Diagnostics)

	// ValidateProviderConfig allows the provider to validate the configuration.
	// The ValidateProviderConfigResponse.PreparedConfig field is unused. The
	// final configuration is not stored in the state, and any modifications
	// that need to be made must be made during the Configure method call.
	ValidateProviderConfig(context.Context, client.ValidateProviderConfigRequest) (*client.ValidateProviderConfigResponse, client.Diagnostics)

	// ValidateResourceConfig allows the provider to validate the resource
	// configuration values.
	ValidateResourceConfig(context.Context, client.ValidateResourceConfigRequest) (*client.ValidateResourceConfigResponse, client.Diagnostics)

	// ValidateDataResourceConfig allows the provider to validate the data source
	// configuration values.
	ValidateDataResourceConfig(context.Context, client.ValidateDataResourceConfigRequest) (*client.ValidateDataResourceConfigResponse, client.Diagnostics)

	// UpgradeResourceState is called when the state loader encounters an
	// instance state whose schema version is less than the one reported by the
	// currently-used version of the corresponding provider, and the upgraded
	// result is used for any further processing.
	UpgradeResourceState(context.Context, client.UpgradeResourceStateRequest) (*client.UpgradeResourceStateResponse, client.Diagnostics)

	// ConfigureProvider configures and initialized the provider.
	ConfigureProvider(context.Context, client.ConfigureProviderRequest) (*client.ConfigureProviderResponse, client.Diagnostics)

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
	ReadResource(context.Context, client.ReadResourceRequest) (*client.ReadResourceResponse, client.Diagnostics)

	// PlanResourceChange takes the current state and proposed state of a
	// resource, and returns the planned final state.
	PlanResourceChange(context.Context, client.PlanResourceChangeRequest) (*client.PlanResourceChangeResponse, client.Diagnostics)

	// ApplyResourceChange takes the planned state for a resource, which may
	// yet contain unknown computed values, and applies the changes returning
	// the final state.
	ApplyResourceChange(context.Context, client.ApplyResourceChangeRequest) (*client.ApplyResourceChangeResponse, client.Diagnostics)

	// ImportResourceState requests that the given resource be imported.
	ImportResourceState(context.Context, client.ImportResourceStateRequest) (*client.ImportResourceStateResponse, client.Diagnostics)

	// ReadDataSource returns the data source's current state.
	ReadDataSource(context.Context, client.ReadDataSourceRequest) (*client.ReadDataSourceResponse, client.Diagnostics)

	// Close shuts down the plugin process if applicable.
	Close()
}
