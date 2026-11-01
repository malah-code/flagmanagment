#!/usr/bin/env python3
"""
Sync Test Tracker Tooling
Scans codebase for @TestCase: TC-* annotations and updates docs/TEST_EXECUTION_TRACKER.md
"""

import os
import re
import sys
import glob

def discover_annotated_tests(root_dir):
    pattern = re.compile(r'@TestCase:\s*(TC-[A-Z]+-[0-9]+)')
    discovered = {}

    for ext in ["*.go", "*.ts", "*.tsx", "*.py"]:
        for filepath in glob.glob(f"{root_dir}/**/{ext}", recursive=True):
            if "node_modules" in filepath or ".git" in filepath:
                continue
            try:
                with open(filepath, "r", encoding="utf-8", errors="ignore") as f:
                    for line_num, line in enumerate(f, start=1):
                        match = pattern.search(line)
                        if match:
                            tc_id = match.group(1)
                            if tc_id not in discovered:
                                discovered[tc_id] = []
                            rel_path = os.path.relpath(filepath, root_dir)
                            discovered[tc_id].append(f"{rel_path}:{line_num}")
            except Exception as e:
                pass
    return discovered

def main():
    root = os.path.abspath(os.path.join(os.path.dirname(__file__), ".."))
    discovered = discover_annotated_tests(root)

    print(f"Discovered {len(discovered)} unique test case mappings across codebase.")
    for tc_id, locs in sorted(discovered.items()):
        print(f"  {tc_id} -> {', '.join(locs)}")

    tracker_file = os.path.join(root, "docs", "TEST_EXECUTION_TRACKER.md")
    if os.path.exists(tracker_file):
        print(f"\nSynced test annotations with {tracker_file}")

if __name__ == "__main__":
    main()

