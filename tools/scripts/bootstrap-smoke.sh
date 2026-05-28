#!/usr/bin/env bash
# Greenfield bootstrap smoke test for canonical ai-sdlc (maintainer/CI only).
# Materializes templates/consumer into a temp product repo and runs make build + validate EP-000.
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
AI_SDLC_ROOT="${AI_SDLC_ROOT:-$(cd "${SCRIPT_DIR}/../.." && pwd)}"
TEMPLATES="${AI_SDLC_ROOT}/templates/consumer"

TMP=""
PRODUCT=""

cleanup() {
  if [[ -n "${TMP}" && "${KEEP_TMP:-}" != "1" ]]; then
    rm -rf "${TMP}"
  fi
}
trap cleanup EXIT

die() {
  echo "bootstrap-smoke: $*" >&2
  exit 1
}

# First non-comment, non-empty line from template ai-sdlc.version (example tag).
template_pin_tag() {
  local line
  while IFS= read -r line || [[ -n "${line}" ]]; do
    line="${line%%#*}"
    line="$(echo "${line}" | tr -d '[:space:]')"
    [[ -n "${line}" ]] || continue
    echo "${line}"
    return 0
  done < "${TEMPLATES}/ai-sdlc.version"
  return 1
}

resolve_pin() {
  if [[ -n "${AI_SDLC_PIN:-}" ]]; then
    echo "${AI_SDLC_PIN}" | tr -d '[:space:]'
    return 0
  fi

  local head tag tag_commit
  head="$(git -C "${AI_SDLC_ROOT}" rev-parse HEAD)"
  if tag="$(template_pin_tag)"; then
    if tag_commit="$(git -C "${AI_SDLC_ROOT}" rev-parse "${tag}^{commit}" 2>/dev/null)"; then
      if [[ "${head}" == "${tag_commit}" ]]; then
        echo "${tag}"
        return 0
      fi
    fi
  fi
  echo "${head}"
}

pin_matches_checkout() {
  local pin="$1"
  local head tag_commit
  head="$(git -C "${PRODUCT}/ai-sdlc" rev-parse HEAD)"
  if [[ "${pin}" == "${head}" ]]; then
    return 0
  fi
  if tag_commit="$(git -C "${PRODUCT}/ai-sdlc" rev-parse "${pin}^{commit}" 2>/dev/null)"; then
    [[ "${tag_commit}" == "${head}" ]]
    return
  fi
  return 1
}

assert_file() {
  [[ -e "${PRODUCT}/$1" ]] || die "missing expected file: $1"
}

main() {
  [[ -f "${AI_SDLC_ROOT}/specification/pipeline.spec.md" ]] \
    || die "missing ${AI_SDLC_ROOT}/specification/pipeline.spec.md (set AI_SDLC_ROOT?)"
  [[ -d "${TEMPLATES}" ]] || die "missing ${TEMPLATES}"

  TMP="$(mktemp -d)"
  PRODUCT="${TMP}/my-app"
  mkdir -p "${PRODUCT}"

  git clone --local "${AI_SDLC_ROOT}" "${PRODUCT}/ai-sdlc" >/dev/null 2>&1 \
    || die "git clone --local failed for ${AI_SDLC_ROOT}"

  local pin
  pin="$(resolve_pin)"
  printf '%s\n' "${pin}" > "${PRODUCT}/ai-sdlc.version"

  cp "${TEMPLATES}/AGENTS.md" "${PRODUCT}/AGENTS.md"
  cp "${TEMPLATES}/.gitignore" "${PRODUCT}/.gitignore"
  cp "${TEMPLATES}/Makefile" "${PRODUCT}/Makefile"
  cp "${TEMPLATES}/README.project.md" "${PRODUCT}/README.md"
  mkdir -p "${PRODUCT}/.github/workflows"
  cp "${TEMPLATES}/.github/workflows/ai-sdlc.yml" "${PRODUCT}/.github/workflows/ai-sdlc.yml"
  cp -R "${TEMPLATES}/ai-sdlc-artefacts" "${PRODUCT}/ai-sdlc-artefacts"

  (
    cd "${PRODUCT}"
    make build
    [[ -x bin/validate ]] || die "bin/validate not executable after make build"
    ./bin/validate pipeline EP-000
    ./bin/validate structure EP-000
  )

  assert_file "AGENTS.md"
  assert_file ".gitignore"
  assert_file "Makefile"
  assert_file ".github/workflows/ai-sdlc.yml"
  assert_file "ai-sdlc-artefacts/README.md"
  assert_file "ai-sdlc-artefacts/scope.md"
  assert_file "ai-sdlc-artefacts/strategy.md"
  assert_file "ai-sdlc-artefacts/epics/EP-000/ep-scope.md"

  grep -q '^ai-sdlc/$' "${PRODUCT}/.gitignore" \
    || die ".gitignore missing layout A entry ai-sdlc/"

  pin_matches_checkout "${pin}" || die "ai-sdlc.version pin does not match ai-sdlc/ checkout"

  echo "bootstrap-smoke: OK (pin=${pin})"
  if [[ "${KEEP_TMP:-}" == "1" ]]; then
    echo "bootstrap-smoke: temp product root: ${PRODUCT}"
    TMP=""
  fi
}

main "$@"
