| report, as the Verifier receives it | main (today)<br>pinned v0.14.2-pre | main (today)<br>v0.15.0 | main + `DisableCertFetching: false`<br>pinned v0.14.2-pre | main + `DisableCertFetching: false`<br>v0.15.0 | #488<br>pinned v0.14.2-pre | #488<br>v0.15.0 | #488 + `DisableCertFetching: false`<br>pinned v0.14.2-pre | #488 + `DisableCertFetching: false`<br>v0.15.0 | #488 + #490<br>pinned v0.14.2-pre | #488 + #490<br>v0.15.0 |
|---|:---:|:---:|:---:|:---:|:---:|:---:|:---:|:---:|:---:|:---:|
| upstream's own sample: Genoa VCEK, PEM table, ratsd token | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| GCP Milan VCEK, the host's table (VCEK+ASK+ARK, DER) | ❌ | ❌ | ❌ | ❌ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| GCP Milan VCEK, a host's table with the VCEK alone | ❌ | ❌ | ❌ | ❌ | ❌ | ❌ | ❌ | ❌ | ✅ | ✅ |
| AWS Milan VLEK, the host's table (VLEK alone) | ❌ | ❌ | ❌ | ❌ | ❌ | ❌ | ❌ | ❌ | ✅ | ✅ |
| AWS Milan VLEK, the guest appended ASVK+ARK | ❌ | ❌ | ❌ | ❌ | ✅ | ❌ | ✅ | ❌ | ✅ | ✅ |
| GCP VCEK alone, then AWS VLEK alone in the same process: the VCEK | ❌ | ❌ | ❌ | ❌ | ❌ | ❌ | ❌ | ❌ | ✅ | ✅ |
| GCP VCEK alone, then AWS VLEK alone in the same process: the VLEK | ❌ | ❌ | ❌ | ❌ | ❌ | ❌ | ❌ | ❌ | ✅ | ✅ |

Why each ❌ (the scheme's error, then the integrity check's):

| variant | go-sev-guest | report | GetTrustAnchorIDs | ValidateEvidenceIntegrity |
|---|---|---|---|---|
| main (today) | pinned v0.14.2-pre | GCP Milan VCEK, the host's table (VCEK+ASK+ARK, DER) | FAIL: can't read ARK to compose TA ID: failed to read certificate | FAIL: evidence Trust Anchor (ARK) doesn't match the provisioned one |
| main (today) | pinned v0.14.2-pre | GCP Milan VCEK, a host's table with the VCEK alone | FAIL: can't read ARK to compose TA ID: empty certificate | FAIL: missing ARK certificate in evidence |
| main (today) | pinned v0.14.2-pre | AWS Milan VLEK, the host's table (VLEK alone) | FAIL: can't read ARK to compose TA ID: empty certificate | FAIL: missing ARK certificate in evidence |
| main (today) | pinned v0.14.2-pre | AWS Milan VLEK, the guest appended ASVK+ARK | FAIL: can't read ARK to compose TA ID: failed to read certificate | FAIL: evidence Trust Anchor (ARK) doesn't match the provisioned one |
| main (today) | pinned v0.14.2-pre | GCP VCEK alone, then AWS VLEK alone in the same process: the VCEK | FAIL: can't read ARK to compose TA ID: empty certificate | FAIL: missing ARK certificate in evidence |
| main (today) | pinned v0.14.2-pre | GCP VCEK alone, then AWS VLEK alone in the same process: the VLEK | FAIL: can't read ARK to compose TA ID: empty certificate | FAIL: missing ARK certificate in evidence |
| main (today) | v0.15.0 | GCP Milan VCEK, the host's table (VCEK+ASK+ARK, DER) | FAIL: can't read ARK to compose TA ID: failed to read certificate | FAIL: evidence Trust Anchor (ARK) doesn't match the provisioned one |
| main (today) | v0.15.0 | GCP Milan VCEK, a host's table with the VCEK alone | FAIL: can't read ARK to compose TA ID: empty certificate | FAIL: missing ARK certificate in evidence |
| main (today) | v0.15.0 | AWS Milan VLEK, the host's table (VLEK alone) | FAIL: can't read ARK to compose TA ID: empty certificate | FAIL: missing ARK certificate in evidence |
| main (today) | v0.15.0 | AWS Milan VLEK, the guest appended ASVK+ARK | FAIL: can't read ARK to compose TA ID: failed to read certificate | FAIL: evidence Trust Anchor (ARK) doesn't match the provisioned one |
| main (today) | v0.15.0 | GCP VCEK alone, then AWS VLEK alone in the same process: the VCEK | FAIL: can't read ARK to compose TA ID: empty certificate | FAIL: missing ARK certificate in evidence |
| main (today) | v0.15.0 | GCP VCEK alone, then AWS VLEK alone in the same process: the VLEK | FAIL: can't read ARK to compose TA ID: empty certificate | FAIL: missing ARK certificate in evidence |
| main + `DisableCertFetching: false` | pinned v0.14.2-pre | GCP Milan VCEK, the host's table (VCEK+ASK+ARK, DER) | FAIL: can't read ARK to compose TA ID: failed to read certificate | FAIL: evidence Trust Anchor (ARK) doesn't match the provisioned one |
| main + `DisableCertFetching: false` | pinned v0.14.2-pre | GCP Milan VCEK, a host's table with the VCEK alone | FAIL: can't read ARK to compose TA ID: empty certificate | FAIL: missing ARK certificate in evidence |
| main + `DisableCertFetching: false` | pinned v0.14.2-pre | AWS Milan VLEK, the host's table (VLEK alone) | FAIL: can't read ARK to compose TA ID: empty certificate | FAIL: missing ARK certificate in evidence |
| main + `DisableCertFetching: false` | pinned v0.14.2-pre | AWS Milan VLEK, the guest appended ASVK+ARK | FAIL: can't read ARK to compose TA ID: failed to read certificate | FAIL: evidence Trust Anchor (ARK) doesn't match the provisioned one |
| main + `DisableCertFetching: false` | pinned v0.14.2-pre | GCP VCEK alone, then AWS VLEK alone in the same process: the VCEK | FAIL: can't read ARK to compose TA ID: empty certificate | FAIL: missing ARK certificate in evidence |
| main + `DisableCertFetching: false` | pinned v0.14.2-pre | GCP VCEK alone, then AWS VLEK alone in the same process: the VLEK | FAIL: can't read ARK to compose TA ID: empty certificate | FAIL: missing ARK certificate in evidence |
| main + `DisableCertFetching: false` | v0.15.0 | GCP Milan VCEK, the host's table (VCEK+ASK+ARK, DER) | FAIL: can't read ARK to compose TA ID: failed to read certificate | FAIL: evidence Trust Anchor (ARK) doesn't match the provisioned one |
| main + `DisableCertFetching: false` | v0.15.0 | GCP Milan VCEK, a host's table with the VCEK alone | FAIL: can't read ARK to compose TA ID: empty certificate | FAIL: missing ARK certificate in evidence |
| main + `DisableCertFetching: false` | v0.15.0 | AWS Milan VLEK, the host's table (VLEK alone) | FAIL: can't read ARK to compose TA ID: empty certificate | FAIL: missing ARK certificate in evidence |
| main + `DisableCertFetching: false` | v0.15.0 | AWS Milan VLEK, the guest appended ASVK+ARK | FAIL: can't read ARK to compose TA ID: failed to read certificate | FAIL: evidence Trust Anchor (ARK) doesn't match the provisioned one |
| main + `DisableCertFetching: false` | v0.15.0 | GCP VCEK alone, then AWS VLEK alone in the same process: the VCEK | FAIL: can't read ARK to compose TA ID: empty certificate | FAIL: missing ARK certificate in evidence |
| main + `DisableCertFetching: false` | v0.15.0 | GCP VCEK alone, then AWS VLEK alone in the same process: the VLEK | FAIL: can't read ARK to compose TA ID: empty certificate | FAIL: missing ARK certificate in evidence |
| #488 | pinned v0.14.2-pre | GCP Milan VCEK, a host's table with the VCEK alone | FAIL: can't read ARK to compose TA ID: empty certificate | FAIL: missing ARK certificate in evidence |
| #488 | pinned v0.14.2-pre | AWS Milan VLEK, the host's table (VLEK alone) | FAIL: can't read ARK to compose TA ID: empty certificate | FAIL: missing ARK certificate in evidence |
| #488 | pinned v0.14.2-pre | GCP VCEK alone, then AWS VLEK alone in the same process: the VCEK | FAIL: can't read ARK to compose TA ID: empty certificate | FAIL: missing ARK certificate in evidence |
| #488 | pinned v0.14.2-pre | GCP VCEK alone, then AWS VLEK alone in the same process: the VLEK | FAIL: can't read ARK to compose TA ID: empty certificate | FAIL: missing ARK certificate in evidence |
| #488 | v0.15.0 | GCP Milan VCEK, a host's table with the VCEK alone | FAIL: can't read ARK to compose TA ID: empty certificate | FAIL: missing ARK certificate in evidence |
| #488 | v0.15.0 | AWS Milan VLEK, the host's table (VLEK alone) | FAIL: can't read ARK to compose TA ID: empty certificate | FAIL: missing ARK certificate in evidence |
| #488 | v0.15.0 | AWS Milan VLEK, the guest appended ASVK+ARK | ok: ARK-Milan | FAIL: VLEK could not be verified by any trusted roots. Last error: root of trust missing interme… |
| #488 | v0.15.0 | GCP VCEK alone, then AWS VLEK alone in the same process: the VCEK | FAIL: can't read ARK to compose TA ID: empty certificate | FAIL: missing ARK certificate in evidence |
| #488 | v0.15.0 | GCP VCEK alone, then AWS VLEK alone in the same process: the VLEK | FAIL: can't read ARK to compose TA ID: empty certificate | FAIL: missing ARK certificate in evidence |
| #488 + `DisableCertFetching: false` | pinned v0.14.2-pre | GCP Milan VCEK, a host's table with the VCEK alone | FAIL: can't read ARK to compose TA ID: empty certificate | FAIL: missing ARK certificate in evidence |
| #488 + `DisableCertFetching: false` | pinned v0.14.2-pre | AWS Milan VLEK, the host's table (VLEK alone) | FAIL: can't read ARK to compose TA ID: empty certificate | FAIL: missing ARK certificate in evidence |
| #488 + `DisableCertFetching: false` | pinned v0.14.2-pre | GCP VCEK alone, then AWS VLEK alone in the same process: the VCEK | FAIL: can't read ARK to compose TA ID: empty certificate | FAIL: missing ARK certificate in evidence |
| #488 + `DisableCertFetching: false` | pinned v0.14.2-pre | GCP VCEK alone, then AWS VLEK alone in the same process: the VLEK | FAIL: can't read ARK to compose TA ID: empty certificate | FAIL: missing ARK certificate in evidence |
| #488 + `DisableCertFetching: false` | v0.15.0 | GCP Milan VCEK, a host's table with the VCEK alone | FAIL: can't read ARK to compose TA ID: empty certificate | FAIL: missing ARK certificate in evidence |
| #488 + `DisableCertFetching: false` | v0.15.0 | AWS Milan VLEK, the host's table (VLEK alone) | FAIL: can't read ARK to compose TA ID: empty certificate | FAIL: missing ARK certificate in evidence |
| #488 + `DisableCertFetching: false` | v0.15.0 | AWS Milan VLEK, the guest appended ASVK+ARK | ok: ARK-Milan | FAIL: VLEK could not be verified by any trusted roots. Last error: root of trust missing interme… |
| #488 + `DisableCertFetching: false` | v0.15.0 | GCP VCEK alone, then AWS VLEK alone in the same process: the VCEK | FAIL: can't read ARK to compose TA ID: empty certificate | FAIL: missing ARK certificate in evidence |
| #488 + `DisableCertFetching: false` | v0.15.0 | GCP VCEK alone, then AWS VLEK alone in the same process: the VLEK | FAIL: can't read ARK to compose TA ID: empty certificate | FAIL: missing ARK certificate in evidence |
