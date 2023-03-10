package planetscale

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/planetscale/planetscale-go/planetscale"
)

// coffeesDataSourceModel maps the data source schema data.
type regionsDataSourceModel struct {
	Organization types.String          `tfsdk:"organization"`
	Regions      []databaseRegionModel `tfsdk:"regions"`
}

func NewRegionsDataSource() datasource.DataSource {
	return &regionsDataSource{}
}

type regionsDataSource struct {
	client *planetscale.Client
}

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &regionsDataSource{}
	_ datasource.DataSourceWithConfigure = &regionsDataSource{}
)

func (d *regionsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_regions"
}

// Schema defines the schema for the data source.
func (d *regionsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "List of regions. This data source is used for listing regions enabled for your organization.",
		Attributes: map[string]schema.Attribute{
			"organization": schema.StringAttribute{
				Required:    true,
				Description: "The name of the organization to list regions for.",
			},
			"regions": schema.ListNestedAttribute{
				Description: "List of regions",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "The name of the region.",
						},
						"slug": schema.StringAttribute{
							Computed:    true,
							Description: "The slug of the region.",
						},
						"location": schema.StringAttribute{
							Computed:    true,
							Description: "The location of the region in humanized form.",
						},
						"enabled": schema.BoolAttribute{
							Computed:    true,
							Description: "Whether the region is enabled or not.",
						},
					},
				},
			},
		},
	}
}

func (d *regionsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state regionsDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	regions, err := d.client.Organizations.ListRegions(ctx, &planetscale.ListOrganizationRegionsRequest{
		Organization: state.Organization.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to read Planetscale regions. Make sure that the organization is correct and that you have"+
				" the correct permissions.",
			err.Error(),
		)
		return
	}

	for _, region := range regions {
		regionObject := databaseRegionModel{
			Name:     types.StringValue(region.Name),
			Slug:     types.StringValue(region.Slug),
			Location: types.StringValue(region.Location),
			Enabled:  types.BoolValue(region.Enabled),
		}
		state.Regions = append(state.Regions, regionObject)
	}

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Configure adds the provider configured client to the data source.
func (d *regionsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*planetscale.Client)
}
