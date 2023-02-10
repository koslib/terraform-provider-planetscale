package planetscale

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/planetscale/planetscale-go/planetscale"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces
var (
	_ provider.Provider = &planetscaleProvider{}
)

// New is a helper function to simplify provider server and testing implementation.
func New() provider.Provider {
	return &planetscaleProvider{}
}

// planetscaleProviderModel maps provider schema data to a Go type.
type planetscaleProviderModel struct {
	ServiceTokenID types.String `tfsdk:"service_token_id"`
	ServiceToken   types.String `tfsdk:"service_token"`
}

// planetscaleProvider is the provider implementation.
type planetscaleProvider struct{}

// Metadata returns the provider type name.
func (p *planetscaleProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "planetscale"
}

// Schema defines the provider-level schema for configuration data.
func (p *planetscaleProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"service_token": schema.StringAttribute{
				Required: true,
			},
			"service_token_id": schema.StringAttribute{
				Required: true,
			},
		},
	}
}

// Configure prepares a Planetscale API client for data sources and resources.
func (p *planetscaleProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Info(ctx, "configuring Planetscale client")

	// Retrieve provider data from configuration
	var config planetscaleProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If practitioner provided a configuration value for any of the
	// attributes, it must be a known value.
	if config.ServiceTokenID.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("service_token_id"),
			"Unknown ServiceTokenID",
			"The provider cannot create the Planetscale API client as there is an unknown configuration value for the Planetscale API service token id. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the PLANETSCALE_SERVICE_TOKEN_ID environment variable.",
		)
	}

	if config.ServiceTokenID.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("service_token"),
			"Unknown ServiceTokenID",
			"The provider cannot create the Planetscale API client as there is an unknown configuration value for the Planetscale API service token. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the PLANETSCALE_SERVICE_TOKEN environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Default values to environment variables, but override
	// with Terraform configuration value if set.

	serviceTokenID := os.Getenv("PLANETSCALE_SERVICE_TOKEN_ID")
	serviceToken := os.Getenv("PLANETSCALE_SERVICE_TOKEN")

	if !config.ServiceTokenID.IsNull() {
		serviceTokenID = config.ServiceTokenID.ValueString()
	}

	if !config.ServiceToken.IsNull() {
		serviceToken = config.ServiceToken.ValueString()
	}

	// If any of the expected configurations are missing, return
	// errors with provider-specific guidance.

	if serviceTokenID == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("service_token_id"),
			"Missing ServiceTokenID",
			"The provider cannot create the Planetscale API client as there is a missing or empty value for the service token id. "+
				"Set the service token id value in the provider configuration or use the PLANETSCALE_SERVICE_TOKEN_ID environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if serviceToken == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("service_token"),
			"Missing ServiceToken",
			"The provider cannot create the Planetscale API client as there is a missing or empty value for the service token. "+
				"Set the service token id value in the provider configuration or use the PLANETSCALE_SERVICE_TOKEN environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	ctx = tflog.SetField(ctx, "service_token_id", serviceTokenID)
	ctx = tflog.SetField(ctx, "service_token", serviceToken)
	ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "service_token")

	tflog.Info(ctx, "creating Planetscale client")

	// Create a new Planetscale API client using the configuration values
	client, err := planetscale.NewClient(planetscale.WithServiceToken(serviceTokenID, serviceToken))
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to create a Planetscale API Client",
			"An unexpected error occurred when creating the Planetscale API client. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"Planetscale Client Error: "+err.Error(),
		)
		return
	}

	// Make the Planetscale client available during DataSource and Resource
	// type Configure methods.
	resp.DataSourceData = client
	resp.ResourceData = client

	tflog.Info(ctx, "Configured Planetscale client", map[string]any{"success": true})
}

// DataSources defines the data sources implemented in the provider.
func (p *planetscaleProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewDatabasesDataSource,
		NewRegionsDataSource,
		NewDatabaseBranchesDataSource,
		NewDatabaseBranchPasswordDataSource,
	}
}

// Resources defines the resources implemented in the provider.
func (p *planetscaleProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewDatabaseResource,
		NewDatabaseBranchResource,
		NewDatabaseBranchPasswordResource,
	}
}
