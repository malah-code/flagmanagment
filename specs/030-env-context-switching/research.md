# Research & Decisions

- **Decision**: Add a localized `isTransitioning` state in `FlagStatesList.tsx` that triggers via a `useEffect` whenever the `environmentId` prop changes.
- **Rationale**: React Query's `isLoading` only fires when there is no cached data. To ensure users perceive a context switch even when data is cached (instant render), an artificial micro-transition is required.
- **Decision**: Use a CSS opacity fade (`opacity-50` with a 200ms `setTimeout`) rather than full skeleton rows.
- **Rationale**: A fade is subtle, elegant, and matches the "100-200ms subtle fade" suggestion in the spec without the heavy boilerplate of building a skeleton table component.
