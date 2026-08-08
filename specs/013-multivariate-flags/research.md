# Phase 0: Research - Multivariate Flags

## MurmurHash3 Implementation in Go
- **Decision**: Use `github.com/spaolacci/murmur3` for MurmurHash3 bucketing.
- **Rationale**: It is the most widely used and tested Go implementation of MurmurHash3. It produces 32-bit and 64-bit non-cryptographic hashes that are highly performant (required for the < 1ms evaluation constraint).
- **Alternatives considered**: `crypto/sha256` (too slow/cryptographic, overkill for bucketing, used already for PII hashing but not ideal for fast numeric bucketing). `hash/crc32` (less uniform distribution than MurmurHash3).

## Percentage Bucketing Logic
- **Decision**: Use `hash(flagKey + identityKey) % 10000` to get a bucket from 0 to 9999. Map variations to these 10,000 buckets (e.g., 33.33% = 3333 buckets).
- **Rationale**: Allowing percentages up to two decimal places (e.g. 33.33%) requires an integer space of at least 10,000 to avoid floating-point inaccuracies.
- **Alternatives considered**: Floating point mapping 0.0-1.0 (prone to IEEE 754 precision issues). modulo 100 (only allows whole numbers).
