export * from './types';
export * from './client';
export * from './provider';

import { FlagWebClient } from './client';
import { WebClientOptions } from './types';

/**
 * Factory method for initializing FlagManagment in browser applications.
 */
export function initFlagManagment(options: WebClientOptions): FlagWebClient {
  const client = new FlagWebClient(options);
  client.connect().catch((err) => {
    console.warn('[FlagManagment Web SDK] initial connect warning:', err);
  });
  return client;
}
