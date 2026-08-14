# Research: Competitor UX Parity

## Technical Context Unknowns

No `NEEDS CLARIFICATION` items were identified in the technical context. The requirements are standard UI enhancements that can be accomplished with React and Tailwind CSS.

## Technology Choices & Best Practices

**Task**: Find best practices for global sidebars and toggle switches in React/Tailwind.

**Decision**: Use a persistent left sidebar with active state highlights. Use a native-feeling toggle switch (like Shadcn's Switch) for the flag states.
**Rationale**: Matches competitor patterns and provides clear, immediate access to environments without hidden dropdowns.
**Alternatives considered**: Top navigation bar (rejected as it scales poorly with many environments).

**Task**: Best practices for tag filtering.

**Decision**: Use a multi-select dropdown or a pill-based filter bar above the table.
**Rationale**: Scalable for many tags, easy to clear, visual feedback.
