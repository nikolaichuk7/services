#!/usr/bin/env python3
"""Render MX| rows (variant|gsg|input|trust-anchor id|integrity|kds calls) as two Markdown tables."""
import sys, collections
rows=[l.rstrip("\n").split("|") for l in open(sys.argv[1]) if l.startswith("MX|")]
VAR={"main":"main (today)","flag":"main + `DisableCertFetching: false`","der":"#488","derflag":"#488 + `DisableCertFetching: false`","roots":"#488 + #490"}
GSG=lambda g: "pinned v0.14.2-pre" if g.startswith("v0.14.2") else g
inputs=[]; cells={}
for _,v,g,inp,ta,integ,kds in rows:
    if inp not in inputs: inputs.append(inp)
    cells[(v,g,inp)]=(ta,integ,kds)
cols=[(v,g) for v in ["main","flag","der","derflag","roots"] for g in sorted({r[2] for r in rows}, key=lambda x: (not x.startswith("v0.14.2"), x))]
def mark(c):
    ta,integ,kds=c
    ok = ta.startswith("ok") and integ=="ok"
    return ("✅" if ok else "❌") + (f" ({kds} KDS)" if kds!="0" else "")
print("| report, as the Verifier receives it | " + " | ".join(f"{VAR[v]}<br>{GSG(g)}" for v,g in cols) + " |")
print("|---|" + "|".join(":---:" for _ in cols) + "|")
for inp in inputs:
    print(f"| {inp} | " + " | ".join(mark(cells[(v,g,inp)]) for v,g in cols) + " |")
print()
print("Why each ❌ (the scheme's error, then the integrity check's):")
print()
print("| variant | go-sev-guest | report | GetTrustAnchorIDs | ValidateEvidenceIntegrity |")
print("|---|---|---|---|---|")
for v,g in cols:
    for inp in inputs:
        ta,integ,kds=cells[(v,g,inp)]
        if not (ta.startswith("ok") and integ=="ok"):
            print(f"| {VAR[v]} | {GSG(g)} | {inp} | {ta} | {integ} |")
