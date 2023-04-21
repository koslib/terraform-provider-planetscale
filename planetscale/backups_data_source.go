package planetscale

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/planetscale/planetscale-go/planetscale"
)

// backupsDataSourceModel maps the data source schema data.
type backupsDataSourceModel struct {
	Organization types.String   `tfsdk:"organization"`
	Database     types.String   `tfsdk:"database"`
	Branch       types.String   `tfsdk:"branch"`
	Backups      []backupsModel `tfsdk:"backups"`
}

// backupsModel maps database backup schema data.
type backupsModel struct {
	Name        types.String `tfsdk:"name"`
	PublicID    types.String `tfsdk:"public_id"`
	State       types.String `tfsdk:"state"`
	StartedAt   types.String `tfsdk:"started_at"`
	UpdatedAt   types.String `tfsdk:"updated_at"`
	CompletedAt types.String `tfsdk:"completed_at"`
	ExpiresAt   types.String `tfsdk:"expires_at"`
	Size        types.Int64  `tfsdk:"size"`
}

func NewBackupsDataSource() datasource.DataSource {
	return &backupsDataSource{}
}

type backupsDataSource struct {
	client *planetscale.Client
}

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &backupsDataSource{}
	_ datasource.DataSourceWithConfigure = &backupsDataSource{}
)

func (d *backupsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_backups"
}

// Schema defines the schema for the data source.
func (d *backupsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "The backups data source provides information about the backups of a database/branch. For more" +
			" information, see the official documentation on Backups and Restore here:" +
			" https://planetscale.com/docs/concepts/back-up-and-restore",
		Attributes: map[string]schema.Attribute{
			"organization": schema.StringAttribute{
				Required:    true,
				Description: "The name of the organization that the backups belong to.",
			},
			"database": schema.StringAttribute{
				Required:    true,
				Description: "The name of the database to get the backups for.",
			},
			"branch": schema.StringAttribute{
				Required:    true,
				Description: "The name of the branch to get the backups for.",
			},
			"backups": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
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
							Description: "The public ID of the backup.",
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
				},
			},
		},
	}
}

func (d *backupsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state backupsDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	ctx = tflog.SetField(ctx, "organization", state.Organization)
	ctx = tflog.SetField(ctx, "database", state.Database)
	ctx = tflog.SetField(ctx, "branch", state.Branch)

	tflog.Info(ctx, "requesting backups listing from Planetscale")
	backups, err := d.client.Backups.List(ctx, &planetscale.ListBackupsRequest{
		Organization: state.Organization.ValueString(),
		Database:     state.Database.ValueString(),
		Branch:       state.Branch.ValueString(),
	},
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to read Planetscale backups. Make sure the database or branch that you requested backups "+
				"for exist and you have access to it.",
			err.Error(),
		)
		return
	}

	for _, backup := range backups {
		backupState := backupsModel{
			Name:        types.StringValue(backup.Name),
			PublicID:    types.StringValue(backup.PublicID),
			State:       types.StringValue(backup.State),
			CompletedAt: types.StringValue(backup.CompletedAt.String()),
			ExpiresAt:   types.StringValue(backup.ExpiresAt.String()),
			Size:        types.Int64Value(backup.Size),
			StartedAt:   types.StringValue(backup.StartedAt.String()),
			UpdatedAt:   types.StringValue(backup.UpdatedAt.String()),
		}
		state.Backups = append(state.Backups, backupState)
	}

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Configure adds the provider configured client to the data source.
func (d *backupsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*planetscale.Client)
}
