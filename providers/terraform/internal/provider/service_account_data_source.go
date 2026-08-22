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

var _ datasource.DataSource = &ServiceAccountDataSource{}
var _ datasource.DataSourceWithConfigure = &ServiceAccountDataSource{}

type ServiceAccountDataSource struct {
	client *client.Client
}

type ServiceAccountDataSourceModel struct {
	ID        types.String `tfsdk:"id"`
	ProjectID types.String `tfsdk:"project_id"`
	Name      types.String `tfsdk:"name"`
	Role      types.String `tfsdk:"role"`
}

func NewServiceAccountDataSource() datasource.DataSource {
	return &ServiceAccountDataSource{}
}

func (d *ServiceAccountDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_service_account"
}

func (d *ServiceAccountDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Data source for FlagManagment Service Account.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Required:    true,
				Description: "Unique identifier of the service account.",
			},
			"project_id": schema.StringAttribute{
				Required:    true,
				Description: "The ID of the project.",
			},
			"name": schema.StringAttribute{
				Computed:    true,
				Description: "Name of the service account.",
			},
			"role": schema.StringAttribute{
				Computed:    true,
				Description: "Role of the service account.",
			},
		},
	}
}

func (d *ServiceAccountDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *ServiceAccountDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state ServiceAccountDataSourceModel

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	sa, err := d.client.GetServiceAccount(ctx, state.ProjectID.ValueString(), state.ID.ValueString())
	if err != nil {
		if errors.Is(err, client.ErrNotFound) {
			resp.Diagnostics.AddError("Service Account not found", "Service account not found")
			return
		}
		resp.Diagnostics.AddError("Error reading service account", err.Error())
		return
	}

	state.Name = types.StringValue(sa.Name)
	state.Role = types.StringValue(sa.Role)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
