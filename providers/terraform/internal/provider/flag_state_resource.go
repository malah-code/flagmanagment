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

var _ resource.Resource = &FlagStateResource{}
var _ resource.ResourceWithConfigure = &FlagStateResource{}

type FlagStateResource struct {
	client *client.Client
}

type VariantModel struct {
	Name  types.String `tfsdk:"name"`
	Value types.String `tfsdk:"value"`
}

type ConditionModel struct {
	Attribute types.String   `tfsdk:"attribute"`
	Operator  types.String   `tfsdk:"operator"`
	Values    []types.String `tfsdk:"values"`
}

type RolloutModel struct {
	Percentages map[string]types.Int64 `tfsdk:"percentages"`
}

type TargetingRuleModel struct {
	Name       types.String     `tfsdk:"name"`
	Variant    types.String     `tfsdk:"variant"`
	Rollout    *RolloutModel    `tfsdk:"rollout"`
	Conditions []ConditionModel `tfsdk:"conditions"`
}

type TelemetryTriggerModel struct {
	MetricName types.String  `tfsdk:"metric_name"`
	Threshold  types.Float64 `tfsdk:"threshold"`
	Action     types.String  `tfsdk:"action"`
}

type FlagStateResourceModel struct {
	ID               types.String           `tfsdk:"id"`
	EnvironmentID    types.String           `tfsdk:"environment_id"`
	FlagID           types.String           `tfsdk:"flag_id"`
	Enabled          types.Bool             `tfsdk:"enabled"`
	DefaultVariant   types.String           `tfsdk:"default_variant"`
	Variants         []VariantModel         `tfsdk:"variants"`
	TargetingRules   []TargetingRuleModel   `tfsdk:"targeting_rules"`
	TelemetryTrigger *TelemetryTriggerModel `tfsdk:"telemetry_trigger"`
}

func NewFlagStateResource() resource.Resource {
	return &FlagStateResource{}
}

func (r *FlagStateResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_flag_state"
}

func (r *FlagStateResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a FlagState in a specific Environment.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"environment_id": schema.StringAttribute{
				Required:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"flag_id": schema.StringAttribute{
				Required:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"enabled": schema.BoolAttribute{
				Required: true,
			},
			"default_variant": schema.StringAttribute{
				Required: true,
			},
			"variants": schema.ListNestedAttribute{
				Optional: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Required: true,
						},
						"value": schema.StringAttribute{
							Required: true,
						},
					},
				},
			},
			"targeting_rules": schema.ListNestedAttribute{
				Optional: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Required: true,
						},
						"variant": schema.StringAttribute{
							Optional: true,
						},
						"rollout": schema.SingleNestedAttribute{
							Optional: true,
							Attributes: map[string]schema.Attribute{
								"percentages": schema.MapAttribute{
									ElementType: types.Int64Type,
									Required:    true,
								},
							},
						},
						"conditions": schema.ListNestedAttribute{
							Required: true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"attribute": schema.StringAttribute{
										Required: true,
									},
									"operator": schema.StringAttribute{
										Required: true,
									},
									"values": schema.ListAttribute{
										ElementType: types.StringType,
										Required:    true,
									},
								},
							},
						},
					},
				},
			},
			"telemetry_trigger": schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{
					"metric_name": schema.StringAttribute{
						Required: true,
					},
					"threshold": schema.Float64Attribute{
						Required: true,
					},
					"action": schema.StringAttribute{
						Required: true,
					},
				},
			},
		},
	}
}

func (r *FlagStateResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *FlagStateResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan FlagStateResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	fs := client.FlagState{
		EnvironmentID:  plan.EnvironmentID.ValueString(),
		FlagID:         plan.FlagID.ValueString(),
		Enabled:        plan.Enabled.ValueBool(),
		DefaultVariant: plan.DefaultVariant.ValueString(),
	}

	for _, v := range plan.Variants {
		fs.Variants = append(fs.Variants, client.Variant{
			Name:  v.Name.ValueString(),
			Value: v.Value.ValueString(),
		})
	}

	for _, tr := range plan.TargetingRules {
		rule := client.TargetingRule{
			Name:    tr.Name.ValueString(),
			Variant: tr.Variant.ValueString(),
		}
		for _, c := range tr.Conditions {
			cond := client.Condition{
				Attribute: c.Attribute.ValueString(),
				Operator:  c.Operator.ValueString(),
			}
			for _, val := range c.Values {
				cond.Values = append(cond.Values, val.ValueString())
			}
			rule.Conditions = append(rule.Conditions, cond)
		}
		if tr.Rollout != nil {
			rule.Rollout = &client.Rollout{
				Percentages: make(map[string]int),
			}
			for k, v := range tr.Rollout.Percentages {
				rule.Rollout.Percentages[k] = int(v.ValueInt64())
			}
		}
		fs.TargetingRules = append(fs.TargetingRules, rule)
	}

	if plan.TelemetryTrigger != nil {
		fs.TelemetryTrigger = &client.TelemetryTrigger{
			MetricName: plan.TelemetryTrigger.MetricName.ValueString(),
			Threshold:  plan.TelemetryTrigger.Threshold.ValueFloat64(),
			Action:     plan.TelemetryTrigger.Action.ValueString(),
		}
	}

	updated, err := r.client.UpdateFlagState(ctx, &fs)
	if err != nil {
		resp.Diagnostics.AddError("Error creating/updating flag state", err.Error())
		return
	}

	plan.ID = types.StringValue(fmt.Sprintf("%s:%s", updated.EnvironmentID, updated.FlagID))

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *FlagStateResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state FlagStateResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	fs, err := r.client.GetFlagState(ctx, state.EnvironmentID.ValueString(), state.FlagID.ValueString())
	if err != nil {
		if errors.Is(err, client.ErrNotFound) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading flag state", err.Error())
		return
	}

	state.Enabled = types.BoolValue(fs.Enabled)
	state.DefaultVariant = types.StringValue(fs.DefaultVariant)

	var variants []VariantModel
	for _, v := range fs.Variants {
		variants = append(variants, VariantModel{
			Name:  types.StringValue(v.Name),
			Value: types.StringValue(v.Value),
		})
	}
	state.Variants = variants

	var targetingRules []TargetingRuleModel
	for _, tr := range fs.TargetingRules {
		rule := TargetingRuleModel{
			Name: types.StringValue(tr.Name),
		}
		if tr.Variant != "" {
			rule.Variant = types.StringValue(tr.Variant)
		} else {
			rule.Variant = types.StringNull()
		}
		
		var conditions []ConditionModel
		for _, c := range tr.Conditions {
			var vals []types.String
			for _, v := range c.Values {
				vals = append(vals, types.StringValue(v))
			}
			conditions = append(conditions, ConditionModel{
				Attribute: types.StringValue(c.Attribute),
				Operator:  types.StringValue(c.Operator),
				Values:    vals,
			})
		}
		rule.Conditions = conditions
		
		if tr.Rollout != nil && len(tr.Rollout.Percentages) > 0 {
			rule.Rollout = &RolloutModel{
				Percentages: make(map[string]types.Int64),
			}
			for k, v := range tr.Rollout.Percentages {
				rule.Rollout.Percentages[k] = types.Int64Value(int64(v))
			}
		}
		targetingRules = append(targetingRules, rule)
	}
	state.TargetingRules = targetingRules

	if fs.TelemetryTrigger != nil && fs.TelemetryTrigger.MetricName != "" {
		state.TelemetryTrigger = &TelemetryTriggerModel{
			MetricName: types.StringValue(fs.TelemetryTrigger.MetricName),
			Threshold:  types.Float64Value(fs.TelemetryTrigger.Threshold),
			Action:     types.StringValue(fs.TelemetryTrigger.Action),
		}
	} else {
		state.TelemetryTrigger = nil
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *FlagStateResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	r.Create(ctx, resource.CreateRequest{Plan: req.Plan}, &resource.CreateResponse{State: resp.State, Diagnostics: resp.Diagnostics})
}

func (r *FlagStateResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state FlagStateResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteFlagState(ctx, state.EnvironmentID.ValueString(), state.FlagID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting flag state", err.Error())
		return
	}
}
