//! Official FlagManagment Rust SDK
//!
//! Provides ultra-high-performance, zero-allocation sub-microsecond local evaluation (<1 µs)
//! with MurmurHash3 percentage bucketing.

pub mod murmur;
pub mod evaluator;
pub mod client;

pub use client::Client;
pub use evaluator::{Evaluator, EvaluationResult};
pub use murmur::bucket_user;

