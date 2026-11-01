package provider

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/malah-code/flagmanagment/providers/terraform/internal/client"
)

var _ datasource.DataSource = &EnvironmentDataSource{}
var _ datasource.DataSourceWithConfigure = &EnvironmentDataSource{}

type EnvironmentDataSource struct {
	client *client.Client
}

type EnvironmentDataSourceModel struct {
	ID          types.String `tfsdk:"id"`
	ProjectID   types.String `tfsdk:"project_id"`
	Name        types.String `tfsdk:"name"`
	IsProtected types.Bool   `tfsdk:"is_protected"`
}

func NewEnvironmentDataSource() datasource.DataSource {
	return &EnvironmentDataSource{}
}

func (d *EnvironmentDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_environment"
}

func (d *EnvironmentDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Data source for FlagManagment Environment.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Required:    true,
				Description: "Unique identifier of the environment.",
			},
			"project_id": schema.StringAttribute{
				Required:    true,
				Description: "The ID of the project.",
			},
			"name": schema.StringAttribute{
				Computed:    true,
				Description: "Human-readable name of the environment.",
			},
			"is_protected": schema.BoolAttribute{
				Computed:    true,
				Description: "Whether modifications generate change requests.",
			},
		},
	}
}

func (d *EnvironmentDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	c, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError("Unexpected Data Source Configure Type", fmt.Sprintf("Expected *client.Client, got: %T", req.ProviderData))
		return
	}

	d.client = c
}

func (d *EnvironmentDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state EnvironmentDataSourceModel

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	env, err := d.client.GetEnvironment(ctx, state.ID.ValueString())
	if err != nil {
		if errors.Is(err, client.ErrNotFound) {
			resp.Diagnostics.AddError("Environment not found", fmt.Sprintf("Environment ID %s not found in Project %s", state.ID.ValueString(), state.ProjectID.ValueString()))
			return
		}
		resp.Diagnostics.AddError("Error reading environment", err.Error())
		return
	}

	state.Name = types.StringValue(env.Name)
	state.IsProtected = types.BoolValue(env.IsProtected)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
