const test = require('node:test');
const assert = require('node:assert');
const app = require('../src/index');

test('GET /health returns ok status', async () => {
  const server = app.listen(0);
  const { port } = server.address();

  try {
    const res = await fetch(`http://localhost:${port}/health`);
    assert.strictEqual(res.status, 200);
    const body = await res.json();
    assert.strictEqual(body.status, 'ok');
  } finally {
    server.close();
  }
});

test('GET /api/checkout returns fallback flow when server offline', async () => {
  const server = app.listen(0);
  const { port } = server.address();

  try {
    const res = await fetch(`http://localhost:${port}/api/checkout?userId=test-user`);
    assert.strictEqual(res.status, 200);
    const body = await res.json();
    assert.ok(body.flow);
  } finally {
    server.close();
  }
});

