# Feature Specification: Environment SDK Key UX & Integration Guide

**Feature Branch**: `[###-feature-name]`

**Created**: 2026-08-10

**Status**: Draft

**Input**: User description: "specify the enhancement in detail before we start tasks and implementation: Environment SDK Key user experience enhancement matching competitor standards (Flagsmith, LaunchDarkly). Always-visible client SDK key in environment settings with 1-click copy, clear distinction between public client keys and private keys, and a modal with ready-to-use SDK code integration snippets for React, Node.js, Python, and Go."

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Always-Visible Client SDK Key (Priority: P1)

As an administrator or developer, I want to view and copy an environment's Client SDK key at any time from Environment Settings, so that I can easily integrate new client applications without being forced to regenerate or lose keys.

**Why this priority**: Eliminates developer friction. Developers frequently need to grab the environment key days or weeks after environment creation.

**Independent Test**: Can be verified by navigating to Environment Settings and clicking a permanent "Copy SDK Key" button next to any environment row.

**Acceptance Scenarios**:

1. **Given** the user is viewing Environment Settings, **When** they look at any environment card/row, **Then** the Client SDK Key (or deterministic client token) is displayed with a 1-click copy button.
2. **Given** the user clicks "Copy SDK Key", **When** the clipboard write succeeds, **Then** a visual confirmation (e.g. "Copied!" checkmark) appears for 2 seconds.

### User Story 2 - Interactive SDK Integration Guide (Priority: P1)

As a developer, I want to click an "Integration Guide" (`< />`) button on any environment to see pre-configured code snippets for React, Node.js, Python, and Go, so that I can copy working SDK initialization code immediately.

**Why this priority**: Speeds up onboarding and reduces SDK integration errors for developers using FlagManagment.

**Independent Test**: Can be verified by clicking the "Integration Guide" button on an environment and switching between React, Node.js, Python, and Go code snippet tabs with pre-populated environment keys.

**Acceptance Scenarios**:

1. **Given** the user clicks the "Integration Guide" (`< />`) button on an environment, **When** the modal opens, **Then** it presents code tabs for React, Node.js, Python, and Go.
2. **Given** a tab is selected, **When** the code snippet renders, **Then** the actual Environment SDK Key is automatically embedded in the code sample.
3. **Given** the user clicks "Copy Code", **When** copied, **Then** the exact snippet with the key is saved to the clipboard.

### User Story 3 - Distinct Public vs. Private Key Management (Priority: P2)

As a security-conscious administrator, I want a clear visual distinction between Public Client SDK Keys (safe for frontend/mobile apps) and Private/Server Admin Keys (protected/hashed), so that I understand which credentials to expose in frontend codebases.

**Why this priority**: Prevents security misconfigurations where admins accidentally put private admin secrets into public browser applications.

**Independent Test**: Can be verified by checking the Environment Settings UI, which clearly badges Client SDK Keys as "Public / Client SDK" and Secret Keys as "Private / Admin".

**Acceptance Scenarios**:

1. **Given** the user views key information for an environment, **When** they inspect the UI badges, **Then** the Client SDK Key is labeled as "Client-side / Public Key (React, Mobile, Web)" and Admin Keys are labeled as "Server Secret / Admin Key".

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST display the Client SDK Key in a readable/copyable format in the Environment Settings list.
- **FR-002**: System MUST provide a dedicated "Integration Guide" modal accessible from each environment row.
- **FR-003**: The Integration Guide modal MUST include tabs for React, Node.js, Python, and Go SDKs.
- **FR-004**: Code snippets in the Integration Guide modal MUST automatically pre-fill the selected environment's Client SDK Key.
- **FR-005**: System MUST visually distinguish between Public Client SDK Keys and Private Server Keys using badges and helper descriptions.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Developers can copy an SDK key and complete SDK initialization setup in under 30 seconds.
- **SC-002**: 100% of environment cards in Environment Settings provide instant 1-click access to the Client SDK Key.
- **SC-003**: 0% accidental usage of Private Admin keys in frontend applications due to clear UI labeling.

## Assumptions

- Client SDK Keys are read-only and intended for flag evaluation.
- The environment API returns the public Client SDK key (or deterministic client token) when listing environments for authenticated project members.
