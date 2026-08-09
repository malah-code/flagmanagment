package provider

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/malah-code/flagmanagment/providers/terraform/internal/client"
)

var _ resource.Resource = &FeatureFlagResource{}
var _ resource.ResourceWithConfigure = &FeatureFlagResource{}

type FeatureFlagResource struct {
	client *client.Client
}

type FeatureFlagResourceModel struct {
	ID           types.String `tfsdk:"id"`
	ProjectID    types.String `tfsdk:"project_id"`
	Key          types.String `tfsdk:"key"`
	Name         types.String `tfsdk:"name"`
	Type         types.String `tfsdk:"type"`
	ParentFlagID types.String `tfsdk:"parent_flag_id"`
}

func NewFeatureFlagResource() resource.Resource {
	return &FeatureFlagResource{}
}

func (r *FeatureFlagResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_feature_flag"
}

func (r *FeatureFlagResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a FlagManagment Feature Flag.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"project_id": schema.StringAttribute{
				Required:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"key": schema.StringAttribute{
				Required:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"name": schema.StringAttribute{
				Required: true,
			},
			"type": schema.StringAttribute{
				Required:      true,
				Description:   "Type of flag: boolean or multivariate",
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"parent_flag_id": schema.StringAttribute{
				Optional:    true,
				Description: "Parent flag ID for sequential dependencies",
			},
		},
	}
}

func (r *FeatureFlagResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	c, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError("Unexpected Resource Configure Type", fmt.Sprintf("Expected *client.Client, got: %T", req.ProviderData))
		return
	}
	r.client = c
}

func (r *FeatureFlagResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan FeatureFlagResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	flag := client.FeatureFlag{
		ProjectID:    plan.ProjectID.ValueString(),
		Key:          plan.Key.ValueString(),
		Name:         plan.Name.ValueString(),
		Type:         plan.Type.ValueString(),
		ParentFlagID: plan.ParentFlagID.ValueString(),
	}

	created, err := r.client.CreateFeatureFlag(ctx, &flag)
	if err != nil {
		resp.Diagnostics.AddError("Error creating feature flag", err.Error())
		return
	}

	plan.ID = types.StringValue(created.ID)
	plan.Name = types.StringValue(created.Name)
	if created.ParentFlagID != "" {
		plan.ParentFlagID = types.StringValue(created.ParentFlagID)
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *FeatureFlagResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state FeatureFlagResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	flag, err := r.client.GetFeatureFlag(ctx, state.ID.ValueString())
	if err != nil {
		if errors.Is(err, client.ErrNotFound) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading feature flag", err.Error())
		return
	}

	state.Name = types.StringValue(flag.Name)
	state.Type = types.StringValue(flag.Type)
	if flag.ParentFlagID != "" {
		state.ParentFlagID = types.StringValue(flag.ParentFlagID)
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *FeatureFlagResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan FeatureFlagResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	flag := client.FeatureFlag{
		ID:           plan.ID.ValueString(),
		ProjectID:    plan.ProjectID.ValueString(),
		Key:          plan.Key.ValueString(),
		Name:         plan.Name.ValueString(),
		Type:         plan.Type.ValueString(),
		ParentFlagID: plan.ParentFlagID.ValueString(),
	}

	updated, err := r.client.UpdateFeatureFlag(ctx, &flag)
	if err != nil {
		resp.Diagnostics.AddError("Error updating feature flag", err.Error())
		return
	}

	plan.Name = types.StringValue(updated.Name)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *FeatureFlagResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state FeatureFlagResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteFeatureFlag(ctx, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting feature flag", err.Error())
		return
	}
}
