package planetscale

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/planetscale/planetscale-go/planetscale"
)

// coffeesDataSourceModel maps the data source schema data.
type regionsDataSourceModel struct {
	Regions []databaseRegionModel `tfsdk:"regions"`
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
		Attributes: map[string]schema.Attribute{
			"regions": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Computed: true,
						},
						"slug": schema.StringAttribute{
							Computed: true,
						},
						"location": schema.StringAttribute{
							Computed: true,
						},
						"enabled": schema.BoolAttribute{
							Computed: true,
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

	tflog.Info(ctx, "requesting regions listing from Planetscale")
	regions, err := d.client.Organizations.ListRegions(ctx, &planetscale.ListOrganizationRegionsRequest{Organization: "koslib"})
	if err != nil {
		resp.Diagnostics.AddError(
			"unable to read Planetscale regions",
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
