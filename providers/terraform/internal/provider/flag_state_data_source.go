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

var _ datasource.DataSource = &FlagStateDataSource{}
var _ datasource.DataSourceWithConfigure = &FlagStateDataSource{}

type FlagStateDataSource struct {
	client *client.Client
}

type FlagStateDataSourceModel struct {
	ID            types.String `tfsdk:"id"`
	EnvironmentID types.String `tfsdk:"environment_id"`
	FeatureFlagID types.String `tfsdk:"feature_flag_id"`
	Enabled       types.Bool   `tfsdk:"enabled"`
	Rules         types.String `tfsdk:"rules"`
	DefaultVariant types.String `tfsdk:"default_variant"`
}

func NewFlagStateDataSource() datasource.DataSource {
	return &FlagStateDataSource{}
}

func (d *FlagStateDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_flag_state"
}

func (d *FlagStateDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Data source for FlagManagment Flag State.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Required:    true,
				Description: "Unique identifier of the flag state.",
			},
			"environment_id": schema.StringAttribute{
				Required:    true,
				Description: "The ID of the environment.",
			},
			"feature_flag_id": schema.StringAttribute{
				Required:    true,
				Description: "The ID of the feature flag.",
			},
			"enabled": schema.BoolAttribute{
				Computed:    true,
				Description: "Whether the flag is enabled.",
			},
			"rules": schema.StringAttribute{
				Computed:    true,
				Description: "JSON representation of targeting rules.",
			},
			"default_variant": schema.StringAttribute{
				Computed:    true,
				Description: "Default variant when no rules match.",
			},
		},
	}
}

func (d *FlagStateDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *FlagStateDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state FlagStateDataSourceModel

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	fs, err := d.client.GetFlagState(ctx, state.EnvironmentID.ValueString(), state.FeatureFlagID.ValueString(), state.ID.ValueString())
	if err != nil {
		if errors.Is(err, client.ErrNotFound) {
			resp.Diagnostics.AddError("Flag State not found", "Flag state not found")
			return
		}
		resp.Diagnostics.AddError("Error reading flag state", err.Error())
		return
	}

	state.Enabled = types.BoolValue(fs.Enabled)
	state.Rules = types.StringValue(fs.Rules)
	state.DefaultVariant = types.StringValue(fs.DefaultVariant)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
