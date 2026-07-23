import { FlagRule, EvaluationContext, EvaluationResult } from './types';
import * as crypto from 'crypto';

import * as murmur from 'murmurhash3js-revisited';

export function hashPII(ctx?: EvaluationContext): EvaluationContext | undefined {
  if (!ctx) return undefined;

  const hashedAttrs: Record<string, any> = {};
  if (ctx.attributes) {
    for (const [k, v] of Object.entries(ctx.attributes)) {
      if (typeof v === 'string') {
        hashedAttrs[k] = crypto.createHash('sha256').update(v).digest('hex');
      } else {
        hashedAttrs[k] = v;
      }
    }
  }

  return {
    identity: ctx.identity ? crypto.createHash('sha256').update(ctx.identity).digest('hex') : undefined,
    attributes: hashedAttrs,
  };
}

export function evaluateFlag(flag: FlagRule, ctx?: EvaluationContext): EvaluationResult {
  if (!flag.enabled) {
    return {
      value: parseValue(flag.defaultVariation, flag.type),
      reason: 'DISABLED',
    };
  }

  // Placeholder for Targeting Rules (Rollouts / Segments)
  // If we had targeting rules, we would use murmur.x86.hash32(`${flag.key}${ctx?.identity || ''}`) % 10000

  return {
    value: parseValue(flag.defaultVariation, flag.type),
    reason: 'DEFAULT',
  };
}

function parseValue(val: string, type: string): any {
  if (type === 'BOOLEAN') {
    return val === 'true';
  }
  if (type === 'NUMBER') {
    const num = Number(val);
    return isNaN(num) ? val : num;
  }
  if (type === 'JSON') {
    try {
      return JSON.parse(val);
    } catch {
      return val;
    }
  }
  return val;
}
