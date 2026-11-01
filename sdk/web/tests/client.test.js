const test = require('node:test');
const assert = require('node:assert');
const { FlagWebClient } = require('../dist/client');

test('FlagWebClient returns default values when cache is empty', () => {
  const client = new FlagWebClient({
    apiKey: 'test-key',
    apiUrl: 'http://localhost:8080',
    enableLocalStorage: false,
  });

  assert.strictEqual(client.getBooleanValue('unknown-flag', false), false);
  assert.strictEqual(client.getBooleanValue('unknown-flag-true', true), true);
  assert.strictEqual(client.getStringValue('theme', 'dark'), 'dark');
  assert.strictEqual(client.getNumberValue('limit', 42), 42);
  assert.deepStrictEqual(client.getObjectValue('config', { a: 1 }), { a: 1 });
});

test('FlagWebClient subscribes to events and emits changes', () => {
  const client = new FlagWebClient({
    apiKey: 'test-key',
    apiUrl: 'http://localhost:8080',
    enableLocalStorage: false,
  });

  let changedFlag = '';
  let newValue = null;

  const unsubscribe = client.on('change', (key, val) => {
    changedFlag = key;
    newValue = val;
  });

  // Simulate updating flags internally
  client['updateFlags']({
    'dark-mode': { value: true, variant: 'on' },
  });

  assert.strictEqual(changedFlag, 'dark-mode');
  assert.strictEqual(newValue, true);
  assert.strictEqual(client.getBooleanValue('dark-mode', false), true);

  unsubscribe();
});

