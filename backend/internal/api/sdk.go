package api

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"time"

	"github.com/flagmanagment/backend/internal/dto"
	"github.com/flagmanagment/backend/internal/models"
	"github.com/flagmanagment/backend/internal/repository"
	"github.com/flagmanagment/backend/internal/sdk"
	"github.com/flagmanagment/backend/internal/services"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type SDKHandler struct {
	store         repository.Store
	metricService services.MetricAggregationService
}

func NewSDKHandler(store repository.Store, metricService services.MetricAggregationService) *SDKHandler {
	return &SDKHandler{store: store, metricService: metricService}
}

func (h *SDKHandler) RegisterRoutes(r chi.Router) {
	// Protected by AuthMiddleware (mounted in main router)
	r.Get("/evaluate/flags", h.EvaluateFlags)
	r.Post("/sdk/evaluate", h.EvaluateSingleFlag)
	r.Post("/sdk/metrics", h.IngestMetrics)
	r.Post("/client/evaluate", h.EvaluateClientFlags)
}

func (h *SDKHandler) EvaluateFlags(w http.ResponseWriter, r *http.Request) {
	env := GetEnvironmentFromContext(r.Context())
	if env == nil {
		RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Fetch all flags for the project to get their keys and types
	flags, _, err := h.store.FlagRepo().ListByProject(r.Context(), env.ProjectID, 10000, 0)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Failed to load project flags")
		return
	}

	// Build a map of flag ID to flag for quick lookup
	flagMap := make(map[string]*models.FeatureFlag)
	for _, f := range flags {
		flagMap[f.ID.String()] = f
	}

	// Fetch all flag states for the environment
	states, err := h.store.FlagStateRepo().ListByEnvironment(r.Context(), env.ID)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Failed to load environment flag states")
		return
	}

	// Calculate ETag
	var timestamps []string
	timestamps = append(timestamps, fmt.Sprintf("env_%d", env.UpdatedAt.UnixNano()))
	for _, s := range states {
		timestamps = append(timestamps, fmt.Sprintf("state_%s_%d", s.FeatureFlagID.String(), s.UpdatedAt.UnixNano()))
	}
	for _, f := range flags {
		timestamps = append(timestamps, fmt.Sprintf("flag_%s_%d", f.ID.String(), f.UpdatedAt.UnixNano()))
	}

	// Sort timestamps to ensure deterministic ETag
	sort.Strings(timestamps)

	hash := sha256.New()
	for _, t := range timestamps {
		hash.Write([]byte(t))
	}
	etag := fmt.Sprintf(`"%s"`, hex.EncodeToString(hash.Sum(nil)))

	if match := r.Header.Get("If-None-Match"); match == etag {
		w.Header().Set("ETag", etag)
		w.WriteHeader(http.StatusNotModified)
		return
	}

	// Build the response payload (excluding ARCHIVED flags)
	resFlags := make(map[string]dto.SDKFlag)
	for _, s := range states {
		if s.LifecycleState == models.LifecycleArchived {
			continue
		}
		f, ok := flagMap[s.FeatureFlagID.String()]
		if !ok {
			continue
		}
		resFlags[f.Key] = dto.SDKFlag{
			Enabled: s.Enabled,
			Type:    string(f.Type),
			Rules:   s.TargetingRules,
			Value:   s.RemoteConfig,
		}
	}

	w.Header().Set("ETag", etag)
	RespondWithJSON(w, http.StatusOK, dto.SDKEvaluationResponse{
		EnvironmentID: env.ID.String(),
		Flags:         resFlags,
	})
}

type EvaluateSingleFlagRequest struct {
	FlagKey string                   `json:"flagKey"`
	Context models.EvaluationContext `json:"context"`
}

func (h *SDKHandler) EvaluateSingleFlag(w http.ResponseWriter, r *http.Request) {
	env := GetEnvironmentFromContext(r.Context())
	if env == nil {
		RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var req EvaluateSingleFlagRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Fetch all flag states for the environment (or we could fetch just one if we join with flags table to get key)
	// For MVP: load project flags to find flag ID from key
	flags, _, err := h.store.FlagRepo().ListByProject(r.Context(), env.ProjectID, 10000, 0)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Failed to load project flags")
		return
	}

	var targetFlag *models.FeatureFlag
	for _, f := range flags {
		if f.Key == req.FlagKey {
			targetFlag = f
			break
		}
	}

	if targetFlag == nil {
		RespondWithError(w, http.StatusNotFound, "Flag not found")
		return
	}

	state, err := h.store.FlagStateRepo().GetByEnvAndFlag(r.Context(), env.ID, targetFlag.ID)
	if err != nil || state.LifecycleState == models.LifecycleArchived {
		RespondWithJSON(w, http.StatusOK, map[string]interface{}{
			"value":  false,
			"reason": "NO_STATE",
		})
		return
	}

	if h.metricService != nil {
		h.metricService.RecordEvaluation(state.ID, time.Now())
	}

	// PII Hashing per Constitution VII
	hashedContext := sdk.HashPII(&req.Context, env.Salt)

	var targetingRules json.RawMessage
	if state.TargetingRules != nil {
		if b, err := json.Marshal(state.TargetingRules); err == nil {
			targetingRules = b
		}
	}

	var rolloutRules json.RawMessage
	if state.RolloutRules != nil {
		if b, err := json.Marshal(state.RolloutRules); err == nil {
			rolloutRules = b
		}
	}

	var variations json.RawMessage
	if targetFlag.Variations != nil {
		if b, err := json.Marshal(targetFlag.Variations); err == nil {
			variations = b
		}
	}

	defaultVar := "false"
	if state.DefaultVariation != "" {
		defaultVar = state.DefaultVariation
	}

	// Create models.FlagRule
	flagRule := &models.FlagRule{
		Key:              targetFlag.Key,
		Type:             string(targetFlag.Type),
		Enabled:          state.Enabled,
		DefaultVariation: defaultVar,
		TargetingRules:   targetingRules,
		RolloutRules:     rolloutRules,
		Variations:       variations,
	}

	// We need to pass rulesMap. Since we are in EvaluateSingleFlag, we can build it.
	rulesMap := make(map[string]*models.FlagRule)
	for _, f := range flags {
		fState, err := h.store.FlagStateRepo().GetByEnvAndFlag(r.Context(), env.ID, f.ID)
		if err == nil && fState.LifecycleState != models.LifecycleArchived {
			var tRules, rRules, vars json.RawMessage
			if fState.TargetingRules != nil {
				b, _ := json.Marshal(fState.TargetingRules)
				tRules = b
			}
			if fState.RolloutRules != nil {
				b, _ := json.Marshal(fState.RolloutRules)
				rRules = b
			}
			if f.Variations != nil {
				b, _ := json.Marshal(f.Variations)
				vars = b
			}
			defVar := "false"
			if fState.DefaultVariation != "" {
				defVar = fState.DefaultVariation
			}
			parentKey := ""
			if f.ParentFlagID != nil {
				// find parent key
				for _, pf := range flags {
					if pf.ID == *f.ParentFlagID {
						parentKey = pf.Key
						break
					}
				}
			}
			rulesMap[f.Key] = &models.FlagRule{
				Key:              f.Key,
				Type:             string(f.Type),
				Enabled:          fState.Enabled,
				DefaultVariation: defVar,
				TargetingRules:   tRules,
				RolloutRules:     rRules,
				Variations:       vars,
				ParentFlagKey:    parentKey,
			}
		}
	}

	result := sdk.EvaluateFlag(flagRule, hashedContext, rulesMap, env.Salt)

	RespondWithJSON(w, http.StatusOK, map[string]interface{}{
		"value":  result.Value,
		"reason": result.Reason,
	})
}

func (h *SDKHandler) IngestMetrics(w http.ResponseWriter, r *http.Request) {
	env := GetEnvironmentFromContext(r.Context())
	if env == nil {
		RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var req dto.SDKMetricsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// For efficiency, list all project flags and all environment states to map them
	flags, _, err := h.store.FlagRepo().ListByProject(r.Context(), env.ProjectID, 10000, 0)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Failed to load project flags")
		return
	}

	states, err := h.store.FlagStateRepo().ListByEnvironment(r.Context(), env.ID)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Failed to load environment flag states")
		return
	}

	// Map FlagID -> StateID
	stateMap := make(map[string]string)
	for _, s := range states {
		stateMap[s.FeatureFlagID.String()] = s.ID.String()
	}

	// Map FlagKey -> FlagID
	flagMap := make(map[string]string)
	for _, f := range flags {
		flagMap[f.Key] = f.ID.String()
	}

	if h.metricService != nil {
		for _, eval := range req.Evaluations {
			flagIDStr, ok := flagMap[eval.FlagKey]
			if !ok {
				continue
			}
			stateIDStr, ok := stateMap[flagIDStr]
			if !ok {
				continue
			}
			
			// Try parsing timestamp, fallback to time.Now
			ts, err := time.Parse(time.RFC3339, eval.Timestamp)
			if err != nil {
				ts = time.Now()
			}
			
			// Parse stateID string to uuid.UUID
			stateUUID, err := uuid.Parse(stateIDStr)
			if err != nil {
				continue
			}

			h.metricService.RecordEvaluation(stateUUID, ts)
		}
	}

	RespondWithJSON(w, http.StatusOK, map[string]interface{}{"status": "ok"})
}

func (h *SDKHandler) EvaluateClientFlags(w http.ResponseWriter, r *http.Request) {
	env := GetEnvironmentFromContext(r.Context())
	if env == nil {
		RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var req dto.EvaluateClientRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	var evalCtx models.EvaluationContext
	if identity, ok := req.Context["identity"].(string); ok {
		evalCtx.Identity = identity
	} else if targetingKey, ok := req.Context["targetingKey"].(string); ok {
		evalCtx.Identity = targetingKey
	}
	evalCtx.Attributes = req.Context

	flags, _, err := h.store.FlagRepo().ListByProject(r.Context(), env.ProjectID, 10000, 0)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Failed to load project flags")
		return
	}

	states, err := h.store.FlagStateRepo().ListByEnvironment(r.Context(), env.ID)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Failed to load environment flag states")
		return
	}

	// PII Hashing
	hashedContext := sdk.HashPII(&evalCtx, env.Salt)

	// Build Rules map
	rulesMap := make(map[string]*models.FlagRule)
	stateMap := make(map[string]*models.EnvironmentFlagState)
	for _, s := range states {
		if s.LifecycleState == models.LifecycleArchived {
			continue
		}
		// find flag
		var flag *models.FeatureFlag
		for _, f := range flags {
			if f.ID == s.FeatureFlagID {
				flag = f
				break
			}
		}
		if flag == nil {
			continue
		}
		stateMap[flag.Key] = s

		var tRules, rRules, vars json.RawMessage
		if s.TargetingRules != nil {
			b, _ := json.Marshal(s.TargetingRules)
			tRules = b
		}
		if s.RolloutRules != nil {
			b, _ := json.Marshal(s.RolloutRules)
			rRules = b
		}
		if flag.Variations != nil {
			b, _ := json.Marshal(flag.Variations)
			vars = b
		}
		defVar := "false"
		if s.DefaultVariation != "" {
			defVar = s.DefaultVariation
		}
		
		parentKey := ""
		if flag.ParentFlagID != nil {
			for _, pf := range flags {
				if pf.ID == *flag.ParentFlagID {
					parentKey = pf.Key
					break
				}
			}
		}
		
		rulesMap[flag.Key] = &models.FlagRule{
			Key:              flag.Key,
			Type:             string(flag.Type),
			Enabled:          s.Enabled,
			DefaultVariation: defVar,
			TargetingRules:   tRules,
			RolloutRules:     rRules,
			Variations:       vars,
			ParentFlagKey:    parentKey,
		}
	}

	results := make(map[string]map[string]interface{})
	for key, rule := range rulesMap {
		res := sdk.EvaluateFlag(rule, hashedContext, rulesMap, env.Salt)
		results[key] = map[string]interface{}{
			"value":  res.Value,
			"reason": res.Reason,
		}
		
		// Record metric for client evaluations
		if h.metricService != nil {
			if st, ok := stateMap[key]; ok {
				h.metricService.RecordEvaluation(st.ID, time.Now())
			}
		}
	}

	RespondWithJSON(w, http.StatusOK, results)
}
