package com.flagmanagment.sdk

import java.nio.charset.StandardCharsets

/**
 * MurmurHash3 x86_32 implementation for deterministic cross-SDK bucketing.
 * Matches Go (spaolacci/murmur3), Java, Python (mmh3), .NET, React, and iOS implementations.
 * Seed is fixed to 0 for cross-language consistency.
 */
object MurmurHash3 {
    private const val C1 = 0xcc9e2d51.toInt()
    private const val C2 = 0x1b873593.toInt()
    private const val SEED = 0

    fun hash32(key: String, seed: Int = SEED): Int {
        val data = key.toByteArray(StandardCharsets.UTF_8)
        val len = data.size
        val nblocks = len / 4
        var h1 = seed

        // body
        for (i in 0 until nblocks) {
            val i4 = i * 4
            var k1 = (data[i4].toInt() and 0xff) or
                    ((data[i4 + 1].toInt() and 0xff) shl 8) or
                    ((data[i4 + 2].toInt() and 0xff) shl 16) or
                    ((data[i4 + 3].toInt() and 0xff) shl 24)

            k1 *= C1
            k1 = Integer.rotateLeft(k1, 15)
            k1 *= C2

            h1 = h1 xor k1
            h1 = Integer.rotateLeft(h1, 13)
            h1 = h1 * 5 + 0xe6546b64.toInt()
        }

        // tail
        val tail = nblocks * 4
        var k1 = 0
        when (len and 3) {
            3 -> {
                k1 = k1 xor ((data[tail + 2].toInt() and 0xff) shl 16)
                k1 = k1 xor ((data[tail + 1].toInt() and 0xff) shl 8)
                k1 = k1 xor (data[tail].toInt() and 0xff)
                k1 *= C1
                k1 = Integer.rotateLeft(k1, 15)
                k1 *= C2
                h1 = h1 xor k1
            }
            2 -> {
                k1 = k1 xor ((data[tail + 1].toInt() and 0xff) shl 8)
                k1 = k1 xor (data[tail].toInt() and 0xff)
                k1 *= C1
                k1 = Integer.rotateLeft(k1, 15)
                k1 *= C2
                h1 = h1 xor k1
            }
            1 -> {
                k1 = k1 xor (data[tail].toInt() and 0xff)
                k1 *= C1
                k1 = Integer.rotateLeft(k1, 15)
                k1 *= C2
                h1 = h1 xor k1
            }
        }

        // finalization
        h1 = h1 xor len
        h1 = fmix32(h1)

        return h1
    }

    fun bucketUser(targetingKey: String): Int {
        val hash = hash32(targetingKey).toLong() and 0xFFFFFFFFL
        return (hash % 100).toInt()
    }

    private fun fmix32(h: Int): Int {
        var h = h
        h = h xor (h ushr 16)
        h *= 0x85ebca6b.toInt()
        h = h xor (h ushr 13)
        h *= 0xc2b2ae35.toInt()
        h = h xor (h ushr 16)
        return h
    }
}
