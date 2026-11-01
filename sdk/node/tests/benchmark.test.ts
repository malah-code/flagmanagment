import { evaluateFlag } from '../src/evaluator';
import { FlagRule } from '../src/types';

describe('Performance Benchmarks (Constitution Principle IV)', () => {
  test('In-memory flag evaluation executes in under 1 millisecond', () => {
    const flag: FlagRule = {
      key: 'perf-flag',
      type: 'BOOLEAN',
      enabled: true,
      defaultVariation: 'true',
    };

    const ctx = {
      identity: 'user-perf-123',
      attributes: { tier: 'enterprise' },
    };

    // Warm-up
    for (let i = 0; i < 1000; i++) {
      evaluateFlag(flag, ctx);
    }

    // Benchmark 10,000 iterations
    const iterations = 10000;
    const start = process.hrtime.bigint();
    for (let i = 0; i < iterations; i++) {
      evaluateFlag(flag, ctx);
    }
    const end = process.hrtime.bigint();

    const totalNs = Number(end - start);
    const avgNs = totalNs / iterations;
    const avgMs = avgNs / 1_000_000;

    // Must be < 1.0 ms
    expect(avgMs).toBeLessThan(1.0);
  });
});
