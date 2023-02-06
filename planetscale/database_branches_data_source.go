package planetscale

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/planetscale/planetscale-go/planetscale"
)

// databaseBranchesDataSource maps the data source schema data.
type databaseBranchesDataSourceModel struct {
	Organization     types.String            `tfsdk:"organization"`
	Database         types.String            `tfsdk:"database"`
	DatabaseBranches []databaseBranchesModel `tfsdk:"database_branches"`
}

// databasesModel maps database-branch schema data.
type databaseBranchesModel struct {
	Name          types.String `tfsdk:"name"`
	ParentBranch  types.String `tfsdk:"parent_branch"`
	HtmlURL       types.String `tfsdk:"html_url"`
	Region        types.String `tfsdk:"region"`
	AccessHostURL types.String `tfsdk:"access_host_url"`
	Production    types.Bool   `tfsdk:"production"`
	Ready         types.Bool   `tfsdk:"ready"`
}

func NewDatabaseBranchesDataSource() datasource.DataSource {
	return &databaseBranchesDataSource{}
}

type databaseBranchesDataSource struct {
	client *planetscale.Client
}

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &databaseBranchesDataSource{}
	_ datasource.DataSourceWithConfigure = &databaseBranchesDataSource{}
)

func (d *databaseBranchesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_database_branches"
}

// Schema defines the schema for the data source.
func (d *databaseBranchesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"organization": schema.StringAttribute{Required: true},
			"database":     schema.StringAttribute{Required: true},
			"database_branches": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Computed: true,
						},
						"parent_branch": schema.StringAttribute{
							Computed: true,
						},
						"html_url": schema.StringAttribute{
							Computed: true,
						},
						"region": schema.StringAttribute{
							Computed: true,
						},
						"access_host_url": schema.StringAttribute{
							Computed: true,
						},
						"production": schema.BoolAttribute{
							Computed: true,
						},
						"ready": schema.BoolAttribute{
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func (d *databaseBranchesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state databaseBranchesDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	ctx = tflog.SetField(ctx, "database", state.Database)

	tflog.Info(ctx, "requesting database branches listing from Planetscale")
	databaseBranches, err := d.client.DatabaseBranches.List(ctx, &planetscale.ListDatabaseBranchesRequest{
		Organization: state.Organization.ValueString(),
		Database:     state.Database.ValueString(),
	},
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"unable to read Planetscale database branches",
			err.Error(),
		)
		return
	}

	for _, databaseBranch := range databaseBranches {
		databaseBranchState := databaseBranchesModel{
			Name:          types.StringValue(databaseBranch.Name),
			ParentBranch:  types.StringValue(databaseBranch.ParentBranch),
			Region:        types.StringValue(databaseBranch.Region.Name),
			HtmlURL:       types.StringValue(databaseBranch.HtmlURL),
			AccessHostURL: types.StringValue(databaseBranch.AccessHostURL),
			Production:    types.BoolValue(databaseBranch.Production),
			Ready:         types.BoolValue(databaseBranch.Ready),
		}
		state.DatabaseBranches = append(state.DatabaseBranches, databaseBranchState)
	}

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Configure adds the provider configured client to the data source.
func (d *databaseBranchesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*planetscale.Client)
}
