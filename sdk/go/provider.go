package sdk

import (
	"context"
	"fmt"

	"github.com/open-feature/go-sdk/openfeature"
)

// Provider implements the openfeature.FeatureProvider interface.
type Provider struct {
	client *Client
}

// NewProvider creates and connects a new FlagManagment provider.
func NewProvider(apiKey string, streamURL string) *Provider {
	client := NewClient(apiKey, streamURL)
	client.Connect()
	return &Provider{client: client}
}

// Metadata returns the provider's metadata.
func (p *Provider) Metadata() openfeature.Metadata {
	return openfeature.Metadata{
		Name: "FlagManagment-Go-Provider",
	}
}

// Init initializes the provider. Connection is already started in NewProvider.
func (p *Provider) Init(evalContext openfeature.EvaluationContext) error {
	return nil
}

// Shutdown gracefully shuts down the SSE streaming goroutine.
func (p *Provider) Shutdown() {
	p.client.Shutdown()
}

// Hooks returns any provider-specific hooks.
func (p *Provider) Hooks() []openfeature.Hook {
	return []openfeature.Hook{}
}

// BooleanEvaluation returns a boolean flag value.
func (p *Provider) BooleanEvaluation(ctx context.Context, flagKey string, defaultValue bool, evalContext openfeature.EvaluationContext) openfeature.BoolResolutionDetail {
	res, err := p.evaluate(flagKey, evalContext)
	if err != nil {
		return openfeature.BoolResolutionDetail{
			Value: defaultValue,
			ProviderResolutionDetail: openfeature.ProviderResolutionDetail{
				ResolutionError: openfeature.NewGeneralResolutionError(err.Error()),
				Reason:          openfeature.ErrorReason,
			},
		}
	}

	if b, ok := res.Value.(bool); ok {
		return openfeature.BoolResolutionDetail{
			Value: b,
			ProviderResolutionDetail: openfeature.ProviderResolutionDetail{
				Reason:  res.Reason,
				Variant: res.Variant,
			},
		}
	}

	return openfeature.BoolResolutionDetail{
		Value: defaultValue,
		ProviderResolutionDetail: openfeature.ProviderResolutionDetail{
			ResolutionError: openfeature.NewTypeMismatchResolutionError("flag value is not a boolean"),
			Reason:          openfeature.ErrorReason,
		},
	}
}

// StringEvaluation returns a string flag value.
func (p *Provider) StringEvaluation(ctx context.Context, flagKey string, defaultValue string, evalContext openfeature.EvaluationContext) openfeature.StringResolutionDetail {
	res, err := p.evaluate(flagKey, evalContext)
	if err != nil {
		return openfeature.StringResolutionDetail{
			Value: defaultValue,
			ProviderResolutionDetail: openfeature.ProviderResolutionDetail{
				ResolutionError: openfeature.NewGeneralResolutionError(err.Error()),
				Reason:          openfeature.ErrorReason,
			},
		}
	}

	if s, ok := res.Value.(string); ok {
		return openfeature.StringResolutionDetail{
			Value: s,
			ProviderResolutionDetail: openfeature.ProviderResolutionDetail{
				Reason:  res.Reason,
				Variant: res.Variant,
			},
		}
	}

	return openfeature.StringResolutionDetail{
		Value: defaultValue,
		ProviderResolutionDetail: openfeature.ProviderResolutionDetail{
			ResolutionError: openfeature.NewTypeMismatchResolutionError("flag value is not a string"),
			Reason:          openfeature.ErrorReason,
		},
	}
}

// FloatEvaluation returns a float flag value.
func (p *Provider) FloatEvaluation(ctx context.Context, flagKey string, defaultValue float64, evalContext openfeature.EvaluationContext) openfeature.FloatResolutionDetail {
	res, err := p.evaluate(flagKey, evalContext)
	if err != nil {
		return openfeature.FloatResolutionDetail{
			Value: defaultValue,
			ProviderResolutionDetail: openfeature.ProviderResolutionDetail{
				ResolutionError: openfeature.NewGeneralResolutionError(err.Error()),
				Reason:          openfeature.ErrorReason,
			},
		}
	}

	if f, ok := res.Value.(float64); ok {
		return openfeature.FloatResolutionDetail{
			Value: f,
			ProviderResolutionDetail: openfeature.ProviderResolutionDetail{
				Reason:  res.Reason,
				Variant: res.Variant,
			},
		}
	}

	return openfeature.FloatResolutionDetail{
		Value: defaultValue,
		ProviderResolutionDetail: openfeature.ProviderResolutionDetail{
			ResolutionError: openfeature.NewTypeMismatchResolutionError("flag value is not a float"),
			Reason:          openfeature.ErrorReason,
		},
	}
}

// IntEvaluation returns an int flag value.
func (p *Provider) IntEvaluation(ctx context.Context, flagKey string, defaultValue int64, evalContext openfeature.EvaluationContext) openfeature.IntResolutionDetail {
	res, err := p.evaluate(flagKey, evalContext)
	if err != nil {
		return openfeature.IntResolutionDetail{
			Value: defaultValue,
			ProviderResolutionDetail: openfeature.ProviderResolutionDetail{
				ResolutionError: openfeature.NewGeneralResolutionError(err.Error()),
				Reason:          openfeature.ErrorReason,
			},
		}
	}

	// JSON unmarshals all numbers to float64
	if f, ok := res.Value.(float64); ok {
		return openfeature.IntResolutionDetail{
			Value: int64(f),
			ProviderResolutionDetail: openfeature.ProviderResolutionDetail{
				Reason:  res.Reason,
				Variant: res.Variant,
			},
		}
	}

	return openfeature.IntResolutionDetail{
		Value: defaultValue,
		ProviderResolutionDetail: openfeature.ProviderResolutionDetail{
			ResolutionError: openfeature.NewTypeMismatchResolutionError("flag value is not an integer"),
			Reason:          openfeature.ErrorReason,
		},
	}
}

// ObjectEvaluation returns an object flag value.
func (p *Provider) ObjectEvaluation(ctx context.Context, flagKey string, defaultValue interface{}, evalContext openfeature.EvaluationContext) openfeature.InterfaceResolutionDetail {
	res, err := p.evaluate(flagKey, evalContext)
	if err != nil {
		return openfeature.InterfaceResolutionDetail{
			Value: defaultValue,
			ProviderResolutionDetail: openfeature.ProviderResolutionDetail{
				ResolutionError: openfeature.NewGeneralResolutionError(err.Error()),
				Reason:          openfeature.ErrorReason,
			},
		}
	}

	return openfeature.InterfaceResolutionDetail{
		Value: res.Value,
		ProviderResolutionDetail: openfeature.ProviderResolutionDetail{
			Reason:  res.Reason,
			Variant: res.Variant,
		},
	}
}

type evaluationResult struct {
	Value   interface{}
	Variant string
	Reason  openfeature.Reason
}

func (p *Provider) evaluate(flagKey string, evalContext openfeature.EvaluationContext) (evaluationResult, error) {
	flag, err := p.client.GetFlag(flagKey)
	if err != nil {
		return evaluationResult{}, fmt.Errorf("flag not found: %s", flagKey)
	}

	val, variant, reason, err := evaluateFlag(flag, evalContext)
	return evaluationResult{
		Value:   val,
		Variant: variant,
		Reason:  reason,
	}, err
}
