#!/usr/bin/env bash
# Module boundary checks — stub until cmd/ and internal/ exist (EP-001+).
# Replace this stub with cycle detection and forbidden-edge rules; set MODULE from: go list -m

set -e

echo "check-module-boundaries: skipped (no product packages yet; enable when cmd/ and internal/ exist)"
exit 0

# --- Example implementation (uncomment and remove stub above when ready) ---
#
# ROOT="$(cd "$(dirname "$0")/.." && pwd)"
# cd "$ROOT"
# MODULE="$(go list -m)"
# packages=$(go list "./cmd/..." "./internal/..." 2>/dev/null | sort -u)
# if [ -z "$packages" ]; then
#   echo "error: no packages under cmd/ or internal/" >&2
#   exit 1
# fi
# # Add cycle detection and forbidden import edges for your module layout.
