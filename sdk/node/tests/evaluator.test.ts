import { evaluateFlag, hashPII } from '../src/evaluator';
import { FlagRule } from '../src/types';

describe('Evaluator Tests', () => {
  test('Evaluates disabled flag to default variation', () => {
    const flag: FlagRule = {
      key: 'test-flag',
      type: 'BOOLEAN',
      enabled: false,
      defaultVariation: 'false',
    };

    const res = evaluateFlag(flag);
    expect(res.value).toBe(false);
    expect(res.reason).toBe('DISABLED');
  });

  test('Evaluates enabled flag to default variation', () => {
    const flag: FlagRule = {
      key: 'test-flag',
      type: 'BOOLEAN',
      enabled: true,
      defaultVariation: 'true',
    };

    const res = evaluateFlag(flag);
    expect(res.value).toBe(true);
    expect(res.reason).toBe('DEFAULT');
  });

  test('Hashes PII attributes using SHA256', () => {
    const ctx = {
      identity: 'user-123',
      attributes: {
        email: 'test@example.com',
        age: 30,
      },
    };

    const hashed = hashPII(ctx);
    expect(hashed?.identity).not.toBe('user-123');
    expect(hashed?.attributes?.email).not.toBe('test@example.com');
    expect(hashed?.attributes?.age).toBe(30);
  });
});
