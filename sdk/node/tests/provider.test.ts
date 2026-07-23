import { FlagManagmentClient } from '../src/client';
import { FlagManagmentProvider } from '../src/provider';
import { OpenFeature } from '@openfeature/server-sdk';

describe('OpenFeature Provider Tests', () => {
  test('Resolves boolean values through OpenFeature API', async () => {
    const client = new FlagManagmentClient({
      environmentToken: 'test-token',
    });

    const provider = new FlagManagmentProvider(client);
    OpenFeature.setProvider(provider);

    const ofClient = OpenFeature.getClient();
    const val = await ofClient.getBooleanValue('non-existent-flag', true);
    expect(val).toBe(true);
  });
});
