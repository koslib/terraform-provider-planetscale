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
	_ resource.Resource              = &deployRequestResource{}
	_ resource.ResourceWithConfigure = &deployRequestResource{}
)

type deployRequestModel struct {
	Organization    types.String `tfsdk:"organization"`
	Database        types.String `tfsdk:"database"`
	Branch          types.String `tfsdk:"branch"`
	IntoBranch      types.String `tfsdk:"into_branch"`
	Notes           types.String `tfsdk:"notes"`
	ID              types.String `tfsdk:"id"`
	State           types.String `tfsdk:"state"`
	DeploymentState types.String `tfsdk:"deployment_state"`
	HTMLURL         types.String `tfsdk:"html_url"`
	CreatedAt       types.String `tfsdk:"created_at"`
	UpdatedAt       types.String `tfsdk:"updated_at"`
	Approved        types.Bool   `tfsdk:"approved"`
	Number          types.Int64  `tfsdk:"number"`
}

// NewDeployRequestResource is a helper function to simplify the provider implementation.
func NewDeployRequestResource() resource.Resource {
	return &deployRequestResource{}
}

// deployRequestResource is the resource implementation.
type deployRequestResource struct {
	client *planetscale.Client
}

// Metadata returns the resource type name.
func (r *deployRequestResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_deploy_request"
}

// Schema defines the schema for the resource.
func (r *deployRequestResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "A Planetscale deploy request allows to create and revert non-blocking schema changes on a database" +
			" branch. More info: https://planetscale.com/docs/concepts/deploy-requests",
		Attributes: map[string]schema.Attribute{
			"organization": schema.StringAttribute{
				Required:    true,
				Description: "The name of the organization.",
			},
			"database": schema.StringAttribute{
				Required:    true,
				Description: "The name of the database.",
			},
			"notes": schema.StringAttribute{
				Optional:    true,
				Description: "The notes for the deploy request.",
			},
			"branch": schema.StringAttribute{
				Required:    true,
				Description: "The name of the branch to start the deploy request onto.",
			},
			"into_branch": schema.StringAttribute{
				Required:    true,
				Description: "The name of the branch to merge the deploy request into.",
			},
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The ID of the deploy request.",
			},
			"number": schema.Int64Attribute{
				Computed:    true,
				Description: "The number of the deploy request.",
			},
			"state": schema.StringAttribute{
				Computed:    true,
				Description: "The state of the deploy request.",
			},
			"deployment_state": schema.StringAttribute{
				Computed:    true,
				Description: "The deployment state of the deploy request.",
			},
			"approved": schema.BoolAttribute{
				Computed:    true,
				Description: "Whether the deploy request has been approved.",
			},
			"html_url": schema.StringAttribute{
				Computed:    true,
				Description: "The URL of the deploy request.",
			},
			"created_at": schema.StringAttribute{
				Computed:    true,
				Description: "The time the deploy request was created.",
			},
			"updated_at": schema.StringAttribute{
				Computed:    true,
				Description: "The time the deploy request was last updated.",
			},
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *deployRequestResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan deployRequestModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// create resource on Planetscale
	deployRequest, err := r.client.DeployRequests.Create(ctx, &planetscale.CreateDeployRequestRequest{
		Organization: plan.Organization.ValueString(),
		Database:     plan.Database.ValueString(),
		Branch:       plan.Branch.ValueString(),
		IntoBranch:   plan.IntoBranch.ValueString(),
		Notes:        plan.Notes.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating deploy request",
			"Could not create deploy request, unexpected error: "+err.Error()+". Make sure you have the correct "+
				"permissions to create a deploy request for this database and branch combination.",
		)
		return
	}

	plan.ID = types.StringValue(deployRequest.ID)
	plan.Number = types.Int64Value(int64(deployRequest.Number))
	plan.State = types.StringValue(deployRequest.State)
	plan.DeploymentState = types.StringValue(deployRequest.DeploymentState)
	plan.Approved = types.BoolValue(deployRequest.Approved)
	plan.HTMLURL = types.StringValue(deployRequest.HtmlURL)
	plan.CreatedAt = types.StringValue(deployRequest.CreatedAt.String())
	plan.UpdatedAt = types.StringValue(deployRequest.UpdatedAt.String())

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *deployRequestResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state deployRequestModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed deploy request info from Planetscale
	deployRequest, err := r.client.DeployRequests.Get(ctx, &planetscale.GetDeployRequestRequest{
		Organization: state.Organization.ValueString(),
		Database:     state.Database.ValueString(),
		Number:       uint64(state.Number.ValueInt64()),
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading deploy request",
			"Could not read info about deploy request for database "+state.Database.ValueString()+" and"+
				" branch "+state.Branch.ValueString()+": "+err.Error(),
		)
		return
	}

	// Overwrite items with refreshed state
	state.State = types.StringValue(deployRequest.State)
	state.DeploymentState = types.StringValue(deployRequest.DeploymentState)
	state.Approved = types.BoolValue(deployRequest.Approved)
	state.UpdatedAt = types.StringValue(deployRequest.UpdatedAt.String())

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *deployRequestResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// do nothing as a deploy-request should not be updated
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *deployRequestResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state deployRequestModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing deploy request
	ctx = tflog.SetField(ctx, "organization", state.Organization.ValueString())
	ctx = tflog.SetField(ctx, "database", state.Database.ValueString())
	ctx = tflog.SetField(ctx, "branch", state.Branch.ValueString())
	ctx = tflog.SetField(ctx, "into_branch", state.IntoBranch.ValueString())
	ctx = tflog.SetField(ctx, "number", state.Number.ValueInt64())

	// A deploy-request cannot be deleted. It can only be closed or reverted. Since in Terraform we might remove a
	// deploy-request resource from the configuration, we will close the deploy request instead of deleting it.
	_, err := r.client.DeployRequests.CloseDeploy(ctx, &planetscale.CloseDeployRequestRequest{
		Organization: state.Organization.ValueString(),
		Database:     state.Database.ValueString(),
		Number:       uint64(state.Number.ValueInt64()),
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting Planetscale database",
			"Could not delete database, unexpected error: "+err.Error(),
		)
		return
	}
	tflog.Debug(ctx, "closes deploy request")
}

// Configure adds the provider configured client to the resource.
func (r *deployRequestResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*planetscale.Client)
}
