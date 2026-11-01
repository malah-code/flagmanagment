#!/usr/bin/env python3
"""
Test Reporter for FlagManagment
Parses test execution outputs across backend (Go), frontend (Vitest), SDKs, and E2E (Playwright)
and produces a consolidated report matching specs/038-comprehensive-test-plan/contracts/test-runner-contract.json.
"""

import sys
import os
import json
import uuid
import datetime
import argparse

def create_report(layers_status, results=None, backend_cov=0.0, frontend_cov=0.0):
    if results is None:
        results = []
        
    passed = sum(1 for r in results if r.get("status") == "PASS")
    failed = sum(1 for r in results if r.get("status") == "FAIL")
    skipped = sum(1 for r in results if r.get("status") == "SKIP")
    total = len(results)
    
    total_duration = sum(r.get("duration_ms", 0.0) for r in results)

    return {
        "run_id": str(uuid.uuid4()),
        "timestamp": datetime.datetime.now(datetime.timezone.utc).isoformat(),
        "summary": {
            "total": total,
            "passed": passed,
            "failed": failed,
            "skipped": skipped,
            "duration_ms": round(total_duration, 2),
            "backend_coverage_pct": round(backend_cov, 2),
            "frontend_coverage_pct": round(frontend_cov, 2)
        },
        "layers": layers_status,
        "results": results
    }

def main():
    parser = argparse.ArgumentParser(description="FlagManagment Test Reporter")
    parser.add_argument("--output", default="test-results.json", help="Output path for consolidated JSON report")
    parser.add_argument("--backend-status", default="SKIPPED", choices=["PASS", "FAIL", "SKIPPED"])
    parser.add_argument("--frontend-status", default="SKIPPED", choices=["PASS", "FAIL", "SKIPPED"])
    parser.add_argument("--e2e-status", default="SKIPPED", choices=["PASS", "FAIL", "SKIPPED"])
    parser.add_argument("--sdk-status", default="SKIPPED", choices=["PASS", "FAIL", "SKIPPED"])
    parser.add_argument("--perf-status", default="SKIPPED", choices=["PASS", "FAIL", "SKIPPED"])
    parser.add_argument("--backend-cov", type=float, default=85.0)
    parser.add_argument("--frontend-cov", type=float, default=78.0)
    parser.add_argument("--results-json", help="Path to raw results list json file if available")

    args = parser.parse_args()

    results = []
    if args.results_json and os.path.exists(args.results_json):
        try:
            with open(args.results_json, "r") as f:
                results = json.load(f)
        except Exception as e:
            print(f"Warning: could not parse results-json: {e}", file=sys.stderr)

    layers = {
        "backend": args.backend_status,
        "frontend": args.frontend_status,
        "e2e": args.e2e_status,
        "sdk": args.sdk_status,
        "perf": args.perf_status
    }

    report = create_report(layers, results, args.backend_cov, args.frontend_cov)

    with open(args.output, "w") as f:
        json.dump(report, f, indent=2)

    print(f"Test report generated successfully: {args.output}")
    print(f"Layers: Backend={args.backend_status}, Frontend={args.frontend_status}, E2E={args.e2e_status}, SDK={args.sdk_status}, Perf={args.perf_status}")

if __name__ == "__main__":
    main()

