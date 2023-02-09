package planetscale

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/planetscale/planetscale-go/planetscale"
)

// databaseBranchPasswordDataSourceModel maps the data source schema data.
type databaseBranchPasswordDataSourceModel struct {
	Organization types.String                  `tfsdk:"organization"`
	Database     types.String                  `tfsdk:"database"`
	Branch       types.String                  `tfsdk:"branch"`
	Passwords    []databaseBranchPasswordModel `tfsdk:"passwords"`
}

// databaseBranchPasswordModel maps the database branch password model schema data.
type databaseBranchPasswordModel struct {
	Name      types.String `tfsdk:"name"`
	Username  types.String `tfsdk:"username"`
	Hostname  types.String `tfsdk:"hostname"`
	PlainText types.String `tfsdk:"plaintext"`
	PublicID  types.String `tfsdk:"public_id"`
	Role      types.String `tfsdk:"role"`
}

func NewDatabaseBranchPasswordDataSource() datasource.DataSource {
	return &databaseBranchPasswordDataSource{}
}

type databaseBranchPasswordDataSource struct {
	client *planetscale.Client
}

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &databaseBranchPasswordDataSource{}
	_ datasource.DataSourceWithConfigure = &databaseBranchPasswordDataSource{}
)

func (d *databaseBranchPasswordDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_database_branch_passwords"
}

// Schema defines the schema for the data source.
func (d *databaseBranchPasswordDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"organization": schema.StringAttribute{
				Optional:    true,
				Description: "The name of the organization that the database belongs to. ",
			},
			"database": schema.StringAttribute{
				Optional:    true,
				Description: "The name of the database that the branch belongs to.",
			},
			"branch": schema.StringAttribute{
				Optional:    true,
				Description: "The name of the branch that the passwords belong to.",
			},
			"passwords": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "The name of the password.",
						},
						"username": schema.StringAttribute{
							Computed:    true,
							Description: "The username for this password.",
						},
						"hostname": schema.StringAttribute{
							Computed:    true,
							Description: "The hostname for this password.",
						},
						"plaintext": schema.StringAttribute{
							Computed:    true,
							Sensitive:   true,
							Description: "The plaintext password for this password. This value is sensitive and will not be stored in the state file.",
						},
						"public_id": schema.StringAttribute{
							Computed:    true,
							Description: "The public ID for this password.",
						},
						"role": schema.StringAttribute{
							Computed:    true,
							Description: "The role for this password.",
						},
					},
				},
			},
		},
	}
}

func (d *databaseBranchPasswordDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state databaseBranchPasswordDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	ctx = tflog.SetField(ctx, "organization", state.Organization)
	ctx = tflog.SetField(ctx, "database", state.Database)
	ctx = tflog.SetField(ctx, "branch", state.Branch)

	// todo: add validation for required fields
	// todo: add validation for organization, database, and branch existence

	tflog.Info(ctx, "requesting database branch passwords listing from Planetscale")
	databaseBranchPasswords, err := d.client.Passwords.List(ctx, &planetscale.ListDatabaseBranchPasswordRequest{
		Organization: state.Organization.ValueString(),
		Database:     state.Database.ValueString(),
		Branch:       state.Branch.ValueString(),
	},
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"unable to read Planetscale database branch passwords",
			err.Error(),
		)
		return
	}

	for _, databaseBranchPassword := range databaseBranchPasswords {
		databaseBranchPasswordsState := databaseBranchPasswordModel{
			Name:      types.StringValue(databaseBranchPassword.Name),
			Username:  types.StringValue(databaseBranchPassword.Username),
			Hostname:  types.StringValue(databaseBranchPassword.Hostname),
			PlainText: types.StringValue(databaseBranchPassword.PlainText),
			PublicID:  types.StringValue(databaseBranchPassword.PublicID),
			Role:      types.StringValue(databaseBranchPassword.Role),
		}
		state.Passwords = append(state.Passwords, databaseBranchPasswordsState)
	}

	tflog.Debug(ctx, "returning database branch passwords listing from Planetscale")

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Configure adds the provider configured client to the data source.
func (d *databaseBranchPasswordDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*planetscale.Client)
}
