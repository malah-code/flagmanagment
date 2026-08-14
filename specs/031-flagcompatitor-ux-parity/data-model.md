# Data Model: Competitor UX Parity

No new backend database entities are introduced. The following existing entities are interacted with:

- **Environment**: Used to populate the sidebar navigation.
- **Flag**: Contains `tags` (array of strings) and `initial_value` (if backend supports setting it upon creation).
- **FlagState**: Toggled directly from the UI lists.
