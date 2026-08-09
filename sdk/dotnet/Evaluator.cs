using System;
using System.Text;

namespace FlagManagment.Sdk
{
    /// <summary>
    /// MurmurHash3 32-bit (x86) implementation for deterministic cross-SDK bucketing.
    /// Matches Go (spaolacci/murmur3), Python (mmh3), and Java implementations.
    /// Seed is fixed to 0 for cross-language consistency.
    /// </summary>
    public static class Evaluator
    {
        private const uint C1 = 0xcc9e2d51;
        private const uint C2 = 0x1b873593;
        private const uint Seed = 0;

        public static int BucketUser(string targetingKey)
        {
            byte[] data = Encoding.UTF8.GetBytes(targetingKey);
            uint hash = MurmurHash3_32(data, Seed);
            return (int)(hash % 100);
        }

        public static uint MurmurHash3_32(byte[] data, uint seed)
        {
            uint h1 = seed;
            int len = data.Length;
            int nblocks = len / 4;

            // body
            for (int i = 0; i < nblocks; i++)
            {
                int i4 = i * 4;
                uint k1 = (uint)(data[i4]
                    | (data[i4 + 1] << 8)
                    | (data[i4 + 2] << 16)
                    | (data[i4 + 3] << 24));

                k1 *= C1;
                k1 = RotateLeft(k1, 15);
                k1 *= C2;

                h1 ^= k1;
                h1 = RotateLeft(h1, 13);
                h1 = h1 * 5 + 0xe6546b64;
            }

            // tail
            int tail = nblocks * 4;
            uint k1Tail = 0;
            switch (len & 3)
            {
                case 3:
                    k1Tail ^= (uint)data[tail + 2] << 16;
                    goto case 2;
                case 2:
                    k1Tail ^= (uint)data[tail + 1] << 8;
                    goto case 1;
                case 1:
                    k1Tail ^= data[tail];
                    k1Tail *= C1;
                    k1Tail = RotateLeft(k1Tail, 15);
                    k1Tail *= C2;
                    h1 ^= k1Tail;
                    break;
            }

            // finalization
            h1 ^= (uint)len;
            h1 = FMix32(h1);

            return h1;
        }

        private static uint RotateLeft(uint value, int count)
        {
            return (value << count) | (value >> (32 - count));
        }

        private static uint FMix32(uint h)
        {
            h ^= h >> 16;
            h *= 0x85ebca6b;
            h ^= h >> 13;
            h *= 0xc2b2ae35;
            h ^= h >> 16;
            return h;
        }
    }
}
