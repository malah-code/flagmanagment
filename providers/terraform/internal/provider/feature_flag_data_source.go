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

var _ datasource.DataSource = &FeatureFlagDataSource{}
var _ datasource.DataSourceWithConfigure = &FeatureFlagDataSource{}

type FeatureFlagDataSource struct {
	client *client.Client
}

type FeatureFlagDataSourceModel struct {
	ID          types.String `tfsdk:"id"`
	ProjectID   types.String `tfsdk:"project_id"`
	Key         types.String `tfsdk:"key"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Type        types.String `tfsdk:"type"`
}

func NewFeatureFlagDataSource() datasource.DataSource {
	return &FeatureFlagDataSource{}
}

func (d *FeatureFlagDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_feature_flag"
}

func (d *FeatureFlagDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Data source for FlagManagment Feature Flag.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Required:    true,
				Description: "Unique identifier of the feature flag.",
			},
			"project_id": schema.StringAttribute{
				Required:    true,
				Description: "The ID of the project.",
			},
			"key": schema.StringAttribute{
				Computed:    true,
				Description: "Unique string key for the feature flag.",
			},
			"name": schema.StringAttribute{
				Computed:    true,
				Description: "Human-readable name of the flag.",
			},
			"description": schema.StringAttribute{
				Computed:    true,
				Description: "Detailed description of the flag.",
			},
			"type": schema.StringAttribute{
				Computed:    true,
				Description: "Type of flag (boolean, string, number, json).",
			},
		},
	}
}

func (d *FeatureFlagDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *FeatureFlagDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state FeatureFlagDataSourceModel

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	flag, err := d.client.GetFeatureFlag(ctx, state.ProjectID.ValueString(), state.ID.ValueString())
	if err != nil {
		if errors.Is(err, client.ErrNotFound) {
			resp.Diagnostics.AddError("Feature Flag not found", fmt.Sprintf("Flag ID %s not found in Project %s", state.ID.ValueString(), state.ProjectID.ValueString()))
			return
		}
		resp.Diagnostics.AddError("Error reading feature flag", err.Error())
		return
	}

	state.Key = types.StringValue(flag.Key)
	state.Name = types.StringValue(flag.Name)
	state.Description = types.StringValue(flag.Description)
	state.Type = types.StringValue(flag.Type)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
