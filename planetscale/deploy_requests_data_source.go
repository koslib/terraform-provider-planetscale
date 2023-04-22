package planetscale

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/planetscale/planetscale-go/planetscale"
)

// deployRequestsDataSourceModel maps the data source schema data.
type deployRequestsDataSourceModel struct {
	Organization  types.String       `tfsdk:"organization"`
	Database      types.String       `tfsdk:"database"`
	Number        types.String       `tfsdk:"number"`
	DeployRequest deployRequestModel `tfsdk:"deploy_request"`
}

func NewDeployRequestsDataSource() datasource.DataSource {
	return &deployRequestsDataSource{}
}

type deployRequestsDataSource struct {
	client *planetscale.Client
}

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &deployRequestsDataSource{}
	_ datasource.DataSourceWithConfigure = &deployRequestsDataSource{}
)

func (d *deployRequestsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_deploy_requests"
}

// Schema defines the schema for the data source.
func (d *deployRequestsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "List of deploy requests for the given database.",
		Attributes: map[string]schema.Attribute{
			"organization": schema.StringAttribute{
				Required:    true,
				Description: "The name of the organization the database belongs to.",
			},
			"database": schema.StringAttribute{
				Required:    true,
				Description: "The name of the database to list deploy requests for.",
			},
			"number": schema.StringAttribute{
				Required:    true,
				Description: "The number of the deploy request to get.",
			},
			"deploy_request": schema.SingleNestedAttribute{
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Computed:    true,
						Description: "The ID of the deploy request.",
					},
					"number": schema.Int64Attribute{
						Computed:    true,
						Description: "The number of the deploy request.",
					},
					"branch": schema.StringAttribute{
						Computed:    true,
						Description: "The name of the branch to start the deploy request onto.",
					},
					"into_branch": schema.StringAttribute{
						Computed:    true,
						Description: "The name of the branch to merge the deploy request into.",
					},
					"notes": schema.StringAttribute{
						Computed:    true,
						Description: "The notes for the deploy request.",
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
			},
		},
	}
}

func (d *deployRequestsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state deployRequestsDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	ctx = tflog.SetField(ctx, "organization", state.Organization.ValueString())
	ctx = tflog.SetField(ctx, "database", state.Database.ValueString())
	ctx = tflog.SetField(ctx, "number", state.Number.ValueString())

	tflog.Info(ctx, "requesting deploy request info")

	deployRequest, err := d.client.DeployRequests.Get(ctx, &planetscale.GetDeployRequestRequest{
		Organization: state.Organization.ValueString(),
		Database:     state.Database.ValueString(),
		Number:       0,
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to read deploy requests. Make sure your credentials are correct and you have access "+
				"to the organization, or that you have the correct permissions for this database.",
			err.Error(),
		)
		return
	}

	state.DeployRequest.ID = types.StringValue(deployRequest.ID)
	state.DeployRequest.Number = types.Int64Value(int64(deployRequest.Number))
	state.DeployRequest.Branch = types.StringValue(deployRequest.Branch)
	state.DeployRequest.IntoBranch = types.StringValue(deployRequest.IntoBranch)
	state.DeployRequest.Notes = types.StringValue(deployRequest.Notes)
	state.DeployRequest.State = types.StringValue(deployRequest.State)
	state.DeployRequest.DeploymentState = types.StringValue(deployRequest.DeploymentState)
	state.DeployRequest.Approved = types.BoolValue(deployRequest.Approved)
	state.DeployRequest.HTMLURL = types.StringValue(deployRequest.HtmlURL)
	state.DeployRequest.CreatedAt = types.StringValue(deployRequest.CreatedAt.String())
	state.DeployRequest.UpdatedAt = types.StringValue(deployRequest.UpdatedAt.String())

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Configure adds the provider configured client to the data source.
func (d *deployRequestsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*planetscale.Client)
}
