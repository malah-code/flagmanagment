package api

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"

	"github.com/flagmanagment/backend/internal/dto"
	"github.com/flagmanagment/backend/internal/models"
	"github.com/flagmanagment/backend/internal/repository"
	"github.com/flagmanagment/backend/internal/sdk"
	"github.com/go-chi/chi/v5"
)

type SDKHandler struct {
	store repository.Store
}

func NewSDKHandler(store repository.Store) *SDKHandler {
	return &SDKHandler{store: store}
}

func (h *SDKHandler) RegisterRoutes(r chi.Router) {
	// Protected by AuthMiddleware (mounted in main router)
	r.Get("/evaluate/flags", h.EvaluateFlags)
	r.Post("/sdk/evaluate", h.EvaluateSingleFlag)
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

	// Build the response payload
	resFlags := make(map[string]dto.SDKFlag)
	for _, s := range states {
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
	if err != nil {
		RespondWithJSON(w, http.StatusOK, map[string]interface{}{
			"value":  false,
			"reason": "NO_STATE",
		})
		return
	}

	// PII Hashing per Constitution VII
	hashedContext := sdk.HashPII(&req.Context)

	var targetingRules json.RawMessage
	if state.TargetingRules != nil {
		if b, err := json.Marshal(state.TargetingRules); err == nil {
			targetingRules = b
		}
	}

	// Create models.FlagRule
	flagRule := &models.FlagRule{
		Key:              targetFlag.Key,
		Type:             string(targetFlag.Type),
		Enabled:          state.Enabled,
		DefaultVariation: "false", // Or remoteConfig base string
		TargetingRules:   targetingRules,
	}

	result := sdk.EvaluateFlag(flagRule, hashedContext)

	RespondWithJSON(w, http.StatusOK, map[string]interface{}{
		"value":  result.Value,
		"reason": result.Reason,
	})
}
