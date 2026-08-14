# Research & Decisions

- **Decision**: Introduce `react-hot-toast` to handle non-blocking global toast notifications.
- **Rationale**: `react-hot-toast` is lightweight, easy to integrate, and fits perfectly for the success/error toasts required by the spec.
- **Alternatives considered**: Building a custom context-based toast provider was considered, but it introduces unnecessary boilerplate and complexity compared to using an established, headless-friendly library.

- **Decision**: Update `handleToggle` in `FlagStatesList.tsx` to handle loading states and fire toast notifications.
- **Rationale**: We can track loading state using `updateMutation.isPending`, or local state if we need per-row spinners. Actually, `useMutation` from `react-query` returns `isPending`, but since it's one mutation for the whole list, clicking a toggle might spin *all* toggles if we just use `updateMutation.isPending`. Wait! The `FlagStatesList` currently disables the toggle button using `disabled={updateMutation.isPending}`. This means all toggles disable when one is clicked. To provide a better experience, we can use a local state array or track `updatingStateId`. For now, just utilizing `updateMutation.isPending` alongside `updateMutation.variables?.id === state.id` allows us to target the specific row's spinner.
