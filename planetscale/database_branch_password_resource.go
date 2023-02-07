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
	_ resource.Resource              = &databaseBranchPasswordResource{}
	_ resource.ResourceWithConfigure = &databaseBranchPasswordResource{}
)

type databaseBranchPasswordResourceModel struct {
	Name         types.String `tfsdk:"name"`
	Branch       types.String `tfsdk:"branch"`
	Database     types.String `tfsdk:"database"`
	Organization types.String `tfsdk:"organization"`
	Role         types.String `tfsdk:"role"`
	PublicID     types.String `tfsdk:"public_id"`
	Username     types.String `tfsdk:"username"`
}

// NewDatabaseBranchPasswordResource is a helper function to simplify the provider implementation.
func NewDatabaseBranchPasswordResource() resource.Resource {
	return &databaseBranchPasswordResource{}
}

// databaseResource is the resource implementation.
type databaseBranchPasswordResource struct {
	client *planetscale.Client
}

// Metadata returns the resource type name.
func (r *databaseBranchPasswordResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_database_branch_password"
}

// Schema defines the schema for the resource.
func (r *databaseBranchPasswordResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Required: true,
			},
			"branch": schema.StringAttribute{
				Required: true,
			},
			"database": schema.StringAttribute{
				Required: true,
			},
			"organization": schema.StringAttribute{
				Required: true,
			},
			"role": schema.StringAttribute{
				Optional: true,
			},
			"public_id": schema.StringAttribute{
				Optional: true,
			},
			"username": schema.StringAttribute{
				Computed: true,
			},
			// todo: add plaintext password as a sensitive string attribute
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *databaseBranchPasswordResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan databaseBranchPasswordResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// create resource on Planetscale
	databaseBranchPassword, err := r.client.Passwords.Create(ctx, &planetscale.DatabaseBranchPasswordRequest{
		Organization: plan.Organization.ValueString(),
		Database:     plan.Database.ValueString(),
		Branch:       plan.Branch.ValueString(),
		Name:         plan.Name.ValueString(),
		Role:         plan.Role.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating database branch password",
			"Could not create database branch password, unexpected error: "+err.Error(),
		)
		return
	}

	plan.PublicID = types.StringValue(databaseBranchPassword.PublicID)
	plan.Username = types.StringValue(databaseBranchPassword.Username)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Read refreshes the Terraform state with the latest data.
func (r *databaseBranchPasswordResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state databaseBranchPasswordResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed database branch password value from Planetscale
	databaseBranchPassword, err := r.client.Passwords.Get(ctx, &planetscale.GetDatabaseBranchPasswordRequest{
		Organization: state.Organization.ValueString(),
		Database:     state.Database.ValueString(),
		Branch:       state.Branch.ValueString(),
		Name:         state.Name.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading Planetscale database branch password",
			"Could not read info about Planetscale database branch password "+state.Name.ValueString()+": "+err.Error(),
		)
		return
	}

	// Overwrite items with refreshed state
	state.PublicID = types.StringValue(databaseBranchPassword.PublicID)
	state.Username = types.StringValue(databaseBranchPassword.Username)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Update updates the resource and sets the updated Terraform state on success.
func (r *databaseBranchPasswordResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// not supported by the golang sdk (https://github.com/planetscale/planetscale-go) yet
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *databaseBranchPasswordResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state databaseBranchPasswordResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing Planetscale database branch password
	ctx = tflog.SetField(ctx, "organization", state.Organization.ValueString())
	ctx = tflog.SetField(ctx, "database", state.Database.ValueString())
	ctx = tflog.SetField(ctx, "branch", state.Branch.ValueString())
	ctx = tflog.SetField(ctx, "name", state.Name.ValueString())

	err := r.client.Passwords.Delete(ctx, &planetscale.DeleteDatabaseBranchPasswordRequest{
		Organization: state.Organization.ValueString(),
		Database:     state.Database.ValueString(),
		Branch:       state.Branch.ValueString(),
		Name:         state.Name.ValueString(),
		PasswordId:   state.PublicID.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error on deleting Planetscale database branch password",
			"Could not delete database branch password, unexpected error: "+err.Error(),
		)
		return
	}

	tflog.Info(ctx, "deleted Planetscale database branch password")
}

// Configure adds the provider configured client to the resource.
func (r *databaseBranchPasswordResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*planetscale.Client)
}
