package planetscale

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/planetscale/planetscale-go/planetscale"
)

// databasesDataSourceModel maps the data source schema data.
type databasesDataSourceModel struct {
	Organization types.String     `tfsdk:"organization"`
	Databases    []databasesModel `tfsdk:"databases"`
}

type databaseRegionModel struct {
	Name     types.String `tfsdk:"name"`
	Slug     types.String `tfsdk:"slug"`
	Location types.String `tfsdk:"location"`
	Enabled  types.Bool   `tfsdk:"enabled"`
}

// databasesModel maps database schema data.
type databasesModel struct {
	Name    types.String        `tfsdk:"name"`
	Notes   types.String        `tfsdk:"notes"`
	Region  databaseRegionModel `tfsdk:"region"`
	HtmlURL types.String        `tfsdk:"html_url"`
	State   types.String        `tfsdk:"state"`
}

func NewDatabasesDataSource() datasource.DataSource {
	return &databasesDataSource{}
}

type databasesDataSource struct {
	client *planetscale.Client
}

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &databasesDataSource{}
	_ datasource.DataSourceWithConfigure = &databasesDataSource{}
)

func (d *databasesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_databases"
}

// Schema defines the schema for the data source.
func (d *databasesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "List of databases in the organization.",
		Attributes: map[string]schema.Attribute{
			"organization": schema.StringAttribute{
				Required:    true,
				Description: "The name of the organization to list databases for.",
			},
			"databases": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description: "Name of the database.",
							Computed:    true,
						},
						"notes": schema.StringAttribute{
							Description: "Notes data for the database.",
							Computed:    true,
						},
						"html_url": schema.StringAttribute{
							Description: "URL to the database in the web UI.",
							Computed:    true,
						},
						"state": schema.StringAttribute{
							Description: "State of the database.",
							Computed:    true,
						},
						"region": schema.SingleNestedAttribute{
							Required:    true,
							Description: "Region of the database.",
							Attributes: map[string]schema.Attribute{
								"name": schema.StringAttribute{
									Description: "Name of the region.",
									Computed:    true,
								},
								"slug": schema.StringAttribute{
									Description: "Slug of the region.",
									Computed:    true,
								},
								"location": schema.StringAttribute{
									Description: "Location of the region.",
									Computed:    true,
								},
								"enabled": schema.BoolAttribute{
									Description: "Whether the region is enabled.",
									Computed:    true,
								},
							},
						},
					},
				},
			},
		},
	}
}

func (d *databasesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state databasesDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	if state.Organization.ValueString() != "" {
		tflog.SetField(ctx, "organization", state.Organization.ValueString())
	}

	tflog.Info(ctx, "requesting database listing from Planetscale")
	databases, err := d.client.Databases.List(ctx, &planetscale.ListDatabasesRequest{
		Organization: state.Organization.ValueString()},
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"unable to read Planetscale databases",
			err.Error(),
		)
		return
	}

	for _, database := range databases {
		dbState := databasesModel{
			Name:  types.StringValue(database.Name),
			Notes: types.StringValue(database.Notes),
			Region: databaseRegionModel{
				Name:     types.StringValue(database.Region.Name),
				Slug:     types.StringValue(database.Region.Slug),
				Location: types.StringValue(database.Region.Location),
				Enabled:  types.BoolValue(database.Region.Enabled),
			},
			HtmlURL: types.StringValue(database.HtmlURL),
			State:   types.StringValue(string(database.State)),
		}
		state.Databases = append(state.Databases, dbState)
	}

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Configure adds the provider configured client to the data source.
func (d *databasesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*planetscale.Client)
}
