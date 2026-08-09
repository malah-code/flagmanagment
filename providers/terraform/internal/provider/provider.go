package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/malah-code/flagmanagment/providers/terraform/internal/client"
)

var _ provider.Provider = &FlagManagmentProvider{}

type FlagManagmentProvider struct {
	version string
}

type FlagManagmentProviderModel struct {
	APIURL               types.String `tfsdk:"api_url"`
	APIKey               types.String `tfsdk:"api_key"`
	BypassChangeRequests types.Bool   `tfsdk:"bypass_change_requests"`
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &FlagManagmentProvider{
			version: version,
		}
	}
}

func (p *FlagManagmentProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "flagmanagment"
	resp.Version = p.version
}

func (p *FlagManagmentProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Interact with FlagManagment feature flag and remote configuration server.",
		Attributes: map[string]schema.Attribute{
			"api_url": schema.StringAttribute{
				Optional:    true,
				Description: "The base URL for FlagManagment server. Defaults to http://localhost:8080 or FLAGMANAGMENT_API_URL.",
			},
			"api_key": schema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				Description: "The Service Account API Key or Bearer Token. Defaults to FLAGMANAGMENT_API_KEY environment variable.",
			},
			"bypass_change_requests": schema.BoolAttribute{
				Optional:    true,
				Description: "Set to true to bypass protected environment change request approvals if authorized. Defaults to false.",
			},
		},
	}
}

func (p *FlagManagmentProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config FlagManagmentProviderModel

	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiURL := config.APIURL.ValueString()
	apiKey := config.APIKey.ValueString()
	bypass := config.BypassChangeRequests.ValueBool()

	c, err := client.NewClient(apiURL, apiKey, bypass)
	if err != nil {
		resp.Diagnostics.AddError("Failed to initialize API client", err.Error())
		return
	}

	resp.DataSourceData = c
	resp.ResourceData = c
}

func (p *FlagManagmentProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewProjectResource,
		NewEnvironmentResource,
		NewFeatureFlagResource,
		NewFlagStateResource,
		NewServiceAccountResource,
	}
}

func (p *FlagManagmentProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}
