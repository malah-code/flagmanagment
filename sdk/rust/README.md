# FlagManagment — Rust SDK

Ultra-high-performance, thread-safe Rust SDK and local flag evaluation engine for **FlagManagment**.

---

## Performance

- **Sub-microsecond Evaluation**: Zero network hops, thread-safe `RwLock` in-memory cached evaluation (<1 µs).
- **MurmurHash3 Bucketing**: Consistent cross-language bucketing with Node, Go, Python, Java, and .NET.

---

## Usage

```rust
use flagmanagment_sdk::{Client, Evaluator};
use serde_json::json;

fn main() {
    let client = Client::new();

    // Load or sync flags
    client.update_flag(
        "new-checkout-v2".to_string(),
        json!({
            "key": "new-checkout-v2",
            "enabled": true,
            "defaultVariant": "off",
            "variants": { "on": true, "off": false },
            "rules": [
                {
                    "rollout": { "on": 50, "off": 50 }
                }
            ]
        }),
    );

    // Evaluate flag locally
    let is_enabled = client.get_boolean_value("new-checkout-v2", false, Some("user-42"));
    println!("Flag enabled: {}", is_enabled);
}
```

