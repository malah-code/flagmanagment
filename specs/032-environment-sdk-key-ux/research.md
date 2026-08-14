# Research: Environment SDK Key UX & Integration Guide

## Technical Context Unknowns

No `NEEDS CLARIFICATION` items were identified.

## Technology Choices & Best Practices

**Task**: Best practices for displaying public Client SDK keys vs Private Admin keys.

**Decision**: Display the Client SDK Key with a green/emerald "Public / Client Key" badge and a 1-click copy icon button. Keep Admin keys hidden behind a reveal/regenerate security threshold.
**Rationale**: Clarifies security boundaries for developers and prevents admin key exposure in client apps.

**Task**: Code snippet modal formatting.

**Decision**: Use dark-themed syntax-highlighted code blocks for React, Node.js, Python, and Go with a 1-click "Copy Snippet" action.
**Rationale**: Standard practice in modern developer platforms (Vercel, Supabase, LaunchDarkly).
