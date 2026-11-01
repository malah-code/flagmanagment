# FlagManagment — React & Vite Sample

This is a complete, runnable Single-Page Application (SPA) demonstrating how to integrate the **FlagManagment React SDK** (`@flagmanagment/react-sdk`) using React 19, TypeScript, and Vite.

---

## Features Demonstrated

1. **`<FlagProvider>` Setup**: Wrapping the React application tree with your client-side environment key.
2. **Context-Aware Targeting**: Dynamic updates when user attributes change (`targetingKey`, subscription `tier`, `country`).
3. **`useFlag<T>()` Hook**: Sub-millisecond synchronous in-memory flag access using React 18/19 `useSyncExternalStore`.
4. **`useFlags()` Store**: Accessing the entire dictionary of evaluated flags with zero flicker.
5. **Interactive Identity Switcher**: Live switching between Free, Pro, and Enterprise users to demonstrate real-time rule evaluation.

---

## Quickstart

### 1. Install Dependencies
```bash
npm install
```

### 2. Run the Development Server
```bash
npm run dev
```
Open your browser at `http://localhost:5173`.

---

## Code Example

```tsx
import { FlagProvider, useFlag } from '@flagmanagment/react-sdk';

function App() {
  return (
    <FlagProvider
      apiKey="flg_client_sample_key"
      apiUrl="http://localhost:8080"
      context={{ targetingKey: "user-123", tier: "premium" }}
    >
      <MyFeature />
    </FlagProvider>
  );
}

function MyFeature() {
  const isEnabled = useFlag<boolean>('new-checkout-v2', false);
  return isEnabled ? <V2Checkout /> : <StandardCheckout />;
}
```

---

## Build for Production
```bash
npm run build
```

