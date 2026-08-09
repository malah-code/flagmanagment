import { useSyncExternalStore, useCallback } from 'react';
import { useFlagContext } from './FlagProvider';
import { FlagClient } from './client';

/**
 * MurmurHash3 x86_32 implementation matching Go (spaolacci/murmur3),
 * Java, Python (mmh3), and .NET implementations.
 * Seed is fixed to 0 for cross-language consistency.
 */
const murmur3_32 = (key: string, seed: number = 0): number => {
  const data = new TextEncoder().encode(key);
  const len = data.length;
  const nblocks = Math.floor(len / 4);
  let h1 = seed >>> 0;

  const C1 = 0xcc9e2d51;
  const C2 = 0x1b873593;

  // body
  for (let i = 0; i < nblocks; i++) {
    const i4 = i * 4;
    let k1 =
      (data[i4] | (data[i4 + 1] << 8) | (data[i4 + 2] << 16) | (data[i4 + 3] << 24)) >>> 0;

    k1 = Math.imul(k1, C1) >>> 0;
    k1 = ((k1 << 15) | (k1 >>> 17)) >>> 0;
    k1 = Math.imul(k1, C2) >>> 0;

    h1 = (h1 ^ k1) >>> 0;
    h1 = ((h1 << 13) | (h1 >>> 19)) >>> 0;
    h1 = (Math.imul(h1, 5) + 0xe6546b64) >>> 0;
  }

  // tail
  const tail = nblocks * 4;
  let k1 = 0;
  switch (len & 3) {
    case 3:
      k1 ^= data[tail + 2] << 16;
    // fallthrough
    case 2:
      k1 ^= data[tail + 1] << 8;
    // fallthrough
    case 1:
      k1 ^= data[tail];
      k1 = Math.imul(k1, C1) >>> 0;
      k1 = ((k1 << 15) | (k1 >>> 17)) >>> 0;
      k1 = Math.imul(k1, C2) >>> 0;
      h1 = (h1 ^ k1) >>> 0;
  }

  // finalization
  h1 = (h1 ^ len) >>> 0;
  h1 ^= h1 >>> 16;
  h1 = Math.imul(h1, 0x85ebca6b) >>> 0;
  h1 ^= h1 >>> 13;
  h1 = Math.imul(h1, 0xc2b2ae35) >>> 0;
  h1 ^= h1 >>> 16;

  return h1 >>> 0;
};

const bucketUser = (key: string): number => {
  return murmur3_32(key) % 100;
};

const evaluateLocally = (client: FlagClient, flagKey: string, defaultValue: any) => {
  const flag = client.getFlag(flagKey);
  if (!flag || !flag.enabled) {
    return defaultValue;
  }

  const userContext = client.getContext();
  const targetingKey = userContext?.targetingKey;

  if (targetingKey && flag.rules && Array.isArray(flag.rules)) {
    for (const rule of flag.rules) {
      if (rule.rollout) {
        const bucket = bucketUser(targetingKey);
        let currentThreshold = 0;
        // Sort keys for deterministic iteration (matches Go, Java, Python, .NET)
        const sortedVariants = Object.keys(rule.rollout).sort();
        for (const variant of sortedVariants) {
          currentThreshold += Number(rule.rollout[variant]);
          if (bucket < currentThreshold) {
            return flag.variants?.[variant] ?? defaultValue;
          }
        }
      }
    }
  }

  const defaultVariant = flag.defaultVariant;
  return flag.variants?.[defaultVariant] ?? defaultValue;
};

export const useFlag = <T = any>(flagKey: string, defaultValue: T): T => {
  const { client } = useFlagContext();

  const subscribe = useCallback(
    (onStoreChange: () => void) => client.subscribe(onStoreChange),
    [client]
  );

  const getSnapshot = useCallback(
    () => client.getFlags(),
    [client]
  );

  useSyncExternalStore(subscribe, getSnapshot, getSnapshot);

  return evaluateLocally(client, flagKey, defaultValue);
};

export const useFlags = () => {
  const { client } = useFlagContext();

  const subscribe = useCallback(
    (onStoreChange: () => void) => client.subscribe(onStoreChange),
    [client]
  );

  const getSnapshot = useCallback(
    () => client.getFlags(),
    [client]
  );

  return useSyncExternalStore(subscribe, getSnapshot, getSnapshot);
};
