package planetscale

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/planetscale/planetscale-go/planetscale"
)

var (
	_ resource.Resource              = &backupResource{}
	_ resource.ResourceWithConfigure = &backupResource{}
)

type backupResourceModel struct {
	Organization types.String `tfsdk:"organization"`
	Database     types.String `tfsdk:"database"`
	Branch       types.String `tfsdk:"branch"`
	PublicID     types.String `tfsdk:"public_id"`
	Name         types.String `tfsdk:"name"`
	State        types.String `tfsdk:"state"`
	Size         types.Int64  `tfsdk:"size"`
	CreatedAt    types.String `tfsdk:"created_at"`
	UpdatedAt    types.String `tfsdk:"updated_at"`
	StartedAt    types.String `tfsdk:"started_at"`
	ExpiresAt    types.String `tfsdk:"expires_at"`
	CompletedAt  types.String `tfsdk:"completed_at"`
}

// NewBackupResource is a helper function to simplify the provider implementation.
func NewBackupResource() resource.Resource {
	return &backupResource{}
}

// backupResource is the resource implementation.
type backupResource struct {
	client *planetscale.Client
}

// Metadata returns the resource type name.
func (r *backupResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_backup"
}

// Schema defines the schema for the resource.
func (r *backupResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "A Planetscale backup. This resource will create a new backup for a database in your Planetscale organization." +
			" The backup will be created for the specified branch.",
		Attributes: map[string]schema.Attribute{
			"organization": schema.StringAttribute{
				Required:    true,
				Description: "The organization where the backup will be created as well as the database/branch belong to.",
			},
			"database": schema.StringAttribute{
				Required:    true,
				Description: "The name of the database to create the backup for.",
			},
			"branch": schema.StringAttribute{
				Required:    true,
				Description: "The name of the branch to create the backup for.",
			},
			"name": schema.StringAttribute{
				Computed:    true,
				Description: "The name of the backup.",
			},
			"state": schema.StringAttribute{
				Computed:    true,
				Description: "The state of the backup. Options are: 'pending', 'running', 'success', 'failed', 'canceled'.",
			},
			"public_id": schema.StringAttribute{
				Computed:    true,
				Optional:    true,
				Description: "The public ID of the backup.",
			},
			"created_at": schema.StringAttribute{
				Computed:    true,
				Description: "The timestamp of when the backup object was created.",
			},
			"started_at": schema.StringAttribute{
				Computed:    true,
				Description: "The timestamp of when the backup started.",
			},
			"updated_at": schema.StringAttribute{
				Computed:    true,
				Description: "Last update timestamp for the backup.",
			},
			"completed_at": schema.StringAttribute{
				Computed:    true,
				Description: "If the backup is completed, this is the timestamp of when it was completed.",
			},
			"expires_at": schema.StringAttribute{
				Computed:    true,
				Description: "If the backup is completed, this is the timestamp of when it will expire.",
			},
			"size": schema.Int64Attribute{
				Computed:    true,
				Description: "The size of the backup.",
			},
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *backupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan backupResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// create resource on Planetscale
	backup, err := r.client.Backups.Create(ctx, &planetscale.CreateBackupRequest{
		Organization: plan.Organization.ValueString(),
		Database:     plan.Database.ValueString(),
		Branch:       plan.Branch.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating backup",
			"Could not create backup, unexpected error: "+err.Error()+". Make sure you have the correct "+
				"permissions to create a backup for this branch or in this organization.",
		)
		return
	}

	// todo: Name field is not supported by the API yet, therefore is only computed. Add support for it once
	//  the golang-sdk is updated.
	plan.Name = types.StringValue(backup.Name)
	plan.PublicID = types.StringValue(backup.PublicID)
	plan.State = types.StringValue(backup.State)
	plan.Size = types.Int64Value(backup.Size)
	plan.CreatedAt = types.StringValue(backup.CreatedAt.String())
	plan.UpdatedAt = types.StringValue(backup.UpdatedAt.String())
	plan.StartedAt = types.StringValue(backup.StartedAt.String())
	plan.ExpiresAt = types.StringValue(backup.ExpiresAt.String())
	plan.CompletedAt = types.StringValue(backup.CompletedAt.String())

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Backup created")
}

// Read refreshes the Terraform state with the latest data.
func (r *backupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state backupResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed backup object info from Planetscale
	backup, err := r.client.Backups.Get(ctx, &planetscale.GetBackupRequest{
		Organization: state.Organization.ValueString(),
		Database:     state.Database.ValueString(),
		Branch:       state.Branch.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading Planetscale backup",
			"Could not read info about Planetscale database "+state.Database.ValueString()+" on branch"+
				""+state.Branch.ValueString()+": "+err.Error(),
		)
		return
	}

	state.State = types.StringValue(backup.State)
	state.Size = types.Int64Value(backup.Size)
	state.UpdatedAt = types.StringValue(backup.UpdatedAt.String())
	state.CompletedAt = types.StringValue(backup.CompletedAt.String())

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *backupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// does not make sense to update a backup, therefore skip
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *backupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state backupResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing Planetscale backup
	ctx = tflog.SetField(ctx, "organization", state.Organization.ValueString())
	ctx = tflog.SetField(ctx, "database", state.Database.ValueString())
	ctx = tflog.SetField(ctx, "branch", state.Branch.ValueString())
	ctx = tflog.SetField(ctx, "backup", state.PublicID.ValueString())

	err := r.client.Backups.Delete(ctx, &planetscale.DeleteBackupRequest{
		Organization: state.Organization.ValueString(),
		Database:     state.Database.ValueString(),
		Branch:       state.Branch.ValueString(),
		Backup:       state.PublicID.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting Planetscale backup",
			"Could not delete backup for database, unexpected error: "+err.Error(),
		)
		tflog.Error(ctx, "could not delete Planetscale backup", map[string]interface{}{"error": err.Error()})
		return
	}

	tflog.Debug(ctx, "deleted Planetscale backup")
}

// Configure adds the provider configured client to the resource.
func (r *backupResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*planetscale.Client)
}
