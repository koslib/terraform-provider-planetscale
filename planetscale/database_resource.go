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
	_ resource.Resource              = &databaseResource{}
	_ resource.ResourceWithConfigure = &databaseResource{}
)

type databaseResourceModel struct {
	Name         types.String `tfsdk:"name"`
	Notes        types.String `tfsdk:"notes"`
	Organization types.String `tfsdk:"organization"`
	Region       types.String `tfsdk:"region"`
	HTMLURL      types.String `tfsdk:"html_url"`
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
		Description: "A Planetscale database. This resource will create a new database in your Planetscale organization." +
			"Once created, you can manage the database using the Planetscale web UI or the Planetscale CLI. For more " +
			"information on Planetscale databases, please see here: https://planetscale.com/docs/concepts/planetscale-workflow.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the database. This must be unique within the organization.",
			},
			"notes": schema.StringAttribute{
				Optional:    true,
				Description: "Notes about the database. These are only visible to you and other members of the organization.",
			},
			"organization": schema.StringAttribute{
				Required:    true,
				Description: "The organization where the database will be created.",
			},
			"region": schema.StringAttribute{
				Optional: true,
				Description: "The region where the database will be created. If not specified, the default region" +
					" for the organization will be used. Values supported are: us-east, us-west, eu-west, ap-southeast," +
					" ap-south, ap-northeast, eu-central, aws-ap-southeast-2, aws-sa-east-1, aws-sa-east-2. For more information" +
					" on regions, please see here: https://planetscale.com/docs/reference/region.",
				Validators: []validator.String{
					stringvalidator.OneOf(
						"us-east",
						"us-west",
						"eu-west",
						"ap-southeast",
						"ap-south",
						"ap-northeast",
						"eu-central",
						"aws-ap-southeast-2",
						"aws-sa-east-1",
						"aws-sa-east-2",
					),
				},
			},
			"html_url": schema.StringAttribute{
				Computed:    true,
				Description: "The URL of the database in the Planetscale web UI.",
			},
			"state": schema.StringAttribute{
				Computed:    true,
				Description: "The state of the database. This will be one of the following: creating, ready, or error.",
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
			"Could not create database, unexpected error: "+err.Error()+". Make sure you have the correct "+
				"permissions to create a database in the organization.",
		)
		return
	}

	plan.HTMLURL = types.StringValue(database.HtmlURL)
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
	state.HTMLURL = types.StringValue(database.HtmlURL)
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
	// not supported by the golang sdk (https://github.com/planetscale/planetscale-go) yet
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
	tflog.Debug(ctx, "deleted Planetscale database")
}

// Configure adds the provider configured client to the resource.
func (r *databaseResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*planetscale.Client)
}

func splitDatabaseResourceID(id string) (teamID, _id string, ok bool) {
	attributes := strings.Split(id, "/")
	requiredAttributesLength := 2
	if len(attributes) != requiredAttributesLength {
		return "", "", false
	}
	return attributes[0], attributes[1], true
}

func (r *databaseResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	organizationName, databaseName, ok := splitDatabaseResourceID(req.ID)
	if !ok {
		resp.Diagnostics.AddError(
			"Error importing database",
			fmt.Sprintf("Invalid id '%s' specified. should be in format \"organization_name/database_name\"", req.ID),
		)
		return
	}

	out, err := r.client.Databases.Get(ctx, &planetscale.GetDatabaseRequest{
		Organization: organizationName,
		Database:     databaseName,
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading project",
			fmt.Sprintf("Could not get database %s %s, unexpected error: %s",
				organizationName,
				databaseName,
				err,
			),
		)
		return
	}

	tflog.Trace(ctx, "imported database", map[string]interface{}{
		organizationName: organizationName,
		databaseName:     databaseName,
	})

	diags := resp.State.Set(ctx, &databaseResourceModel{
		Name:         types.StringValue(out.Name),
		Notes:        types.StringValue(out.Notes),
		Organization: types.StringValue(organizationName),
		Region:       types.StringValue(out.Region.Slug),
		HTMLURL:      types.StringValue(out.HtmlURL),
		State:        types.StringValue(string(out.State)),
	})
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
