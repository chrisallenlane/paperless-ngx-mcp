#!/usr/bin/env bash
set -euo pipefail

# Duration per fuzz target (default: 30s).
FUZZTIME="${FUZZTIME:-30s}"

# Each entry is "package:FuzzFunctionName".
targets=(
  "./internal/tools/:FuzzValidateFilePath"
  "./internal/tools/:FuzzParseIDArg"
  "./internal/tools/:FuzzParsePatchArgs"
  "./internal/tools/:FuzzBuildListPath"
  "./internal/tools/:FuzzBuildDocumentListPath"
  "./internal/tools/:FuzzBuildTaskListPath"
  "./internal/tools/:FuzzValidateMatchableCreate"
  "./internal/tools/:FuzzValidateCreateTag"
  "./internal/tools/:FuzzValidateCreateStoragePath"
  "./internal/tools/:FuzzValidateCreateCustomField"
  "./internal/tools/:FuzzValidateCreateSavedView"
  "./internal/tools/:FuzzFormatDate"
  "./internal/tools/:FuzzFormatFileSize"
  "./internal/tools/:FuzzFormatStatistics"
)

passed=0
failed=0

for entry in "${targets[@]}"; do
  pkg="${entry%%:*}"
  func="${entry##*:}"

  echo "=== ${func} (${pkg}, ${FUZZTIME}) ==="
  if go test -fuzz="${func}" -fuzztime="${FUZZTIME}" "${pkg}"; then
    passed=$((passed + 1))
  else
    failed=$((failed + 1))
  fi
  echo ""
done

echo "=== Summary: ${passed} passed, ${failed} failed ==="
exit "${failed}"
