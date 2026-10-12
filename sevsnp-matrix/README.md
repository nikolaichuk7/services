# sevsnp-matrix

One test, `matrix_test.go`, run against five variants of veraison/services, each on the pinned go-sev-guest
(v0.14.2-0.20251119154202-af1c107a648f) and on v0.15.0, with upstream's own sample and two real reports:

| variant | code |
|---|---|
| main | veraison/services main b1cd67e, the base of #488 and #490 |
| main + flag | the same with `DisableCertFetching: false` in `scheme/sevsnp/scheme.go` |
| #488 | 283df9a, DER certificates in the host certificate table |
| #488 + flag | the same with `DisableCertFetching: false` |
| #488 + #490 | d73ae83, the chain verified up to AMD's roots held by the scheme |

Reports (`evidence/`): an AWS EC2 c6a.large (Milan) report signed with a VLEK, 6 October 2026, with the
table as the host supplies it (the VLEK alone) and as the guest appended AMD's ASVK and ARK; a Google Cloud
(Milan) report signed with a VCEK, 7 October 2026, with the host's table (VCEK, ASK, ARK, DER) and cut down
to the VCEK alone; and the two leaf-alone reports through one process in a row. Each cell runs
`GetTrustAnchorIDs` then `ValidateEvidenceIntegrity` with AMD's ARK of the product line as the provisioned
trust anchor, and counts the requests to `kdsintf.amd.com` at `http.DefaultTransport`.

Run `./matrix.sh` (needs Go and network for modules), then `./render.py results/matrix.txt`. The workflow
`.github/workflows/sevsnp-matrix.yml` does the same on a GitHub runner and writes the tables into the job
summary. `results/` holds the rows and tables of the run of 12 October 2026.
