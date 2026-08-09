import { useSyncExternalStore, useCallback } from 'react';
import { useFlagContext } from './FlagProvider';
import { FlagClient } from './client';

import { evaluateLocally } from './evaluator';

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
