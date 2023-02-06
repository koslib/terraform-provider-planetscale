package planetscale

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/planetscale/planetscale-go/planetscale"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

var (
	_ resource.Resource              = &databaseResource{}
	_ resource.ResourceWithConfigure = &databaseResource{}
)

type databaseResourceModel struct {
	Name         types.String `tfsdk:"name"`
	Notes        types.String `tfsdk:"notes"`
	Organization types.String `tfsdk:"organization"`
	Region       types.String `tfsdk:"region"`
	HtmlURL      types.String `tfsdk:"html_url"`
	State        types.String `tfsdk:"state"`
}

// NewDatabaseResource is a helper function to simplify the provider implementation.
func NewDatabaseResource() resource.Resource {
	return &databaseResource{}
}

// databaseResource is the resource implementation.
type databaseResource struct {
	client *planetscale.Client
}

// Metadata returns the resource type name.
func (r *databaseResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_database"
}

// Schema defines the schema for the resource.
func (r *databaseResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Required: true,
			},
			"notes": schema.StringAttribute{
				Optional: true,
			},
			"organization": schema.StringAttribute{
				Required: true,
			},
			"region": schema.StringAttribute{
				Optional: true,
			},
			"html_url": schema.StringAttribute{
				Computed: true,
			},
			"state": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *databaseResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan databaseResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// create resource on Planetscale
	database, err := r.client.Databases.Create(ctx, &planetscale.CreateDatabaseRequest{
		Organization: plan.Organization.ValueString(),
		Name:         plan.Name.ValueString(),
		Notes:        plan.Notes.ValueString(),
		Region:       plan.Region.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating database",
			"Could not create database, unexpected error: "+err.Error(),
		)
		return
	}

	plan.HtmlURL = types.StringValue(database.HtmlURL)
	plan.State = types.StringValue(string(database.State))

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Read refreshes the Terraform state with the latest data.
func (r *databaseResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state databaseResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed database value from Planetscale
	database, err := r.client.Databases.Get(ctx, &planetscale.GetDatabaseRequest{
		Organization: state.Organization.ValueString(),
		Database:     state.Name.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading Planetscale database",
			"Could not read info about Planetscale database "+state.Name.ValueString()+": "+err.Error(),
		)
		return
	}

	// Overwrite items with refreshed state
	state.HtmlURL = types.StringValue(database.HtmlURL)
	state.State = types.StringValue(string(database.State))

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Update updates the resource and sets the updated Terraform state on success.
func (r *databaseResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *databaseResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state databaseResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing Planetscale database
	ctx = tflog.SetField(ctx, "organization", state.Organization.ValueString())
	ctx = tflog.SetField(ctx, "database", state.Name.ValueString())

	_, err := r.client.Databases.Delete(ctx, &planetscale.DeleteDatabaseRequest{
		Organization: state.Organization.ValueString(),
		Database:     state.Name.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting Planetscale database",
			"Could not delete database, unexpected error: "+err.Error(),
		)
		return
	}
	tflog.Info(ctx, "deleted Planetscale database")
}

// Configure adds the provider configured client to the resource.
func (r *databaseResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*planetscale.Client)
}
