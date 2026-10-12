#!/bin/bash
# One test against five variants of this repository, each on the pinned go-sev-guest and on v0.15.0.
# Prints MX|variant|go-sev-guest|report|GetTrustAnchorIDs|ValidateEvidenceIntegrity|KDS requests, one line per cell.
set -u
HERE=$(cd "$(dirname "$0")" && pwd)
OUT=${1:-$HERE/results/matrix.txt}
WORK=${MX_WORK:-$(mktemp -d)}
FORK=https://github.com/nikolaichuk7/services.git
MAIN=b1cd67e3b2644d0aa52ccd93d6e02b356d7a13ab   # veraison/services main, the base of #488 and #490
DER=283df9a74e86f2b6126b96c321b65d19f6ff10b0    # the head of #488
ROOTS=d73ae831e16faaa5ffb30a64d73551efe90218b4  # the head of #490
mkdir -p "$(dirname "$OUT")"; : > "$OUT"
git fetch -q https://github.com/veraison/services.git "$MAIN"
git fetch -q "$FORK" sevsnp-der-cert-table sevsnp-amd-roots
for v in "main:$MAIN:" "flag:$MAIN:flag" "der:$DER:" "derflag:$DER:flag" "roots:$ROOTS:"; do
  IFS=: read -r name sha flag <<< "$v"
  d="$WORK/$name"
  git worktree add -q --detach "$d" "$sha" || { echo "MX|$name|-|worktree failed|||" >> "$OUT"; continue; }
  cp "$HERE/evidence/"*.cbor "$d/scheme/sevsnp/test/evidence/"
  cp "$HERE/matrix_test.go" "$d/scheme/sevsnp/matrix_test.go"
  if [ "$flag" = flag ]; then sed -i.bak 's/DisableCertFetching: true,/DisableCertFetching: false,/' "$d/scheme/sevsnp/scheme.go" && rm -f "$d/scheme/sevsnp/scheme.go.bak"; fi
  for gsg in pin v0.15.0; do
    ( cd "$d" || exit 1
      git checkout -q -- go.mod go.sum
      if [ "$gsg" = v0.15.0 ]; then go get github.com/google/go-sev-guest@v0.15.0 >/dev/null 2>&1 || { echo "MX|$name|v0.15.0|go get failed|||"; exit 0; }; fi
      ver=$(go list -m -f '{{.Version}}' github.com/google/go-sev-guest)
      MX_VARIANT=$name MX_GSG="$ver" go test ./scheme/sevsnp/ -run '^Test_Matrix$' -count=1 -v 2>&1 | grep '^MX|' || echo "MX|$name|$ver|the test did not run|||"
      git checkout -q -- go.mod go.sum
    ) >> "$OUT"
  done
  git worktree remove --force "$d"
done
echo "rows: $(grep -c '^MX|' "$OUT")"
