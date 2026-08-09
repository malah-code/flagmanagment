import Foundation

/// MurmurHash3 x86_32 implementation for deterministic cross-SDK bucketing.
/// Matches Go (spaolacci/murmur3), Java, Python (mmh3), .NET, and React implementations.
/// Seed is fixed to 0 for cross-language consistency.
public struct MurmurHash3 {
    
    private static let c1: UInt32 = 0xcc9e2d51
    private static let c2: UInt32 = 0x1b873593
    
    public static func hash32(_ key: String, seed: UInt32 = 0) -> UInt32 {
        let data = Array(key.utf8)
        let len = data.count
        let nblocks = len / 4
        var h1 = seed
        
        // body
        for i in 0..<nblocks {
            let i4 = i * 4
            var k1 = UInt32(data[i4])
                | (UInt32(data[i4 + 1]) << 8)
                | (UInt32(data[i4 + 2]) << 16)
                | (UInt32(data[i4 + 3]) << 24)
            
            k1 = k1 &* c1
            k1 = (k1 << 15) | (k1 >> 17)
            k1 = k1 &* c2
            
            h1 ^= k1
            h1 = (h1 << 13) | (h1 >> 19)
            h1 = h1 &* 5 &+ 0xe6546b64
        }
        
        // tail
        let tail = nblocks * 4
        var k1: UInt32 = 0
        switch len & 3 {
        case 3:
            k1 ^= UInt32(data[tail + 2]) << 16
            fallthrough
        case 2:
            k1 ^= UInt32(data[tail + 1]) << 8
            fallthrough
        case 1:
            k1 ^= UInt32(data[tail])
            k1 = k1 &* c1
            k1 = (k1 << 15) | (k1 >> 17)
            k1 = k1 &* c2
            h1 ^= k1
        default:
            break
        }
        
        // finalization
        h1 ^= UInt32(len)
        h1 = fmix32(h1)
        
        return h1
    }
    
    public static func bucketUser(_ targetingKey: String) -> Int {
        return Int(hash32(targetingKey) % 100)
    }
    
    private static func fmix32(_ h: UInt32) -> UInt32 {
        var h = h
        h ^= h >> 16
        h = h &* 0x85ebca6b
        h ^= h >> 13
        h = h &* 0xc2b2ae35
        h ^= h >> 16
        return h
    }
}
