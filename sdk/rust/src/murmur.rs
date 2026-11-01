use std::io::Cursor;

/// Deterministic 32-bit MurmurHash3 bucketing (0..99).
pub fn bucket_user(targeting_key: &str) -> u32 {
    let mut cursor = Cursor::new(targeting_key.as_bytes());
    let hash = murmur3::murmur3_32(&mut cursor, 0).unwrap_or(0);
    hash % 100
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_consistent_bucketing() {
        let b1 = bucket_user("user-42");
        let b2 = bucket_user("user-42");
        assert_eq!(b1, b2);
        assert!(b1 < 100);
    }
}
