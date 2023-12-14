package planetscale

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/planetscale/planetscale-go/planetscale"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

var (
	_ resource.Resource              = &databaseBranchResource{}
	_ resource.ResourceWithConfigure = &databaseBranchResource{}
)

type databaseBranchResourceModel struct {
	Name         types.String `tfsdk:"name"`
	Database     types.String `tfsdk:"database"`
	Organization types.String `tfsdk:"organization"`
	Region       types.String `tfsdk:"region"`
	ParentBranch types.String `tfsdk:"parent_branch"`
	BackupID     types.String `tfsdk:"backup_id"`
	SeedData     types.String `tfsdk:"seed_data"`
	HTMLURL      types.String `tfsdk:"html_url"`
	Production   types.Bool   `tfsdk:"production"`
	Ready        types.Bool   `tfsdk:"ready"`
}

// NewDatabaseBranchResource is a helper function to simplify the provider implementation.
func NewDatabaseBranchResource() resource.Resource {
	return &databaseBranchResource{}
}

// databaseResource is the resource implementation.
type databaseBranchResource struct {
	client *planetscale.Client
}

// Metadata returns the resource type name.
func (r *databaseBranchResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_database_branch"
}

// Schema defines the schema for the resource.
func (r *databaseBranchResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "A database branch resource. This resource allows you to create and manage database branches. " +
			"A database branch is a copy of a database that can be used for development, testing, or other purposes." +
			"For more information on database branches please see here: https://planetscale.com/docs/concepts/branching.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the database branch.",
			},
			"database": schema.StringAttribute{
				Required:    true,
				Description: "The name of the database to create the branch for.",
			},
			"organization": schema.StringAttribute{
				Required:    true,
				Description: "The name of the organization to create the database branch in.",
			},
			"region": schema.StringAttribute{
				Optional: true,
				Description: "The region where the database will be created. If not specified, the default region" +
					" for the organization will be used. Currently, following regions are supported: " +
					"ap-northeast, ap-south, ap-southeast, aws-ap-southeast-2, eu-central, eu-west, aws-eu-west-2, " +
					"aws-sa-east-1, us-east, aws-us-east-2, us-west, gcp-us-central1, gcp-us-east4, " +
					"gcp-northamerica-northeast1, gcp-asia-northeast3. For more information on regions, " +
					"please see here: https://planetscale.com/docs/concepts/regions.",
				Validators: []validator.String{
					stringvalidator.OneOf(
						"ap-northeast",
						"ap-south",
						"ap-southeast",
						"aws-ap-southeast-2",
						"eu-central",
						"eu-west",
						"aws-eu-west-2",
						"aws-sa-east-1",
						"us-east",
						"aws-us-east-2",
						"us-west",
						"gcp-us-central1",
						"gcp-us-east4",
						"gcp-northamerica-northeast1",
						"gcp-asia-northeast3",
					),
				},
			},
			"parent_branch": schema.StringAttribute{
				Optional:    true,
				Description: "The name of the parent branch to create the database branch from. If not specified, the database's default branch will be used.",
			},
			"backup_id": schema.StringAttribute{
				Optional:    true,
				Description: "The ID of the backup to create the database branch from. If not specified, the database's default branch will be used.",
			},
			"seed_data": schema.StringAttribute{
				Optional:    true,
				Description: "The name of the database branch to seed the new database branch with. If not specified, the database's default branch will be used.",
			},
			"html_url": schema.StringAttribute{
				Computed:    true,
				Description: "The URL to the database branch in the Planetscale UI.",
			},
			"production": schema.BoolAttribute{
				Computed:    true,
				Description: "Whether the database branch is a production branch.",
			},
			"ready": schema.BoolAttribute{
				Computed:    true,
				Description: "Whether the database branch is ready to be used.",
			},
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *databaseBranchResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan databaseBranchResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// create resource on Planetscale
	databaseBranch, err := r.client.DatabaseBranches.Create(ctx, &planetscale.CreateDatabaseBranchRequest{
		Organization: plan.Organization.ValueString(),
		Database:     plan.Database.ValueString(),
		Region:       plan.Region.ValueString(),
		Name:         plan.Name.ValueString(),
		ParentBranch: plan.ParentBranch.ValueString(),
		BackupID:     plan.BackupID.ValueString(),
		SeedData:     plan.SeedData.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating database branch",
			"Could not create database branch, unexpected error: "+err.Error()+". Make sure the database branch "+
				"name is unique and that you have the correct permissions.",
		)
		return
	}

	plan.HTMLURL = types.StringValue(databaseBranch.HtmlURL)
	plan.Production = types.BoolValue(databaseBranch.Production)
	plan.Ready = types.BoolValue(databaseBranch.Ready)

	tflog.Debug(ctx, "database branch created")

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *databaseBranchResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state databaseBranchResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed database branch value from Planetscale
	databaseBranch, err := r.client.DatabaseBranches.Get(ctx, &planetscale.GetDatabaseBranchRequest{
		Organization: state.Organization.ValueString(),
		Database:     state.Database.ValueString(),
		Branch:       state.Name.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading Planetscale database branch",
			"Could not read info about Planetscale database branch "+state.Name.ValueString()+": "+err.Error(),
		)
		return
	}

	// Overwrite items with refreshed state
	state.HTMLURL = types.StringValue(databaseBranch.HtmlURL)
	state.Production = types.BoolValue(databaseBranch.Production)
	state.Ready = types.BoolValue(databaseBranch.Ready)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *databaseBranchResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// not supported by the golang sdk (https://github.com/planetscale/planetscale-go) yet
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *databaseBranchResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state databaseBranchResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing Planetscale database branch
	ctx = tflog.SetField(ctx, "organization", state.Organization.ValueString())
	ctx = tflog.SetField(ctx, "database", state.Database.ValueString())
	ctx = tflog.SetField(ctx, "branch", state.Name.ValueString())

	err := r.client.DatabaseBranches.Delete(ctx, &planetscale.DeleteDatabaseBranchRequest{
		Organization: state.Organization.ValueString(),
		Database:     state.Database.ValueString(),
		Branch:       state.Name.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting Planetscale database branch",
			"Could not delete database branch, unexpected error: "+err.Error(),
		)
		return
	}
	tflog.Info(ctx, "deleted Planetscale database branch")
}

// Configure adds the provider configured client to the resource.
func (r *databaseBranchResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*planetscale.Client)
}

func splitDatabaseBranchResourceID(id string) (teamID, _id string, branchName string, ok bool) {
	attributes := strings.Split(id, "/")
	requiredAttributesLength := 3
	if len(attributes) != requiredAttributesLength {
		return "", "", "", false
	}
	return attributes[0], attributes[1], attributes[2], true
}

func (r *databaseBranchResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	organizationName, databaseName, branchName, ok := splitDatabaseBranchResourceID(req.ID)
	if !ok {
		resp.Diagnostics.AddError(
			"Error importing database",
			fmt.Sprintf("Invalid input '%s' provided. should be in format \"organization_name/database_name/branch_name\"", req.ID),
		)
		return
	}

	out, err := r.client.DatabaseBranches.Get(ctx, &planetscale.GetDatabaseBranchRequest{
		Organization: organizationName,
		Database:     databaseName,
		Branch:       branchName,
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading database branch",
			fmt.Sprintf("Could not get database branch %s %s %s, unexpected error: %v",
				organizationName,
				databaseName,
				branchName,
				err,
			),
		)
		return
	}

	tflog.Trace(ctx, "imported database branch", map[string]interface{}{
		organizationName: organizationName,
		databaseName:     databaseName,
		branchName:       branchName,
	})

	diags := resp.State.Set(ctx, &databaseBranchResourceModel{
		Name:         types.StringValue(out.Name),
		Organization: types.StringValue(organizationName),
		Region:       types.StringValue(out.Region.Slug),
		HTMLURL:      types.StringValue(out.HtmlURL),
		Database:     types.StringValue(databaseName),
		ParentBranch: types.StringValue(out.ParentBranch),
		BackupID:     types.StringUnknown(),
		SeedData:     types.StringUnknown(),
		Production:   types.BoolValue(out.Production),
		Ready:        types.BoolValue(out.Ready),
	})
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
