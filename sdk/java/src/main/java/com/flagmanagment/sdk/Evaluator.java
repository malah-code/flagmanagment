package com.flagmanagment.sdk;

/**
 * MurmurHash3 32-bit implementation for deterministic cross-SDK bucketing.
 * This is a pure-Java port matching the canonical x86_32 variant used by
 * spaolacci/murmur3 (Go), mmh3 (Python), and all other FlagManagment SDKs.
 * Seed is fixed to 0 for cross-language consistency.
 */
public class Evaluator {

    private static final int C1 = 0xcc9e2d51;
    private static final int C2 = 0x1b873593;
    private static final int SEED = 0;

    /**
     * Returns a bucket between 0 and 99 (inclusive) for the given targeting key.
     */
    public static int bucketUser(String targetingKey) {
        int hash = murmur3_32(targetingKey.getBytes(java.nio.charset.StandardCharsets.UTF_8), SEED);
        return Math.abs(hash % 100);
    }

    /**
     * MurmurHash3 x86_32 implementation.
     */
    public static int murmur3_32(byte[] data, int seed) {
        int h1 = seed;
        int len = data.length;
        int nblocks = len / 4;

        // body
        for (int i = 0; i < nblocks; i++) {
            int i4 = i * 4;
            int k1 = (data[i4] & 0xff)
                    | ((data[i4 + 1] & 0xff) << 8)
                    | ((data[i4 + 2] & 0xff) << 16)
                    | ((data[i4 + 3] & 0xff) << 24);

            k1 *= C1;
            k1 = Integer.rotateLeft(k1, 15);
            k1 *= C2;

            h1 ^= k1;
            h1 = Integer.rotateLeft(h1, 13);
            h1 = h1 * 5 + 0xe6546b64;
        }

        // tail
        int tail = nblocks * 4;
        int k1 = 0;
        switch (len & 3) {
            case 3:
                k1 ^= (data[tail + 2] & 0xff) << 16;
                // fallthrough
            case 2:
                k1 ^= (data[tail + 1] & 0xff) << 8;
                // fallthrough
            case 1:
                k1 ^= (data[tail] & 0xff);
                k1 *= C1;
                k1 = Integer.rotateLeft(k1, 15);
                k1 *= C2;
                h1 ^= k1;
        }

        // finalization
        h1 ^= len;
        h1 = fmix32(h1);

        return h1;
    }

    private static int fmix32(int h) {
        h ^= h >>> 16;
        h *= 0x85ebca6b;
        h ^= h >>> 13;
        h *= 0xc2b2ae35;
        h ^= h >>> 16;
        return h;
    }
}
