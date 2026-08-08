# Research & Decisions: Contextual Targeting Engine

## 1. Rule Storage Format

**Decision**: Store rules as a JSON array within the existing `targeting_rules` (JSONB) column in `environment_flag_states`.
**Rationale**: PostgreSQL JSONB is perfectly suited for flexible, schema-less rule storage. It allows adding new operators in the future without database migrations.
**Alternatives considered**: Creating a separate relational table for rules. Rejected because it would require complex joins for flag evaluation, impacting SDK payload generation speed.

## 2. Rule Evaluation Engine

**Decision**: Implement the rule evaluation logic natively in Go within the Server-Side SDK and backend API (for testing).
**Rationale**: Go's native capabilities for string comparison and regex matching are highly performant. A custom, lightweight evaluator avoids the overhead of a full third-party rules engine.
**Alternatives considered**: Using a rules engine like `go-rodeo` or `grule`. Rejected because they are too heavy for simple flag targeting (Equals, Contains, Regex) and would bloat the SDK.

## 3. Regex Performance

**Decision**: Pre-compile and cache regex patterns upon flag creation/update, or compile them once when the SDK downloads the ruleset.
**Rationale**: `regexp.Compile` is relatively expensive in Go. Recompiling a regex for every evaluation would violate the `< 1ms` performance constraint.
**Alternatives considered**: Using `regexp.MatchString` directly. Rejected because it compiles the regex on every call.

## 4. Rule Structure (AND/OR logic)

**Decision**: A single `TargetingRule` object will contain a list of `Conditions` (which act as AND). The overall flag will have a list of `TargetingRules` (which act as OR).
**Rationale**: This is a standard industry pattern (e.g., LaunchDarkly). "First match wins" for the rules (OR). Within a rule, all conditions must be true (AND).
**Alternatives considered**: Complex AST-like boolean expressions. Rejected because they are unnecessarily complex for 99% of feature flag use cases.
